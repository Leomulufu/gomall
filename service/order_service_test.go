package service

import (
	"context"
	"order_service/dao"
	"order_service/model"
	"testing"
)

func TestOrderService(t *testing.T) {
	// 初始化数据库连接
	if err := dao.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	dao.DB.AutoMigrate(&model.Order{}, &model.OrderItem{})

	s := &OrderService{}
	ctx := context.Background()

	// 测试创建订单
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

	address := &model.Address{
		StreetAddress: "123 抖音路",
		City:          "北京",
		State:         "北京",
		Country:       "中国",
		ZipCode:       100000,
	}

	// 测试创建订单
	order, err := s.PlaceOrder(ctx, 1001, "CNY", address, "test@example.com", items)
	if err != nil {
		t.Errorf("创建订单失败: %v", err)
		return
	}
	if order.OrderID == "" {
		t.Error("创建的订单ID为空")
	}

	// 测试查询订单列表
	orders, err := s.ListOrders(ctx, 1001)
	if err != nil {
		t.Errorf("查询订单列表失败: %v", err)
		return
	}
	if len(orders) == 0 {
		t.Error("订单列表为空")
	}

	// 测试标记订单已支付
	err = s.MarkOrderPaid(ctx, order.OrderID, "支付宝")
	if err != nil {
		t.Errorf("标记订单已支付失败: %v", err)
		return
	}

	// 测试标记订单已发货
	err = s.ShipOrder(ctx, order.OrderID)
	if err != nil {
		t.Errorf("标记订单已发货失败: %v", err)
		return
	}

	// 测试标记订单已送达
	err = s.DeliverOrder(ctx, order.OrderID)
	if err != nil {
		t.Errorf("标记订单已送达失败: %v", err)
		return
	}

	// 测试创建订单失败的情况
	t.Run("PlaceOrderFailed", func(t *testing.T) {
		// 模拟数据库错误
		dao.DB.Migrator().DropTable(&model.Order{})
		items := []model.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Cost:      99.9,
			},
		}
		_, err := s.PlaceOrder(ctx, 1001, "CNY", address, "test@example.com", items)
		if err == nil {
			t.Error("Expected to fail to place order")
		}
		// 恢复数据库
		dao.InitDB()
		dao.DB.AutoMigrate(&model.Order{}, &model.OrderItem{})
	})

	// 测试查询订单列表为空的情况
	t.Run("ListOrdersEmpty", func(t *testing.T) {
		orders, err := s.ListOrders(ctx, 9999)
		if err != nil {
			t.Errorf("Failed to list orders: %v", err)
		}
		if len(orders) != 0 {
			t.Error("Expected to find no orders")
		}
	})

	// 测试标记订单支付失败的情况
	t.Run("MarkOrderPaidFailed", func(t *testing.T) {
		err := s.MarkOrderPaid(ctx, "9999", "支付宝")
		if err == nil {
			t.Error("Expected to fail to mark order as paid")
		}
	})
}
