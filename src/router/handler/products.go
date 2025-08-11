package handler

import (
	"encoding/json"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"github.com/gin-gonic/gin"
)

// @Title		 获取商品列表
// @Description  获取当前所有上架的商品
// @Response     json
// @Router       /v1/eshop_api/product/list [get]
func GetProductList(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("GetProductList 请求参数, req:%s", string(req))

	// 获取商品信息
	resList, err := dao.GetProductsByStatus(model.ProductStatusOn)
	if err != nil {
		log.Errorf("GetProductList GetProductsByStatus fail, err:%v", err)
		Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 格式化返回结果
	var resUserList []*model.ProductUserView
	for _, v := range resList {
		resUserList = append(resUserList, v.UserViewFormat())
	}

	// 返回数据
	dataMap["result"] = resUserList
	dataMap["len"] = len(resUserList)
	Success(c, dataMap)
}

// @Title		 获取商品列表(后台)
// @Description  获取当前所有的商品
// @Response     json
// @Router       /v1/eshop_api/admin/product/list [get]
func AdminGetProductList(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("AdminGetProductList 请求参数, req:%s", string(req))

	// 获取商品信息
	// TODO 分页查询
	resList, err := dao.GetAllProducts(1, 50, "created_at", "desc")
	if err != nil {
		log.Errorf("AdminGetProductList GetAllProducts fail, err:%v", err)
		Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 返回数据
	dataMap["result"] = resList
	dataMap["len"] = len(resList)
	Success(c, dataMap)
}

// @Title		 上架商品
// @Description	 创建新商品并上架
// @Accept       json model.CreateProductReq
// @Response     json
// @Router       /v1/eshop_api/product/create [post]
func AdminCreateProduct(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("CreateProduct 请求参数, req:%s", string(req))

	// JSON解析
	var reqbody model.Products
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("CreateProduct json解析失败, error:%v", err)
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 参数判断和预处理
	if reqbody.Id == "" || reqbody.Title == "" || reqbody.Price <= 0 {
		log.Errorf("CreateProduct 商品基本参数无效, reqbody:%+v", reqbody)
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":基本参数有误")
		return
	}
	if reqbody.ExternalId == "" || reqbody.ExternalLink == "" { // TODO 使用默认ID
		log.Errorf("CreateProduct 商品第三方参数无效, reqbody:%+v", reqbody)
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":第三方参数有误")
		return
	}
	if reqbody.ImageUrl == "" { // TODO 使用默认图片
		reqbody.ImageUrl = model.ProductImageUrlDefault
	}

	// 创建上架商品
	res, err := dao.CreateProduct(&reqbody)
	if err != nil {
		log.Errorf("CreateProduct 创建商品失败, reqbody:%+v, err:%v", reqbody, err)
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	// TODO 飞书通知

	// 返回数据
	dataMap["result"] = res
	Success(c, dataMap)
}

// @Title		 下架商品
// @Description  下架商品使页面不可见
// @Param        product_id
// @Response     json
// @Router       /v1/eshop_api/product/remove [put]
func AdminRemoveProduct(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("GetProductList 请求参数, req:%s", string(req))

	// 查询商品ID信息
	ProductID := c.Param("id")

	// 查询数据库中的商品信息
	res, err := dao.GetProductById(ProductID)
	if err != nil {
		log.Errorf("RemoveProduct GetProductById fail, ProductID:%s, err:%v", ProductID, err)
		Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 更新商品状态，将 status 从 1 更新为 0
	res.Status = model.ProductStatusOff
	res, err = dao.UpdateProductsByField(res, []string{"status"})
	if err != nil {
		log.Errorf("RemoveProduct 更新商品状态失败, m:%+v, err:%v", res, err)
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	// TODO 飞书通知下架

	// 返回成功响应
	dataMap["result"] = res
	Success(c, dataMap)
}
