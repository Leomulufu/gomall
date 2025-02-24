package task

import (
	"order_service/dao"
	"order_service/model"
	"time"
)

// 启动定时取消订单任务
func StartOrderCancelTask() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			cancelTimeoutOrders()
		}
	}()
}

// 取消超时订单
func cancelTimeoutOrders() {
	timeout := time.Now().Add(-30 * time.Minute)
	dao.DB.Model(&model.Order{}).
		Where("order_status = ? AND created_at < ?", "pending", timeout).
		Update("order_status", "cancelled")
} 