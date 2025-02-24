package service

import (
	"order_service/dao"
	"order_service/model"
	"testing"
)

func TestOrderService(t *testing.T) {
	// 初始化数据库连接
	if err := dao.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	s := &OrderService{}

	// 测试创建订单
	t.Run("PlaceOrder", func(t *testing.T) {
		items := []model.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Cost:      99.9,
			},
		}
		order, err := s.PlaceOrder(1001, items)
		if err != nil {
			t.Errorf("Failed to place order: %v", err)
		}
		if order.OrderID == 0 {
			t.Error("Expected order ID to be non-zero")
		}
	})

	// 测试查询订单
	t.Run("ListOrders", func(t *testing.T) {
		orders, err := s.ListOrders(1001)
		if err != nil {
			t.Errorf("Failed to list orders: %v", err)
		}
		if len(orders) == 0 {
			t.Error("Expected to find orders")
		}
	})

	// 测试标记订单支付
	t.Run("MarkOrderPaid", func(t *testing.T) {
		// 先创建一个订单
		items := []model.OrderItem{{ProductID: 1, Quantity: 1, Cost: 10.0}}
		order, _ := s.PlaceOrder(1001, items)

		err := s.MarkOrderPaid(order.OrderID)
		if err != nil {
			t.Errorf("Failed to mark order as paid: %v", err)
		}
	})
}
