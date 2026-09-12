package ginlong

import (
	"fmt"
	"strings"
)

// AddStationRequest 新增电站绑定逆变器SN请求
type AddStationRequest struct {
	InverterSN      string `json:"inverterSn,omitempty"`      // 逆变器SN，填写即可绑定
	StationName     string `json:"stationName"`               // 电站名称，必填
	UserID          *Int64 `json:"userId,omitempty"`          // 业主Id，默认操作账号
	Mobile          Str    `json:"mobile,omitempty"`          // 手机号，重复字段仅保留此一个
	Capacity        string `json:"capacity"`                  // 装机容量，必填（单位kWp）
	Latitude        string `json:"latitude,omitempty"`        // 纬度
	Longitude       string `json:"longitude,omitempty"`       // 经度
	Dip             *Num   `json:"dip,omitempty"`             // 倾角
	Azimuth         *Num   `json:"azimuth,omitempty"`         // 方位角
	Money           string `json:"money"`                     // 货币种类，必填
	Addr            string `json:"addr"`                      // 电站详细地址，必填
	GdAreaCode      string `json:"gdAreaCode,omitempty"`      // 高德地址代码
	CountryStr      string `json:"countryStr,omitempty"`      // 国家名称
	RegionStr       string `json:"regionStr,omitempty"`       // 区域名称
	CityStr         string `json:"cityStr,omitempty"`         // 城市名称
	Price           Num    `json:"price"`                     // 每度电收益，必填
	Offset          *Num   `json:"offset,omitempty"`          // 时区偏移量
	Module          *int   `json:"module,omitempty"`          // 组件数量
	InstallerEmail  string `json:"installerEmail,omitempty"`  // 安装商邮箱
	InstallerMobile Str    `json:"installerMobile,omitempty"` // 安装商电话
	NMICode         string `json:"nmiCode,omitempty"`         // nmi码，唯一不可重复
}

// AddStationResult 新增电站绑定逆变器SN结果
type AddStationResult struct {
	rawHolder
	StationID Int64 `json:"stationId"` // 电站id，字段表字段名
	ID        Int64 `json:"id"`        // 电站id，示例返回字段名
}

// AddStation 新增电站并绑定逆变器SN
func (c *SolisSDK) AddStation(req AddStationRequest) (*AddStationResult, error) {
	if strings.TrimSpace(req.StationName) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 stationName 不能为空", PathAddStation)
	}
	if strings.TrimSpace(req.Capacity) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 capacity 不能为空", PathAddStation)
	}
	if strings.TrimSpace(req.Money) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 money 不能为空", PathAddStation)
	}
	if strings.TrimSpace(req.Addr) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 addr 不能为空", PathAddStation)
	}

	var out AddStationResult
	if err := c.do(PathAddStation, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// StationUpdateRequest 修改电站信息请求
type StationUpdateRequest struct {
	ID              *Int64 `json:"id,omitempty"`              // 电站id
	StationName     string `json:"stationName"`               // 电站名称，必填
	Mobile          Str    `json:"mobile,omitempty"`          // 用户手机号，不填默认安装商
	Capacity        string `json:"capacity"`                  // 装机容量，必填（单位kWp）
	Latitude        string `json:"latitude,omitempty"`        // 纬度
	Longitude       string `json:"longitude,omitempty"`       // 经度
	Dip             *Num   `json:"dip,omitempty"`             // 倾角
	Azimuth         *Num   `json:"azimuth,omitempty"`         // 方位角
	Money           string `json:"money,omitempty"`           // 货币种类
	Price           Num    `json:"price"`                     // 每度电收益，必填
	Addr            string `json:"addr"`                      // 电站详细地址，必填
	GdAreaCode      string `json:"gdAreaCode,omitempty"`      // 高德地址代码
	Country         *int   `json:"country,omitempty"`         // 国家，有高德地址代码可不填
	Region          *int   `json:"region,omitempty"`          // 区域，有高德地址代码可不填
	City            *int   `json:"city,omitempty"`            // 城市，有高德地址代码可不填
	CountryStr      string `json:"countryStr,omitempty"`      // 国家名称，其他地图使用
	RegionStr       string `json:"regionStr,omitempty"`       // 区域名称，其他地图使用
	CityStr         string `json:"cityStr,omitempty"`         // 城市名称，其他地图使用
	Module          *int   `json:"module,omitempty"`          // 组件数量
	InstallerEmail  string `json:"installerEmail,omitempty"`  // 安装商邮箱
	InstallerMobile Str    `json:"installerMobile,omitempty"` // 安装商电话
	NMICode         string `json:"nmiCode,omitempty"`         // nmi码
	DaylightSwitch  *int   `json:"daylightSwitch,omitempty"`  // 夏令时开关，0关1开
}

// StationUpdate 修改电站信息，成功时 data 为 null
func (c *SolisSDK) StationUpdate(req StationUpdateRequest) error {
	if strings.TrimSpace(req.StationName) == "" {
		return fmt.Errorf("[Ginlong]%s 必填参数 stationName 不能为空", PathStationUpdate)
	}
	if strings.TrimSpace(req.Capacity) == "" {
		return fmt.Errorf("[Ginlong]%s 必填参数 capacity 不能为空", PathStationUpdate)
	}
	if strings.TrimSpace(req.Addr) == "" {
		return fmt.Errorf("[Ginlong]%s 必填参数 addr 不能为空", PathStationUpdate)
	}

	return c.do(PathStationUpdate, req, nil)
}

// AddStationBindCollectorRequest 新增电站绑定新采集器请求
type AddStationBindCollectorRequest struct {
	StationID           *Int64               `json:"stationId,omitempty"`           // 已存在电站传ID绑定SN，其余字段无效
	SN                  string               `json:"sn,omitempty"`                  // 采集器SN，多个用英文逗号分隔
	StationName         string               `json:"stationName"`                   // 电站名称，新建电站必填
	UserID              *Int64               `json:"userId,omitempty"`              // 业主Id，不填默认安装商
	Capacity            string               `json:"capacity"`                      // 装机容量，新建电站必填（kWp）
	PicName             string               `json:"picName,omitempty"`             // 图片
	Latitude            string               `json:"latitude,omitempty"`            // 纬度
	Longitude           string               `json:"longitude,omitempty"`           // 经度
	Dip                 *Num                 `json:"dip,omitempty"`                 // 倾角
	Azimuth             *Num                 `json:"azimuth,omitempty"`             // 方位角
	Money               string               `json:"money,omitempty"`               // 货币种类
	Addr                string               `json:"addr,omitempty"`                // 电站详细地址
	GdAreaCode          string               `json:"gdAreaCode,omitempty"`          // 高德地址代码
	Country             *int                 `json:"country,omitempty"`             // 国家，有高德地址代码可不填
	Region              *int                 `json:"region,omitempty"`              // 区域，有高德地址代码可不填
	City                *int                 `json:"city,omitempty"`                // 城市，有高德地址代码可不填
	CountryStr          string               `json:"countryStr,omitempty"`          // 国家名称
	RegionStr           string               `json:"regionStr,omitempty"`           // 区域名称
	CityStr             string               `json:"cityStr,omitempty"`             // 城市
	Price               *Num                 `json:"price,omitempty"`               // 每度电收益
	Offset              *Num                 `json:"offset,omitempty"`              // 时区偏移量
	Type                *StationType         `json:"type,omitempty"`                // 电站类型
	Contribution        *int                 `json:"contribution,omitempty"`        // 出资方式
	SynchronizationType *SynchronizationType `json:"synchronizationType,omitempty"` // 并网类型
	InstallTime         string               `json:"installTime,omitempty"`         // 安装时间
	Module              *int                 `json:"module,omitempty"`              // 组件数量
	Mobile              Str                  `json:"mobile,omitempty"`              // 电站联系人
	InstallerEmail      string               `json:"installerEmail,omitempty"`      // 安装商邮箱
	InstallerMobile     Str                  `json:"installerMobile,omitempty"`     // 安装商电话
	NMICode             string               `json:"nmiCode,omitempty"`             // nmi码
	DaylightSwitch      *int                 `json:"daylightSwitch,omitempty"`      // 夏令时开关，0关1开
}

// AddStationBindCollectorResult 新增电站绑定采集器结果，data 为新增电站ID
type AddStationBindCollectorResult struct {
	rawHolder
	StationID Int64 `json:"-"` // data 标量本身，新增电站ID
}

// UnmarshalJSON data 为标量ID，直接解入 StationID
func (r *AddStationBindCollectorResult) UnmarshalJSON(b []byte) error {
	return r.StationID.UnmarshalJSON(b)
}

// AddStationBindCollector 新增电站绑定新采集器，返回新增电站ID
func (c *SolisSDK) AddStationBindCollector(req AddStationBindCollectorRequest) (*AddStationBindCollectorResult, error) {
	// 传 stationId 时为绑定已存在电站，文档注明其余字段均无效
	if req.StationID == nil || req.StationID.Int() == 0 {
		if strings.TrimSpace(req.StationName) == "" {
			return nil, fmt.Errorf("[Ginlong]%s 必填参数 stationName 不能为空", PathAddStationBindCollector)
		}
		if strings.TrimSpace(req.Capacity) == "" {
			return nil, fmt.Errorf("[Ginlong]%s 必填参数 capacity 不能为空", PathAddStationBindCollector)
		}
	}

	var out AddStationBindCollectorResult
	if err := c.do(PathAddStationBindCollector, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DelCollectorRequest 电站解绑采集器请求
type DelCollectorRequest struct {
	SN           string `json:"sn,omitempty"` // 采集器SN
	DeleteInvert int    `json:"deleteInvert"` // 1连带删除所有逆变器，0不删除
}

// DelCollector 电站解绑采集器，成功时 data 为 null
func (c *SolisSDK) DelCollector(req DelCollectorRequest) error {
	return c.do(PathDelCollector, req, nil)
}

// AddDeviceRequest 电站绑定逆变器请求
type AddDeviceRequest struct {
	ID      *Int64 `json:"id,omitempty"`      // 电站id，和nmiCode二选一
	SN      string `json:"sn"`                // 逆变器SN，多个用英文逗号分隔
	NMICode string `json:"nmiCode,omitempty"` // nmi码，和id二选一
}

// AddDeviceItem 绑定结果项，字段未列出，可从 Raw 读取
type AddDeviceItem struct {
	rawHolder
}

// AddDevice 电站绑定逆变器
func (c *SolisSDK) AddDevice(req AddDeviceRequest) ([]AddDeviceItem, error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 sn 不能为空", PathAddDevice)
	}

	var out []AddDeviceItem
	if err := c.do(PathAddDevice, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
