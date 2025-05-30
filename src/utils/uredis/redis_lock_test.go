package uredis

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestRedisLock_Lock(t *testing.T) {
	client := New("127.0.0.1:6379", "", 0)
	key := "lock"
	rl := NewRedisLock(key, client)
	go func() {
		err := rl.LockNoWait()
		if err != nil {
			log.Print("cuowu", err)
		}
		fmt.Println("111", err)
	}()

	time.Sleep(1 * time.Second) //锁住

	err := rl.LockNoWait() //此时 会获取新的lock Block
	fmt.Println("222", err)

	time.Sleep(2 * time.Second) //解锁

	err = rl.LockNoWait()
	fmt.Println("333", err)

	rl.Unlock()

	err = rl.LockNoWait() //成功
	fmt.Println("444", err)

	time.Sleep(5 * time.Second)
}

func TestRedisLock_LockWait(t *testing.T) {
	client := New("127.0.0.1:6379", "", 0)
	key := "lock"
	rl := NewRedisLock(key, client)
	go func() {
		err := rl.Lock()
		if err != nil {
			log.Print("cuowu", err)
		}
		fmt.Println("111", err)
	}()

	time.Sleep(time.Duration(100) * time.Millisecond) //锁住

	err := rl.Lock() //此时 会获取新的lock Block
	fmt.Println("222", err)

	time.Sleep(1 * time.Second) //解锁

	err = rl.Lock()
	fmt.Println("333", err)

	rl.Unlock()

	err = rl.Lock() //成功
	fmt.Println("444", err)

	time.Sleep(5 * time.Second)
}
