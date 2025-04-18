package config

import (
	"context"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/dreamilk/rpc_server/log"
)

var DefaultConf DeployConfig

func init() {
	ctx := context.Background()

	DefaultConf = initDeployConfig(ctx)

	log.Info(ctx, "deploy config", zap.Any("config", DefaultConf))
}

type DeployConfig struct {
	AppName string `yaml:"app_name"`
	Addr    string `yaml:"addr"`
	Port    int    `yaml:"port"`
	Id      string `yaml:"id"`
	Consul  struct {
		Addr string `yaml:"addr"`
	} `yaml:"consul"`
}

func initDeployConfig(ctx context.Context) DeployConfig {
	config := DeployConfig{
		AppName: "app_name",
		Id:      "id",
		Port:    9000,
	}

	f, err := os.ReadFile("./deploy.yml")
	if err != nil {
		log.Error(ctx, "init deploy config failed", zap.Error(err))
		return config
	}

	if err := yaml.Unmarshal(f, &config); err != nil {
		log.Error(ctx, "init deploy config failed", zap.Error(err), zap.String("json", string(f)))
		return config
	}
	return config
}
