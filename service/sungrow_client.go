package service

import (
	"log"
	"strings"

	"ps-sdk/model"
	"ps-sdk/sdk/sungrow"
)

// sungrowAdapter 阳光云（iSolarCloud）适配器
type sungrowAdapter struct {
	sdk *sungrow.SungrowSDK
}

func init() {
	registerAdapter(CodeSungrow, newSungrowAdapter)
}

func newSungrowAdapter(cred Credential, opts Options) (Adapter, error) {
	appID := strings.TrimSpace(cred.AppID)
	appSecret := strings.TrimSpace(cred.AppSecret)
	account := firstNonEmpty(cred.Account, appID)
	password := firstNonEmpty(cred.Password, appSecret)
	if appID == "" || appSecret == "" || account == "" || password == "" {
		return nil, wrapErr(cred.PlatformID, CodeSungrow, "构造适配器", ErrCredentialMissing)
	}

	options := []sungrow.Option{sungrow.WithTimeout(opts.Timeout)}
	if opts.Debug {
		options = append(options, sungrow.WithDebugf(log.Printf))
	}
	if opts.DevMode {
		options = append(options, sungrow.WithDevMode(true))
	}
	if opts.MaxAttempts > 0 {
		options = append(options, sungrow.WithMaxAttempts(opts.MaxAttempts))
	}
	if opts.RetryDelay > 0 {
		options = append(options, sungrow.WithRetryBaseDelay(opts.RetryDelay))
	}

	sdk, err := sungrow.NewSungrowSDK(sungrow.Credentials{
		AppID:        appID,
		AppSecret:    appSecret,
		UserAccount:  account,
		UserPassword: password,
		BaseURL:      strings.TrimSpace(cred.APIURL),
	}, options...)
	if err != nil {
		return nil, wrapErr(cred.PlatformID, CodeSungrow, "构造适配器", err)
	}
	return &sungrowAdapter{sdk: sdk}, nil
}

func (a *sungrowAdapter) Platform() Platform {
	p, _ := PlatformByCode(CodeSungrow)
	return p
}

func (a *sungrowAdapter) ListPowerStations() ([]model.PowerStation, error) {
	stations := make([]model.PowerStation, 0, stationPageSize)
	for curPage := 1; ; curPage++ {
		res, err := a.sdk.GetPowerStationList(sungrow.PowerStationListRequest{
			PageRequest: sungrow.PageRequest{CurPage: curPage, Size: stationPageSize},
		})
		if err != nil {
			return nil, wrapErr(PlatformSungrow, CodeSungrow, "列出电站", err)
		}
		for i := range res.PageList {
			stations = append(stations, *sungrowStation(&res.PageList[i]))
		}
		if len(res.PageList) == 0 || len(stations) >= res.RowCount {
			return stations, nil
		}
	}
}

func (a *sungrowAdapter) ListPowerDevices(stationID uint64) ([]model.PowerDevice, error) {
	ps, err := findStation(PlatformSungrow, stationID)
	if err != nil {
		return nil, wrapErr(PlatformSungrow, CodeSungrow, "读取电站", err)
	}

	devices := make([]model.PowerDevice, 0, devicePageSize)
	for curPage := 1; ; curPage++ {
		res, err := a.sdk.GetDeviceListByPsID(sungrow.DeviceListByPsIDRequest{
			PageRequest: sungrow.PageRequest{CurPage: curPage, Size: devicePageSize},
			PsID:        ps.StationIDOrigin,
		})
		if err != nil {
			return nil, wrapErr(PlatformSungrow, CodeSungrow, "列出设备", err)
		}
		for i := range res.PageList {
			devices = append(devices, *sungrowDevice(&res.PageList[i], stationID))
		}
		if len(res.PageList) == 0 || len(devices) >= res.RowCount {
			return devices, nil
		}
	}
}

// sungrowStation 转换为统一电站模型
func sungrowStation(item *sungrow.PowerStation) *model.PowerStation {
	return &model.PowerStation{
		PlatformID:      PlatformSungrow,
		StationIDOrigin: item.PsId.String(),
		StationName:     item.PsName,
		CapacityKwp:     float64Ptr(parseFloat(item.TotalCapcity.Value.String())),
		Province:        item.ProvinceName,
		City:            item.CityName,
		Address:         item.PsLocation,
		Longitude:       item.Longitude.Float(),
		Latitude:        item.Latitude.Float(),
		GridConnectedAt: unixMillis(item.GridConnectionTime),
		Status:          sungrowStatus(int(item.PsStatus)),
	}
}

// sungrowDevice 转换为统一设备模型
func sungrowDevice(item *sungrow.Device, stationID uint64) *model.PowerDevice {
	// device_type 的文档类型与真实返回不一致，统一转字符串后再归一
	origin := firstNonEmpty(item.DeviceSn, formatInt64(item.Uuid.Int()))
	typeOrigin := item.DeviceType.String()
	deviceType := resolveDeviceType(CodeSungrow, typeOrigin)

	device := newPowerDevice(PlatformSungrow, stationID, deviceType.ID,
		origin, item.DeviceSn, item.DeviceName, item.DeviceModelCode, item.FactoryName)
	device.DeviceTypeOrigin = typeOrigin
	device.InstalledAt = dateTime(item.GridConnectionDate)
	device.Status = sungrowDevStatus(item.DevStatus.String(), item.DevFaultStatus)
	return device
}

// sungrowStatus 电站状态：1 在线 0 离线，缺省按停运处理
func sungrowStatus(status int) int8 {
	if status == 1 {
		return StationRunning
	}
	return StationStopped
}

// sungrowDevStatus 设备状态：dev_status 为文本，故障状态优先判定告警
func sungrowDevStatus(status string, fault int64) int8 {
	if fault != 0 {
		return DeviceAlarm
	}
	switch strings.TrimSpace(status) {
	case "1":
		return DeviceOnline
	case "2", "0":
		return DeviceOffline
	default:
		return DeviceUnknown
	}
}
