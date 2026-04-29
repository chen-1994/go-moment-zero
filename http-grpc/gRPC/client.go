package main

import (
	"context"
	. "http-grpc/gRPC/hello"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 連接服務端
	conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	c := NewGreeterClient(conn)
	// 調用遠程方法
	r, err := c.SayHello(context.Background(), &HelloRequest{Name: "Gin 學習者"})
	if err != nil {
		log.Fatalf("調用失敗: %v", err)
	}
	log.Printf("收到回應: %s", r.GetMessage())
}
