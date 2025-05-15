package dao

import (
	"eshop_server/src/router/model"
	"eshop_server/utils/db"
	"eshop_server/utils/log"
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
