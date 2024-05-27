package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/dreamilk/rpc_server/config"
	"github.com/dreamilk/rpc_server/log"
)

var s *grpc.Server

func init() {
	s = grpc.NewServer()
}

func Serve() error {
	var ctx = context.Background()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.DefaultConf.Port))
	if err != nil {
		log.Errorf(ctx, "failed to listen: %v", err)
		return err
	}

	return s.Serve(lis)
}

func RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.RegisterService(desc, impl)
}
