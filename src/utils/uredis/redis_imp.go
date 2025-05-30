package uredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
)

// IRedis redis client interface
type IRedis interface {
	// KEYS get patten key array
	KEYS(patten string) ([]string, error)

	// SCAN get patten key array
	SCAN(patten string) ([]string, error)

	// DEL delete k-v
	DEL(key string) (int, error)

	// DELALL delete key array
	DELALL(key []string) (int, error)

	// GET get k-v
	GET(key string) (string, error)

	// SET set k-v
	//SET(key string, value string) (int64, error)

	// SETEX set k-v expire seconds
	SETEX(key string, sec int, value string) (int64, error)

	// EXPIRE set key expire seconds
	EXPIRE(key string, sec int64) (int64, error)

	// HGETALL get map of key
	HGETALL(key string) (map[string]string, error)

	// HGET get value of key-field
	HGET(key string, field string) (string, error)

	// HSET set value of key-field
	//HSET(key string, field string, value string) (int64, error)

	// Write 向redis中写入多组数据
	//Write(data RedisDataArray)
}

// RedisClient redis client instance
type RedisClient struct {
	con     *redis.Client
	connOpt RedisConnOpt
}

// RedisConnOpt connect redis options
type RedisConnOpt struct {
	Enable   bool
	Host     string
	Port     int32
	Password string
	Index    int32
	TTL      int32
	Db       int
}

// NewRedis new redis client
func NewRedis(opt RedisConnOpt) *RedisClient {
	con := redis.NewClient(&redis.Options{
		PoolSize: 100,
		Addr:     strings.TrimSpace(opt.Host),
		Password: opt.Password,
		DB:       opt.Db,
	})

	ctx := context.Background()
	_, err := con.Ping(ctx).Result()
	if err != nil {
		panic("redis初始化失败, err:" + err.Error())
	}

	return &RedisClient{
		connOpt: opt,
		con:     con,
	}
}
