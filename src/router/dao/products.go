package dao

import (
	"eshop_server/src/router/model"
	"eshop_server/src/utils/db"
	"eshop_server/src/utils/log"
	"go.uber.org/zap"
	"time"
)

// @Title   获取数据记录
// @Description 商品id
// @Author  AInoriex  (2024/07/22 18:05)
func GetProductById(id string) (res *model.Products, err error) {
	err = db.MysqlCon.Where("id = ?", id).First(&res).Error
	if err != nil {
		log.Error("GetProductById fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   检查商品是否有效
// @Description 商品id
// @Author  AInoriex  (2024/07/22 18:05)
func CheckProductById(id string) (res *model.Products, err error) {
	err = db.MysqlCon.Where("id = ? and status = ?", id, model.ProductStatusOn).First(&res).Error
	if err != nil {
		log.Error("CheckProductById fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取数据记录
// @Description 商品状态status
// @Author  AInoriex  (2024/07/22 18:05)
func GetProductsByStatus(status int32) (res []*model.Products, err error) {
	err = db.MysqlCon.Where("status = ?", status).Find(&res).Error
	if err != nil {
		log.Error("GetProductsByStatus fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   获取所有数据记录
// @Description 商品ID不为空
// @Author  AInoriex  (2025/06/26 16:22)
func GetAllProducts() (res []*model.Products, err error) {
	err = db.MysqlCon.Where("id != ''").Find(&res).Error
	if err != nil {
		log.Error("GetAllProducts fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   创建数据记录
// @Description desc
// @Author  AInoriex  (2024/07/22 18:05)
func CreateProduct(m *model.Products) (res *model.Products, err error) {
	log.Info("CreateProduct", zap.Any("req", m))
	m.CreateAt = time.Now()

	err, _ = db.Create(db.MysqlCon, &m)
	if err != nil {
		log.Error("CreateProduct fail", zap.Error(err))
		return m, err
	}

	return m, nil
}

// @Title   更新数据记录
// @Description 特定字段
// @Author  AInoriex  (2024/07/22 18:05)
func UpdateProductsByField(m *model.Products, field []string) (res *model.Products, err error) {
	log.Infof("UpdateProductsByField params, m:%+v, field:%+v", m, field)
	// Select 除 Omit() 外的所有字段（包括零值字段的所有字段）
	err = db.MysqlCon.Model(&model.Products{}).Select(field).Omit("id").
		Where("id = ?", m.Id).Updates(m).Error
	if err != nil {
		log.Error("UpdateProductsByField fail ", zap.Any("m", m))
		return nil, err
	}

	return m, nil
}

// @Title   删除数据记录
// @Description desc
// @Author  AInoriex  (2024/07/22 18:05)
func DeleteProducts(m *model.Products) (res *model.Products, err error) {
	log.Info("DeleteProducts params", zap.Any("m", m))
	reply, err := GetProductById(m.Id)
	if err != nil {
		return
	}

	// 带额外条件的删除
	err = db.MysqlCon.Where("id = ?", m.Id).Delete(&m).Error
	if err != nil {
		log.Error("DeleteProducts fail", zap.Error(err))
		return reply, err
	}

	return reply, nil
}

// @Title   replace数据记录
// @Description desc
// @Author  AInoriex  (2024/07/22 18:05)
// @Param   m *model.Products
// @Param   field []string
// @Return  *model.Products, error
// @Detail  如果id对应的记录不存在，则进行Insert操作，否则进行Update操作
func ReplaceProducts(m *model.Products, field []string) (res *model.Products, err error) {
	l, err := GetProductById(m.Id)
	if err != nil || l == nil {
		res, err = CreateProduct(m)
	} else {
		m.Id = l.Id
		res, err = UpdateProductsByField(m, field)
	}
	return
}

// @Title   更新status
// @Description source_type来源, language语言, old_stauts旧状态, new_status新状态
// @Return	rows, err
// @Author  AInoriex  (2024/11/13 17:34)
func UpdateProductsStatus(old_status int32, new_status int32) (rows int64, err error) {
	log.Infof("UpdateProductsStatus params, old_status:%v, new_status:%v", old_status, new_status)
	// 使用事务确保操作的原子性
	tx := db.GetDb().Begin()

	// 执行更新操作
	result := tx.Model(&model.Products{}).Where("status = ?", old_status).Updates(map[string]interface{}{"status": new_status})
	// 检查事务执行情况
	if result.Error != nil {
		log.Errorf("UpdateProductsStatus tx update error, err:%s", result.Error.Error())
		tx.Rollback()
		return 0, result.Error
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		log.Errorf("UpdateProductsStatus tx.Commit error, err:%s", err.Error())
		tx.Rollback()
		return 0, err
	}

	log.Infof("UpdateProductsStatus success, affect rows:%v, old_status:%v, new_status:%v", result.RowsAffected, old_status, new_status)
	return result.RowsAffected, nil
}
