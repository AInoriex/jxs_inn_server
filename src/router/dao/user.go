package dao

import (
	"eshop_server/src/router/model"
	"eshop_server/src/utils/db"
	"eshop_server/src/utils/log"
	"go.uber.org/zap"
	"time"
)

// @Title   根据用户Id查询用户信息
// @Description 用户id
// @Author  AInoriex  (2025/05/06 15:07)
func GetUserById(id string) (res *model.User, err error) {
	err = db.MysqlCon.Where("id = ?", id).First(&res).Error
	if err != nil {
		log.Error("GetUserById fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   根据用户Id查询非黑名单用户信息
// @Description 用户id
// @Author  AInoriex  (2025/06/25 17:50)
func GetValidUserById(id string) (res *model.User, err error) {
	err = db.MysqlCon.Where("id = ?", id).Where("status = ?", model.UserStatusNormal).First(&res).Error
	if err != nil {
		log.Error("GetValidUserById fail", zap.Error(err))
		return nil, err
	}

	return
}

// @Title   根据邮箱获取用户信息
// @Description 邮箱email
// @Author  AInoriex  (2024/07/22 18:05)
func GetUserByEmail(email string) (res *model.User, err error) {
	err = db.MysqlCon.Where("email = ?", email).First(&res).Error
	if err != nil {
		log.Error("GetUserByEmail fail", zap.Error(err))
		return res, err
	}

	return res, nil
}

// @Title   根据邮箱获取非黑名单用户信息
// @Description 邮箱email
// @Author  AInoriex  (2025/06/25 17:50)
func GetValidUserByEmail(email string) (res *model.User, err error) {
	err = db.MysqlCon.Where("email = ?", email).Where("status = ?", model.UserStatusNormal).First(&res).Error
	if err != nil {
		log.Error("GetValidUserByEmail fail", zap.Error(err))
		return res, err
	}

	return res, nil
}

// @Title   创建新用户
// @Description *model.User
// @Author  AInoriex  (2024/07/22 18:05)
func CreateUser(m *model.User) (res *model.User, err error) {
	log.Info("CreateUser", zap.Any("req", m))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	err, _ = db.Create(db.MysqlCon, &m)
	if err != nil {
		log.Error("CreateUser fail", zap.Error(err))
		return m, err
	}

	return m, nil
}

// @Title   更新用户记录
// @Description 特定字段
// @Author  AInoriex  (2025/06/17 15:49)
func UpdateUserByField(m *model.User, field []string) (res *model.User, err error) {
	log.Infof("UpdateUserByField params, m:%+v, field:%+v", m, field)
	// Select 除 Omit() 外的所有字段
	err = db.MysqlCon.Model(&model.User{}).Select(field).Omit("id").Omit("create_at").
		Where("id = ?", m.Id).Updates(m).Error
	if err != nil {
		log.Error("UpdateUserByField fail ", zap.Any("m", m))
		return nil, err
	}

	return m, nil
}
