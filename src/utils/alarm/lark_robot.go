package alarm

import (
	"bytes"
	"encoding/json"
	"eshop_server/src/utils/utime"
	"fmt"
	"net/http"
)

// 飞书告警文本信息
func PostFeiShu(level, url string, info string) error {
	var retryCount int = 3
	msg := FeiShuMsg{
		MsgType: "text",
	}
	msg.Content.Text = fmt.Sprintf("[%s] %s", level, info)
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 包含初始请求的总重试循环（初始请求+retryCount次重试）
	for i := 0; i <= retryCount; i++ {
		req, err := http.NewRequest("POST", url, bytes.NewReader(marshal))
		if err != nil {
			return err // 请求创建失败无法重试，直接返回错误
		}
		resp, err := http.DefaultClient.Do(req)
		defer func() {
			if resp != nil { _ = resp.Body.Close() }
		}()
		// 请求成功
		if err == nil {
			return nil
		}
		// 请求失败尝试重试
		if i < retryCount {
			utime.RandomSleep(3, 5)
		}
	}
	return fmt.Errorf("请求失败，已尝试%d次", retryCount + 1)
}
