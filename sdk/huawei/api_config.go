package huawei

import (
	"fmt"
	"strings"
)

// BatteryModeQueryRequest 储能工作模式查询请求
type BatteryModeQueryRequest struct {
	PlantCode string `json:"plantCode"` // 电站编号，一次仅支持一个电站
}

// BatteryModeQueryResult 储能工作模式查询结果
type BatteryModeQueryResult struct {
	PlantCode                   string               `json:"plantCode"`                   // 电站 DN
	OperationMode               BatteryOperationMode `json:"operationMode"`               // 工作模式
	TOUOperationModeParamDetail *TOUOperationDetail  `json:"touOperationModeParamDetail"` // TOU 参数，非 TOU 时为 null
}

// TOUOperationDetail TOU 工作模式参数
type TOUOperationDetail struct {
	RedundantPVEnergyPriority        PVEnergyPriority            `json:"redundantPVEnergyPriority"`        // 多余 PV 能量优先级
	AllowedAcChargePower             Float64                     `json:"allowedAcChargePower"`             // 电网充电最大功率(kW)
	ChargingAndDischargingTimeWindow []ChargeDischargeTimeWindow `json:"chargingAndDischargingTimeWindow"` // 充放电时间窗口
}

// GetBatteryMode 查询电站储能工作模式，仅支持单一组网电站
func (sdk *FusionSolarSDK) GetBatteryMode(req BatteryModeQueryRequest) (*BatteryModeQueryResult, error) {
	if strings.TrimSpace(req.PlantCode) == "" {
		return nil, fmt.Errorf("plantCode 不能为空")
	}
	var out BatteryModeQueryResult
	if err := sdk.do(PathBatteryModeQuery, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
