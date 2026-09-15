package service

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"

	"ps-sdk/model"
)

// Client 面向业务的平台客户端：持有适配器与数据库句柄，负责把平台数据落库。
// 同一平台只应创建一个，适配器内部维护 Token，重复创建会触发多余登录。
type Client struct {
	adapter Adapter
}

// Open 按平台短码创建客户端，凭据优先取环境变量，其次取 platform_auth 表
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
	return NewClient(adapter), nil
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
	return NewClient(adapter), nil
}

// NewClient 用现成的适配器创建客户端，便于测试与自定义实现
func NewClient(adapter Adapter) *Client {
	return &Client{adapter: adapter}
}

// Adapter 返回底层适配器
func (c *Client) Adapter() Adapter { return c.adapter }

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

	for i := range stations {
		saved, err := UpsertStation(&stations[i])
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("电站 %s(%s) 落库失败: %w",
				stations[i].StationName, stations[i].StationIDOrigin, err))
			continue
		}
		result.SavedStations++

		devices, err := listDevices(c.adapter, saved.StationID, result, stations[i].StationName)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("电站 %s 设备拉取失败: %w",
				stations[i].StationName, err))
			continue
		}

		for j := range devices {
			device := &devices[j]
			device.StationID = saved.StationID
			if device.DeviceName == "" {
				device.DeviceName = device.DeviceSN
			}
			if err := UpsertDevice(device); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("设备 %s 落库失败: %w",
					device.DeviceIDOrigin, err))
				continue
			}
			result.SavedDevices++
		}
		result.Devices += len(devices)
		result.DeviceCounts[saved.StationID] = len(devices)
	}

	Logf("[%s] 同步完成：电站 %d/%d，设备 %d/%d，失败 %d 项",
		platform.Code, result.SavedStations, result.Stations,
		result.SavedDevices, result.Devices, len(result.Errors))
	return result, nil
}

// listDevices 拉取设备列表，把「部分设备族不可用」降级为告警并返回已拿到的设备
func listDevices(adapter Adapter, stationID uint64, result *SyncResult, stationName string) ([]model.PowerDevice, error) {
	devices, err := adapter.ListPowerDevices(stationID)
	if err == nil {
		return devices, nil
	}

	var partial *PartialError
	if errors.As(err, &partial) {
		result.Warnings = append(result.Warnings, fmt.Errorf("电站 %s: %w", stationName, partial))
		return partial.Devices, nil
	}
	return nil, err
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

	for i := range devices {
		devices[i].StationID = stationID
		if devices[i].DeviceName == "" {
			devices[i].DeviceName = devices[i].DeviceSN
		}
		if err := UpsertDevice(&devices[i]); err != nil {
			return nil, fmt.Errorf("设备 %s 落库失败: %w", devices[i].DeviceIDOrigin, err)
		}
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
				log.Printf("[%s] 跳过：%v", platform.Code, err)
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
