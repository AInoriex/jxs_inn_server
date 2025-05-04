package handler

import (
	"eshop_server/utils/log"
	"eshop_server/utils/qrcode"
	"eshop_server/utils/utime"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

func YltOrderHandler() {
	phone := "13292535169"
	password := "1234qwer"

	gt_token, cookie, err := YltUserLogin(phone, password)
	if err != nil {
		log.Error("YltOrderHandler登陆失败", zap.Error(err))
		return
	}
	time.Sleep(2 * time.Second)

	orderId, base64, err := YltCreateOrder(gt_token, cookie)
	if err != nil {
		log.Error("YltOrderHandler创建订单失败", zap.Error(err))
		return
	}
	log.Infof("YltOrderHandler创建订单成功", zap.String("订单ID", orderId), zap.String("支付二维码", base64))
	temp_filename := fmt.Sprintf("%v.jpg", utime.TimeToStrWindows(time.Now()))
	if err = qrcode.DecodeBase64ToImage(base64, temp_filename); err != nil {
		log.Error("YltOrderHandler生成二维码图片失败", zap.Error(err))
		return
	}
	defer os.Remove(temp_filename)
	fmt.Printf("YltOrderHandler二维码已生成在: %s ，请完成支付操作。", temp_filename)

	for {
		time.Sleep(3 * time.Second)
		payOk, err := YltCheckOrder(gt_token, cookie, orderId)
		if err != nil {
			log.Error("YltOrderHandler查询订单失败", zap.Error(err))
			continue
		}

		if !payOk {
			fmt.Println("未完成支付，等待下一轮查询...")
			continue
		} else {
			log.Info("YltOrderHandler查询订单已支付")
			break
		}
	}
	log.Success("YltOrderHandler用户已购买商品")
}
