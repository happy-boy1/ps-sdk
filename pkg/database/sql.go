package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() {
	log.Println("初始化数据库连接...")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"wsljj",
		"127.0.0.1",
		"3306",
		"rtm",
	)

	gormlog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
			Logger:         gormlog,
		},
	)

	if err != nil {
		panic(err)
	}

	sqldb, err := db.DB()

	if err != nil {
		panic(err)
	}

	sqldb.SetMaxIdleConns(12)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("数据库连接完成")
}

func Close() {
	db, err := DB.DB()
	if err != nil {
		log.Printf("[Error] 获取sql.DB错误: %v\n", err)
		return
	}
	db.Close()
	log.Println("关闭数据库连接成功")
}
