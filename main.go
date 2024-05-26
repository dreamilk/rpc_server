package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/dreamilk/rpc_server/api"
	"github.com/dreamilk/rpc_server/config"
	"github.com/dreamilk/rpc_server/handler"
	"github.com/dreamilk/rpc_server/log"
)

func main() {
	ctx := context.Background()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.DefaultConf.Port))
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
