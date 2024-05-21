package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"

	pb "github.com/dreamilk/rpc_server/api"
	"github.com/dreamilk/rpc_server/handler"
)

func initGrpc() {
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

func regeisterConsul() {
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println("api.NewClient err:", err)
		return
	}

	registerService := api.AgentServiceRegistration{
		ID:      "id01",
		Tags:    []string{"grpc"},
		Name:    "greeter",
		Address: "127.0.0.1",
		Port:    9989,
		Check: &api.AgentServiceCheck{
			HTTP:     "http://host.docker.internal:8888/health",
			Timeout:  "1s",
			Interval: "5s",
		},
	}

	err = consulClient.Agent().ServiceRegister(&registerService)
	if err != nil {
		fmt.Println("consulClient.Agent().ServiceRegister err:", err)
		return
	}
}

func initCheck() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {

	go regeisterConsul()

	go initGrpc()

	go initCheck()

	select {}

}
