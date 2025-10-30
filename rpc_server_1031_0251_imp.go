// 代码生成时间: 2025-10-31 02:51:14
package main

import (
    "fmt"
    "net"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

// 定义一个RPC服务接口
type MyService interface {
    ComputeSum(args *SumArgs) (*SumReply, error)
}

// 实现MyService接口
type MyServiceImpl struct{}

// SumArgs是求和方法的参数
type SumArgs struct {
    A, B int
}

// SumReply是求和方法的返回结果
type SumReply struct {
    Result int
}

// ComputeSum实现了MyService接口的ComputeSum方法
func (s *MyServiceImpl) ComputeSum(args *SumArgs) (*SumReply, error) {
    result := args.A + args.B
    return &SumReply{Result: result}, nil
}

// 启动RPC服务器
func startGRPCServer(port string, service MyService) {
    // 创建gRPC服务器
    server := grpc.NewServer()

    // 注册服务
    RegisterMyServiceServer(server, service)

    // 创建监听器
    listener, err := net.Listen("tcp", port)
    if err != nil {
        panic(fmt.Sprintf("Failed to listen: %v", err))
    }

    // 启动服务器
    fmt.Printf("Server listening on %s
", port)
    if err := server.Serve(listener); err != nil {
        panic(fmt.Sprintf("Failed to serve: %v", err))
    }
}

// main函数是程序的入口点
func main() {
    // 启动gRPC服务器
    startGRPCServer(":50051", &MyServiceImpl{})

    // 反射服务，用于gRPC的自省
    reflection.Register(server)
}

// RegisterMyServiceServer将MyService服务注册到gRPC服务器
func RegisterMyServiceServer(s *grpc.Server, srv MyService) {
    RegisterMyServiceServer(s, srv)
}
