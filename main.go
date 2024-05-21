package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/dreamilk/rpc_server/api"
	"github.com/dreamilk/rpc_server/handler"
)

func main() {
	lis, err := net.Listen("tcp", ":9989")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &handler.GreeterHandler{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
