package service

import (
	"fmt"
	"order_service/dao"
	"order_service/model"
	"time"

	"gorm.io/gorm"
)

type OrderService struct{}

// PlaceOrder creates a new order
func (s *OrderService) PlaceOrder(userID int64, items []model.OrderItem) (*model.Order, error) {
	order := &model.Order{
		UserID:      userID,
		OrderStatus: "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Use transaction to ensure data consistency
	return order, dao.DB.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Create order items
		for i := range items {
			items[i].OrderID = order.OrderID
			if err := tx.Create(&items[i]).Error; err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}
		return nil
	})
}

// ListOrders retrieves orders for a given user ID
func (s *OrderService) ListOrders(userID int64) ([]model.Order, error) {
	var orders []model.Order
	result := dao.DB.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list orders for user %d: %w", userID, result.Error)
	}
	return orders, nil
}

// MarkOrderPaid updates the order status to paid
func (s *OrderService) MarkOrderPaid(orderID int64) error {
	result := dao.DB.Model(&model.Order{}).
		Where("order_id = ?", orderID).
		Update("order_status", "paid")
	if result.Error != nil {
		return fmt.Errorf("failed to mark order %d as paid: %w", orderID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order %d not found", orderID)
	}
	return nil
}
