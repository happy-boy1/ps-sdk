package service

// 电站状态。与 model.PowerStation.Status 对应，语义由表结构固定，
// 各平台适配器负责把自己平台的状态值换算到本域。
const (
	StationStopped    int8 = 0 // 停运
	StationRunning    int8 = 1 // 运行
	StationUnderConst int8 = 2 // 在建
)

// 设备状态。与 model.PowerDevice.Status 对应
const (
	DeviceUnknown int8 = 0 // 未知
	DeviceOnline  int8 = 1 // 在线
	DeviceOffline int8 = 2 // 离线
	DeviceAlarm   int8 = 3 // 告警
)

// StationStatusName 返回电站状态中文名
func StationStatusName(status int8) string {
	switch status {
	case StationRunning:
		return "运行"
	case StationStopped:
		return "停运"
	case StationUnderConst:
		return "在建"
	default:
		return "未知"
	}
}

// DeviceStatusName 返回设备状态中文名
func DeviceStatusName(status int8) string {
	switch status {
	case DeviceOnline:
		return "在线"
	case DeviceOffline:
		return "离线"
	case DeviceAlarm:
		return "告警"
	default:
		return "未知"
	}
}

// stationStatusFromOnline 由「是否在线」推导电站状态
func stationStatusFromOnline(online bool) int8 {
	if online {
		return StationRunning
	}
	return StationStopped
}
