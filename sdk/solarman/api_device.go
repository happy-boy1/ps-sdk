package solarman

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// 3.1 设备报警明细 POST /device/v1.0/alertDetail
// ---------------------------------------------------------------------------

// AlertDetailRequest 3.1 设备报警明细请求
type AlertDetailRequest struct {
	AlertID  Int64  `json:"alertId"`            // 报警的 ID，必填
	DeviceID Int64  `json:"deviceId,omitempty"` // 设备在平台内的 ID，可选
	DeviceSN string `json:"deviceSn"`           // 设备的唯一标识，必填
}

// AlertDetailResult 3.1 设备报警明细返回
type AlertDetailResult struct {
	Response               // 公共响应字段
	rawHolder              // 原始数据兜底
	Addr        string     `json:"addr"`        // 报警协议的名称
	AlertCode   string     `json:"alertCode"`   // 报警协议的代码
	AlertID     Int64      `json:"alertId"`     // 报警的 ID
	AlertName   string     `json:"alertName"`   // 平台定义的报警名称
	AlertTime   Int64      `json:"alertTime"`   // 最后一次触发时间戳
	Description string     `json:"description"` // 报警描述
	DeviceSN    string     `json:"deviceSn"`    // 设备的唯一标识
	DeviceType  DeviceType `json:"deviceType"`  // 设备类型
	Influence   int        `json:"influence"`   // 影响范围 0-3
	Level       int        `json:"level"`       // 报警等级 0-2
	Reason      string     `json:"reason"`      // 报警原因
	Solution    string     `json:"solution"`    // 报警解决方案
}

// AlertDetail 3.1 设备报警明细，仅返回设备原始报警数据
func (sdk *SolarmanSDK) AlertDetail(req AlertDetailRequest) (*AlertDetailResult, error) {
	if req.AlertID <= 0 {
		return nil, fmt.Errorf("[Solarman] AlertDetail: alertId 必填")
	}
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] AlertDetail: deviceSn 必填")
	}

	var out AlertDetailResult
	if err := sdk.do(PathAlertDetail, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 3.2 设备报警列表 POST /device/v1.0/alertList
// ---------------------------------------------------------------------------

// AlertListRequest 3.2 设备报警列表请求
type AlertListRequest struct {
	PageRequest           // 分页参数
	DeviceID       Int64  `json:"deviceId,omitempty"` // 设备在平台内的 ID，可选
	DeviceSN       string `json:"deviceSn"`           // 设备的唯一标识，必填
	EndTimestamp   Int64  `json:"endTimestamp"`       // 结束时间戳，必填
	StartTimestamp Int64  `json:"startTimestamp"`     // 开始时间戳，必填
}

// AlertListItem 3.2 单条设备报警
type AlertListItem struct {
	rawHolder
	Addr      string `json:"addr"`      // 报警协议的名称
	AlertID   Int64  `json:"alertId"`   // 报警的 ID
	AlertName string `json:"alertName"` // 平台定义的报警名称
	AlertTime Int64  `json:"alertTime"` // 最后一次触发时间戳
	Code      string `json:"code"`      // 报警协议的代码
	Influence int    `json:"influence"` // 影响范围 0-3
	Level     int    `json:"level"`     // 报警等级 0-2
}

// AlertListResult 3.2 设备报警列表返回
type AlertListResult struct {
	Response                   // 公共响应字段
	rawHolder                  // 原始数据兜底
	AlertList  []AlertListItem `json:"alertList"`  // 报警列表
	DeviceID   Int64           `json:"deviceId"`   // 设备在平台内的 ID
	DeviceSN   string          `json:"deviceSn"`   // 设备的唯一标识
	DeviceType DeviceType      `json:"deviceType"` // 设备类型
	Total      int             `json:"total"`      // 结果总数
}

// AlertList 3.2 设备报警列表，时间跨度最大支持三个月
func (sdk *SolarmanSDK) AlertList(req AlertListRequest) (*AlertListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] AlertList: deviceSn 必填")
	}
	if req.StartTimestamp <= 0 {
		return nil, fmt.Errorf("[Solarman] AlertList: startTimestamp 必填")
	}
	if req.EndTimestamp <= 0 {
		return nil, fmt.Errorf("[Solarman] AlertList: endTimestamp 必填")
	}

	var out AlertListResult
	if err := sdk.do(PathAlertList, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 3.3 设备实时数据 POST /device/v1.0/currentData
// ---------------------------------------------------------------------------

// DeviceDataPoint 数据点，key 的含义由设备点表决定
type DeviceDataPoint struct {
	rawHolder
	Key   string `json:"key"`   // 参数 key
	Value string `json:"value"` // 参数值，统一为字符串
	Unit  Str    `json:"unit"`  // 参数单位，可为 null
	Name  string `json:"name"`  // 参数名称
}

// CurrentDataRequest 3.3 设备实时数据请求
type CurrentDataRequest struct {
	DeviceSN string `json:"deviceSn"`           // 设备的序列号，必填
	DeviceID Int64  `json:"deviceId,omitempty"` // 设备 ID，传了优先于序列号
}

// CurrentDataResult 3.3 设备实时数据返回
type CurrentDataResult struct {
	Response                         // 公共响应字段
	rawHolder                        // 原始数据兜底
	DeviceSN       string            `json:"deviceSn"`       // 设备的唯一标识
	DeviceID       Int64             `json:"deviceId"`       // 设备在平台内的 ID
	DeviceType     DeviceType        `json:"deviceType"`     // 设备类型
	DeviceState    DeviceState       `json:"deviceState"`    // 设备当前状态
	DataList       []DeviceDataPoint `json:"dataList"`       // 数据点表
	CollectionTime Int64             `json:"collectionTime"` // 数据更新时间戳
}

// CurrentData 3.3 设备实时数据，dataList 为设备点表数据
func (sdk *SolarmanSDK) CurrentData(req CurrentDataRequest) (*CurrentDataResult, error) {
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] CurrentData: deviceSn 必填")
	}

	var out CurrentDataResult
	if err := sdk.do(PathCurrentData, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 3.4 设备历史数据 POST /device/v1.0/historical
// ---------------------------------------------------------------------------

// DeviceHistoricalRequest 3.4 设备历史数据请求
type DeviceHistoricalRequest struct {
	DeviceSN  string `json:"deviceSn"`           // 设备的唯一标识，必填
	DeviceID  Int64  `json:"deviceId,omitempty"` // 设备在平台内的 ID，可选
	StartTime string `json:"startTime"`          // 开始时间，格式随 timeType
	EndTime   string `json:"endTime"`            // 结束时间，格式随 timeType
	TimeType  int    `json:"timeType"`           // 统计维度 2日/3月/4年/5帧
}

// DeviceHistoricalParam 3.4 单个时间点的数据分组
type DeviceHistoricalParam struct {
	rawHolder
	CollectTime string            `json:"collectTime"` // 时间，格式随 timeType
	DataList    []DeviceDataPoint `json:"dataList"`    // 数据点列表
}

// DeviceHistoricalResult 3.4 设备历史数据返回
type DeviceHistoricalResult struct {
	Response                              // 公共响应字段
	rawHolder                             // 原始数据兜底
	DeviceID      Int64                   `json:"deviceId"`      // 设备在平台内的 ID
	DeviceSN      string                  `json:"deviceSn"`      // 设备的唯一标识
	DeviceType    DeviceType              `json:"deviceType"`    // 设备类型
	ParamDataList []DeviceHistoricalParam `json:"paramDataList"` // 历史数据列表
	TimeType      int                     `json:"timeType"`      // 统计维度
}

// DeviceHistorical 3.4 设备历史数据，按 timeType 决定时间格式
func (sdk *SolarmanSDK) DeviceHistorical(req DeviceHistoricalRequest) (*DeviceHistoricalResult, error) {
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] DeviceHistorical: deviceSn 必填")
	}
	if strings.TrimSpace(req.StartTime) == "" {
		return nil, fmt.Errorf("[Solarman] DeviceHistorical: startTime 必填")
	}
	if strings.TrimSpace(req.EndTime) == "" {
		return nil, fmt.Errorf("[Solarman] DeviceHistorical: endTime 必填")
	}
	if req.TimeType <= 0 {
		return nil, fmt.Errorf("[Solarman] DeviceHistorical: timeType 必填")
	}

	var out DeviceHistoricalResult
	if err := sdk.do(PathHistorical, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 3.5 设备库-设备列表 POST /device/v1.0/list
// ---------------------------------------------------------------------------

// DeviceListRequest 3.5 设备列表请求
type DeviceListRequest struct {
	PageRequest // 分页参数
}

// DeviceListItem 3.5 单个设备
type DeviceListItem struct {
	rawHolder
	DeviceID    Int64       `json:"deviceId"`    // 设备在平台内的 ID
	DeviceSN    string      `json:"deviceSn"`    // 设备的唯一标识
	ProductID   Str         `json:"productId"`   // 设备的产品编号
	DeviceState DeviceState `json:"deviceState"` // 设备当前状态
	DeviceType  DeviceType  `json:"deviceType"`  // 设备类型
	UpdateTime  Int64       `json:"updateTime"`  // 最后一条数据更新时间
}

// DeviceListResult 3.5 设备列表返回
type DeviceListResult struct {
	Response                    // 公共响应字段
	rawHolder                   // 原始数据兜底
	DeviceList []DeviceListItem `json:"deviceList"` // 设备列表
	Total      int              `json:"total"`      // 结果总数
}

// DeviceList 3.5 设备库-设备列表，仅设备商账号可调用
func (sdk *SolarmanSDK) DeviceList(req DeviceListRequest) (*DeviceListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()

	var out DeviceListResult
	if err := sdk.do(PathDeviceList, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 3.6 获取网关设备 SIM 卡信息 POST /device/v1.0/simInfo
// ---------------------------------------------------------------------------

// DeviceSimInfoRequest 3.6 网关 SIM 卡信息请求
type DeviceSimInfoRequest struct {
	DeviceSN string `json:"deviceSn"`           // 设备的唯一标识，必填
	DeviceID Int64  `json:"deviceId,omitempty"` // 设备在平台内的 ID，可选
}

// DeviceSimInfoResult 3.6 网关 SIM 卡信息返回
type DeviceSimInfoResult struct {
	Response         // 公共响应字段
	rawHolder        // 原始数据兜底
	LimitTime string `json:"limitTime"` // 流量截止日期 yyyy-MM-dd
	Status    int    `json:"status"`    // 通讯状态 1-6
}

// DeviceSimInfo 3.6 获取网关设备 SIM 卡信息
func (sdk *SolarmanSDK) DeviceSimInfo(req DeviceSimInfoRequest) (*DeviceSimInfoResult, error) {
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] DeviceSimInfo: deviceSn 必填")
	}

	var out DeviceSimInfoResult
	if err := sdk.do(PathSimInfo, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 3.8 设备通讯关系 POST /device/v1.0/communication
// ---------------------------------------------------------------------------

// DeviceCommunicationRequest 3.8 设备通讯关系请求
type DeviceCommunicationRequest struct {
	DeviceSN string `json:"deviceSn"`           // 设备的唯一标识，必填
	DeviceID Int64  `json:"deviceId,omitempty"` // 设备在平台内的 ID，可选
}

// DeviceCommunication 通讯关系节点，子设备递归同结构
type DeviceCommunication struct {
	rawHolder
	DeviceSN    string                `json:"deviceSn"`    // 设备的唯一标识
	DeviceID    Int64                 `json:"deviceId"`    // 设备在平台内的 ID
	ParentSN    string                `json:"parentSn"`    // 父设备唯一标识
	DeviceType  DeviceType            `json:"deviceType"`  // 设备类型
	DeviceState DeviceState           `json:"deviceState"` // 设备当前状态
	UpdateTime  Int64                 `json:"updateTime"`  // 最后数据更新时间戳
	TimeZone    string                `json:"timeZone"`    // 设备时区
	ChildList   []DeviceCommunication `json:"childList"`   // 子设备列表
}

// DeviceCommunicationResult 3.8 设备通讯关系返回
type DeviceCommunicationResult struct {
	Response                           // 公共响应字段
	rawHolder                          // 原始数据兜底
	Communication *DeviceCommunication `json:"communication"` // 连接关系
}

// DeviceCommunication 3.8 设备通讯关系，返回网关与子设备层级
func (sdk *SolarmanSDK) DeviceCommunication(req DeviceCommunicationRequest) (*DeviceCommunicationResult, error) {
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] DeviceCommunication: deviceSn 必填")
	}

	var out DeviceCommunicationResult
	if err := sdk.do(PathCommunication, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 3.9 自定义指令（指令透传） POST /device/v1.0/customControl
// ---------------------------------------------------------------------------

// CustomControlRequest 3.9 自定义指令请求
type CustomControlRequest struct {
	CallBackURL    string `json:"callBackUrl"`              // 指令结果回调地址，必填
	Content        string `json:"content"`                  // 透传指令内容，十六进制，必填
	DeviceSN       string `json:"deviceSN"`                 // 设备的序列号，必填
	DeviceID       Int64  `json:"deviceId,omitempty"`       // 设备 ID，传了优先于序列号
	TimeoutSeconds int    `json:"timeoutSeconds,omitempty"` // 超时秒数 10-600，默认 600
}

// CustomControlResult 3.9 自定义指令返回
type CustomControlResult struct {
	Response              // 公共响应字段
	rawHolder             // 原始数据兜底
	CollectionTime string `json:"collectionTime"` // 最新数据接收时间戳
	ConnectStatus  int    `json:"connectStatus"`  // 通讯状态 0离线 1在线
	OrderID        string `json:"orderId"`        // 指令任务的 ID
}

// CustomControlCallback 3.9 指令结果回调请求体，回调方需返回 SUCCESS
type CustomControlCallback struct {
	rawHolder
	Ack        string `json:"ack"`        // 设备反馈的原始数值
	Content    string `json:"content"`    // 平台下发的指令内容
	CreateTime Str    `json:"createTime"` // 指令任务创建时间
	DeviceID   Int64  `json:"deviceId"`   // 设备在平台内的 ID
	DeviceSN   string `json:"deviceSn"`   // 设备的唯一标识
	OrderID    string `json:"orderId"`    // 指令任务的 ID
	Success    bool   `json:"success"`    // 指令任务反馈结果
	UpdateTime Str    `json:"updateTime"` // 指令任务反馈时间
	ErrorCode  string `json:"errorCode"`  // 错误码 530 下发失败/573 边缘通信失败
}

// CustomControl 3.9 自定义指令透传，结果通过回调地址异步返回
func (sdk *SolarmanSDK) CustomControl(req CustomControlRequest) (*CustomControlResult, error) {
	if strings.TrimSpace(req.CallBackURL) == "" {
		return nil, fmt.Errorf("[Solarman] CustomControl: callBackUrl 必填")
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("[Solarman] CustomControl: content 必填")
	}
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] CustomControl: deviceSN 必填")
	}

	var out CustomControlResult
	if err := sdk.do(PathCustomControl, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.13 添加网关（采集器或 DTU） POST /device/v1.0/register
// ---------------------------------------------------------------------------

// DeviceRegisterRequest 4.13 添加网关请求
type DeviceRegisterRequest struct {
	DeviceSN  string `json:"deviceSn"`  // 设备的唯一标识，必填
	IsAuto    bool   `json:"isAuto"`    // 是否自动发现，必填
	StationID Int64  `json:"stationId"` // 电站 ID，必填
}

// DeviceRegisterResult 4.13 添加网关返回
type DeviceRegisterResult struct {
	Response        // 公共响应字段
	rawHolder       // 原始数据兜底
	GatewayID Int64 `json:"gatewayId"` // 网关 ID
}

// DeviceRegister 4.13 添加网关（采集器或 DTU）到指定电站
func (sdk *SolarmanSDK) DeviceRegister(req DeviceRegisterRequest) (*DeviceRegisterResult, error) {
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] DeviceRegister: deviceSn 必填")
	}
	if req.StationID <= 0 {
		return nil, fmt.Errorf("[Solarman] DeviceRegister: stationId 必填")
	}

	var out DeviceRegisterResult
	if err := sdk.do(PathDeviceRegister, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.14 删除设备 POST /device/v1.0/delete
// ---------------------------------------------------------------------------

// DeviceDeleteRequest 4.14 删除设备请求
type DeviceDeleteRequest struct {
	DeviceSN  string `json:"deviceSn"`  // 设备的唯一标识，必填
	StationID Int64  `json:"stationId"` // 电站 ID，必填
}

// DeviceDeleteResult 4.14 删除设备返回，仅含公共响应字段
type DeviceDeleteResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// DeviceDelete 4.14 从指定电站删除设备
func (sdk *SolarmanSDK) DeviceDelete(req DeviceDeleteRequest) (*DeviceDeleteResult, error) {
	if strings.TrimSpace(req.DeviceSN) == "" {
		return nil, fmt.Errorf("[Solarman] DeviceDelete: deviceSn 必填")
	}
	if req.StationID <= 0 {
		return nil, fmt.Errorf("[Solarman] DeviceDelete: stationId 必填")
	}

	var out DeviceDeleteResult
	if err := sdk.do(PathDeviceDelete, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}
