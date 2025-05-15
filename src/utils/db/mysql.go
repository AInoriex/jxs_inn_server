package db

import (
	"fmt"
	"gorm.io/gorm"
)

type TypeMysql int32

const (
	Con_Main TypeMysql = 1
	Con_Log  TypeMysql = 2
)

var (
	MysqlCon    *gorm.DB //1
	MysqlLogCon *gorm.DB //2
	MysqlActCon *gorm.DB //3

)

func InitMysqlAll(host, db string, maxCon int, cate TypeMysql, enable bool, arg ...interface{}) {
	dsn := fmt.Sprintf("%v/%v?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local", host, db)
	con := NewMysql(dsn, maxCon, arg, enable)
	switch cate {
	case Con_Main:
		MysqlCon = con
	case Con_Log:
		MysqlLogCon = con
	}
	// log.Info("初始化mysql完成", zap.Any("cate", cate), zap.Any("arg", arg))
}

func InTransaction(db *gorm.DB, withTransaction func(transaction *gorm.DB) (interface{}, error)) (interface{}, error) {
	tran := db.Begin()
	obj, err := withTransaction(tran)
	if err == nil {
		tran.Commit()
	} else {
		tran.Rollback()
	}
	return obj, err
}

func GetDb() *gorm.DB {
	return MysqlCon
}

func GetDbLog() *gorm.DB {
	return MysqlLogCon
}
