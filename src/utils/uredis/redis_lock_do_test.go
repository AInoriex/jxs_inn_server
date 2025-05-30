package uredis

import (
	"fmt"
	"testing"
	"time"
)

func TestRedisLock_GetString(t *testing.T) {
	client := New("127.0.0.1:6379", "", 0)
	key := "lock"

	SetString(client, key, "123")

	rl := NewRedisLock(key, client)
	go func() {
		b, err := rl.GetString()
		if err != nil {
			t.Error("cuowu", err)
		}
		fmt.Printf("str=%s, err=%v \n", b, err)
	}()

	time.Sleep(time.Duration(100) * time.Millisecond) //锁住
	b, err := rl.GetString()
	if err == nil {
		t.Error(err)
	}
	fmt.Printf("str=%s, err=%v \n", b, err)

	rl.SetString("456") //设置并解锁

	b, err = rl.GetString()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("str=%s, err=%v \n", b, err)
}
