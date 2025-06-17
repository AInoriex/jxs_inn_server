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

// @Title      获取商品列表
// @Description  获取当前所有上架的商品
// @Param        sign
// @Response     json
// @Router       /v1/eshop_api/product/list [get]
func GetProductList(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	// attrMap := make(map[string]interface{})
	log.Info("GetProductList 请求参数", zap.String("body", string(req)))

	// TODO 参数签名解析
	// sign := c.DefaultQuery("sign", "0")
	// if !CheckSignParam(sign) {
	// 	log.Infof("GetProductList sign NOT pass. sign:%v", sign)
	// 	// Success(c, dataMap)
	// 	FailTrack(c, uerrors.Parse(uerrors.ErrApiParamSignNotPass.Error()).Code, uerrors.Parse(uerrors.ErrApiParamSignNotPass.Error()).Detail+"，别在这搞事哈", dataMap)
	// 	return
	// }

	// 获取商品信息
	resList, err := dao.GetProductsByStatus(model.ProductStatusOn)
	if err != nil {
		log.Error("GetProductList GetProductsByStatus fail", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 返回数据
	dataMap["result"] = resList
	Success(c, dataMap)
}

// @Title      上架商品
// @Description  创建新商品并上架
// @Accept       json model.Products
// @Response     json
// @Router       /v1/eshop_api/product/create [post]
func CreateProduct(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	// attrMap := make(map[string]interface{})
	log.Info("CreateProduct 请求参数", zap.String("body", string(req)))

	// JSON解析
	var reqbody model.Products
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("CreateProduct json解析失败, error:%v", err)
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 参数判断和预处理
	if reqbody.Title == "" || reqbody.Description == "" || reqbody.Price == 0 {
		log.Error("CreateProduct 商品参数无效", zap.Any("reqbody", reqbody))
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail)
		return
	}
	if reqbody.ImageUrl == "" { // TODO 使用默认图片
		reqbody.ImageUrl = ""
	}
	reqbody.Status = model.ProductStatusOn

	// 创建上架商品
	res, err := dao.CreateProduct(&reqbody)
	if err != nil {
		log.Error("CreateProduct 数据库创建商品失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	// TODO 飞书通知

	// 返回数据
	dataMap["result"] = res
	Success(c, dataMap)
}

// @Title      下架商品
// @Description  下架商品使页面不可见
// @Param        product_id
// @Response     json
// @Router       /v1/eshop_api/product/remove [put]
func RemoveProduct(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	// attrMap := make(map[string]interface{})
	log.Info("GetProductList 请求参数", zap.String("body", string(req)))

	// 查询商品ID信息
	ProductID := c.Param("id")

	// 查询数据库中的商品信息
	res, err := dao.GetProductById(ProductID)
	if err != nil {
		log.Error("RemoveProduct GetProductById fail", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 更新商品状态，将 status 从 1 更新为 0
	res.Status = model.ProductStatusOff
	res, err = dao.UpdateProductsByField(res, []string{"status"})
	if err != nil {
		log.Error("RemoveProduct 更新商品状态失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	// TODO 飞书通知下架

	// 返回成功响应
	dataMap["result"] = res
	Success(c, dataMap)
}

// TODO 获取商品详情信息
// func GetProductInfo(*gin.Context) {
// }
