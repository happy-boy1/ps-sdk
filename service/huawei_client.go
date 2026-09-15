package service

import (
	"log"
	"strings"

	"ps-sdk/model"
	"ps-sdk/sdk/huawei"
)

// huaweiStationBatch 电站列表接口每页上限，同时决定设备列表的装机批量
const huaweiStationBatch = 100

// huaweiAdapter 华为 FusionSolar 北向接口适配器
type huaweiAdapter struct {
	limitBackoff
	sdk *huawei.FusionSolarSDK
}

func init() {
	registerAdapter(CodeFusionSolar, newHuaweiAdapter)
}

func newHuaweiAdapter(cred Credential, opts Options) (Adapter, error) {
	userName := firstNonEmpty(cred.Account, cred.AppID)
	systemCode := firstNonEmpty(cred.Password, cred.AppSecret)
	if userName == "" || systemCode == "" {
		return nil, wrapErr(cred.PlatformID, CodeFusionSolar, "构造适配器", ErrCredentialMissing)
	}

	options := []huawei.Option{huawei.WithTimeout(opts.Timeout)}
	if opts.Debug {
		options = append(options, huawei.WithDebugf(log.Printf))
	}
	if opts.DevMode {
		options = append(options, huawei.WithDevMode(true))
	}
	if opts.MaxAttempts > 0 {
		options = append(options, huawei.WithMaxAttempts(opts.MaxAttempts))
	}
	if opts.RetryDelay > 0 {
		options = append(options, huawei.WithRetryBaseDelay(opts.RetryDelay))
	}

	sdk, err := huawei.NewFusionSolarSDK(huawei.Credentials{
		UserName:   userName,
		SystemCode: systemCode,
		BaseURL:    strings.TrimSpace(cred.APIURL),
	}, options...)
	if err != nil {
		return nil, wrapErr(cred.PlatformID, CodeFusionSolar, "构造适配器", err)
	}
	return &huaweiAdapter{sdk: sdk}, nil
}

func (a *huaweiAdapter) Platform() Platform {
	p, _ := PlatformByCode(CodeFusionSolar)
	return p
}

func (a *huaweiAdapter) ListPowerStations() ([]model.PowerStation, error) {
	stations := make([]model.PowerStation, 0, huaweiStationBatch)
	for pageNo := 1; ; pageNo++ {
		var res *huawei.StationListResult
		err := a.throttle(isRateLimitError, func() error {
			result, err := a.sdk.GetStationList(huawei.StationListRequest{PageNo: pageNo})
			res = result
			return err
		})
		if err != nil {
			return nil, wrapErr(PlatformFusionSolar, CodeFusionSolar, "列出电站", err)
		}
		for i := range res.List {
			stations = append(stations, *huaweiStation(&res.List[i]))
		}
		if len(res.List) == 0 || len(stations) >= int(res.Total) {
			return stations, nil
		}
	}
}

func (a *huaweiAdapter) ListPowerDevices(stationID uint64) ([]model.PowerDevice, error) {
	ps, err := findStation(PlatformFusionSolar, stationID)
	if err != nil {
		return nil, wrapErr(PlatformFusionSolar, CodeFusionSolar, "读取电站", err)
	}

	var devices []huawei.Device
	err = a.throttleRetry(isRateLimitError, deviceListAttempts, func() error {
		result, err := a.sdk.GetDevList(huawei.DeviceListRequest{StationCodes: ps.StationIDOrigin})
		devices = result
		return err
	})
	if err != nil {
		return nil, wrapErr(PlatformFusionSolar, CodeFusionSolar, "列出设备", err)
	}

	out := make([]model.PowerDevice, 0, len(devices))
	for i := range devices {
		out = append(out, *huaweiDevice(&devices[i], stationID))
	}
	return out, nil
}

// huaweiStation 转换为统一电站模型
func huaweiStation(item *huawei.Station) *model.PowerStation {
	return &model.PowerStation{
		PlatformID:      PlatformFusionSolar,
		StationIDOrigin: item.PlantCode,
		StationName:     item.PlantName,
		CapacityKwp:     float64Ptr(item.Capacity.Float()),
		Address:         item.PlantAddress,
		Longitude:       item.Longitude.Float(),
		Latitude:        item.Latitude.Float(),
		GridConnectedAt: dateTime(item.GridConnectionDate),
		// 电站列表不返回运行状态
		Status: StationRunning,
	}
}

// huaweiDevice 转换为统一设备模型
func huaweiDevice(item *huawei.Device, stationID uint64) *model.PowerDevice {
	origin := firstNonEmpty(item.DevDn, item.EsnCode, formatInt64(item.ID))
	deviceType := resolveDeviceType(CodeFusionSolar, formatInt64(int64(item.DevTypeID)))

	device := newPowerDevice(PlatformFusionSolar, stationID, deviceType.ID,
		origin, item.EsnCode, item.DevName, item.Model, item.InvType)
	device.DeviceTypeOrigin = formatInt64(int64(item.DevTypeID))
	return device
}
