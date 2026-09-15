package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"ps-sdk/pkg/config"
)

var (
	// DB 全局数据库句柄，Init 成功后可用
	DB *gorm.DB

	mu sync.Mutex
)

// ErrNotInitialized 尚未初始化数据库
var ErrNotInitialized = errors.New("数据库尚未初始化")

// Init 按配置初始化数据库连接，重复调用直接返回已有连接
func Init(cfg config.Database) error {
	mu.Lock()
	defer mu.Unlock()

	if DB != nil {
		return nil
	}

	database, err := Open(cfg)
	if err != nil {
		return err
	}
	DB = database
	return nil
}

// Open 建立并校验一个数据库连接，不改动全局 DB，便于测试与多数据源场景
func Open(cfg config.Database) (*gorm.DB, error) {
	log.Printf("[DB] 初始化数据库连接: %s/%s", cfg.Host, cfg.Name)

	database, err := gorm.Open(mysql.Open(cfg.DSNOrBuild()), &gorm.Config{
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
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifeDuration())

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连通性检查失败: %w", err)
	}

	log.Printf("[DB] 连接完成: %s/%s", cfg.Host, cfg.Name)
	return database, nil
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
