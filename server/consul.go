package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"github.com/dreamilk/rpc_server/config"
	"github.com/dreamilk/rpc_server/log"
)

func init() {
	register(&config.DefaultConf)
}

func register(conf *config.DeployConfig) {
	ctx := context.Background()

	go healthCheck()

	// deregister
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c
		if err := deregeisterService(ctx, conf.Id); err != nil {
			log.Error(ctx, "deregister service err", zap.Error(err))
		}
		log.Info(ctx, "", zap.Any("sig", sig))

		os.Exit(0)
	}()

	if err := registerService(ctx, conf.Id, conf.AppName, conf.Port, conf.Addr, conf.Consul.Addr); err != nil {
		log.Error(ctx, "register service err", zap.Error(err))
	}
}

func registerService(ctx context.Context, id string, name string, port int, serviceAddr string, consulAddr string) error {
	consulConfig := api.DefaultConfig()

	consulConfig.Address = consulAddr

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Error(ctx, "new consul client err", zap.Error(err))
		return err
	}

	registerService := api.AgentServiceRegistration{
		ID:      id,
		Tags:    []string{"grpc"},
		Name:    name,
		Address: serviceAddr,
		Port:    port,
		Check: &api.AgentServiceCheck{
			HTTP:     "http://" + serviceAddr + ":8888/health",
			Timeout:  "1s",
			Interval: "5s",
		},
	}

	err = consulClient.Agent().ServiceRegister(&registerService)
	if err != nil {
		log.Error(ctx, "consul register service err", zap.Error(err))
		return err
	}
	return nil
}

func deregeisterService(ctx context.Context, id string) error {
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Error(ctx, "api.NewClient err", zap.Error(err))
		return err
	}

	err = consulClient.Agent().ServiceDeregister(id)
	if err != nil {
		log.Error(ctx, "consulClient.Agent().ServiceDeregister err", zap.Error(err))
		return err
	}
	return nil
}
