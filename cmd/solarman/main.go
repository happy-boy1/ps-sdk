// SolarMan（小麦智电 / 小麦商家版）开放平台调用示例。
// 凭据取自 pkg/config/config.toml 的 [platform.solarman] 段，
// 也可用环境变量临时覆盖（PS_SOLARMAN_PASSWORD / PS_SOLARMAN_ACCOUNT 等）。
package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"ps-sdk/pkg/config"
	"ps-sdk/sdk/solarman"
	"ps-sdk/service"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal(err)
	}
	cred, ok := cfg.Secrets(service.CodeSolarman)
	if !ok {
		log.Fatal("config.toml 缺少 [platform.solarman] 配置")
	}

	sdk, err := solarman.NewSolarmanSDK(
		solarman.Credentials{
			AppID:       cred.AppKey,
			AppSecret:   cred.AppSecret,
			Email:       cred.Email,
			Mobile:      cred.Mobile,
			UserName:    cred.UserName,
			CountryCode: cred.CountryCode,
			Password:    cred.Password,
			OrgID:       cred.OrgID,
			BaseURL:     cred.APIURL,
		},
		solarman.WithLanguage(solarman.LangZh),
		solarman.WithTimeout(cfg.HTTP.Timeout.Duration()),
		solarman.WithDebugf(log.Printf),
	)
	if err != nil {
		fatal(err)
	}

	// 1. 获取 Token（也可不显式调用，首次请求会自动获取）
	res, err := sdk.AcquireToken(solarman.TokenRequest{})
	if err != nil {
		fatal(err)
	}
	fmt.Printf("Token 获取成功，uid=%d，有效期至 %s\n",
		res.UID.Int(), sdk.TokenExpiry().Format(time.RFC3339))

	// 2. 账号商家关系（商家版用它拿 orgId；拿到后可 sdk.LoginWithOrg(id) 切换）
	orgs, err := sdk.AccountInfo()
	if err != nil {
		fatal(err)
	}
	for _, o := range orgs.OrgInfoList {
		fmt.Printf("  商家 %d %s 角色 %s\n", o.CompanyID.Int(), o.CompanyName, o.RoleName)
	}

	// 3. 账号下电站列表
	stations, err := sdk.StationList(solarman.StationListRequest{
		PageRequest: solarman.PageRequest{Page: 1, Size: 20},
	})
	if err != nil {
		fatal(err)
	}
	fmt.Printf("共 %d 个电站\n", stations.Total)
	for _, s := range stations.StationList {
		fmt.Printf("  %d %s %s %.2fkW %s %s\n",
			stationID(s.ID), s.Name, s.LocationAddress,
			s.InstalledCapacity.Float(), s.Type, s.NetworkStatus)
	}
	if len(stations.StationList) == 0 {
		return
	}
	station := stations.StationList[0]

	// 4. 电站实时数据
	rt, err := sdk.StationRealTime(solarman.StationRealTimeRequest{StationID: stationID(station.ID)})
	if err != nil {
		fatal(err)
	}
	fmt.Printf("电站实时: 发电功率 %.2f 并网功率 %.2f 用电功率 %.2f\n",
		rt.GenerationPower.Float(), rt.GridPower.Float(), rt.UsePower.Float())

	// 5. 电站下的设备列表
	devices, err := sdk.StationDeviceList(solarman.StationDeviceListRequest{
		PageRequest: solarman.PageRequest{Page: 1, Size: 50},
		StationID:   stationID(station.ID),
	})
	if err != nil {
		fatal(err)
	}
	fmt.Printf("设备 %d 台\n", devices.Total)
	for _, d := range devices.DeviceListItems {
		fmt.Printf("  %s %s 状态 %d 更新 %.0f\n", d.DeviceSN, d.DeviceType, d.ConnectStatus, d.CollectionTime.Float())
	}
	if len(devices.DeviceListItems) == 0 {
		return
	}
	device := devices.DeviceListItems[0]

	// 6. 设备实时数据（返回数据点表：key / value / unit / name）
	cur, err := sdk.CurrentData(solarman.CurrentDataRequest{
		DeviceSN: device.DeviceSN,
		DeviceID: device.DeviceID,
	})
	if err != nil {
		fatal(err)
	}
	fmt.Printf("设备 %s 实时数据 %d 个数据点\n", cur.DeviceSN, len(cur.DataList))
	for i, item := range cur.DataList {
		if i >= 10 {
			fmt.Printf("  ... 其余 %d 个\n", len(cur.DataList)-10)
			break
		}
		fmt.Printf("  %-16s %-12s %s\n", item.Key, item.Value, item.Unit)
	}

	// 7. 设备历史数据（TimeType：2 日 / 3 月 / 4 年 / 5 帧，时间格式随之变化）
	hist, err := sdk.DeviceHistorical(solarman.DeviceHistoricalRequest{
		DeviceSN:  device.DeviceSN,
		DeviceID:  device.DeviceID,
		StartTime: time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
		EndTime:   time.Now().Format("2006-01-02"),
		TimeType:  2,
	})
	if err != nil {
		fatal(err)
	}
	fmt.Printf("历史数据 %d 组\n", len(hist.ParamDataList))

	// 8. 设备通讯关系（网关与子设备的层级结构）
	comm, err := sdk.DeviceCommunication(solarman.DeviceCommunicationRequest{DeviceSN: device.DeviceSN})
	if err != nil {
		log.Printf("通讯关系: %v", err)
	} else if comm.Communication != nil {
		fmt.Printf("通讯关系根节点: %+v\n", comm.Communication)
	}

	// 9. 电站报警列表
	alerts, err := sdk.StationAlert(solarman.StationAlertRequest{
		PageRequest: solarman.PageRequest{Page: 1, Size: 20},
		StationID:   stationID(station.ID),
		StartTime:   time.Now().AddDate(0, 0, -30).Format("2006-01-02"),
		EndTime:     time.Now().Format("2006-01-02"),
	})
	if err != nil {
		log.Printf("电站报警: %v", err)
	} else {
		fmt.Printf("近 30 天报警 %d 条\n", alerts.Total)
	}

	// 10. 电站实时天气
	if w, err := sdk.StationWeather(solarman.StationWeatherRequest{StationID: stationID(station.ID)}); err != nil {
		log.Printf("天气: %v", err)
	} else if v, ok := w.Raw["temperature"]; ok {
		fmt.Printf("电站所在地温度: %v\n", v)
	}

	// 11. APPID 剩余可调用次数（调用会扣次数，勿频繁调用）
	balance, err := sdk.AppIDBalance(solarman.AppIDBalanceRequest{AppID: sdk.Credentials().AppID})
	if err != nil {
		log.Printf("余量查询: %v", err)
	} else {
		fmt.Printf("APPID 剩余可调用次数: %d\n", balance.Total.Int())
	}

	// 12. 未建模字段从 Raw 兜底（如 4.4 示例里的 batterySoc）
	if v, ok := station.Raw.Float64("batterySoc"); ok {
		fmt.Printf("原始数据兜底: batterySoc = %.1f\n", v)
	}

	// 13. 调用任意接口拿原始数据
	raw, err := sdk.Call(solarman.PathStationList, map[string]any{"page": 1, "size": 5})
	if err == nil {
		fmt.Printf("原始响应字段数: %d\n", len(raw))
	}

	// 14. 写操作（会变更账号数据，按需打开）
	// _, err = sdk.StationCreate(solarman.StationCreateRequest{...})
	// _, err = sdk.DeviceRegister(solarman.DeviceRegisterRequest{StationID: stationID(station.ID), DeviceSN: "xxx", IsAuto: true})
}

// stationID 把电站 ID 从响应里的字符串形态转成请求需要的数值形态。
// 4.4 的 id 字段类型不固定，SDK 用 Str 兜底，发送请求时需显式转换
func stationID(id solarman.Str) solarman.Int64 {
	v, err := strconv.ParseInt(id.String(), 10, 64)
	if err != nil {
		log.Fatalf("电站 ID %q 无法解析为数值: %v", id.String(), err)
	}
	return solarman.Int64(v)
}

func fatal(err error) {
	var apiErr *solarman.APIError
	if errors.As(err, &apiErr) {
		log.Fatalf("code=%s msg=%s hint=%s", apiErr.Code, apiErr.Msg,
			solarman.CodeHint(string(apiErr.Code)))
	}
	log.Fatal(err)
}
