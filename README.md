# 抖音电商订单服务

这是抖音电商平台的订单服务模块，负责处理订单的创建、支付、发货、查询等功能。

## 功能特性

- 创建订单
- 查询订单详情
- 查询用户订单列表
- 标记订单为已支付
- 标记订单为已发货
- 标记订单为已送达
- 取消订单
- 自动取消超时未支付订单
- 发送订单支付提醒

## 技术栈

- Go 1.20+
- GORM (MySQL)
- gRPC
- 定时任务

## 项目结构

```
order_service/
├── dao/           # 数据访问层
├── model/         # 数据模型
├── service/       # 业务逻辑层
├── task/          # 定时任务
├── main.go        # 主程序入口
├── go.mod         # Go模块定义
└── README.md      # 项目说明
```

## 安装和运行

### 前置条件

- Go 1.20 或更高版本
- MySQL 5.7 或更高版本

### 安装 Go

1. 访问 [Go 官方下载页面](https://golang.org/dl/)
2. 下载适合您操作系统的安装包
3. 按照安装向导完成安装
4. 验证安装：打开命令行，输入 `go version`，应该显示已安装的 Go 版本

### 安装 MySQL

1. 访问 [MySQL 官方下载页面](https://dev.mysql.com/downloads/mysql/)
2. 下载适合您操作系统的安装包
3. 按照安装向导完成安装
4. 创建数据库：

```sql
CREATE DATABASE gomall CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 配置项目

1. 克隆项目（如果是从 Git 仓库获取）：

```bash
git clone https://github.com/yourusername/order_service.git
cd order_service
```

2. 修改数据库连接配置

编辑 `dao/db.go` 文件，更新数据库连接字符串：

```go
dsn := "用户名:密码@tcp(主机:端口)/gomall?charset=utf8mb4&parseTime=True&loc=Local"
```

### 安装依赖

```bash
go mod tidy
```

### 编译和运行

```bash
go build -o order_service
./order_service
```

或者直接运行：

```bash
go run main.go
```

### 使用 Docker 运行（可选）

如果您已安装 Docker，可以使用以下命令构建和运行：

1. 构建 Docker 镜像：

```bash
docker build -t order_service .
```

2. 运行容器：

```bash
docker run -p 50051:50051 order_service
```

## API 接口

服务通过 gRPC 提供以下接口：

- `PlaceOrder`: 创建订单
- `GetOrder`: 获取订单详情
- `ListOrders`: 获取用户订单列表
- `MarkOrderPaid`: 标记订单为已支付
- `ShipOrder`: 标记订单为已发货
- `DeliverOrder`: 标记订单为已送达
- `CancelOrder`: 取消订单

## 定时任务

- 自动取消超时未支付订单（每10分钟执行一次）
- 发送订单支付提醒（每30分钟执行一次）

## 订单状态流转

订单状态按照以下流程流转：

1. `pending`: 待支付
2. `paid`: 已支付
3. `shipped`: 已发货
4. `delivered`: 已送达

特殊状态：
- `cancelled`: 已取消
- `refunded`: 已退款

## 测试

运行单元测试：

```bash
go test ./...
```

## 故障排除

### 常见问题

1. 数据库连接失败
   - 检查数据库连接字符串是否正确
   - 确保 MySQL 服务正在运行
   - 检查防火墙设置

2. 依赖问题
   - 运行 `go mod tidy` 更新依赖
   - 检查 Go 版本是否满足要求

3. 端口冲突
   - 如果 50051 端口已被占用，可以修改 `main.go` 中的端口号

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件 