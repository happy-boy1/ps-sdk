package sungrow

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// env 组装统一响应信封，data 为空表示 result_data 为 null
func env(code ResultCode, msg, data string) string {
	if data == "" {
		data = "null"
	}
	return `{"req_serial_num":"1","result_code":"` + string(code) +
		`","result_msg":"` + msg + `","result_data":` + data + `}`
}

// loginData 组装登录成功的 result_data
func loginData(token string) string {
	return fmt.Sprintf(`{"token":"%s","login_state":"1","user_id":"9527"}`, token)
}

// authBody 业务请求体中与鉴权相关的字段
type authBody struct {
	AppKey string `json:"appkey"`
	Token  string `json:"token"`
}

// decodeAuthBody 读取并解析请求体，解析失败仅报错不中断 handler
func decodeAuthBody(t *testing.T, r *http.Request) authBody {
	t.Helper()

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("读取请求体失败: %v", err)
		return authBody{}
	}

	var body authBody
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Errorf("解析请求体失败: %s: %v", raw, err)
	}
	return body
}

// newTestSDK 启动模拟网关，默认把重试退避压到 1ms 以便快速断言重试次数
func newTestSDK(t *testing.T, handler http.HandlerFunc, opts ...Option) *SungrowSDK {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	sdk, err := NewSungrowSDK(Credentials{
		AppID:        "app",
		AppSecret:    "secret",
		UserAccount:  "user",
		UserPassword: "pass",
		BaseURL:      srv.URL,
	}, append([]Option{WithRetryBaseDelay(time.Millisecond)}, opts...)...)
	if err != nil {
		t.Fatalf("NewSungrowSDK: %v", err)
	}
	return sdk
}

func listReq() PowerStationListRequest {
	return PowerStationListRequest{PageRequest: PageRequest{CurPage: 1, Size: 10}}
}

// TestAutoAcquireAndCacheToken token 为空时首次调用先取 token，之后复用缓存
func TestAutoAcquireAndCacheToken(t *testing.T) {
	var logins, biz int32

	var mu sync.Mutex
	var bodies []authBody

	sdk := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PathLogin:
			atomic.AddInt32(&logins, 1)
			if got := decodeAuthBody(t, r); got.AppKey != "app" {
				t.Errorf("登录请求 appkey = %q, 期望 app", got.AppKey)
			}
			_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", loginData("tk-1"))))
		case PathGetPsList:
			atomic.AddInt32(&biz, 1)
			mu.Lock()
			bodies = append(bodies, decodeAuthBody(t, r))
			mu.Unlock()
			_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", `{"total":0,"list":[]}`)))
		default:
			http.NotFound(w, r)
		}
	})

	for i := 1; i <= 2; i++ {
		if err := sdk.GetPowerStationList(listReq()); err != nil {
			t.Fatalf("第 %d 次 GetPowerStationList: %v", i, err)
		}
	}

	if got := atomic.LoadInt32(&logins); got != 1 {
		t.Fatalf("登录次数 = %d, 期望 1（第二次调用应复用缓存 token）", got)
	}
	if got := atomic.LoadInt32(&biz); got != 2 {
		t.Fatalf("业务调用次数 = %d, 期望 2", got)
	}

	mu.Lock()
	defer mu.Unlock()
	for i, b := range bodies {
		if b.AppKey != "app" || b.Token != "tk-1" {
			t.Fatalf("第 %d 次业务请求体 = %+v, 期望 appkey=app token=tk-1", i+1, b)
		}
	}
	if sdk.UID() != 9527 {
		t.Fatalf("UID = %d, 期望 9527", sdk.UID())
	}
	if sdk.TokenExpiry().Before(time.Now()) {
		t.Fatalf("token 过期时间应晚于当前时间")
	}
}

// TestReacquireOnInvalidToken token 失效（E00003）时丢弃缓存、重新获取后重试成功
func TestReacquireOnInvalidToken(t *testing.T) {
	var logins, biz int32

	var mu sync.Mutex
	var tokens []string

	sdk := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PathLogin:
			n := atomic.AddInt32(&logins, 1)
			_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", loginData(fmt.Sprintf("tk-%d", n)))))
		case PathGetPsList:
			body := decodeAuthBody(t, r)
			mu.Lock()
			tokens = append(tokens, body.Token)
			mu.Unlock()

			if atomic.AddInt32(&biz, 1) == 1 {
				// 首次携带缓存 token 调用时，平台判定 token 已失效
				_, _ = w.Write([]byte(env(ResultCodeErTokenLoginInvalid, "token invalid", "")))
				return
			}
			_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", `{"total":0,"list":[]}`)))
		}
	})

	if err := sdk.GetPowerStationList(listReq()); err != nil {
		t.Fatalf("token 失效后应重新获取并重试成功: %v", err)
	}

	if got := atomic.LoadInt32(&logins); got != 2 {
		t.Fatalf("登录次数 = %d, 期望 2（token 失效后重新获取）", got)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(tokens) != 2 || tokens[0] != "tk-1" || tokens[1] != "tk-2" {
		t.Fatalf("业务请求携带的 token = %v, 期望 [tk-1 tk-2]", tokens)
	}
}

// TestPresetTokenReused 预置且未过期的 token 直接使用，不触发登录
func TestPresetTokenReused(t *testing.T) {
	var logins, biz int32

	var mu sync.Mutex
	var token string

	sdk := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLogin {
			atomic.AddInt32(&logins, 1)
		}
		atomic.AddInt32(&biz, 1)

		mu.Lock()
		token = decodeAuthBody(t, r).Token
		mu.Unlock()

		_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", `{"total":0,"list":[]}`)))
	}, WithAccessToken("tk-preset", "", time.Hour))

	if err := sdk.GetPowerStationList(listReq()); err != nil {
		t.Fatalf("GetPowerStationList: %v", err)
	}

	if got := atomic.LoadInt32(&logins); got != 0 {
		t.Fatalf("登录次数 = %d, 期望 0（有效 token 应直接复用）", got)
	}

	mu.Lock()
	defer mu.Unlock()
	if token != "tk-preset" {
		t.Fatalf("业务请求 token = %q, 期望 tk-preset", token)
	}
}

// TestExpiredPresetTokenTriggersLogin 预置 token 已过期时重新获取
func TestExpiredPresetTokenTriggersLogin(t *testing.T) {
	var logins int32

	sdk := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLogin {
			atomic.AddInt32(&logins, 1)
			_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", loginData("tk-new"))))
			return
		}
		if got := decodeAuthBody(t, r).Token; got != "tk-new" {
			t.Errorf("业务请求 token = %q, 期望 tk-new", got)
		}
		_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", `{"total":0,"list":[]}`)))
	}, WithAccessToken("tk-expired", "", 0))

	if err := sdk.GetPowerStationList(listReq()); err != nil {
		t.Fatalf("GetPowerStationList: %v", err)
	}

	if got := atomic.LoadInt32(&logins); got != 1 {
		t.Fatalf("登录次数 = %d, 期望 1（过期 token 应重新获取）", got)
	}
}

// TestConcurrentCallsLoginOnce 并发调用只登录一次，且都携带同一 token
func TestConcurrentCallsLoginOnce(t *testing.T) {
	var logins int32

	var mu sync.Mutex
	tokens := map[string]int{}

	sdk := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLogin {
			atomic.AddInt32(&logins, 1)
			time.Sleep(10 * time.Millisecond) // 拉长登录窗口，暴露并发重复登录
			_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", loginData("tk-shared"))))
			return
		}

		mu.Lock()
		tokens[decodeAuthBody(t, r).Token]++
		mu.Unlock()

		_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", `{"total":0,"list":[]}`)))
	})

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := sdk.GetPowerStationList(listReq()); err != nil {
				t.Errorf("GetPowerStationList: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := atomic.LoadInt32(&logins); got != 1 {
		t.Fatalf("登录次数 = %d, 期望 1（并发下应串行化并复用缓存）", got)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(tokens) != 1 || tokens["tk-shared"] != 8 {
		t.Fatalf("业务请求 token 分布 = %v, 期望全部为 tk-shared", tokens)
	}
}

// TestDisableAutoLoginWithoutToken 关闭自动登录且无 token 时直接报错且不发请求
func TestDisableAutoLoginWithoutToken(t *testing.T) {
	var calls int32

	sdk := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		http.NotFound(w, r)
	}, WithDisableAutoLogin())

	err := sdk.GetPowerStationList(listReq())
	if !errors.Is(err, ErrNotLoggedIn) {
		t.Fatalf("期望 ErrNotLoggedIn, 实际 %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("HTTP 请求次数 = %d, 期望 0", got)
	}
}

// TestLoginFailureNotRetried 密码错误导致登录失败时不重试，避免锁定账户
func TestLoginFailureNotRetried(t *testing.T) {
	var logins int32

	sdk := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&logins, 1)
		_, _ = w.Write([]byte(env(ResultCodeSuccess, "success",
			`{"token":"","login_state":"0"}`)))
	})

	err := sdk.GetPowerStationList(listReq())
	if err == nil {
		t.Fatal("期望登录失败错误")
	}
	if !strings.Contains(err.Error(), "密码错误") {
		t.Fatalf("错误信息应包含登录失败原因, 实际 %v", err)
	}
	if got := atomic.LoadInt32(&logins); got != 1 {
		t.Fatalf("登录次数 = %d, 期望 1（密码错误不可重试）", got)
	}
}

// TestBodyMustBePointer 请求体不是指针时显式报错，避免静默漏传 token
func TestBodyMustBePointer(t *testing.T) {
	sdk := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(env(ResultCodeSuccess, "success", `{"total":0,"list":[]}`)))
	})

	var out any
	err := sdk.do(PathGetPsList, PowerStationListRequest{}, &out)
	if err == nil || !strings.Contains(err.Error(), "结构体指针") {
		t.Fatalf("期望请求体类型错误, 实际 %v", err)
	}
}

// TestParseLoginState 登录状态解析与说明
func TestParseLoginState(t *testing.T) {
	cases := []struct {
		in   string
		want LoginState
		ok   bool
		hint string
	}{
		{"1", LoginStateOK, true, "登录成功"},
		{" 0 ", LoginStateWrongPassword, true, "密码错误"},
		{"-1", LoginStateAccountNotExists, true, "用户账户不存在"},
		{"2", LoginStateAccountLockByPassword, true, "多次密码输入错误导致账户被锁定"},
		{"5", LoginStateAccountLockByAdmin, true, "该账号被管理员锁定"},
		{"", 0, false, ""},
		{"ok", 0, false, ""},
	}

	for _, c := range cases {
		got, ok := ParseLoginState(c.in)
		if ok != c.ok || got != c.want {
			t.Fatalf("ParseLoginState(%q) = %v %v, 期望 %v %v", c.in, got, ok, c.want, c.ok)
		}
		if ok && got.String() != c.hint {
			t.Fatalf("LoginState(%d).String() = %q, 期望 %q", got, got.String(), c.hint)
		}
	}
}
