package main

import (
	"fmt"
	"log"
	"ps-sdk/sdk/solarman"
	"time"
)

func main() {
	sdk, err := solarman.NewSolarmanSDK(
		solarman.Credentials{
			AppID:       "2024071755661944",
			AppSecret:   "916c65304351bed64942e3de056289d1",
			Email:       "",
			Mobile:      "13092743221",
			CountryCode: "86",
			UserName:    "",
			Password:    "20010923wsljj..",
			OrgID:       1574,
		},
		solarman.WithLanguage(solarman.LangZh),
		solarman.WithTimeout(30*time.Second),
		solarman.WithDebugf(log.Printf),
	)

	if err != nil {
		panic(err)
	}

	stationListReq := solarman.StationListRequest{
		PageRequest: solarman.PageRequest{
			Page: 1,
			Size: 200,
		},
	}
	stations, err := sdk.StationList(stationListReq)
	if err != nil {
		log.Println("[Solarman]: 获取电站列表失败")
	}

	for _, station := range stations.StationList {
		stationAlertReq := solarman.StationAlertV2Request{
			PageRequest: solarman.PageRequest{
				Page: 1,
				Size: 100,
			},
			StationID: station.ID,
			StartTime: "2026-09-12",
			EndTime:   "2026-09-12",
		}
		alertList, err := sdk.StationAlertV2(stationAlertReq)
		if err != nil {
			log.Println("[SolarMan]: 获取电站的告警列表失败")
			continue
		}
		fmt.Println(station.Name)
		for _, alert := range alertList.StationAlertItems {
			fmt.Println(alert)
		}
		fmt.Println("-----------------------------------------------------------------------------")
	}
}
