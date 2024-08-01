package goredissvc

import (
	"github.com/go-redis/redis"
	"github.com/norniastar/go-infra/contract"
)

// NewSingleOption  单节点选项
func NewSingleOption(o redis.Options) contract.RedisOption {
	return func(redis contract.IRedis) {
		redis.(*redisAdapter).options = &o
	}
}

// NewClusterOption  集群选项
func NewClusterOption(o redis.ClusterOptions) contract.RedisOption {
	return func(redis contract.IRedis) {
		redis.(*redisAdapter).clusterOptions = &o
	}
}
