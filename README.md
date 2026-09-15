# ps-sdk

各光伏云平台 API 实现集合。

| 平台 | 包 | 文档依据 | 鉴权方式 |
| --- | --- | --- | --- |
| 华为 FusionSolar SmartPVMS | `sdk/huawei` | SmartPVMS 26.2.0 北向接口参考 | API 账户登录 + XSRF-TOKEN |
| 锦浪 Ginlong / SolisCloud | `sdk/ginlong` | 锦浪云平台 API 文档 V2.0.3 | KeyID/KeySecret + HMAC-SHA1 签名 |
| SolarMan 小麦智电 / 小麦商家版 | `sdk/solarman` | SolarMan OpenApi 在线文档 V2.0.3 | OAuth2 Bearer Token（60 天） |
| 阳光云 Sungrow iSolarCloud | `sdk/sungrow` | Sungrow OpenApi 文档 | appkey + x-access-key + Token |

`service` 包在四个 SDK 之上提供统一的数据接入层：把各平台的电站与设备归一后写入数据库。
运行参数与平台凭据由 `pkg/config` 从 TOML 文件加载。

## 目录结构

```
pkg/config/               配置加载（TOML + 环境变量覆盖）
  config.go               配置结构、加载、平台凭据（实现 service.ConfigProvider）
  duration.go             "30s" 形式的时长字段
  env.go                  PS_* 环境变量覆盖
  config.example.toml     配置模板（随代码分发，go:embed 进二进制）
  config.toml             本机实际配置，含口令，已加入 .gitignore

service/                  统一接入层（平台适配、字段归一、落库、同步编排）
  doc.go                  包说明
  platform.go             平台常量与元信息（ID / 短码 / 默认地址）
  devicetype.go           统一设备类型目录与各平台类型编码映射
  status.go               电站与设备状态枚举
  errors.go               统一错误、日志出口、部分成功错误、运行参数
  adapter.go              适配器接口、构造器注册表、请求节流
  factory.go              平台客户端构造函数
  wire.go                 配置对接（OptionsFrom / OpenFromConfig）
  resolve.go              凭据解析（配置文件 → 数据库 → 环境变量）
  store.go                基础数据写入与电站/设备 UPSERT
  sync.go                 同步编排、SyncAll、Client
  syncresult.go           同步结果
  convert.go              时间与数值的容错转换
  *_client.go             四个平台的适配器实现
  *_test.go               单元测试

pkg/database/sql.go       数据库连接（按配置连接，返回错误而非 panic）
pkg/tools/jsonfile.go     JSON 文件读写

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

sdk/sungrow/
  client.go               客户端 SungrowSDK、凭证、Token 管理、请求内核（自动登录）
  consts.go               路径、返回码、设备类型/电站类型/在线状态枚举
  types.go                Envelope 响应封装、Num/Int64/Str 容错标量、Entity、ItemMap
  api_token.go            登录、登录状态枚举、共享类型
  api_monitor.go          电站列表、设备列表、设备实时测点数据

main.go                   命令行入口：同步或查看各平台数据
cmd/ginlong/main.go       锦浪 SDK 调用示例
cmd/solarman/main.go      SolarMan SDK 调用示例
```

---

# 配置（pkg/config）

## 三层优先级

```
代码默认值  <  config.toml  <  环境变量（PS_ 前缀）
```

- **平台凭据**：`app_id`（AppKey）与 `app_secret`（AppSecret）**以数据库 `platform_auth` 表为准**；
  接口地址、登录账号、密码、商家 ID 等数据库未建模的信息由 `config.toml` 提供；
  环境变量可覆盖任意字段，便于临时切换。
- 数据库不可用或缺行时，凭据退化为「config.toml + 环境变量」，不会中断同步。

## config.toml

模板见 `pkg/config/config.example.toml`，复制为 `pkg/config/config.toml` 后填写。
`config.toml` 含账号口令，已加入 `.gitignore`。

```toml
[http]
timeout = "30s"        # 单次请求超时
max_attempts = 3       # 含首调在内的最大尝试次数
retry_delay = "2s"     # 重试基础退避
debug = false          # 打印请求日志
dev_mode = false       # 打印完整请求/响应报文

[sync]
concurrency = 1        # 同步设备时的并发电站数（被限流时会自动串行退避）

[database]
host = "127.0.0.1"
port = 3306
user = "root"
password = ""
name = "rtm"
# dsn 非空时直接使用，忽略其余分项

[platform.solarman]
api_url = "https://api.solarmanpv.com"
account = ""           # 登录身份三选一：account / email / user_name
country_code = "86"
password = ""
org_id = 0             # 商家版填商家 ID

[platform.sungrow]
api_url = "https://gateway.isolarcloud.com"
account = ""           # 登录账号（appkey 不能当账号用）
password = ""

[platform.ginlong]
api_url = "https://api.ginlong.com:13333"

[platform.huawei]
api_url = "https://intl.fusionsolar.huawei.com"
account = ""           # API 账户名（app_id 亦可）
password = ""          # API 账户密码（app_secret 亦可）
```

加载顺序：`-config` 参数 → `PS_CONFIG` 环境变量 → `./config.toml` → `./pkg/config/config.toml`；
都不存在时按模板在工作目录写出 `config.toml` 并继续使用默认值。

## 环境变量

| 用途 | 变量 |
| --- | --- |
| 配置文件路径 | `PS_CONFIG` |
| 数据库 | `PS_DB_DSN` / `PS_DB_HOST` / `PS_DB_PORT` / `PS_DB_USER` / `PS_DB_PASSWORD` / `PS_DB_NAME` / `PS_DB_CHARSET` / `PS_DB_MAX_IDLE` / `PS_DB_MAX_OPEN` / `PS_DB_CONN_MAX_LIFE` |
| 请求与同步 | `PS_HTTP_TIMEOUT` / `PS_MAX_ATTEMPTS` / `PS_RETRY_DELAY` / `PS_DEBUG` / `PS_DEV_MODE` / `PS_SYNC_CONCURRENCY` |
| solarman | `PS_SOLARMAN_APPID` / `PS_SOLARMAN_APPSECRET` / `PS_SOLARMAN_ACCOUNT` / `PS_SOLARMAN_PASSWORD` / `PS_SOLARMAN_ORGID` / `PS_SOLARMAN_COUNTRYCODE` / `PS_SOLARMAN_EMAIL` / `PS_SOLARMAN_USERNAME` |
| sungrow | `PS_SUNGROW_APPID` / `PS_SUNGROW_APPSECRET` / `PS_SUNGROW_ACCOUNT` / `PS_SUNGROW_PASSWORD` |
| ginlong | `PS_GINLONG_APPID` / `PS_GINLONG_APPSECRET` |
| huawei | `PS_HUAWEI_APPID`（API 账户名）/ `PS_HUAWEI_PASSWORD` |
| 通用 | `PS_<平台>_APIURL` 覆盖接口地址，`PS_<平台>_APP_ID` 等带下划线写法同样识别 |

## 代码中的用法

```go
cfg, err := config.Load("")            // 加载配置（含环境变量覆盖）
if err != nil {
    log.Fatal(err)
}
if err := database.Init(cfg.Database); err != nil {
    log.Fatal(err)
}
defer database.Close()

// 注入配置：平台凭据与运行参数都从 cfg 取
if err := service.InitWith(database.DB, cfg); err != nil {
    log.Fatal(err)
}

client, err := service.PlatformClient("solarman", service.OptionsFrom(cfg))
if err != nil {
    log.Fatal(err)
}
result, err := client.Sync()
log.Println(result)
```

不使用配置文件时用 `service.Init(database.DB)`，凭据只依赖环境变量与数据库。

---

# service 统一接入层

## 能力

- **统一适配**：四个平台的鉴权方式、分页字段、设备类型编码差异全部收敛在适配器内部，
  上层只依赖 `service.Adapter`。
- **字段归一**：平台原始电站/设备 ID 落在 `station_id_origin` / `device_id_origin`，
  与 `platform_id` 组成幂等键；统一设备类型见 `service/devicetype.go`。
- **幂等落库**：按「平台 + 原始 ID」UPSERT，重复同步不会产生脏数据。
- **局部容错**：单个电站或单个设备族失败不影响其余数据，结果中按「失败 / 告警」分别汇报。

## 快速开始

```go
package main

import (
    "log"

    "ps-sdk/pkg/config"
    "ps-sdk/pkg/database"
    "ps-sdk/service"
)

func main() {
    cfg, err := config.Load("") // 配置见「配置（pkg/config）」章节
    if err != nil {
        log.Fatal(err)
    }
    if err := database.Init(cfg.Database); err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    if err := service.InitWith(database.DB, cfg); err != nil { // 自动补表 + 写入平台/设备类型基础数据
        log.Fatal(err)
    }

    client, err := service.PlatformClient("solarman", service.OptionsFrom(cfg))
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.Sync()
    log.Println(result) // 小麦智电: 电站 96/96，设备 1296/1296，失败 0 项，告警 0 项
    if err != nil {
        log.Fatal(err)
    }
}
```

命令行方式：

```bash
go run . -platform solarman     # 同步指定平台
go run . -platform all          # 同步全部已配置凭据的平台
go run . -list                  # 只读库，打印已同步的电站与设备
go run . -seed                  # 只写入平台/设备类型基础数据
go run . -config ./my.toml      # 指定配置文件
go run . -platform ginlong -debug
```

## 接口

| 方法 | 说明 |
| --- | --- |
| `Init(db)` | 只注入数据库句柄，自动补表并写入基础数据 |
| `InitWith(db, cfg)` | 额外注入配置来源（`*config.Config`），平台凭据从中读取 |
| `Seed()` | 仅写入基础数据，可重复调用 |
| `PlatformClient(code, opts)` | 按平台短码创建客户端 |
| `OptionsFrom(cfg)` | 从配置来源取运行参数 |
| `NewAdapter(cred, opts)` | 用显式凭据创建适配器 |
| `NewClient(adapter, opts)` | 用自定义适配器创建客户端，便于测试 |
| `(*Client).Sync()` | 拉取电站与设备并落库，返回 `*SyncResult` |
| `(*Client).SyncStations()` | 只拉取电站并落库 |
| `(*Client).SyncDevices(stationID)` | 只拉取指定电站的设备并落库 |
| `(*Client).Stations()` / `Devices(id)` | 只读库，stationID 为 0 表示全部电站 |
| `ListStations(platformID)` / `ListDevices(platformID, stationID)` | 读取已落库数据，platformID 为 0 表示全部平台 |
| `UpsertStation(s)` / `UpsertDevice(d)` | 手动写入单条记录 |
| `SyncAll(opts)` | 遍历全部平台同步，单平台失败不影响其余 |

平台短码：`solarman`、`sungrow`、`ginlong`、`huawei`（即 `service.CodeXxx`）。

`SyncResult` 区分两类问题：`Errors` 是电站或设备落库失败等硬失败，`Warnings` 是
「设备族接口无权限」这类不影响其余数据的告警，两者都会在命令行输出。

## 设备类型归一

统一设备类型的 ID 段位：`11xx` 光伏类、`12xx` 储能类、`13xx` 计量与采集类、`14xx` 环境与其他。
每个平台自己的类型编码保留在 `PowerDevice.DeviceTypeOrigin`，映射表见 `service/devicetype.go`，
未登记的平台编码归入 `UNKNOWN`，可通过 `service.KnownDeviceTypes("solarman")` 查看已登记编码。

平台状态值统一换算为：

| 域 | 取值 |
| --- | --- |
| 电站状态（`model.PowerStation.Status`） | 0 停运、1 运行、2 在建 |
| 设备状态（`model.PowerDevice.Status`） | 0 未知、1 在线、2 离线、3 告警 |

## 已知平台限制

- 锦浪云没有统一设备列表接口，按设备族（逆变器 / 采集器 / EPM / 电表 / 气象仪）逐个拉取。
  某一族接口未开通权限时（`R0000`）返回 `*PartialError`，其余设备照常入库，
  该情况计入 `SyncResult.Warnings` 而非失败。
- 华为设备列表需先取电站编号并按 100 个电站批量查询，电站列表不返回在线状态。
- 阳光云设备列表按电站 `ps_id` 过滤，设备类型编码与文档标注的类型不一致（两种形态都兼容）。
- 平台限流阈值多数未公开，`service` 对每个平台内置请求节流：请求按最小间隔串行发出，
  被限流时自适应加倍退避（上限 5s），设备列表还会自动重试；因此把 `[sync] concurrency`
  调大不会压垮平台，但增大到 4 以上通常收益有限。

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

凭证允许留空构造，只有**真正发起请求时**才会返回 `huawei.ErrMissingCredentials`；
也可以用 `sdk.SetCredentials(...)` 后续补充（会同时清空已缓存的 token）。

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
c, err := ginlong.NewSolisSDK(ginlong.Credentials{
    APIID:     "",   // KeyID，锦浪云 WEB 端「服务 - API 管理」获取
    APISecret: "",   // KeySecret
    // BaseURL: ginlong.DefaultBaseURL, // 缺省 https://api.ginlong.com:13333
}, ginlong.WithDebugf(log.Printf))
if err != nil {
    log.Fatal(err)
}

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
sdk, err := solarman.NewSolarmanSDK(solarman.Credentials{
    AppID:     "", // 应用 APPID
    AppSecret: "", // 应用密钥
    Email:     "", // 登录身份三选一：Email / Mobile(+CountryCode) / UserName
    Password:  "", // 明文密码，SDK 内部做 SHA256；也可直接填 PasswordSHA256
    OrgID:     0,  // 商家版填商家 ID，C 端留 0
    // BaseURL: solarman.BaseURLGlobal, // 国际数据中心
}, solarman.WithDebugf(log.Printf))
if err != nil {
    log.Fatal(err)
}

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
    fmt.Println(s.ID.String(), s.Name, s.InstalledCapacity.Float())
}
```

凭证允许留空构造，只有**真正发起请求时**才会返回 `solarman.ErrMissingCredentials`；也可以用 `sdk.SetCredentials(...)` 后续补充。

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
| `ErrMissingCredentials` | 凭证未填写 |

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
    Response                      // 公共响应字段
    rawHolder                     // 原始数据兜底
    Total       int               `json:"total"`
    StationList []StationListItem `json:"stationList"`
}
```

成功时 `code` 为 `null`，用 `res.OK()` 判断即可。容错标量与锦浪同款：

| 类型 | 兼容形态 | 用法 |
| --- | --- | --- |
| `Num` | 数值 / 字符串 / null | `.Float()` |
| `Int64` | 数值 / 字符串 / 空串 / null | `.Int()` |
| `Str` | 字符串 / 数值 / null | `.String()` |

`StationListItem.ID` 用 `Str` 承接（文档标注 `Number`，实际形态不固定），
需要按数值下发时自行转换，见 `cmd/solarman/main.go` 的 `stationID` 辅助函数。

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

# 阳光云 Sungrow iSolarCloud

## 快速开始

```go
sdk, err := sungrow.NewSungrowSDK(sungrow.Credentials{
    AppID:        "", // 应用 appkey
    AppSecret:    "", // x-access-key
    UserAccount:  "", // 登录账号
    UserPassword: "", // 登录密码
    // BaseURL: sungrow.BaseURLGlobal, // 国际站缺省为中国站
}, sungrow.WithDebugf(log.Printf))
if err != nil {
    log.Fatal(err)
}

// 登录是隐式的：首次调用业务接口时自动登录并缓存 token
stations, err := sdk.GetPowerStationList(sungrow.PowerStationListRequest{
    PageRequest: sungrow.PageRequest{CurPage: 1, Size: 100},
})
if err != nil {
    log.Fatal(err)
}
for _, ps := range stations.PageList {
    fmt.Println(ps.PsId.String(), ps.PsName, ps.TotalCapcity.String())
}
```

凭证允许留空构造，只有**真正发起请求时**才会返回 `sungrow.ErrMissingCredentials`；
也可以用 `sdk.SetCredentials(...)` 后续补充（会同时清空已缓存的 token）。

### 可选项

| Option | 说明 |
| --- | --- |
| `WithClient(*req.Client)` / `WithTimeout(d)` / `WithDevMode(true)` | 同其它平台 |
| `WithAccessToken(token, expiresIn)` | 复用已有 Token；有效期未知时按 24 小时估算 |
| `WithMaxAttempts(n)` | 含首调在内的最大尝试次数，默认 3 |
| `WithRetryBaseDelay(d)` | 重试退避基数，默认 2s |
| `WithDisableAutoLogin()` | 关闭自动登录，无 token 时直接返回错误 |

## 鉴权与请求内核

- 每个请求都带 `x-access-key: {appSecret}`、`sys_code: 901`，请求体内注入 `appkey` 与 `token`。
  业务请求体都是「内嵌 `Request` 的结构体指针」，由 `callOnce` 统一注入，接口方法内部不关心 token。
- `POST /openapi/login` 成功后拿 `result_data.token`，SDK 按 `TokenTTL`（24 小时）缓存；
  剩余不足 5 秒时视为过期并重新登录。登录错误码 `E00003` 会触发丢弃缓存并重试。
- `E901`（调用频繁）固定等待 1 分钟重试；其余可重试错误按 `WithRetryBaseDelay` 指数退避。
- 登录失败（账号不存在、密码错误、账户锁定）**不会重试**，避免连续输错导致锁定。

## 错误处理

```go
_, err := sdk.GetPowerStationList(req)
if err != nil {
    var apiErr *sungrow.APIError
    if errors.As(err, &apiErr) {
        log.Printf("resultCode=%s hint=%s retryable=%v",
            apiErr.ResultCode, apiErr.ResultCode.Hint(), apiErr.Retryable())
    }
}
```

| 错误类型 | 说明 |
| --- | --- |
| `*APIError` | 业务错误（`result_code != "1"`），`Hint()` 给出中文提示 |
| `*BadResponseError` | 响应不是预期的 JSON |
| `ErrNotLoggedIn` | 关闭自动登录且本地无 token |
| `ErrMissingCredentials` | 凭证未填写 |

关键码：`1` 成功、`E00003` token 失效、`E901` 调用频繁、`E999`/`E998` 小时/月调用次数上限、
`E918`/`E919` 白名单限制、`E916` 登录频繁、`E913` 时间戳偏差过大。

## 接口清单

| 方法 | 接口 | 路径 |
| --- | --- | --- |
| `Login` / `Logout` | 登录 / 注销 | `/openapi/login` |
| `GetPowerStationList` | 电站列表（分页） | `/openapi/getPowerStationList` |
| `GetDeviceListByPsID` | 电站下设备列表（按 `ps_id` 过滤，分页） | `/openapi/getDeviceList` |
| `GetDeviceRealTimeData` | 设备实时测点数据（点表 key 为 `p*`） | `/openapi/getDeviceRealTimeData` |

`consts.go` 中还登记了电站详情、故障告警、测点历史、参数设置、数据订阅等路径常量，
其中标注「暂时不实现」的部分未提供方法，可按同样模式扩展。

## 类型与容错

| 类型 | 兼容形态 | 用法 |
| --- | --- | --- |
| `Num` | 数值 / 字符串 / null | `.Float()` |
| `Int64` | 数值 / 字符串 / 空串 / null | `.Int()` |
| `Str` | 字符串 / 数值 / 布尔 / null | `.String()` |
| `Entity` | `{"unit":"kW","value":"1.5"}` | `.String()` → `1.5 (kW)` |

枚举（`DevTypeID` / `PsType` / `PsOnlineStatus`）兼容 `1` 与 `"1"` 两种形态，
`DevTypeID.String()` / `PsType.String()` 返回中文名称。

`DevicePointInner` 会把响应里 `p*` 形式的测点收进 `Points`，其余字段保持强类型：

```go
rtd, _ := sdk.GetDeviceRealTimeData(req)
for _, dp := range rtd.DevicePointList {
    fmt.Println(dp.DevicePoint.PsKey, dp.DevicePoint.Points["p1"])
}
```

## 注意事项

- 所有接口均为 **POST**，`Content-Type: application/json`。
- 国际站与国内站地址二选一：`BaseURLChina` / `BaseURLGlobal`。
- 设备类型编码（`device_type`）文档标注与真实返回不一致，SDK 用 `Str` 承接，`service` 层再做归一。
- 电站列表用 `PsId`（字符串）标识电站，设备列表返回的 `ps_id` 是数值，两者需按字符串比对。

---

## 未实现

- 华为 OAuth 2.0 接入方式（文档 3.1 节）：授权码流程、获取/刷新/注销 AT、以 `Authorization: Bearer` 调用 OpenAPI。
- 锦浪 OAuth 2.0 多账号授权（该部分需与锦浪签署合作协议后单独获取文档）。
- SolarMan 的 Token「延长」接口：文档只在概述里提到，未给出接口定义；SDK 通过重新获取 Token 达到同样效果（平台允许多次获取且旧 Token 不失效）。
- SolarMan 设备固件升级接口：`/device/v1.0/upgrade` 只在文档的链接残留里出现过，没有参数说明，`PathDeviceUpgrade` 常量已备好但未实现。
- 阳光云的参数设置、只读参数、故障告警、测点历史与数据订阅接口：路径常量已登记，方法未实现。

## 开发

```bash
go build ./...
go vet ./...
gofmt -l .       # 应为空
go test ./...    # sdk/huawei 与 service 使用 httptest/单元测试，不访问真实接口
```

调用示例：华为见各 SDK 文档示例，锦浪见 `cmd/ginlong/main.go`，SolarMan 见 `cmd/solarman/main.go`，
四个平台的统一接入见 `service` 与 `main.go`（`go run . -h` 查看参数）。
