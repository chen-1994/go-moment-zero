package main

import (
	"context"
	. "http-grpc/gRPC2/hello"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type HelloServiceImpl struct {
	UnimplementedHelloServiceServer // 必須包含這一行
}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *String,
) (*String, error) {
	reply := &String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func (p *HelloServiceImpl) Channel(stream HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &String{Value: "hello:" + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

// gRPC流
func main() {
	server := grpc.NewServer()
	RegisterHelloServiceServer(server, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC PubSub server running on :1234")
	server.Serve(lis)
}
