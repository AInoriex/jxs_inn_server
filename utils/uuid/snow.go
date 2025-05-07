package uuid

import (
	"fmt"
	"sync"
	"time"
)

const (
	// 雪花算法的起始时间戳（自定义基准时间）
	snowStartTime = 1609459200000 // 2021-01-01 00:00:00 UTC 的毫秒时间戳

	// 位移量
	timeStampShift    = 22
	dataCenterIDShift = 17 // 占位5位，取值范围 [0, 2^5-1]
	machineIDShift    = 12 // 占位5位，取值范围 [0, 2^5-1]
	sequenceMask      = 0xFFF // 12位序列号掩码
)

// Snowflake 结构体
type Snowflake struct {
	// 开始时间戳，用于计算时间差
	startTime int64
	// 数据中心ID（5位）
	dataCenterID int64
	// 机器ID（5位）
	machineID int64
	// 序列号（12位）
	sequence int64
	// 锁，用于并发控制
	lock sync.Mutex
}

// NewSnowflake 创建一个新的雪花算法实例
func NewSnowflake(dataCenterID, machineID int64) *Snowflake {
	return &Snowflake{
		startTime:    snowStartTime, // 2021-01-01 00:00:00 UTC 的毫秒时间戳
		dataCenterID: dataCenterID,
		machineID:    machineID,
		sequence:     0,
	}
}

// Generate 生成一个雪花ID
func (s *Snowflake) Generate() (int64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// 获取当前时间的毫秒级时间戳
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)

	// 如果当前时间小于上次时间戳，说明发生了时钟回拨
	if currentTime < s.startTime {
		return 0, fmt.Errorf("clock moved backwards. Refusing to generate id for %d milliseconds", s.startTime-currentTime)
	}

	// 如果是同一毫秒的生成，序列号自增
	if currentTime == s.startTime {
		s.sequence = (s.sequence + 1) & sequenceMask // 12位序列号
		// 如果序列号溢出，等待下一毫秒
		if s.sequence == 0 {
			currentTime = s.tilNextMillis(s.startTime)
		}
	} else {
		// 重置序列号
		s.sequence = 0
	}

	// 更新上次时间戳
	s.startTime = currentTime

	// 按位拼接生成ID
	id := ((currentTime - 1609459200000) << timeStampShift) | // 时间戳部分（41位）
		(s.dataCenterID << dataCenterIDShift) | // 数据中心ID部分（5位）
		(s.machineID << machineIDShift) | // 机器ID部分（5位）
		s.sequence // 序列号部分（12位）

	return id, nil
}

// tilNextMillis 等待下一毫秒，直到时间戳改变
func (s *Snowflake) tilNextMillis(lastTime int64) int64 {
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	for timeStamp <= lastTime {
		timeStamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timeStamp
}

// ParseSnowflakeID 反序列化雪花ID，解析出时间戳、数据中心ID、机器ID和序列号
func ParseSnowflakeID(id int64) (int64, int64, int64, int64) {
	// 提取时间戳部分
	timestamp := (id >> timeStampShift) + snowStartTime
	// 提取数据中心ID部分
	dataCenterID := (id >> dataCenterIDShift) & 0x1F // 5位
	// 提取机器ID部分
	machineID := (id >> machineIDShift) & 0x1F // 5位
	// 提取序列号部分
	sequence := id & sequenceMask // 12位

	return timestamp, dataCenterID, machineID, sequence
}

func exampleGenerateSnowIds() {
	// 创建一个雪花算法实例，数据中心ID为1，机器ID为1
	snowflake := NewSnowflake(1, 1)

	// 生成10个雪花ID
	for i := 0; i < 10; i++ {
		id, err := snowflake.Generate()
		if err != nil {
			fmt.Println("Error generating ID:", err)
		} else {
			fmt.Println("Generated ID:", id)
		}
	}
}

func exampleParseSnowIds() {
	// 示例雪花ID
	id := int64(4611686018427387905)

	// 解析雪花ID
	t, dataCenterID, machineID, sequence := ParseSnowflakeID(id)

	// 打印解析结果
	fmt.Printf("ID: %d\n", id)
	fmt.Printf("Timestamp: %s\n", t)
	fmt.Printf("Data Center ID: %d\n", dataCenterID)
	fmt.Printf("Machine ID: %d\n", machineID)
	fmt.Printf("Sequence: %d\n", sequence)
}
