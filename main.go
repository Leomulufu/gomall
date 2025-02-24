package main

import (
	"log"
	"order_service/dao"
	"order_service/model"
	"order_service/service"
	"order_service/task"
	"time"
)

func testOrderService(s *service.OrderService) {
	// 测试创建订单
	items := []model.OrderItem{
		{
			ProductID: 1,
			Quantity:  2,
			Cost:      99.9,
		},
	}

	order, err := s.PlaceOrder(1001, items)
	if err != nil {
		log.Printf("Failed to place order: %v", err)
		return
	}
	log.Printf("Order created successfully: %v", order.OrderID)

	// 测试查询订单
	orders, err := s.ListOrders(1001)
	if err != nil {
		log.Printf("Failed to list orders: %v", err)
		return
	}
	log.Printf("Found %d orders", len(orders))

	// 测试标记订单支付
	if err := s.MarkOrderPaid(order.OrderID); err != nil {
		log.Printf("Failed to mark order as paid: %v", err)
		return
	}
	log.Printf("Order %d marked as paid", order.OrderID)
}

func main() {
	// 初始化数据库连接
	if err := dao.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connected successfully!")

	// 自动迁移数据库表
	if err := dao.DB.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed!")

	// 启动定时任务
	task.StartOrderCancelTask()
	log.Println("Order cancel task started!")

	// 测试订单服务
	orderService := &service.OrderService{}
	testOrderService(orderService)

	// 等待一段时间观察定时任务
	time.Sleep(time.Minute)
	log.Println("Test completed!")

	// TODO: 启动RPC服务器

	// 保持程序运行
	select {}
}
