package cache

import (
	"fmt"
)

const (
	// jxs用户登录态
	KeyJxsUserToken        = "JxsUser:%v"
	KeyJxsUserTokenTimeout = 30 * 60 // 用户Token有效时长30mins

	// ylt登录态
	KeyYltUserToken        = "YltUser:%v"
	KeyYltUserTokenTimeout = 3 * 60 * 60 // 用户Token有效时长1hours
)

// jxs用户登录态Key
func GetJxsUserTokenKey(userId string) string {
	return fmt.Sprintf(KeyJxsUserToken, userId)
}

// ylt用户登录态Key
func GetYltUserTokenKey(phone string) string {
	return fmt.Sprintf(KeyYltUserToken, phone)
}