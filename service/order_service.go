package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"order_service/dao"
	"order_service/model"
	"time"

	"gorm.io/gorm"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidInput  = errors.New("invalid input parameters")
	rng              = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// 生成订单ID
func generateOrderID(userID uint32) string {
	timestamp := time.Now().Format("060102150405")
	randomStr := fmt.Sprintf("%03d", rng.Intn(1000))
	return fmt.Sprintf("%s%d%s", timestamp, userID, randomStr)
}

type OrderService struct{}

// PlaceOrder 创建新订单
func (s *OrderService) PlaceOrder(ctx context.Context, userID uint32, userCurrency string, address *model.Address, email string, items []model.OrderItem) (*model.Order, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("%w: invalid user ID", ErrInvalidInput)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("%w: order must contain at least one item", ErrInvalidInput)
	}

	// 生成订单ID
	orderID := generateOrderID(userID)

	// 计算订单总金额
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Cost
	}

	order := &model.Order{
		OrderID:       orderID,
		UserID:        userID,
		UserCurrency:  userCurrency,
		TotalAmount:   totalAmount,
		ItemCount:     len(items),
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		Email:         email,
		OrderStatus:   model.OrderStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PaidAt:        time.Time{},
		ShippedAt:     time.Time{},
		DeliveredAt:   time.Time{},
	}

	// 使用事务确保数据一致性
	err := dao.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// 创建订单项
		for i := range items {
			items[i].OrderID = orderID
			if err := tx.Create(&items[i]).Error; err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error in PlaceOrder: %v", err)
		return nil, err
	}

	return order, nil
}

// GetOrder 获取单个订单详情
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*model.Order, []model.OrderItem, error) {
	if orderID == "" {
		return nil, nil, fmt.Errorf("%w: invalid order ID", ErrInvalidInput)
	}

	var order model.Order
	if err := dao.DB.WithContext(ctx).Where("order_id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("%w: order %s not found", ErrOrderNotFound, orderID)
		}
		log.Printf("Error in GetOrder: %v", err)
		return nil, nil, fmt.Errorf("failed to get order %s: %w", orderID, err)
	}

	var items []model.OrderItem
	if err := dao.DB.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		log.Printf("Error in GetOrder items: %v", err)
		return nil, nil, fmt.Errorf("failed to get order items for %s: %w", orderID, err)
	}

	return &order, items, nil
}

// ListOrders 获取用户的所有订单
func (s *OrderService) ListOrders(ctx context.Context, userID uint32) ([]model.Order, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("%w: invalid user ID", ErrInvalidInput)
	}

	var orders []model.Order
	result := dao.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&orders)
	if result.Error != nil {
		log.Printf("Error in ListOrders: %v", result.Error)
		return nil, fmt.Errorf("failed to list orders for user %d: %w", userID, result.Error)
	}
	return orders, nil
}

// MarkOrderPaid 将订单标记为已支付
func (s *OrderService) MarkOrderPaid(ctx context.Context, orderID string, paymentMethod string) error {
	if orderID == "" {
		return fmt.Errorf("%w: invalid order ID", ErrInvalidInput)
	}

	now := time.Now()
	result := dao.DB.WithContext(ctx).Model(&model.Order{}).
		Where("order_id = ? AND order_status = ?", orderID, model.OrderStatusPending).
		Updates(map[string]interface{}{
			"order_status":   model.OrderStatusPaid,
			"payment_method": paymentMethod,
			"paid_at":        now,
			"updated_at":     now,
		})

	if result.Error != nil {
		log.Printf("Error in MarkOrderPaid: %v", result.Error)
		return fmt.Errorf("failed to mark order %s as paid: %w", orderID, result.Error)
	}

	if result.RowsAffected == 0 {
		// 检查订单是否存在
		var count int64
		dao.DB.WithContext(ctx).Model(&model.Order{}).Where("order_id = ?", orderID).Count(&count)
		if count == 0 {
			return fmt.Errorf("%w: order %s not found", ErrOrderNotFound, orderID)
		}
		// 订单存在但状态不是pending
		return fmt.Errorf("order %s cannot be marked as paid due to its current status", orderID)
	}

	return nil
}

// ShipOrder 将订单标记为已发货
func (s *OrderService) ShipOrder(ctx context.Context, orderID string) error {
	if orderID == "" {
		return fmt.Errorf("%w: invalid order ID", ErrInvalidInput)
	}

	now := time.Now()
	result := dao.DB.WithContext(ctx).Model(&model.Order{}).
		Where("order_id = ? AND order_status = ?", orderID, model.OrderStatusPaid).
		Updates(map[string]interface{}{
			"order_status": model.OrderStatusShipped,
			"shipped_at":   now,
			"updated_at":   now,
		})

	if result.Error != nil {
		log.Printf("Error in ShipOrder: %v", result.Error)
		return fmt.Errorf("failed to mark order %s as shipped: %w", orderID, result.Error)
	}

	if result.RowsAffected == 0 {
		// 检查订单是否存在
		var count int64
		dao.DB.WithContext(ctx).Model(&model.Order{}).Where("order_id = ?", orderID).Count(&count)
		if count == 0 {
			return fmt.Errorf("%w: order %s not found", ErrOrderNotFound, orderID)
		}
		// 订单存在但状态不是paid
		return fmt.Errorf("order %s cannot be shipped due to its current status", orderID)
	}

	return nil
}

// DeliverOrder 将订单标记为已送达
func (s *OrderService) DeliverOrder(ctx context.Context, orderID string) error {
	if orderID == "" {
		return fmt.Errorf("%w: invalid order ID", ErrInvalidInput)
	}

	now := time.Now()
	result := dao.DB.WithContext(ctx).Model(&model.Order{}).
		Where("order_id = ? AND order_status = ?", orderID, model.OrderStatusShipped).
		Updates(map[string]interface{}{
			"order_status": model.OrderStatusDelivered,
			"delivered_at": now,
			"updated_at":   now,
		})

	if result.Error != nil {
		log.Printf("Error in DeliverOrder: %v", result.Error)
		return fmt.Errorf("failed to mark order %s as delivered: %w", orderID, result.Error)
	}

	if result.RowsAffected == 0 {
		// 检查订单是否存在
		var count int64
		dao.DB.WithContext(ctx).Model(&model.Order{}).Where("order_id = ?", orderID).Count(&count)
		if count == 0 {
			return fmt.Errorf("%w: order %s not found", ErrOrderNotFound, orderID)
		}
		// 订单存在但状态不是shipped
		return fmt.Errorf("order %s cannot be delivered due to its current status", orderID)
	}

	return nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, orderID string) error {
	if orderID == "" {
		return fmt.Errorf("%w: invalid order ID", ErrInvalidInput)
	}

	now := time.Now()
	result := dao.DB.WithContext(ctx).Model(&model.Order{}).
		Where("order_id = ? AND order_status IN ?", orderID, []string{model.OrderStatusPending, model.OrderStatusPaid}).
		Updates(map[string]interface{}{
			"order_status": model.OrderStatusCancelled,
			"updated_at":   now,
		})

	if result.Error != nil {
		log.Printf("Error in CancelOrder: %v", result.Error)
		return fmt.Errorf("failed to cancel order %s: %w", orderID, result.Error)
	}

	if result.RowsAffected == 0 {
		// 检查订单是否存在
		var count int64
		dao.DB.WithContext(ctx).Model(&model.Order{}).Where("order_id = ?", orderID).Count(&count)
		if count == 0 {
			return fmt.Errorf("%w: order %s not found", ErrOrderNotFound, orderID)
		}
		// 订单存在但状态不允许取消
		return fmt.Errorf("order %s cannot be cancelled due to its current status", orderID)
	}

	return nil
}
