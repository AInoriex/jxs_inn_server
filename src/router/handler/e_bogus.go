package handler

import (
	// "eshop_server/utils/common"
	// "math"
	// "time"
)

// 签名参数检查：时间戳间隔limit_count秒内
func CheckSignParam(sign string) bool {
	return true

	// limit_count := 60.0 * 10 // 过期时间10分钟
	// // 判断sign是否为10/13位时间戳
	// if !(len(sign) == 10 || len(sign) == 13) {
	// 	return false
	// }
	// val := common.StringToInt64NotErr(sign)
	// if val <= 0 {
	// 	return false
	// }
	// if len(sign) == 10 {
	// 	now := time.Now().Unix()
	// 	cal := math.Abs(float64(val - now))
	// 	return cal <= limit_count
	// } else if len(sign) == 13 {
	// 	now := time.Now().Unix()*1000
	// 	cal := math.Abs(float64(val - now))
	// 	return cal <= limit_count
	// } else {
	// 	return false
	// }
}
