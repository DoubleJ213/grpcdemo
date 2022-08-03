package main

import (
	"flag"
	"fmt"
	proto2 "grpcdemo/proto"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 10089, "the port to serve on")

type ecServer struct {
	proto2.UnimplementedEchoServer
}

func (e *ecServer) BidirectionalStreamingEcho(stream proto2.Echo_BidirectionalStreamingEchoServer) error {
	var (
		waitGroup sync.WaitGroup
		msgCh     = make(chan string)
	)
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for v := range msgCh {
			err := stream.Send(&proto2.EchoResponse{Message: v})
			if err != nil {
				fmt.Println("Send error:", err)
				continue
			}
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("recv error:%v", err)
			}
			fmt.Printf("Recved :%v \n", req.GetMessage())
			msgCh <- req.GetMessage()
		}
		close(msgCh)
	}()
	waitGroup.Wait()

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
