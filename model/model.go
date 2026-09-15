package model

import "time"

// 公共时间字段，可嵌入各结构体
type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;autoUpdateTime"`
}

// 时间忽略字段（从 GORM 映射中剔除，保持结构）
type TimestampOnly struct {
	CreatedAt time.Time `gorm:"->;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"->;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type PlatformInfo struct {
	PlatformID     int16  `gorm:"column:platform_id;type:smallint;primaryKey;not null;comment:平台ID"`
	PlatformCode   string `gorm:"column:platform_code;type:varchar(64);uniqueIndex:uk_platform_code;not null;comment:平台英文标识(小写短码)"`
	PlatformNameEn string `gorm:"column:platform_name_en;type:varchar(256);uniqueIndex:uk_platform_name_en;not null;comment:平台英文名称"`
	PlatformNameZh string `gorm:"column:platform_name_zh;type:varchar(256);uniqueIndex:uk_platform_name_zh;not null;comment:平台中文名称"`

	// 关联:一个平台对应一条凭据记录
	Auth *PlatformAuth `gorm:"foreignKey:PlatformID;references:PlatformID"`
	BaseModel
}

func (PlatformInfo) TableName() string { return "platform_info" }

type PlatformAuth struct {
	AuthID     uint32 `gorm:"column:auth_id;type:int unsigned;autoIncrement;primaryKey;comment:凭据ID"`
	PlatformID int16  `gorm:"column:platform_id;type:smallint;uniqueIndex:uk_platform;not null;comment:平台ID"`
	ApiURL     string `gorm:"column:api_url;type:varchar(256);default:'';comment:开放接口地址"`
	AppKey     string `gorm:"column:app_key;type:varchar(256);comment:接口Key"`
	AppSecret  string `gorm:"column:app_secret;type:varchar(512);comment:接口Secret"`
	Enabled    int8   `gorm:"column:enabled;type:tinyint;default:1;comment:是否启用"`
	Remark     string `gorm:"column:remark;type:varchar(256);comment:备注"`
	BaseModel

	// 外键关联
	Platform *PlatformInfo `gorm:"foreignKey:PlatformID;references:PlatformID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (PlatformAuth) TableName() string { return "platform_auth" }

type PowerStation struct {
	StationID        uint64     `gorm:"column:station_id;type:bigint unsigned;autoIncrement;primaryKey;comment:电站内部ID"`
	PlatformID       int16      `gorm:"column:platform_id;type:smallint;uniqueIndex:uk_platform_station,priority:1;index:idx_platform;not null;comment:所属平台ID"`
	StationIDOrigin  string     `gorm:"column:station_id_origin;type:varchar(64);uniqueIndex:uk_platform_station,priority:2;not null;comment:平台原始电站ID"`
	StationName      string     `gorm:"column:station_name;type:varchar(256);not null;default:'';index:idx_station_name;comment:电站名称"`
	StationShortName string     `gorm:"column:station_short_name;type:varchar(128);comment:电站简称"`
	CapacityKwp      *float64   `gorm:"column:capacity_kwp;type:decimal(12,3);comment:装机容量(kWp);check:capacity_kwp > 0"`
	Province         string     `gorm:"column:province;type:varchar(64);comment:省份"`
	City             string     `gorm:"column:city;type:varchar(64);comment:城市"`
	Address          string     `gorm:"column:address;type:varchar(256);comment:详细地址"`
	Longitude        float64    `gorm:"column:longitude;type:decimal(10,6);comment:经度"`
	Latitude         float64    `gorm:"column:latitude;type:decimal(9,6);comment:纬度"`
	GridConnectedAt  *time.Time `gorm:"column:grid_connected_at;type:date;comment:并网日期"`
	Status           int8       `gorm:"column:status;type:tinyint;not null;default:1;comment:电站状态:0停运,1运行,2在建"`
	BaseModel

	Platform *PlatformInfo `gorm:"foreignKey:PlatformID;references:PlatformID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Devices  []PowerDevice `gorm:"foreignKey:StationID;references:StationID"`
}

func (PowerStation) TableName() string { return "power_station" }

type DeviceType struct {
	DeviceTypeID  int16  `gorm:"column:device_type_id;type:smallint;primaryKey;not null;comment:设备类型ID"`
	TypeCode      string `gorm:"column:type_code;type:varchar(32);uniqueIndex:uk_type_code;not null;comment:类型编码"`
	TypeName      string `gorm:"column:type_name;type:varchar(64);not null;comment:类型名称"`
	TypeShortName string `gorm:"column:type_short_name;type:varchar(32);not null;comment:类型简称"`
	BaseModel
}

func (DeviceType) TableName() string { return "device_type" }

type PowerDevice struct {
	DeviceID         uint64     `gorm:"column:device_id;type:bigint unsigned;autoIncrement;primaryKey;comment:设备内部ID"`
	PlatformID       int16      `gorm:"column:platform_id;type:smallint;uniqueIndex:uk_platform_device,priority:1;not null;comment:所属平台ID"`
	StationID        uint64     `gorm:"column:station_id;type:bigint unsigned;index:idx_station_id;not null;comment:从属电站"`
	DeviceTypeID     int16      `gorm:"column:device_type_id;type:smallint;index:idx_device_type_id;not null;comment:设备类型ID"`
	DeviceIDOrigin   string     `gorm:"column:device_id_origin;type:varchar(64);uniqueIndex:uk_platform_device,priority:2;not null;comment:平台原始设备ID"`
	DeviceSN         string     `gorm:"column:device_sn;type:varchar(64);index:idx_device_sn;comment:设备SN号"`
	DeviceTypeOrigin string     `gorm:"column:device_type_origin;type:varchar(32);comment:平台原始设备类型编码"`
	DeviceName       string     `gorm:"column:device_name;type:varchar(256);comment:设备名称"`
	DeviceAlias      string     `gorm:"column:device_alias;type:varchar(256);comment:设备别称"`
	DeviceModel      string     `gorm:"column:device_model;type:varchar(128);comment:设备型号"`
	Brand            string     `gorm:"column:brand;type:varchar(64);comment:品牌"`
	RatedPowerKw     *float64   `gorm:"column:rated_power_kw;type:decimal(12,3);comment:额定功率(kW)"`
	Status           int8       `gorm:"column:status;type:tinyint;not null;default:0;comment:设备状态:0未知,1在线,2离线,3告警"`
	InstalledAt      *time.Time `gorm:"column:installed_at;type:date;comment:投运日期"`
	BaseModel

	// 外键关联
	Station    *PowerStation `gorm:"foreignKey:StationID;references:StationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DeviceType *DeviceType   `gorm:"foreignKey:DeviceTypeID;references:DeviceTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (PowerDevice) TableName() string { return "power_device" }
