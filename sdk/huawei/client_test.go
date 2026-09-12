package huawei

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// newTestSDK 启动一个模拟北向接口的服务端
func newTestSDK(t *testing.T, handler http.HandlerFunc) (*FusionSolarSDK, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	sdk, err := NewFusionSolarSDK(Credentials{
		UserName:   "u",
		SystemCode: "p",
		BaseURL:    srv.URL,
	}, WithRetryBaseDelay(time.Millisecond))
	if err != nil {
		t.Fatalf("NewFusionSolarSDK: %v", err)
	}
	return sdk, srv
}

func TestLoginAndStationList(t *testing.T) {
	var logins, stations int32

	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PathLogin:
			atomic.AddInt32(&logins, 1)
			w.Header().Set(HeaderToken, "token-1")
			_, _ = w.Write([]byte(`{"success":true,"failCode":0,"data":null}`))
		case PathStations:
			atomic.AddInt32(&stations, 1)
			if r.Header.Get(HeaderToken) != "token-1" {
				t.Errorf("缺少 XSRF-TOKEN 请求头")
			}
			_, _ = w.Write([]byte(`{"success":true,"failCode":0,"data":{
				"total":1,"pageCount":1,"pageNo":1,"pageSize":100,
				"list":[{"plantCode":"NE=1","plantName":"p1","capacity":100.5}]}}`))
		default:
			http.NotFound(w, r)
		}
	})

	res, err := sdk.GetStationList(StationListRequest{PageNo: 1})
	if err != nil {
		t.Fatalf("GetStationList: %v", err)
	}
	if res.Total != 1 || len(res.List) != 1 || res.List[0].PlantCode != "NE=1" {
		t.Fatalf("结果异常: %+v", res)
	}
	if got := atomic.LoadInt32(&logins); got != 1 {
		t.Fatalf("登录次数 = %d, 期望 1", got)
	}

	// 第二次调用应复用 token
	if _, err := sdk.GetStationList(StationListRequest{PageNo: 1}); err != nil {
		t.Fatalf("GetStationList #2: %v", err)
	}
	if got := atomic.LoadInt32(&logins); got != 1 {
		t.Fatalf("登录次数 = %d, 期望 1（应复用 token）", got)
	}
}

func TestReloginOn305(t *testing.T) {
	var logins, stations int32

	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PathLogin:
			n := atomic.AddInt32(&logins, 1)
			w.Header().Set(HeaderToken, "token-"+string(rune('0'+n)))
			_, _ = w.Write([]byte(`{"success":true,"failCode":0}`))
		case PathStations:
			if atomic.AddInt32(&stations, 1) == 1 {
				_, _ = w.Write([]byte(`{"success":false,"failCode":305,"message":"not login"}`))
				return
			}
			_, _ = w.Write([]byte(`{"success":true,"failCode":0,"data":{"total":0,"list":[]}}`))
		}
	})

	if _, err := sdk.GetStationList(StationListRequest{PageNo: 1}); err != nil {
		t.Fatalf("GetStationList: %v", err)
	}
	if got := atomic.LoadInt32(&logins); got != 2 {
		t.Fatalf("登录次数 = %d, 期望 2（305 后重新登录）", got)
	}
}

func TestBusinessError(t *testing.T) {
	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLogin {
			w.Header().Set(HeaderToken, "t")
			_, _ = w.Write([]byte(`{"success":true,"failCode":0}`))
			return
		}
		_, _ = w.Write([]byte(`{"success":false,"failCode":401,"message":"no permission"}`))
	})

	_, err := sdk.GetStationList(StationListRequest{PageNo: 1})
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("期望 *APIError, 实际 %T: %v", err, err)
	}
	if apiErr.FailCode != FailCodeNoPermission {
		t.Fatalf("failCode = %d", apiErr.FailCode)
	}
	if apiErr.Retryable() {
		t.Fatalf("401 不应重试")
	}
}

func TestTaskPartialSuccess(t *testing.T) {
	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLogin {
			w.Header().Set(HeaderToken, "t")
			_, _ = w.Write([]byte(`{"success":true,"failCode":0}`))
			return
		}
		_, _ = w.Write([]byte(`{"success":true,"failCode":1,"message":"partial","data":{
			"requestID":1107541629693466,
			"result":[{"plantCode":"NE=1","dispatchResult":0,"subTaskId":"s1"},
			          {"plantCode":"NE=2","dispatchResult":1,"description":"illegal"}]}}`))
	})

	res, err := sdk.SubmitChargeDischargeTask(ChargeDischargeTaskRequest{Tasks: []ChargeDischargeTask{
		{PlantCode: "NE=1", DispatchSwitch: DispatchCharge, ControlType: ControlByTime, DispatchTime: 600, PowerDispatch: 5000},
		{PlantCode: "NE=2", DispatchSwitch: DispatchCharge, ControlType: ControlByTime, DispatchTime: 600, PowerDispatch: 5000},
	}})
	if err != nil {
		t.Fatalf("部分成功不应返回错误: %v", err)
	}
	if res.RequestID != 1107541629693466 || len(res.Result) != 2 {
		t.Fatalf("结果异常: %+v", res)
	}
}

func TestItemMapHelpers(t *testing.T) {
	// 华为实时数据以字符串/null 返回
	m := ItemMap{
		KpiDayPower:        "10000",
		KpiTotalPower:      nil,
		KpiRealHealthState: "3",
		"active_power":     12.5,
	}
	if v, ok := m.Float64(KpiDayPower); !ok || v != 10000 {
		t.Fatalf("day_power = %v %v", v, ok)
	}
	if _, ok := m.Float64(KpiTotalPower); ok {
		t.Fatalf("null 应视为无效值")
	}
	if v, ok := m.Int(KpiRealHealthState); !ok || v != 3 {
		t.Fatalf("real_health_state = %v %v", v, ok)
	}
	if v, ok := m.Float64("active_power"); !ok || v != 12.5 {
		t.Fatalf("active_power = %v %v", v, ok)
	}

	kpi := ItemMap{KpiDayPower: "123.45", KpiRealHealthState: "2"}.ParseStationRealtime()
	if kpi.DayPower != 123.45 || kpi.RealHealthState != HealthFault {
		t.Fatalf("ParseStationRealtime = %+v", kpi)
	}
}

func TestValidation(t *testing.T) {
	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {})

	if _, err := sdk.GetDevRealKpi(DeviceRealKpiRequest{DevTypeID: DevTypeInverter}); err == nil {
		t.Fatal("devIds/sns 均空应报错")
	}
	if _, err := sdk.GetStationKpiDay(StationKpiRequest{StationCodes: "NE=1"}); err == nil {
		t.Fatal("collectTime 为空应报错")
	}
	if _, err := sdk.GetAlarmList(AlarmListRequest{BeginTime: 1, EndTime: 2, Language: string(LangZhCN)}); err == nil {
		t.Fatal("stationCodes/sns 均空应报错")
	}
	if _, err := sdk.GetBatteryMode(BatteryModeQueryRequest{}); err == nil {
		t.Fatal("plantCode 为空应报错")
	}
}

func TestLoginFailureNotRetried(t *testing.T) {
	var logins int32

	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&logins, 1)
		_, _ = w.Write([]byte(`{"success":false,"failCode":20400,"message":"user.login.user_or_value_invalid"}`))
	})

	if _, err := sdk.GetStationList(StationListRequest{PageNo: 1}); err == nil {
		t.Fatal("期望登录失败错误")
	}
	if got := atomic.LoadInt32(&logins); got != 1 {
		t.Fatalf("登录请求次数 = %d, 期望 1（密码错误不可重试，否则会锁定账户）", got)
	}
}

func TestExceptionNotRetried(t *testing.T) {
	var calls int32

	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLogin {
			w.Header().Set(HeaderToken, "t")
			_, _ = w.Write([]byte(`{"success":true,"failCode":0}`))
			return
		}
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`{"exceptionId":"framwork.remote.Paramerror",
			"exceptionType":"ROA_EXFRAME_EXCEPTION","descArgs":null,
			"reasonArgs":["plantCode"],"detailArgs":["plantCode may not be null"]}`))
	})

	_, err := sdk.GetStationList(StationListRequest{PageNo: 1})
	var exc *Exception
	if !errors.As(err, &exc) {
		t.Fatalf("期望 *Exception, 实际 %T: %v", err, err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("业务接口调用次数 = %d, 期望 1（网关异常不可重试）", got)
	}
}

func TestRetryOnServerBusy(t *testing.T) {
	var calls int32

	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLogin {
			w.Header().Set(HeaderToken, "t")
			_, _ = w.Write([]byte(`{"success":true,"failCode":0}`))
			return
		}
		if atomic.AddInt32(&calls, 1) < 3 {
			_, _ = w.Write([]byte(`{"success":false,"failCode":20004,"message":"server error"}`))
			return
		}
		_, _ = w.Write([]byte(`{"success":true,"failCode":0,"data":{"total":0,"list":[]}}`))
	})

	if _, err := sdk.GetStationList(StationListRequest{PageNo: 1}); err != nil {
		t.Fatalf("20004 应重试至成功: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Fatalf("调用次数 = %d, 期望 3", got)
	}
}

func TestFloat64Forms(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{`1.5`, 1.5},
		{`"1.5"`, 1.5},
		{`"32.326424"`, 32.326424},
		{`null`, 0},
		{`0`, 0},
	}
	for _, c := range cases {
		var f Float64
		if err := json.Unmarshal([]byte(c.in), &f); err != nil {
			t.Fatalf("解析 %s 失败: %v", c.in, err)
		}
		if f.Float() != c.want {
			t.Fatalf("%s = %v, 期望 %v", c.in, f, c.want)
		}
	}
	// 序列化仍输出数值，不影响请求体
	b, err := json.Marshal(Float64(1.5))
	if err != nil || string(b) != "1.5" {
		t.Fatalf("Marshal = %s, err = %v", b, err)
	}
}

// TestStationRealResponse 复现线上响应：Double 字段以字符串返回
func TestStationRealResponse(t *testing.T) {
	sdk, _ := newTestSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLogin {
			w.Header().Set(HeaderToken, "t")
			_, _ = w.Write([]byte(`{"data":null,"success":true,"failCode":0,"params":{}}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":{"list":[{"capacity":6700.68,"contactMethod":"15315609546",` +
			`"gridConnectionDate":"2026-08-26T10:45:42+08:00","latitude":"32.326424","longitude":"118.420860",` +
			`"plantAddress":"安徽省滁州市琅琊区镇江路198号","plantCode":"NE=371670079","plantName":"滁州东泰光伏站"}],` +
			`"pageCount":1,"pageNo":1,"pageSize":100,"total":1},"failCode":0,` +
			`"message":"get plant list success","success":true}`))
	})

	res, err := sdk.GetStationList(StationListRequest{PageNo: 1})
	if err != nil {
		t.Fatalf("GetStationList: %v", err)
	}
	if res.Total != 1 || len(res.List) != 1 {
		t.Fatalf("分页信息异常: %+v", res)
	}
	s := res.List[0]
	if s.Latitude.Float() != 32.326424 || s.Longitude.Float() != 118.420860 {
		t.Fatalf("经纬度解析异常: %v %v", s.Latitude, s.Longitude)
	}
	if s.Capacity.Float() != 6700.68 {
		t.Fatalf("capacity = %v", s.Capacity)
	}
	if s.PlantAddress == "" || s.PlantName != "滁州东泰光伏站" {
		t.Fatalf("字段解析异常: %+v", s)
	}
}

func TestJoinCodesTruncates(t *testing.T) {
	codes := make([]string, MaxBatch+10)
	for i := range codes {
		codes[i] = "NE=1"
	}
	if n := len(SplitCodes(JoinCodes(codes))); n != MaxBatch {
		t.Fatalf("JoinCodes 截断后数量 = %d, 期望 %d", n, MaxBatch)
	}
}
