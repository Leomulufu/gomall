package model

import "time"

type Order struct {
	OrderID       int64     `gorm:"primaryKey;column:order_id;type:bigint;autoIncrement"`
	UserID        int64     `gorm:"column:user_id;type:bigint"`
	UserCurrency  string    `gorm:"column:user_currency;type:varchar(20)"`
	StreetAddress string    `gorm:"column:street_address;type:varchar(200)"`
	City          string    `gorm:"column:city;type:varchar(50)"`
	State         string    `gorm:"column:state;type:varchar(50)"`
	Country       string    `gorm:"column:country;type:varchar(50)"`
	ZipCode       string    `gorm:"column:zip_code;type:varchar(10)"`
	Email         string    `gorm:"column:email;type:varchar(100)"`
	OrderStatus   string    `gorm:"column:order_status;type:varchar(20)"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

type OrderItem struct {
	OrderItemID int64   `gorm:"primaryKey;column:order_item_id;type:bigint;autoIncrement"`
	OrderID     int64   `gorm:"column:order_id;type:bigint;foreignKey:OrderID;references:OrderID"`
	ProductID   int64   `gorm:"column:product_id;type:bigint"`
	Quantity    int     `gorm:"column:quantity"`
	Cost        float32 `gorm:"column:cost"`
}
