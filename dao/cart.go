package dao

import (
	"eshop_server/model"
	"eshop_server/utils/db"
	"eshop_server/utils/log"
	"go.uber.org/zap"
	"time"
)

// @Title   获取用户所有购物车记录
// @Description 用户id
// @Author  AInoriex  (2025/05/06 14:11)
func GetCartItemsByUserId(userId string) (res []*model.CartItem, err error) {
	err = db.MysqlCon.Where("user_id = ?", userId).Find(&res).Error
	if err != nil {
		log.Error("GetCartItemsByUserId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取用户单个商品购物车记录
// @Description 用户id，商品id
// @Author  AInoriex  (2025/05/06 14:11)
func GetCartItemByUserIdAndProductId(userId string, productId string) (res *model.CartItem, err error) {
	err = db.MysqlCon.Where("user_id = ? and product_id = ?", userId, productId).First(&res).Error
	if err != nil {
		log.Error("GetCartItemByUserIdAndProductId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   创建数据记录
// @Description desc
// @Author  AInoriex  (2025/05/06 14:11)
func CreateCartItem(m *model.CartItem) (res *model.CartItem, err error) {
	log.Infof("CreateCartItem params, m:%+v", m)
	m.CreatedAt = time.Now()

	err, _ = db.Create(db.MysqlCon, &m)
	if err != nil {
		log.Error("CreateCartItem fail", zap.Error(err))
		return m, err
	}

	return m, nil
}

// @Title   更新购物车记录
// @Description 特定字段
// @Author  AInoriex  (2025/05/08 14:30)
func UpdateCartItemByField(m *model.CartItem, field []string) (res *model.CartItem, err error) {
	log.Infof("UpdateCartItemByField params, m:%+v, field:%+v", m, field)
	// Select 除 Omit() 外的所有字段（包括零值字段的所有字段）
	err = db.MysqlCon.Model(&model.CartItem{}).Select(field).Omit("id").
		Where("id = ?", m.Id).Updates(m).Error
	if err != nil {
		log.Error("UpdateCartItemByField fail ", zap.Any("m", m))
		return nil, err
	}

	return m, nil
}

// @Title		删除购物车记录
// @Description	移除用户某件商品的购物车记录
// @Param1	userId 用户Id
// @Param2	productId 商品Id
// @Author	AInoriex  (2025/05/08 14:30)
func RemoveUserCartProduct(userId string, productId string) (err error) {
	log.Infof("DeleteCartItemById id:%+v", productId)
	err = db.MysqlCon.Where("product_id = ? and user_id = ?", productId, userId).Delete(&model.CartItem{}).Error
	if err != nil {
		log.Error("DeleteCartItemById fail ", zap.String("userId", userId), zap.String("productId", productId), zap.Error(err))
		return err
	}
	return nil
}
