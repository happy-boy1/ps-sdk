
CREATE DATABASE IF NOT EXISTS rtm DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;

USE rtm;

-- ============================================================
-- 一、平台与主数据
-- ============================================================
CREATE TABLE IF NOT EXISTS platform_info (
	platform_id SMALLINT NOT NULL COMMENT '平台ID',
	platform_code VARCHAR(64) NOT NULL COMMENT '平台英文标识(程序内使用的小写短码,如xiaomai/ginlong/sungrow/huawei)',
	platform_name_en VARCHAR(256) NOT NULL COMMENT '平台英文名称',
	platform_name_zh VARCHAR(256) NOT NULL COMMENT '平台中文名称',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (platform_id),
	UNIQUE KEY uk_platform_code (platform_code),
	UNIQUE KEY uk_platform_name_en (platform_name_en),
	UNIQUE KEY uk_platform_name_zh (platform_name_zh)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '平台信息表';

-- 各平台开放接口的访问凭据（app_key/app_secret 建议应用层加密后存储；
-- 一个平台对应一组主账号凭据，若多账号可加 owner 字段扩展）
CREATE TABLE IF NOT EXISTS platform_auth (
	auth_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	platform_id SMALLINT NOT NULL COMMENT '平台ID',
	api_url VARCHAR(256) NOT NULL DEFAULT '' COMMENT '开放接口地址',
	app_key VARCHAR(256) COMMENT '接口Key',
	app_secret VARCHAR(512) COMMENT '接口Secret',
	enabled TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用:0停用,1启用',
	remark VARCHAR(256) COMMENT '备注(账号归属等)',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (auth_id),
	UNIQUE KEY uk_platform (platform_id),
	CONSTRAINT fk_pa_platform_id FOREIGN KEY (platform_id) REFERENCES platform_info (platform_id) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '平台采集凭据表';

-- 电站主数据
CREATE TABLE IF NOT EXISTS power_station (
	station_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '电站内部ID(下游表统一使用)',
	platform_id SMALLINT NOT NULL COMMENT '所属平台ID',
	station_id_origin VARCHAR(64) NOT NULL COMMENT '平台原始电站ID(幂等键)',
	station_name VARCHAR(256) NOT NULL DEFAULT '' COMMENT '电站名称(平台原始值)',
	station_short_name VARCHAR(128) COMMENT '电站简称(本系统展示用)',
	capacity_kwp DECIMAL(12, 3) COMMENT '装机容量(kWp)',
	province VARCHAR(64) COMMENT '省份',
	city VARCHAR(64) COMMENT '城市',
	address VARCHAR(256) COMMENT '详细地址',
	longitude DECIMAL(10, 6) COMMENT '经度',
	latitude DECIMAL(9, 6) COMMENT '纬度',
	grid_connected_at DATE COMMENT '并网日期',
	status TINYINT NOT NULL DEFAULT 1 COMMENT '电站状态:0停运,1运行,2在建',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (station_id),
	UNIQUE KEY uk_platform_station (platform_id, station_id_origin),
	KEY idx_platform (platform_id),
	KEY idx_station_name (station_name),
	CONSTRAINT fk_ps_platform_id FOREIGN KEY (platform_id) REFERENCES platform_info (platform_id) ON UPDATE CASCADE ON DELETE RESTRICT,
	CONSTRAINT ck_station_capacity CHECK (capacity_kwp > 0)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '电站信息表';

-- 设备主数据
CREATE TABLE IF NOT EXISTS power_device (
	device_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '设备内部ID(下游表统一使用)',
	platform_id SMALLINT NOT NULL COMMENT '所属平台ID(冗余自电站,采集层按此定位接口)',
	station_id BIGINT UNSIGNED NOT NULL COMMENT '从属电站(内部ID)',
	device_type_id SMALLINT NOT NULL COMMENT '设备类型ID',
	device_id_origin VARCHAR(64) NOT NULL COMMENT '平台原始设备ID/通信地址(幂等键)',
	device_sn VARCHAR(64) COMMENT '设备SN号/序列号',
	device_type_origin VARCHAR(32) COMMENT '平台原始设备类型编码(部分平台的实时接口要求回传,如阳光云 device_type=8/11/23)',
	device_name VARCHAR(256) COMMENT '设备名称(平台原始值)',
	device_alias VARCHAR(256) COMMENT '设备别称',
	device_model VARCHAR(128) COMMENT '设备型号',
	brand VARCHAR(64) COMMENT '品牌',
	rated_power_kw DECIMAL(12, 3) COMMENT '额定功率(kW),逆变器/储能等填写',
	installed_at DATE COMMENT '投运/安装日期',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (device_id),
	UNIQUE KEY uk_platform_device (platform_id, device_id_origin),
	KEY idx_station_id (station_id),
	KEY idx_device_type_id (device_type_id),
	KEY idx_device_sn (device_sn),
	CONSTRAINT fk_pd_station_id FOREIGN KEY (station_id) REFERENCES power_station (station_id) ON UPDATE CASCADE ON DELETE RESTRICT,
	CONSTRAINT fk_pd_device_type_id FOREIGN KEY (device_type_id) REFERENCES device_type (device_type_id) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '电站设备表';

-- 设备类型：全局字典（跨平台归一化，采集层负责把各平台类型映射到 type_code）
CREATE TABLE IF NOT EXISTS device_type (
	device_type_id SMALLINT NOT NULL COMMENT '设备类型ID',
	type_code VARCHAR(32) NOT NULL COMMENT '类型编码(程序内使用,如inverter/data_logger)',
	type_name VARCHAR(64) NOT NULL COMMENT '类型名称',
	type_short_name VARCHAR(32) NOT NULL COMMENT '类型简称',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (device_type_id),
	UNIQUE KEY uk_type_code (type_code)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '设备类型字典表(全局,不区分平台)';