package sungrow

type LoginState int

const (
	LoginStateAccountNotExists      LoginState = -1 // 用户账户不存在
	LoginStateWrongPassword         LoginState = 0  // 密码错误
	LoginStateOK                    LoginState = 1  // 登录成功
	LoginStateAccountLockByPassword LoginState = 2  // 由于多次密码输入错误导致账户被锁定
	LoginStateAccountLockByAdmin    LoginState = 5  // 该账号被管理员锁定
)

type ShareType int

const (
	ShareTypeGlanceOver ShareType = 1 // 浏览权限
	ShareTypeManage     ShareType = 2 // 管理权限
)
