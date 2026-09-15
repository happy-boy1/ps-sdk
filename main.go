package main

import (
	"fmt"
	"ps-sdk/service"
)

func main() {
	c := service.NewSolarmanClient()
	ps, err := c.GetPowerStationList(1, 200)
	if err != nil {
		panic(err)
	}
	for _, p := range ps {
		fmt.Println(p)
	}

	// if err != nil {
	// 	panic(err)
	// }

	// ps := pwls.PageList[1]

	// devs, err := sdk.GetDeviceListByPsID(
	// 	sungrow.DeviceListByPsIDRequest{
	// 		PageRequest: sungrow.PageRequest{
	// 			CurPage: 1,
	// 			Size:    100,
	// 		},
	// 		PsID: ps.PsId.String(),
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// dev := devs.PageList[0]

	// rtd, err := sdk.GetDeviceRealTimeData(
	// 	sungrow.DeviceRtdRequest{
	// 		PsKeyList:   []string{dev.PsKey},
	// 		PointIDList: []string{"1"},
	// 		DeviceType:  dev.DeviceType,
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// tools.WriteJsonFile("test.json", rtd)
}
