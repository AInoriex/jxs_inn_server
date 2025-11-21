package dao

import (
	"eshop_server/src/router/model"
	"eshop_server/src/utils/db"
	"eshop_server/src/utils/log"
	"go.uber.org/zap"
	"time"
)

// @Title   获取指定单条订单记录
// @Description 根据订单ID获取订单记录
// @Author  AInoriex  (2025/05/16 16:43)
func GetOrderByOrderId(orderId string) (res *model.Order, err error) {
	err = db.MysqlCon.Where("id =?", orderId).First(&res).Error
	if err != nil {
		log.Error("GetOrderByOrderId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取用户所有订单记录
// @Description 用户id
// @Author  AInoriex  (2025/05/06 14:11)
func GetOrdersByUserId(userId string) (res []*model.Order, err error) {
	err = db.MysqlCon.Where("user_id = ?", userId).Find(&res).Error
	if err != nil {
		log.Error("GetOrdersByUserId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取用户单个订单记录
// @Description 用户id，订单id
// @Author  AInoriex  (2025/05/06 14:11)
func GetOrderByUserIdAndProductId(userId string, orderId string) (res *model.Order, err error) {
	err = db.MysqlCon.Where("user_id = ? and id = ?", userId, orderId).First(&res).Error
	if err != nil {
		log.Error("GetOrderByUserIdAndProductId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title	获取指定单个订单记录
// @Description 根据支付ID获取订单记录
// @Author	AInoriex (2025/05/15 20:00)
func GetOrderByPaymentId(paymentId string) (res *model.Order, err error) {
	err = db.MysqlCon.Where("payment_id =?", paymentId).First(&res).Error
	if err != nil {
		log.Error("GetOrderByPaymentId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title 获取全部订单记录
// @Description 管理后台获取全部订单信息
// @Author AInoriex (2025/05/16 16:43)
func GetAllOrders(pageNum int, pageSize int, orderBy string, orderType string) (res []*model.Order, err error) {
	err = db.MysqlCon.Find(&res).Order(orderBy + " " + orderType).
		Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	if err != nil {
		log.Error("GetAllOrders fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   创建数据记录
// @Description desc
// @Author  AInoriex  (2025/05/06 14:11)
func CreateOrder(m *model.Order) (res *model.Order, err error) {
	log.Infof("CreateOrder params, m:%+v", m)
	m.CreatedAt = time.Now()

	err, _ = db.Create(db.MysqlCon, &m)
	if err != nil {
		log.Error("CreateOrder fail", zap.Error(err))
		return m, err
	}

	return m, nil
}

// @Title   更新订单记录
// @Description 特定字段
// @Author  AInoriex  (2025/05/08 14:30)
func UpdateOrderByField(m *model.Order, field []string) (res *model.Order, err error) {
	m.UpdatedAt = time.Now()
	log.Infof("UpdateOrderByField params, m:%+v, field:%+v", m, field)
	// Select 除 Omit() 外的所有字段（包括零值字段的所有字段）
	err = db.MysqlCon.Model(&model.Order{}).Select(field).Omit("id").
		Where("id = ?", m.Id).Updates(m).Error
	if err != nil {
		log.Error("UpdateOrderByField fail ", zap.Any("m", m))
		return nil, err
	}

	return m, nil
}
