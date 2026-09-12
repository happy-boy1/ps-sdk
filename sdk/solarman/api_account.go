package solarman

// 账号相关接口：2.2 账号商家关系 / 2.3 账号内权限 / 2.4 注册帐号 / 2.5 获取账号信息 /
// 2.6 修改账号信息 / 2.7 修改绑定信息 / 2.8 重置密码 / 2.9 修改密码 /
// 2.10 账号注销校验 / 2.11 账号注销 / 5.1 生成验证码 / 13 查询 APPID 剩余可调用次数。
//
// 响应结构为平铺式：code / msg / success / requestId 与业务字段同层，
// 因此每个结果结构体匿名嵌入 Response（公共字段）与 rawHolder（原始数据兜底）。

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// 2.2 账号商家关系
// ---------------------------------------------------------------------------

// AccountInfoOrg 账号归属的商家及角色
type AccountInfoOrg struct {
	rawHolder
	CompanyID   Int64  `json:"companyId"`   // 公司 ID
	CompanyName string `json:"companyName"` // 公司名称
	RoleName    string `json:"roleName"`    // 角色名称
}

// AccountInfoResult 2.2 账号商家关系结果
type AccountInfoResult struct {
	Response                     // 公共响应字段
	rawHolder                    // 原始数据兜底
	OrgInfoList []AccountInfoOrg `json:"orgInfoList"` // 账号商家列表
}

// AccountInfo 2.2 查询账号归属的商家及角色
func (sdk *SolarmanSDK) AccountInfo() (*AccountInfoResult, error) {
	var out AccountInfoResult
	if err := sdk.do(PathInfo, nil, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.3 账号内权限
// ---------------------------------------------------------------------------

// AccountRoleResult 2.3 账号内权限结果，0 无权限 1 有权限
type AccountRoleResult struct {
	Response                // 公共响应字段
	rawHolder               // 原始数据兜底
	ViewDeviceAlertData int `json:"viewDeviceAlertData"` // 查看报警详情
	ViewDeviceAlertList int `json:"viewDeviceAlertList"` // 查看报警信息
	ViewDeviceData      int `json:"viewDeviceData"`      // 查看设备详情
	ViewDeviceList      int `json:"viewDeviceList"`      // 查看设备列表
	ViewPlantAlertList  int `json:"viewPlantAlertList"`  // 查看报警列表
	ViewPlantData       int `json:"viewPlantData"`       // 查看电站详情
	ViewPlantDeviceList int `json:"viewPlantDeviceList"` // 查看子系统/设备
	ViewPlantInfo       int `json:"viewPlantInfo"`       // 查看关于电站
	ViewPlantList       int `json:"viewPlantList"`       // 查看电站列表
}

// AccountRole 2.3 查询账号在商家下的内部权限
func (sdk *SolarmanSDK) AccountRole() (*AccountRoleResult, error) {
	var out AccountRoleResult
	if err := sdk.do(PathRole, nil, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.4 注册帐号
// ---------------------------------------------------------------------------

// RegisterUserRequest 2.4 注册帐号请求。
// appId/appSecret 由客户端从 Credentials 取并放到 query 上，故不在此结构体中
type RegisterUserRequest struct {
	Captcha           string `json:"captcha"`                     // 验证码，必填
	Email             string `json:"email,omitempty"`             // 邮箱
	Nickname          string `json:"nickname,omitempty"`          // 昵称，缺省为 User+时间戳
	OldUserID         string `json:"oldUserId,omitempty"`         // 客户自有平台 Uid
	OriginalPhotoURL  string `json:"originalPhotoUrl,omitempty"`  // 头像原图链接
	Password          string `json:"password,omitempty"`          // 密码，需 SHA256 小写密文
	PhoneNumber       string `json:"phoneNumber,omitempty"`       // 手机号
	PhoneNumberPrefix string `json:"phoneNumberPrefix,omitempty"` // 手机号前缀，如 86
	PhotoURL          string `json:"photoUrl,omitempty"`          // 头像压缩图链接
}

// RegisterUserResult 2.4 注册帐号结果
type RegisterUserResult struct {
	Response        // 公共响应字段
	rawHolder       // 原始数据兜底
	UserID    Int64 `json:"userId"` // 新注册的用户 ID
}

// RegisterUser 2.4 注册新用户。该接口用 appId+appSecret（query）鉴权，不需要 Token
func (sdk *SolarmanSDK) RegisterUser(req RegisterUserRequest) (*RegisterUserResult, error) {
	if strings.TrimSpace(req.Captcha) == "" {
		return nil, fmt.Errorf("[Solarman] 注册帐号: captcha 不能为空")
	}

	var out RegisterUserResult
	if err := sdk.doAppAuth(PathUser, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.5 获取账号信息
// ---------------------------------------------------------------------------

// UserInfoResult 2.5 获取账号信息结果
type UserInfoResult struct {
	Response                 // 公共响应字段
	rawHolder                // 原始数据兜底
	Email             string `json:"email"`             // 邮箱号
	LastLoginTime     Str    `json:"lastLoginTime"`     // 最近登录时间，可能是数值
	Nickname          string `json:"nickname"`          // 昵称
	OldUserID         string `json:"oldUserId"`         // 客户自有平台 Uid
	OriginalPhotoURL  string `json:"originalPhotoUrl"`  // 原始图片地址
	PhoneNumber       string `json:"phoneNumber"`       // 手机号
	PhoneNumberPrefix string `json:"phoneNumberPrefix"` // 手机号前缀
	PhotoURL          string `json:"photoUrl"`          // 头像地址
	RegTime           Str    `json:"regTime"`           // 注册时间，可能为空
	UserID            Int64  `json:"userId"`            // 用户 ID
	Username          string `json:"username"`          // 用户名称
}

// UserInfo 2.5 获取我的账号信息
func (sdk *SolarmanSDK) UserInfo() (*UserInfoResult, error) {
	var out UserInfoResult
	if err := sdk.do(PathUserInfo, nil, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.6 修改账号信息
// ---------------------------------------------------------------------------

// UpdateUserInfoRequest 2.6 修改账号信息请求，字段均可选
type UpdateUserInfoRequest struct {
	NewNickname      string `json:"newNickname,omitempty"`      // 新昵称
	OriginalPhotoURL string `json:"originalPhotoUrl,omitempty"` // 头像原图链接
	PhotoURL         string `json:"photoUrl,omitempty"`         // 头像压缩图链接
}

// UpdateUserInfoResult 2.6 修改账号信息结果，仅公共字段
type UpdateUserInfoResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// UpdateUserInfo 2.6 修改我的账号信息
func (sdk *SolarmanSDK) UpdateUserInfo(req UpdateUserInfoRequest) (*UpdateUserInfoResult, error) {
	var out UpdateUserInfoResult
	if err := sdk.do(PathUserInfoUpd, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.7 修改绑定信息
// ---------------------------------------------------------------------------

// UpdateBindInfoRequest 2.7 修改绑定信息请求，字段均可选
type UpdateBindInfoRequest struct {
	Captcha           string `json:"captcha,omitempty"`           // 验证码
	NewEmail          string `json:"newEmail,omitempty"`          // 新邮箱
	NewPhoneNumber    string `json:"newPhoneNumber,omitempty"`    // 新手机号
	NewUsername       string `json:"newUsername,omitempty"`       // 新用户名
	PhoneNumberPrefix string `json:"phoneNumberPrefix,omitempty"` // 手机号前缀，如 86
}

// UpdateBindInfoResult 2.7 修改绑定信息结果，仅公共字段
type UpdateBindInfoResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// UpdateBindInfo 2.7 修改用户的手机号/邮箱/用户名
func (sdk *SolarmanSDK) UpdateBindInfo(req UpdateBindInfoRequest) (*UpdateBindInfoResult, error) {
	var out UpdateBindInfoResult
	if err := sdk.do(PathBindInfo, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.8 重置密码
// ---------------------------------------------------------------------------

// ResetPasswordRequest 2.8 重置密码请求，仅能重置本 APPID 注册的账号
type ResetPasswordRequest struct {
	Captcha           string `json:"captcha,omitempty"`           // 验证码
	Email             string `json:"email,omitempty"`             // 邮箱
	NewPassword       string `json:"newPassword,omitempty"`       // 新密码，需 SHA256 小写密文
	PhoneNumber       string `json:"phoneNumber,omitempty"`       // 手机号
	PhoneNumberPrefix string `json:"phoneNumberPrefix,omitempty"` // 手机号前缀，如 86
}

// ResetPasswordResult 2.8 重置密码结果，仅公共字段
type ResetPasswordResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// ResetPassword 2.8 重置用户密码。该接口用 appId+appSecret（query）鉴权，不需要 Token
func (sdk *SolarmanSDK) ResetPassword(req ResetPasswordRequest) (*ResetPasswordResult, error) {
	var out ResetPasswordResult
	if err := sdk.doAppAuth(PathPasswordRst, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.9 修改密码
// ---------------------------------------------------------------------------

// UpdatePasswordRequest 2.9 修改密码请求
type UpdatePasswordRequest struct {
	NewPassword string `json:"newPassword,omitempty"` // 新密码，需 SHA256 小写密文
	OldPassword string `json:"oldPassword,omitempty"` // 旧密码，需 SHA256 小写密文
}

// UpdatePasswordResult 2.9 修改密码结果，仅公共字段
type UpdatePasswordResult struct {
	Response  // 公共响应字段
	rawHolder // 原始数据兜底
}

// UpdatePassword 2.9 使用旧密码修改我的密码
func (sdk *SolarmanSDK) UpdatePassword(req UpdatePasswordRequest) (*UpdatePasswordResult, error) {
	var out UpdatePasswordResult
	if err := sdk.do(PathPasswordUpd, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.10 账号注销校验
// ---------------------------------------------------------------------------

// CancelCheckOrg 注销校验中的商家信息
type CancelCheckOrg struct {
	rawHolder
	AdminID       Int64  `json:"adminId"`       // 管理员用户 ID
	AreaID        Int64  `json:"areaId"`        // 地区 ID
	BusinessType  string `json:"businessType"`  // 业务板块视图，逗号分隔
	Category      int    `json:"category"`      // 分类，1 安装运维商等
	ID            Int64  `json:"id"`            // 商家 ID
	Logo          string `json:"logo"`          // LOGO 地址
	Name          string `json:"name"`          // 名称
	OperateObject string `json:"operateObject"` // 运维对象，1 光伏 2 风电
	OriginalLogo  string `json:"originalLogo"`  // LOGO 原始图片地址
	System        string `json:"system"`        // 所属系统
	Timezone      string `json:"timezone"`      // 时区
	TopGroupID    Int64  `json:"topGroupId"`    // 顶级组 ID
	Type          int    `json:"type"`          // 类型，1 企业 2 个体
}

// CancelCheckResult 2.10 账号注销校验结果
type CancelCheckResult struct {
	Response                                   // 公共响应字段
	rawHolder                                  // 原始数据兜底
	ExisitDeviceOrgList       []CancelCheckOrg `json:"exisitDeviceOrgList"`       // 存在设备的商家
	ExisitGroupAndUserOrgList []CancelCheckOrg `json:"exisitGroupAndUserOrgList"` // 存在组织成员的商家
	ExisitPlantOrgList        []CancelCheckOrg `json:"exisitPlantOrgList"`        // 存在电站的商家
	ExisitUserPlantRel        bool             `json:"exisitUserPlantRel"`        // 账号下是否存在电站
	Pass                      bool             `json:"pass"`                      // 注销检查是否通过
}

// CancelCheck 2.10 账号注销校验，通过后方可注销
func (sdk *SolarmanSDK) CancelCheck() (*CancelCheckResult, error) {
	var out CancelCheckResult
	if err := sdk.do(PathCancelCheck, nil, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 2.11 账号注销
// ---------------------------------------------------------------------------

// CancelAccountRequest 2.11 账号注销请求
type CancelAccountRequest struct {
	Captcha      string `json:"captcha"`                // 验证码，必填
	Identifier   string `json:"identifier,omitempty"`   // 标识，手机号或邮箱
	IdentityType int    `json:"identityType,omitempty"` // 标识类型，1 手机号 2 邮箱
	UserID       Int64  `json:"userId,omitempty"`       // 用户 ID
}

// CancelAccountResult 2.11 账号注销结果，商家列表字段与 2.10 一致
type CancelAccountResult struct {
	Response                                   // 公共响应字段
	rawHolder                                  // 原始数据兜底
	CancelSuccess             bool             `json:"cancelSuccess"`             // 是否注销成功
	ExisitDeviceOrgList       []CancelCheckOrg `json:"exisitDeviceOrgList"`       // 存在设备的商家
	ExisitGroupAndUserOrgList []CancelCheckOrg `json:"exisitGroupAndUserOrgList"` // 存在组织成员的商家
	ExisitPlantOrgList        []CancelCheckOrg `json:"exisitPlantOrgList"`        // 存在电站的商家
	ExisitUserPlantRel        bool             `json:"exisitUserPlantRel"`        // 账号下是否存在电站
	Pass                      bool             `json:"pass"`                      // 注销检查是否通过
}

// CancelAccount 2.11 账号注销，操作不可恢复
func (sdk *SolarmanSDK) CancelAccount(req CancelAccountRequest) (*CancelAccountResult, error) {
	if strings.TrimSpace(req.Captcha) == "" {
		return nil, fmt.Errorf("[Solarman] 账号注销: captcha 不能为空")
	}

	var out CancelAccountResult
	if err := sdk.do(PathCancel, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 5.1 生成验证码
// ---------------------------------------------------------------------------

// CaptchaRequest 5.1 生成验证码请求，仅生成不发送
type CaptchaRequest struct {
	Email             string `json:"email,omitempty"`             // 邮箱
	PhoneNumber       string `json:"phoneNumber,omitempty"`       // 手机号
	PhoneNumberPrefix string `json:"phoneNumberPrefix,omitempty"` // 手机号前缀，如 86
	Purpose           string `json:"purpose"`                     // 用途，如 REG/BIND/CANCEL
}

// CaptchaResult 5.1 生成验证码结果
type CaptchaResult struct {
	Response         // 公共响应字段
	rawHolder        // 原始数据兜底
	Captcha   string `json:"captcha"`  // 验证码
	Validity  int    `json:"validity"` // 有效期，单位秒
}

// Captcha 5.1 根据手机号或邮箱生成验证码。该接口用 appId+appSecret（query）鉴权，不需要 Token
func (sdk *SolarmanSDK) Captcha(req CaptchaRequest) (*CaptchaResult, error) {
	if strings.TrimSpace(req.Purpose) == "" {
		return nil, fmt.Errorf("[Solarman] 生成验证码: purpose 不能为空")
	}

	var out CaptchaResult
	if err := sdk.doAppAuth(PathCaptcha, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 13 查询 APPID 剩余可调用次数
// ---------------------------------------------------------------------------

// AppIDBalanceRequest 13 查询余量请求
type AppIDBalanceRequest struct {
	AppID string `json:"appId"` // appId，必填
}

// AppIDBalanceResult 13 查询余量结果
type AppIDBalanceResult struct {
	Response        // 公共响应字段
	rawHolder       // 原始数据兜底
	Total     Int64 `json:"total"` // 剩余调用次数
}

// AppIDBalance 13 查询 APPID 剩余可调用次数，本调用也计数
func (sdk *SolarmanSDK) AppIDBalance(req AppIDBalanceRequest) (*AppIDBalanceResult, error) {
	if strings.TrimSpace(req.AppID) == "" {
		return nil, fmt.Errorf("[Solarman] 查询余量: appId 不能为空")
	}

	var out AppIDBalanceResult
	if err := sdk.do(PathBalance, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}
