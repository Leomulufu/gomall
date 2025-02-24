# 订单服务 (Order Service)

## 功能描述
订单服务负责处理订单的创建、查询和支付状态更新，以及超时订单的自动取消。

## 已实现功能
1. 订单管理
   - 创建订单 (PlaceOrder)
   - 查询订单列表 (ListOrders)
   - 标记订单支付 (MarkOrderPaid)
2. 自动化任务
   - 超时订单自动取消（30分钟未支付）

## 技术栈
- Go
- MySQL
- GORM
- Kitex (准备集成)

## 测试覆盖
- 单元测试：service/order_service_test.go
- 功能测试：main.go 中的测试函数

## 数据库设计
详见 model/order.go 中的结构定义 