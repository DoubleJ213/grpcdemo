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
