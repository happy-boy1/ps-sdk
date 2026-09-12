package ginlong

import (
	"fmt"
	"strings"
)

// ---------- 3.1 获取账号下逆变器列表 /v1/api/inverterList ----------

// InverterListRequest 账号下逆变器列表请求
type InverterListRequest struct {
	PageRequest          // 分页参数
	StationID   Int64    `json:"stationId,omitempty"` // 电站id
	NmiCode     string   `json:"nmiCode,omitempty"`   // nmi编码
	SNList      []string `json:"snList,omitempty"`    // 逆变器SN列表
}

// InverterStatusVo 逆变器状态统计
type InverterStatusVo struct {
	All     int `json:"all"`     // 逆变器总数
	Normal  int `json:"normal"`  // 正常逆变器数
	Offline int `json:"offline"` // 离线逆变器数
	Fault   int `json:"fault"`   // 故障逆变器数
}

// InverterListResult 账号下逆变器列表结果
type InverterListResult struct {
	rawHolder
	Page             Page[InverterListItem] `json:"page"`             // 分页列表
	InverterStatusVo InverterStatusVo       `json:"inverterStatusVo"` // 状态统计
}

// InverterListItem 逆变器列表项
type InverterListItem struct {
	rawHolder
	ID                 Int64       `json:"id"`                 // 逆变器id
	SN                 string      `json:"sn"`                 // 逆变器SN
	StationID          Int64       `json:"stationId"`          // 电站id
	StationName        string      `json:"stationName"`        // 电站名称
	UserID             Int64       `json:"userId"`             // 业主id
	Power              Num         `json:"power"`              // 装机容量
	PowerStr           string      `json:"powerStr"`           // 装机容量单位
	Etoday             Num         `json:"etoday"`             // 当日能量
	Etoday1            Num         `json:"etoday1"`            // 当日发电量原始值
	EtodayStr          string      `json:"etodayStr"`          // 当日能量单位
	Etotal             Num         `json:"etotal"`             // 总能量
	Etotal1            Num         `json:"etotal1"`            // 累计发电量原始值
	EtotalStr          string      `json:"etotalStr"`          // 总能量单位
	FullHour           Num         `json:"fullHour"`           // 满发小时数
	Pac                Num         `json:"pac"`                // 功率
	PacStr             string      `json:"pacStr"`             // 功率单位
	State              DeviceState `json:"state"`              // 逆变器状态
	DataTimestamp      Int64       `json:"dataTimestamp"`      // 更新时间
	CollectorSn        string      `json:"collectorSn"`        // 采集器SN
	ProductModel       string      `json:"productModel"`       // 逆变器类型
	DcInputType        int         `json:"dcInputType"`        // 直流输入路数
	AcOutputType       int         `json:"acOutputType"`       // 交流输出类型
	Series             string      `json:"series"`             // 逆变器系列
	Name               string      `json:"name"`               // 逆变器名称
	Addr               string      `json:"addr"`               // 电站地址
	CollectorState     int         `json:"collectorState"`     // 采集器状态
	StateExceptionFlag int         `json:"stateExceptionFlag"` // 离线状态标志
	TotalFullHour      Num         `json:"totalFullHour"`      // 累计满发小时数
	InverterMeterModel MeterType   `json:"inverterMeterModel"` // 逆变器电表类型
	CreateDate         Int64       `json:"createDate"`         // 创建时间
	UpdateShelfEndTime Int64       `json:"updateShelfEndTime"` // 质保结束时间
}

// InverterList 获取账号下逆变器列表
func (c *SolisSDK) InverterList(req InverterListRequest) (*InverterListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out InverterListResult
	if err := c.do(PathInverterList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 3.2 获取单台逆变器详情 /v1/api/inverterDetail ----------

// InverterDetailRequest 单台逆变器详情请求
type InverterDetailRequest struct {
	ID Int64  `json:"id,omitempty"` // 逆变器id
	SN string `json:"sn,omitempty"` // 逆变器SN
}

// InverterDetailResult 单台逆变器详情结果
type InverterDetailResult struct {
	rawHolder
	ID                             Int64       `json:"id"`                             // 逆变器id
	SN                             string      `json:"sn"`                             // 逆变器SN
	StationID                      Int64       `json:"stationId"`                      // 电站id
	StationName                    string      `json:"stationName"`                    // 电站名称
	UserID                         Int64       `json:"userId"`                         // 业主id
	CollectorID                    Int64       `json:"collectorId"`                    // 采集器id
	CollectorName                  string      `json:"collectorName"`                  // 采集器名称
	Collectorsn                    string      `json:"collectorsn"`                    // 采集器SN
	CurrentState                   string      `json:"currentState"`                   // 当前状态
	EToday                         Num         `json:"eToday"`                         // 当日能量
	ETodayStr                      string      `json:"eTodayStr"`                      // 当日能量单位
	EMonth                         Num         `json:"eMonth"`                         // 当月能量
	EMonthStr                      string      `json:"eMonthStr"`                      // 当月能量单位
	EYear                          Num         `json:"eYear"`                          // 当年能量
	EYearStr                       string      `json:"eYearStr"`                       // 当年能量单位
	ETotal                         Num         `json:"eTotal"`                         // 总能量
	ETotalStr                      string      `json:"eTotalStr"`                      // 总能量单位
	Fac                            Num         `json:"fac"`                            // 电网频率
	FacStr                         string      `json:"facStr"`                         // 电网频率单位
	Pac                            Num         `json:"pac"`                            // 实时功率
	PacStr                         string      `json:"pacStr"`                         // 实时功率单位
	PacPec                         Num         `json:"pacPec"`                         // 功率百分比
	FullHour                       Num         `json:"fullHour"`                       // 满发小时数
	PicName                        string      `json:"picName"`                        // 图片
	Power                          Num         `json:"power"`                          // 装机容量
	PowerStr                       string      `json:"powerStr"`                       // 装机容量单位
	IAc1                           Num         `json:"iAc1"`                           // 交流电流R
	IAc2                           Num         `json:"iAc2"`                           // 交流电流S
	IAc3                           Num         `json:"iAc3"`                           // 交流电流T
	UAc1                           Num         `json:"uAc1"`                           // 交流电压R
	UAc2                           Num         `json:"uAc2"`                           // 交流电压S
	UAc3                           Num         `json:"uAc3"`                           // 交流电压T
	State                          DeviceState `json:"state"`                          // 逆变器状态
	DataTimestamp                  Int64       `json:"dataTimestamp"`                  // 更新时间
	InverterTemperature            Num         `json:"inverterTemperature"`            // 逆变器温度
	NationalStandardstr            string      `json:"nationalStandardstr"`            // 国家标准
	AcOutputType                   int         `json:"acOutputType"`                   // 交流输出类型
	DcInputType                    int         `json:"dcInputType"`                    // 直流输入路数
	PowerFactor                    Num         `json:"powerFactor"`                    // 功率因数
	BatteryPower                   Num         `json:"batteryPower"`                   // 电池功率
	BatteryPowerStr                string      `json:"batteryPowerStr"`                // 电池功率单位
	BatteryPowerPec                Num         `json:"batteryPowerPec"`                // 电池功率百分比
	BatteryCapacitySoc             Num         `json:"batteryCapacitySoc"`             // 电池SOC
	BatteryHealthSoh               Num         `json:"batteryHealthSoh"`               // 电池SOH
	SocDischargeSet                Num         `json:"socDischargeSet"`                // 过放SOC
	SocChargingSet                 Num         `json:"socChargingSet"`                 // 强充SOC
	BatteryType                    string      `json:"batteryType"`                    // 当前运行电池型号
	BatteryVoltage                 Num         `json:"batteryVoltage"`                 // 电池电压
	BatteryVoltageStr              string      `json:"batteryVoltageStr"`              // 电池电压单位
	BstteryCurrent                 Num         `json:"bstteryCurrent"`                 // 电池电流
	BstteryCurrentStr              string      `json:"bstteryCurrentStr"`              // 电池电流单位
	BatteryTodayChargeEnergy       Num         `json:"batteryTodayChargeEnergy"`       // 当日电池充电电量
	BatteryTodayChargeEnergyStr    string      `json:"batteryTodayChargeEnergyStr"`    // 当日电池充电电量单位
	BatteryMonthChargeEnergy       Num         `json:"batteryMonthChargeEnergy"`       // 当月电池充电电量
	BatteryMonthChargeEnergyStr    string      `json:"batteryMonthChargeEnergyStr"`    // 当月电池充电电量单位
	BatteryYearChargeEnergy        Num         `json:"batteryYearChargeEnergy"`        // 当年电池充电电量
	BatteryYearChargeEnergyStr     string      `json:"batteryYearChargeEnergyStr"`     // 当年电池充电电量单位
	BatteryTotalChargeEnergy       Num         `json:"batteryTotalChargeEnergy"`       // 电池总充电电量
	BatteryTotalChargeEnergyStr    string      `json:"batteryTotalChargeEnergyStr"`    // 电池总充电电量单位
	BatteryTodayDischargeEnergy    Num         `json:"batteryTodayDischargeEnergy"`    // 当日电池放电电量
	BatteryTodayDischargeEnergyStr string      `json:"batteryTodayDischargeEnergyStr"` // 当日电池放电电量单位
	BatteryMonthDischargeEnergy    Num         `json:"batteryMonthDischargeEnergy"`    // 当月电池放电电量
	BatteryMonthDischargeEnergyStr string      `json:"batteryMonthDischargeEnergyStr"` // 当月电池放电电量单位
	BatteryYearDischargeEnergy     Num         `json:"batteryYearDischargeEnergy"`     // 当年电池放电电量
	BatteryYearDischargeEnergyStr  string      `json:"batteryYearDischargeEnergyStr"`  // 当年电池放电电量单位
	BatteryTotalDischargeEnergy    Num         `json:"batteryTotalDischargeEnergy"`    // 电池总放电量
	BatteryTotalDischargeEnergyStr string      `json:"batteryTotalDischargeEnergyStr"` // 电池总放电量单位
	GridPurchasedTodayEnergy       Num         `json:"gridPurchasedTodayEnergy"`       // 当日电表买电
	GridPurchasedTodayEnergyStr    string      `json:"gridPurchasedTodayEnergyStr"`    // 当日电表买电单位
	GridPurchasedMonthEnergy       Num         `json:"gridPurchasedMonthEnergy"`       // 当月电表买电
	GridPurchasedMonthEnergyStr    string      `json:"gridPurchasedMonthEnergyStr"`    // 当月电表买电单位
	GridPurchasedYearEnergy        Num         `json:"gridPurchasedYearEnergy"`        // 当年电表买电
	GridPurchasedYearEnergyStr     string      `json:"gridPurchasedYearEnergyStr"`     // 当年电表买电单位
	GridPurchasedTotalEnergy       Num         `json:"gridPurchasedTotalEnergy"`       // 累计电表买电
	GridPurchasedTotalEnergyStr    string      `json:"gridPurchasedTotalEnergyStr"`    // 累计电表买电单位
	GridSellTodayEnergy            Num         `json:"gridSellTodayEnergy"`            // 当日电表卖电
	GridSellTodayEnergyStr         string      `json:"gridSellTodayEnergyStr"`         // 当日电表卖电单位
	GridSellMonthEnergy            Num         `json:"gridSellMonthEnergy"`            // 当月电表卖电
	GridSellMonthEnergyStr         string      `json:"gridSellMonthEnergyStr"`         // 当月电表卖电单位
	GridSellYearEnergy             Num         `json:"gridSellYearEnergy"`             // 当年电表卖电
	GridSellYearEnergyStr          string      `json:"gridSellYearEnergyStr"`          // 当年电表卖电单位
	GridSellTotalEnergy            Num         `json:"gridSellTotalEnergy"`            // 累计电表卖电
	GridSellTotalEnergyStr         string      `json:"gridSellTotalEnergyStr"`         // 累计电表卖电单位
	FamilyLoadPower                Num         `json:"familyLoadPower"`                // 家庭负载功率
	FamilyLoadPowerStr             string      `json:"familyLoadPowerStr"`             // 家庭负载功率单位
	BypassLoadPower                Num         `json:"bypassLoadPower"`                // 旁路负载功率
	BypassLoadPowerStr             string      `json:"bypassLoadPowerStr"`             // 旁路负载功率单位
	PSum                           Num         `json:"pSum"`                           // 电网总有功功率
	PSumStr                        string      `json:"pSumStr"`                        // 电网总有功功率单位
	PsumPec                        Num         `json:"psumPec"`                        // 电网总有功功率百分比
	HomeLoadTodayEnergy            Num         `json:"homeLoadTodayEnergy"`            // 当日负载用电电量
	HomeLoadTodayEnergyStr         string      `json:"homeLoadTodayEnergyStr"`         // 当日负载用电电量单位
	HomeLoadTotalEnergy            Num         `json:"homeLoadTotalEnergy"`            // 累计负载用电量
	HomeLoadTotalEnergyStr         string      `json:"homeLoadTotalEnergyStr"`         // 累计负载用电量单位
	Model                          string      `json:"model"`                          // 逆变器model号
	Type                           int         `json:"type"`                           // 逆变器类型
	Name                           string      `json:"name"`                           // 逆变器名称
	InverterMeterModel             MeterType   `json:"inverterMeterModel"`             // 逆变器电表类型
	StateExceptionFlag             int         `json:"stateExceptionFlag"`             // 离线状态标志
	CollectorState                 int         `json:"collectorState"`                 // 采集器状态
	CollectorModel                 string      `json:"collectorModel"`                 // 采集器model
	WarningInfoData                int         `json:"warningInfoData"`                // 警告信息故障数据
	ProductModel                   string      `json:"productModel"`                   // 产品型号
	NationalStandards              string      `json:"nationalStandards"`              // 执行的国家标准
	Version                        string      `json:"version"`                        // 逆变器软件版本
	ReactivePower                  Num         `json:"reactivePower"`                  // 逆变器无功功率
	ReactivePowerStr               string      `json:"reactivePowerStr"`               // 无功功率单位
	ApparentPower                  Num         `json:"apparentPower"`                  // 逆变器视在功率
	ApparentPowerStr               string      `json:"apparentPowerStr"`               // 视在功率单位
	DcPac                          Num         `json:"dcPac"`                          // 总直流输入功率
	DcPacStr                       string      `json:"dcPacStr"`                       // 总直流输入功率单位
	UpdateShelfEndTime             Int64       `json:"updateShelfEndTime"`             // 质保结束时间
	IA                             Num         `json:"iA"`                             // 电表A相电流
	UA                             Num         `json:"uA"`                             // 电表A相电压
	ALookedPower                   Num         `json:"aLookedPower"`                   // 电表A相视在功率
	AReactivePower                 Num         `json:"aReactivePower"`                 // 电表A相无功功率
	AphasePowerFactor              Num         `json:"aPhasePowerFactor"`              // 电表A相有功功率（示例拼写，表为 aphasePowerFactor）
	AveragePowerFactor             Num         `json:"averagePowerFactor"`             // 电表功率因数
	IB                             Num         `json:"iB"`                             // 电表B相电流
	UB                             Num         `json:"uB"`                             // 电表B相电压
	BLookedPower                   Num         `json:"bLookedPower"`                   // 电表B相视在功率
	BReactivePower                 Num         `json:"bReactivePower"`                 // 电表B相无功功率
	BphasePowerFactor              Num         `json:"bPhasePowerFactor"`              // 电表B相有功功率（示例拼写，表为 bphasePowerFactor）
	IC                             Num         `json:"iC"`                             // 电表C相电流
	UC                             Num         `json:"uC"`                             // 电表C相电压
	CLookedPower                   Num         `json:"cLookedPower"`                   // 电表C相视在功率
	CReactivePower                 Num         `json:"cReactivePower"`                 // 电表C相无功功率
	CphasePowerFactor              Num         `json:"cPhasePowerFactor"`              // 电表C相有功功率（示例拼写，表为 cphasePowerFactor）
	FAc                            Num         `json:"fAc"`                            // 电表电网频率
	MpptIpv1                       Num         `json:"mpptIpv1"`                       // mppt电流1
	MpptUpv1                       Num         `json:"mpptUpv1"`                       // mppt电压1
	MpptIpv2                       Num         `json:"mpptIpv2"`                       // mppt电流2
	MpptUpv2                       Num         `json:"mpptUpv2"`                       // mppt电压2
	MpptIpv3                       Num         `json:"mpptIpv3"`                       // mppt电流3
	MpptUpv3                       Num         `json:"mpptUpv3"`                       // mppt电压3
	MpptIpv4                       Num         `json:"mpptIpv4"`                       // mppt电流4
	MpptUpv4                       Num         `json:"mpptUpv4"`                       // mppt电压4
	MpptIpv5                       Num         `json:"mpptIpv5"`                       // mppt电流5
	MpptUpv5                       Num         `json:"mpptUpv5"`                       // mppt电压5
	MpptIpv6                       Num         `json:"mpptIpv6"`                       // mppt电流6
	MpptUpv6                       Num         `json:"mpptUpv6"`                       // mppt电压6
	MpptIpv7                       Num         `json:"mpptIpv7"`                       // mppt电流7
	MpptUpv7                       Num         `json:"mpptUpv7"`                       // mppt电压7
	MpptIpv8                       Num         `json:"mpptIpv8"`                       // mppt电流8
	MpptUpv8                       Num         `json:"mpptUpv8"`                       // mppt电压8
	MpptIpv9                       Num         `json:"mpptIpv9"`                       // mppt电流9
	MpptUpv9                       Num         `json:"mpptUpv9"`                       // mppt电压9
	MpptIpv10                      Num         `json:"mpptIpv10"`                      // mppt电流10
	MpptUpv10                      Num         `json:"mpptUpv10"`                      // mppt电压10
	MpptIpv11                      Num         `json:"mpptIpv11"`                      // mppt电流11
	MpptUpv11                      Num         `json:"mpptUpv11"`                      // mppt电压11
	MpptIpv12                      Num         `json:"mpptIpv12"`                      // mppt电流12
	MpptUpv12                      Num         `json:"mpptUpv12"`                      // mppt电压12
	MpptIpv13                      Num         `json:"mpptIpv13"`                      // mppt电流13
	MpptUpv13                      Num         `json:"mpptUpv13"`                      // mppt电压13
	MpptIpv14                      Num         `json:"mpptIpv14"`                      // mppt电流14
	MpptUpv14                      Num         `json:"mpptUpv14"`                      // mppt电压14
	MpptIpv15                      Num         `json:"mpptIpv15"`                      // mppt电流15
	MpptUpv15                      Num         `json:"mpptUpv15"`                      // mppt电压15
	MpptIpv16                      Num         `json:"mpptIpv16"`                      // mppt电流16
	MpptUpv16                      Num         `json:"mpptUpv16"`                      // mppt电压16
	MpptIpv17                      Num         `json:"mpptIpv17"`                      // mppt电流17
	MpptUpv17                      Num         `json:"mpptUpv17"`                      // mppt电压17
	MpptIpv18                      Num         `json:"mpptIpv18"`                      // mppt电流18
	MpptUpv18                      Num         `json:"mpptUpv18"`                      // mppt电压18
	MpptIpv19                      Num         `json:"mpptIpv19"`                      // mppt电流19
	MpptUpv19                      Num         `json:"mpptUpv19"`                      // mppt电压19
	MpptIpv20                      Num         `json:"mpptIpv20"`                      // mppt电流20
	MpptUpv20                      Num         `json:"mpptUpv20"`                      // mppt电压20
	DcInputTypeMppt                Num         `json:"dcInputTypeMppt"`                // MPPT直流输入路数
}

// InverterDetail 获取单台逆变器详情
func (c *SolisSDK) InverterDetail(req InverterDetailRequest) (*InverterDetailResult, error) {
	if req.ID.Int() == 0 && strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("获取单台逆变器详情: id 和 sn 不能同时为空")
	}
	var out InverterDetailResult
	if err := c.do(PathInverterDetail, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 3.3 获取多台逆变器详情 /v1/api/inverterDetailList ----------

// InverterDetailListRequest 多台逆变器详情请求
type InverterDetailListRequest struct {
	PageRequest             // 分页参数
	SNList      []string    `json:"snList,omitempty"`      // 逆变器SN列表
	StationType StationType `json:"stationType,omitempty"` // 电站类型
}

// InverterDetailListItem 多台逆变器详情项
type InverterDetailListItem struct {
	rawHolder
	ID                             Int64       `json:"id"`                             // 逆变器id
	SN                             string      `json:"sn"`                             // 逆变器SN
	StationID                      Int64       `json:"stationId"`                      // 电站id
	StationName                    string      `json:"stationName"`                    // 电站名称
	UserID                         Int64       `json:"userId"`                         // 业主id
	CollectorID                    Int64       `json:"collectorId"`                    // 采集器id
	CollectorName                  string      `json:"collectorName"`                  // 采集器名称
	Collectorsn                    string      `json:"collectorsn"`                    // 采集器SN
	CurrentState                   string      `json:"currentState"`                   // 当前状态
	EToday                         Num         `json:"eToday"`                         // 当日能量
	ETodayStr                      string      `json:"eTodayStr"`                      // 当日能量单位
	EMonth                         Num         `json:"eMonth"`                         // 当月能量
	EMonthStr                      string      `json:"eMonthStr"`                      // 当月能量单位
	EYear                          Num         `json:"eYear"`                          // 当年能量
	EYearStr                       string      `json:"eYearStr"`                       // 当年能量单位
	ETotal                         Num         `json:"eTotal"`                         // 总能量
	ETotalStr                      string      `json:"eTotalStr"`                      // 总能量单位
	Fac                            Num         `json:"fac"`                            // 电网频率
	FacStr                         string      `json:"facStr"`                         // 电网频率单位
	Pac                            Num         `json:"pac"`                            // 实时功率
	PacStr                         string      `json:"pacStr"`                         // 实时功率单位
	PacPec                         Num         `json:"pacPec"`                         // 功率百分比
	FullHour                       Num         `json:"fullHour"`                       // 满发小时数
	PicName                        string      `json:"picName"`                        // 图片
	Power                          Num         `json:"power"`                          // 装机容量
	PowerStr                       string      `json:"powerStr"`                       // 装机容量单位
	IAc1                           Num         `json:"iAc1"`                           // 交流电流R
	IAc2                           Num         `json:"iAc2"`                           // 交流电流S
	IAc3                           Num         `json:"iAc3"`                           // 交流电流T
	UAc1                           Num         `json:"uAc1"`                           // 交流电压R
	UAc2                           Num         `json:"uAc2"`                           // 交流电压S
	UAc3                           Num         `json:"uAc3"`                           // 交流电压T
	State                          DeviceState `json:"state"`                          // 逆变器状态
	DataTimestamp                  Int64       `json:"dataTimestamp"`                  // 更新时间
	InverterTemperature            Num         `json:"inverterTemperature"`            // 逆变器温度
	NationalStandardstr            string      `json:"nationalStandardstr"`            // 国家标准
	AcOutputType                   int         `json:"acOutputType"`                   // 交流输出类型
	DcInputType                    int         `json:"dcInputType"`                    // 直流输入路数
	PowerFactor                    Num         `json:"powerFactor"`                    // 功率因数
	BatteryPower                   Num         `json:"batteryPower"`                   // 电池功率
	BatteryPowerStr                string      `json:"batteryPowerStr"`                // 电池功率单位
	BatteryPowerPec                Num         `json:"batteryPowerPec"`                // 电池功率百分比
	BatteryCapacitySoc             Num         `json:"batteryCapacitySoc"`             // 电池SOC
	BatteryHealthSoh               Num         `json:"batteryHealthSoh"`               // 电池SOH
	SocDischargeSet                Num         `json:"socDischargeSet"`                // 过放SOC
	SocChargingSet                 Num         `json:"socChargingSet"`                 // 强充SOC
	BatteryType                    string      `json:"batteryType"`                    // 当前运行电池型号
	BatteryVoltage                 Num         `json:"batteryVoltage"`                 // 电池电压
	BatteryVoltageStr              string      `json:"batteryVoltageStr"`              // 电池电压单位
	BstteryCurrent                 Num         `json:"bstteryCurrent"`                 // 电池电流
	BstteryCurrentStr              string      `json:"bstteryCurrentStr"`              // 电池电流单位
	BatteryTodayChargeEnergy       Num         `json:"batteryTodayChargeEnergy"`       // 当日电池充电电量
	BatteryTodayChargeEnergyStr    string      `json:"batteryTodayChargeEnergyStr"`    // 当日电池充电电量单位
	BatteryMonthChargeEnergy       Num         `json:"batteryMonthChargeEnergy"`       // 当月电池充电电量
	BatteryMonthChargeEnergyStr    string      `json:"batteryMonthChargeEnergyStr"`    // 当月电池充电电量单位
	BatteryYearChargeEnergy        Num         `json:"batteryYearChargeEnergy"`        // 当年电池充电电量
	BatteryYearChargeEnergyStr     string      `json:"batteryYearChargeEnergyStr"`     // 当年电池充电电量单位
	BatteryTotalChargeEnergy       Num         `json:"batteryTotalChargeEnergy"`       // 电池总充电电量
	BatteryTotalChargeEnergyStr    string      `json:"batteryTotalChargeEnergyStr"`    // 电池总充电电量单位
	BatteryTodayDischargeEnergy    Num         `json:"batteryTodayDischargeEnergy"`    // 当日电池放电电量
	BatteryTodayDischargeEnergyStr string      `json:"batteryTodayDischargeEnergyStr"` // 当日电池放电电量单位
	BatteryMonthDischargeEnergy    Num         `json:"batteryMonthDischargeEnergy"`    // 当月电池放电电量
	BatteryMonthDischargeEnergyStr string      `json:"batteryMonthDischargeEnergyStr"` // 当月电池放电电量单位
	BatteryYearDischargeEnergy     Num         `json:"batteryYearDischargeEnergy"`     // 当年电池放电电量
	BatteryYearDischargeEnergyStr  string      `json:"batteryYearDischargeEnergyStr"`  // 当年电池放电电量单位
	BatteryTotalDischargeEnergy    Num         `json:"batteryTotalDischargeEnergy"`    // 电池总放电量
	BatteryTotalDischargeEnergyStr string      `json:"batteryTotalDischargeEnergyStr"` // 电池总放电量单位
	GridPurchasedTodayEnergy       Num         `json:"gridPurchasedTodayEnergy"`       // 当日电表买电
	GridPurchasedTodayEnergyStr    string      `json:"gridPurchasedTodayEnergyStr"`    // 当日电表买电单位
	GridPurchasedMonthEnergy       Num         `json:"gridPurchasedMonthEnergy"`       // 当月电表买电
	GridPurchasedMonthEnergyStr    string      `json:"gridPurchasedMonthEnergyStr"`    // 当月电表买电单位
	GridPurchasedYearEnergy        Num         `json:"gridPurchasedYearEnergy"`        // 当年电表买电
	GridPurchasedYearEnergyStr     string      `json:"gridPurchasedYearEnergyStr"`     // 当年电表买电单位
	GridPurchasedTotalEnergy       Num         `json:"gridPurchasedTotalEnergy"`       // 累计电表买电
	GridPurchasedTotalEnergyStr    string      `json:"gridPurchasedTotalEnergyStr"`    // 累计电表买电单位
	GridSellTodayEnergy            Num         `json:"gridSellTodayEnergy"`            // 当日电表卖电
	GridSellTodayEnergyStr         string      `json:"gridSellTodayEnergyStr"`         // 当日电表卖电单位
	GridSellMonthEnergy            Num         `json:"gridSellMonthEnergy"`            // 当月电表卖电
	GridSellMonthEnergyStr         string      `json:"gridSellMonthEnergyStr"`         // 当月电表卖电单位
	GridSellYearEnergy             Num         `json:"gridSellYearEnergy"`             // 当年电表卖电
	GridSellYearEnergyStr          string      `json:"gridSellYearEnergyStr"`          // 当年电表卖电单位
	GridSellTotalEnergy            Num         `json:"gridSellTotalEnergy"`            // 累计电表卖电
	GridSellTotalEnergyStr         string      `json:"gridSellTotalEnergyStr"`         // 累计电表卖电单位
	FamilyLoadPower                Num         `json:"familyLoadPower"`                // 家庭负载功率
	FamilyLoadPowerStr             string      `json:"familyLoadPowerStr"`             // 家庭负载功率单位
	BypassLoadPower                Num         `json:"bypassLoadPower"`                // 旁路负载功率
	BypassLoadPowerStr             string      `json:"bypassLoadPowerStr"`             // 旁路负载功率单位
	PSum                           Num         `json:"pSum"`                           // 电网总有功功率
	PSumStr                        string      `json:"pSumStr"`                        // 电网总有功功率单位
	PsumPec                        Num         `json:"psumPec"`                        // 电网总有功功率百分比
	HomeLoadTodayEnergy            Num         `json:"homeLoadTodayEnergy"`            // 当日负载用电电量
	HomeLoadTodayEnergyStr         string      `json:"homeLoadTodayEnergyStr"`         // 当日负载用电电量单位
	HomeLoadTotalEnergy            Num         `json:"homeLoadTotalEnergy"`            // 累计负载用电量
	HomeLoadTotalEnergyStr         string      `json:"homeLoadTotalEnergyStr"`         // 累计负载用电量单位
	Model                          string      `json:"model"`                          // 逆变器Model
	Type                           int         `json:"type"`                           // 逆变器类型
	Name                           string      `json:"name"`                           // 逆变器名称
	InverterMeterModel             MeterType   `json:"inverterMeterModel"`             // 逆变器电表类型
	StateExceptionFlag             int         `json:"stateExceptionFlag"`             // 离线状态标志
	CollectorState                 int         `json:"collectorState"`                 // 采集器状态
	CollectorModel                 string      `json:"collectorModel"`                 // 采集器model
	WarningInfoData                int         `json:"warningInfoData"`                // 警告信息故障数据
	ProductModel                   string      `json:"productModel"`                   // 产品型号
	NationalStandards              string      `json:"nationalStandards"`              // 执行的国家标准
	Version                        string      `json:"version"`                        // 逆变器软件版本
	ReactivePower                  Num         `json:"reactivePower"`                  // 逆变器无功功率
	ReactivePowerStr               string      `json:"reactivePowerStr"`               // 无功功率单位
	ApparentPower                  Num         `json:"apparentPower"`                  // 逆变器视在功率
	ApparentPowerStr               string      `json:"apparentPowerStr"`               // 视在功率单位
	DcPac                          Num         `json:"dcPac"`                          // 总直流输入功率
	DcPacStr                       string      `json:"dcPacStr"`                       // 总直流输入功率单位
	UpdateShelfEndTime             Int64       `json:"updateShelfEndTime"`             // 质保结束时间
	MpptIpv1                       Num         `json:"mpptIpv1"`                       // mppt电流1
	MpptUpv1                       Num         `json:"mpptUpv1"`                       // mppt电压1
	MpptIpv2                       Num         `json:"mpptIpv2"`                       // mppt电流2
	MpptUpv2                       Num         `json:"mpptUpv2"`                       // mppt电压2
	MpptIpv3                       Num         `json:"mpptIpv3"`                       // mppt电流3
	MpptUpv3                       Num         `json:"mpptUpv3"`                       // mppt电压3
	MpptIpv4                       Num         `json:"mpptIpv4"`                       // mppt电流4
	MpptUpv4                       Num         `json:"mpptUpv4"`                       // mppt电压4
	MpptIpv5                       Num         `json:"mpptIpv5"`                       // mppt电流5
	MpptUpv5                       Num         `json:"mpptUpv5"`                       // mppt电压5
	MpptIpv6                       Num         `json:"mpptIpv6"`                       // mppt电流6
	MpptUpv6                       Num         `json:"mpptUpv6"`                       // mppt电压6
	MpptIpv7                       Num         `json:"mpptIpv7"`                       // mppt电流7
	MpptUpv7                       Num         `json:"mpptUpv7"`                       // mppt电压7
	MpptIpv8                       Num         `json:"mpptIpv8"`                       // mppt电流8
	MpptUpv8                       Num         `json:"mpptUpv8"`                       // mppt电压8
	MpptIpv9                       Num         `json:"mpptIpv9"`                       // mppt电流9
	MpptUpv9                       Num         `json:"mpptUpv9"`                       // mppt电压9
	MpptIpv10                      Num         `json:"mpptIpv10"`                      // mppt电流10
	MpptUpv10                      Num         `json:"mpptUpv10"`                      // mppt电压10
	MpptIpv11                      Num         `json:"mpptIpv11"`                      // mppt电流11
	MpptUpv11                      Num         `json:"mpptUpv11"`                      // mppt电压11
	MpptIpv12                      Num         `json:"mpptIpv12"`                      // mppt电流12
	MpptUpv12                      Num         `json:"mpptUpv12"`                      // mppt电压12
	MpptIpv13                      Num         `json:"mpptIpv13"`                      // mppt电流13
	MpptUpv13                      Num         `json:"mpptUpv13"`                      // mppt电压13
	MpptIpv14                      Num         `json:"mpptIpv14"`                      // mppt电流14
	MpptUpv14                      Num         `json:"mpptUpv14"`                      // mppt电压14
	MpptIpv15                      Num         `json:"mpptIpv15"`                      // mppt电流15
	MpptUpv15                      Num         `json:"mpptUpv15"`                      // mppt电压15
	MpptIpv16                      Num         `json:"mpptIpv16"`                      // mppt电流16
	MpptUpv16                      Num         `json:"mpptUpv16"`                      // mppt电压16
	MpptIpv17                      Num         `json:"mpptIpv17"`                      // mppt电流17
	MpptUpv17                      Num         `json:"mpptUpv17"`                      // mppt电压17
	MpptIpv18                      Num         `json:"mpptIpv18"`                      // mppt电流18
	MpptUpv18                      Num         `json:"mpptUpv18"`                      // mppt电压18
	MpptIpv19                      Num         `json:"mpptIpv19"`                      // mppt电流19
	MpptUpv19                      Num         `json:"mpptUpv19"`                      // mppt电压19
	MpptIpv20                      Num         `json:"mpptIpv20"`                      // mppt电流20
	MpptUpv20                      Num         `json:"mpptUpv20"`                      // mppt电压20
	DcInputTypeMppt                Num         `json:"dcInputTypeMppt"`                // MPPT直流输入路数
}

// InverterDetailList 获取多台逆变器详情
func (c *SolisSDK) InverterDetailList(req InverterDetailListRequest) ([]InverterDetailListItem, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out []InverterDetailListItem
	if err := c.do(PathInverterDetailList, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 3.4 获取单台逆变器某日的实时数据 /v1/api/inverterDay ----------

// InverterDayRequest 逆变器某日实时数据请求
type InverterDayRequest struct {
	ID       Int64  `json:"id,omitempty"` // 逆变器id
	SN       string `json:"sn,omitempty"` // 逆变器SN
	Money    string `json:"money"`        // 电站货币单位
	Time     string `json:"time"`         // 查询日期 yyyy-MM-dd
	TimeZone int    `json:"timeZone"`     // 电站时区
}

// InverterDayItem 逆变器某日实时数据项
type InverterDayItem struct {
	rawHolder
	DataTimestamp               Int64  `json:"dataTimestamp"`               // 更新时间戳
	TimeStr                     string `json:"timeStr"`                     // 更新时间字符串
	EToday                      Num    `json:"eToday"`                      // 当日能量
	ETotal                      Num    `json:"eTotal"`                      // 总能量
	Fac                         Num    `json:"fac"`                         // 电网频率
	Pac                         Num    `json:"pac"`                         // 实时功率
	PacStr                      string `json:"pacStr"`                      // 实时功率单位
	PacPec                      Num    `json:"pacPec"`                      // 功率百分比
	IAc1                        Num    `json:"iAc1"`                        // 交流电流R
	IAc2                        Num    `json:"iAc2"`                        // 交流电流S
	IAc3                        Num    `json:"iAc3"`                        // 交流电流T
	UAc1                        Num    `json:"uAc1"`                        // 交流电压R
	UAc2                        Num    `json:"uAc2"`                        // 交流电压S
	UAc3                        Num    `json:"uAc3"`                        // 交流电压T
	InverterTemperature         Num    `json:"inverterTemperature"`         // 逆变器温度
	AcOutputType                int    `json:"acOutputType"`                // 交流输出类型
	DcInputType                 int    `json:"dcInputType"`                 // 直流输入类型
	PowerFactor                 Num    `json:"powerFactor"`                 // 功率因数
	BatteryCapacitySoc          Num    `json:"batteryCapacitySoc"`          // 电池SOC
	BatteryHealthSoh            Num    `json:"batteryHealthSoh"`            // 电池SOH
	SocDischargeSet             Num    `json:"socDischargeSet"`             // 过放SOC
	SocChargingSet              Num    `json:"socChargingSet"`              // 强充SOC
	BatteryVoltage              Num    `json:"batteryVoltage"`              // 电池电压
	BstteryCurrent              Num    `json:"bstteryCurrent"`              // 电池电流
	BatteryPower                Num    `json:"batteryPower"`                // 电池功率
	BatteryTodayChargeEnergy    Num    `json:"batteryTodayChargeEnergy"`    // 当日电池充电电量
	BatteryTotalChargeEnergy    Num    `json:"batteryTotalChargeEnergy"`    // 电池总充电电量
	BatteryTodayDischargeEnergy Num    `json:"batteryTodayDischargeEnergy"` // 当日电池放电电量
	BatteryTotalDischargeEnergy Num    `json:"batteryTotalDischargeEnergy"` // 电池总放电量
	GridPurchasedTodayEnergy    Num    `json:"gridPurchasedTodayEnergy"`    // 当日电表买电
	GridPurchasedTotalEnergy    Num    `json:"gridPurchasedTotalEnergy"`    // 累计电表买电
	GridSellTodayEnergy         Num    `json:"gridSellTodayEnergy"`         // 当日电表卖电
	GridSellTotalEnergy         Num    `json:"gridSellTotalEnergy"`         // 累计电表卖电
	FamilyLoadPower             Num    `json:"familyLoadPower"`             // 家庭负载功率
	BypassLoadPower             Num    `json:"bypassLoadPower"`             // 旁路负载功率
	PSum                        Num    `json:"pSum"`                        // 电网总有功功率
	HomeLoadTodayEnergy         Num    `json:"homeLoadTodayEnergy"`         // 当日负载用电电量
	HomeLoadTotalEnergy         Num    `json:"homeLoadTotalEnergy"`         // 累计负载用电量
}

// InverterDay 获取单台逆变器某日的实时数据
func (c *SolisSDK) InverterDay(req InverterDayRequest) ([]InverterDayItem, error) {
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("获取单台逆变器某日实时数据: money 不能为空")
	}
	if strings.TrimSpace(req.Time) == "" {
		return nil, fmt.Errorf("获取单台逆变器某日实时数据: time 不能为空")
	}
	if req.TimeZone == 0 {
		return nil, fmt.Errorf("获取单台逆变器某日实时数据: timeZone 不能为空")
	}
	var out []InverterDayItem
	if err := c.do(PathInverterDay, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 3.5 获取单台逆变器某月的日数据 /v1/api/inverterMonth ----------

// InverterMonthRequest 逆变器某月日数据请求
type InverterMonthRequest struct {
	ID       Int64  `json:"id,omitempty"`       // 逆变器id
	SN       string `json:"sn,omitempty"`       // 逆变器SN
	Money    string `json:"money"`              // 收益货币单位
	Month    string `json:"month"`              // 查询月份 yyyy-MM
	TimeZone int    `json:"timeZone,omitempty"` // 电站时区
}

// InverterMonthItem 逆变器某月日数据项
type InverterMonthItem struct {
	rawHolder
	Energy                 Num    `json:"energy"`                 // 发电量
	EnergyStr              string `json:"energyStr"`              // 发电量单位
	Date                   Int64  `json:"date"`                   // 时间戳
	DateStr                string `json:"dateStr"`                // 时间字符串
	Money                  Num    `json:"money"`                  // 收益
	MoneyStr               string `json:"moneyStr"`               // 收益单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电
}

// InverterMonth 获取单台逆变器某月的日数据
func (c *SolisSDK) InverterMonth(req InverterMonthRequest) ([]InverterMonthItem, error) {
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("获取单台逆变器某月日数据: money 不能为空")
	}
	if strings.TrimSpace(req.Month) == "" {
		return nil, fmt.Errorf("获取单台逆变器某月日数据: month 不能为空")
	}
	var out []InverterMonthItem
	if err := c.do(PathInverterMonth, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 3.6 获取单台逆变器某年的月数据 /v1/api/inverterYear ----------

// InverterYearRequest 逆变器某年月数据请求
type InverterYearRequest struct {
	ID       Int64  `json:"id,omitempty"`       // 逆变器id
	SN       string `json:"sn,omitempty"`       // 逆变器SN
	Money    string `json:"money"`              // 收益货币单位
	Year     string `json:"year"`               // 查询年份 yyyy
	TimeZone int    `json:"timeZone,omitempty"` // 电站时区
}

// InverterYearItem 逆变器某年月数据项
type InverterYearItem struct {
	rawHolder
	Energy                 Num    `json:"energy"`                 // 发电量
	EnergyStr              string `json:"energyStr"`              // 发电量单位
	Date                   Int64  `json:"date"`                   // 时间戳
	DateStr                string `json:"dateStr"`                // 时间字符串
	Money                  Num    `json:"money"`                  // 收益
	MoneyStr               string `json:"moneyStr"`               // 收益单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电
}

// InverterYear 获取单台逆变器某年的月数据
func (c *SolisSDK) InverterYear(req InverterYearRequest) ([]InverterYearItem, error) {
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("获取单台逆变器某年月数据: money 不能为空")
	}
	if strings.TrimSpace(req.Year) == "" {
		return nil, fmt.Errorf("获取单台逆变器某年月数据: year 不能为空")
	}
	var out []InverterYearItem
	if err := c.do(PathInverterYear, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 3.7 获取单台逆变器的年数据 /v1/api/inverterAll ----------

// InverterAllRequest 逆变器年数据请求
type InverterAllRequest struct {
	ID       Int64  `json:"id,omitempty"`       // 逆变器id
	SN       string `json:"sn,omitempty"`       // 逆变器SN
	Money    string `json:"money"`              // 收益货币单位
	TimeZone int    `json:"timeZone,omitempty"` // 电站时区
}

// InverterAllResult 逆变器年数据结果
type InverterAllResult struct {
	rawHolder
	Year                   int    `json:"year"`                   // 年
	Energy                 Num    `json:"energy"`                 // 发电量
	EnergyStr              string `json:"energyStr"`              // 发电量单位
	BatteryDischargeEnergy Num    `json:"batteryDischargeEnergy"` // 电池放电电量
	BatteryChargeEnergy    Num    `json:"batteryChargeEnergy"`    // 电池充电电量
	GridPurchasedEnergy    Num    `json:"gridPurchasedEnergy"`    // 电表买电
	GridSellEnergy         Num    `json:"gridSellEnergy"`         // 电表卖电
}

// InverterAll 获取单台逆变器的年数据
func (c *SolisSDK) InverterAll(req InverterAllRequest) (*InverterAllResult, error) {
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("获取单台逆变器年数据: money 不能为空")
	}
	var out InverterAllResult
	if err := c.do(PathInverterAll, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 3.8 获取多台逆变器质保数据 /v1/api/inverter/shelfTime ----------

// InverterShelfTimeRequest 多台逆变器质保数据请求
type InverterShelfTimeRequest struct {
	PageRequest        // 分页参数
	SN          string `json:"sn,omitempty"` // 逆变器SN，多个逗号分隔
}

// InverterShelfTimeResult 多台逆变器质保数据结果
type InverterShelfTimeResult struct {
	rawHolder
	Page Page[InverterShelfTimeItem] `json:"page"` // 分页列表
}

// UnmarshalJSON 文档表格写 page 嵌套、官方示例是 data 下平铺，两种形态都兼容
func (r *InverterShelfTimeResult) UnmarshalJSON(data []byte) error {
	return unmarshalPaged(data, &r.Page)
}

// InverterShelfTimeItem 逆变器质保数据项
type InverterShelfTimeItem struct {
	rawHolder
	ID             Int64  `json:"id"`             // 逆变器id
	SN             string `json:"sn"`             // 逆变器SN
	ShelfBeginTime Int64  `json:"shelfBeginTime"` // 质保开始时间
	ShelfEndTime   Int64  `json:"shelfEndTime"`   // 质保结束时间
	ShelfTime      Num    `json:"shelfTime"`      // 质保年限
	ShelfState     Str    `json:"shelfState"`     // 质保状态
}

// InverterShelfTime 获取多台逆变器质保数据
func (c *SolisSDK) InverterShelfTime(req InverterShelfTimeRequest) (*InverterShelfTimeResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out InverterShelfTimeResult
	if err := c.do(PathInverterShelfTime, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 3.9 获取账号下设备报警列表 /v1/api/alarmList ----------

// AlarmListRequest 账号下设备报警列表请求
type AlarmListRequest struct {
	PageRequest           // 分页参数
	StationID      Int64  `json:"stationId,omitempty"`      // 电站id
	AlarmDeviceSn  string `json:"alarmDeviceSn,omitempty"`  // 逆变器SN
	AlarmBeginTime string `json:"alarmBeginTime,omitempty"` // 报警开始日期
	AlarmEndTime   string `json:"alarmEndTime,omitempty"`   // 报警结束日期
	NmiCode        string `json:"nmiCode,omitempty"`        // nmi编码
	State          int    `json:"State,omitempty"`          // 报警状态筛选
}

// AlarmListResult 账号下设备报警列表结果
type AlarmListResult struct {
	rawHolder
	Page Page[AlarmListItem] `json:"page"` // 分页列表
}

// UnmarshalJSON 文档表格写 page 嵌套、官方示例是 data 下平铺，两种形态都兼容
func (r *AlarmListResult) UnmarshalJSON(data []byte) error {
	return unmarshalPaged(data, &r.Page)
}

// AlarmListItem 设备报警列表项
type AlarmListItem struct {
	rawHolder
	StationID       Int64  `json:"stationId"`       // 电站id
	StationName     string `json:"stationName"`     // 电站名称
	AlarmDeviceSn   string `json:"alarmDeviceSn"`   // 逆变器SN
	AlarmCode       string `json:"alarmCode"`       // 报警代码
	AlarmLevel      string `json:"alarmLevel"`      // 报警等级
	AlarmBeginTime  Int64  `json:"alarmBeginTime"`  // 报警开始时间
	AlarmEndTime    Int64  `json:"alarmEndTime"`    // 报警结束时间
	AlarmMsg        string `json:"alarmMsg"`        // 报警内容
	Advice          string `json:"advice"`          // 报警处理建议
	State           string `json:"state"`           // 报警状态
	WarningInfoData int    `json:"warningInfoData"` // 子报警代码
}

// AlarmList 获取账号下设备报警列表
func (c *SolisSDK) AlarmList(req AlarmListRequest) (*AlarmListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out AlarmListResult
	if err := c.do(PathAlarmList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
