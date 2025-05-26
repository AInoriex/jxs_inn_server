package handler

import (
	"eshop_server/src/utils/config"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/qrcode"
	"eshop_server/src/utils/utime"
	"fmt"
	"math/rand"
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
		log.Error("YltOrderHandler 登陆失败", zap.Error(err))
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

	// 重新登陆获取新token，并更新缓存
	gt_token, cookie, err := YltUserLogin(phone, password)
	if err != nil {
		log.Error("YltCreateOrderHandler 登陆失败", zap.Error(err))
		return "", "", err
	}
	time.Sleep(2 * time.Second)

	orderId, base64, err := YltCreateOrder(gt_token, cookie, productId, price)
	if err != nil {
		log.Error("YltCreateOrderHandler 创建订单失败", zap.Error(err))
		return "", "", err
	}
	log.Infof("YltCreateOrderHandler 创建订单成功", zap.String("订单ID", orderId), zap.String("支付二维码", base64))
	return orderId, base64, err
}

// @Title		获取YLT代理账号
// @Description	读取配置文件随机获取账号
func GetYltRandomAgent() (string, string, error) {
	accounts := config.CommonConfig.YltAccount
	rand.NewSource(utime.GetNowUnix())
	index := rand.Intn(len(accounts)) // 生成随机索引
	for phone, password := range accounts {
		if index == 0 {
			return phone, password, nil
		}
		index--
	}
	return "", "", fmt.Errorf("索引:%v没有找到有效的账号", index)
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

// @Title		获取YLT代理账号token, cookie
// @Description	从缓存获取账号对应token, cookie
func GetYltAgentTokenCookie(phone string) (string, string, error) {
	// 从Redis缓存获取gt_token, cookie
	// HARDCODE 目前从配置文件读取
	// TODO 从redis读取&更新
	returnMap := map[string]map[string]string{
		"17803152032": {
			"gt_token": "",
			"cookie": "",
		},
		"16762397443": {
			"gt_token": "",
			"cookie": "",
		},
		"13292535169": {
			"gt_token": "",
			"cookie": "",
		},
		"17118521660": {
			"gt_token": "",
			"cookie": "",
		},
		"17054870744": {
			"gt_token": "",
			"cookie": "",
		},
	}
	if _, ok := returnMap[phone]; !ok {
		return "", "", fmt.Errorf("账号:%v没有找到对应的token, cookie", phone)	
	}
	return returnMap[phone]["gt_token"], returnMap[phone]["cookie"], nil
}
