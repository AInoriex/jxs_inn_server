package handler

import (
	"errors"
	"eshop_server/dao"
	"eshop_server/model"
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
func QrcodeOrderPaymentHandler(reqbody model.CreateOrderReq, order model.Order) (qrcode string, err error) {
	// 参数判断
	if reqbody.PaymentMethod != model.PaymentMethodQrcode {
		log.Error("QrcodeOrderPaymentHandler 支付方式无效:非扫码支付", zap.String("payment_method", reqbody.PaymentMethod))
		return "", errors.New("")
	}
	if reqbody.PaymentGatewayType != model.PaymentGatewayTypeYlt {
		log.Error("QrcodeOrderPaymentHandler 支付网关类别无效:非原力通", zap.Int32("payment_gateway_type", reqbody.PaymentGatewayType))
		return "", errors.New("")
	}

	// 获取YLT账号
	phone, password := "", ""

	// 调用接口创建YLT订单
	yltOrderId, qrcode, err := YltCreateOrderHandler(phone, password)

	// 创建payment
	paymentId := uuid.GetUuid()
	payment := &model.Payment{
		ID:          paymentId,
		OrderID:     order.Id,                   // 订单ID
		FinalAmount: order.FinalAmount,          // 订单金额
		GatewayType: reqbody.PaymentGatewayType, // 支付网关类别
		Method:      model.PaymentMethodQrcode,  // 支付方式
		Status:      model.PaymentStatusToPay,   // 支付状态
		GatewayID:   yltOrderId,                 // YLT订单ID
	}
	dao.CreatePayment(payment)

	// 更新order表信息
	order.PaymentId = paymentId
	order.PaymentStatus = payment.Status
	if _, err = dao.UpdateOrderByField(order, []string{"payment_id", "payment_status"}); err != nil {
		log.Error("QrcodeOrderPaymentHandler 更新订单信息失败", zap.Error(err))
		return "", errors.New("")
	}

	return
}
