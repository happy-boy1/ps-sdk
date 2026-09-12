package huawei

import (
	"fmt"
	"strings"
)

// 有功功率控制方式取值（5.2.7 下发接口）
const (
	ActivePowerControlModeNoLimit = "0" // 无限制
	ActivePowerControlModeLimitKW = "6" // 限功率并网(kW)
)

// 限制方式取值（5.2.7 下发接口）
const (
	ActivePowerControlLimitTotalPower       = "0" // 总功率
	ActivePowerControlLimitSinglePhasePower = "1" // 单相功率
)

// ActivePowerControlMaxTasks 一次任务最多下发的电站数
const ActivePowerControlMaxTasks = 10

// ActivePowerControlTaskItem 单个电站的有功功率设置任务
type ActivePowerControlTaskItem struct {
	PlantCode   string                   `json:"plantCode"`             // 电站DN
	ControlMode string                   `json:"controlMode"`           // 控制方式：0无限制/6限功率并网
	ControlInfo *ActivePowerControlParam `json:"controlInfo,omitempty"` // 限功率并网(kW)时的设置参数
}

// ActivePowerControlParam 有功功率设置参数
type ActivePowerControlParam struct {
	MaxGridFeedInPower *Float64 `json:"maxGridFeedInPower,omitempty"` // 最大馈送电网功率(kW)，0 为有效值
	LimitationMode     string   `json:"limitationMode,omitempty"`     // 限制方式：0总功率/1单相功率
}

// ActivePowerControlTaskRequest 逆变器有功功率设置任务下发请求
type ActivePowerControlTaskRequest struct {
	Tasks []ActivePowerControlTaskItem `json:"tasks"` // 任务列表，最多10个电站
}

// ActivePowerControlTaskDispatch 单个电站的任务下发结果
type ActivePowerControlTaskDispatch struct {
	PlantCode string     `json:"plantCode"` // 电站DN
	Status    TaskStatus `json:"status"`    // 下发状态：RUNNING/FAIL
	Message   string     `json:"message"`   // 下发结果描述
}

// ActivePowerControlTaskResult 逆变器有功功率设置任务下发结果
type ActivePowerControlTaskResult struct {
	TaskID string                           `json:"taskId"` // 任务唯一ID，用于结果查询
	Result []ActivePowerControlTaskDispatch `json:"result"` // 各电站下发结果列表
}

// ActivePowerControlTask 下发逆变器有功功率设置任务，一次最多10个电站；
// failCode=1 为部分成功，failCode=2 返回错误但结果仍会填充
func (sdk *FusionSolarSDK) ActivePowerControlTask(req ActivePowerControlTaskRequest) (*ActivePowerControlTaskResult, error) {
	if n := len(req.Tasks); n == 0 {
		return nil, fmt.Errorf("tasks 不能为空")
	} else if n > ActivePowerControlMaxTasks {
		return nil, fmt.Errorf("tasks 一次最多 %d 个电站，当前 %d 个", ActivePowerControlMaxTasks, n)
	}
	for i, task := range req.Tasks {
		if strings.TrimSpace(task.PlantCode) == "" {
			return nil, fmt.Errorf("tasks[%d].plantCode 不能为空", i)
		}
		switch task.ControlMode {
		case ActivePowerControlModeNoLimit, ActivePowerControlModeLimitKW:
		default:
			return nil, fmt.Errorf("tasks[%d].controlMode 只支持 %q(无限制) 或 %q(限功率并网kW)",
				i, ActivePowerControlModeNoLimit, ActivePowerControlModeLimitKW)
		}
	}

	var out ActivePowerControlTaskResult
	if err := sdk.doTask(PathActivePowerControlTask, req, &out); err != nil {
		// failCode=2 时任务级结果仍已填充，便于定位失败电站
		return &out, err
	}
	return &out, nil
}

// ActivePowerControlTaskInfoRequest 逆变器有功功率设置任务查询请求
type ActivePowerControlTaskInfoRequest struct {
	TaskID string `json:"taskId"` // 任务唯一ID，16位，由下发接口返回
}

// ActivePowerControlDispatchResult 查询返回的单个电站任务执行情况
type ActivePowerControlDispatchResult struct {
	PlantCode   string                   `json:"plantCode"`             // 电站DN
	ControlMode string                   `json:"controlMode"`           // 控制方式：0无限制/6限功率并网
	Status      TaskStatus               `json:"status"`                // 任务状态：RUNNING/SUCCESS/FAIL
	Message     string                   `json:"message,omitempty"`     // 失败原因，仅FAIL时返回
	ControlInfo *ActivePowerControlParam `json:"controlInfo,omitempty"` // 控制方式为6时返回
}

// ActivePowerControlTaskInfoResult 逆变器有功功率设置任务查询结果
type ActivePowerControlTaskInfoResult struct {
	DispatchResult []ActivePowerControlDispatchResult `json:"dispatchResult"`    // 各电站任务执行情况
	StartTime      Time                               `json:"startTime"`         // 收到任务的时间(含时区)
	EndTime        Time                               `json:"endTime,omitempty"` // 任务完成时间，未完成返回null
}

// ActivePowerControlTaskInfo 按taskId查询逆变器有功功率设置任务执行情况
func (sdk *FusionSolarSDK) ActivePowerControlTaskInfo(req ActivePowerControlTaskInfoRequest) (*ActivePowerControlTaskInfoResult, error) {
	if strings.TrimSpace(req.TaskID) == "" {
		return nil, fmt.Errorf("taskId 不能为空")
	}

	var out ActivePowerControlTaskInfoResult
	if err := sdk.do(PathActivePowerControlTaskIf, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ActivePowerControlQueryMode 查询返回的逆变器有功功率控制模式
type ActivePowerControlQueryMode string

const (
	ActivePowerControlQueryNoLimit        ActivePowerControlQueryMode = "noLimit"                 // 无限制
	ActivePowerControlQueryZeroExport     ActivePowerControlQueryMode = "zeroExportLimitation"    // 零功率并网
	ActivePowerControlQueryLimitedPowerKW ActivePowerControlQueryMode = "limitedPowerGridKW"      // 限功率并网(kW)
	ActivePowerControlQueryLimitedPercent ActivePowerControlQueryMode = "limitedPowerGridPercent" // 限功率并网(%)
	ActivePowerControlQueryOther          ActivePowerControlQueryMode = "other"                   // 其他控制模式
)

// ActivePowerControlLimitMode 查询返回的限制方式
type ActivePowerControlLimitMode string

const (
	ActivePowerControlLimitModeTotalPower       ActivePowerControlLimitMode = "totalPower"       // 总功率
	ActivePowerControlLimitModeSinglePhasePower ActivePowerControlLimitMode = "singlePhasePower" // 单相功率
)

// ActivePowerControlZeroExportParam 零功率并网模式详细参数
type ActivePowerControlZeroExportParam struct {
	LimitationMode ActivePowerControlLimitMode `json:"limitationMode"` // 限制方式
}

// ActivePowerControlPercentParam 限功率并网(%)模式详细参数
type ActivePowerControlPercentParam struct {
	LimitationMode            ActivePowerControlLimitMode `json:"limitationMode"`            // 限制方式
	MaxGridFeedInPowerPercent Float64                     `json:"maxGridFeedInPowerPercent"` // 最大馈送电网功率(%)
}

// ActivePowerControlValueParam 限功率并网(kW)模式详细参数
type ActivePowerControlValueParam struct {
	LimitationMode          ActivePowerControlLimitMode `json:"limitationMode"`          // 限制方式
	MaxGridFeedInPowerValue Float64                     `json:"maxGridFeedInPowerValue"` // 最大馈送电网功率(kW)
}

// ActivePowerControlModeRequest 逆变器有功功率控制模式查询请求
type ActivePowerControlModeRequest struct {
	PlantCode string `json:"plantCode"` // 电站编号
}

// ActivePowerControlModeResult 逆变器有功功率控制模式查询结果
type ActivePowerControlModeResult struct {
	PlantCode                    string                             `json:"plantCode"`                    // 电站DN
	ControlMode                  ActivePowerControlQueryMode        `json:"controlMode"`                  // 有功功率控制模式
	ZeroExportLimitationParam    *ActivePowerControlZeroExportParam `json:"zeroExportLimitationParam"`    // 零功率并网参数，否则null
	LimitedPowerGridPercentParam *ActivePowerControlPercentParam    `json:"limitedPowerGridPercentParam"` // 限功率并网(%)参数，否则null
	LimitedPowerGridValueParam   *ActivePowerControlValueParam      `json:"limitedPowerGridValueParam"`   // 限功率并网(kW)参数，否则null
}

// ActivePowerControlMode 查询电站下逆变器的有功功率控制模式
func (sdk *FusionSolarSDK) ActivePowerControlMode(req ActivePowerControlModeRequest) (*ActivePowerControlModeResult, error) {
	if strings.TrimSpace(req.PlantCode) == "" {
		return nil, fmt.Errorf("plantCode 不能为空")
	}

	var out ActivePowerControlModeResult
	if err := sdk.do(PathActivePowerControlMode, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
