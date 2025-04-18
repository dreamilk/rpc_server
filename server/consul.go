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
		<-c
		if err := deregeisterService(ctx, conf.Id, conf.Consul); err != nil {
			log.Error(ctx, "deregister service err", zap.Error(err))
		}
		log.Info(ctx, "deregister service")

		os.Exit(0)
	}()

	if err := registerService(ctx, conf.Id, conf.AppName, conf.Port, conf.Addr, conf.Consul); err != nil {
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

	healthCheckUrl := "http://" + serviceAddr + ":8888/health"
	log.Info(ctx, "health check url", zap.String("url", healthCheckUrl))

	registerService := api.AgentServiceRegistration{
		ID:      id,
		Tags:    []string{"grpc"},
		Name:    name,
		Address: serviceAddr,
		Port:    port,
		Check: &api.AgentServiceCheck{
			HTTP:     healthCheckUrl,
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

func deregeisterService(ctx context.Context, id string, consulAddr string) error {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddr

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Error(ctx, "new consul client err", zap.Error(err))
		return err
	}

	err = consulClient.Agent().ServiceDeregister(id)
	if err != nil {
		log.Error(ctx, "deregister service err", zap.Error(err))
		return err
	}
	return nil
}
