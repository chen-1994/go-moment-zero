package main

import (
	"context"
	. "http-grpc/gRPC/hello"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	UnimplementedGreeterServer
}

// 實作 SayHello 邏輯
func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "你好 " + in.Name}, nil
}

func main() {
	lis, _ := net.Listen("tcp", ":50051") // 監聽 50051 端口
	s := grpc.NewServer()
	RegisterGreeterServer(s, &server{}) // 註冊服務
	log.Println("gRPC 服務啟動在 :50051...")
	s.Serve(lis)
}
