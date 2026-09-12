package solarman

import (
	"fmt"
	"strings"
)

// 本文件实现 SolarMan 开放平台电站接口（文档 4.1 ~ 4.16）。
// 平台响应把业务字段与 code/msg/success/requestId 平铺在同一层，
// 因此每个结果结构体都匿名嵌入 Response 与 rawHolder。

// ---------------------------------------------------------------------------
// 4.1 查询电站基础信息 /station/v1.0/base
// ---------------------------------------------------------------------------

// StationRegion 电站行政区划
type StationRegion struct {
	rawHolder
	NationID int    `json:"nationId"`         // 国家ID
	Level1   int    `json:"level1,omitempty"` // 行政区1
	Level2   int    `json:"level2,omitempty"` // 行政区2
	Level3   int    `json:"level3,omitempty"` // 行政区3
	Level4   int    `json:"level4,omitempty"` // 行政区4
	Level5   int    `json:"level5,omitempty"` // 行政区5
	Timezone string `json:"timezone"`         // 时区，如 PRC
}

// StationImage 电站图片
type StationImage struct {
	rawHolder
	ID          Int64  `json:"id"`          // 图片ID
	Name        string `json:"name"`        // 图片名称
	Description string `json:"description"` // 图片描述
	URL         string `json:"url"`         // 图片url
}

// StationBaseRequest 4.1 请求参数
type StationBaseRequest struct {
	StationID Int64 `json:"stationId"`           // 电站ID，必填
	CompanyID Int64 `json:"companyId,omitempty"` // 商家ID，商家登录必传
	UserID    Int64 `json:"userId,omitempty"`    // 账号ID
}

// StationBaseResult 4.1 响应结果
type StationBaseResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	ID                       Int64                   `json:"id"`                       // 电站ID
	Name                     string                  `json:"name"`                     // 电站名称
	LocationLat              Num                     `json:"locationLat"`              // 纬度
	LocationLng              Num                     `json:"locationLng"`              // 经度
	LocationAddress          string                  `json:"locationAddress"`          // 详细地址
	Region                   StationRegion           `json:"region"`                   // 行政区划
	Type                     StationType             `json:"type"`                     // 电站类型
	GridInterconnectionType  GridInterconnectionType `json:"gridInterconnectionType"`  // 并网类型
	InstalledCapacity        Num                     `json:"installedCapacity"`        // 装机容量
	InstallationAzimuthAngle Num                     `json:"installationAzimuthAngle"` // 方位角
	InstallationTiltAngle    Num                     `json:"installationTiltAngle"`    // 倾角
	StartOperatingTime       Num                     `json:"startOperatingTime"`       // 开始运行时间，UNIX 秒
	Operating                bool                    `json:"operating"`                // 并网状态
	TotalRunningDayCount     Num                     `json:"totalRunningDayCount"`     // 累计运行天数
	Currency                 string                  `json:"currency"`                 // 货币单位
	OwnerName                string                  `json:"ownerName"`                // 业主姓名
	OwnerCompany             string                  `json:"ownerCompany"`             // 业主工作单位
	ContactPhone             string                  `json:"contactPhone"`             // 联系电话
	MergeElectricPrice       Num                     `json:"mergeElectricPrice"`       // 上网电价（元/kWh）
	ConstructionCost         Num                     `json:"constructionCost"`         // 总成本（元）
	TheoryConsumeProportion  Num                     `json:"theoryConsumeProportion"`  // 计划自发自用率
	StationImage             string                  `json:"stationImage"`             // 电站封面
	StationImages            []StationImage          `json:"stationImages"`            // 电站图片列表
	CreatedDate              Num                     `json:"createdDate"`              // 创建时间，UNIX 秒
}

// StationBase 4.1 查询电站基础信息
func (sdk *SolarmanSDK) StationBase(req StationBaseRequest) (*StationBaseResult, error) {
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}

	var out StationBaseResult
	if err := sdk.do(PathStationBase, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.2 获取电站设备列表 /station/v1.0/device
// ---------------------------------------------------------------------------

// StationDeviceItem 电站设备项
type StationDeviceItem struct {
	rawHolder
	DeviceSN       string     `json:"deviceSn"`       // 设备SN
	DeviceID       Int64      `json:"deviceId"`       // 平台内设备ID
	DeviceType     DeviceType `json:"deviceType"`     // 设备类型
	ConnectStatus  int        `json:"connectStatus"`  // 状态：0离线1在线2报警
	CollectionTime Num        `json:"collectionTime"` // 最后数据更新时间，UNIX 秒
}

// StationDeviceListRequest 4.2 请求参数
type StationDeviceListRequest struct {
	PageRequest            // 分页参数
	StationID   Int64      `json:"stationId"`            // 电站ID，必填
	DeviceType  DeviceType `json:"deviceType,omitempty"` // 设备类型，不传取全部
}

// StationDeviceListResult 4.2 响应结果
type StationDeviceListResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	Total           int                 `json:"total"`           // 设备总数
	DeviceListItems []StationDeviceItem `json:"deviceListItems"` // 设备列表
}

// StationDeviceList 4.2 获取电站设备列表
func (sdk *SolarmanSDK) StationDeviceList(req StationDeviceListRequest) (*StationDeviceListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}

	var out StationDeviceListResult
	if err := sdk.do(PathStationDevice, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.3 获取电站历史数据 /station/v1.0/history
// ---------------------------------------------------------------------------

// StationHistoryItem 电站历史数据项
type StationHistoryItem struct {
	rawHolder
	GenerationPower       Num    `json:"generationPower"`       // 发电功率
	UsePower              Num    `json:"usePower"`              // 用电功率
	GridPower             Num    `json:"gridPower"`             // 并网功率
	PurchasePower         Num    `json:"purchasePower"`         // 购电功率
	WirePower             Num    `json:"wirePower"`             // 电网功率
	ChargePower           Num    `json:"chargePower"`           // 充电功率
	DischargePower        Num    `json:"dischargePower"`        // 放电功率
	BatteryPower          Num    `json:"batteryPower"`          // 电池功率
	BatterySoc            Num    `json:"batterySoc"`            // 电池容量
	IrradiateIntensity    Num    `json:"irradiateIntensity"`    // 辐照强度
	GenerationValue       Num    `json:"generationValue"`       // 发电量
	GenerationRatio       Num    `json:"generationRatio"`       // 自发自用比率
	GridRatio             Num    `json:"gridRatio"`             // 发电并网比例
	ChargeRatio           Num    `json:"chargeRatio"`           // 自发充电比例
	UseValue              Num    `json:"useValue"`              // 用电量
	UseRatio              Num    `json:"useRatio"`              // 用电来自发电比例
	BuyRatio              Num    `json:"buyRatio"`              // 电网购电比例
	UseDischargeRatio     Num    `json:"useDischargeRatio"`     // 用电放电比例
	GridValue             Num    `json:"gridValue"`             // 并网量
	BuyValue              Num    `json:"buyValue"`              // 购电量
	ChargeValue           Num    `json:"chargeValue"`           // 充电量
	DischargeValue        Num    `json:"dischargeValue"`        // 放电量
	FullPowerHours        Num    `json:"fullPowerHours"`        // 满发小时数
	Irradiate             Num    `json:"irradiate"`             // 气象站辐照量
	TheoreticalGeneration Num    `json:"theoreticalGeneration"` // 理论发电量
	PR                    Num    `json:"pr"`                    // PR
	CPR                   Num    `json:"cpr"`                   // cpr
	DateTime              string `json:"dateTime"`              // 数据时间
	Year                  int    `json:"year"`                  // 数据处理年份
	Month                 int    `json:"month"`                 // 数据处理月份
	Day                   int    `json:"day"`                   // 数据处理日期
}

// StationHistoryRequest 4.3 请求参数
type StationHistoryRequest struct {
	StationID Int64  `json:"stationId"`         // 电站ID，必填
	TimeType  int    `json:"timeType"`          // 时间类型：1帧2日3月4年
	StartTime string `json:"startTime"`         // 开始时间，必填
	EndTime   string `json:"endTime,omitempty"` // 结束时间
}

// StationHistoryResult 4.3 响应结果
type StationHistoryResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	Total            int                  `json:"total"`            // 结果总数
	StationDataItems []StationHistoryItem `json:"stationDataItems"` // 结果列表
}

// StationHistory 4.3 获取电站历史数据
func (sdk *SolarmanSDK) StationHistory(req StationHistoryRequest) (*StationHistoryResult, error) {
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}
	if req.TimeType < 1 || req.TimeType > 4 {
		return nil, fmt.Errorf("[Solarman] 时间类型(timeType) 必填，取值为 1~4")
	}
	if strings.TrimSpace(req.StartTime) == "" {
		return nil, fmt.Errorf("[Solarman] 开始时间(startTime) 必填")
	}

	var out StationHistoryResult
	if err := sdk.do(PathStationHistory, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.4 获取账号下电站列表 /station/v1.0/list
// ---------------------------------------------------------------------------

// StationListItem 账号下电站项
type StationListItem struct {
	rawHolder
	ID                      Int64                   `json:"id"`                      // 电站ID
	Name                    string                  `json:"name"`                    // 电站名称
	LocationLat             Num                     `json:"locationLat"`             // 纬度
	LocationLng             Num                     `json:"locationLng"`             // 经度
	LocationAddress         string                  `json:"locationAddress"`         // 详细地址
	RegionNationID          int                     `json:"regionNationId"`          // 国家
	RegionLevel1            int                     `json:"regionLevel1"`            // 行政区1
	RegionLevel2            int                     `json:"regionLevel2"`            // 行政区2
	RegionLevel3            int                     `json:"regionLevel3"`            // 行政区3
	RegionLevel4            int                     `json:"regionLevel4"`            // 行政区4
	RegionLevel5            int                     `json:"regionLevel5"`            // 行政区5
	RegionTimezone          string                  `json:"regionTimezone"`          // 时区
	Type                    StationType             `json:"type"`                    // 电站类型
	GridInterconnectionType GridInterconnectionType `json:"gridInterconnectionType"` // 并网类型
	InstalledCapacity       Num                     `json:"installedCapacity"`       // 装机容量
	StartOperatingTime      Num                     `json:"startOperatingTime"`      // 开始运行时间，UNIX 秒
	CreatedDate             Num                     `json:"createdDate"`             // 创建时间，UNIX 秒
	LastUpdateTime          Num                     `json:"lastUpdateTime"`          // 最后更新时间，UNIX 秒
	BatterySoc              Num                     `json:"batterySoc"`              // 电池容量
	NetworkStatus           string                  `json:"networkStatus"`           // 通信状态
	GenerationPower         Num                     `json:"generationPower"`         // 发电功率
	StationImage            string                  `json:"stationImage"`            // 电站封面
	ContactPhone            string                  `json:"contactPhone"`            // 联系电话
	OwnerName               string                  `json:"ownerName"`               // 业主姓名
	ConstructionCost        Num                     `json:"constructionCost"`        // 总成本（元）
}

// StationListRequest 4.4 请求参数
type StationListRequest struct {
	PageRequest // 分页参数
}

// StationListResult 4.4 响应结果
type StationListResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	Total       int               `json:"total"`       // 结果总数
	StationList []StationListItem `json:"stationList"` // 电站列表
}

// StationList 4.4 获取账号下电站列表
func (sdk *SolarmanSDK) StationList(req StationListRequest) (*StationListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()

	var out StationListResult
	if err := sdk.do(PathStationList, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.5 获取电站实时数据 /station/v1.0/realTime
// ---------------------------------------------------------------------------

// StationRealTimeRequest 4.5 请求参数
type StationRealTimeRequest struct {
	StationID Int64 `json:"stationId"` // 电站ID，必填
}

// StationRealTimeResult 4.5 响应结果
type StationRealTimeResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	GenerationPower    Num `json:"generationPower"`    // 发电功率
	UsePower           Num `json:"usePower"`           // 用电功率
	GridPower          Num `json:"gridPower"`          // 并网功率
	PurchasePower      Num `json:"purchasePower"`      // 购电功率
	WirePower          Num `json:"wirePower"`          // 电网功率
	ChargePower        Num `json:"chargePower"`        // 充电功率
	DischargePower     Num `json:"dischargePower"`     // 放电功率
	BatteryPower       Num `json:"batteryPower"`       // 电池功率
	BatterySoc         Num `json:"batterySoc"`         // 电池容量
	IrradiateIntensity Num `json:"irradiateIntensity"` // 辐照强度
	LastUpdateTime     Num `json:"lastUpdateTime"`     // 最新更新时间，UNIX 秒
}

// StationRealTime 4.5 获取电站实时数据
func (sdk *SolarmanSDK) StationRealTime(req StationRealTimeRequest) (*StationRealTimeResult, error) {
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}

	var out StationRealTimeResult
	if err := sdk.do(PathStationRealTime, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.6 获取电站操作权限 /station/v1.0/role
// ---------------------------------------------------------------------------

// StationRoleRequest 4.6 请求参数
type StationRoleRequest struct {
	StationID Int64 `json:"stationId"` // 电站ID，必填
}

// StationRoleResult 4.6 响应结果，1 有权限 0 无权限
type StationRoleResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	ViewPlantAlert    int `json:"viewPlantAlert"`    // 查看电站报警权限
	ViewPlantInfo     int `json:"viewPlantInfo"`     // 查看基本信息权限
	ViewPlantDevice   int `json:"viewPlantDevice"`   // 查看下设备权限
	EditPlant         int `json:"editPlant"`         // 编辑电站权限
	DeletePlant       int `json:"deletePlant"`       // 删除电站权限
	SetPlant          int `json:"setPlant"`          // 设置电站权限
	AddPlantDevice    int `json:"addPlantDevice"`    // 添加设备权限
	DeletePlantDevice int `json:"deletePlantDevice"` // 删除设备权限
}

// StationRole 4.6 获取电站操作权限
func (sdk *SolarmanSDK) StationRole(req StationRoleRequest) (*StationRoleResult, error) {
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}

	var out StationRoleResult
	if err := sdk.do(PathStationRole, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.7 获取电站报警列表 /station/v1.0/alert
// ---------------------------------------------------------------------------

// StationAlertItem 电站报警项
type StationAlertItem struct {
	rawHolder
	DeviceSN   string     `json:"deviceSn"`   // 设备唯一标识
	DeviceID   Int64      `json:"deviceId"`   // 平台内设备唯一标识
	DeviceType DeviceType `json:"deviceType"` // 设备类型
	RuleID     Int64      `json:"ruleId"`     // 报警的id
	ShowName   string     `json:"showName"`   // PAAS 定义的报警名称
	Addr       string     `json:"addr"`       // 报警协议的名称
	Code       string     `json:"code"`       // 报警协议的代码
	Level      int        `json:"level"`      // 报警等级：提示/警告/故障
	Influence  int        `json:"influence"`  // 影响范围
	AlertTime  Num        `json:"alertTime"`  // 最后触发时间，UNIX 秒
}

// StationAlertRequest 4.7 请求参数
type StationAlertRequest struct {
	PageRequest        // 分页参数
	StationID   Int64  `json:"stationId"` // 电站ID，必填
	StartTime   string `json:"startTime"` // 开始时间，必填 yyyy-MM-dd
	EndTime     string `json:"endTime"`   // 结束时间，必填 yyyy-MM-dd
}

// StationAlertResult 4.7 响应结果
type StationAlertResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	Total             int                `json:"total"`             // 当前报警总数
	StationAlertItems []StationAlertItem `json:"stationAlertItems"` // 报警结果列表
}

// StationAlert 4.7 获取电站报警列表
func (sdk *SolarmanSDK) StationAlert(req StationAlertRequest) (*StationAlertResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}
	if strings.TrimSpace(req.StartTime) == "" {
		return nil, fmt.Errorf("[Solarman] 开始时间(startTime) 必填")
	}
	if strings.TrimSpace(req.EndTime) == "" {
		return nil, fmt.Errorf("[Solarman] 结束时间(endTime) 必填")
	}

	var out StationAlertResult
	if err := sdk.do(PathStationAlert, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.8 创建电站 /station/v1.0/create
// ---------------------------------------------------------------------------

// StationCreateRequest 4.8 请求参数
type StationCreateRequest struct {
	Name                     string                  `json:"name"`                               // 电站名称，必填
	Type                     StationType             `json:"type"`                               // 电站类型，必填
	GridInterconnectionType  GridInterconnectionType `json:"gridInterconnectionType"`            // 并网类型，必填
	Currency                 string                  `json:"currency"`                           // 货币单位，必填
	InstalledCapacity        Num                     `json:"installedCapacity"`                  // 装机容量，必填
	LocationAddress          string                  `json:"locationAddress"`                    // 电站地址，必填
	LocationLat              Num                     `json:"locationLat"`                        // 纬度，必填
	LocationLng              Num                     `json:"locationLng"`                        // 经度，必填
	Region                   StationRegion           `json:"region"`                             // 电站区域信息，必填
	InstallationAzimuthAngle Num                     `json:"installationAzimuthAngle,omitempty"` // 方位角
	InstallationTiltAngle    Num                     `json:"installationTiltAngle,omitempty"`    // 倾角
	MergeElectricPrice       Num                     `json:"mergeElectricPrice,omitempty"`       // 度电收益（元/kWh）
	ConstructionCost         Num                     `json:"constructionCost,omitempty"`         // 总成本（元）
	StartOperatingTime       Num                     `json:"startOperatingTime,omitempty"`       // 开始运行时间，UNIX 秒
	OwnerName                string                  `json:"ownerName,omitempty"`                // 业主姓名
	OwnerCompany             string                  `json:"ownerCompany,omitempty"`             // 业主工作单位
	ContactPhone             string                  `json:"contactPhone,omitempty"`             // 联系电话
	StationImage             string                  `json:"stationImage,omitempty"`             // 电站封面
}

// StationCreateResult 4.8 响应结果
type StationCreateResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	ID Int64 `json:"id"` // 新电站ID
}

// StationCreate 4.8 创建电站
func (sdk *SolarmanSDK) StationCreate(req StationCreateRequest) (*StationCreateResult, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("[Solarman] 电站名称(name) 必填")
	}
	if strings.TrimSpace(string(req.Type)) == "" {
		return nil, fmt.Errorf("[Solarman] 电站类型(type) 必填")
	}
	if strings.TrimSpace(string(req.GridInterconnectionType)) == "" {
		return nil, fmt.Errorf("[Solarman] 并网类型(gridInterconnectionType) 必填")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return nil, fmt.Errorf("[Solarman] 货币单位(currency) 必填")
	}
	if req.InstalledCapacity == 0 {
		return nil, fmt.Errorf("[Solarman] 装机容量(installedCapacity) 必填")
	}
	if strings.TrimSpace(req.LocationAddress) == "" {
		return nil, fmt.Errorf("[Solarman] 电站地址(locationAddress) 必填")
	}
	if req.LocationLat == 0 {
		return nil, fmt.Errorf("[Solarman] 电站纬度(locationLat) 必填")
	}
	if req.LocationLng == 0 {
		return nil, fmt.Errorf("[Solarman] 电站经度(locationLng) 必填")
	}
	if req.Region.NationID == 0 {
		return nil, fmt.Errorf("[Solarman] 国家ID(region.nationId) 必填")
	}
	if strings.TrimSpace(req.Region.Timezone) == "" {
		return nil, fmt.Errorf("[Solarman] 时区(region.timezone) 必填")
	}

	var out StationCreateResult
	if err := sdk.do(PathStationCreate, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.9 修改电站 /station/v1.0/update
// ---------------------------------------------------------------------------

// StationUpdateRequest 4.9 请求参数
type StationUpdateRequest struct {
	StationID                Int64                   `json:"stationId"`                          // 电站ID，必填
	Name                     string                  `json:"name"`                               // 电站名称，必填
	Type                     StationType             `json:"type"`                               // 电站类型，必填
	GridInterconnectionType  GridInterconnectionType `json:"gridInterconnectionType"`            // 并网类型，必填
	Currency                 string                  `json:"currency"`                           // 货币单位，必填
	InstalledCapacity        Num                     `json:"installedCapacity"`                  // 装机容量，必填
	LocationAddress          string                  `json:"locationAddress"`                    // 电站地址，必填
	LocationLat              Num                     `json:"locationLat"`                        // 纬度，必填
	LocationLng              Num                     `json:"locationLng"`                        // 经度，必填
	Region                   StationRegion           `json:"region"`                             // 电站区域信息，必填
	InstallationAzimuthAngle Num                     `json:"installationAzimuthAngle,omitempty"` // 方位角
	InstallationTiltAngle    Num                     `json:"installationTiltAngle,omitempty"`    // 倾角
	MergeElectricPrice       Num                     `json:"mergeElectricPrice,omitempty"`       // 度电收益（元/kWh）
	ConstructionCost         Num                     `json:"constructionCost,omitempty"`         // 总成本（元）
	StartOperatingTime       Num                     `json:"startOperatingTime,omitempty"`       // 开始运行时间，UNIX 秒
	OwnerName                string                  `json:"ownerName,omitempty"`                // 业主姓名
	OwnerCompany             string                  `json:"ownerCompany,omitempty"`             // 业主工作单位
	ContactPhone             string                  `json:"contactPhone,omitempty"`             // 联系电话
	StationImage             string                  `json:"stationImage,omitempty"`             // 电站封面
}

// StationUpdateResult 4.9 响应结果
type StationUpdateResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// StationUpdate 4.9 修改电站
func (sdk *SolarmanSDK) StationUpdate(req StationUpdateRequest) (*StationUpdateResult, error) {
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("[Solarman] 电站名称(name) 必填")
	}
	if strings.TrimSpace(string(req.Type)) == "" {
		return nil, fmt.Errorf("[Solarman] 电站类型(type) 必填")
	}
	if strings.TrimSpace(string(req.GridInterconnectionType)) == "" {
		return nil, fmt.Errorf("[Solarman] 并网类型(gridInterconnectionType) 必填")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return nil, fmt.Errorf("[Solarman] 货币单位(currency) 必填")
	}
	if req.InstalledCapacity == 0 {
		return nil, fmt.Errorf("[Solarman] 装机容量(installedCapacity) 必填")
	}
	if strings.TrimSpace(req.LocationAddress) == "" {
		return nil, fmt.Errorf("[Solarman] 电站地址(locationAddress) 必填")
	}
	if req.LocationLat == 0 {
		return nil, fmt.Errorf("[Solarman] 电站纬度(locationLat) 必填")
	}
	if req.LocationLng == 0 {
		return nil, fmt.Errorf("[Solarman] 电站经度(locationLng) 必填")
	}
	if req.Region.NationID == 0 {
		return nil, fmt.Errorf("[Solarman] 国家ID(region.nationId) 必填")
	}
	if strings.TrimSpace(req.Region.Timezone) == "" {
		return nil, fmt.Errorf("[Solarman] 时区(region.timezone) 必填")
	}

	var out StationUpdateResult
	if err := sdk.do(PathStationUpdate, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.10 删除电站 /station/v1.0/delete
// ---------------------------------------------------------------------------

// StationDeleteRequest 4.10 请求参数
type StationDeleteRequest struct {
	StationID Int64 `json:"stationId"` // 电站ID，必填
}

// StationDeleteResult 4.10 响应结果，仅公共字段
type StationDeleteResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// StationDelete 4.10 删除电站
func (sdk *SolarmanSDK) StationDelete(req StationDeleteRequest) (*StationDeleteResult, error) {
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}

	var out StationDeleteResult
	if err := sdk.do(PathStationDelete, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.11 累计发电量的计算方法设置 /station/v1.0/metering
// ---------------------------------------------------------------------------

// StationMeteringRequest 4.11 请求参数
type StationMeteringRequest struct {
	StationID           Int64               `json:"stationId"`           // 电站ID，必填
	TotalProductionType TotalProductionType `json:"totalProductionType"` // 计算方法1或2，必填
}

// StationMeteringResult 4.11 响应结果，仅公共字段
type StationMeteringResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// StationMetering 4.11 累计发电量的计算方法设置
func (sdk *SolarmanSDK) StationMetering(req StationMeteringRequest) (*StationMeteringResult, error) {
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}
	if req.TotalProductionType != TotalByDeviceTotal && req.TotalProductionType != TotalByDailySum {
		return nil, fmt.Errorf("[Solarman] 累计发电量计算方法(totalProductionType) 必填，取值为 1 或 2")
	}

	var out StationMeteringResult
	if err := sdk.do(PathStationMetering, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.12 设置偏移量 /station/v1.0/offset
// ---------------------------------------------------------------------------

// StationOffsetRequest 4.12 请求参数
type StationOffsetRequest struct {
	StationID   Int64  `json:"stationId"`      // 电站ID，必填
	OffsetType  int    `json:"offsetType"`     // 1累计偏移量2日偏移量，必填
	OffsetValue Num    `json:"offsetValue"`    // 偏移数值，正增负减，必填
	Date        string `json:"date,omitempty"` // 日偏移量日期 yyyy-MM-dd
}

// StationOffsetResult 4.12 响应结果，仅公共字段
type StationOffsetResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// StationOffset 4.12 设置偏移量
func (sdk *SolarmanSDK) StationOffset(req StationOffsetRequest) (*StationOffsetResult, error) {
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}
	if req.OffsetType != 1 && req.OffsetType != 2 {
		return nil, fmt.Errorf("[Solarman] 偏移量设置类型(offsetType) 必填，取值为 1 或 2")
	}
	if req.OffsetValue == 0 {
		return nil, fmt.Errorf("[Solarman] 偏移量数值(offsetValue) 必填")
	}
	if req.OffsetType == 2 && strings.TrimSpace(req.Date) == "" {
		return nil, fmt.Errorf("[Solarman] 日偏移量日期(date) 在 offsetType=2 时必填")
	}

	var out StationOffsetResult
	if err := sdk.do(PathStationOffset, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.16 电站报警列表 V2.0 /station/v1.0/alertV2
// ---------------------------------------------------------------------------

// StationAlertV2Item 电站报警项 V2
type StationAlertV2Item struct {
	rawHolder
	ID             Int64      `json:"id"`             // 报警记录唯一ID
	DeviceSN       string     `json:"deviceSn"`       // 设备唯一标识
	DeviceID       string     `json:"deviceId"`       // 平台内设备唯一标识
	DeviceType     DeviceType `json:"deviceType"`     // 设备类型
	RuleID         Int64      `json:"ruleId"`         // 报警的id
	AlertName      string     `json:"alertName"`      // PAAS 定义的报警名称
	Addr           string     `json:"addr"`           // 报警协议的名称
	Code           string     `json:"code"`           // 报警协议的代码
	Level          int        `json:"level"`          // 报警等级：提示/警告/故障
	Status         int        `json:"status"`         // 状态：1发生中0已恢复
	AlertStartTime Num        `json:"alertStartTime"` // 报警开始时间，UNIX 秒
	AlertEndTime   Num        `json:"alertEndTime"`   // 报警结束时间，UNIX 秒
}

// StationAlertV2Request 4.16 请求参数
type StationAlertV2Request struct {
	PageRequest        // 分页参数
	StationID   Int64  `json:"stationId"` // 电站ID，必填
	StartTime   string `json:"startTime"` // 开始时间，必填 yyyy-MM-dd
	EndTime     string `json:"endTime"`   // 结束时间，必填 yyyy-MM-dd
}

// StationAlertV2Result 4.16 响应结果
type StationAlertV2Result struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底

	Total             int                  `json:"total"`             // 当前报警总数
	StationAlertItems []StationAlertV2Item `json:"stationAlertItems"` // 报警结果列表
}

// StationAlertV2 4.16 电站报警列表 V2.0
func (sdk *SolarmanSDK) StationAlertV2(req StationAlertV2Request) (*StationAlertV2Result, error) {
	req.PageRequest = req.PageRequest.Normalize()
	if req.StationID == 0 {
		return nil, fmt.Errorf("[Solarman] 电站ID(stationId) 必填")
	}
	if strings.TrimSpace(req.StartTime) == "" {
		return nil, fmt.Errorf("[Solarman] 开始时间(startTime) 必填")
	}
	if strings.TrimSpace(req.EndTime) == "" {
		return nil, fmt.Errorf("[Solarman] 结束时间(endTime) 必填")
	}

	var out StationAlertV2Result
	if err := sdk.do(PathStationAlertV2, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}
