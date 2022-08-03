package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto2 "grpcdemo/proto"
	"log"
)

var addr = flag.String("addr", "localhost:10089", "the address to connect to")

func clientStream(client proto2.EchoClient) {
	stream, err := client.ClientStreamingEcho(context.Background())
	if err != nil {
		log.Fatalf("Sum() error: %v", err)
	}
	for i := int64(0); i < 2; i++ {
		err := stream.Send(&proto2.EchoRequest{Message: "Hello World"})
		if err != nil {
			log.Printf("send error: %v", err)
			continue
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv() error: %v", err)
	}
	log.Printf("sum: %v", resp.GetMessage())

}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	clientStream(proto2.NewEchoClient(conn))
}
