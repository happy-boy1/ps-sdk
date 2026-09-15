// Package service 把 sdk 下各平台的实现收敛成统一的数据接入层。
//
// 职责划分：
//
//	service/adapter.go    统一适配器接口，屏蔽四个平台的鉴权与分页差异
//	service/*_client.go   各平台适配器实现，负责字段映射与设备类型归一
//	service/store.go      电站/设备落库（按 平台+原始ID 幂等 UPSERT）
//	service/sync.go       同步编排，单项失败不影响整体
//	service/resolve.go    凭据解析：环境变量优先，其次 platform_auth 表
//
// 典型用法：
//
//	if err := database.Init(); err != nil { ... }
//	if err := service.Init(database.DB); err != nil { ... }
//
//	client, err := service.PlatformClient("solarman", service.Options{Debug: true})
//	if err != nil { ... }
//	result, err := client.Sync()
//
// 平台凭据的环境变量命名规则见 CredentialFromEnv 的说明。
package service
