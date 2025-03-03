package model

import "time"

// Order 订单模型
type Order struct {
	OrderID       string    `gorm:"primaryKey;column:order_id;type:varchar(64)"`
	UserID        uint32    `gorm:"column:user_id;type:int;index"`
	UserCurrency  string    `gorm:"column:user_currency;type:varchar(20)"`
	TotalAmount   float64   `gorm:"column:total_amount;type:decimal(10,2)"`
	ItemCount     int       `gorm:"column:item_count;type:int"`
	StreetAddress string    `gorm:"column:street_address;type:varchar(200)"`
	City          string    `gorm:"column:city;type:varchar(50)"`
	State         string    `gorm:"column:state;type:varchar(50)"`
	Country       string    `gorm:"column:country;type:varchar(50)"`
	ZipCode       int32     `gorm:"column:zip_code;type:int"`
	Email         string    `gorm:"column:email;type:varchar(100)"`
	OrderStatus   string    `gorm:"column:order_status;type:varchar(20);index"`
	PaymentMethod string    `gorm:"column:payment_method;type:varchar(20)"`
	CreatedAt     time.Time `gorm:"column:created_at;index;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
	PaidAt        time.Time `gorm:"column:paid_at;default:null"`
	ShippedAt     time.Time `gorm:"column:shipped_at;default:null"`
	DeliveredAt   time.Time `gorm:"column:delivered_at;default:null"`
}

// OrderItem 订单项模型
type OrderItem struct {
	OrderItemID int64   `gorm:"primaryKey;column:order_item_id;type:bigint;autoIncrement"`
	OrderID     string  `gorm:"column:order_id;type:varchar(64);index"`
	ProductID   uint32  `gorm:"column:product_id;type:int"`
	ProductName string  `gorm:"column:product_name;type:varchar(100)"`
	Quantity    int32   `gorm:"column:quantity"`
	UnitPrice   float64 `gorm:"column:unit_price;type:decimal(10,2)"`
	Cost        float64 `gorm:"column:cost;type:decimal(10,2)"`
}

// TableName 设置Order表名
func (Order) TableName() string {
	return "orders"
}

// TableName 设置OrderItem表名
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderStatus 订单状态常量
const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCancelled = "cancelled"
	OrderStatusRefunded  = "refunded"
)
