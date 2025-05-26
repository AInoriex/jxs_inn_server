package handler

import (
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Title			查询用户藏品
// @Description		查询用户购买历史，封装商品信息数据返回
// @Router			/v1/eshop_api/user/inventory/list [get]
// @Response		json
func GetInventoryList(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Error("GetInventoryList 非法用户", zap.Error(err))
		FailWithAuthorization(c)
		return
	}
	log.Infof("GetInventoryList 请求参数, user_id: %s", user.Id)

	// 查询用户购买历史
	purchaseList, err := dao.GetPurchaseHistorysByUserId(user.Id)
	if err != nil {
		log.Errorf("GetInventoryList 查询购买历史记录失败, user_id: %s, error: %s", user.Id, err.Error())
		Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 构建返回数据
	GetInventoryListResponseList := make([]model.GetInventoryListResponse, 0)
	ExistProductIdMap := make(map[string]bool)
	for _, purchase := range purchaseList {
		// 过滤重复商品
		if _, ok := ExistProductIdMap[purchase.ProductId]; ok {
			continue
		}
		// 从商品表中查询商品信息，构建返回数据
		product, err := dao.GetProductById(purchase.ProductId)
		if err != nil {
			log.Errorf("GetInventoryList 查询商品信息失败, product_id: %s, error: %s", purchase.ProductId, err.Error())
			continue
		}
		GetInventoryListResponse := model.GetInventoryListResponse{
			ProductId:   product.Id,
			Title:       product.Title,
			Description: product.Description,
			ImageUrl:    product.ImageUrl,
			PurchaseAt:  purchase.PurchasedAt,
		}
		GetInventoryListResponseList = append(GetInventoryListResponseList, GetInventoryListResponse)
		ExistProductIdMap[product.Id] = true
	}

	dataMap["purchase_list"] = GetInventoryListResponseList
	log.Infof("GetInventoryList user_id: %s 查询作品成功，返回数据: %+v", user.Id, dataMap)
	Success(c, dataMap)
}
