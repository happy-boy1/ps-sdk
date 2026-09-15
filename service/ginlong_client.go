package service

import (
	"log"
	"strings"

	"ps-sdk/model"
	"ps-sdk/sdk/ginlong"
)

// ginlongAdapter 锦浪云（SolisCloud）适配器。
// 平台没有统一的设备列表接口，按设备族逐个拉取后合并。
type ginlongAdapter struct {
	limitBackoff
	sdk *ginlong.SolisSDK
}

func init() {
	registerAdapter(CodeGinlong, newGinlongAdapter)
}

func newGinlongAdapter(cred Credential, opts Options) (Adapter, error) {
	if strings.TrimSpace(cred.AppID) == "" || strings.TrimSpace(cred.AppSecret) == "" {
		return nil, wrapErr(cred.PlatformID, CodeGinlong, "构造适配器", ErrCredentialMissing)
	}

	options := []ginlong.Option{ginlong.WithTimeout(opts.Timeout)}
	if opts.Debug {
		options = append(options, ginlong.WithDebugf(log.Printf))
	}
	if opts.DevMode {
		options = append(options, ginlong.WithDevMode(true))
	}

	sdk, err := ginlong.NewSolisSDK(ginlong.Credentials{
		APIID:     strings.TrimSpace(cred.AppID),
		APISecret: strings.TrimSpace(cred.AppSecret),
		BaseURL:   strings.TrimSpace(cred.APIURL),
	}, options...)
	if err != nil {
		return nil, wrapErr(cred.PlatformID, CodeGinlong, "构造适配器", err)
	}
	return &ginlongAdapter{sdk: sdk}, nil
}

func (a *ginlongAdapter) Platform() Platform {
	p, _ := PlatformByCode(CodeGinlong)
	return p
}

func (a *ginlongAdapter) ListPowerStations() ([]model.PowerStation, error) {
	stations := make([]model.PowerStation, 0, stationPageSize)
	for pageNo := 1; ; pageNo++ {
		var res *ginlong.UserStationListResult
		err := a.throttle(isRateLimitError, func() error {
			result, err := a.sdk.UserStationList(ginlong.UserStationListRequest{
				PageRequest: ginlong.PageRequest{PageNo: pageNo, PageSize: stationPageSize},
			})
			res = result
			return err
		})
		if err != nil {
			return nil, wrapErr(PlatformGinlong, CodeGinlong, "列出电站", err)
		}
		for i := range res.Page.Records {
			stations = append(stations, *ginlongStation(&res.Page.Records[i]))
		}
		if len(res.Page.Records) == 0 || len(stations) >= res.Page.Total {
			return stations, nil
		}
	}
}

func (a *ginlongAdapter) ListPowerDevices(stationID uint64) ([]model.PowerDevice, error) {
	ps, err := findStation(PlatformGinlong, stationID)
	if err != nil {
		return nil, wrapErr(PlatformGinlong, CodeGinlong, "读取电站", err)
	}
	origin := parseInt64(ps.StationIDOrigin)
	page := ginlong.PageRequest{PageNo: 1, PageSize: devicePageSize}

	// 平台按设备族提供接口，任一族的权限缺失不应阻断其余设备
	families := []struct {
		name  string
		fetch func() ([]model.PowerDevice, error)
	}{
		{"逆变器", func() ([]model.PowerDevice, error) {
			res, err := a.sdk.InverterList(ginlong.InverterListRequest{PageRequest: page, StationID: ginlong.Int64(origin)})
			if err != nil {
				return nil, err
			}
			out := make([]model.PowerDevice, 0, len(res.Page.Records))
			for i := range res.Page.Records {
				out = append(out, *ginlongInverter(&res.Page.Records[i], stationID))
			}
			return out, nil
		}},
		{"采集器", func() ([]model.PowerDevice, error) {
			res, err := a.sdk.CollectorList(ginlong.CollectorListRequest{PageRequest: page, StationID: ginlong.Int64(origin)})
			if err != nil {
				return nil, err
			}
			out := make([]model.PowerDevice, 0, len(res.Page.Records))
			for i := range res.Page.Records {
				item := &res.Page.Records[i]
				out = append(out, *newPowerDevice(PlatformGinlong, stationID, DeviceTypeCollector,
					item.SN, item.SN, item.Name, item.Model, ""))
			}
			return out, nil
		}},
		{"EPM", func() ([]model.PowerDevice, error) {
			res, err := a.sdk.EpmList(ginlong.EpmListRequest{PageRequest: page, StationID: ps.StationIDOrigin})
			if err != nil {
				return nil, err
			}
			out := make([]model.PowerDevice, 0, len(res.Page.Records))
			for i := range res.Page.Records {
				item := &res.Page.Records[i]
				out = append(out, *newPowerDevice(PlatformGinlong, stationID, DeviceTypeEnergyManagement,
					item.SN, item.SN, item.SN, "", ""))
			}
			return out, nil
		}},
		{"电表", func() ([]model.PowerDevice, error) {
			res, err := a.sdk.AmmeterList(ginlong.AmmeterListRequest{PageRequest: page, StationID: ginlong.Int64(origin)})
			if err != nil {
				return nil, err
			}
			out := make([]model.PowerDevice, 0, len(res.Page.Records))
			for i := range res.Page.Records {
				item := &res.Page.Records[i]
				out = append(out, *newPowerDevice(PlatformGinlong, stationID, DeviceTypeMeter,
					item.SN, item.SN, item.Name, "", ""))
			}
			return out, nil
		}},
		{"气象仪", func() ([]model.PowerDevice, error) {
			res, err := a.sdk.WeatherList(ginlong.WeatherListRequest{PageRequest: page, StationID: ginlong.Int64(origin)})
			if err != nil {
				return nil, err
			}
			out := make([]model.PowerDevice, 0, len(res.Page.Records))
			for i := range res.Page.Records {
				item := &res.Page.Records[i]
				out = append(out, *newPowerDevice(PlatformGinlong, stationID, DeviceTypeWeatherStation,
					item.CollectorSN, item.CollectorSN, item.Name, item.WeatherModel, ""))
			}
			return out, nil
		}},
	}

	devices := make([]model.PowerDevice, 0, devicePageSize)
	warnings := make([]error, 0, len(families))
	for _, family := range families {
		var items []model.PowerDevice
		err := a.throttle(isRateLimitError, func() error {
			result, err := family.fetch()
			items = result
			return err
		})
		if err != nil {
			// 常见于接口未开通权限（R0000），记录后继续拉取其余设备族
			warnings = append(warnings, wrapErr(PlatformGinlong, CodeGinlong, "列出"+family.name, err))
			Logf("[ginlong] 电站 %s 的%s列表不可用: %v", ps.StationName, family.name, err)
			continue
		}
		devices = append(devices, items...)
	}

	if len(warnings) == 0 {
		return devices, nil
	}
	return nil, &PartialError{Devices: devices, Warnings: warnings}
}

// ginlongStation 转换为统一电站模型
func ginlongStation(item *ginlong.UserStationListItem) *model.PowerStation {
	return &model.PowerStation{
		PlatformID:       PlatformGinlong,
		StationIDOrigin:  formatInt64(item.ID.Int()),
		StationName:      item.StationName,
		StationShortName: item.Sno,
		CapacityKwp:      float64Ptr(item.Capacity.Float()),
		Province:         item.RegionStr,
		City:             item.CityStr,
		Address:          item.Addr,
		Longitude:        parseFloat(item.Longitude),
		Latitude:         parseFloat(item.Latitude),
		GridConnectedAt:  unixMillis(item.ConnectTime.Int()),
		Status:           ginlongStationStatus(int(item.State)),
	}
}

// ginlongInverter 逆变器列表项转换为统一设备模型
func ginlongInverter(item *ginlong.InverterListItem, stationID uint64) *model.PowerDevice {
	device := newPowerDevice(PlatformGinlong, stationID, DeviceTypeInverter,
		item.SN, item.SN, firstNonEmpty(item.Name, item.SN), item.ProductModel, "")
	device.RatedPowerKw = float64Ptr(item.Power.Float())
	device.Status = ginlongDeviceStatus(int(item.State))
	return device
}

// ginlongStationStatus 锦浪电站状态：1 在线 2 离线 3 报警
func ginlongStationStatus(state int) int8 {
	switch state {
	case 1:
		return StationRunning
	case 3:
		return StationRunning
	default:
		return StationStopped
	}
}

// ginlongDeviceStatus 锦浪设备状态：1 在线 2 离线 3 报警
func ginlongDeviceStatus(state int) int8 {
	switch state {
	case 1:
		return DeviceOnline
	case 2:
		return DeviceOffline
	case 3:
		return DeviceAlarm
	default:
		return DeviceUnknown
	}
}

// newPowerDevice 组装统一设备模型，origin 为空时用 SN 兜底
func newPowerDevice(platformID int16, stationID uint64, typeID int16, origin, sn, name, modelName, brand string) *model.PowerDevice {
	deviceType, ok := DeviceTypeDefByID(typeID)
	if !ok {
		deviceType, _ = DeviceTypeDefByID(DeviceTypeUnknown)
	}

	origin = strings.TrimSpace(origin)
	sn = strings.TrimSpace(sn)
	if origin == "" {
		origin = sn
	}

	return &model.PowerDevice{
		PlatformID:       platformID,
		StationID:        stationID,
		DeviceTypeID:     deviceType.ID,
		DeviceIDOrigin:   origin,
		DeviceSN:         sn,
		DeviceTypeOrigin: deviceType.Code,
		DeviceName:       strings.TrimSpace(name),
		DeviceModel:      strings.TrimSpace(modelName),
		Brand:            strings.TrimSpace(brand),
	}
}
