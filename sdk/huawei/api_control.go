package huawei

import "fmt"

// ExecStatus 储能充放电任务执行状态
type ExecStatus int

const (
	ExecStatusDone    ExecStatus = 0 // 完成
	ExecStatusRunning ExecStatus = 1 // 执行中
	ExecStatusFailed  ExecStatus = 2 // 超时或失败
)

// ChargeDischargeTask 单个电站的储能充放电任务
type ChargeDischargeTask struct {
	PlantCode      string              `json:"plantCode"`               // 电站编码
	DispatchSwitch DispatchSwitch      `json:"dispatchSwitch"`          // 0停止 1强制充电 2强制放电
	ControlType    DispatchControlType `json:"controlType,omitempty"`   // 1 SOC 控制 2 时间控制
	TargetSOC      Float64             `json:"targetSOC,omitempty"`     // 目标 SOC(%)，SOC 控制必填
	DispatchTime   int                 `json:"dispatchTime,omitempty"`  // 时长(min)，时间控制必填
	PowerDispatch  int                 `json:"powerDispatch,omitempty"` // 功率(W)，充电为正放电为负
}

// ChargeDischargeTaskRequest 储能充放电任务下发请求
type ChargeDischargeTaskRequest struct {
	Tasks []ChargeDischargeTask `json:"tasks"` // 任务列表，一次最多 100 个电站
}

// ChargeDischargeDispatchItem 单电站下发结果
type ChargeDischargeDispatchItem struct {
	PlantCode      string `json:"plantCode"`      // 电站编号
	Sn             string `json:"sn"`             // 逆变器 SN
	DispatchResult int    `json:"dispatchResult"` // 0 接收成功 1 接收失败
	SubTaskID      string `json:"subTaskId"`      // 子任务唯一 ID
	Description    string `json:"description"`    // 下发结果描述
}

// ChargeDischargeTaskResult 储能充放电任务下发结果
type ChargeDischargeTaskResult struct {
	RequestID int64                         `json:"requestID"` // 请求任务唯一 ID，用于查询结果
	Result    []ChargeDischargeDispatchItem `json:"result"`    // 各电站下发结果
}

// SubmitChargeDischargeTask 下发储能充放电任务。不改变储能原始工作模式，任务结束即恢复
func (sdk *FusionSolarSDK) SubmitChargeDischargeTask(req ChargeDischargeTaskRequest) (*ChargeDischargeTaskResult, error) {
	if len(req.Tasks) == 0 {
		return nil, fmt.Errorf("tasks 不能为空")
	}
	if len(req.Tasks) > MaxBatch {
		return nil, fmt.Errorf("tasks 最多 %d 个电站，当前 %d 个", MaxBatch, len(req.Tasks))
	}
	for i, t := range req.Tasks {
		if t.PlantCode == "" {
			return nil, fmt.Errorf("tasks[%d].plantCode 不能为空", i)
		}
		if t.DispatchSwitch < DispatchStop || t.DispatchSwitch > DispatchDischg {
			return nil, fmt.Errorf("tasks[%d].dispatchSwitch 只能为 0/1/2", i)
		}
		if t.DispatchSwitch == DispatchStop {
			continue
		}
		switch t.ControlType {
		case ControlBySOC:
			if t.TargetSOC <= 0 {
				return nil, fmt.Errorf("tasks[%d] SOC 控制需填 targetSOC", i)
			}
		case ControlByTime:
			if t.DispatchTime <= 0 {
				return nil, fmt.Errorf("tasks[%d] 时间控制需填 dispatchTime", i)
			}
		default:
			return nil, fmt.Errorf("tasks[%d].controlType 只能为 1(SOC)/2(时间)", i)
		}
	}
	var out ChargeDischargeTaskResult
	if err := sdk.doTask(PathChargeDischargeTask, req, &out); err != nil {
		return &out, err
	}
	return &out, nil
}

// ChargeDischargeStatusRequest 储能充放电任务查询请求
type ChargeDischargeStatusRequest struct {
	RequestID int64 `json:"requestID"` // 下发接口返回的任务唯一 ID
}

// ChargeDischargeStatus 储能充放电任务执行情况
type ChargeDischargeStatus struct {
	PlantCode          string     `json:"plantCode"`          // 电站编号
	Sn                 string     `json:"sn"`                 // 逆变器 SN
	RemoteID           string     `json:"remoteID"`           // 子任务唯一 ID
	Status             ExecStatus `json:"status"`             // 0完成 1执行中 2超时或失败
	ChargedCapacity    Float64    `json:"chargedCapacity"`    // 已强制充电电量(kWh)
	DischargedCapacity Float64    `json:"dischargedCapacity"` // 已强制放电电量(kWh)
	ExecStartTime      Time       `json:"execStartTime"`      // 收到任务的时间
	ExecEndTime        Time       `json:"execEndTime"`        // 任务完成时间(未完成返回 null)
	Message            string     `json:"message"`            // 任务结果描述
}

// GetChargeDischargeStatus 查询储能充放电任务执行情况。状态每 3 分钟刷新，24 小时未完成即超时
func (sdk *FusionSolarSDK) GetChargeDischargeStatus(req ChargeDischargeStatusRequest) ([]ChargeDischargeStatus, error) {
	if req.RequestID <= 0 {
		return nil, fmt.Errorf("requestID 不能为空")
	}
	var out []ChargeDischargeStatus
	if err := sdk.do(PathChargeDischargeStatus, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
