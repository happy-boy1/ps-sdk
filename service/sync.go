package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"

	"ps-sdk/model"
)

// Client 面向业务的平台客户端：持有适配器与数据库句柄，负责把平台数据落库。
// 同一平台只应创建一个，适配器内部维护 Token，重复创建会触发多余登录。
type Client struct {
	adapter Adapter
	opts    Options
}

// Open 按平台短码创建客户端。
// 凭据来源优先级由低到高：config.toml → platform_auth 表 → 环境变量。
func Open(code string, opts Options) (*Client, error) {
	platform, ok := PlatformByCode(code)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrPlatformUnknown, code)
	}

	cred, err := CredentialFromDB(platform.ID)
	if err != nil {
		return nil, err
	}
	adapter, err := NewAdapter(cred, opts)
	if err != nil {
		return nil, err
	}
	return NewClient(adapter, opts), nil
}

// OpenWith 使用调用方提供的凭据创建客户端，不查库
func OpenWith(cred Credential, opts Options) (*Client, error) {
	if cred.Code == "" {
		if platform, ok := PlatformByID(cred.PlatformID); ok {
			cred.Code = platform.Code
		}
	}
	adapter, err := NewAdapter(cred, opts)
	if err != nil {
		return nil, err
	}
	return NewClient(adapter, opts), nil
}

// NewClient 用现成的适配器创建客户端，便于测试与自定义实现
func NewClient(adapter Adapter, opts Options) *Client {
	return &Client{adapter: adapter, opts: opts.WithDefaults()}
}

// Adapter 返回底层适配器
func (c *Client) Adapter() Adapter { return c.adapter }

// Options 返回客户端使用的运行参数
func (c *Client) Options() Options { return c.opts }

// Platform 返回平台元信息
func (c *Client) Platform() Platform {
	if c.adapter == nil {
		return Platform{}
	}
	return c.adapter.Platform()
}

// Sync 从平台拉取全部电站及设备并写入数据库。
// 单个电站的设备拉取失败不影响其余电站，失败信息收集在 SyncResult.Errors 中。
func (c *Client) Sync() (*SyncResult, error) {
	if c.adapter == nil {
		return nil, errors.New("客户端未绑定平台适配器")
	}
	if db == nil {
		return nil, ErrDatabaseNotReady
	}

	platform := c.adapter.Platform()
	result := newSyncResult(platform)

	stations, err := c.adapter.ListPowerStations()
	if err != nil {
		return result, err
	}
	result.Stations = len(stations)

	saved, err := c.saveStations(stations, result)
	if err != nil {
		return result, err
	}
	c.syncDevices(saved, result)

	Logf("[%s] 同步完成：电站 %d/%d，设备 %d/%d，失败 %d 项，告警 %d 项",
		platform.Code, result.SavedStations, result.Stations,
		result.SavedDevices, result.Devices, len(result.Errors), len(result.Warnings))
	return result, nil
}

// saveStations 逐条落库电站，返回落库后的记录
func (c *Client) saveStations(stations []model.PowerStation, result *SyncResult) ([]model.PowerStation, error) {
	saved := make([]model.PowerStation, 0, len(stations))
	for i := range stations {
		station, err := UpsertStation(&stations[i])
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("电站 %s(%s) 落库失败: %w",
				stations[i].StationName, stations[i].StationIDOrigin, err))
			continue
		}
		saved = append(saved, *station)
		result.SavedStations++
	}
	return saved, nil
}

// syncDevices 拉取并落库各电站的设备。
// Concurrency > 1 时并发拉取（网络在平台侧耗时），落库仍串行执行。
func (c *Client) syncDevices(stations []model.PowerStation, result *SyncResult) {
	if len(stations) == 0 {
		return
	}

	type pulled struct {
		devices    []model.PowerDevice
		warnings   []error
		fetchError error
		station    model.PowerStation
	}

	results := make([]pulled, len(stations))
	workers := c.opts.Concurrency
	if workers < 1 {
		workers = 1
	}

	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	for i := range stations {
		wg.Add(1)
		sem <- struct{}{}
		go func(index int) {
			defer wg.Done()
			defer func() { <-sem }()

			item := pulled{station: stations[index]}
			devices, err := c.adapter.ListPowerDevices(stations[index].StationID)
			switch {
			case err == nil:
				item.devices = devices
			default:
				var partial *PartialError
				if errors.As(err, &partial) {
					item.devices = partial.Devices
					item.warnings = partial.Warnings
				} else {
					item.fetchError = err
				}
			}
			results[index] = item
		}(i)
	}
	wg.Wait()

	for _, item := range results {
		name := item.station.StationName
		for _, warning := range item.warnings {
			result.Warnings = append(result.Warnings, fmt.Errorf("电站 %s: %w", name, warning))
		}
		if item.fetchError != nil {
			result.Errors = append(result.Errors, fmt.Errorf("电站 %s 设备拉取失败: %w", name, item.fetchError))
			continue
		}
		c.saveDevices(item.station.StationID, item.devices, result)
	}
}

// saveDevices 落库单站设备，同一原始 ID 只保留最后一条，与库内 UPSERT 语义一致
func (c *Client) saveDevices(stationID uint64, devices []model.PowerDevice, result *SyncResult) {
	lastIndex := make(map[string]int, len(devices))
	for i := range devices {
		key := devices[i].DeviceIDOrigin
		if prev, ok := lastIndex[key]; ok {
			devices[prev].DeviceIDOrigin = ""
		}
		lastIndex[key] = i
	}

	saved := 0
	for i := range devices {
		device := &devices[i]
		if device.DeviceIDOrigin == "" {
			continue
		}
		device.StationID = stationID
		if device.DeviceName == "" {
			device.DeviceName = device.DeviceSN
		}
		if err := UpsertDevice(device); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("设备 %s 落库失败: %w", device.DeviceIDOrigin, err))
			continue
		}
		saved++
		result.SavedDevices++
	}
	result.Devices += saved
	result.DeviceCounts[stationID] = saved
}

// Stations 读取本平台已落库的电站
func (c *Client) Stations() ([]model.PowerStation, error) {
	return ListStations(c.Platform().ID)
}

// Devices 读取本平台已落库的设备，stationID 为 0 时返回全部
func (c *Client) Devices(stationID uint64) ([]model.PowerDevice, error) {
	return ListDevices(c.Platform().ID, stationID)
}

// SyncStations 只同步电站（含设备列表所需的主键回填），不拉取设备
func (c *Client) SyncStations() ([]model.PowerStation, error) {
	if c.adapter == nil {
		return nil, errors.New("客户端未绑定平台适配器")
	}
	if db == nil {
		return nil, ErrDatabaseNotReady
	}

	stations, err := c.adapter.ListPowerStations()
	if err != nil {
		return nil, err
	}

	saved := make([]model.PowerStation, 0, len(stations))
	for i := range stations {
		station, err := UpsertStation(&stations[i])
		if err != nil {
			return saved, fmt.Errorf("电站 %s(%s) 落库失败: %w",
				stations[i].StationName, stations[i].StationIDOrigin, err)
		}
		saved = append(saved, *station)
	}
	return saved, nil
}

// SyncDevices 拉取并落库指定电站的设备
func (c *Client) SyncDevices(stationID uint64) ([]model.PowerDevice, error) {
	if c.adapter == nil {
		return nil, errors.New("客户端未绑定平台适配器")
	}
	if db == nil {
		return nil, ErrDatabaseNotReady
	}
	if stationID == 0 {
		return nil, errors.New("stationID 不能为 0")
	}

	devices, err := c.adapter.ListPowerDevices(stationID)
	if err != nil {
		return nil, err
	}

	result := newSyncResult(c.Platform())
	c.saveDevices(stationID, devices, result)
	if len(result.Errors) > 0 {
		return devices, errors.Join(result.Errors...)
	}
	return devices, nil
}

// SyncAll 同步所有已实现适配器的平台，单个平台失败不影响其余平台
func SyncAll(opts Options) ([]*SyncResult, error) {
	results := make([]*SyncResult, 0, len(Platforms))
	errs := make([]error, 0)

	for _, platform := range Platforms {
		client, err := Open(platform.Code, opts)
		if err != nil {
			if errors.Is(err, ErrCredentialMissing) || errors.Is(err, ErrDatabaseNotReady) {
				Logf("[%s] 跳过：%v", platform.Code, err)
				continue
			}
			errs = append(errs, fmt.Errorf("[%s] 打开失败: %w", platform.Code, err))
			continue
		}

		result, err := client.Sync()
		if result != nil {
			results = append(results, result)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return results, errors.Join(errs...)
	}
	return results, nil
}

// Close 释放数据库连接
func Close() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	db = nil
	return sqlDB.Close()
}

// EnsureSchema 自动迁移服务层依赖的表结构，只补字段不建外键。
// 外键约束由 DDL 维护：交给 GORM 生成会与其他表互相引用，进而迁移失败。
func EnsureSchema(database *gorm.DB) error {
	if database == nil {
		return ErrDatabaseNotReady
	}
	return database.Session(&gorm.Session{}).AutoMigrate(
		&model.PlatformInfo{},
		&model.PlatformAuth{},
		&model.DeviceType{},
		&model.PowerStation{},
		&model.PowerDevice{},
	)
}

// NormalizeCode 归一化平台短码，便于调用方直接传入用户输入
func NormalizeCode(code string) string { return strings.TrimSpace(normalizeCode(code)) }
