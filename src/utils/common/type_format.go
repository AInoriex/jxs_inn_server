package common

import (
	"eshop_server/src/utils/log"
	"go.uber.org/zap"
	"sort"
	"strconv"
)

// 字符串转int64
func StringToInt32NotErr(numString string) int32 {
	num := StringToInt64NotErr(numString)
	return int32(num)
}

// 字符串转int64
func StringToInt64NotErr(numString string) int64 {
	num, err := strconv.ParseInt(numString, 10, 64)
	if err != nil {
		log.Error("字符串转int64", zap.Any("numString", numString))
		return 0
	}
	return num
}

// int32转字符串
func Int32ToString(num int32) string {
	return Int64ToString(int64(num))
}

// int64转字符串
func Int64ToString(num int64) string {
	numString := strconv.FormatInt(num, 10)
	return numString
}

// 重新包装code
func RespCodeToInt32(code interface{}) int32 {
	switch code.(type) {
	case int:
		return int32(code.(int))
	case string:
		data, err := strconv.ParseInt(code.(string), 10, 32)
		if err == nil {
			return int32(data)
		}
		return int32(-1)
	case int32:
		return int32(code.(int32))
	case int64:
		return int32(code.(int64))
	case float64:
		return int32(code.(float64))
	case float32:
		return int32(code.(float32))
	default:
		log.Error("is an unknown type.", zap.Any("code", code))
		return int32(-1)
	}
}

func StringInSlice(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	//index的取值：[0,len(str_array)]
	if index < len(strArray) && strArray[index] == target { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}
