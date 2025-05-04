package alarm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// 飞书告警文本信息
func PostFeiShu(level, url string, info string) error {
	msg := FeiShuMsg{
		MsgType: "text",
	}
	msg.Content.Text = fmt.Sprintf("[%s] %s", level, info)
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}
