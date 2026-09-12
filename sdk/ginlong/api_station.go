package ginlong

import (
	"encoding/json"
	"fmt"
	"strings"
)

// unmarshalPaged 解析分页数据，兼容 page 嵌套与 data 下平铺两种返回形态。
// 4.4/4.5/4.6 的文档表格把 total/records 画在 page 对象下，官方示例却是平铺的，
// 两种形态都写入同一个 Page 字段，调用方始终读 res.Page.Records
func unmarshalPaged[T any](data []byte, out *Page[T]) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if raw, ok := fields["page"]; ok && string(raw) != "null" {
		return json.Unmarshal(raw, out)
	}
	return json.Unmarshal(data, out)
}

// ---------- 4.1 获取账号下电站列表 /v1/api/userStationList ----------

// UserStationListRequest 账号下电站列表请求
type UserStationListRequest struct {
	PageRequest         // 分页参数
	NmiCode     string  `json:"nmiCode,omitempty"` // nmi码
	IDList      []Int64 `json:"idList,omitempty"`  // 批量电站id，空为全部
}

// StationStatusVo 电站状态统计
type StationStatusVo struct {
	All     int `json:"all"`     // 电站总数
	Normal  int `json:"normal"`  // 正常电站数
	Offline int `json:"offline"` // 离线电站数
	Fault   int `json:"fault"`   // 故障电站数
}

// UserStationListResult 账号下电站列表结果
type UserStationListResult struct {
	rawHolder
	Page            Page[UserStationListItem] `json:"page"`            // 分页列表
	StationStatusVo StationStatusVo           `json:"stationStatusVo"` // 状态统计
}

// UserStationListItem 电站列表项
type UserStationListItem struct {
	rawHolder
	ID                          Int64               `json:"id"`                          // 电站id
	StationName                 string              `json:"stationName"`                 // 电站名称
	Addr                        string              `json:"addr"`                        // 电站地址
	UserID                      Int64               `json:"userId"`                      // 业主Id
	UserName                    string              `json:"userName"`                    // 业主姓名
	Mobile                      Str                 `json:"mobile"`                      // 手机号
	Capacity                    Num                 `json:"capacity"`                    // 装机容量
	CapacityStr                 string              `json:"capacityStr"`                 // 装机容量单位
	Capacity1                   Num                 `json:"capacity1"`                   // 装机容量不进位
	FullHour                    Num                 `json:"fullHour"`                    // 满发小时数
	PicName                     string              `json:"picName"`                     // 图片
	Latitude                    string              `json:"latitude"`                    // 纬度
	Longitude                   string              `json:"longitude"`                   // 经度
	InstallerID                 Int64               `json:"installerId"`                 // 安装商组织id
	Installer                   string              `json:"installer"`                   // 安装商组织
	InstallerMobile             Str                 `json:"installerMobile"`             // 安装商电话
	InstallerEmail              string              `json:"installerEmail"`              // 安装商邮箱
	UserMobile                  Str                 `json:"userMobile"`                  // 业主电话
	UserEmail                   string              `json:"userEmail"`                   // 业主邮箱
	Email                       string              `json:"email"`                       // 邮箱
	Sno                         string              `json:"sno"`                         // 电站短ID
	Country                     Int64               `json:"country"`                     // 国家id
	CountryStr                  string              `json:"countryStr"`                  // 国家名称
	Region                      Int64               `json:"region"`                      // 区域id
	RegionStr                   string              `json:"regionStr"`                   // 区域名称
	City                        Int64               `json:"city"`                        // 城市id
	CityStr                     string              `json:"cityStr"`                     // 城市名称
	County                      Int64               `json:"county"`                      // 区id
	CountyStr                   string              `json:"countyStr"`                   // 区名称
	Dip                         Num                 `json:"dip"`                         // 倾角
	Azimuth                     Num                 `json:"azimuth"`                     // 方位角
	TimeZone                    Num                 `json:"timeZone"`                    // 时区
	TimeZoneName                string              `json:"timeZoneName"`                // 时区名称
	TimeZoneStr                 string              `json:"timeZoneStr"`                 // 时区格式化字符串
	TimeZoneID                  Int64               `json:"timeZoneId"`                  // 时区id
	Daylight                    int                 `json:"daylight"`                    // 夏令时
	Price                       Num                 `json:"price"`                       // 每度电收益
	Module                      int                 `json:"module"`                      // 组件数量
	Pic1URL                     string              `json:"pic1Url"`                     // 电站照片url
	Power                       Num                 `json:"power"`                       // 功率
	PowerStr                    string              `json:"powerStr"`                    // 功率单位
	DayEnergy                   Num                 `json:"dayEnergy"`                   // 当日能量
	DayEnergyStr                string              `json:"dayEnergyStr"`                // 当日能量单位
	DayIncome                   Num                 `json:"dayIncome"`                   // 日收益
	DayIncomeUnit               string              `json:"dayIncomeUnit"`               // 日收益单位
	MonthEnergy                 Num                 `json:"monthEnergy"`                 // 月能量
	MonthEnergyStr              string              `json:"monthEnergyStr"`              // 月能量单位
	YearEnergy                  Num                 `json:"yearEnergy"`                  // 年能量
	YearEnergyStr               string              `json:"yearEnergyStr"`               // 年能量单位
	AllEnergy                   Num                 `json:"allEnergy"`                   // 总能量
	AllEnergyStr                string              `json:"allEnergyStr"`                // 总能量单位
	AllEnergy1                  Num                 `json:"allEnergy1"`                  // 累计能量原始值
	AllIncome                   Num                 `json:"allIncome"`                   // 累计收益
	AllIncomeUnit               string              `json:"allIncomeUnit"`               // 累计收益单位
	SynchronizationType         SynchronizationType `json:"synchronizationType"`         // 并网类型
	StationTypeNew              StationType         `json:"stationTypeNew"`              // 电站类型
	BatteryTotalDischargeEnergy Num                 `json:"batteryTotalDischargeEnergy"` // 电池累计放电量
	BatteryTotalChargeEnergy    Num                 `json:"batteryTotalChargeEnergy"`    // 电池累计充电量
	GridPurchasedTotalEnergy    Num                 `json:"gridPurchasedTotalEnergy"`    // 电表累计买电量
	GridSellTotalEnergy         Num                 `json:"gridSellTotalEnergy"`         // 电表累计卖电量
	HomeLoadTotalEnergy         Num                 `json:"homeLoadTotalEnergy"`         // 负载总用电量
	OneSelf                     Num                 `json:"oneSelf"`                     // 自发自用
	BatteryTodayDischargeEnergy Num                 `json:"batteryTodayDischargeEnergy"` // 电池当日放电量
	BatteryTodayChargeEnergy    Num                 `json:"batteryTodayChargeEnergy"`    // 电池当日充电量
	GridPurchasedTodayEnergy    Num                 `json:"gridPurchasedTodayEnergy"`    // 电表当日买电量
	GridSellTodayEnergy         Num                 `json:"gridSellTodayEnergy"`         // 电表当日卖电量
	HomeLoadTodayEnergy         Num                 `json:"homeLoadTodayEnergy"`         // 负载当日用电量
	Money                       string              `json:"money"`                       // 货币单位
	AccessTime                  Int64               `json:"accessTime"`                  // 接入平台时间
	ConnectTime                 Int64               `json:"connectTime"`                 // 并网时间
	CreateDate                  Int64               `json:"createDate"`                  // 创建时间
	Remark1                     string              `json:"remark1"`                     // 备注1
	Remark2                     string              `json:"remark2"`                     // 备注2
	Remark3                     string              `json:"remark3"`                     // 备注3
	State                       DeviceState         `json:"state"`                       // 电站状态
	DataTimestamp               Int64               `json:"dataTimestamp"`               // 更新时间
	InverterPower               Num                 `json:"inverterPower"`               // 逆变器额定功率和
	NmiCode                     string              `json:"nmiCode"`                     // nmi码
}

// UserStationList 获取账号下电站列表
func (c *SolisSDK) UserStationList(req UserStationListRequest) (*UserStationListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out UserStationListResult
	if err := c.do(PathUserStationList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 4.2 获取单个电站详情 /v1/api/stationDetail ----------

// StationDetailRequest 单个电站详情请求
type StationDetailRequest struct {
	ID      Int64  `json:"id,omitempty"`      // 电站id
	NmiCode string `json:"nmiCode,omitempty"` // nmi码
}

// StationDetailResult 单个电站详情结果
type StationDetailResult struct {
	rawHolder
	ID                          Int64       `json:"id"`                          // 电站id
	StationName                 string      `json:"stationName"`                 // 电站名称
	Addr                        string      `json:"addr"`                        // 电站地址
	UserID                      Int64       `json:"userId"`                      // 业主Id
	UserName                    string      `json:"userName"`                    // 业主姓名
	UserMobile                  Str         `json:"userMobile"`                  // 业主电话
	UserEmail                   string      `json:"userEmail"`                   // 业主邮箱
	Capacity                    Num         `json:"capacity"`                    // 装机容量
	CapacityPec                 string      `json:"capacityPec"`                 // 装机容量倍率
	CapacityStr                 string      `json:"capacityStr"`                 // 装机容量单位
	DayEnergy                   Num         `json:"dayEnergy"`                   // 当日能量
	DayEnergyStr                string      `json:"dayEnergyStr"`                // 当日能量单位
	MonthEnergy                 Num         `json:"monthEnergy"`                 // 当月能量
	MonthEnergyStr              string      `json:"monthEnergyStr"`              // 当月能量单位
	YearEnergy                  Num         `json:"yearEnergy"`                  // 当年能量
	YearEnergyStr               string      `json:"yearEnergyStr"`               // 当年能量单位
	AllEnergy                   Num         `json:"allEnergy"`                   // 总能量
	AllEnergyStr                string      `json:"allEnergyStr"`                // 总能量单位
	DayInCome                   Num         `json:"dayInCome"`                   // 当日收益
	DayInComeUnit               string      `json:"dayInComeUnit"`               // 当日收益单位
	MonthInCome                 Num         `json:"monthInCome"`                 // 当月收益
	MonthInComeUnit             string      `json:"monthInComeUnit"`             // 当月收益单位
	YearInCome                  Num         `json:"yearInCome"`                  // 当年收益
	YearInComeUnit              string      `json:"yearInComeUnit"`              // 当年收益单位
	AllInCome                   Num         `json:"allInCome"`                   // 总收益
	AllInComeUnit               string      `json:"allInComeUnit"`               // 总收益单位
	FullHour                    Num         `json:"fullHour"`                    // 满发小时数
	PicName                     string      `json:"picName"`                     // 图片
	Power                       Num         `json:"power"`                       // 功率
	PowerPec                    Num         `json:"powerPec"`                    // 功率倍率
	PowerStr                    string      `json:"powerStr"`                    // 功率单位
	Latitude                    string      `json:"latitude"`                    // 纬度
	Longitude                   string      `json:"longitude"`                   // 经度
	Dip                         Num         `json:"dip"`                         // 倾角
	Azimuth                     Num         `json:"azimuth"`                     // 方位角
	Price                       Num         `json:"price"`                       // 每度电收益
	State                       DeviceState `json:"state"`                       // 电站状态
	DataTimestamp               Int64       `json:"dataTimestamp"`               // 更新时间
	Money                       string      `json:"money"`                       // 货币种类
	Brand                       string      `json:"brand"`                       // 品牌
	CondTxtN                    string      `json:"condTxtN"`                    // 晚上天气
	CondTxtD                    string      `json:"condTxtD"`                    // 白天天气
	TmpMax                      string      `json:"tmpMax"`                      // 最高温度
	TmpMin                      string      `json:"tmpMin"`                      // 最低温度
	TmpUnit                     string      `json:"tmpUnit"`                     // 温度单位
	Sr                          string      `json:"sr"`                          // 日出时间
	Ss                          string      `json:"ss"`                          // 日落时间
	Hum                         string      `json:"hum"`                         // 湿度
	WindDir                     string      `json:"windDir"`                     // 风向
	WindSpd                     string      `json:"windSpd"`                     // 风速
	WeatherUpdateDate           string      `json:"weatherUpdateDate"`           // 天气更新时间
	PowerStationNumTree         Num         `json:"powerStationNumTree"`         // 等效植树量
	PowerStationNumTreeUnit     string      `json:"powerStationNumTreeUnit"`     // 植物单位
	PowerStationAvoidedCo2      Num         `json:"powerStationAvoidedCo2"`      // 二氧化碳减排
	PowerStationAvoidedCo2Unit  string      `json:"powerStationAvoidedCo2Unit"`  // 减排单位
	Module                      int         `json:"module"`                      // 组件数量
	Mobile                      Str         `json:"mobile"`                      // 电站联系人电话
	InstallerEmail              string      `json:"installerEmail"`              // 安装商邮箱
	InstallerMobile             Str         `json:"installerMobile"`             // 安装商电话
	BatteryPower                Num         `json:"batteryPower"`                // 电池充放电功率
	BatteryPowerStr             string      `json:"batteryPowerStr"`             // 电池充放电功率单位
	BatteryPowerPec             Num         `json:"batteryPowerPec"`             // 电池充放电功率百分比
	BatteryDischargeEnergy      Num         `json:"batteryDischargeEnergy"`      // 电池放电量
	BatteryDischargeEnergyStr   string      `json:"batteryDischargeEnergyStr"`   // 电池放电量单位
	BatteryChargeEnergy         Num         `json:"batteryChargeEnergy"`         // 电池充电量
	BatteryChargeEnergyStr      string      `json:"batteryChargeEnergyStr"`      // 电池充电量单位
	BatteryPercent              Num         `json:"batteryPercent"`              // 电池SOC
	Psum                        Num         `json:"psum"`                        // 电表功率
	PsumStr                     string      `json:"psumStr"`                     // 电表功率单位
	PsumPec                     Num         `json:"psumPec"`                     // 电表功率百分比
	GridPurchasedDayEnergy      Num         `json:"gridPurchasedDayEnergy"`      // 电表当日买电
	GridPurchasedDayEnergyStr   string      `json:"gridPurchasedDayEnergyStr"`   // 电表当日买电单位
	GridPurchasedMonthEnergy    Num         `json:"gridPurchasedMonthEnergy"`    // 电表当月买电
	GridPurchasedMonthEnergyStr string      `json:"gridPurchasedMonthEnergyStr"` // 电表当月买电单位
	GridPurchasedYearEnergy     Num         `json:"gridPurchasedYearEnergy"`     // 电表当年买电
	GridPurchasedYearEnergyStr  string      `json:"gridPurchasedYearEnergyStr"`  // 电表当年买电单位
	GridPurchasedTotalEnergy    Num         `json:"gridPurchasedTotalEnergy"`    // 电表总买电
	GridPurchasedTotalEnergyStr string      `json:"gridPurchasedTotalEnergyStr"` // 电表总买电单位
	GridSellDayEnergy           Num         `json:"gridSellDayEnergy"`           // 电表当日卖电
	GridSellDayEnergyStr        string      `json:"gridSellDayEnergyStr"`        // 电表当日卖电单位
	GridSellMonthEnergy         Num         `json:"gridSellMonthEnergy"`         // 电表当月卖电
	GridSellMonthEnergyStr      string      `json:"gridSellMonthEnergyStr"`      // 电表当月卖电单位
	GridSellYearEnergy          Num         `json:"gridSellYearEnergy"`          // 电表当年卖电
	GridSellYearEnergyStr       string      `json:"gridSellYearEnergyStr"`       // 电表当年卖电单位
	GridSellTotalEnergy         Num         `json:"gridSellTotalEnergy"`         // 电表总卖电
	GridSellTotalEnergyStr      string      `json:"gridSellTotalEnergyStr"`      // 电表总卖电单位
	FamilyLoadPower             Num         `json:"familyLoadPower"`             // 负载功率
	FamilyLoadPowerStr          string      `json:"familyLoadPowerStr"`          // 负载功率单位
	FamilyLoadPowerPec          Num         `json:"familyLoadPowerPec"`          // 负载功率百分比
	HomeLoadEnergy              Num         `json:"homeLoadEnergy"`              // 当日负载电量
	HomeLoadEnergyStr           string      `json:"homeLoadEnergyStr"`           // 当日负载电量单位
	InverterPower               Num         `json:"inverterPower"`               // 逆变器额定功率和
	NmiCode                     string      `json:"nmiCode"`                     // nmi码
	Country                     Int64       `json:"country"`                     // 国家id
	CountryStr                  string      `json:"countryStr"`                  // 国家名称
	Region                      Int64       `json:"region"`                      // 区域id
	RegionStr                   string      `json:"regionStr"`                   // 区域名称
	City                        Int64       `json:"city"`                        // 城市id
	CityStr                     string      `json:"cityStr"`                     // 城市名称
	County                      Int64       `json:"county"`                      // 区id
	CountyStr                   string      `json:"countyStr"`                   // 区名称
	TimeZone                    Num         `json:"timeZone"`                    // 时区
	TimeZoneName                string      `json:"timeZoneName"`                // 时区名称
	TimeZoneStr                 string      `json:"timeZoneStr"`                 // 时区格式化字符串
	TimeZoneID                  Int64       `json:"timeZoneId"`                  // 时区id
	Daylight                    int         `json:"daylight"`                    // 夏令时
	StationTypeNew              StationType `json:"stationTypeNew"`              // 电站类型
	AccessTime                  Int64       `json:"accessTime"`                  // 接入平台时间
	ConnectTime                 Int64       `json:"connectTime"`                 // 并网时间
	CreateDate                  Int64       `json:"createDate"`                  // 创建时间
	DaylightSwitch              int         `json:"daylightSwitch"`              // 夏令时开关，0关1开
	InverterBatteryCapacity     Num         `json:"inverterBatteryCapacity"`     // 电池容量(逆变器)
	BatteryCapacity             Num         `json:"batteryCapacity"`             // 电池容量(用户设置)
	BatteryCapacityEnergy       Num         `json:"batteryCapacityEnergy"`       // 电池容量(自研电池)
	TimeZoneStandardId          string      `json:"timeZoneStandardId"`          // 时区标准名称
}

// StationDetail 获取单个电站详情
func (c *SolisSDK) StationDetail(req StationDetailRequest) (*StationDetailResult, error) {
	if req.ID.Int() == 0 && strings.TrimSpace(req.NmiCode) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 id 和 nmiCode 不能同时为空", PathStationDetail)
	}
	var out StationDetailResult
	if err := c.do(PathStationDetail, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 4.3 获取多个电站详情 /v1/api/stationDetailList ----------

// StationDetailListRequest 多个电站详情请求
type StationDetailListRequest struct {
	PageRequest             // 分页参数
	IDList      []Int64     `json:"idList,omitempty"`      // 批量电站id，空为全部
	StationType StationType `json:"stationType,omitempty"` // 默认不传，访客传3
}

// StationDetailListItem 多个电站详情项
type StationDetailListItem struct {
	rawHolder
	ID                          Int64       `json:"id"`                          // 电站id
	StationName                 string      `json:"stationName"`                 // 电站名称
	Addr                        string      `json:"addr"`                        // 电站地址
	UserID                      Int64       `json:"userId"`                      // 业主Id
	UserName                    string      `json:"userName"`                    // 业主姓名
	UserMobile                  Str         `json:"userMobile"`                  // 业主电话
	UserEmail                   string      `json:"userEmail"`                   // 业主邮箱
	Capacity                    Num         `json:"capacity"`                    // 装机容量
	CapacityPec                 string      `json:"capacityPec"`                 // 装机容量倍率
	CapacityStr                 string      `json:"capacityStr"`                 // 装机容量单位
	DayEnergy                   Num         `json:"dayEnergy"`                   // 当日能量
	DayEnergyStr                string      `json:"dayEnergyStr"`                // 当日能量单位
	MonthEnergy                 Num         `json:"monthEnergy"`                 // 当月能量
	MonthEnergyStr              string      `json:"monthEnergyStr"`              // 当月能量单位
	YearEnergy                  Num         `json:"yearEnergy"`                  // 当年能量
	YearEnergyStr               string      `json:"yearEnergyStr"`               // 当年能量单位
	AllEnergy                   Num         `json:"allEnergy"`                   // 总能量
	AllEnergyStr                string      `json:"allEnergyStr"`                // 总能量单位
	FullHour                    Num         `json:"fullHour"`                    // 满发小时数
	PicName                     string      `json:"picName"`                     // 图片
	Power                       Num         `json:"power"`                       // 功率
	PowerPec                    Num         `json:"powerPec"`                    // 功率倍率
	PowerStr                    string      `json:"powerStr"`                    // 功率单位
	Latitude                    string      `json:"latitude"`                    // 纬度
	Longitude                   string      `json:"longitude"`                   // 经度
	Dip                         Num         `json:"dip"`                         // 倾角
	Azimuth                     Num         `json:"azimuth"`                     // 方位角
	Price                       Num         `json:"price"`                       // 每度电收益
	State                       DeviceState `json:"state"`                       // 电站状态
	DataTimestamp               Int64       `json:"dataTimestamp"`               // 更新时间
	Money                       string      `json:"money"`                       // 货币种类
	Brand                       string      `json:"brand"`                       // 品牌
	CondTxtN                    string      `json:"condTxtN"`                    // 晚上天气
	CondTxtD                    string      `json:"condTxtD"`                    // 白天天气
	TmpMax                      string      `json:"tmpMax"`                      // 最高温度
	TmpMin                      string      `json:"tmpMin"`                      // 最低温度
	TmpUnit                     string      `json:"tmpUnit"`                     // 温度单位
	Sr                          string      `json:"sr"`                          // 日出时间
	Ss                          string      `json:"ss"`                          // 日落时间
	Hum                         string      `json:"hum"`                         // 湿度
	WindDir                     string      `json:"windDir"`                     // 风向
	WindSpd                     string      `json:"windSpd"`                     // 风速
	WeatherUpdateDate           string      `json:"weatherUpdateDate"`           // 天气更新时间
	PowerStationNumTree         Num         `json:"powerStationNumTree"`         // 等效植树量
	PowerStationNumTreeUnit     string      `json:"powerStationNumTreeUnit"`     // 植物单位
	PowerStationAvoidedCo2      Num         `json:"powerStationAvoidedCo2"`      // 二氧化碳减排
	PowerStationAvoidedCo2Unit  string      `json:"powerStationAvoidedCo2Unit"`  // 减排单位
	Module                      int         `json:"module"`                      // 组件数量
	Mobile                      Str         `json:"mobile"`                      // 电站联系人电话
	InstallerEmail              string      `json:"installerEmail"`              // 安装商邮箱
	InstallerMobile             Str         `json:"installerMobile"`             // 安装商电话
	BatteryPower                Num         `json:"batteryPower"`                // 电池充放电功率
	BatteryPowerStr             string      `json:"batteryPowerStr"`             // 电池充放电功率单位
	BatteryPowerPec             Num         `json:"batteryPowerPec"`             // 电池充放电功率百分比
	BatteryDischargeEnergy      Num         `json:"batteryDischargeEnergy"`      // 电池放电量
	BatteryDischargeEnergyStr   string      `json:"batteryDischargeEnergyStr"`   // 电池放电量单位
	BatteryChargeEnergy         Num         `json:"batteryChargeEnergy"`         // 电池充电量
	BatteryChargeEnergyStr      string      `json:"batteryChargeEnergyStr"`      // 电池充电量单位
	BatteryPercent              Num         `json:"batteryPercent"`              // 电池SOC
	Psum                        Num         `json:"psum"`                        // 电表功率
	PsumStr                     string      `json:"psumStr"`                     // 电表功率单位
	PsumPec                     Num         `json:"psumPec"`                     // 电表功率百分比
	GridPurchasedDayEnergy      Num         `json:"gridPurchasedDayEnergy"`      // 电表当日买电
	GridPurchasedDayEnergyStr   string      `json:"gridPurchasedDayEnergyStr"`   // 电表当日买电单位
	GridPurchasedMonthEnergy    Num         `json:"gridPurchasedMonthEnergy"`    // 电表当月买电
	GridPurchasedMonthEnergyStr string      `json:"gridPurchasedMonthEnergyStr"` // 电表当月买电单位
	GridPurchasedYearEnergy     Num         `json:"gridPurchasedYearEnergy"`     // 电表当年买电
	GridPurchasedYearEnergyStr  string      `json:"gridPurchasedYearEnergyStr"`  // 电表当年买电单位
	GridPurchasedTotalEnergy    Num         `json:"gridPurchasedTotalEnergy"`    // 电表总买电
	GridPurchasedTotalEnergyStr string      `json:"gridPurchasedTotalEnergyStr"` // 电表总买电单位
	GridSellDayEnergy           Num         `json:"gridSellDayEnergy"`           // 电表当日卖电
	GridSellDayEnergyStr        string      `json:"gridSellDayEnergyStr"`        // 电表当日卖电单位
	GridSellMonthEnergy         Num         `json:"gridSellMonthEnergy"`         // 电表当月卖电
	GridSellMonthEnergyStr      string      `json:"gridSellMonthEnergyStr"`      // 电表当月卖电单位
	GridSellYearEnergy          Num         `json:"gridSellYearEnergy"`          // 电表当年卖电
	GridSellYearEnergyStr       string      `json:"gridSellYearEnergyStr"`       // 电表当年卖电单位
	GridSellTotalEnergy         Num         `json:"gridSellTotalEnergy"`         // 电表总卖电
	GridSellTotalEnergyStr      string      `json:"gridSellTotalEnergyStr"`      // 电表总卖电单位
	FamilyLoadPower             Num         `json:"familyLoadPower"`             // 负载功率
	FamilyLoadPowerStr          string      `json:"familyLoadPowerStr"`          // 负载功率单位
	FamilyLoadPowerPec          Num         `json:"familyLoadPowerPec"`          // 负载功率百分比
	HomeLoadEnergy              Num         `json:"homeLoadEnergy"`              // 当日负载电量
	HomeLoadEnergyStr           string      `json:"homeLoadEnergyStr"`           // 当日负载电量单位
	InverterPower               Num         `json:"inverterPower"`               // 逆变器额定功率和
	NmiCode                     string      `json:"nmiCode"`                     // nmi码
	Country                     Int64       `json:"country"`                     // 国家id
	CountryStr                  string      `json:"countryStr"`                  // 国家名称
	Region                      Int64       `json:"region"`                      // 区域id
	RegionStr                   string      `json:"regionStr"`                   // 区域名称
	City                        Int64       `json:"city"`                        // 城市id
	CityStr                     string      `json:"cityStr"`                     // 城市名称
	County                      Int64       `json:"county"`                      // 区id
	CountyStr                   string      `json:"countyStr"`                   // 区名称
	TimeZone                    Num         `json:"timeZone"`                    // 时区
	TimeZoneName                string      `json:"timeZoneName"`                // 时区名称
	TimeZoneStr                 string      `json:"timeZoneStr"`                 // 时区格式化字符串
	TimeZoneID                  Int64       `json:"timeZoneId"`                  // 时区id
	Daylight                    int         `json:"daylight"`                    // 夏令时
	StationTypeNew              StationType `json:"stationTypeNew"`              // 电站类型
	AccessTime                  Int64       `json:"accessTime"`                  // 接入平台时间
	ConnectTime                 Int64       `json:"connectTime"`                 // 并网时间
	CreateDate                  Int64       `json:"createDate"`                  // 创建时间
	DaylightSwitch              int         `json:"daylightSwitch"`              // 夏令时开关，0关1开
	InverterBatteryCapacity     Num         `json:"inverterBatteryCapacity"`     // 电池容量(逆变器)
	BatteryCapacity             Num         `json:"batteryCapacity"`             // 电池容量(用户设置)
	BatteryCapacityEnergy       Num         `json:"batteryCapacityEnergy"`       // 电池容量(自研电池)
	TimeZoneStandardId          string      `json:"timeZoneStandardId"`          // 时区标准名称
}

// StationDetailList 获取多个电站详情
func (c *SolisSDK) StationDetailList(req StationDetailListRequest) ([]StationDetailListItem, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out []StationDetailListItem
	if err := c.do(PathStationDetailList, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 4.4 获取多个电站某日的实时数据 /v1/api/stationDayEnergyList ----------

// StationDayEnergyListRequest 多电站某日实时数据请求
type StationDayEnergyListRequest struct {
	PageRequest         // 分页参数
	Time        string  `json:"time"`                 // 查询某天 yyyy-MM-dd
	StationIds  string  `json:"stationIds,omitempty"` // 电站id，英文逗号分隔
	IDList      []Int64 `json:"idList,omitempty"`     // 批量电站id，空为全部
}

// StationDayEnergyListItem 多电站某日实时数据项
type StationDayEnergyListItem struct {
	rawHolder
	ID                     Int64  `json:"id"`                     // 电站id
	Energy                 Num    `json:"energy"`                 // 发电量
	EnergyStr              string `json:"energyStr"`              // 发电量单位
	Date                   Int64  `json:"date"`                   // 时间戳
	Money                  Num    `json:"money"`                  // 收益
	MoneyStr               string `json:"moneyStr"`               // 收益单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电
	TotalRKwh              Num    `json:"totalRKwh"`              // 辐射量 kWh/㎡
}

// StationDayEnergyListResult 多电站某日实时数据结果
type StationDayEnergyListResult struct {
	rawHolder
	Page Page[StationDayEnergyListItem] `json:"page"` // 分页列表
}

// UnmarshalJSON 兼容 page 嵌套与平铺两种形态
func (r *StationDayEnergyListResult) UnmarshalJSON(data []byte) error {
	return unmarshalPaged(data, &r.Page)
}

// StationDayEnergyList 获取多个电站某日的实时数据
func (c *SolisSDK) StationDayEnergyList(req StationDayEnergyListRequest) (*StationDayEnergyListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	if strings.TrimSpace(req.Time) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 time 不能为空", PathStationDayEnergyList)
	}
	var out StationDayEnergyListResult
	if err := c.do(PathStationDayEnergyList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 4.5 获取多个电站的月数据 /v1/api/stationMonthEnergyList ----------

// StationMonthEnergyListRequest 多电站月数据请求
type StationMonthEnergyListRequest struct {
	PageRequest         // 分页参数
	Time        string  `json:"time"`                 // 查询某月 yyyy-MM
	StationIds  string  `json:"stationIds,omitempty"` // 电站id，英文逗号分隔
	NmiCode     string  `json:"nmiCode,omitempty"`    // nmi码
	IDList      []Int64 `json:"idList,omitempty"`     // 批量电站id，空为全部
}

// StationMonthEnergyListItem 多电站月数据项
type StationMonthEnergyListItem struct {
	rawHolder
	ID                     Int64  `json:"id"`                     // 电站id
	Energy                 Num    `json:"energy"`                 // 发电量
	Date                   Int64  `json:"date"`                   // 时间戳(+8时区)
	DateStr                string `json:"dateStr"`                // 更新时间字符串
	Money                  Num    `json:"money"`                  // 收益
	MoneyStr               string `json:"moneyStr"`               // 收益单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电量
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电量
	TotalRKwh              Num    `json:"totalRKwh"`              // 辐射量 kWh/㎡
}

// StationMonthEnergyListResult 多电站月数据结果
type StationMonthEnergyListResult struct {
	rawHolder
	Page Page[StationMonthEnergyListItem] `json:"page"` // 分页列表
}

// UnmarshalJSON 兼容 page 嵌套与平铺两种形态
func (r *StationMonthEnergyListResult) UnmarshalJSON(data []byte) error {
	return unmarshalPaged(data, &r.Page)
}

// StationMonthEnergyList 获取多个电站的月数据
func (c *SolisSDK) StationMonthEnergyList(req StationMonthEnergyListRequest) (*StationMonthEnergyListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	if strings.TrimSpace(req.Time) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 time 不能为空", PathStationMonthEnergyList)
	}
	var out StationMonthEnergyListResult
	if err := c.do(PathStationMonthEnergyList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 4.6 获取多个电站的年数据 /v1/api/stationYearEnergyList ----------

// StationYearEnergyListRequest 多电站年数据请求
type StationYearEnergyListRequest struct {
	PageRequest         // 分页参数
	Time        string  `json:"time"`                 // 查询某年 yyyy
	StationIds  string  `json:"stationIds,omitempty"` // 电站id，英文逗号分隔
	NmiCode     string  `json:"nmiCode,omitempty"`    // nmi码
	IDList      []Int64 `json:"idList,omitempty"`     // 批量电站id，空为全部
}

// StationYearEnergyListItem 多电站年数据项
type StationYearEnergyListItem struct {
	rawHolder
	ID                     Int64  `json:"id"`                     // 电站id
	Energy                 Num    `json:"energy"`                 // 发电量
	Year                   int    `json:"year"`                   // 年
	Money                  Num    `json:"money"`                  // 收益
	MoneyStr               string `json:"moneyStr"`               // 收益单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电量
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电量
	TotalRKwh              Num    `json:"totalRKwh"`              // 辐射量 kWh/㎡
}

// StationYearEnergyListResult 多电站年数据结果
type StationYearEnergyListResult struct {
	rawHolder
	Page Page[StationYearEnergyListItem] `json:"page"` // 分页列表
}

// UnmarshalJSON 兼容 page 嵌套与平铺两种形态
func (r *StationYearEnergyListResult) UnmarshalJSON(data []byte) error {
	return unmarshalPaged(data, &r.Page)
}

// StationYearEnergyList 获取多个电站的年数据
func (c *SolisSDK) StationYearEnergyList(req StationYearEnergyListRequest) (*StationYearEnergyListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	if strings.TrimSpace(req.Time) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 time 不能为空", PathStationYearEnergyList)
	}
	var out StationYearEnergyListResult
	if err := c.do(PathStationYearEnergyList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 4.7 获取单个电站某日的实时数据 /v1/api/stationDay ----------

// StationDayRequest 单电站某日实时数据请求
type StationDayRequest struct {
	ID       Int64  `json:"id,omitempty"`      // 电站id
	Money    string `json:"money"`             // 电站货币单位，例：CNY
	Time     string `json:"time"`              // 查询日期 yyyy-MM-dd
	TimeZone int    `json:"timeZone"`          // 电站时区，例：8
	NmiCode  string `json:"nmiCode,omitempty"` // nmi码
}

// StationDayItem 单电站某日实时数据项
type StationDayItem struct {
	rawHolder
	Power              Num    `json:"power"`              // 功率
	PowerStr           string `json:"powerStr"`           // 功率单位
	Time               Int64  `json:"time"`               // 时间戳
	TotalR             Num    `json:"totalR"`             // 瞬时辐射值
	BatteryCapacitySoc Num    `json:"batteryCapacitySoc"` // 电池SOC
}

// StationDay 获取单个电站某日的实时数据
func (c *SolisSDK) StationDay(req StationDayRequest) ([]StationDayItem, error) {
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 money 不能为空", PathStationDay)
	}
	if strings.TrimSpace(req.Time) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 time 不能为空", PathStationDay)
	}
	if req.TimeZone == 0 {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 timeZone 不能为空", PathStationDay)
	}
	if req.ID.Int() == 0 && strings.TrimSpace(req.NmiCode) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 id 和 nmiCode 不能同时为空", PathStationDay)
	}
	var out []StationDayItem
	if err := c.do(PathStationDay, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 4.8 获取单个电站某月的日数据 /v1/api/stationMonth ----------

// StationMonthRequest 单电站某月日数据请求
type StationMonthRequest struct {
	ID       Int64  `json:"id,omitempty"`      // 电站id
	Money    string `json:"money"`             // 电站货币单位，例：CNY
	Month    string `json:"month"`             // 查询月份 yyyy-MM
	TimeZone int    `json:"timeZone"`          // 电站时区，例：8
	NmiCode  string `json:"nmiCode,omitempty"` // nmi码
}

// StationMonthItem 单电站某月日数据项
type StationMonthItem struct {
	rawHolder
	Energy                 Num    `json:"energy"`                 // 发电量
	EnergyStr              string `json:"energyStr"`              // 发电量单位
	Date                   Int64  `json:"date"`                   // 时间戳
	Money                  Num    `json:"money"`                  // 收益
	MoneyStr               string `json:"moneyStr"`               // 收益单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电
	TotalRKwh              Num    `json:"totalRKwh"`              // 辐射量 kWh/㎡
}

// StationMonth 获取单个电站某月的日数据
func (c *SolisSDK) StationMonth(req StationMonthRequest) ([]StationMonthItem, error) {
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 money 不能为空", PathStationMonth)
	}
	if strings.TrimSpace(req.Month) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 month 不能为空", PathStationMonth)
	}
	if req.TimeZone == 0 {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 timeZone 不能为空", PathStationMonth)
	}
	if req.ID.Int() == 0 && strings.TrimSpace(req.NmiCode) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 id 和 nmiCode 不能同时为空", PathStationMonth)
	}
	var out []StationMonthItem
	if err := c.do(PathStationMonth, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 4.9 获取单个电站某年的月数据 /v1/api/stationYear ----------

// StationYearRequest 单电站某年月数据请求
type StationYearRequest struct {
	ID       Int64  `json:"id,omitempty"`      // 电站id
	Money    string `json:"money"`             // 电站货币单位，例：CNY
	Year     string `json:"year"`              // 查询年份 yyyy
	TimeZone int    `json:"timeZone"`          // 电站时区，例：8
	NmiCode  string `json:"nmiCode,omitempty"` // nmi码
}

// StationYearItem 单电站某年月数据项
type StationYearItem struct {
	rawHolder
	Energy                 Num    `json:"energy"`                 // 发电量
	EnergyStr              string `json:"energyStr"`              // 发电量单位
	Date                   Int64  `json:"date"`                   // 时间戳
	Money                  Num    `json:"money"`                  // 收益
	MoneyStr               string `json:"moneyStr"`               // 收益单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电
	TotalRKwh              Num    `json:"totalRKwh"`              // 辐射量 kWh/㎡
}

// StationYear 获取单个电站某年的月数据
func (c *SolisSDK) StationYear(req StationYearRequest) ([]StationYearItem, error) {
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 money 不能为空", PathStationYear)
	}
	if strings.TrimSpace(req.Year) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 year 不能为空", PathStationYear)
	}
	if req.TimeZone == 0 {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 timeZone 不能为空", PathStationYear)
	}
	if req.ID.Int() == 0 && strings.TrimSpace(req.NmiCode) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 id 和 nmiCode 不能同时为空", PathStationYear)
	}
	var out []StationYearItem
	if err := c.do(PathStationYear, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 4.10 获取某个电站的年数据 /v1/api/stationAll ----------

// StationAllRequest 单电站年数据请求
type StationAllRequest struct {
	ID       Int64  `json:"id,omitempty"`      // 电站id
	Money    string `json:"money"`             // 电站货币单位，例：CNY
	TimeZone int    `json:"timeZone"`          // 电站时区，例：8
	NmiCode  string `json:"nmiCode,omitempty"` // nmi码
}

// StationAllItem 单电站年数据项
type StationAllItem struct {
	rawHolder
	Energy                 Num    `json:"energy"`                 // 发电量
	EnergyStr              string `json:"energyStr"`              // 发电量单位
	Date                   Int64  `json:"date"`                   // 时间戳
	Money                  Num    `json:"money"`                  // 收益
	MoneyStr               string `json:"moneyStr"`               // 收益单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电
	TotalRKwh              Num    `json:"totalRKwh"`              // 辐射量 kWh/㎡
}

// StationAll 获取某个电站的年数据
func (c *SolisSDK) StationAll(req StationAllRequest) ([]StationAllItem, error) {
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 money 不能为空", PathStationAll)
	}
	if req.TimeZone == 0 {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 timeZone 不能为空", PathStationAll)
	}
	if req.ID.Int() == 0 && strings.TrimSpace(req.NmiCode) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 id 和 nmiCode 不能同时为空", PathStationAll)
	}
	var out []StationAllItem
	if err := c.do(PathStationAll, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
