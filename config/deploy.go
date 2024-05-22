package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DeployConfig struct {
	Port    int    `yml:"port"`
	AppName string `yml:"app_name"`
	Id      string `yml:"id"`
}

func ReadDeploy() *DeployConfig {
	config := DeployConfig{
		AppName: "app_name",
		Id:      "id",
		Port:    9000,
	}

	f, err := os.ReadFile("./deploy.yml")
	if err != nil {
		return &config
	}

	if err := yaml.Unmarshal(f, &config); err != nil {
		return &config
	}
	return &config
}
