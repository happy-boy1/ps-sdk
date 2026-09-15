// Package service 把 sdk 下各平台的实现收敛成统一的数据接入层。
//
// 职责划分：
//
//	service/adapter.go    统一适配器接口，屏蔽四个平台的鉴权与分页差异，内置请求节流
//	service/*_client.go   各平台适配器实现，负责字段映射与设备类型归一
//	service/store.go      电站/设备落库（按 平台+原始ID 幂等 UPSERT）
//	service/sync.go       同步编排，单项失败不影响整体，支持并发电站
//	service/resolve.go    凭据解析：config.toml → platform_auth 表 → 环境变量
//
// 典型用法：
//
//	cfg, err := config.Load("")            // pkg/config
//	if err != nil { ... }
//	if err := database.Init(cfg.Database); err != nil { ... }
//	if err := service.InitWith(database.DB, cfg); err != nil { ... }
//
//	client, err := service.PlatformClient("solarman", service.OptionsFrom(cfg))
//	if err != nil { ... }
//	result, err := client.Sync()
//
// 不使用配置文件时用 service.Init(database.DB)，凭据只依赖环境变量与数据库；
// 平台凭据的环境变量命名规则见 pkg/config 与 CredentialFromDB 的说明。
package service
