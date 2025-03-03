package task

import (
	"context"
	"log"
	"order_service/dao"
	"order_service/model"
	"sync"
	"time"
)

var (
	// 确保任务不会并发运行
	cancelMutex = &sync.Mutex{}
	reminderMutex = &sync.Mutex{}
)

// StartOrderTasks 启动所有订单相关的定时任务
func StartOrderTasks() {
	// 启动取消超时订单任务
	StartOrderCancelTask()
	
	// 启动订单支付提醒任务
	StartOrderPaymentReminderTask()
}

// StartOrderCancelTask 启动定时取消订单任务
func StartOrderCancelTask() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		
		log.Println("Order cancellation task started")
		
		// 立即执行一次，不等待第一个ticker
		if err := cancelTimeoutOrders(); err != nil {
			log.Printf("Error cancelling timeout orders: %v", err)
		}
		
		for range ticker.C {
			if err := cancelTimeoutOrders(); err != nil {
				log.Printf("Error cancelling timeout orders: %v", err)
			}
		}
	}()
}

// StartOrderPaymentReminderTask 启动订单支付提醒任务
func StartOrderPaymentReminderTask() {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		
		log.Println("Order payment reminder task started")
		
		// 立即执行一次
		if err := sendPaymentReminders(); err != nil {
			log.Printf("Error sending payment reminders: %v", err)
		}
		
		for range ticker.C {
			if err := sendPaymentReminders(); err != nil {
				log.Printf("Error sending payment reminders: %v", err)
			}
		}
	}()
}

// cancelTimeoutOrders 取消超时订单
func cancelTimeoutOrders() error {
	// 使用互斥锁确保任务不会并发运行
	cancelMutex.Lock()
	defer cancelMutex.Unlock()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// 超过30分钟未支付的订单将被取消
	timeout := time.Now().Add(-30 * time.Minute)
	
	result := dao.DB.WithContext(ctx).Model(&model.Order{}).
		Where("order_status = ? AND created_at < ?", model.OrderStatusPending, timeout).
		Update("order_status", model.OrderStatusCancelled)
		
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected > 0 {
		log.Printf("Cancelled %d timeout orders", result.RowsAffected)
	}
	
	return nil
}

// sendPaymentReminders 发送支付提醒
func sendPaymentReminders() error {
	// 使用互斥锁确保任务不会并发运行
	reminderMutex.Lock()
	defer reminderMutex.Unlock()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// 创建时间在15-25分钟之间的未支付订单需要发送提醒
	startTime := time.Now().Add(-25 * time.Minute)
	endTime := time.Now().Add(-15 * time.Minute)
	
	var orders []model.Order
	if err := dao.DB.WithContext(ctx).
		Where("order_status = ? AND created_at BETWEEN ? AND ?", 
			model.OrderStatusPending, startTime, endTime).
		Find(&orders).Error; err != nil {
		return err
	}
	
	for _, order := range orders {
		// 这里应该调用邮件服务发送提醒
		// 在实际项目中，这里应该调用邮件服务的API
		log.Printf("Sending payment reminder for order %s to %s", order.OrderID, order.Email)
	}
	
	if len(orders) > 0 {
		log.Printf("Sent %d payment reminders", len(orders))
	}
	
	return nil
} 