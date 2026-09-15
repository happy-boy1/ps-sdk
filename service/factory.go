package service

import (
	"errors"
	"log"
	"ps-sdk/model"
	"ps-sdk/sdk/ginlong"
	"ps-sdk/sdk/huawei"
	"ps-sdk/sdk/solarman"
	"ps-sdk/sdk/sungrow"
)

type SolarmanClient struct {
	*solarman.SolarmanSDK
}

type SungrowClient struct {
	sungrow.SungrowSDK
}

type SolisClient struct {
	ginlong.SolisSDK
}

type FusionSolarClient struct {
	huawei.FusionSolarSDK
}

type ClientMethod interface {
	GetPowerStationList(page, size int) ([]model.PowerStation, error)
}

func (c *SolarmanClient) GetPowerStationList(page, size int) ([]model.PowerStation, error) {
	req := solarman.StationListRequest{
		PageRequest: solarman.PageRequest{
			Page: page,
			Size: size,
		},
	}
	stations, err := c.StationList(req)
	if err != nil {
		return nil, errors.New("获取电站列表时出错")
	}

	var powerStations []model.PowerStation
	for _, station := range stations.StationList {
		powerStation := station.PowerStation()
		powerStations = append(powerStations, *powerStation)
	}
	return powerStations, nil
}

func NewSolarmanClient() ClientMethod {
	sdk, err := solarman.NewSolarmanSDK(
		solarman.Credentials{
			AppID:       "2024071755661944",
			AppSecret:   "916c65304351bed64942e3de056289d1",
			Mobile:      "13092743221",
			CountryCode: "86",
			Password:    "20010923wsljj..",
			OrgID:       1574,
		},
		solarman.WithDebugf(log.Printf),
	)
	if err != nil {
		return nil
	}

	return &SolarmanClient{
		SolarmanSDK: sdk,
	}
}
