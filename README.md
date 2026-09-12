# ps-sdk

各光伏云平台 API 实现集合。

| 平台 | 包 | 文档依据 | 鉴权方式 |
| --- | --- | --- | --- |
| 华为 FusionSolar SmartPVMS | `sdk/huawei` | SmartPVMS 26.2.0 北向接口参考 | API 账户登录 + XSRF-TOKEN |
| 锦浪 Ginlong / SolisCloud | `sdk/ginlong` | 锦浪云平台 API 文档 V2.0.3 | KeyID/KeySecret + HMAC-SHA1 签名 |
| SolarMan 小麦智电 / 小麦商家版 | `sdk/solarman` | SolarMan OpenApi 在线文档 V2.0.3 | OAuth2 Bearer Token（60 天） |

## 目录结构

```
sdk/huawei/
  client.go               客户端、凭证、Token 管理、请求内核（305 重登录 / 限流重试）
  consts.go               路径、错误码、设备类型、业务枚举
  types.go                统一响应封装、分页、ItemMap、时间与工具函数
  kpi.go                  电站实时/报表指标常量与解析
  api_station.go          电站列表、设备列表
  api_monitor.go          电站实时、设备实时、设备历史
  api_alarm.go            活动告警
  api_report.go           电站小时/日/月/年、设备日/月/年
  api_config.go           储能工作模式查询
  api_control.go          储能充放电任务下发/查询（v2）
  api_control_battery.go  储能工作模式设置、储能参数设置（下发/查询）
  api_control_power.go    逆变器有功功率设置（下发/查询）、有功功率控制模式查询
  api_control_dispatch.go 储能调度充放电任务（下发/查询）

sdk/ginlong/
  client.go               客户端、凭证、HMAC-SHA1 签名、请求内核
  consts.go               路径、返回码、状态与类型枚举
  types.go                响应封装、Num/Int64/Str 容错标量、分页、Chart、ItemMap
  api_inverter.go         逆变器列表/详情/批量详情/日/月/年/年数据/质保、报警列表
  api_collector.go        采集器列表/详情/信号值、气象仪列表/详情
  api_epm.go              EPM 列表/详情/日/月/年/年数据
  api_ammeter.go          电表列表/详情
  api_station.go          电站列表/详情/批量详情、日/月/年/年数据
  api_station_write.go    新增电站、修改电站、绑定采集器、解绑采集器、绑定逆变器

sdk/solarman/
  client.go               客户端 SolarmanSDK、凭证、Token 管理、请求内核（三种鉴权）
  consts.go               数据中心地址、全部路径常量、响应码、枚举
  types.go                Response 公共响应、Num/Int64/Str 容错标量、分页、ItemMap
  api_account.go          2.1~2.11 账号接口 + 5.1 验证码 + 13 APPID 余量
  api_device.go           3.1~3.9 设备接口 + 4.13 添加网关、4.14 删除设备
  api_station.go          4.1~4.12、4.16 电站查询与写操作
  api_weather.go          4.17~4.19 天气接口

cmd/ginlong/main.go       锦浪调用示例（凭证留空，测试时填写）
cmd/solarman/main.go      SolarMan 调用示例（凭证留空，测试时填写）
main.go                   华为调用示例
```

---

# 华为 FusionSolar SmartPVMS

## 快速开始

```go
sdk, err := huawei.NewFusionSolarSDK(
    huawei.Credentials{
        UserName:   "apiUser",              // API 账户名
        SystemCode: "password",             // API 账户密码
        BaseURL:    "https://xxx.fusionsolar.huawei.com", // 缺省为 intl 站点
    },
    huawei.WithTimeout(30*time.Second),
    huawei.WithDebugf(log.Printf),          // 可选：打印请求/响应
)
if err != nil {
    log.Fatal(err)
}

// 登录是隐式的：首次调用任意接口时自动登录
stations, err := sdk.GetStationList(huawei.StationListRequest{PageNo: 1})
```

### 可选项

| Option | 说明 |
| --- | --- |
| `WithClient(*req.Client)` | 自定义 HTTP 客户端（代理、TLS 等） |
| `WithTimeout(d)` | 单次请求超时，默认 30s |
| `WithToken(token)` | 复用已有 XSRF-TOKEN，跳过首次登录 |
| `WithDebugf(f)` | 调试日志 |
| `WithDevMode(true)` | 打印完整请求/响应明细 |
| `WithMaxAttempts(n)` | 含首调在内的最大尝试次数，默认 3 |
| `WithRetryBaseDelay(d)` | 重试退避基数，默认 2s |
| `WithDisableAutoLogin()` | 关闭自动登录，token 失效时直接返回 305 |

## 认证与 Token

- 调用 `/thirdData/login` 登录成功后在**响应头**返回 `XSRF-TOKEN`，有效期 **30 分钟**；期间持续调用会自动续期。SDK 在 token 剩余不足 5 分钟时主动重新登录。
- 同一 API 账户**只允许一个在线会话**，重复登录会使先前 token 失效。SDK 用互斥锁串行化登录；**多实例部署时请共享 token 或改用固定实例**，否则会互相踢下线。
- 任意接口返回 `failCode=305` 时，SDK 会自动重新登录并重试一次。
- 连续 5 次密码错误会锁定账户 30 分钟；登录接口限流为每 10 分钟 5 次。

## 错误处理

```go
stations, err := sdk.GetStationList(req)
if err != nil {
    var apiErr *huawei.APIError
    if errors.As(err, &apiErr) {
        log.Printf("failCode=%d hint=%s retryable=%v",
            apiErr.FailCode, apiErr.FailCode.Hint(), apiErr.Retryable())
    }
}
```

- `*APIError` —— 业务错误（`failCode != 0`），`Hint()` 给出中文排查提示。
- `*BadResponseError` —— 响应不是预期的 JSON（如网关 HTML 错误页）。
- `*Exception` —— 网关参数校验异常（`exceptionId` 形式）。
- 传输层错误直接透传，已用 `%w` 包装。

**日志与限流**：`407`（单用户限流）、`429`/`20200`/`20004`（系统级）会自动退避重试，其中 `429` 固定等待 1 分钟后重试；`305` 触发重新登录。其余业务错误不重试。

## 接口清单

### 基础类

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `GetStationList` | 电站列表（分页，每页 ≤100） | `POST /thirdData/stations` |
| `GetDevList` | 设备列表（≤100 个电站） | `POST /thirdData/getDevList` |
| `GetStationRealKpi` | 电站实时数据（≤100 个电站） | `POST /thirdData/getStationRealKpi` |
| `GetDevRealKpi` | 设备实时数据（1 种类型 ×≤100 台） | `POST /thirdData/getDevRealKpi` |
| `GetDevHistory` | 设备历史数据（5 分钟粒度，1 台 ≤24h） | `POST /rest/openapi/pvms/nbi/v1/device/history` |
| `GetAlarmList` | 活动告警（不含历史告警） | `POST /thirdData/getAlarmList` |
| `GetStationKpiHour` | 电站小时数据 | `POST /thirdData/getKpiStationHour` |
| `GetStationKpiDay` | 电站日数据 | `POST /thirdData/getKpiStationDay` |
| `GetStationKpiMonth` | 电站月数据 | `POST /thirdData/getKpiStationMonth` |
| `GetStationKpiYear` | 电站年数据 | `POST /thirdData/getKpiStationYear` |
| `GetDevKpiDay` | 设备日数据 | `POST /thirdData/getDevKpiDay` |
| `GetDevKpiMonth` | 设备月数据 | `POST /thirdData/getDevKpiMonth` |
| `GetDevKpiYear` | 设备年数据 | `POST /thirdData/getDevKpiYear` |
| `GetBatteryMode` | 储能工作模式查询 | `POST .../nbi/v1/configuration/battery-mode` |
| `ActivePowerControlMode` | 逆变器有功功率控制模式查询 | `POST .../nbi/v1/configuration/active-power-control-mode` |

### 控制类

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `SubmitChargeDischargeTask` | 储能充放电任务下发（≤100 电站） | `POST .../nbi/v2/control/charge-and-discharge/async-task` |
| `GetChargeDischargeStatus` | 储能充放电任务查询 | `POST /rest/openapi/pvms/v1/vpp/chargeAndDischargeStatus` |
| `BatteryModeTask` | 储能工作模式设置下发（≤10 电站） | `POST .../nbi/v1/control/battery/mode/async-task` |
| `BatteryModeTaskInfo` | 储能工作模式设置查询 | `POST .../nbi/v1/control/battery/mode/task-info` |
| `BatteryConfigTask` | 储能参数设置下发（≤10 电站） | `POST .../nbi/v1/control/battery/configuration/async-task` |
| `BatteryConfigTaskInfo` | 储能参数设置查询 | `POST .../nbi/v1/control/battery/configuration/task-info` |
| `ActivePowerControlTask` | 逆变器有功功率设置下发（≤10 电站） | `POST .../nbi/v2/control/active-power-control/async-task` |
| `ActivePowerControlTaskInfo` | 逆变器有功功率设置查询 | `POST .../nbi/v2/control/active-power-control/task-info` |
| `BatteryDispatchTask` | 储能调度充放电任务下发 | `POST .../nbi/v1/control/battery/battery-dispatch/async-task` |
| `BatteryDispatchTaskInfo` | 储能调度充放电任务查询 | `POST .../nbi/v1/control/battery/battery-dispatch/task-info` |

控制类接口的 `failCode` 约定：`0` 成功、`1` 部分成功（**不算错误**，结果正常返回）、`2` 失败。SDK 在 `doTask` 中对 `1` 不报错；`2` 时返回 `*APIError`，但**结果结构体仍已填充**，可用于定位失败电站。

## 数据项（KPI）

华为把实时/报表指标放在 `dataItemMap` / `dataItems` 里，值可能是**字符串或 `null`**。SDK 用 `ItemMap` 承载并提供容错取值：

```go
kpi := stationRealKpi.KPI()               // 已解析的电站实时指标
dayPower, ok := item.DataItemMap.Float64(huawei.KpiDayPower)
```

- `ItemMap.Float64 / Int / Int64 / String / Bool / Raw`
- `ItemMap.ParseStationRealtime()` → `StationRealtime`
- `ItemMap.ParseStationReport()` → `StationReport`
- 指标 key 常量见 `kpi.go`（`KpiDayPower`、`KpiPVYield`、`KpiRealHealthState` 等）。

设备级指标（逆变器、储能、电表等）key 较多且随设备类型变化，按需用 `ItemMap` 取值，key 名称见文档「设备实时数据列表」「设备历史数据列表」。

## 数值字段容错

华为文档标注为 `Double` 的字段**实际可能以字符串返回**（实测：电站列表的 `longitude`/`latitude` 返回 `"118.420860"`）。因此 SDK 中所有 `Double` 字段使用 `huawei.Float64` 类型，它兼容三种形态：

```jsonc
{"longitude": 118.42086}    // 数值
{"longitude": "118.42086"}  // 字符串（实测形态）
{"longitude": null}         // 空值 → 0
```

`Float64` 底层就是 `float64`，序列化时输出数值，所以构造请求不受影响：

```go
lat := station.Latitude.Float()          // 取 float64
fmt.Printf("%.4f", station.Latitude)     // 直接打印也可以
req := huawei.BatteryConfigTaskRequest{Tasks: []huawei.BatteryConfigTaskItem{{
    PlantCode: "NE=1",
    BatteryConfigurationInfo: huawei.BatteryConfigurationInfo{
        EndOfChargeSoc: huawei.Float64Ptr(98), // 可选数值字段用 Float64Ptr
    },
}}}
```

`ItemMap` 已内置同样的容错（`Float64`/`Int`/`Int64`/`String`/`Bool`），无需额外处理。

## 辅助工具

```go
huawei.JoinCodes(codes)      // 拼成英文逗号分隔，超过 100 个自动截断
huawei.SplitCodes(s)         // 拆分并去空
huawei.Millis(t)             // time.Time → 毫秒时间戳
huawei.Float64Ptr(v)         // 可选数值字段取指针
type huawei.Float64 float64  // 兼容 数值/字符串/null 的 Double
type huawei.Time struct{ time.Time } // 解析 2020-02-06T00:00:00+08:00，兼容 null
```

## 注意事项

- 单次查询上限：电站列表/设备列表/实时/报表类为 **100** 个电站或设备；控制类下发为 **10** 个电站（充放电任务为 100 个电站）。
- 时间参数一律为**毫秒时间戳**；文档中带时区的字符串时间用 `huawei.Time` 接收。
- 报表接口的 `collectTime` 决定统计范围：小时=该自然日、日=该自然月、月=该自然年。
- 控制类接口会变更设备运行参数，请先确认设备在线并注意调用频率。

---

# 锦浪 Ginlong / SolisCloud

依据《锦浪云平台 API 文档 V2.0.3》，覆盖文档全部 **37 个接口**（设备 22 个 + 电站 15 个）。

## 快速开始

```go
c := ginlong.NewClient(ginlong.Credentials{
    APIID:     "",   // KeyID，锦浪云 WEB 端「服务 - API 管理」获取
    APISecret: "",   // KeySecret
    // BaseURL: ginlong.DefaultBaseURL, // 缺省 https://api.ginlong.com:13333
}, ginlong.WithDebugf(log.Printf))

// 账号下电站列表
res, err := c.UserStationList(ginlong.UserStationListRequest{
    PageRequest: ginlong.PageRequest{PageNo: 1, PageSize: 20},
})
if err != nil {
    log.Fatal(err)
}
for _, s := range res.Page.Records {
    fmt.Println(s.ID.Int(), s.StationName, s.DayEnergy.Float())
}
```

凭证允许留空构造，只有**真正发起请求时**才会返回 `ginlong.ErrMissingCredentials`；也可以用 `c.SetCredentials(...)` 后续补充。

### 可选项

| Option | 说明 |
| --- | --- |
| `WithClient(*req.Client)` | 自定义 HTTP 客户端 |
| `WithTimeout(d)` | 单次请求超时，默认 30s |
| `WithDebugf(f)` | 调试日志 |
| `WithDevMode(true)` | 打印完整请求/响应明细 |
| `WithContentType(ct)` | 覆盖请求头与签名用的 Content-Type，默认 `application/json` |
| `WithClock(now)` | 注入时钟，便于测试签名 |

## 签名与请求头

每个请求都必须带 `Content-MD5` / `Content-Type` / `Date` / `Authorization` 四个头：

```
Content-MD5   = base64(md5(body))
Date          = GMT 时间，格式 Mon, 02 Jan 2006 15:04:05 GMT（与当前时间相差不能超过 ±15 分钟）
Authorization = "API " + apiId + ":" + sign
sign          = base64(HmacSHA1(apiSecret,
                  "POST\n" + Content-MD5 + "\n" + application/json + "\n" + Date + "\n" + 接口路径))
```

SDK 已封装，直接调用方法即可。两处容易踩坑的细节已按**官方 Java 示例 + 文档调用实例**对齐：

- 参与签名的 Content-Type 是**固定值 `application/json`**，不是请求头里那个 `application/json;charset=UTF-8`。SDK 默认把请求头也设成 `application/json`，保证「服务端按固定值算」和「服务端按头值算」两种实现都能通过；如遇异常可用 `WithContentType("application/json;charset=UTF-8")` 切换。
- `Authorization` 中 `API` 与 apiId 之间有**一个空格**（文档 4.1 的调用实例即为 `Authorization: API 1300386381676565707:...`）。

`ginlong.ContentMD5` / `ginlong.Sign` / `ginlong.Authorization` 都是导出函数，可以单独用来做联调自检。

## 错误处理

```go
res, err := c.InverterList(req)
if err != nil {
    var apiErr *ginlong.APIError
    if errors.As(err, &apiErr) {
        log.Printf("code=%s msg=%s hint=%s",
            apiErr.Code, apiErr.Msg, ginlong.CodeHint(string(apiErr.Code)))
    }
}
```

返回体固定为 `{success, code, msg, data}`，`code` 为 **String**（`"0"` 表示成功），SDK 的 `Envelope.Code` 同时兼容字符串与数值两种形态。

| 错误类型 | 说明 |
| --- | --- |
| `*APIError` | 业务错误（`code != "0"`），`Hint()` 给出附录 1 的中文提示 |
| `*BadResponseError` | 响应不是预期的 JSON |
| `ErrMissingCredentials` | 凭证未填写 |

附录 1 返回码：`R0000` 无权限操作、`B0001` 已绑定其他用户、`I0003` 请输入 SN、`B0049` 采集器不存在或无权限、`I0000` 必要参数为空、`B0011` 用户不存在、`I0012` 账号或密码错误。

## 接口清单

### 设备接口

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `InverterList` | 账号下逆变器列表 | `/v1/api/inverterList` |
| `InverterDetail` | 单台逆变器详情 | `/v1/api/inverterDetail` |
| `InverterDetailList` | 多台逆变器详情 | `/v1/api/inverterDetailList` |
| `InverterDay` | 逆变器某日实时数据 | `/v1/api/inverterDay` |
| `InverterMonth` | 逆变器某月日数据 | `/v1/api/inverterMonth` |
| `InverterYear` | 逆变器某年月数据 | `/v1/api/inverterYear` |
| `InverterAll` | 逆变器年数据 | `/v1/api/inverterAll` |
| `InverterShelfTime` | 多台逆变器质保数据 | `/v1/api/inverter/shelfTime` |
| `AlarmList` | 账号下设备报警列表 | `/v1/api/alarmList` |
| `CollectorList` | 账号下采集器列表 | `/v1/api/collectorList` |
| `CollectorDetail` | 单台采集器详情 | `/v1/api/collectorDetail` |
| `CollectorDay` | 单台采集器信号值 | `/v1/api/collector/day` |
| `EpmList` | 账号下 EPM 列表 | `/v1/api/epmList` |
| `EpmDetail` | 单台 EPM 详情 | `/v1/api/epmDetail` |
| `EpmDay` | EPM 某日实时数据 | `/v1/api/epm/day` |
| `EpmMonth` | EPM 某月日数据 | `/v1/api/epm/month` |
| `EpmYear` | EPM 某年月数据 | `/v1/api/epm/year` |
| `EpmAll` | EPM 年数据 | `/v1/api/epm/all` |
| `WeatherList` | 账号下气象仪列表 | `/v1/api/weatherList` |
| `WeatherDetail` | 单台气象仪详情 | `/v1/api/weatherDetail` |
| `AmmeterList` | 账号下电表列表 | `/v1/api/ammeterList` |
| `AmmeterDetail` | 单台电表详情 | `/v1/api/ammeterDetail` |

### 电站接口

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `UserStationList` | 账号下电站列表 | `/v1/api/userStationList` |
| `StationDetail` | 单个电站详情 | `/v1/api/stationDetail` |
| `StationDetailList` | 多个电站详情 | `/v1/api/stationDetailList` |
| `StationDayEnergyList` | 多电站某日实时数据 | `/v1/api/stationDayEnergyList` |
| `StationMonthEnergyList` | 多电站月数据 | `/v1/api/stationMonthEnergyList` |
| `StationYearEnergyList` | 多电站年数据 | `/v1/api/stationYearEnergyList` |
| `StationDay` | 单电站某日实时数据 | `/v1/api/stationDay` |
| `StationMonth` | 单电站某月日数据 | `/v1/api/stationMonth` |
| `StationYear` | 单电站某年月数据 | `/v1/api/stationYear` |
| `StationAll` | 单电站年数据 | `/v1/api/stationAll` |
| `AddStation` | 新增电站绑定逆变器 SN | `/v1/api/addStation` |
| `StationUpdate` | 修改电站信息 | `/v1/api/stationUpdate` |
| `AddStationBindCollector` | 新增电站绑定新采集器 | `/v1/api/addStationBindCollector` |
| `DelCollector` | 电站解绑采集器 | `/v1/api/delCollector` |
| `AddDevice` | 电站绑定逆变器 | `/v1/api/addDevice` |

## 类型与容错

锦浪的字段类型标注与真实返回**经常不一致**，SDK 用容错标量兜住：

| 类型 | 兼容形态 | 用法 |
| --- | --- | --- |
| `Num` | 数值 / 字符串 / null（如 `"1"`、`"null"`） | `.Float()` |
| `Int64` | 数值 / 字符串 / 空串 / null | `.Int()`，序列化输出**字符串** |
| `Str` | 字符串 / 数值 / null | `.String()` |

实测证据：`id` / `stationId` / `userId` / `collectorId` / `dataTimestamp` 在文档的所有响应示例里都是字符串（`"id":"1308675217944611083"`、`"collectorId":""`），而文档表格却标注为 `Number`。`Int64` 序列化时输出字符串，与官方请求示例 `{"id":"1308675217944611083"}` 一致。

枚举（`DeviceState` / `StationType` / `SynchronizationType` / `MeterType`）同样兼容 `1` 与 `"1"`。

### 未建模字段的兜底

每个结果结构体都嵌入了原始数据，未建模的字段（编号族 `iPv1..32`、`uPv1..32`、`pow1..32`，以及文档示例里出现但表格未列出的字段）都能取到：

```go
res, _ := c.InverterDetail(req)
res.ID.Int()                       // 强类型字段
res.Raw["iPv17"]                   // 原始数据兜底
records, _ := res.Raw["page"].(map[string]any)["records"].([]any)
```

也可以直接调任意接口拿原始数据：

```go
data, err := c.Call(ginlong.PathInverterDetail, ginlong.InverterDetailRequest{SN: "120B40198150131"})
list, err := c.CallList(ginlong.PathInverterDay, req)
```

### 图表接口的两种形态

文档表格与官方示例对 `epm/day` 的 `data` 描述冲突（表格说对象数组，示例是「指标名 → 并列数组」的对象）。该接口返回 `*ginlong.Chart[ginlong.EpmDayItem]`，两种都能读：

```go
ch, _ := c.EpmDay(req)
ch.Items          // data 为对象数组时
ch.Raw["uAc1"]    // data 为「指标名→并列数组」对象时（键名与 searchinfo 一致）
ch.Raw["daylightSwitch"]
```

### 分页与形态兼容

部分接口的**文档表格与官方示例自相矛盾**，SDK 统一做了兼容，调用方只需读 `res.Page.Records`：

| 情况 | 接口 | 处理 |
| --- | --- | --- |
| 表格写 `data.page.records`，示例是 `data.records` 平铺 | `InverterShelfTime`、`AlarmList`、`StationDayEnergyList`、`StationMonthEnergyList`、`StationYearEnergyList` | 两种形态都解析进同一个 `Page` |
| 表格写 `data Object`，示例是数组 | `InverterDetailList`、`StationDetailList`、`InverterDay`、`InverterMonth`、`InverterYear`、`StationDay`、`StationMonth`、`StationYear`、`StationAll`、`CollectorDay` | 按示例用切片 |
| 表格写 `data Array`，示例是对象 | `EpmDay` | 用 `Chart[T]` 兼容两种 |

未建模字段在所有层级都能取到——`Page` 会逐条填充记录的 `Raw`：

```go
for _, inv := range res.Page.Records {
    inv.SN               // 强类型
    inv.Raw["iPv17"]     // 该条记录的原始数据兜底
}
```

## 注意事项

- 所有接口均为 **POST**，`Content-Type: application/json;charset=UTF-8`，返回 `application/json`。
- 所有接口数据更新频率为 **5 分钟**，频率限制 **2 次/秒**。
- `Date` 头与服务器时间相差超过 **±15 分钟**会调用失败（SDK 每次请求实时生成）。
- 列表接口 `pageSize` 最大 **100**，SDK 的 `PageRequest.Normalize()` 会自动补默认值 1/20 并截断上限。
- 请求体里的 `id` / `stationId` 官方示例均以**字符串**下发，`Int64` 已按此序列化。
- 写操作接口（`AddStation` / `StationUpdate` / `AddStationBindCollector` / `DelCollector` / `AddDevice`）会变更账号数据，请谨慎调用。

---

# SolarMan 小麦智电 / 小麦商家版

依据 [SolarMan OpenApi 在线文档](https://doc.solarmanpv.com/)，覆盖文档全部 **39 个接口**（账号 13 + 设备 10 + 电站 13 + 天气 3）。主结构体为 `SolarmanSDK`。

## 快速开始

```go
sdk := solarman.NewSolarmanSDK(solarman.Credentials{
    AppID:     "", // 应用 APPID
    AppSecret: "", // 应用密钥
    Email:     "", // 登录身份三选一：Email / Mobile(+CountryCode) / UserName
    Password:  "", // 明文密码，SDK 内部做 SHA256；也可直接填 PasswordSHA256
    OrgID:     0,  // 商家版填商家 ID，C 端留 0
    // BaseURL: solarman.BaseURLGlobal, // 国际数据中心
}, solarman.WithDebugf(log.Printf))

// 首次调用任意接口时会自动获取 Token，也可以显式获取
if _, err := sdk.AcquireToken(solarman.TokenRequest{}); err != nil {
    log.Fatal(err)
}

stations, err := sdk.StationList(solarman.StationListRequest{
    PageRequest: solarman.PageRequest{Page: 1, Size: 20},
})
if err != nil {
    log.Fatal(err)
}
for _, s := range stations.StationList {
    fmt.Println(s.ID.Int(), s.Name, s.InstalledCapacity.Float())
}
```

### 可选项

| Option | 说明 |
| --- | --- |
| `WithClient(*req.Client)` / `WithTimeout(d)` / `WithDevMode(true)` | 同其它平台 |
| `WithLanguage("en")` | query 参数 `language`，默认 `zh` |
| `WithAccessToken(at, rt, expires)` | 复用已有 Token，跳过首次获取 |
| `WithAppIDInQuery(true)` | 在非 Token 接口的 query 中也带上 `appId` |
| `WithMaxAttempts(n)` | 鉴权失效时的最大尝试次数，默认 2 |
| `WithDisableAutoLogin()` | 关闭自动获取 Token |

## 三种鉴权方式

文档里这三类接口的鉴权方式**不一样**，SDK 已按各自要求处理：

| 方式 | 适用接口 | 说明 |
| --- | --- | --- |
| **Bearer Token** | 绝大多数接口 | 先 `POST /account/v1.0/token?appId=xxx`，body 带 `appSecret` + `password`(SHA256 小写) + 三种登录身份之一，拿 `access_token`；之后每个请求带 `Authorization: bearer {token}` |
| **appId + appSecret（query）** | 2.4 注册帐号、2.8 重置密码、5.1 生成验证码 | 文档中没有 `authorization` 行，`appId`/`appSecret` 都放 query，**不**走 Token。SDK 的 `doAppAuth` 自动补 query |
| **无需鉴权** | 2.1 获取 Token | 只要 `appId` query + body |

细节：

- `Authorization` 的值是 `"bearer " + token`，**小写且带一个空格**（文档专门强调过）。
- `access_token` 有效期约 **60 天**（`expires_in` 约 5183999 秒）；多次获取不会让旧 Token 失效，所以 SDK 在剩余不足 5 分钟时直接重新获取。
- 收到 `2101019`(auth invalid token) / `2101017`(auth token not found) 等鉴权错误时，SDK 会自动重新获取 Token 并重试一次。
- 商家版需要 **两次获取 Token**：先不带 `orgId` 拿 C 端 Token → 调 2.2 拿到 `companyId` → 再用 `sdk.LoginWithOrg(orgId)` 换成商家 Token。
- 重置密码、修改用户角色会让原有 Token 失效。

## 错误处理

```go
if _, err := sdk.StationList(req); err != nil {
    var apiErr *solarman.APIError
    if errors.As(err, &apiErr) {
        log.Printf("code=%s msg=%s hint=%s retryable=%v",
            apiErr.Code, apiErr.Msg,
            solarman.CodeHint(string(apiErr.Code)), apiErr.Retryable())
    }
}
```

| 错误类型 | 说明 |
| --- | --- |
| `*APIError` | 业务错误（`success=false` 或 `code` 非空），`Hint()` 给出中文提示 |
| `*BadResponseError` | 响应不是预期的 JSON |
| `ErrNotLoggedIn` | 关闭自动获取 Token 且本地无 Token |

`APIError.Retryable()` 对 `3201001`/`2101002`/`3501004` 返回 true；`IsAuthError()` 覆盖各类鉴权失效码。常见码：`2101010` 调用次数用完、`2101009` 接口未开通写权限、`2101026` 分页过大、`2101012`/`2101013` 时间范围超限、`2101022` 请求越权（账号与 Home/Pro 端不一致）。

## 接口清单

### 账号接口（`/account/v1.0/*`）

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `AcquireToken` / `Login` | 2.1 获取 Token | `/account/v1.0/token` |
| `AccountInfo` | 2.2 账号商家关系（拿 orgId） | `/account/v1.0/info` |
| `AccountRole` | 2.3 账号内权限 | `/account/v1.0/role` |
| `RegisterUser` | 2.4 注册帐号 | `/account/v1.0/user` |
| `UserInfo` | 2.5 获取账号信息 | `/account/v1.0/user-info` |
| `UpdateUserInfo` | 2.6 修改账号信息 | `/account/v1.0/user-info-update` |
| `UpdateBindInfo` | 2.7 修改绑定信息 | `/account/v1.0/bind-info` |
| `ResetPassword` | 2.8 重置密码 | `/account/v1.0/password-reset` |
| `UpdatePassword` | 2.9 修改密码 | `/account/v1.0/password-update` |
| `CancelCheck` | 2.10 账号注销校验 | `/account/v1.0/cancelCheck` |
| `CancelAccount` | 2.11 账号注销 | `/account/v1.0/cancel` |
| `Captcha` | 5.1 生成验证码 | `/account/v1.0/captcha` |
| `AppIDBalance` | 13 查询 APPID 剩余调用次数 | `/account/v1.0/balance` |

### 设备接口（`/device/v1.0/*`）

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `AlertDetail` | 3.1 设备报警明细 | `/device/v1.0/alertDetail` |
| `AlertList` | 3.2 设备报警列表 | `/device/v1.0/alertList` |
| `CurrentData` | 3.3 设备实时数据（数据点表） | `/device/v1.0/currentData` |
| `DeviceHistorical` | 3.4 设备历史数据 | `/device/v1.0/historical` |
| `DeviceList` | 3.5 设备库-设备列表 | `/device/v1.0/list` |
| `DeviceSimInfo` | 3.6 网关设备 SIM 卡信息 | `/device/v1.0/simInfo` |
| `DeviceCommunication` | 3.8 设备通讯关系 | `/device/v1.0/communication` |
| `CustomControl` | 3.9 自定义指令（指令透传） | `/device/v1.0/customControl` |
| `DeviceRegister` | 4.13 添加网关（采集器或 DTU） | `/device/v1.0/register` |
| `DeviceDelete` | 4.14 删除设备 | `/device/v1.0/delete` |

### 电站接口（`/station/v1.0/*`）

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `StationBase` | 4.1 查询电站基础信息 | `/station/v1.0/base` |
| `StationDeviceList` | 4.2 获取电站设备列表 | `/station/v1.0/device` |
| `StationHistory` | 4.3 获取电站历史数据 | `/station/v1.0/history` |
| `StationList` | 4.4 获取账号下电站列表 | `/station/v1.0/list` |
| `StationRealTime` | 4.5 获取电站实时数据 | `/station/v1.0/realTime` |
| `StationRole` | 4.6 获取电站操作权限 | `/station/v1.0/role` |
| `StationAlert` | 4.7 获取电站报警列表 | `/station/v1.0/alert` |
| `StationCreate` | 4.8 创建电站 | `/station/v1.0/create` |
| `StationUpdate` | 4.9 修改电站 | `/station/v1.0/update` |
| `StationDelete` | 4.10 删除电站 | `/station/v1.0/delete` |
| `StationMetering` | 4.11 累计发电量计算方法设置 | `/station/v1.0/metering` |
| `StationOffset` | 4.12 设置偏移量 | `/station/v1.0/offset` |
| `StationAlertV2` | 4.16 电站报警列表 V2.0 | `/station/v1.0/alertV2` |

### 天气接口

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `StationWeather` | 4.17 查询电站实时天气 | `/station/v1.0/weather/info` |
| `CurrentWeather` | 4.18 查询实时天气 | `/weather/v1.0/getCurrentWeather` |
| `ForecastWeather` | 4.19 查询未来天气 | `/weather/v1.0/getForecastWeather` |

## 响应结构与容错

SolarMan 的响应是**平铺**的——`code` / `msg` / `success` / `requestId` 与业务字段在同一层，没有 `data` 包裹。SDK 用内嵌 `Response` 承载公共字段：

```go
type StationListResult struct {
    Response                    // 公共响应字段
    rawHolder                   // 原始数据兜底
    Total       int             `json:"total"`
    StationList []StationListItem `json:"stationList"`
}
```

成功时 `code` 为 `null`，用 `res.OK()` 判断即可。容错标量与锦浪同款：

| 类型 | 兼容形态 | 用法 |
| --- | --- | --- |
| `Num` | 数值 / 字符串 / null | `.Float()` |
| `Int64` | 数值 / 字符串 / 空串 / null | `.Int()` |
| `Str` | 字符串 / 数值 / null | `.String()` |

**未建模字段的兜底**：结果结构体、列表元素与嵌套结构体都嵌入了原始响应，`fillRaw` 会递归填充：

```go
station.Raw["batterySoc"]              // 字段表没列但示例返回了的字段
devices.DeviceListItems[0].Raw         // 列表元素同样有
sdk.Call(solarman.PathStationList, req) // 任意接口拿原始响应
```

## 注意事项

- 所有接口均为 **POST**，`Content-Type: application/json`，参数 utf-8。
- 常规接口限流 **300 次/10 秒**（单 OpenApi 账号）；设备控制类（3.7）**50 次/分钟**。
- 数据中心二选一：中国区 `https://api.solarmanpv.com`、国际区 `https://globalapi.solarmanpv.com`，用 `Credentials.BaseURL` 指定。
- `password` / `newPassword` / `oldPassword` 都要传 **SHA256 小写密文**，可用 `solarman.SHA256Hex`。
- 时间字段（`startOperatingTime`/`createdDate`/`lastUpdateTime` 等）实际是 **UNIX 秒**，SDK 用 `Num` 承接。
- 分页接口 `page` 必填、`size` 可选，SDK 的 `PageRequest.Normalize()` 会补默认页码。
- `13 查询 APPID 剩余可调用次数` **本身也计费**，不要频繁调用。
- 时间范围限制：多数查询不超过 30 天（部分接口 12 个月）。
- 写操作接口（注册、创建/修改/删除电站、增删设备、改密码、注销）会变更账号数据，请谨慎调用。

---

## 未实现

- 华为 OAuth 2.0 接入方式（文档 3.1 节）：授权码流程、获取/刷新/注销 AT、以 `Authorization: Bearer` 调用 OpenAPI。
- 锦浪 OAuth 2.0 多账号授权（该部分需与锦浪签署合作协议后单独获取文档）。
- SolarMan 的 Token「延长」接口：文档只在概述里提到，未给出接口定义；SDK 通过重新获取 Token 达到同样效果（平台允许多次获取且旧 Token 不失效）。
- SolarMan 设备固件升级接口：`/device/v1.0/upgrade` 只在文档的链接残留里出现过，没有参数说明，`PathDeviceUpgrade` 常量已备好但未实现。

## 开发

```bash
go build ./...
go vet ./...
go test ./...    # 使用 httptest 模拟服务端，不访问真实接口
```

调用示例：华为见 `main.go`，锦浪见 `cmd/ginlong/main.go`，SolarMan 见 `cmd/solarman/main.go`（后两者凭证留空，测试时自行填写）。
