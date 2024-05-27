package kv

import (
	"context"
	"errors"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"github.com/dreamilk/rpc_server/log"
)

func Get(ctx context.Context, key string) (string, error) {
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Error(ctx, "new client error", zap.Error(err))
		return "", err
	}

	pair, _, err := consulClient.KV().Get(key, nil)
	if err != nil {
		log.Error(ctx, "get kv from client failed", zap.Error(err))
		return "", err
	}

	if pair == nil {
		return "", errors.New("no found kv")
	}

	return string(pair.Value), nil
}
