package dao

import (
	"eshop_server/src/router/model"
	"eshop_server/utils/db"
	"eshop_server/utils/log"
	"go.uber.org/zap"
)

// @Title   获取用户所有支付记录
// @Description 用户id
// @Author  AInoriex  (2025/05/12 17:16)
func GetPurchaseHistorysByUserId(userId string) (res []*model.PurchaseHistory, err error) {
	log.Infof("GetPurchaseHistorysByUserId params, userId:%s", userId)
	err = db.MysqlCon.Where("user_id = ?", userId).Find(&res).Error
	if err != nil {
		log.Error("GetPurchaseHistorysByUserId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取用户单个支付记录
// @Description 用户id，支付id
// @Author  AInoriex  (2025/05/12 17:16)
func GetPurchaseHistoryByUserIdAndProductIdAndPaymentId(userId string, productId string, paymentId string) (res *model.PurchaseHistory, err error) {
	log.Infof("GetPurchaseHistoryByUserIdAndProductIdAndPaymentId params, userId:%s, productId:%s, paymentId:%s", userId, productId, paymentId)
	err = db.MysqlCon.Where("user_id = ? and product_id = ? and payment_id = ?", userId, productId, paymentId).First(&res).Error
	if err != nil {
		log.Error("GetPurchaseHistoryByUserIdAndProductIdAndPaymentId fail", zap.Error(err))
		return nil, err
	}
	
	return
}

// @Title   创建数据记录
// @Description desc
// @Author  AInoriex  (2025/05/12 17:16)
func CreatePurchaseHistory(m *model.PurchaseHistory) (res *model.PurchaseHistory, err error) {
	log.Infof("CreatePurchaseHistory params, m:%+v", m)
	err, _ = db.Create(db.MysqlCon, &m)
	if err != nil {
		log.Error("CreatePurchaseHistory fail", zap.Error(err))
		return m, err
	}

	return m, nil
}

// @Title   更新支付记录
// @Description 特定字段
// @Author  AInoriex  (2025/05/08 14:30)
func UpdatePurchaseHistoryByField(m *model.PurchaseHistory, field []string) (res *model.PurchaseHistory, err error) {
	log.Infof("UpdatePurchaseHistoryByField params, m:%+v, field:%+v", m, field)
	// Select 除 Omit() 外的所有字段（包括零值字段的所有字段）
	err = db.MysqlCon.Model(&model.PurchaseHistory{}).Select(field).Omit("id").
		Where("id = ?", m.Id).Updates(m).Error
	if err != nil {
		log.Error("UpdatePurchaseHistoryByField fail ", zap.Any("m", m))
		return nil, err
	}

	return m, nil
}
