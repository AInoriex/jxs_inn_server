package dao

import (
	"eshop_server/src/router/model"
	"eshop_server/src/utils/db"
	"eshop_server/src/utils/log"
	"go.uber.org/zap"
	"time"
)

// @Title   获取用户所有支付记录
// @Params	用户id
// @Author  AInoriex  (2025/05/06 14:11)
func GetPaymentsById(id string) (res *model.Payment, err error) {
	log.Infof("GetPaymentsById params, id:%+v", id)
	err = db.MysqlCon.Where("id = ?", id).First(&res).Error
	if err != nil {
		log.Error("GetPaymentsById fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取用户单个支付记录
// @Params	用户id，支付id
// @Author  AInoriex  (2025/05/06 14:11)
func GetPaymentByUserIdAndProductId(userId string, paymentId string) (res *model.Payment, err error) {
	log.Infof("GetPaymentByUserIdAndProductId params, userId:%+v, paymentId:%+v", userId, paymentId)
	err = db.MysqlCon.Where("user_id = ? and id = ?", userId, paymentId).First(&res).Error
	if err != nil {
		log.Error("GetPaymentByUserIdAndProductId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取用户所有支付记录
// @Params	用户id
// @Author  AInoriex  (2025/05/06 14:11)
func GetPaymentsByUserId(userId string) (res []*model.Payment, err error) {
	log.Infof("GetPaymentsByUserId params, userId:%+v", userId)
	err = db.MysqlCon.Where("user_id = ?", userId).Find(&res).Error
	if err != nil {
		log.Error("GetPaymentsByUserId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取指定支付状态的所有支付记录
// @Params	支付状态status
// @Author  AInoriex  (2025/05/15 11:15)
func GetPaymentsByStatus(status int32) (res []*model.Payment, err error) {
	log.Infof("GetPaymentsByStatus params, status:%+v", status)
	err = db.MysqlCon.Where("status = ?", status).Find(&res).Error
	if err != nil {
		log.Error("GetPaymentsByStatus fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   创建数据记录
// @Params	desc
// @Author  AInoriex  (2025/05/06 14:11)
func CreatePayment(m *model.Payment) (res *model.Payment, err error) {
	log.Infof("CreatePayment params, m:%+v", m)
	m.CreatedAt = time.Now()

	err, _ = db.Create(db.MysqlCon, &m)
	if err != nil {
		log.Error("CreatePayment fail", zap.Error(err))
		return m, err
	}

	return m, nil
}

// @Title   更新支付记录
// @Params	特定字段
// @Author  AInoriex  (2025/05/08 14:30)
func UpdatePaymentByField(m *model.Payment, field []string) (res *model.Payment, err error) {
	m.UpdatedAt = time.Now()
	log.Infof("UpdatePaymentByField params, m:%+v, field:%+v", m, field)
	// Select 除 Omit() 外的所有字段（包括零值字段的所有字段）
	err = db.MysqlCon.Model(&model.Payment{}).Select(field).Omit("id").
		Where("id = ?", m.Id).Updates(m).Error
	if err != nil {
		log.Error("UpdatePaymentByField fail ", zap.Any("m", m))
		return nil, err
	}

	return m, nil
}
