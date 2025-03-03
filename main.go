package main

import (
	"context"
	"log"
	"order_service/dao"
	"order_service/grpc"
	"order_service/model"
	"order_service/service"
	"order_service/task"
)

// 测试订单服务
func testOrderService(ctx context.Context, s *service.OrderService) {
	log.Println("开始测试订单服务...")

	// 创建测试地址
	address := &model.Address{
		StreetAddress: "123 抖音路",
		City:          "北京",
		State:         "北京",
		Country:       "中国",
		ZipCode:       100000,
	}

	// 测试创建订单
	items := []model.OrderItem{
		{
			ProductID:   1,
			ProductName: "抖音限定T恤",
			Quantity:    2,
			UnitPrice:   49.9,
			Cost:        99.8,
		},
		{
			ProductID:   2,
			ProductName: "抖音定制手机壳",
			Quantity:    1,
			UnitPrice:   29.9,
			Cost:        29.9,
		},
	}

	order, err := s.PlaceOrder(ctx, 1001, "CNY", address, "test@example.com", items)
	if err != nil {
		log.Printf("创建订单失败: %v", err)
		return
	}
	log.Printf("订单创建成功: %v", order.OrderID)

	// 测试查询订单
	orders, err := s.ListOrders(ctx, 1001)
	if err != nil {
		log.Printf("查询订单失败: %v", err)
		return
	}
	log.Printf("找到 %d 个订单", len(orders))

	// 测试标记订单支付
	if err := s.MarkOrderPaid(ctx, order.OrderID, "支付宝"); err != nil {
		log.Printf("标记订单已支付失败: %v", err)
		return
	}
	log.Printf("订单 %s 已标记为已支付", order.OrderID)

	// 测试标记订单发货
	if err := s.ShipOrder(ctx, order.OrderID); err != nil {
		log.Printf("标记订单已发货失败: %v", err)
		return
	}
	log.Printf("订单 %s 已标记为已发货", order.OrderID)

	// 测试标记订单送达
	if err := s.DeliverOrder(ctx, order.OrderID); err != nil {
		log.Printf("标记订单已送达失败: %v", err)
		return
	}
	log.Printf("订单 %s 已标记为已送达", order.OrderID)

	log.Println("订单服务测试完成!")
}

func main() {
	log.Println("启动抖音电商订单服务...")

	// 初始化数据库连接
	if err := dao.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("数据库连接成功!")

	// 执行数据库迁移
	if err := dao.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 启动定时任务
	task.StartOrderTasks()
	log.Println("订单相关定时任务已启动!")

	// 测试创建订单
	log.Println("开始测试订单服务...")
	orderService := &service.OrderService{}
	testOrder := &model.Order{
		UserID:        1001,
		UserCurrency:  "CNY",
		TotalAmount:   129.70,
		ItemCount:     2,
		StreetAddress: "123 抖音路",
		City:          "北京",
		State:         "北京",
		Country:       "中国",
		ZipCode:       100000,
		Email:         "test@example.com",
		OrderStatus:   model.OrderStatusPending,
	}

	items := []model.OrderItem{
		{
			ProductID:   1,
			ProductName: "测试商品1",
			Quantity:    1,
			UnitPrice:   59.90,
			Cost:        59.90,
		},
		{
			ProductID:   2,
			ProductName: "测试商品2",
			Quantity:    1,
			UnitPrice:   69.80,
			Cost:        69.80,
		},
	}

	_, err := orderService.PlaceOrder(nil, testOrder.UserID, testOrder.UserCurrency, &model.Address{
		StreetAddress: testOrder.StreetAddress,
		City:          testOrder.City,
		State:         testOrder.State,
		Country:       testOrder.Country,
		ZipCode:       testOrder.ZipCode,
	}, testOrder.Email, items)

	if err != nil {
		log.Printf("创建订单失败: %v", err)
	}

	// 启动 gRPC 服务器
	if err := grpc.StartServer(); err != nil {
		log.Fatalf("gRPC服务器启动失败: %v", err)
	}

	log.Println("服务已启动. 按 Ctrl+C 退出.")
	select {}
}
