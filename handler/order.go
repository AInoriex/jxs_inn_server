package handler

import (
	"encoding/json"
	"eshop_server/dao"
	"eshop_server/model"
	uerrors "eshop_server/utils/errors"
	"eshop_server/utils/log"
	"eshop_server/utils/uuid"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary      创建订单
// @Description	 用户购物车结算并创建订单
// @Body		 json
// @Response     json
// @Router       /v1/eshop_api/user/order/create [post]
func CreateOrder(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})
	req := GetGinBody(c)
	log.Info("CreateOrder 请求参数", zap.String("body", string(req)))

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Error("CreateOrder 非法用户", zap.Error(err))
		FailWithAuthorization(c)
		return
	}

	// TODO 参数签名解析

	// JSON解析
	var reqbody model.CreateOrderReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("CreateOrder json解析失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 校验参数
	if len(reqbody.ItemList) <= 0 {
		log.Error("CreateOrder 商品列表为空")
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":商品列表为空")
		return
	}
	if reqbody.PaymentGatewayType <= 0|| reqbody.PaymentMethod == "" {
		log.Error("CreateOrder 支付参数错误", zap.Int32("payment_gateway", reqbody.PaymentGatewayType), zap.String("payment_method", reqbody.PaymentMethod))
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":支付参数无效")
		return
	}
	// HARDNEED 当前只支持同种且单件商品结算
	if len(reqbody.ItemList) > 1 {
		log.Error("CreateOrder 商品列表过多，当前仅支持同种且单件商品结算", zap.Int("item_list", len(reqbody.ItemList)))
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":当前仅支持同种单件商品结算")
		return
	}
	for _, item := range reqbody.ItemList {
		if item.ProductId == "" || item.Quantity <= 0 {
			log.Error("CreateOrder 商品参数错误", zap.String("product_id", item.ProductId), zap.Int32("quantity", item.Quantity))
			Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":商品信息无效")
			return
		}
		// HARDNEED 商品数量校验，数量只支持1个
		if item.Quantity != 1 {
			log.Error("CreateOrder 商品数量错误", zap.Int32("quantity", item.Quantity))
			Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":当前仅支持同种单件商品结算")
			return
		}
	}
	// HARDNEED 支付方式校验，当前仅支持扫码支付+YLT支付
	if reqbody.PaymentMethod != model.PaymentMethodQrcode || reqbody.PaymentGatewayType != model.PaymentGatewayTypeYlt {
		log.Error("CreateOrder 支付参数无效", zap.String("payment_method", reqbody.PaymentMethod), zap.Int32("payment_gateway", reqbody.PaymentGatewayType))
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":支付参数无效")
		return
	}
	// 校验商品是否有效&计算总价
	var totalAmount float64 = 0.00
	for _, item := range reqbody.ItemList {
		p, err := dao.CheckProductById(item.ProductId)
		if err != nil {
			log.Error("CreateOrder 获取商品信息失败", zap.String("product_id", item.ProductId))
			Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail+":商品不存在")
			return
		}
		totalAmount += p.Price * float64(item.Quantity)
	}

	// 创建订单商品
	orderItemId := uuid.GetUuid()
	for _, item := range reqbody.ItemList {
		OrderItem := &model.OrderItem{
			Id:        orderItemId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
		if _, err = dao.CreateOrderItem(OrderItem); err != nil {
			log.Error("CreateOrder 创建订单商品失败", zap.Error(err))
			Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
			return
		}
	}

	// 创建订单
	order := &model.Order{
		Id:            uuid.GetUuid(),
		UserId:        user.Id,
		ItemId:        orderItemId,
		TotalAmount:   totalAmount,
		PaymentId:     "",
		PaymentStatus: model.PaymentStatusToPay,
	}
	order.Discount = 0.00 // TODO 优惠券折扣计算
	order.FinalAmount = order.TotalAmount - order.Discount
	if order.FinalAmount <= 0 { // 防止金额越界
		order.FinalAmount = 0.00
	}

	// 创建支付
	payment := &model.Payment{
		ID:          uuid.GetUuid(),
		OrderID:     order.Id,
		FinalAmount: order.FinalAmount,		
	}

	// 创建支付流程，获取二维码
	qrcode_base64, err := QrcodeOrderPaymentHandler(reqbody)
	if err!= nil {
		log.Error("CreateOrder 创建二维码支付流程失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail+":创建二维码支付流程失败")
		return
	}

	// 创建订单
	order, err = dao.CreateOrder(order)
	if err != nil {
		log.Error("CreateOrder 创建订单失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail+":创建订单失败")
		return
	}

	// 检查用户购物车中是否存在该商品，存在则移除
	go func() {
		for _, item := range reqbody.ItemList {
			cart, err := dao.GetCartItemByUserIdAndProductId(user.Id, item.ProductId)	
			if err == nil && cart.Id != 0 {
				if err = dao.RemoveUserCartProduct(user.Id, item.ProductId); err != nil {
					log.Error("CreateOrder 移除用户购物车商品失败", zap.String("userId", user.Id), zap.String("productId", item.ProductId), zap.Error(err))
				}
			}
		}
	}()

	// 返回数据
	dataMap["order_id"] = order.Id
	dataMap["qrcode"] = qrcode_base64
	Success(c, dataMap)
}
