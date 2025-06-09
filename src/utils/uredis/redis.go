package uredis

import (
	errors "eshop_server/src/utils/errors"
	log "eshop_server/src/utils/log"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"time"
)

const (
	RepeatedTimes    = 5  //重试次数
	RepeatedInterval = 20 //重试间隔：毫秒

	DefaultTime = time.Duration(-1) * time.Second // 默认过期时间
)

type dialOptions func(*redis.Options)

func New(host string, pwd string, db int) *redis.Client {
	configs := strings.Split(host, ",")
	con := redis.NewClient(&redis.Options{
		PoolSize: 100,
		Addr:     strings.TrimSpace(configs[0]),
		Password: pwd,
		DB:       db,
	})
	ctx := context.Background()
	_, err := con.Ping(ctx).Result()
	if err != nil {
		panic("redis初始化失败, err:" + err.Error())
	}
	return con
}

func Dial(host string, options ...dialOptions) *redis.Client {
	dial := &redis.Options{
		Addr:     host,
		PoolSize: 100,
	}
	for _, option := range options {
		option(dial)
	}
	ctx := context.Background()
	con := redis.NewClient(dial)
	_, err := con.Ping(ctx).Result()
	if err != nil {
		panic("redis初始化失败, err:" + err.Error())
	}
	return con
}

func DialPassword(password string) dialOptions {
	return func(options *redis.Options) {
		options.Password = password
	}
}

func DialDB(db int) dialOptions {
	return func(options *redis.Options) {
		options.DB = db
	}
}

// 设置读超时
func DialReadTimeout(timeout time.Duration) dialOptions {
	return func(options *redis.Options) {
		options.ReadTimeout = timeout
	}
}

// 设置写超时
func DialWriteTimeout(timeout time.Duration) dialOptions {
	return func(options *redis.Options) {
		options.WriteTimeout = timeout
	}
}

// 设置过期时间
func Expire(con *redis.Client, key string, t int64) error {
	err := con.Expire(context.Background(), key, time.Duration(t)*time.Second).Err()
	return err
}

// 删除key
func DelKey(con *redis.Client, key ...string) error {
	err := con.Del(context.Background(), key...).Err()
	return err
}

// 获取string key
func GetString(con *redis.Client, key string) ([]byte, error) {
	count := 1
	for {
		result, err := con.Get(context.Background(), key).Bytes()
		if err == nil {
			return result, nil
		}
		if err == redis.Nil {
			return nil, nil
		}
		count++
		if count > RepeatedTimes {
			log.Error("GetString失败 ", zap.Any("key", key), zap.Error(err))
			return nil, err
		} else {
			log.Error(fmt.Sprintf("GetString请求key=%v数据失败，err = %s, "+
				"发起第%v请求", key, err, count))
			time.Sleep(RepeatedInterval * time.Millisecond) //重试间隔
		}
	}
}


func SetString(con *redis.Client, key string, value interface{}, ex ...int64) error {
	count := 1
	var t = DefaultTime
	if len(ex) > 0 {
		t = time.Duration(ex[0]) * time.Second
	}
	for {
		err := con.Set(context.Background(), key, value, t).Err()
		if err == nil {
			return nil
		}
		count++
		if count > RepeatedTimes {
			log.Error("SetString失败", zap.Any("key", key), zap.Error(err))
			return err
		} else {
			log.Error(fmt.Sprintf("SetString请求key=%v数据失败，err = %s, "+
				"发起第%v请求", key, err, count))
			time.Sleep(RepeatedInterval * time.Millisecond) //重试间隔
		}
	}
	return nil
}

// 获取hash key
func GetHash(con *redis.Client, key string, field string) ([]byte, error) {
	count := 1
	for {
		result, err := con.HGet(context.Background(), key, field).Bytes()
		if err == nil {
			return result, nil
		}
		if err == redis.Nil {
			return nil, nil
		}
		count++
		if count > RepeatedTimes {
			log.Error(fmt.Sprintf("GetHash请求key=%v数据失败 err=%s", key, err))
			return nil, err
		} else {
			log.Error(fmt.Sprintf("GetHash请求key=%v数据失败，err = %s, "+
				"发起第%v请求", key, err, count))
			time.Sleep(RepeatedInterval * time.Millisecond) //重试间隔
		}
	}
}

func MGetHash(con *redis.Client, key string, fields []string) (map[string]*string, error) {
	count := 1
	for {
		result, err := con.HMGet(context.Background(), key, fields...).Result()
		if err == nil && len(fields) == len(result) {
			m := make(map[string]*string, len(fields))
			for index, field := range fields {
				v := result[index]
				switch v.(type) {
				case string:
					v2 := v.(string)
					m[field] = &v2
				default:
					m[field] = nil
				}
			}

			return m, nil
		}

		count++
		if count > RepeatedTimes {
			log.Error(fmt.Sprintf("MGetHash请求key=%v数据失败 err=%s", key, err.Error()))
			return nil, err
		} else {
			log.Error(fmt.Sprintf("GetHash请求key=%v数据失败，err = %s, "+
				"发起第%v请求", key, err, count))
			time.Sleep(RepeatedInterval * time.Millisecond) //重试间隔
		}
	}
}

func GetHashAll(con *redis.Client, key string) (map[string]string, error) {
	count := 1
	for {
		result, err := con.HGetAll(context.Background(), key).Result()
		if err == nil {
			return result, nil
		}
		count++
		if count > RepeatedTimes {
			log.Error(fmt.Sprintf("GetHash请求key=%v数据失败 err=%s", key, err))
			return nil, err
		} else {
			log.Error(fmt.Sprintf("GetHash请求key=%v数据失败，err = %s, "+
				"发起第%v请求", key, err, count))
			time.Sleep(RepeatedInterval * time.Millisecond) //重试间隔
		}
	}
}

// 获取hash的长度
func GetHashLen(con *redis.Client, key string) (int64, error) {
	result, err := con.HLen(context.Background(), key).Result()
	return result, err
}

func SetHash(con *redis.Client, key, field string, value []byte) error {
	count := 1
	//if len(ex) > 0 {
	//	t = time.Duration(ex[0]) * time.Second
	//}
	for {
		err := con.HSet(context.Background(), key, field, value).Err()
		if err == nil {
			return nil
		}
		count++
		if count > RepeatedTimes {
			log.Error(fmt.Sprintf("SetHash请求key=%v数据失败 err=%s", key, err))
			return err
		} else {
			log.Error(fmt.Sprintf("SetHash请求key=%v数据失败，err = %s, "+
				"发起第%v请求", key, err, count))
			time.Sleep(RepeatedInterval * time.Millisecond) //重试间隔
		}
	}
	return nil
}

// 删除hash
func DelHash(con *redis.Client, key string, field ...string) error {
	err := con.HDel(context.Background(), key, field...).Err()
	return err
}

// 有序集合
func SetSortSet(con *redis.Client, key string, score int64, value interface{}) error {
	err := con.ZAdd(context.Background(), key, &redis.Z{Score: float64(score), Member: value}).Err()
	if err != nil {
		log.Error(fmt.Sprintf("设置有序集合数据失败, key=%s, score=%v, value=%v, err=%s", key, score, value, err))
	}
	return err
}

func ZcardSortSet(con *redis.Client, key string) int64 {
	return con.ZCard(context.Background(), key).Val()
}

func SetAllSortSet(con *redis.Client, key string, value ...*redis.Z) error {
	err := con.ZAdd(context.Background(), key, value...).Err()
	if err != nil {
		log.Error(fmt.Sprintf("设置有序集合数据失败, key=%s, alue=%v, err=%s", key, value, err))
	}
	return err
}

func ListSortSet(con *redis.Client, key string, num ...int64) ([]string, error) {
	var n1, n2 int64 = 0, -1
	if len(num) > 0 {
		n1 = num[0]
	}
	if len(num) > 1 {
		n2 = num[1]
	}
	arr, err := con.ZRange(context.Background(), key, n1, n2).Result()
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// 获取从大到小的有序集合
func ListSortSetRev(con *redis.Client, key string, num ...int64) ([]string, error) {
	var n1, n2 int64 = 0, -1
	if len(num) > 0 {
		n1 = num[0]
	}
	if len(num) > 1 {
		n2 = num[1]
	}
	arr, err := con.ZRevRange(context.Background(), key, n1, n2).Result()
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func DelSortSet(con *redis.Client, key string, start, end int64) error {
	s := strconv.FormatInt(start, 10)
	e := strconv.FormatInt(end, 10)
	err := con.ZRemRangeByLex(context.Background(), key, s, e).Err()
	return err
}

func Lrange(con *redis.Client, key string, start int64, stop int64) ([]string, error) {
	arr, err := con.LRange(context.Background(), key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func SAdd(con *redis.Client, key string, value interface{}, ex ...int64) error {
	err := con.SAdd(context.Background(), key, value).Err()
	if err == nil && len(ex) > 0 {
		var t = DefaultTime
		t = time.Duration(ex[0]) * time.Second
		err = con.Expire(context.Background(), key, t).Err()
		return err
	}
	return err
}

func SMembers(con *redis.Client, key string) ([]string, error) {
	arr, err := con.SMembers(context.Background(), key).Result()
	return arr, err
}

func SRem(con *redis.Client, key string, members ...interface{}) (int64, error) {
	res, err := con.SRem(context.Background(), key, members).Result()
	return res, err
}

func SIsMember(con *redis.Client, key string, member interface{}) (bool, error) {
	return con.SIsMember(context.Background(), key, member).Result()
}

func SetNx(con *redis.Client, key string, value interface{}, ex ...int64) (bool, error) {
	if len(ex) > 0 {
		t := time.Duration(ex[0]) * time.Second
		return con.SetNX(context.Background(), key, value, t).Result()
	} else {
		return con.SetNX(context.Background(), key, value, 0).Result()
	}
}

func SetNxWait(con *redis.Client, key string, value interface{}, ex ...int64) (bool, error) {
	inval := 10
	flag := false
	err := errors.ErrRedis
	for i := 0; i <= 5; i++ {
		flag, err = SetNx(con, key, value, ex...)
		if err != nil || flag {
			fmt.Println(err)
			return flag, err
		}
		fmt.Println(i)
		time.Sleep(time.Duration(inval) * time.Millisecond)
	}
	if !flag && err == nil {
		return flag, nil
	}
	return flag, err
}

func Keys(con *redis.Client, key string) ([]string, error) {
	return con.Keys(context.Background(), key).Result()
}

func ZRange(con *redis.Client, key string, start int64, end int64) ([]string, error) {
	return con.ZRange(context.Background(), key, start, end).Result()
}

func ZRem(con *redis.Client, key string, members ...interface{}) error {
	return con.ZRem(context.Background(), key, members...).Err()
}

func ZCard(con *redis.Client, key string) (int64, error) {
	return con.ZCard(context.Background(), key).Result()
}

func ZRemRangeByRank(con *redis.Client, key string, start int64, end int64) error {
	return con.ZRemRangeByRank(context.Background(), key, start, end).Err()
}

func ZScore(con *redis.Client, key string, member string) (float64, error) {
	return con.ZScore(context.Background(), key, member).Result()
}

func ZAdd(con *redis.Client, key string, score float64, member string) error {
	return con.ZAdd(context.Background(), key, &redis.Z{Score: score, Member: member}).Err()
}

// 增加的score
func ZIncr(con *redis.Client, key string, score float64, member string) error {
	return con.ZIncr(context.Background(), key, &redis.Z{Score: score, Member: member}).Err()
}

func ZRangeWithScores(con *redis.Client, key string, start, stop int64) ([]redis.Z, error) {
	return con.ZRangeWithScores(context.Background(), key, start, stop).Result()
}

func ZRevRange(con *redis.Client, key string, start, stop int64) ([]redis.Z, error) {
	return con.ZRevRangeWithScores(context.Background(), key, start, stop).Result()
}

func LPush(con *redis.Client, key string, value interface{}) error {
	return con.LPush(context.Background(), key, value).Err()
}

func LTrim(con *redis.Client, key string, start, stop int64) error {
	return con.LTrim(context.Background(), key, start, stop).Err()
}

func LPop(con *redis.Client, key string) (string, error) {
	return con.LPop(context.Background(), key).Result()
}

func RPop(con *redis.Client, key string) (string, error) {
	return con.RPop(context.Background(), key).Result()
}

func IsExists(redisClient *redis.Client, name string, key string) (bool, error) {
	k := name + key
	ENum, err := redisClient.Exists(context.Background(), k).Result()
	if err == nil {
		if ENum > 0 {
			return true, err
		}
	}
	return false, err
}

func IncrKey(redisClient *redis.Client, name string) (int64, error) {
	k := name
	i, err := redisClient.Incr(context.Background(), k).Result()
	if err != nil {
		return 0, err
	}
	return i, err
}
func DecrKey(redisClient *redis.Client, name string) (int64, error) {
	k := name
	i, err := redisClient.Decr(context.Background(), k).Result()
	if err != nil {
		return -1, err
	}
	return i, err
}

func Incr(redisClient *redis.Client, name string, key string) error {
	k := name + key
	err := redisClient.Incr(context.Background(), k).Err()
	return err
}

func IncrBy(redisClient *redis.Client, name string, key string, val int64) (int64, error) {
	k := name + key
	val, err := redisClient.IncrBy(context.Background(), k, val).Result()
	return val, err
}

func Decr(redisClient *redis.Client, name string, key string) error {
	k := name + key
	err := redisClient.Decr(context.Background(), k).Err()
	return err
}

///////////////string
//func (cli *RedisClient) GetString(ctx context.Context,  key string) ([]byte, error) {
//	count := 1
//	for {
//		result, err := cli.con.Get(context.Background(), key).Bytes()
//		if err == nil {
//			return result, nil
//		}
//		if err == redis.Nil {
//			return nil, nil
//		}
//		count++
//		if count > Repeated_Times {
//			log.Error("GetString失败 ", zap.Any("key", key), zap.Error(err))
//			return nil, err
//		} else {
//			log.Error(fmt.Sprintf("GetString请求key=%v数据失败，err = %s, "+
//				"发起第%v请求", key, err, count))
//			time.Sleep(Repeated_Interval * time.Millisecond) //重试间隔
//		}
//	}
//}
