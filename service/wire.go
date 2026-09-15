package service

// OpenFromConfig 用配置来源创建客户端。
// 配置未提供时，运行参数取代码默认值，凭据仍可来自环境变量与数据库。
func OpenFromConfig(code string, cfg ConfigProvider) (*Client, error) {
	return Open(code, OptionsFrom(cfg))
}

// SyncAllFromConfig 用同一份配置同步全部平台
func SyncAllFromConfig(cfg ConfigProvider) ([]*SyncResult, error) {
	return SyncAll(OptionsFrom(cfg))
}

// OptionsFrom 取配置来源中的运行参数，未配置时返回默认参数
func OptionsFrom(cfg ConfigProvider) Options {
	if cfg == nil {
		return Options{}.WithDefaults()
	}
	return cfg.Options().WithDefaults()
}
