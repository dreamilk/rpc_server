package main

import (
	pb "github.com/dreamilk/rpc_server/api"
	"github.com/dreamilk/rpc_server/handler"
	"github.com/dreamilk/rpc_server/server"
)

func main() {
	server.RegisterService(&pb.Greeter_ServiceDesc, &handler.GreeterHandler{})
	if err := server.Serve(); err != nil {
		return
	}
}
