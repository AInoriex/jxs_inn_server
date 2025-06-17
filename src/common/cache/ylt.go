package cache

import (
	"encoding/json"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uredis"
	"eshop_server/src/utils/common"
)

// YLT用户登录态信息缓存结构
type YltUserToken struct {
	GtToken string `json:"gt_token"`
	Cookie  string `json:"cookie"`
}

// 打包YLT用户登录态信息数据
func MarshalYltUserToken(gt_token string, cookie string) (string, error) {
	tokenPack := YltUserToken{
		GtToken: gt_token,
		Cookie:  cookie,
	}
	rawBytes, err := json.Marshal(tokenPack)
	if err != nil {
		return "", err
	}
	return string(rawBytes), nil
}

// 解包YLT用户登录态信息数据
func UnmarshalYltUserToken(rawBytes string) (tokenPack YltUserToken, err error) {
	err = json.Unmarshal([]byte(rawBytes), &tokenPack)
	return tokenPack, err
}

// 根据账号获取YLT用户登录态信息
// @Return	flag, gt_token, cookie
func GetYltUserToken(phone string) (bool, string, string) {
	key := GetYltUserTokenKey(phone)
	rawBytes, err := uredis.GetString(uredis.RedisCon, key)
	if err != nil {
		//log.Errorf("GetYltUserToken redis错误", zap.Any("data", b), zap.Error(err))
		return false, "", ""
	} else {
		if rawBytes == nil {
			//log.Errorf("GetYltUserToken redis为空", zap.Any("data", b), zap.Error(err))
			return false, "", ""
		}
		token, err := UnmarshalYltUserToken(string(rawBytes))
		if err != nil {
			log.Errorf("GetYltUserToken UnmarshalYltUserToken错误, rawBytes:%s, err:%v", rawBytes, err)
			return false, "", ""
		}
		return true, token.GtToken, token.Cookie
	}
}

// 随机获取YLT用户登录态信息
// @Return	flag, phone, gt_token, cookie
func GetRandomYltUserToken() (bool, string, string, string) {
	keys, err := uredis.GetAllKey(uredis.RedisCon, KeyYltUserPrefix+"*")
	if err != nil {
		log.Errorf("GetRandomYltUserToken GetAllKey错误, err:%v", err)
		return false, "", "", ""
	} 
	if len(keys) <= 0 {
		log.Errorf("GetRandomYltUserToken 未找到YLT用户登录态信息")
		return false, "", "", ""
	}
	// 随机取keys中的一个
	randKey := common.GetRandomElement(keys)
	// 字符串处理 YltUser:17803152032 -> 17803152032
	phone := randKey[len(KeyYltUserPrefix)+1:]
	flag, gt_token, cookie := GetYltUserToken(phone)
	log.Debugf("GetRandomYltUserToken params, phone:%s, gt_token:%s, cookie:%s", phone, gt_token, cookie)
	return flag, phone, gt_token, cookie
}

// 保存YLT用户登录态
func SaveYltUserToken(phone string, gt_token string, cookie string) error {
	key := GetYltUserTokenKey(phone)
	rawBytes, err := MarshalYltUserToken(gt_token, cookie)
	if err != nil {
		log.Errorf("SaveYltUserToken MarshalYltUserToken错误, phone:%s, err:%v", phone, err)
		return err
	}
	// 随机偏移时间 1-10 mins，防止下次刷新token频繁登录
	randOffset := int64(common.RandomInt(60, 600))
	err = uredis.SetString(uredis.RedisCon, key, string(rawBytes), KeyYltUserTokenTimeout+randOffset)
	log.Debugf("SaveYltUserToken params, phone:%s, err:%v", phone, err)
	return err
}

// 删除YLT用户登录态
func DelYltUserToken(phone string) bool {
	key := GetYltUserTokenKey(phone)
	err := uredis.DelKey(uredis.RedisCon, key)
	log.Debugf("DelYltUserToken params, phone:%s, err:%v", phone, err)
	return err == nil
}
