package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 连接参数默认值，可用环境变量覆盖（PS_DB_DSN 优先于分散参数）
const (
	envDSN      = "PS_DB_DSN"
	envHost     = "PS_DB_HOST"
	envPort     = "PS_DB_PORT"
	envUser     = "PS_DB_USER"
	envPassword = "PS_DB_PASSWORD"
	envName     = "PS_DB_NAME"

	defaultHost     = "127.0.0.1"
	defaultPort     = "3306"
	defaultUser     = "root"
	defaultName     = "rtm"
	defaultMaxIdle  = 12
	defaultMaxOpen  = 100
	defaultConnLife = time.Hour
)

var (
	// DB 全局数据库句柄，Init 成功后可用
	DB *gorm.DB

	mu sync.Mutex
)

// ErrNotInitialized 尚未初始化数据库
var ErrNotInitialized = errors.New("数据库尚未初始化")

// Init 初始化数据库连接，重复调用直接返回已有连接
func Init() error {
	mu.Lock()
	defer mu.Unlock()

	if DB != nil {
		return nil
	}

	database, err := open()
	if err != nil {
		return err
	}
	DB = database
	return nil
}

func open() (*gorm.DB, error) {
	log.Println("[DB] 初始化数据库连接")

	database, err := gorm.Open(mysql.Open(dsn()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         gormLogger(),
		// 外键约束由 DDL 维护，交给 GORM 迁移会与其他表互相引用而失败
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("获取 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(envInt("PS_DB_MAX_IDLE", defaultMaxIdle))
	sqlDB.SetMaxOpenConns(envInt("PS_DB_MAX_OPEN", defaultMaxOpen))
	sqlDB.SetConnMaxLifetime(defaultConnLife)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连通性检查失败: %w", err)
	}

	log.Printf("[DB] 连接完成: %s/%s", envOr(envHost, defaultHost), envOr(envName, defaultName))
	return database, nil
}

// dsn 组装 MySQL DSN，PS_DB_DSN 存在时直接使用
func dsn() string {
	if v := envOr(envDSN, ""); v != "" {
		return v
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		envOr(envUser, defaultUser),
		envOr(envPassword, "wsljj"),
		envOr(envHost, defaultHost),
		envOr(envPort, defaultPort),
		envOr(envName, defaultName),
	)
}

func gormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

// Close 关闭数据库连接
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if DB == nil {
		return ErrNotInitialized
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取 sql.DB 失败: %w", err)
	}

	DB = nil
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}
	log.Println("[DB] 连接已关闭")
	return nil
}

// envOr 读取环境变量，为空时返回默认值
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// envInt 读取整型环境变量，非法值返回默认值
func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
