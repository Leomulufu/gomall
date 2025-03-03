package dao

import (
	"log"
	"order_service/model"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	log.Println("开始数据库迁移...")

	// 删除现有的表（如果存在）
	err := DB.Migrator().DropTable(&model.Order{}, &model.OrderItem{})
	if err != nil {
		log.Printf("删除表失败: %v", err)
		return err
	}

	// 创建表
	err = DB.AutoMigrate(&model.Order{}, &model.OrderItem{})
	if err != nil {
		log.Printf("创建表失败: %v", err)
		return err
	}

	log.Println("数据库迁移完成!")
	return nil
}
