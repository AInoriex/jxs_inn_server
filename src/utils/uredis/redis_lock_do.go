package uredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

///////////////string
func (this *RedisLock) GetString() ([]byte, error) {
	ctx := context.Background()
	if err := this.Lock(); err != nil {
		return nil, err
	}
	result, err := this.con.Get(ctx, this.Key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (this *RedisLock) SetString(value interface{}, ex ...int64) error {
	var t = DefaultTime
	if len(ex) > 0 {
		t = time.Duration(ex[0]) * time.Second
	}
	err := this.con.Set(context.Background(), this.Key, value, t).Err()
	this.Unlock()
	return err
}

///////////////set
func (this *RedisLock) GetHash(con *redis.Client) ([]byte, error) {
	if err := this.Lock(); err != nil {
		return nil, err
	}
	result, err := con.HGet(context.Background(), this.Key, this.Field).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (this *RedisLock) SetHash(value interface{}, ex ...int64) error {
	err := this.con.HSet(context.Background(), this.Key, this.Field, value).Err()
	this.Unlock()
	return err
}
