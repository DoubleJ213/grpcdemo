package main

import (
	"context"
	"flag"
	"fmt"
	proto2 "grpcdemo/proto"
	"log"
	"net"
	"os/exec"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var port = flag.Int("port", 10087, "the port to serve on")

type ecServer struct {
	proto2.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *proto2.EchoRequest) (*proto2.EchoResponse, error) {
	fmt.Printf("%s \n", req.Message)
	return &proto2.EchoResponse{Message: req.Message}, nil
}

func main() {
	c := exec.Command("bash", "-c", "cd demo2; sh init.sh")
	output, err := c.CombinedOutput()
	fmt.Println(string(output))

	flag.Parse()
	log.Printf("server starting on port %d...\n", *port)

	creds, err := credentials.NewServerTLSFromFile("demo2/server.crt", "demo2/server.key")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	proto2.RegisterEchoServer(s, &ecServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
