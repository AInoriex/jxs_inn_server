package dao

import (
	"errors"
	"eshop_server/src/router/model"
	"eshop_server/src/utils/db"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uuid"
	"time"
)

// @Title   获取数据记录
// @Description 播放器id
// @Author  AInoriex  (2025/08/12 11:04)
func GetProductsPlayerById(id string) (res *model.ProductsPlayer, err error) {
	err = db.MysqlCon.Where("id = ?", id).First(&res).Error
	if err != nil {
		log.Errorf("GetProductsPlayerById fail, id:%s, err:%+v", id, err)
		return nil, err
	}

	return
}

// @Title   获取数据记录列表
// @Description 商品id
// @Author  AInoriex  (2025/08/12 11:04)
func GetAllProductsPlayerByProductId(product_id string) (res []*model.ProductsPlayer, err error) {
	err = db.MysqlCon.Where("product_id = ?", product_id).Find(&res).Error
	if err != nil {
		log.Errorf("GetProductsPlayerByProductId fail, product_id:%s, err:%+v", product_id, err)
		return nil, err
	}

	return
}

// @Title   获取指定数据记录列表
// @Description 商品id, 状态status
// @Author  AInoriex  (2025/08/12 11:04)
func GetSelectedProductsPlayerByProductId(product_id string, status int32, order string, limit int32) (res []*model.ProductsPlayer, err error) {
	err = db.MysqlCon.Where("product_id = ? and status = ?", product_id, status).
		Order(order).Limit(int(limit)).Find(&res).Error
	if err != nil {
		log.Errorf("GetSelectedProductsPlayerByProductId fail, product_id:%s, err:%+v", product_id, err)
		return nil, err
	}

	return
}

// @Title   创建数据记录
// @Description desc
// @Author  AInoriex  (2025/08/12 11:04)
func CreateProductsPlayer(m *model.ProductsPlayer) (res *model.ProductsPlayer, err error) {
	log.Infof("CreateProductsPlayer params, m:%+v", m)
	m.CreateAt = time.Now()
	m.UpdateAt = time.Now()
	err, _ = db.Create(db.MysqlCon, &m)
	if err != nil {
		log.Errorf("CreateProductsPlayer fail, err:%+v", err)
		return m, err
	}

	return m, nil
}

// @Title   更新数据记录
// @Description 特定字段
// @Author  AInoriex  (2025/08/12 11:04)
func UpdateProductsPlayerByField(m *model.ProductsPlayer, field []string) (res *model.ProductsPlayer, err error) {
	m.UpdateAt = time.Now()
	log.Infof("UpdateProductsPlayerByField params, m:%+v, field:%+v", m, field)
	// Select 除 Omit() 外的所有字段（包括零值字段的所有字段）
	err = db.MysqlCon.Model(&model.ProductsPlayer{}).Select(field).Omit("id").
		Where("id = ?", m.Id).Updates(m).Error
	if err != nil {
		log.Errorf("UpdateProductsPlayerByField fail, m:%+v, err:%+v ", m, err)
		return nil, err
	}

	return m, nil
}

// @Title   删除数据记录
// @Description desc
// @Author  AInoriex  (2025/08/12 11:04)
func DeleteProductsPlayer(m *model.ProductsPlayer) (res *model.ProductsPlayer, err error) {
	log.Infof("CreateProductsPlayer params, m:%+v", m)
	record, err := GetProductsPlayerById(m.Id)
	if err != nil {
		return
	}

	// 带额外条件的删除
	err = db.MysqlCon.Where("id = ?", m.Id).Delete(&m).Error
	if err != nil {
		log.Errorf("DeleteProductsPlayer fail, err:%+v", err)
		return record, err
	}

	return record, nil
}

// @Title   replace数据记录
// @Description 根据player_id检索记录，若不存在则新建，否则更新
// @Author  AInoriex  (2025/11/20 18:38)
// @Param   m *model.ProductsPlayer
// @Param   field []string
// @Return  *model.ProductsPlayer, error
func ReplaceProductsPlayer(m *model.ProductsPlayer, field []string) (res *model.ProductsPlayer, err error) {
	l, err := GetProductsPlayerById(m.Id)
	if err != nil || l == nil {
		res, err = CreateProductsPlayer(m)
	} else {
		m.Id = l.Id
		res, err = UpdateProductsPlayerByField(m, field)
	}
	return
}

// @Title   replace数据记录
// @Description 根据product_id检索所有记录，若不存在则新建，否则更新
// @Author  AInoriex  (2025/11/20 18:38)
// @Param   m *model.ProductsPlayer
// @Param   field []string
// @Return  *model.ProductsPlayer, error
func ReplaceProductsPlayerByProductId(m *model.ProductsPlayer, field []string) (res *model.ProductsPlayer, err error) {
	if m.ProductId == "" {
		log.Errorf("ReplaceProductsPlayerByProductId 请求参数错误, product_id不能为空")
		return nil, errors.New("product_id不能为空")
	}
	if m.Id == "" {
		m.Id = uuid.GetUuid()
	}

	return ReplaceProductsPlayer(m, field)
}