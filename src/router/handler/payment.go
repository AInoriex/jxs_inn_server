package handler

import (
	"errors"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uuid"
	"go.uber.org/zap"
)

// @Author	AInoriex
// @Desc	扫码支付流程
// @HARDNEED	当前仅支持YLT支付
func QrcodeOrderPaymentHandler(reqbody model.CreateOrderReq, order *model.Order, product *model.Products) (qrcode string, err error) {
	// 参数判断
	if reqbody.PaymentMethod != model.PaymentMethodQrcode {
		log.Error("QrcodeOrderPaymentHandler 支付方式无效:非扫码支付", zap.String("payment_method", reqbody.PaymentMethod))
		return "", errors.New("参数错误：支付方式无效")
	}
	if reqbody.PaymentGatewayType != model.PaymentGatewayTypeYlt {
		log.Error("QrcodeOrderPaymentHandler 支付网关类别无效:非原力通", zap.Int32("payment_gateway_type", reqbody.PaymentGatewayType))
		return "", errors.New("参数错误：支付网关无效")
	}
	if product.ExternalId == "" {
		log.Error("QrcodeOrderPaymentHandler 商品关联ID为空")
		return "", errors.New("参数错误：商品关联ID为空")
	}

	// 创建YLT订单
	var yltOrderId, phone, password string
	var retry, retryLimit = 3, 3
	for {
		if retry <= 0 {
			log.Error("QrcodeOrderPaymentHandler 创建YLT订单重试失败")
			return "", errors.New("创建订单失败，请联系客服")
		} else if retry < retryLimit {
			log.Infof("QrcodeOrderPaymentHandler 创建YLT订单当前重试次数:retry:%v", retry)
		}

		// 随机获取YLT账号
		phone, password, err = GetYltConfigRandomAccount()
		if err != nil || phone == "" || password == "" {
			log.Errorf("QrcodeOrderPaymentHandler 获取YLT账号失败, phone:%v, password:%v, error:%v", phone, password, err)
			continue
		}

		// 调用接口创建YLT订单
		yltOrderId, qrcode, err = YltCreateOrderHandler(phone, password, product.ExternalId, product.Price)
		if err != nil || yltOrderId == "" || qrcode == "" {
			log.Errorf("QrcodeOrderPaymentHandler 创建YLT订单失败, yltOrderId:%v, qrcode is null?:%v, error:%v", yltOrderId, (qrcode == ""), err)
			retry--
			continue
		} else {
			log.Infof("QrcodeOrderPaymentHandler 创建YLT订单成功, yltOrderId:%v, qrcode is null?:%v", yltOrderId, (qrcode == ""))
			break
		}
	}

	// 创建payment
	paymentId := uuid.GetUuid()
	payment := &model.Payment{
		Id:          paymentId,
		OrderId:     order.Id,                   // 订单ID
		FinalAmount: order.FinalAmount,          // 订单金额
		GatewayType: reqbody.PaymentGatewayType, // 支付网关类别
		Method:      model.PaymentMethodQrcode,  // 支付方式
		Status:      model.PaymentStatusPaying,  // 支付状态
		GatewayID:   yltOrderId,                 // YLT订单ID
		Agent:       phone,                      // 支付代理账号
	}
	payment, err = dao.CreatePayment(payment)
	if err != nil {
		log.Error("QrcodeOrderPaymentHandler 创建支付记录失败", zap.Error(err))
		return "", errors.New("创建支付失败")
	}

	// 更新order表信息
	order.PaymentId = paymentId
	order.PaymentStatus = payment.Status
	if _, err = dao.UpdateOrderByField(order, []string{"payment_id", "payment_status"}); err != nil {
		log.Error("QrcodeOrderPaymentHandler 更新订单信息失败", zap.Error(err))
		return "", errors.New("更新订单信息失败")
	}

	return qrcode, nil
}
