package main

import (
	"context"
	"net"

	"github.com/dreamilk/rpc_gateway/log"
	"google.golang.org/grpc"

	pb "github.com/dreamilk/rpc_server/api"
	"github.com/dreamilk/rpc_server/consul"
	"github.com/dreamilk/rpc_server/handler"
)

func init() {
	consul.Register("app123", "server", 9989)
}

func main() {
	ctx := context.Background()

	lis, err := net.Listen("tcp", ":9989")
	if err != nil {
		log.Errorf(ctx, "failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &handler.GreeterHandler{})

	log.Infof(ctx, "server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Errorf(ctx, "failed to serve: %v", err)
	}
}
