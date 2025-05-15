package dao

import (
	"eshop_server/src/router/model"
	"eshop_server/utils/db"
	"eshop_server/utils/log"
	"go.uber.org/zap"
	"time"
)

// @Title   获取订单下所有物品记录
// @Description 订单物品id
// @Author  AInoriex  (2025/05/06 14:11)
func GetOrderItemsById(id string) (res []*model.OrderItem, err error) {
	err = db.MysqlCon.Where("id = ?", id).Find(&res).Error
	if err != nil {
		log.Error("GetOrderItemsById fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取用户单个订单订单记录
// @Description 用户id，订单id
// @Author  AInoriex  (2025/05/06 14:11)
func GetOrderItemByUserIdAndProductId(userId string, orderId string) (res *model.OrderItem, err error) {
	err = db.MysqlCon.Where("user_id = ? and id = ?", userId, orderId).First(&res).Error
	if err != nil {
		log.Error("GetOrderItemByUserIdAndProductId fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   创建数据记录
// @Description desc
// @Author  AInoriex  (2025/05/06 14:11)
func CreateOrderItem(m *model.OrderItem) (res *model.OrderItem, err error) {
	log.Infof("CreateOrderItem params, m:%+v", m)
	m.CreatedAt = time.Now()

	err, _ = db.Create(db.MysqlCon, &m)
	if err != nil {
		log.Error("CreateOrderItem fail", zap.Error(err))
		return m, err
	}

	return m, nil
}

// @Title   更新订单记录
// @Description 特定字段
// @Author  AInoriex  (2025/05/08 14:30)
func UpdateOrderItemByField(m *model.OrderItem, field []string) (res *model.OrderItem, err error) {
	log.Infof("UpdateOrderItemByField params, m:%+v, field:%+v", m, field)
	// Select 除 Omit() 外的所有字段（包括零值字段的所有字段）
	err = db.MysqlCon.Model(&model.OrderItem{}).Select(field).Omit("id").
		Where("id = ?", m.Id).Updates(m).Error
	if err != nil {
		log.Error("UpdateOrderItemByField fail ", zap.Any("m", m))
		return nil, err
	}

	return m, nil
}
