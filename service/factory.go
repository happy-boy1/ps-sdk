package service

// PlatformClient 创建指定平台的客户端，凭据取环境变量与 platform_auth 表
func PlatformClient(code string, opts Options) (*Client, error) {
	return Open(code, opts)
}

// SolarmanClient 创建 SolarMan（小麦智电 / 小麦商家版）客户端
func SolarmanClient(opts Options) (*Client, error) {
	return Open(CodeSolarman, opts)
}

// SungrowClient 创建阳光云（iSolarCloud）客户端
func SungrowClient(opts Options) (*Client, error) {
	return Open(CodeSungrow, opts)
}

// GinlongClient 创建锦浪云（SolisCloud）客户端
func GinlongClient(opts Options) (*Client, error) {
	return Open(CodeGinlong, opts)
}

// FusionSolarClient 创建华为 FusionSolar 客户端
func FusionSolarClient(opts Options) (*Client, error) {
	return Open(CodeFusionSolar, opts)
}
