package grpc

import (
	"log"
	"net"
)

// StartServer 启动 gRPC 服务器
func StartServer() error {
	// 监听端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Printf("监听端口失败: %v", err)
		return err
	}

	log.Printf("gRPC 服务器正在监听: %v", lis.Addr())

	// TODO: 实现完整的 gRPC 服务器
	select {} // 暂时阻塞，防止程序退出

	return nil
}
