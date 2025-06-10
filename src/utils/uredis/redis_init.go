package uredis

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"eshop_server/src/utils/log"
)

var (
	RedisCon *redis.Client
)

func GetRedis() *redis.Client {
	return RedisCon
}

//初始化redis
func InitRedis(host string, pwd string, db int) {
	RedisCon = New(host, pwd, db)
	log.Info("初始化redis完成", zap.Any("RedisCon", RedisCon))
}
