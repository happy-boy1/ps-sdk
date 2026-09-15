package service

import (
	"testing"
	"time"
)

func TestPlatformByCode(t *testing.T) {
	cases := []struct {
		in   string
		want int16
	}{
		{"solarman", PlatformSolarman},
		{"SolarMan", PlatformSolarman},
		{"SOLARMAN", PlatformSolarman},
		{" huawei ", PlatformFusionSolar},
		{"fusionsolar", 0}, // 短码是 huawei，fusionsolar 只是名称
		{"sungrow", PlatformSungrow},
		{"ginlong", PlatformGinlong},
		{"huawei", PlatformFusionSolar},
	}

	for _, c := range cases {
		got, ok := PlatformByCode(c.in)
		if c.want == 0 {
			if ok {
				t.Fatalf("%q 不应被识别为平台", c.in)
			}
			continue
		}
		if !ok || got.ID != c.want {
			t.Fatalf("PlatformByCode(%q) = %v, %v; 期望 ID %d", c.in, got, ok, c.want)
		}
	}
}

func TestNormalizeCode(t *testing.T) {
	cases := map[string]string{
		"SolarMan":       "solarman",
		"fusion-solar":   "fusionsolar",
		" fusion_solar ": "fusionsolar",
		"Ginlong":        "ginlong",
	}
	for in, want := range cases {
		if got := normalizeCode(in); got != want {
			t.Fatalf("normalizeCode(%q) = %q, 期望 %q", in, got, want)
		}
	}
}

func TestDeviceTypeCatalog(t *testing.T) {
	defs := DeviceTypes()
	if len(defs) == 0 {
		t.Fatal("设备类型目录为空")
	}

	seenCode := map[string]bool{}
	for _, d := range defs {
		if d.Code == "" || d.Name == "" || d.ShortName == "" {
			t.Fatalf("设备类型字段不完整: %+v", d)
		}
		if seenCode[d.Code] {
			t.Fatalf("设备类型编码重复: %s", d.Code)
		}
		seenCode[d.Code] = true

		// 主键为 smallint，超出范围会导致建表失败
		if d.ID < 0 || d.ID > 32767 {
			t.Fatalf("设备类型 ID 超出 smallint 范围: %d", d.ID)
		}
		if d.ID != DeviceTypeUnknown {
			if _, ok := DeviceTypeDefByID(d.ID); !ok {
				t.Fatalf("DeviceTypeDefByID(%d) 未命中", d.ID)
			}
		}
		if _, ok := DeviceTypeDefByCode(d.Code); !ok {
			t.Fatalf("DeviceTypeDefByCode(%s) 未命中", d.Code)
		}
	}
}

func TestResolveDeviceType(t *testing.T) {
	cases := []struct {
		platform string
		origin   string
		want     int16
	}{
		{CodeSolarman, "INVERTER", DeviceTypeInverter},
		{CodeSolarman, "inverter", DeviceTypeInverter},
		{CodeSolarman, "DTU", DeviceTypeCollector},
		{CodeSolarman, "METER", DeviceTypeMeter},
		{CodeSungrow, "1", DeviceTypeInverter},
		{CodeSungrow, "9", DeviceTypeCollector},
		{CodeSungrow, "43", DeviceTypeBattery},
		{CodeFusionSolar, "1", DeviceTypeInverter},
		{CodeFusionSolar, "38", DeviceTypeInverter},
		{CodeFusionSolar, "60044", DeviceTypePVModule},
		{CodeFusionSolar, "17", DeviceTypeGridMeter},
		{CodeGinlong, "METER", DeviceTypeMeter},
		{CodeGinlong, "EMS", DeviceTypeEnergyManagement},
		{CodeSolarman, "NOT_A_TYPE", DeviceTypeUnknown},
		{CodeSolarman, "", DeviceTypeUnknown},
	}

	for _, c := range cases {
		got := resolveDeviceType(c.platform, c.origin)
		if got.ID != c.want {
			t.Fatalf("resolveDeviceType(%s, %q) = %d(%s), 期望 %d",
				c.platform, c.origin, got.ID, got.Code, c.want)
		}
	}
}

func TestStatusName(t *testing.T) {
	deviceCases := map[int8]string{
		DeviceOnline:  "在线",
		DeviceOffline: "离线",
		DeviceAlarm:   "告警",
		DeviceUnknown: "未知",
		99:            "未知",
	}
	for in, want := range deviceCases {
		if got := DeviceStatusName(in); got != want {
			t.Fatalf("DeviceStatusName(%d) = %q, 期望 %q", in, got, want)
		}
	}

	stationCases := map[int8]string{
		StationRunning:    "运行",
		StationStopped:    "停运",
		StationUnderConst: "在建",
		99:                "未知",
	}
	for in, want := range stationCases {
		if got := StationStatusName(in); got != want {
			t.Fatalf("StationStatusName(%d) = %q, 期望 %q", in, got, want)
		}
	}

	if stationStatusFromOnline(true) != StationRunning || stationStatusFromOnline(false) != StationStopped {
		t.Fatal("stationStatusFromOnline 映射异常")
	}
}

func TestTimeHelpers(t *testing.T) {
	if unixSeconds(0) != nil || unixSeconds(-1) != nil {
		t.Fatal("非正数秒时间戳应返回 nil")
	}
	if got := unixSeconds(1700000000); got == nil || !got.Equal(time.Unix(1700000000, 0)) {
		t.Fatalf("unixSeconds = %v", got)
	}
	if unixMillis(0) != nil {
		t.Fatal("0 毫秒时间戳应返回 nil")
	}
	if got := unixMillis(1700000000000); got == nil || !got.Equal(time.UnixMilli(1700000000000)) {
		t.Fatalf("unixMillis = %v", got)
	}

	for _, in := range []string{"2024-01-02", "2024-01-02 15:04:05", "2024-01-02T15:04:05+08:00"} {
		if dateTime(in) == nil {
			t.Fatalf("dateTime(%q) 解析失败", in)
		}
	}
	for _, in := range []string{"", "not-a-date", "2024/01/02"} {
		if dateTime(in) != nil {
			t.Fatalf("dateTime(%q) 应返回 nil", in)
		}
	}
}

func TestParseHelpers(t *testing.T) {
	if got := parseInt64(" 1308675217944611083 "); got != 1308675217944611083 {
		t.Fatalf("parseInt64 = %d", got)
	}
	if got := parseInt64("abc"); got != 0 {
		t.Fatalf("非法整数应返回 0，实际 %d", got)
	}
	if got := parseFloat("118.420860"); got != 118.42086 {
		t.Fatalf("parseFloat = %v", got)
	}
	if got := parseFloat(""); got != 0 {
		t.Fatalf("空串应返回 0，实际 %v", got)
	}
	if formatInt64(-1) != "-1" {
		t.Fatal("formatInt64 异常")
	}
	if !isDigits("130") || isDigits("13a") || isDigits("") {
		t.Fatal("isDigits 判定异常")
	}
	if firstNonEmpty("", "  ", "x", "y") != "x" {
		t.Fatal("firstNonEmpty 取值异常")
	}
	if trimOr("  ", "fallback") != "fallback" {
		t.Fatal("trimOr 取值异常")
	}
}

func TestOptionsDefaults(t *testing.T) {
	if got := (Options{}).defaults().Timeout; got != 30*time.Second {
		t.Fatalf("默认超时 = %v", got)
	}
	if got := (Options{Timeout: time.Minute}).defaults().Timeout; got != time.Minute {
		t.Fatalf("显式超时被覆盖: %v", got)
	}
}

func TestSupportedAdapters(t *testing.T) {
	// 四个平台都必须注册适配器，否则 service.Platforms 与实现会失配
	for _, p := range Platforms {
		if _, ok := adapterBuilders[p.Code]; !ok {
			t.Fatalf("平台 %s 未注册适配器", p.Code)
		}
	}
	if len(adapterBuilders) != len(Platforms) {
		t.Fatalf("适配器数量 %d 与平台数量 %d 不一致", len(adapterBuilders), len(Platforms))
	}
}

func TestNewAdapterRejectsBadInput(t *testing.T) {
	if _, err := NewAdapter(Credential{Code: "unknown"}, Options{}); err == nil {
		t.Fatal("未知平台应返回错误")
	}
	if _, err := NewAdapter(Credential{Code: CodeSolarman}, Options{}); err == nil {
		t.Fatal("凭据缺失应返回错误")
	}
	if _, err := NewAdapter(Credential{Code: CodeGinlong}, Options{}); err == nil {
		t.Fatal("凭据缺失应返回错误")
	}
}

func TestSyncResultString(t *testing.T) {
	result := newSyncResult(Platforms[0])
	result.Stations = 2
	result.SavedStations = 1
	result.Devices = 5
	result.SavedDevices = 3
	result.Errors = append(result.Errors, errForTest{})

	if result.OK() {
		t.Fatal("存在失败项时 OK 应为 false")
	}
	if got := result.String(); got == "" {
		t.Fatal("String 为空")
	}
}

type errForTest struct{}

func (errForTest) Error() string { return "test" }
