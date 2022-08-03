package main

import (
	"context"
	"flag"
	"fmt"
	proto2 "grpcdemo/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr = flag.String("addr", "localhost:10086", "the address to connect to")

func callUnaryEcho(client proto2.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &proto2.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func main() {
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile(
		"demo1/server.crt", "localhost",
	)

	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	callUnaryEcho(proto2.NewEchoClient(conn), "hello demo1")
}
