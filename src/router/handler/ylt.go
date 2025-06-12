package handler

import (
	"eshop_server/src/common/cache"
	"eshop_server/src/utils/common"
	"eshop_server/src/utils/config"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/qrcode"
	"eshop_server/src/utils/utime"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// @Title	本地调试使用-YLT完整下单付款流程
func YltOrderFullHandler(phone string, password string) {
	// 指定商品 https://yuanlitui.com/a/ar55
	productId, customPrice := string("5517"), float64(0.5)
	// productId, customPrice := string("5517"), float64(1.5) // 逆天bug，能自定义金额生成收款码

	gt_token, cookie, err := YltUserLogin(phone, password)
	if err != nil {
		log.Error("YltOrderHandler 登录失败", zap.Error(err))
		return
	}
	time.Sleep(2 * time.Second)

	orderId, base64, err := YltCreateOrder(gt_token, cookie, productId, customPrice)
	if err != nil {
		log.Error("YltOrderHandler 创建订单失败", zap.Error(err))
		return
	}
	log.Infof("YltOrderHandler 创建订单成功", zap.String("订单ID", orderId), zap.String("支付二维码", base64))
	temp_filename := fmt.Sprintf("%v.jpg", utime.TimeToStrWindows(time.Now()))
	if err = qrcode.DecodeBase64ToImage(base64, temp_filename); err != nil {
		log.Error("YltOrderHandler 生成二维码图片失败", zap.Error(err))
		return
	}
	defer os.Remove(temp_filename)
	fmt.Printf("YltOrderHandler 二维码已生成在: %s ，请完成支付操作。", temp_filename)

	for {
		time.Sleep(3 * time.Second)
		payOk, err := YltCheckOrder(gt_token, cookie, orderId)
		if err != nil {
			log.Error("YltOrderHandler 查询订单失败", zap.Error(err))
			continue
		}

		if !payOk {
			fmt.Println("未完成支付，等待下一轮查询...")
			continue
		} else {
			log.Info("YltOrderHandler 查询订单已支付")
			break
		}
	}
	log.Success("YltOrderHandler 用户已购买商品")
}

// @Title		创建YLT订单
// @Description	创建订单并返回订单ID和支付二维码
// @Param		phone		手机号
// @Param		password	密码
// @Param		productId	商品ID
// @Param		price		自定义价格
// @Return		orderId		订单ID
// @Return		base64		支付二维码
// @Return		err			错误信息
func YltCreateOrderHandler(phone string, password string, productId string, price float64) (string, string, error) {
	log.Infof("YltCreateOrderHandler 开始创建订单: phone: %s, password: %s, productId: %s", phone, password, productId)
	// 尝试从缓存获取gt_token, cookie
	flag, gt_token, cookie := cache.GetYltUserToken(phone)
	if !flag {
		log.Infof("YltCreateOrderHandler 从缓存获取gt_token, cookie失败")
		// 重新登录获取新token，并更新缓存
		// gt_token, cookie, err := YltUserLogin(phone, password)
		// if err != nil {
		// 	log.Error("YltCreateOrderHandler 登录失败", zap.Error(err))
		// 	return "", "", err
		// }
		return "", "", fmt.Errorf("从缓存获取%s账号登陆Token信息失败", phone)
	}
	log.Infof("YltCreateOrderHandler 从缓存获取YLT用户Token信息成功, gt_token: %s, cookie: %s", gt_token, cookie)
	orderId, base64, err := YltCreateOrder(gt_token, cookie, productId, price)
	if err != nil {
		log.Errorf("YltCreateOrderHandler 创建YLT订单失败, phone: %s, password: %s, productId: %s, error: %v", phone, password, productId, err)
		return "", "", err
	}
	log.Infof("YltCreateOrderHandler 创建YLT订单成功, orderId: %s, qrcode:%s", orderId, base64)
	return orderId, base64, err
}

// @Title		获取YLT代理账号
// @Description	读取配置文件随机获取账号
func GetYltConfigRandomAccount() (string, string, error) {
	accounts := config.CommonConfig.YltAccount
	randint := common.RandomInt(0, len(accounts))
	for phone, password := range accounts {
		if randint == 0 {
			return phone, password, nil	
		}	
		randint--
	}
	return "", "", fmt.Errorf("GetYltConfigRandomAccount failed")
}

// @Title		获取YLT代理账号密码
// @Description	读取配置文件获取账号对应密码
func GetYltAgentPassword(phone string) (string, error) {
	accounts := config.CommonConfig.YltAccount
	password, ok := accounts[phone]
	if !ok {
		return "", fmt.Errorf("账号:%v没有找到对应的密码", phone)
	}
	return password, nil
}
