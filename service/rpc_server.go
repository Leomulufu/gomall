package service

import (
	"context"
	"log"
	"order_service/model"
)

// OrderRPCServer 订单RPC服务器
type OrderRPCServer struct {
	orderService *OrderService
}

// NewOrderRPCServer 创建新的订单RPC服务器
func NewOrderRPCServer() *OrderRPCServer {
	return &OrderRPCServer{
		orderService: &OrderService{},
	}
}

// PlaceOrder RPC方法：创建订单
func (s *OrderRPCServer) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*PlaceOrderResponse, error) {
	// 转换请求参数
	address := &model.Address{
		StreetAddress: req.Address.StreetAddress,
		City:          req.Address.City,
		State:         req.Address.State,
		Country:       req.Address.Country,
		ZipCode:       req.Address.ZipCode,
	}

	var orderItems []model.OrderItem
	for _, item := range req.OrderItems {
		orderItems = append(orderItems, model.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Cost:        item.Cost,
		})
	}

	// 调用服务方法
	order, err := s.orderService.PlaceOrder(ctx, req.UserID, req.UserCurrency, address, req.Email, orderItems)
	if err != nil {
		log.Printf("RPC PlaceOrder error: %v", err)
		return nil, err
	}

	// 构造响应
	return &PlaceOrderResponse{
		OrderID: order.OrderID,
	}, nil
}

// GetOrder RPC方法：获取订单详情
func (s *OrderRPCServer) GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error) {
	// 调用服务方法
	order, items, err := s.orderService.GetOrder(ctx, req.OrderID)
	if err != nil {
		log.Printf("RPC GetOrder error: %v", err)
		return nil, err
	}

	// 转换订单项
	var responseItems []*OrderItemResponse
	for _, item := range items {
		responseItems = append(responseItems, &OrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Cost:        item.Cost,
		})
	}

	// 构造响应
	return &GetOrderResponse{
		Order: &OrderResponse{
			OrderID:       order.OrderID,
			UserID:        order.UserID,
			UserCurrency:  order.UserCurrency,
			TotalAmount:   order.TotalAmount,
			ItemCount:     int32(order.ItemCount),
			OrderStatus:   order.OrderStatus,
			CreatedAt:     order.CreatedAt.Unix(),
			StreetAddress: order.StreetAddress,
			City:          order.City,
			State:         order.State,
			Country:       order.Country,
			ZipCode:       order.ZipCode,
			Email:         order.Email,
			Items:         responseItems,
		},
	}, nil
}

// ListOrders RPC方法：获取用户订单列表
func (s *OrderRPCServer) ListOrders(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error) {
	// 调用服务方法
	orders, err := s.orderService.ListOrders(ctx, req.UserID)
	if err != nil {
		log.Printf("RPC ListOrders error: %v", err)
		return nil, err
	}

	// 转换订单列表
	var responseOrders []*OrderResponse
	for _, order := range orders {
		responseOrders = append(responseOrders, &OrderResponse{
			OrderID:       order.OrderID,
			UserID:        order.UserID,
			UserCurrency:  order.UserCurrency,
			TotalAmount:   order.TotalAmount,
			ItemCount:     int32(order.ItemCount),
			OrderStatus:   order.OrderStatus,
			CreatedAt:     order.CreatedAt.Unix(),
			StreetAddress: order.StreetAddress,
			City:          order.City,
			State:         order.State,
			Country:       order.Country,
			ZipCode:       order.ZipCode,
			Email:         order.Email,
		})
	}

	// 构造响应
	return &ListOrdersResponse{
		Orders: responseOrders,
	}, nil
}

// MarkOrderPaid RPC方法：标记订单为已支付
func (s *OrderRPCServer) MarkOrderPaid(ctx context.Context, req *MarkOrderPaidRequest) (*MarkOrderPaidResponse, error) {
	// 调用服务方法
	err := s.orderService.MarkOrderPaid(ctx, req.OrderID, req.PaymentMethod)
	if err != nil {
		log.Printf("RPC MarkOrderPaid error: %v", err)
		return nil, err
	}

	// 构造响应
	return &MarkOrderPaidResponse{
		Success: true,
	}, nil
}

// ShipOrder RPC方法：标记订单为已发货
func (s *OrderRPCServer) ShipOrder(ctx context.Context, req *ShipOrderRequest) (*ShipOrderResponse, error) {
	// 调用服务方法
	err := s.orderService.ShipOrder(ctx, req.OrderID)
	if err != nil {
		log.Printf("RPC ShipOrder error: %v", err)
		return nil, err
	}

	// 构造响应
	return &ShipOrderResponse{
		Success: true,
	}, nil
}

// CancelOrder RPC方法：取消订单
func (s *OrderRPCServer) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CancelOrderResponse, error) {
	// 调用服务方法
	err := s.orderService.CancelOrder(ctx, req.OrderID)
	if err != nil {
		log.Printf("RPC CancelOrder error: %v", err)
		return nil, err
	}

	// 构造响应
	return &CancelOrderResponse{
		Success: true,
	}, nil
}

// RPC请求和响应结构体定义

// AddressRequest 地址请求
type AddressRequest struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}

// OrderItemRequest 订单项请求
type OrderItemRequest struct {
	ProductID   uint32
	ProductName string
	Quantity    int32
	UnitPrice   float64
	Cost        float64
}

// PlaceOrderRequest 创建订单请求
type PlaceOrderRequest struct {
	UserID       uint32
	UserCurrency string
	Address      *AddressRequest
	Email        string
	OrderItems   []*OrderItemRequest
}

// PlaceOrderResponse 创建订单响应
type PlaceOrderResponse struct {
	OrderID string
}

// GetOrderRequest 获取订单请求
type GetOrderRequest struct {
	OrderID string
}

// OrderItemResponse 订单项响应
type OrderItemResponse struct {
	ProductID   uint32
	ProductName string
	Quantity    int32
	UnitPrice   float64
	Cost        float64
}

// OrderResponse 订单响应
type OrderResponse struct {
	OrderID       string
	UserID        uint32
	UserCurrency  string
	TotalAmount   float64
	ItemCount     int32
	OrderStatus   string
	CreatedAt     int64
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
	Email         string
	Items         []*OrderItemResponse
}

// GetOrderResponse 获取订单响应
type GetOrderResponse struct {
	Order *OrderResponse
}

// ListOrdersRequest 获取订单列表请求
type ListOrdersRequest struct {
	UserID uint32
}

// ListOrdersResponse 获取订单列表响应
type ListOrdersResponse struct {
	Orders []*OrderResponse
}

// MarkOrderPaidRequest 标记订单已支付请求
type MarkOrderPaidRequest struct {
	OrderID       string
	PaymentMethod string
}

// MarkOrderPaidResponse 标记订单已支付响应
type MarkOrderPaidResponse struct {
	Success bool
}

// ShipOrderRequest 标记订单已发货请求
type ShipOrderRequest struct {
	OrderID string
}

// ShipOrderResponse 标记订单已发货响应
type ShipOrderResponse struct {
	Success bool
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	OrderID string
}

// CancelOrderResponse 取消订单响应
type CancelOrderResponse struct {
	Success bool
}
