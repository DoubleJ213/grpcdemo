package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	proto2 "grpcdemo/proto"
	"io/ioutil"
	"log"
	"net"
	"os/exec"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var port = flag.Int("port", 10088, "the port to serve on")

type ecServer struct {
	proto2.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *proto2.EchoRequest) (*proto2.EchoResponse, error) {
	fmt.Printf("%s \n", req.Message)
	return &proto2.EchoResponse{Message: req.Message}, nil
}

func main() {
	c := exec.Command("bash", "-c", "cd demo3; sh init.sh")
	output, err := c.CombinedOutput()
	fmt.Println(string(output))

	flag.Parse()
	log.Printf("server starting on port %d...\n", *port)

	certificate, err := tls.LoadX509KeyPair("demo3/server.crt", "demo3/server.key")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("demo3/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert, // NOTE: this is optional!
		ClientCAs:    certPool,
	})

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
