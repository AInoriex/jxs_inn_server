package handler

import (
	"eshop_server/src/common/cache"
	router_handler "eshop_server/src/router/handler"
	"eshop_server/src/utils/config"
	"eshop_server/src/utils/log"
	"time"
)

// 轮询YLT账号登录状态
func YltLoginCronjob() {
	ylt_accounts := config.CommonConfig.YltAccount
	if len(ylt_accounts) <= 0 {
		log.Warnf("YltLoginCronjob 未配置YLT账号列表")
		return
	}

	for phone, password := range ylt_accounts {
		time.Sleep(time.Second * 2)
		// 检查账号缓存
		flag, _, _ := cache.GetYltUserToken(phone)
		if flag {
			log.Infof("YltLoginCronjob 账号 %s 已缓存", phone)
			continue
		}

		// 登录账号
		gt_token, cookie, err := router_handler.YltUserLogin(phone, password)
		if err != nil {
			log.Errorf("YltCreateOrderHandler 登录失败, phone:%s, password:%s, error:%v", phone, password, err)
		}

		// 缓存存储账号登录状态
		err = cache.SaveYltUserToken(phone, gt_token, cookie)
		if err != nil {
			log.Errorf("YltCreateOrderHandler 缓存账号登录状态失败, phone:%s, password:%s, error:%v", phone, password, err)
		}
	}
}
