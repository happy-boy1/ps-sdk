package service

import (
	"sort"
	"strings"

	"ps-sdk/model"
)

// 统一设备类型 ID。各平台的原始类型编码经 DeviceTypeOrigin 保留，统一 ID 用于跨平台聚合。
//
// ID 分段（device_type_id 为 smallint，上限 32767）：
//
//	11xx 光伏类，12xx 储能类，13xx 计量与采集类，14xx 环境与其他
const (
	DeviceTypeUnknown int16 = 0

	DeviceTypeInverter       int16 = 1101 // 光伏逆变器（含组串式、集中式、户用）
	DeviceTypeMicroInverter  int16 = 1102 // 微型逆变器
	DeviceTypeOptimizer      int16 = 1103 // 功率优化器
	DeviceTypePVModule       int16 = 1104 // 光伏组件
	DeviceTypeCombinerBox    int16 = 1105 // 汇流箱
	DeviceTypeString         int16 = 1106 // 组串
	DeviceTypeTrackingSystem int16 = 1107 // 跟踪支架

	DeviceTypePCS                int16 = 1201 // 储能变流器
	DeviceTypeBattery            int16 = 1202 // 电池簇
	DeviceTypeBatteryClusterUnit int16 = 1203 // 电池簇管理单元
	DeviceTypeBMS                int16 = 1204 // 电池管理系统
	DeviceTypeEnergyManagement   int16 = 1205 // 能量管理系统
	DeviceTypeLocalController    int16 = 1206 // 本地控制器
	DeviceTypeContainer          int16 = 1207 // 集装箱
	DeviceTypeResidentialStorage int16 = 1208 // 户用/一体式储能机

	DeviceTypeCollector   int16 = 1301 // 采集器 / 数采 / DTU / 网关
	DeviceTypeMeter       int16 = 1302 // 电表
	DeviceTypeGridMeter   int16 = 1303 // 关口电表
	DeviceTypeUPS         int16 = 1304 // UPS
	DeviceTypeTransformer int16 = 1305 // 变压器
	DeviceTypeGridPoint   int16 = 1306 // 并网点

	DeviceTypeWeatherStation int16 = 1401 // 气象站 / 环境监测仪
	DeviceTypeIrradiance     int16 = 1402 // 辐照仪
	DeviceTypeTemperature    int16 = 1403 // 温湿度传感器
	DeviceTypeFan            int16 = 1404 // 风机
	DeviceTypeCharger        int16 = 1405 // 充电桩
	DeviceTypeDiesel         int16 = 1406 // 柴油发电机
)

// DeviceTypeDef 统一设备类型定义
type DeviceTypeDef struct {
	ID        int16
	Code      string // 统一类型编码，全局唯一
	Name      string // 类型名称
	ShortName string // 类型简称，用于紧凑展示
}

// 统一类型编码
const (
	CodeTypeInverter       = "INVERTER"
	CodeTypeMicroInverter  = "MICRO_INVERTER"
	CodeTypeOptimizer      = "OPTIMIZER"
	CodeTypePVModule       = "PV_MODULE"
	CodeTypeCombinerBox    = "COMBINER_BOX"
	CodeTypeString         = "STRING"
	CodeTypeTrackingSystem = "TRACKING_SYSTEM"

	CodeTypePCS                = "PCS"
	CodeTypeBattery            = "BATTERY"
	CodeTypeBatteryClusterUnit = "BATTERY_CLUSTER_UNIT"
	CodeTypeBMS                = "BMS"
	CodeTypeEnergyManagement   = "EMS"
	CodeTypeLocalController    = "LOCAL_CONTROLLER"
	CodeTypeContainer          = "CONTAINER"
	CodeTypeResidentialStorage = "RESIDENTIAL_STORAGE"

	CodeTypeCollector   = "COLLECTOR"
	CodeTypeMeter       = "METER"
	CodeTypeGridMeter   = "GRID_METER"
	CodeTypeUPS         = "UPS"
	CodeTypeTransformer = "TRANSFORMER"
	CodeTypeGridPoint   = "GRID_POINT"

	CodeTypeWeatherStation = "WEATHER_STATION"
	CodeTypeIrradiance     = "IRRADIANCE"
	CodeTypeTemperature    = "TEMPERATURE"
	CodeTypeFan            = "FAN"
	CodeTypeCharger        = "CHARGER"
	CodeTypeDiesel         = "DIESEL_GENERATOR"

	CodeTypeUnknown = "UNKNOWN"
)

// deviceTypeDefs 统一设备类型目录，保证 DeviceTypeID 与 TypeCode 双向稳定
var deviceTypeDefs = []DeviceTypeDef{
	{ID: DeviceTypeInverter, Code: CodeTypeInverter, Name: "逆变器", ShortName: "INV"},
	{ID: DeviceTypeMicroInverter, Code: CodeTypeMicroInverter, Name: "微型逆变器", ShortName: "MICRO_INV"},
	{ID: DeviceTypeOptimizer, Code: CodeTypeOptimizer, Name: "功率优化器", ShortName: "OPT"},
	{ID: DeviceTypePVModule, Code: CodeTypePVModule, Name: "光伏组件", ShortName: "PV"},
	{ID: DeviceTypeCombinerBox, Code: CodeTypeCombinerBox, Name: "汇流箱", ShortName: "COMB"},
	{ID: DeviceTypeString, Code: CodeTypeString, Name: "组串", ShortName: "STR"},
	{ID: DeviceTypeTrackingSystem, Code: CodeTypeTrackingSystem, Name: "跟踪支架", ShortName: "TRACK"},

	{ID: DeviceTypePCS, Code: CodeTypePCS, Name: "储能变流器", ShortName: "PCS"},
	{ID: DeviceTypeBattery, Code: CodeTypeBattery, Name: "电池簇", ShortName: "BAT"},
	{ID: DeviceTypeBatteryClusterUnit, Code: CodeTypeBatteryClusterUnit, Name: "电池簇管理单元", ShortName: "BCU"},
	{ID: DeviceTypeBMS, Code: CodeTypeBMS, Name: "电池管理系统", ShortName: "BMS"},
	{ID: DeviceTypeEnergyManagement, Code: CodeTypeEnergyManagement, Name: "能量管理系统", ShortName: "EMS"},
	{ID: DeviceTypeLocalController, Code: CodeTypeLocalController, Name: "本地控制器", ShortName: "LC"},
	{ID: DeviceTypeContainer, Code: CodeTypeContainer, Name: "储能集装箱", ShortName: "CONT"},
	{ID: DeviceTypeResidentialStorage, Code: CodeTypeResidentialStorage, Name: "一体式储能机", ShortName: "ESS"},

	{ID: DeviceTypeCollector, Code: CodeTypeCollector, Name: "采集器", ShortName: "DLOG"},
	{ID: DeviceTypeMeter, Code: CodeTypeMeter, Name: "电表", ShortName: "METER"},
	{ID: DeviceTypeGridMeter, Code: CodeTypeGridMeter, Name: "关口电表", ShortName: "GMETER"},
	{ID: DeviceTypeUPS, Code: CodeTypeUPS, Name: "UPS", ShortName: "UPS"},
	{ID: DeviceTypeTransformer, Code: CodeTypeTransformer, Name: "变压器", ShortName: "TRANS"},
	{ID: DeviceTypeGridPoint, Code: CodeTypeGridPoint, Name: "并网点", ShortName: "GRID"},

	{ID: DeviceTypeWeatherStation, Code: CodeTypeWeatherStation, Name: "气象站", ShortName: "WST"},
	{ID: DeviceTypeIrradiance, Code: CodeTypeIrradiance, Name: "辐照仪", ShortName: "IRR"},
	{ID: DeviceTypeTemperature, Code: CodeTypeTemperature, Name: "温湿度传感器", ShortName: "THS"},
	{ID: DeviceTypeFan, Code: CodeTypeFan, Name: "风机", ShortName: "FAN"},
	{ID: DeviceTypeCharger, Code: CodeTypeCharger, Name: "充电桩", ShortName: "CHG"},
	{ID: DeviceTypeDiesel, Code: CodeTypeDiesel, Name: "柴油发电机", ShortName: "DIESEL"},

	{ID: DeviceTypeUnknown, Code: CodeTypeUnknown, Name: "未知设备", ShortName: "UNK"},
}

var (
	deviceTypeByID   = make(map[int16]DeviceTypeDef, len(deviceTypeDefs))
	deviceTypeByCode = make(map[string]DeviceTypeDef, len(deviceTypeDefs))
)

func init() {
	for _, d := range deviceTypeDefs {
		deviceTypeByID[d.ID] = d
		deviceTypeByCode[d.Code] = d
	}
}

// DeviceTypes 返回统一设备类型目录的副本
func DeviceTypes() []DeviceTypeDef {
	out := make([]DeviceTypeDef, len(deviceTypeDefs))
	copy(out, deviceTypeDefs)
	return out
}

// DeviceTypeDefByID 按统一 ID 查设备类型
func DeviceTypeDefByID(id int16) (DeviceTypeDef, bool) {
	d, ok := deviceTypeByID[id]
	return d, ok
}

// DeviceTypeDefByCode 按统一编码查设备类型
func DeviceTypeDefByCode(code string) (DeviceTypeDef, bool) {
	d, ok := deviceTypeByCode[strings.ToUpper(strings.TrimSpace(code))]
	return d, ok
}

// DeviceTypeModels 返回可直接写入 device_type 表的记录
func DeviceTypeModels() []model.DeviceType {
	out := make([]model.DeviceType, 0, len(deviceTypeDefs))
	for _, d := range deviceTypeDefs {
		out = append(out, model.DeviceType{
			DeviceTypeID:  d.ID,
			TypeCode:      d.Code,
			TypeName:      d.Name,
			TypeShortName: d.ShortName,
		})
	}
	return out
}

// typeAlias 平台原始类型编码 → 统一设备类型。
// key 为「平台短码 + 原始编码」，编码统一转大写后匹配。
var typeAlias = map[string]int16{}

// registerAlias 注册平台原始编码到统一类型的映射
func registerAlias(typeID int16, platformCode string, ids ...string) {
	for _, id := range ids {
		id = strings.ToUpper(strings.TrimSpace(id))
		if id == "" {
			continue
		}
		typeAlias[platformCode+":"+id] = typeID
	}
}

func init() {
	// SolarMan 设备类型为字符串枚举
	registerAlias(DeviceTypeInverter, CodeSolarman, "INVERTER")
	registerAlias(DeviceTypeMicroInverter, CodeSolarman, "MICRO_INVERTER")
	registerAlias(DeviceTypePVModule, CodeSolarman, "PV_MODULE")
	registerAlias(DeviceTypeCollector, CodeSolarman, "COLLECTOR", "DTU", "REPEATER")
	registerAlias(DeviceTypeMeter, CodeSolarman, "METER")
	registerAlias(DeviceTypeWeatherStation, CodeSolarman, "WEATHER_STATION")
	registerAlias(DeviceTypeBattery, CodeSolarman, "BATTERY")
	registerAlias(DeviceTypeFan, CodeSolarman, "FAN")
	registerAlias(DeviceTypeGridPoint, CodeSolarman, "SWITCHGEAR")

	// 阳光云设备类型为数字编码
	registerAlias(DeviceTypeInverter, CodeSungrow, "1")
	registerAlias(DeviceTypeContainer, CodeSungrow, "2")
	registerAlias(DeviceTypeGridPoint, CodeSungrow, "3")
	registerAlias(DeviceTypeCombinerBox, CodeSungrow, "4")
	registerAlias(DeviceTypeWeatherStation, CodeSungrow, "5")
	registerAlias(DeviceTypeTransformer, CodeSungrow, "6")
	registerAlias(DeviceTypeMeter, CodeSungrow, "7")
	registerAlias(DeviceTypeUPS, CodeSungrow, "8")
	registerAlias(DeviceTypeCollector, CodeSungrow, "9", "15", "22")
	registerAlias(DeviceTypeString, CodeSungrow, "10")
	registerAlias(DeviceTypePCS, CodeSungrow, "37")
	registerAlias(DeviceTypeOptimizer, CodeSungrow, "41")
	registerAlias(DeviceTypeBattery, CodeSungrow, "43")
	registerAlias(DeviceTypeBatteryClusterUnit, CodeSungrow, "44")
	registerAlias(DeviceTypeLocalController, CodeSungrow, "45")
	registerAlias(DeviceTypeCharger, CodeSungrow, "51")
	registerAlias(DeviceTypeResidentialStorage, CodeSungrow, "14", "52")
	registerAlias(DeviceTypeMicroInverter, CodeSungrow, "55")
	registerAlias(DeviceTypeDiesel, CodeSungrow, "63")
	registerAlias(DeviceTypeBMS, CodeSungrow, "23", "24")
	registerAlias(DeviceTypeEnergyManagement, CodeSungrow, "16", "26")
	registerAlias(DeviceTypeTrackingSystem, CodeSungrow, "27")
	registerAlias(DeviceTypeTemperature, CodeSungrow, "18")

	// 华为设备类型为数字编码
	registerAlias(DeviceTypeInverter, CodeFusionSolar, "1", "38")
	registerAlias(DeviceTypeCollector, CodeFusionSolar, "2", "37", "63")
	registerAlias(DeviceTypeTransformer, CodeFusionSolar, "8")
	registerAlias(DeviceTypeWeatherStation, CodeFusionSolar, "10")
	registerAlias(DeviceTypeGridMeter, CodeFusionSolar, "17")
	registerAlias(DeviceTypePVModule, CodeFusionSolar, "60043", "60044")
	registerAlias(DeviceTypeBattery, CodeFusionSolar, "39", "60014")
	registerAlias(DeviceTypePCS, CodeFusionSolar, "60092")
	registerAlias(DeviceTypeOptimizer, CodeFusionSolar, "46")
	registerAlias(DeviceTypeMeter, CodeFusionSolar, "47")
	registerAlias(DeviceTypeResidentialStorage, CodeFusionSolar, "23093")

	// 锦浪的设备按调用接口区分，编码取统一编码本身
	for _, d := range deviceTypeDefs {
		registerAlias(d.ID, CodeGinlong, d.Code)
	}
}

// resolveDeviceType 把平台原始类型编码解析为统一设备类型，未登记时返回未知类型
func resolveDeviceType(platformCode, origin string) DeviceTypeDef {
	key := platformCode + ":" + strings.ToUpper(strings.TrimSpace(origin))
	if id, ok := typeAlias[key]; ok {
		if d, ok := deviceTypeByID[id]; ok {
			return d
		}
	}
	return deviceTypeByID[DeviceTypeUnknown]
}

// KnownDeviceTypes 返回已登记的类型编码，便于排查未映射的平台类型
func KnownDeviceTypes(platformCode string) []string {
	prefix := normalizeCode(platformCode) + ":"
	out := make([]string, 0, len(typeAlias))
	for k := range typeAlias {
		if strings.HasPrefix(k, prefix) {
			out = append(out, strings.TrimPrefix(k, prefix))
		}
	}
	sort.Strings(out)
	return out
}
