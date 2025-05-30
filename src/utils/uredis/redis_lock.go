package uredis

import (
	"context"
	"encoding/base64"
	"errors"
	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"eshop_server/src/utils/log"
	"math/rand"
	"time"
)

const (
	Lock_Key_Pre = "redis_lock:"
	Lock_Expire  = 2000 //2000毫秒

	Err_Is_Exist  = "get lock fail"
	Lock_Interval = 5  //毫秒
	Lock_Repeated = 20 //重复次数
)

//单点redis，多点注意
type RedisLock struct {
	lockKey   string
	lockValue string
	Key       string
	Field     string
	Expire    int64 //毫秒
	con       *redis.Client
	seq       int
	start     time.Time
	end       time.Time
}

//保证原子性（redis是单线程），避免del删除了，其他client获得的lock
var delScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)

// 第三个参数，单位是毫秒
func NewRedisLock(key string, con *redis.Client, lockExpire ...int64) *RedisLock {
	lock := new(RedisLock)
	lock.lockKey = Lock_Key_Pre + key
	lock.Key = key
	lock.con = con
	lock.seq = rand.Int()
	b := make([]byte, 16)
	rand.Read(b)
	lock.lockValue = base64.StdEncoding.EncodeToString(b)

	lock.Expire = Lock_Expire
	if len(lockExpire) > 0 {
		if lockExpire[0] > 0 {
			lock.Expire = lockExpire[0]
		}
	}
	return lock
}

//加锁 获取锁失败马上返回
func (this *RedisLock) LockNoWait() error {
	lockReply, err := this.con.SetNX(context.Background(), this.lockKey, this.lockValue, time.Duration(this.Expire)*time.Millisecond).Result()
	if err != nil {
		return errors.New("redis fail")
	}
	if !lockReply {
		return errors.New(Err_Is_Exist)
	}
	return nil
}

//等待锁 等待50ms
func (this *RedisLock) Lock() error {
	this.start = time.Now()
	var b = false
	for i := 0; i < Lock_Repeated; i++ {
		err := this.LockNoWait()
		if err != nil {
			if err.Error() != Err_Is_Exist {
				return err
			}
			time.Sleep(Lock_Interval * time.Millisecond) //重试间隔
		} else {
			b = true
			break
		}
	}
	if !b {
		return errors.New(Err_Is_Exist)
	}
	return nil
}

//等待锁 循环一次5ms
func (this *RedisLock) LockNum(num int) error {
	this.start = time.Now()
	var b = false
	for i := 0; i < num; i++ {
		err := this.LockNoWait()
		if err != nil {
			if err.Error() != Err_Is_Exist {
				return err
			}
			time.Sleep(Lock_Interval * time.Millisecond) //重试间隔
		} else {
			b = true
			break
		}
	}
	if !b {
		return errors.New(Err_Is_Exist)
	}
	return nil
}

//解锁
func (this *RedisLock) Unlock() error {
	this.end = time.Now()
	log.Info("redis 锁序列号：", zap.Int("this.seq", this.seq), zap.Int64("开始时间：", this.start.UnixNano()), zap.Int64("结束时间", this.end.UnixNano()), zap.Any("共耗时", time.Since(this.start)))
	return delScript.Run(context.Background(), this.con, []string{this.lockKey}, this.lockValue).Err()
}

// 分布式环境下的一次性任务
func SerializeExecDelay(uniqueTaskName string, cli *redis.Client, do func()) func() {
	return func() {
		lock := NewRedisLock(uniqueTaskName, cli, 5000)
		err := lock.Lock()
		if err != nil {
			log.Info("SerializeExecDelay lock err", zap.Error(err), zap.String("name", uniqueTaskName))
			return
		}
		log.Info("SerializeExecDelay lock OK", zap.Error(err), zap.String("name", uniqueTaskName))
		defer lock.Unlock()
		do()
		time.Sleep(time.Second * 2) // 故意延长时间，避免完成太快被其他节点抢到锁
	}
}

func SerializeExec(uniqueTaskName string, cli *redis.Client, do func()) func() {
	return func() {
		lock := NewRedisLock(uniqueTaskName, cli, 5000)
		err := lock.Lock()
		if err != nil {
			log.Info("SerializeExecDelay lock err", zap.Error(err), zap.String("name", uniqueTaskName))
			return
		}
		log.Info("SerializeExecDelay lock OK", zap.Error(err), zap.String("name", uniqueTaskName))
		defer lock.Unlock()
		do()
	}
}
