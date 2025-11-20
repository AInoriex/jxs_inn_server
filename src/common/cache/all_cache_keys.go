package cache

import (
	"fmt"
)

const (
	// jxs用户登录态
	KeyJxsUserToken         string = "JxsUser:%v" // userId
	KeyJxsUserTokenTimeout         = 30 * 60      // 用户Token有效时长30分钟
	KeyJxsAdminTokenTimeout        = 120 * 60     // 后台用户Token有效时长120分钟

	// jxs邮箱验证
	KeyJxsVerifyMailCode          string = "JxsVEmailCode:%v:%v" // ip:toEmail
	KeyJxsVerifyMailCodeMinsLimit        = 5
	KeyJxsVerifyMailCodeTimeout          = KeyJxsVerifyMailCodeMinsLimit * 60 // 邮箱验证码有效时长5分钟

	// ylt登录态
	KeyYltUserPrefix       string = "YltUser"
	KeyYltUserToken        string = KeyYltUserPrefix + ":%v" // phone
	KeyYltUserTokenTimeout        = 3 * 60 * 60              // YLT Token有效时长3小时
)

// jxs用户登录态Key
func GetJxsUserTokenKey(userId string) string {
	return fmt.Sprintf(KeyJxsUserToken, userId)
}

// jxs邮箱验证Key
func GetJxsVerifyMailCodeKey(ip string, toEmail string) string {
	return fmt.Sprintf(KeyJxsVerifyMailCode, ip, toEmail)
}

// ylt用户登录态Key
func GetYltUserTokenKey(phone string) string {
	return fmt.Sprintf(KeyYltUserToken, phone)
}
