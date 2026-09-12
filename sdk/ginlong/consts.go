package ginlong

// 平台常量
const (
	PlatformName   = "Ginlong"
	DefaultBaseURL = "https://api.ginlong.com:13333"
)

// 签名与请求头常量
const (
	// SignContentType 参与签名的 Content-Type，官方示例中固定为 application/json
	SignContentType = "application/json"
	// HeaderContentType 请求头 Content-Type，与签名值保持一致最稳妥
	HeaderContentType = "application/json;charset=UTF-8"
	// TimeLayout Date 头格式，GMT 时区，等价于 Java "EEE, dd MMM yyyy HH:mm:ss 'GMT'"
	TimeLayout = "Mon, 02 Jan 2006 15:04:05 GMT"
	// MaxPageSize 单页最大条数
	MaxPageSize = 100
	// DefaultPageSize 默认单页条数
	DefaultPageSize = 20
)

// 请求路径
const (
	PathInverterList       = "/v1/api/inverterList"       // 账号下逆变器列表
	PathInverterDetail     = "/v1/api/inverterDetail"     // 单台逆变器详情
	PathInverterDetailList = "/v1/api/inverterDetailList" // 多台逆变器详情
	PathInverterDay        = "/v1/api/inverterDay"        // 逆变器某日实时数据
	PathInverterMonth      = "/v1/api/inverterMonth"      // 逆变器某月日数据
	PathInverterYear       = "/v1/api/inverterYear"       // 逆变器某年月数据
	PathInverterAll        = "/v1/api/inverterAll"        // 逆变器年数据
	PathInverterShelfTime  = "/v1/api/inverter/shelfTime" // 多台逆变器质保数据
	PathAlarmList          = "/v1/api/alarmList"          // 设备报警列表

	PathCollectorList   = "/v1/api/collectorList"   // 采集器列表
	PathCollectorDetail = "/v1/api/collectorDetail" // 单台采集器详情
	PathCollectorDay    = "/v1/api/collector/day"   // 单台采集器信号值

	PathEpmList   = "/v1/api/epmList"   // EPM 列表
	PathEpmDetail = "/v1/api/epmDetail" // 单台 EPM 详情
	PathEpmDay    = "/v1/api/epm/day"   // EPM 某日实时数据
	PathEpmMonth  = "/v1/api/epm/month" // EPM 某月日数据
	PathEpmYear   = "/v1/api/epm/year"  // EPM 某年月数据
	PathEpmAll    = "/v1/api/epm/all"   // EPM 年数据

	PathWeatherList   = "/v1/api/weatherList"   // 气象仪列表
	PathWeatherDetail = "/v1/api/weatherDetail" // 单台气象仪详情

	PathAmmeterList   = "/v1/api/ammeterList"   // 电表列表
	PathAmmeterDetail = "/v1/api/ammeterDetail" // 单台电表详情

	PathUserStationList         = "/v1/api/userStationList"         // 账号下电站列表
	PathStationDetail           = "/v1/api/stationDetail"           // 单个电站详情
	PathStationDetailList       = "/v1/api/stationDetailList"       // 多个电站详情
	PathStationDayEnergyList    = "/v1/api/stationDayEnergyList"    // 多电站某日实时数据
	PathStationMonthEnergyList  = "/v1/api/stationMonthEnergyList"  // 多电站月数据
	PathStationYearEnergyList   = "/v1/api/stationYearEnergyList"   // 多电站年数据
	PathStationDay              = "/v1/api/stationDay"              // 单电站某日实时数据
	PathStationMonth            = "/v1/api/stationMonth"            // 单电站某月日数据
	PathStationYear             = "/v1/api/stationYear"             // 单电站某年月数据
	PathStationAll              = "/v1/api/stationAll"              // 单电站年数据
	PathAddStation              = "/v1/api/addStation"              // 新增电站绑定逆变器SN
	PathStationUpdate           = "/v1/api/stationUpdate"           // 修改电站信息
	PathAddStationBindCollector = "/v1/api/addStationBindCollector" // 新增电站绑定新采集器
	PathDelCollector            = "/v1/api/delCollector"            // 电站解绑采集器
	PathAddDevice               = "/v1/api/addDevice"               // 电站绑定逆变器
)

// 返回码（附录 1）
const (
	CodeOK             = "0"     // 成功
	CodeNoPermission   = "R0000" // 无权限操作
	CodeAlreadyBound   = "B0001" // 已绑定其他用户
	CodeSNRequired     = "I0003" // 请输入 SN 号
	CodeCollectorGone  = "B0049" // 该采集器已不存在或无权限
	CodeParamEmpty     = "I0000" // 必要参数为空
	CodeUserNotExist   = "B0011" // 该用户不存在
	CodeBadCredentials = "I0012" // 账号或者密码错误
)

// codeHints 返回码提示
var codeHints = map[string]string{
	CodeOK:             "成功",
	CodeNoPermission:   "无权限操作，请检查 API 账户权限",
	CodeAlreadyBound:   "设备已绑定其他用户",
	CodeSNRequired:     "请输入 SN 号",
	CodeCollectorGone:  "该采集器已不存在或无权限，无法查看",
	CodeParamEmpty:     "必要参数为空",
	CodeUserNotExist:   "该用户不存在",
	CodeBadCredentials: "账号或者密码错误，请重新输入",
}

// CodeHint 返回码排查提示
func CodeHint(code string) string {
	if h, ok := codeHints[code]; ok {
		return h
	}
	return "详见锦浪云平台 API 文档附录 1 错误码"
}

// DeviceState 设备/电站状态
type DeviceState int

const (
	StateOnline  DeviceState = 1 // 在线
	StateOffline DeviceState = 2 // 离线
	StateAlarm   DeviceState = 3 // 报警
)

func (s DeviceState) String() string {
	switch s {
	case StateOnline:
		return "在线"
	case StateOffline:
		return "离线"
	case StateAlarm:
		return "报警"
	default:
		return "未知"
	}
}

// StationType 电站类型（附录 2）
type StationType int

const (
	StationGridTied       StationType = 0  // 并网电站
	StationStorage        StationType = 1  // 储能电站
	StationACCouple       StationType = 2  // ACCouple 电站
	StationEPM            StationType = 3  // EPM 电站（并网+电表）
	StationInnerMeter     StationType = 4  // 内置电表（并网+电表）
	StationOuterMeter     StationType = 5  // 外置电表（显示电表）
	StationS5OffGrid      StationType = 6  // S5 离网并机储能
	StationS5GridTied     StationType = 7  // S5 并网并机储能
	StationGridACCouple   StationType = 8  // 并网+ACCouple 电站
	StationOffGridStorage StationType = 9  // 离网储能
	StationS6GridTied     StationType = 10 // S6 并网并机储能
	StationS6OffGrid      StationType = 11 // S6 离网并机储能
)

// SynchronizationType 并网类型
type SynchronizationType int

const (
	SyncFullToGrid   SynchronizationType = 0 // 全额上网
	SyncSelfConsumed SynchronizationType = 1 // 自发自用
	SyncOffGrid      SynchronizationType = 2 // 离网
)

// MeterType 逆变器电表类型（附录 3）
type MeterType int

const (
	MeterGridTied          MeterType = 1    // 并网
	MeterGridLoadSide      MeterType = 2    // 并网加负载侧电表
	MeterGridGridSide      MeterType = 3    // 并网加电网侧电表
	MeterStorageLoadSide   MeterType = 4    // 储能加负载侧电表
	MeterStorageGridSide   MeterType = 5    // 储能加电网侧电表
	MeterReserved          MeterType = 6    // 保留
	MeterOffGridStorage    MeterType = 7    // 离网储能
	MeterGridStorageDouble MeterType = 8    // 并网储能双电表
	MeterACCoupleNoCT      MeterType = 1001 // ACCouple（不带 CT）
	MeterACCoupleWithCT    MeterType = 1002 // ACCouple（带 CT）
)

// 以下枚举都兼容 数值 / 字符串 两种返回形态

func (s *DeviceState) UnmarshalJSON(b []byte) error {
	var e intEnum
	*s = DeviceState(e.unmarshal(b))
	return nil
}

func (s *StationType) UnmarshalJSON(b []byte) error {
	var e intEnum
	*s = StationType(e.unmarshal(b))
	return nil
}

func (s *SynchronizationType) UnmarshalJSON(b []byte) error {
	var e intEnum
	*s = SynchronizationType(e.unmarshal(b))
	return nil
}

func (s *MeterType) UnmarshalJSON(b []byte) error {
	var e intEnum
	*s = MeterType(e.unmarshal(b))
	return nil
}
