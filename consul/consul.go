package consul

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dreamilk/rpc_gateway/log"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"github.com/dreamilk/rpc_server/config"
)

func Register(conf *config.DeployConfig) {
	ctx := context.Background()

	go healthCheck()

	// deregister
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c
		if err := deregeisterService(conf.Id); err != nil {
			log.Error(ctx, "de", zap.Error(err))
		}
		log.Info(ctx, "", zap.Any("sig", sig))

		os.Exit(0)
	}()

	if err := registerService(conf.Id, conf.AppName, conf.Port); err != nil {
		log.Error(ctx, "", zap.Error(err))
	}
}

func registerService(id string, name string, port int) error {
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println("api.NewClient err:", err)
		return err
	}

	registerService := api.AgentServiceRegistration{
		ID:      id,
		Tags:    []string{"grpc"},
		Name:    name,
		Address: "127.0.0.1",
		Port:    port,
		Check: &api.AgentServiceCheck{
			HTTP:     "http://host.docker.internal:8888/health",
			Timeout:  "1s",
			Interval: "5s",
		},
	}

	err = consulClient.Agent().ServiceRegister(&registerService)
	if err != nil {
		fmt.Println("consulClient.Agent().ServiceRegister err:", err)
		return err
	}
	return nil
}

func deregeisterService(id string) error {
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println("api.NewClient err:", err)
		return err
	}

	err = consulClient.Agent().ServiceDeregister(id)
	if err != nil {
		fmt.Println("consulClient.Agent().ServiceDeregister err:", err)
		return err
	}
	return nil
}
