package main

import (
	"context"
	"flag"
	"fmt"
	proto2 "grpcdemo/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 10089, "the port to serve on")

type ecServer struct {
	proto2.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *proto2.EchoRequest) (*proto2.EchoResponse, error) {
	fmt.Printf("%s \n", req.Message)
	return &proto2.EchoResponse{Message: req.Message}, nil
}

func (e *ecServer) ServerStreamingEcho(req *proto2.EchoRequest, stream proto2.Echo_ServerStreamingEchoServer) error {
	log.Printf("Recved %v", req.GetMessage())
	// 具体返回多少个response根据业务逻辑调整
	for i := 0; i < 2; i++ {
		// 通过 send 方法不断推送数据
		err := stream.Send(&proto2.EchoResponse{Message: "Reply Hello World"})
		if err != nil {
			log.Fatalf("Send error:%v", err)
			return err
		}
	}
	// 返回nil表示已经完成响应
	return nil
}

func main() {
	flag.Parse()
	log.Printf("server starting on port %d...\n", *port)

	s := grpc.NewServer()
	proto2.RegisterEchoServer(s, &ecServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
