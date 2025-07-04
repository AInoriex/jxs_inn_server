package cache

import (
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uredis"
)

// 获取jxs用户登录态
func GetJxsUserToken(userId string) (bool, string) {
	key := GetJxsUserTokenKey(userId)
	b, err := uredis.GetString(uredis.RedisCon, key)
	if err != nil {
		//log.Errorf("GetJxsUserToken redis错误", zap.Any("data", b), zap.Error(err))
		return false, ""
	} else {
		if b == nil {
			//log.Errorf("GetJxsUserToken redis为空", zap.Any("data", b), zap.Error(err))
			return false, ""
		}
		return true, string(b)
	}
}

// 保存jxs用户登录态
func SaveJxsUserToken(userId string, token string) error {
	key := GetJxsUserTokenKey(userId)
	err := uredis.SetString(uredis.RedisCon, key, token, KeyJxsUserTokenTimeout)
	log.Debugf("SaveJxsUserToken params, userId:%s, err:%v", userId, err)
	return err
}

// 删除jxs用户登录态
func DelJxsUserToken(userId string) bool {
	key := GetJxsUserTokenKey(userId)
	err := uredis.DelKey(uredis.RedisCon, key)
	log.Debugf("DelJxsUserToken params, userId:%s, err:%v", userId, err)
	return err == nil
}

// 获取jxs邮箱验证
func GetJxsVerifyMailCode(ip string, toEmail string) (bool, string) {
	key := GetJxsVerifyMailCodeKey(ip, toEmail)
	b, err := uredis.GetString(uredis.RedisCon, key)
	if err != nil {
		//log.Errorf("GetJxsVerifyMailCode redis错误", zap.Any("data", b), zap.Error(err))
		return false, ""
	} else {
		if b == nil {
			//log.Errorf("GetJxsVerifyMailCode redis为空", zap.Any("data", b), zap.Error(err))
			return false, ""
		}
		return true, string(b)
	}
}

// 保存jxs邮箱验证
func SaveJxsVerifyMailCode(ip string, toEmail string, code string) error {
	key := GetJxsVerifyMailCodeKey(ip, toEmail)
	err := uredis.SetString(uredis.RedisCon, key, code, KeyJxsVerifyMailCodeTimeout)
	log.Debugf("SaveJxsVerifyMailCode params, ip:%s, toEmail:%s, err:%v", ip, toEmail, err)
	return err
}

// 删除jxs邮箱验证
func DelJxsVerifyMailCode(ip string, toEmail string) bool {
	key := GetJxsVerifyMailCodeKey(ip, toEmail)
	err := uredis.DelKey(uredis.RedisCon, key)
	log.Debugf("DelJxsVerifyMailCode params, ip:%s, toEmail:%s, err:%v", ip, toEmail, err)
	return err == nil
}
