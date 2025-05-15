package alarm

import (
	// "time"
)

// 文本信息
type FeiShuMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

var (
	AlarmLevelPanic string = "PANIC" // 飞书告警-宕机级别
	AlarmLevelError string = "ERROR" // 飞书告警-错误级别
	AlarmLevelWarn  string = "WARN"  // 飞书告警-提醒级别
	AlarmLevelInfo  string = "INFO"  // 飞书告警-信息级别
)
