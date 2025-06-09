package cache

import (
	"fmt"
)

const (
	// jxs用户登陆态
	KeyJxsUserToken        = "JxsUser:%v"
	KeyJxsUserTokenTimeout = 0.5 * 60 * 60 // 0.5小时
)

// jxs用户登陆态Key
func GetJxsUserTokenKey(userId string) string {
	return fmt.Sprintf(KeyJxsUserToken, userId)
}
