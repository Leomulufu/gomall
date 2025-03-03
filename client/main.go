package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OrderService 定义了订单服务的接口
type OrderService interface {
	PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (string, error)
	ListOrders(ctx context.Context, req *ListOrderRequest) ([]*Order, error)
	MarkOrderPaid(ctx context.Context, req *MarkOrderPaidRequest) error
}

// OrderServiceClient 实现了 OrderService 接口
type OrderServiceClient struct {
	conn *grpc.ClientConn
}

// NewOrderServiceClient 创建一个新的订单服务客户端
func NewOrderServiceClient(conn *grpc.ClientConn) *OrderServiceClient {
	return &OrderServiceClient{conn: conn}
}

// 地址信息
type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}

// 订单项
type OrderItem struct {
	ProductID   uint32
	ProductName string
	Quantity    int32
	UnitPrice   float32
	Cost        float32
}

// Order 表示一个订单
type Order struct {
	OrderID      string
	UserID       uint32
	UserCurrency string
	Address      *Address
	Email        string
	OrderItems   []*OrderItem
	Status       string
	CreatedAt    int32
}

// PlaceOrderRequest 创建订单的请求
type PlaceOrderRequest struct {
	UserID       uint32
	UserCurrency string
	Address      *Address
	Email        string
	OrderItems   []*OrderItem
}

// ListOrderRequest 查询订单列表的请求
type ListOrderRequest struct {
	UserID uint32
}

// MarkOrderPaidRequest 标记订单已支付的请求
type MarkOrderPaidRequest struct {
	UserID  uint32
	OrderID string
}

// PlaceOrder 实现创建订单
func (c *OrderServiceClient) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (string, error) {
	// TODO: 实现 gRPC 调用
	return "test-order-id", nil
}

// ListOrders 实现查询订单列表
func (c *OrderServiceClient) ListOrders(ctx context.Context, req *ListOrderRequest) ([]*Order, error) {
	// TODO: 实现 gRPC 调用
	return []*Order{
		{
			OrderID:      "test-order-id",
			UserID:       req.UserID,
			UserCurrency: "CNY",
			Status:       "pending",
		},
	}, nil
}

// MarkOrderPaid 实现标记订单已支付
func (c *OrderServiceClient) MarkOrderPaid(ctx context.Context, req *MarkOrderPaidRequest) error {
	// TODO: 实现 gRPC 调用
	return nil
}

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := NewOrderServiceClient(conn)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("开始测试订单服务...")

	// 1. 测试创建订单
	address := &Address{
		StreetAddress: "123 抖音路",
		City:          "北京",
		State:         "北京",
		Country:       "中国",
		ZipCode:       100000,
	}

	items := []*OrderItem{
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

	orderID, err := client.PlaceOrder(ctx, &PlaceOrderRequest{
		UserID:       1001,
		UserCurrency: "CNY",
		Address:      address,
		Email:        "test@example.com",
		OrderItems:   items,
	})
	if err != nil {
		log.Printf("创建订单失败: %v", err)
	} else {
		log.Printf("订单创建成功，订单ID: %s", orderID)
	}

	// 2. 测试查询订单列表
	orders, err := client.ListOrders(ctx, &ListOrderRequest{
		UserID: 1001,
	})
	if err != nil {
		log.Printf("查询订单列表失败: %v", err)
	} else {
		log.Printf("找到 %d 个订单", len(orders))
		for _, order := range orders {
			log.Printf("订单ID: %s, 状态: %s", order.OrderID, order.Status)
		}
	}

	// 3. 测试标记订单已支付
	err = client.MarkOrderPaid(ctx, &MarkOrderPaidRequest{
		UserID:  1001,
		OrderID: orderID,
	})
	if err != nil {
		log.Printf("标记订单已支付失败: %v", err)
	} else {
		log.Printf("订单 %s 已标记为已支付", orderID)
	}

	log.Println("订单服务测试完成!")
}
