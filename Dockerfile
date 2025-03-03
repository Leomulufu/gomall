FROM golang:1.20-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o order_service .

# 使用轻量级的alpine镜像
FROM alpine:latest

WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/order_service .

# 暴露gRPC端口
EXPOSE 50051

# 运行服务
CMD ["./order_service"] 