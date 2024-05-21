package handler

import (
	"context"

	pb "github.com/dreamilk/rpc_server/api"
)

type GreeterHandler struct {
	pb.UnimplementedGreeterServer
}

func (GreeterHandler) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "Hello " + req.GetName(),
	}, nil
}
