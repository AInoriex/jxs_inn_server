package common

import (
	"math/rand"
	"time"
)

const (

)

// 随机种子
func GetRandSeed() *rand.Rand {
	// Create a new random number generator with a custom seed (e.g., current time)
	source := rand.NewSource(time.Now().UnixNano())
	rs := rand.New(source)
	return rs
}

// 获取随机数，x∈[min, max)
func RandomInt(min, max int) int {
	rs := GetRandSeed()
	return min + rs.Intn(max-min+1)	
}

// 获取随机元素
func GetRandomElement(slice []string) string {
	rs := GetRandSeed()
	index := rs.Intn(len(slice))
	return slice[index]	
}
