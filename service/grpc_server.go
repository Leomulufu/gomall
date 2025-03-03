package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"order_service/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer 是gRPC服务器
type GRPCServer struct {
	server  *grpc.Server
	address string
	rpcImpl *OrderRPCServer
}

// NewGRPCServer 创建新的gRPC服务器
func NewGRPCServer(address string) *GRPCServer {
	return &GRPCServer{
		address: address,
		rpcImpl: NewOrderRPCServer(),
	}
}

// Start 启动gRPC服务器
func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.server = grpc.NewServer()

	// 注册服务实现
	// 在实际项目中，这里应该使用生成的gRPC代码注册服务
	// RegisterOrderServiceServer(s.server, s.rpcImpl)

	// 启用反射服务，便于调试
	reflection.Register(s.server)

	log.Printf("gRPC server listening on %s", s.address)

	return s.server.Serve(lis)
}

// Stop 停止gRPC服务器
func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
		log.Println("gRPC server stopped")
	}
}

// 以下是gRPC服务接口实现
// 在实际项目中，这些方法应该实现生成的gRPC接口

// PlaceOrder 实现创建订单接口
func (s *OrderRPCServer) PlaceOrderGRPC(ctx context.Context, req *PlaceOrderGRPCRequest) (*PlaceOrderGRPCResponse, error) {
	// 转换请求
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
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   float64(item.UnitPrice),
			Cost:        float64(item.Cost),
		})
	}

	// 调用服务
	order, err := s.orderService.PlaceOrder(ctx, req.UserId, req.UserCurrency, address, req.Email, orderItems)
	if err != nil {
		return nil, err
	}

	// 返回响应
	return &PlaceOrderGRPCResponse{
		Order: &OrderResultGRPC{
			OrderId: order.OrderID,
		},
	}, nil
}

// ListOrders 实现获取订单列表接口
func (s *OrderRPCServer) ListOrderGRPC(ctx context.Context, req *ListOrderGRPCRequest) (*ListOrderGRPCResponse, error) {
	// 调用服务
	orders, err := s.orderService.ListOrders(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var grpcOrders []*OrderGRPC
	for _, order := range orders {
		grpcOrders = append(grpcOrders, &OrderGRPC{
			OrderId:      order.OrderID,
			UserId:       order.UserID,
			UserCurrency: order.UserCurrency,
			CreatedAt:    int32(order.CreatedAt.Unix()),
			Address: &AddressGRPC{
				StreetAddress: order.StreetAddress,
				City:          order.City,
				State:         order.State,
				Country:       order.Country,
				ZipCode:       order.ZipCode,
			},
			Email: order.Email,
		})
	}

	return &ListOrderGRPCResponse{
		Orders: grpcOrders,
	}, nil
}

// MarkOrderPaid 实现标记订单已支付接口
func (s *OrderRPCServer) MarkOrderPaidGRPC(ctx context.Context, req *MarkOrderPaidGRPCRequest) (*MarkOrderPaidGRPCResponse, error) {
	// 调用服务
	err := s.orderService.MarkOrderPaid(ctx, req.OrderId, "credit_card")
	if err != nil {
		return nil, err
	}

	return &MarkOrderPaidGRPCResponse{}, nil
}

// gRPC请求和响应结构体定义
// 在实际项目中，这些结构体应该由protobuf生成

type AddressGRPC struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}

type OrderItemGRPC struct {
	ProductId   uint32
	ProductName string
	Quantity    int32
	UnitPrice   float32
	Cost        float32
}

type PlaceOrderGRPCRequest struct {
	UserId       uint32
	UserCurrency string
	Address      *AddressGRPC
	Email        string
	OrderItems   []*OrderItemGRPC
}

type OrderResultGRPC struct {
	OrderId string
}

type PlaceOrderGRPCResponse struct {
	Order *OrderResultGRPC
}

type ListOrderGRPCRequest struct {
	UserId uint32
}

type OrderGRPC struct {
	OrderItems   []*OrderItemGRPC
	OrderId      string
	UserId       uint32
	UserCurrency string
	Address      *AddressGRPC
	Email        string
	CreatedAt    int32
}

type ListOrderGRPCResponse struct {
	Orders []*OrderGRPC
}

type MarkOrderPaidGRPCRequest struct {
	UserId  uint32
	OrderId string
}

type MarkOrderPaidGRPCResponse struct{}
