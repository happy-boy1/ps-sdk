package service

import (
	"log"
	"strings"

	"ps-sdk/model"
	"ps-sdk/sdk/solarman"
)

// 适配器分页参数
const (
	stationPageSize = 100 // 电站列表每页条数
	devicePageSize  = 100 // 电站设备列表每页条数
)

// solarmanAdapter SolarMan（小麦智电 / 小麦商家版）适配器
type solarmanAdapter struct {
	limitBackoff
	sdk *solarman.SolarmanSDK
}

func init() {
	registerAdapter(CodeSolarman, newSolarmanAdapter)
}

func newSolarmanAdapter(cred Credential, opts Options) (Adapter, error) {
	appSecret := strings.TrimSpace(cred.AppSecret)
	password := strings.TrimSpace(cred.Password)
	if strings.TrimSpace(cred.AppID) == "" || appSecret == "" {
		return nil, wrapErr(cred.PlatformID, CodeSolarman, "构造适配器", ErrCredentialMissing)
	}
	if password == "" {
		return nil, wrapErr(cred.PlatformID, CodeSolarman, "构造适配器", ErrCredentialMissing)
	}

	// 登录身份：邮箱按 email 下发，纯数字按手机号下发（需国家码），其余按用户名
	email := strings.TrimSpace(cred.Email)
	mobile := strings.TrimSpace(cred.Mobile)
	userName := strings.TrimSpace(cred.UserName)
	if identity := cred.Identity(); identity != "" {
		switch {
		case strings.Contains(identity, "@"):
			email = identity
		case isDigits(identity):
			mobile = identity
		default:
			userName = identity
		}
	}

	sdk, err := solarman.NewSolarmanSDK(solarman.Credentials{
		AppID:       strings.TrimSpace(cred.AppID),
		AppSecret:   appSecret,
		Email:       email,
		Mobile:      mobile,
		CountryCode: strings.TrimSpace(cred.CountryCode),
		UserName:    userName,
		Password:    password,
		OrgID:       cred.OrgID,
		BaseURL:     strings.TrimSpace(cred.APIURL),
	}, solarmanOptions(opts)...)
	if err != nil {
		return nil, wrapErr(cred.PlatformID, CodeSolarman, "构造适配器", err)
	}

	return &solarmanAdapter{sdk: sdk}, nil
}

func (a *solarmanAdapter) Platform() Platform {
	p, _ := PlatformByCode(CodeSolarman)
	return p
}

func (a *solarmanAdapter) ListPowerStations() ([]model.PowerStation, error) {
	stations := make([]model.PowerStation, 0, stationPageSize)
	for page := 1; ; page++ {
		var res *solarman.StationListResult
		err := a.throttle(isRateLimitError, func() error {
			result, err := a.sdk.StationList(solarman.StationListRequest{
				PageRequest: solarman.PageRequest{Page: page, Size: stationPageSize},
			})
			res = result
			return err
		})
		if err != nil {
			return nil, wrapErr(PlatformSolarman, CodeSolarman, "列出电站", err)
		}
		for i := range res.StationList {
			stations = append(stations, *solarmanStation(&res.StationList[i]))
		}
		if len(res.StationList) == 0 || len(stations) >= res.Total {
			return stations, nil
		}
	}
}

func (a *solarmanAdapter) ListPowerDevices(stationID uint64) ([]model.PowerDevice, error) {
	station, err := a.stationOrigin(stationID)
	if err != nil {
		return nil, err
	}

	devices := make([]model.PowerDevice, 0, devicePageSize)
	for page := 1; ; page++ {
		var res *solarman.StationDeviceListResult
		err := a.throttleRetry(isRateLimitError, deviceListAttempts, func() error {
			result, err := a.sdk.StationDeviceList(solarman.StationDeviceListRequest{
				PageRequest: solarman.PageRequest{Page: page, Size: devicePageSize},
				StationID:   station,
			})
			res = result
			return err
		})
		if err != nil {
			return nil, wrapErr(PlatformSolarman, CodeSolarman, "列出设备", err)
		}
		for i := range res.DeviceListItems {
			devices = append(devices, *solarmanDevice(&res.DeviceListItems[i], stationID))
		}
		if len(res.DeviceListItems) == 0 || len(devices) >= res.Total {
			return devices, nil
		}
	}
}

// stationOrigin 把内部电站 ID 还原成平台原始电站 ID
func (a *solarmanAdapter) stationOrigin(stationID uint64) (solarman.Int64, error) {
	ps, err := findStation(PlatformSolarman, stationID)
	if err != nil {
		return 0, wrapErr(PlatformSolarman, CodeSolarman, "读取电站", err)
	}
	return solarman.Int64(parseInt64(ps.StationIDOrigin)), nil
}

// solarmanStation 转换为统一电站模型
func solarmanStation(item *solarman.StationListItem) *model.PowerStation {
	// 电站列表不返回运行状态，能拉取到即视为运行中
	return &model.PowerStation{
		PlatformID:      PlatformSolarman,
		StationIDOrigin: item.ID.String(),
		StationName:     item.Name,
		CapacityKwp:     float64Ptr(item.InstalledCapacity.Float()),
		Address:         item.LocationAddress,
		Longitude:       item.LocationLng.Float(),
		Latitude:        item.LocationLat.Float(),
		GridConnectedAt: unixSeconds(int64(item.StartOperatingTime)),
		Status:          StationRunning,
	}
}

// solarmanDevice 转换为统一设备模型
func solarmanDevice(item *solarman.StationDeviceItem, stationID uint64) *model.PowerDevice {
	origin := firstNonEmpty(formatInt64(item.DeviceID.Int()), item.DeviceSN)
	deviceType := resolveDeviceType(CodeSolarman, string(item.DeviceType))

	device := newPowerDevice(PlatformSolarman, stationID, deviceType.ID,
		origin, item.DeviceSN, "", "", "")
	device.DeviceTypeOrigin = string(item.DeviceType)
	device.Status = solarmanStatus(item.ConnectStatus)
	return device
}

// solarmanStatus 设备通讯状态：0 离线 1 在线 2 报警
func solarmanStatus(status int) int8 {
	switch status {
	case 0:
		return DeviceOffline
	case 1:
		return DeviceOnline
	case 2:
		return DeviceAlarm
	default:
		return DeviceUnknown
	}
}

func solarmanOptions(opts Options) []solarman.Option {
	options := []solarman.Option{solarman.WithTimeout(opts.Timeout)}
	if opts.Debug {
		options = append(options, solarman.WithDebugf(log.Printf))
	}
	if opts.DevMode {
		options = append(options, solarman.WithDevMode(true))
	}
	if opts.MaxAttempts > 0 {
		options = append(options, solarman.WithMaxAttempts(opts.MaxAttempts))
	}
	return options
}
