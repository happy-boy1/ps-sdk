package main

import (
	"log"
	"ps-sdk/pkg/tools"
	"ps-sdk/sdk/sungrow"
	"time"
)

func main() {
	sdk, err := sungrow.NewSungrowSDK(
		sungrow.Credentials{
			AppID:        "EFCACCA7D56F186180F0AB1B708662D9",
			AppSecret:    "0dq39edbjub2rhq7hff9r7qqpeguxk5p",
			UserAccount:  "js",
			UserPassword: "fxdlyw12.",
		},
		sungrow.WithDevMode(false),
		sungrow.WithDebugf(log.Printf),
		sungrow.WithAccessToken(
			"53753_sfxfqg4pp9g4xn5f2tunw5ytmry789sv658jk9z2y3gytt7kb5iyhb9ngqyitj6zmxw39qymh1in54n7ysu2hj1hvmna0cn8iwuepp7s3vqqa9vpn1v40exx7s6vkre4",
			"",
			time.Duration(24)*time.Hour,
		),
	)

	if err != nil {
		panic(err)
	}

	pwls, err := sdk.GetPowerStationList(sungrow.PowerStationListRequest{
		PageRequest: sungrow.PageRequest{
			CurPage: 1,
			Size:    100,
		},
	})

	if err != nil {
		panic(err)
	}

	ps := pwls.PageList[1]

	devs, err := sdk.GetDeviceListByPsID(
		sungrow.DeviceListByPsIDRequest{
			PageRequest: sungrow.PageRequest{
				CurPage: 1,
				Size:    100,
			},
			PsID: ps.PsId.String(),
		},
	)
	if err != nil {
		panic(err)
	}

	dev := devs.PageList[0]

	rtd, err := sdk.GetDeviceRealTimeData(
		sungrow.DeviceRtdRequest{
			PsKeyList:   []string{dev.PsKey},
			PointIDList: []string{"1"},
			DeviceType:  dev.DeviceType,
		},
	)
	if err != nil {
		panic(err)
	}
	tools.WriteJsonFile("test.json", rtd)
}
