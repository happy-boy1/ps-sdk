package service

import "fmt"

// SyncResult 一次同步的结果
type SyncResult struct {
	Platform      Platform       // 平台
	Stations      int            // 平台返回的电站数
	SavedStations int            // 成功落库的电站数
	Devices       int            // 拉取到的设备总数
	SavedDevices  int            // 成功落库的设备数
	DeviceCounts  map[uint64]int // 各电站的设备数，key 为电站内部 ID
	Errors        []error        // 单项失败明细，不影响整体同步
	Warnings      []error        // 非致命告警，如某类设备接口未开通权限
}

// OK 是否无失败项
func (r *SyncResult) OK() bool { return len(r.Errors) == 0 }

// String 便于直接打印
func (r *SyncResult) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: 电站 %d/%d，设备 %d/%d，失败 %d 项，告警 %d 项",
		r.Platform.NameZH, r.SavedStations, r.Stations,
		r.SavedDevices, r.Devices, len(r.Errors), len(r.Warnings))
}

// newSyncResult 构造带初始化 map 的结果
func newSyncResult(platform Platform) *SyncResult {
	return &SyncResult{Platform: platform, DeviceCounts: map[uint64]int{}}
}
