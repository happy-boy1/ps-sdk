// 锦浪云平台 API 调用示例。
// 凭据取自 pkg/config/config.toml 的 [platform.ginlong] 段（app_id/app_secret 由数据库提供，
// 也可用 PS_GINLONG_APPID / PS_GINLONG_APPSECRET 环境变量临时覆盖）。
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"ps-sdk/pkg/config"
	"ps-sdk/sdk/ginlong"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal(err)
	}
	cred, ok := cfg.Credential("ginlong")
	if !ok {
		log.Fatal("config.toml 缺少 [platform.ginlong] 配置")
	}

	sdk, err := ginlong.NewSolisSDK(
		ginlong.Credentials{
			APIID:     cred.AppID,
			APISecret: cred.AppSecret,
			BaseURL:   cred.APIURL,
		},
		ginlong.WithTimeout(cfg.HTTP.Timeout.Duration()),
		ginlong.WithDebugf(log.Printf),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 1. 账号下电站列表
	stations, err := sdk.UserStationList(ginlong.UserStationListRequest{
		PageRequest: ginlong.PageRequest{PageNo: 1, PageSize: 20},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("共 %d 个电站（在线 %d / 离线 %d / 故障 %d）\n",
		stations.StationStatusVo.All, stations.StationStatusVo.Normal,
		stations.StationStatusVo.Offline, stations.StationStatusVo.Fault)
	for _, s := range stations.Page.Records {
		fmt.Printf("  %d %s %s %.2f%s 当日 %.2f%s 状态 %s\n",
			s.ID.Int(), s.StationName, s.Addr,
			s.Capacity.Float(), s.CapacityStr,
			s.DayEnergy.Float(), s.DayEnergyStr, s.State)
	}
	if len(stations.Page.Records) == 0 {
		return
	}
	st := stations.Page.Records[0]

	// 2. 电站详情
	detail, err := sdk.StationDetail(ginlong.StationDetailRequest{ID: st.ID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("电站详情: %s 电站类型 %v 装机 %.2f%s\n",
		detail.StationName, detail.StationTypeNew, detail.Capacity.Float(), detail.CapacityStr)

	// 3. 该电站下的逆变器列表
	inverters, err := sdk.InverterList(ginlong.InverterListRequest{
		PageRequest: ginlong.PageRequest{PageNo: 1, PageSize: 100},
		StationID:   st.ID,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("逆变器 %d 台（在线 %d）\n", inverters.InverterStatusVo.All, inverters.InverterStatusVo.Normal)
	for _, inv := range inverters.Page.Records {
		fmt.Printf("  %s %s %s %.2f%s 功率 %.2f%s 状态 %s\n",
			inv.SN, inv.Name, inv.ProductModel,
			inv.Etoday.Float(), inv.EtodayStr,
			inv.Pac.Float(), inv.PacStr, inv.State)
	}
	if len(inverters.Page.Records) == 0 {
		return
	}
	inv := inverters.Page.Records[0]

	// 4. 单台逆变器详情（含储能/电表数据）
	d, err := sdk.InverterDetail(ginlong.InverterDetailRequest{ID: inv.ID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("逆变器详情: %s 当日 %.2f%s 总发电 %.2f%s 功率 %.2fkW 温度 %.1f℃ SOC %.0f%%\n",
		d.SN, d.EToday.Float(), d.ETodayStr, d.ETotal.Float(), d.ETotalStr,
		d.Pac.Float(), d.InverterTemperature.Float(), d.BatteryCapacitySoc.Float())

	// 未建模字段（iPv1..32、uPv1..32、pow1..32 等）从原始数据取
	if v, ok := d.Raw.Float64("iPv17"); ok {
		fmt.Printf("  原始数据兜底: iPv17 = %.2f\n", v)
	}

	// 5. 某日曲线（time 为 yyyy-MM-dd，timeZone 为电站时区）
	today := time.Now().Format("2006-01-02")
	curve, err := sdk.InverterDay(ginlong.InverterDayRequest{
		ID: inv.ID, Money: "CNY", Time: today, TimeZone: 8,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("当日曲线 %d 个采样点\n", len(curve))
	if len(curve) > 0 {
		last := curve[len(curve)-1]
		fmt.Printf("  最新采样 %s 功率 %.2f%s 当日电量 %.2f\n",
			last.TimeStr, last.Pac.Float(), last.PacStr, last.EToday.Float())
	}

	// 6. 电站某日曲线
	day, err := sdk.StationDay(ginlong.StationDayRequest{
		ID: st.ID, Money: "CNY", Time: today, TimeZone: 8,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("电站当日曲线 %d 个采样点\n", len(day))

	// 7. 设备报警列表
	alarms, err := sdk.AlarmList(ginlong.AlarmListRequest{
		PageRequest:    ginlong.PageRequest{PageNo: 1, PageSize: 20},
		StationID:      st.ID,
		AlarmBeginTime: time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
		AlarmEndTime:   today,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("近 7 天报警 %d 条\n", len(alarms.Page.Records))
	for _, a := range alarms.Page.Records {
		fmt.Printf("  [%s] %s %s %s\n", a.AlarmLevel, a.AlarmDeviceSn, a.AlarmCode, a.AlarmMsg)
	}

	// 8. 任意接口拿原始数据（未建模的接口/字段）
	raw, err := sdk.Call(ginlong.PathInverterList, ginlong.InverterListRequest{
		PageRequest: ginlong.PageRequest{PageNo: 1, PageSize: 10},
	})
	if err != nil {
		var apiErr *ginlong.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("code=%s hint=%s", apiErr.Code, ginlong.CodeHint(string(apiErr.Code)))
		}
		log.Fatal(err)
	}
	fmt.Printf("原始数据 keys: %d 个\n", len(raw))
}
