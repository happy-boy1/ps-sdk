package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ps-sdk/model"
)

// db 服务层使用的数据库句柄，由 Init 注入
var db *gorm.DB

// ErrDatabaseNotReady 数据库未初始化
var ErrDatabaseNotReady = errors.New("数据库未初始化，请先调用 service.Init")

// Init 注入数据库句柄，自动补齐表结构并写入平台、设备类型基础数据，可重复调用
func Init(database *gorm.DB) error {
	if database == nil {
		return ErrDatabaseNotReady
	}
	if err := EnsureSchema(database); err != nil {
		return err
	}
	db = database
	return Seed()
}

// DB 返回当前数据库句柄
func DB() *gorm.DB { return db }

// Seed 幂等写入平台与设备类型基础数据，新增平台或类型时重新调用即可
func Seed() error {
	if db == nil {
		return ErrDatabaseNotReady
	}
	if err := upsertPlatforms(PlatformInfos()); err != nil {
		return err
	}
	return upsertDeviceTypes(DeviceTypeModels())
}

// upsertPlatforms 按主键写入平台，已存在则更新名称
func upsertPlatforms(platforms []model.PlatformInfo) error {
	const stmt = "INSERT INTO platform_info (platform_id, platform_code, platform_name_en, platform_name_zh)" +
		" VALUES (?, ?, ?, ?)" +
		" ON DUPLICATE KEY UPDATE platform_code = VALUES(platform_code)," +
		" platform_name_en = VALUES(platform_name_en), platform_name_zh = VALUES(platform_name_zh)"

	return db.Transaction(func(tx *gorm.DB) error {
		for _, p := range platforms {
			if err := tx.Exec(stmt, p.PlatformID, p.PlatformCode, p.PlatformNameEn, p.PlatformNameZh).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// upsertDeviceTypes 按主键写入设备类型，已存在则更新名称
func upsertDeviceTypes(types []model.DeviceType) error {
	if len(types) == 0 {
		return nil
	}

	const stmt = "INSERT INTO device_type (device_type_id, type_code, type_name, type_short_name)" +
		" VALUES (?, ?, ?, ?)" +
		" ON DUPLICATE KEY UPDATE type_code = VALUES(type_code)," +
		" type_name = VALUES(type_name), type_short_name = VALUES(type_short_name)"

	return db.Transaction(func(tx *gorm.DB) error {
		for _, t := range types {
			if err := tx.Exec(stmt, t.DeviceTypeID, t.TypeCode, t.TypeName, t.TypeShortName).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpsertStation 按「平台 + 平台原始电站 ID」写入或更新电站，返回落库后的记录
func UpsertStation(station *model.PowerStation) (*model.PowerStation, error) {
	if db == nil {
		return nil, ErrDatabaseNotReady
	}
	if station == nil || strings.TrimSpace(station.StationIDOrigin) == "" {
		return nil, errors.New("电站缺少平台原始 ID")
	}

	// StationID 为自增主键，交给数据库分配
	station.StationID = 0
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "platform_id"}, {Name: "station_id_origin"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"station_name", "station_short_name", "capacity_kwp",
			"province", "city", "address", "longitude", "latitude",
			"grid_connected_at", "status", "updated_at",
		}),
	}).Create(station).Error
	if err != nil {
		return nil, err
	}

	var saved model.PowerStation
	if err := db.Where("platform_id = ? AND station_id_origin = ?",
		station.PlatformID, station.StationIDOrigin).First(&saved).Error; err != nil {
		return nil, err
	}
	return &saved, nil
}

// UpsertDevice 按「平台 + 平台原始设备 ID」写入或更新设备
func UpsertDevice(device *model.PowerDevice) error {
	if db == nil {
		return ErrDatabaseNotReady
	}
	if device == nil || strings.TrimSpace(device.DeviceIDOrigin) == "" {
		return errors.New("设备缺少平台原始 ID")
	}

	device.DeviceID = 0
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "platform_id"}, {Name: "device_id_origin"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"station_id", "device_type_id", "device_sn", "device_type_origin",
			"device_name", "device_alias", "device_model", "brand",
			"rated_power_kw", "status", "installed_at", "updated_at",
		}),
	}).Create(device).Error
}

// FindStationByOrigin 按平台原始电站 ID 查库内电站
func FindStationByOrigin(platformID int16, origin string) (*model.PowerStation, error) {
	if db == nil {
		return nil, ErrDatabaseNotReady
	}
	var station model.PowerStation
	if err := db.Where("platform_id = ? AND station_id_origin = ?", platformID, origin).
		First(&station).Error; err != nil {
		return nil, err
	}
	return &station, nil
}

// ListStations 查询指定平台已落库的电站，platformID 为 0 时返回全部平台
func ListStations(platformID int16) ([]model.PowerStation, error) {
	if db == nil {
		return nil, ErrDatabaseNotReady
	}
	query := db.Order("station_id ASC")
	if platformID > 0 {
		query = query.Where("platform_id = ?", platformID)
	}
	var stations []model.PowerStation
	if err := query.Find(&stations).Error; err != nil {
		return nil, err
	}
	return stations, nil
}

// ListDevices 查询已落库的设备，stationID 为 0 时返回全部电站
func ListDevices(platformID int16, stationID uint64) ([]model.PowerDevice, error) {
	if db == nil {
		return nil, ErrDatabaseNotReady
	}
	query := db.Order("device_id ASC")
	if platformID > 0 {
		query = query.Where("platform_id = ?", platformID)
	}
	if stationID > 0 {
		query = query.Where("station_id = ?", stationID)
	}
	var devices []model.PowerDevice
	if err := query.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// resolveStationID 把内存中的电站记录映射为数据库主键。
// 电站刚写入时 StationID 已由 UpsertStation 回填，未回填的记录按业务键补查一次。
func resolveStationID(platformID int16, station *model.PowerStation) (uint64, error) {
	if station.StationID > 0 {
		return station.StationID, nil
	}
	saved, err := FindStationByOrigin(platformID, station.StationIDOrigin)
	if err != nil {
		return 0, err
	}
	station.StationID = saved.StationID
	return saved.StationID, nil
}

// findStation 按主键读取电站，供适配器还原平台原始电站 ID
func findStation(platformID int16, stationID uint64) (*model.PowerStation, error) {
	if db == nil {
		return nil, ErrDatabaseNotReady
	}
	var station model.PowerStation
	if err := db.Where("station_id = ? AND platform_id = ?", stationID, platformID).
		First(&station).Error; err != nil {
		return nil, err
	}
	return &station, nil
}
