package handler

import (
	"encoding/json"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Title      获取购物车列表
// @Description  获取用户购物车列表
// @Response     json
// @Router       /v1/eshop_api/user/cart/list [get]
func GetCartList(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Error("GetCartList 非法用户", zap.Error(err))
		FailWithAuthorization(c)
		return
	}

	// TODO 参数签名解析
	// sign := c.DefaultQuery("sign", "0")
	// if !CheckSignParam(sign) {
	// 	log.Infof("GetCartList sign NOT pass. sign:%v", sign)
	// 	// Success(c, dataMap)
	// 	FailTrack(c, uerrors.Parse(uerrors.ErrApiParamSignNotPass.Error()).Code, uerrors.Parse(uerrors.ErrApiParamSignNotPass.Error()).Detail+"，别在这搞事哈", dataMap)
	// 	return
	// }

	// 获取购物车信息
	cartList, err := dao.GetCartItemsByUserId(user.Id)
	if err != nil {
		log.Error("GetCartList GetCartItemsByUserId fail", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 商品信息聚合
	var resList []model.GetCartListItemResponse
	for _, cartItem := range cartList {
		if cartItem.ProductId == "" {
			log.Error("GetCartList 获取商品详情信息失败，商品ID为空", zap.Int64("cart_id", cartItem.Id))
			continue
		}
		p, err := dao.GetProductById(cartItem.ProductId)
		if err != nil {
			log.Error("GetCartList 获取商品详情信息失败", zap.String("product_id", cartItem.ProductId), zap.Error(err))
			continue
		}
		retItem := model.GetCartListItemResponse{
			Id:       p.Id,
			Title:    p.Title,
			Price:    p.Price,
			Quantity: cartItem.Quantity,
			Image:    p.ImageUrl,
		}
		resList = append(resList, retItem)
	}

	// 返回数据
	dataMap["result"] = resList
	Success(c, dataMap)
}

// @Title      创建购物车商品
// @Description	 用户挑选商品并放入购物车
// @Body		 json
// @Response     json
// @Router       /v1/eshop_api/user/cart/create [post]
func CreateCart(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})
	req := GetGinBody(c)
	log.Info("CreateProduct 请求参数", zap.String("body", string(req)))

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Error("CreateCart 非法用户", zap.Error(err))
		FailWithAuthorization(c)
		return
	}

	// TODO 参数签名解析

	// JSON解析
	var reqbody model.CreateCartItemReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("CreateCart json解析失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 校验参数
	if reqbody.ProductId == "" || reqbody.Quantity <= 0 {
		log.Error("CreateCart 商品参数错误", zap.String("product_id", reqbody.ProductId), zap.Int32("quantity", reqbody.Quantity))
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail)
		return
	}
	// HARDNEED 商品数量校验，数量只支持1个
	if reqbody.Quantity != 1 {
		log.Error("CreateCart 商品数量错误", zap.Int32("quantity", reqbody.Quantity))
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":当前仅支持同种单件商品结算")
		return
	}

	// 校验商品是否有效
	_, err = dao.CheckProductById(reqbody.ProductId)
	if err!= nil {
		log.Error("CreateCart 获取商品信息失败", zap.String("product_id", reqbody.ProductId))
		Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail+":商品不存在")
		return	
	}
	// 校验用户购物车商品是否存在
	cartItem, err := dao.GetCartItemByUserIdAndProductId(user.Id, reqbody.ProductId)
	if err == nil && cartItem.Id > 0 {
		log.Warn("CreateCart 用户购物车商品已存在", zap.String("userId", cartItem.UserId), zap.String("productId", cartItem.ProductId))
		Success(c, dataMap)
		return
	}

	// 创建商品
	cartItem = &model.CartItem{
		UserId:    user.Id,
		ProductId: reqbody.ProductId,
		Quantity:  reqbody.Quantity,
	}
	_, err = dao.CreateCartItem(cartItem)
	if err != nil {
		log.Error("CreateCart 创建购物车商品失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	// 返回数据
	Success(c, dataMap)
}

// @Title      移除购物车商品
// @Description	 用户移除购物车不需要的商品
// @Body		 json
// @Response     json
// @Router       /v1/eshop_api/user/cart/remove [post]
func RemoveCart(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})
	req := GetGinBody(c)
	log.Info("RemoveCart 请求参数", zap.String("body", string(req)))

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Error("RemoveCart 非法用户", zap.Error(err))
		FailWithAuthorization(c)
		return
	}

	// TODO 参数签名解析

	// JSON解析
	var reqbody model.RemoveCartItemReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("RemoveCart json解析失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 校验参数
	if reqbody.ProductId == "" {
		log.Error("RemoveCart json参数错误", zap.String("cart_id", reqbody.ProductId))
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail)
		return
	}

	// 移除商品
	err = dao.RemoveUserCartProduct(user.Id, reqbody.ProductId)
	if err != nil {
		log.Error("RemoveCart 移除购物车商品失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	// 返回数据
	Success(c, dataMap)
}
