package handler

import (
	"errors"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	"eshop_server/utils/log"
	"eshop_server/utils/uuid"
	"go.uber.org/zap"
)

const (
	YltAccount1 = "10000000000000000000" // YLT测试账号1
	YltAccount2 = "10000000000000000001" // YLT测试账号2
)

// @Author	AInoriex
// @Desc	扫码支付流程
// HARDNEED	当前仅支持YLT支付
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

	// 获取YLT账号
	phone, password, err := GetYltRandomAccount()
	if err != nil || phone == "" || password == "" {
		log.Error("QrcodeOrderPaymentHandler 获取YLT账号失败", zap.Error(err))
		return "", errors.New("创建订单失败，请联系客服")
	}

	// 获取YLT关联商品信息

	// 调用接口创建YLT订单
	yltOrderId, qrcode, err := YltCreateOrderHandler(phone, password, product.ExternalId, product.Price)
	if err != nil || yltOrderId == "" || qrcode == "" {
		log.Error("QrcodeOrderPaymentHandler 创建YLT订单失败", zap.String("ylt订单id", yltOrderId), zap.Bool("二维码是否为空", (qrcode == "")), zap.Error(err))
		return "", errors.New("创建订单失败，请咨询客服")
	}

	// 创建payment
	paymentId := uuid.GetUuid()
	payment := &model.Payment{
		Id:          paymentId,
		OrderID:     order.Id,                   // 订单ID
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
