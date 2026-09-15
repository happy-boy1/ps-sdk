package sungrow

import (
	"strconv"
	"strings"
)

type LoginState int

const (
	LoginStateAccountNotExists      LoginState = -1 // 用户账户不存在
	LoginStateWrongPassword         LoginState = 0  // 密码错误
	LoginStateOK                    LoginState = 1  // 登录成功
	LoginStateAccountLockByPassword LoginState = 2  // 由于多次密码输入错误导致账户被锁定
	LoginStateAccountLockByAdmin    LoginState = 5  // 该账号被管理员锁定
)

// String 返回登录状态说明
func (s LoginState) String() string {
	switch s {
	case LoginStateAccountNotExists:
		return "用户账户不存在"
	case LoginStateWrongPassword:
		return "密码错误"
	case LoginStateOK:
		return "登录成功"
	case LoginStateAccountLockByPassword:
		return "多次密码输入错误导致账户被锁定"
	case LoginStateAccountLockByAdmin:
		return "该账号被管理员锁定"
	default:
		return "未知登录状态"
	}
}

// ParseLoginState 解析登录响应里的 login_state 字符串，解析失败返回 false
func ParseLoginState(v string) (LoginState, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, false
	}
	return LoginState(n), true
}

type ShareType int

const (
	ShareTypeGlanceOver ShareType = 1 // 浏览权限
	ShareTypeManage     ShareType = 2 // 管理权限
)
