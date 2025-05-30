package uredis

import (
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/utime"
	"go.uber.org/zap"
	"math/rand"
	"runtime/debug"
	"sync"

	redis "github.com/go-redis/redis/v8"
)

// 高频数据缓存
var RedisHighDataMap = new(sync.Map)

// 高频数据结构体
type RedisHighData struct {
	Data        []byte // Redis返回数据
	Err         error  // Redis返回err
	Timestamp   int64  // 过期时间戳
	DataHashAll map[string]string
}

// 默认本地缓存时间/秒
var DefaultTimeLen int64 = 60

// 优先从服务器获取Redis数据，在服务器保存指定时间，拉取失败会从Redis获取
func GetStringByTimeDefault(redisCon *redis.Client, key string) ([]byte, error) {
	return GetStringByTime(redisCon, key, DefaultTimeLen)
}

func GetStringByTime(redisCon *redis.Client, key string, timeLen int64) ([]byte, error) {
	dataKey := key
	nowTimestamp := utime.GetNowUnix()
	if v, ok := RedisHighDataMap.Load(dataKey); ok {
		highDataOld := v.(*RedisHighData)
		// 本地缓存的数据还没有过期
		if highDataOld.Timestamp >= nowTimestamp {
			return highDataOld.Data, highDataOld.Err
		}
	}
	// 服务器本地缓存已经过期重新从服务器获取
	highData := GetStringRedisHighData(redisCon, key, timeLen)
	if highData.Err == nil {
		RedisHighDataMap.Store(dataKey, highData)
	}
	return highData.Data, highData.Err
}

// 重新获取Redis数据
func GetStringRedisHighData(redisCon *redis.Client, key string, timeLen int64) *RedisHighData {
	highData := new(RedisHighData)
	nowTimestamp := utime.GetNowUnix()
	// 从服务器获取
	data, err := GetString(redisCon, key)
	highData.Data = data
	highData.Err = err
	highData.Timestamp = nowTimestamp + timeLen
	log.Info("highredis高频缓存", zap.Any("highData", highData), zap.Any("key", key))
	return highData
}

// 优先从服务器获取Redis数据，在服务器保存指定时间，拉取失败会从Redis获取
func GetHashByTimeDefault(redisCon *redis.Client, key, fieId string) ([]byte, error) {
	if rand.Int31n(100) == 0 {
		stack := string(debug.Stack())
		log.Info("GetHashByTimeDefault", zap.String("stack", stack))
		//qyweixin.SendMessage(stack, env.QYWeiXinRpc)
	}

	return GetHashByTime(redisCon, key, fieId, DefaultTimeLen)
}

func GetHashByTime(redisCon *redis.Client, key, fieId string, timeLen int64) ([]byte, error) {
	dataKey := key + fieId
	nowTimestamp := utime.GetNowUnix()
	if v, ok := RedisHighDataMap.Load(dataKey); ok {
		highDataOld := v.(*RedisHighData)
		log.Debug("highredis高频缓存", zap.Any("highData", highDataOld), zap.Any("key", key), zap.Any("fieId", fieId))
		// 本地缓存的数据还没有过期
		if highDataOld.Timestamp >= nowTimestamp {
			return highDataOld.Data, highDataOld.Err
		}
	}
	// 服务器本地缓存已经过期重新从服务器获取
	highData := GetHashRedisHighData(redisCon, key, fieId, timeLen)
	if highData.Err == nil {
		RedisHighDataMap.Store(dataKey, highData)
	}
	log.Debug("highredis高频缓存", zap.Any("highData", highData), zap.Any("key", key), zap.Any("fieId", fieId))
	return highData.Data, highData.Err
}

// 重新获取Redis数据
func GetHashRedisHighData(redisCon *redis.Client, key, fieId string, timeLen int64) *RedisHighData {
	highData := new(RedisHighData)
	nowTimestamp := utime.GetNowUnix()
	// 从服务器获取
	data, err := GetHash(redisCon, key, fieId)
	highData.Data = data
	highData.Err = err
	highData.Timestamp = nowTimestamp + timeLen
	log.Info("highredis高频缓存", zap.Any("highData", highData), zap.Any("key", key), zap.Any("fieId", fieId))
	return highData
}

// 优先从服务器获取Redis数据，在服务器保存指定时间，拉取失败会从Redis获取
func GetHashAllByTimeDefault(redisCon *redis.Client, key string) (map[string]string, error) {
	return GetHashAllByTime(redisCon, key, DefaultTimeLen)
}

func GetHashAllByTime(redisCon *redis.Client, key string, timeLen int64) (map[string]string, error) {
	dataKey := key
	nowTimestamp := utime.GetNowUnix()
	if v, ok := RedisHighDataMap.Load(dataKey); ok {
		highDataOld := v.(*RedisHighData)
		log.Debug("highredis高频缓存", zap.Any("highData", highDataOld), zap.Any("key", key))
		// 本地缓存的数据还没有过期
		if highDataOld.Timestamp >= nowTimestamp {
			return highDataOld.DataHashAll, highDataOld.Err
		}
	}
	// 服务器本地缓存已经过期重新从服务器获取
	highData := GetHashAllRedisHighData(redisCon, key, timeLen)
	if highData.Err == nil {
		RedisHighDataMap.Store(dataKey, highData)
	}
	log.Debug("highredis高频缓存", zap.Any("highData", highData), zap.Any("key", key))
	return highData.DataHashAll, highData.Err
}

// 重新获取Redis数据
func GetHashAllRedisHighData(redisCon *redis.Client, key string, timeLen int64) *RedisHighData {
	highData := new(RedisHighData)
	nowTimestamp := utime.GetNowUnix()
	// 从服务器获取
	data, err := GetHashAll(redisCon, key)
	highData.DataHashAll = data
	highData.Err = err
	highData.Timestamp = nowTimestamp + timeLen
	log.Info("highredis高频缓存 highData", zap.Any("key", key))
	return highData
}
