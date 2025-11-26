package handler

import (
	"encoding/json"
	"errors"
	"eshop_server/src/common/api"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	"eshop_server/src/utils/alarm"
	"eshop_server/src/utils/common"
	"eshop_server/src/utils/config"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/utime"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Title		 获取商品列表(前台商品售卖页)
// @Description  获取当前所有上架的商品
// @Response     json
// @Router       /v1/eshop_api/product/list?pageNum=1&pageSize=20 [get]
func GetProductList(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("GetProductList 请求参数, req:%s", string(req))

	// 获取商品信息
	pageNum := common.StringToIntNotErr(c.Query("pageNum"))
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := common.StringToIntNotErr(c.Query("pageSize"))
	if pageSize <= 0 {
		pageSize = 20
	}
	resList, err := dao.GetProductsByStatus(model.ProductStatusOn, pageNum, pageSize, "created_at", "desc")
	if err != nil {
		log.Errorf("GetProductList GetProductsByStatus fail, err:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
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
	api.Success(c, dataMap)
}

// @Title		 获取商品列表(后台)
// @Description  获取当前所有的商品
// @Response     json
// @Router       /v1/eshop_api/admin/product/list [get]
func AdminGetProductList(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("AdminGetProductList 请求参数, req:%s", string(req))

	// 获取商品信息
	// TODO 分页查询
	var resultProductList []model.CreateProductReq
	productList, err := dao.GetAllProducts(1, 50, "created_at", "desc")
	if err != nil {
		log.Errorf("AdminGetProductList GetAllProducts fail, err:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 对于每个商品获取商品播放列表
	_selectlimit := int32(10)
	_orderby := "updated_at desc"
	for _, v := range productList {
		var tmp model.CreateProductReq = model.CreateProductReq{
			Products: *v,
		}
		// 获取播放信息
		tmp.PP = make([]model.ProductsPlayer, 0)
		// playerList, err := dao.GetProductsPlayerByProductId(v.Id)
		playerList, err := dao.GetSelectedProductsPlayerByProductId(v.Id, model.ProductsPlayerStatusOk, _orderby, _selectlimit)
		if err != nil {
			log.Errorf("AdminGetProductList GetSelectedProductsPlayerByProductId fail, product_id:%s, err:%v", v.Id, err)
			api.Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
			return
		}
		for _, p := range playerList {
			if p.PlayUrl != "" || p.Status == model.ProductsPlayerStatusOk {
				tmp.PP = append(tmp.PP, *p)
			}
		}
		resultProductList = append(resultProductList, tmp)
	}

	// 返回数据
	dataMap["result"] = resultProductList
	dataMap["len"] = len(resultProductList)
	api.Success(c, dataMap)
}

// @Title		 创建商品
// @Description	 创建新商品并上架
// @Accept       json model.CreateProductReq
// @Response     json
// @Router       /v1/eshop_api/admin/product/create [post]
func AdminCreateProduct(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("AdminCreateProduct 请求参数, req:%s", string(req))

	// JSON解析
	var reqbody *model.CreateProductReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("AdminCreateProduct json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 参数判断和预处理
	if reqbody.Id == "" || reqbody.Title == "" || reqbody.Price <= 0 {
		log.Errorf("AdminCreateProduct 商品基本参数无效, reqbody:%+v", reqbody)
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":基本参数有误")
		return
	}
	if reqbody.ExternalId == "" || reqbody.ExternalLink == "" {
		log.Errorf("AdminCreateProduct 商品第三方参数无效, reqbody:%+v", reqbody)
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":第三方参数有误")
		return
	}
	if reqbody.ImageUrl == "" {
		reqbody.ImageUrl = model.ProductImageUrlDefault
	}

	// 创建上架商品
	res, err := dao.CreateProduct(&reqbody.Products)
	if err != nil {
		log.Errorf("AdminCreateProduct 创建商品失败, reqbody:%+v, err:%v", reqbody, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	// 创建商品资源映射
	for _, v := range reqbody.PP {
		if v.ProductId == "" {
			v.ProductId = res.Id
		}
		// 查找商品资源映射是否存在，不存在则创建
		if _, err = dao.GetAllProductsPlayerByProductId(v.ProductId); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("AdminCreateProduct 商品资源映射不存在即将创建映射关系, product_id:%s", v.ProductId)
			_, createErr := dao.CreateProductsPlayer(&v)
			if createErr != nil {
				log.Errorf("AdminCreateProduct 创建商品资源映射失败, reqbody:%+v, err:%v", reqbody, createErr)
				api.Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
				return
			}
		} else if err != nil {
			log.Errorf("AdminCreateProduct GetProductsPlayerByProductId failed, product_id:%s, err:%v", v.ProductId, err)
			api.Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail+":查询商品资源映射失败")
			return
		} else {
			log.Infof("AdminCreateProduct 商品资源映射已存在跳过创建, product_id:%s", v.ProductId)
		}
	}

	// 飞书告警
	if err = alarm.PostFeiShu(
		"info", config.CommonConfig.LarkAlarm.InfoBotWebhook,
		fmt.Sprintf("[JXS管理后台] 新商品上架成功 \n\t 环境:%s \n\t商品信息:%+v \n\t通知时间:%s", config.CommonConfig.Env, res, utime.TimeToStr(utime.GetNow())),
	); err != nil {
		log.Errorf("AdminCreateProduct 飞书通知失败, reqbody:%+v, err:%v", reqbody, err)
	}

	// 返回数据
	dataMap["result"] = res
	api.Success(c, dataMap)
}

// @Title		 下架商品
// @Description  下架商品使页面不可见
// @Param        product_id
// @Response     json
// @Router       /v1/eshop_api/admin/product/remove [put]
func AdminRemoveProduct(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("AdminRemoveProduct 请求参数, req:%s", string(req))

	// 查询商品ID信息
	ProductID := c.Param("id")

	// 查询数据库中的商品信息
	res, err := dao.GetProductById(ProductID)
	if err != nil {
		log.Errorf("AdminRemoveProduct GetProductById fail, ProductID:%s, err:%v", ProductID, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail+":查询商品失败")
		return
	}

	// 更新商品状态，将 status 从 1 更新为 0
	res.Status = model.ProductStatusOff
	res, err = dao.UpdateProductsByField(res, []string{"status"})
	if err != nil {
		log.Errorf("AdminRemoveProduct 更新商品状态失败, m:%+v, err:%v", res, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail+":更新商品状态失败")
		return
	}

	// TODO 飞书通知下架
	if err = alarm.PostFeiShu(
		"info", config.CommonConfig.LarkAlarm.InfoBotWebhook,
		fmt.Sprintf("[JXS管理后台] 商品下架成功 \n\t 环境:%s \n\t商品信息:%+v \n\t通知时间:%s", config.CommonConfig.Env, res, utime.TimeToStr(utime.GetNow())),
	); err != nil {
		log.Errorf("AdminRemoveProduct 飞书通知失败, res:%+v, err:%v", res, err)
	}

	// 返回成功响应
	dataMap["result"] = res
	api.Success(c, dataMap)
}

// @Title		 搜索商品
// @Description  搜索商品
// @Param        external_id
// @Response     json
// @Router       /v1/eshop_api/admin/product/search/external_id/:external_id [get]
func AdminSearchProductsByExternalId(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})
	
	// 获取路由参数中的 external_id
	externalID := c.Param("external_id")
	if externalID == "" {
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":external_id不能为空")
		return
	}
	log.Infof("AdminSearchProductsByExternalId 请求参数, external_id:%s", externalID)

	// 查询数据库中的商品信息
	res, err := dao.GetProductByExternalId(externalID)
	if err != nil {
		log.Errorf("AdminSearchProductsByExternalId GetProductByExternalId fail, externalID:%s, err:%v", externalID, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail+":查询商品失败")
		return
	}

	// 返回成功响应
	dataMap["result"] = res
	api.Success(c, dataMap)
}
