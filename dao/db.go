package dao

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gomall?charset=utf8mb4&parseTime=True&loc=Local"
	
	// 配置GORM日志
	gormLogger := logger.Default.LogMode(logger.Info)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}
	
	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间
	
	log.Println("Database connection pool configured")
	
	DB = db
	return nil
}
