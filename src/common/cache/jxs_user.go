package cache

import (
	"eshop_server/src/utils/common"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uredis"
)

// 获取jxs用户登陆态
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

// 保存jxs用户登陆态
func SaveJxsUserToken(userId string, code int32) error {
	key := GetJxsUserTokenKey(userId)
	err := uredis.SetString(uredis.RedisCon, key, common.Int32ToString(code), KeyJxsUserTokenTimeout)
	log.Infof("SaveJxsUserToken params, userId:%s, err:%s", userId, err.Error())
	return err
}

// 删除jxs用户登陆态
func DelJxsUserToken(userId string) bool {
	key := GetJxsUserTokenKey(userId)
	err := uredis.DelKey(uredis.RedisCon, key)
	log.Debugf("DelJxsUserToken params, userId:%s, err:%s", userId, err.Error())
	return err == nil
}
