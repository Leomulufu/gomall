package service

import (
	"order_service/dao"
	"order_service/model"
	"time"

	"gorm.io/gorm"
)

type OrderService struct{}

// 创建订单
func (s *OrderService) PlaceOrder(userID int64, items []model.OrderItem) (*model.Order, error) {
	order := &model.Order{
		UserID:      userID,
		OrderStatus: "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 使用事务确保数据一致性
	return order, dao.DB.Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 创建订单项
		for i := range items {
			items[i].OrderID = order.OrderID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// 查询订单列表
func (s *OrderService) ListOrders(userID int64) ([]model.Order, error) {
	var orders []model.Order
	err := dao.DB.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

// 标记订单已支付
func (s *OrderService) MarkOrderPaid(orderID int64) error {
	return dao.DB.Model(&model.Order{}).
		Where("order_id = ?", orderID).
		Update("order_status", "paid").Error
}
