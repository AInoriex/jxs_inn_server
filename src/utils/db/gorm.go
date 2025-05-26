package db

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"eshop_server/src/utils/config"
	"os"
	"time"
)

type Db struct {
	*gorm.DB
	Tx   *gorm.DB
	Data interface{}
}

// 定义自己的Writer
type LogWriter struct {
	Log *zap.Logger
}

func (sqlDb *Db) Begins() {
	sqlDb.Tx = sqlDb.DB.Begin()
}

func (sqlDb *Db) Commits() error {
	err := sqlDb.Tx.Commit().Error
	sqlDb.Tx = nil
	return err
}

func (sqlDb *Db) Rollbacks() error {
	err := sqlDb.Tx.Rollback().Error
	sqlDb.Tx = nil
	return err
}

type dialOptions func(db *gorm.DB)

func NewLogWriter(l *zap.Logger) *LogWriter {
	return &LogWriter{Log: l}
}

// 实现gorm/logger.Writer接口
func (m *LogWriter) Printf(format string, v ...interface{}) {
	m.Log.Info(fmt.Sprintf(format, v...))
}

// @Title   初始化Mysql
// @Description mysql基于zap写入日志文件, 但目前日志未切割
// @Author  wzj  (2022/8/10 18:05)
// @Param	args string 数据库dsn
// @Param	maxCon int 最大连接数
// @Return
func NewMysql(args string, maxCon int, arr []interface{}, enable bool) *gorm.DB {
	var con *gorm.DB
	var err error
	var gormCfg = &gorm.Config{
		//Logger:logger.Default.LogMode(logger.Info), //开启sql日志
	}

	if enable {
		//开启sql日志
		var logPath string = config.CommonConfig.Log.SavePath
		if err := os.MkdirAll(logPath, 0755); err != nil {
			panic(err)
		}

		now := time.Now()
		hook := &lumberjack.Logger{
			// 日志存储位置
			Filename:   fmt.Sprintf("%s/%04d%02d%02d-%02d/%s.log", logPath, now.Year(), now.Month(), now.Day(), now.Hour(), "db"),
			// 日志文件大小单位: M
			MaxSize:    500,                                       
			// 备份数
			MaxBackups: 50,                                        
			// days
			MaxAge:     30,                                        
			// disabled by default
			Compress:   true,  
		}
		defer hook.Close()
		writer := zapcore.AddSync(hook)
		core := zapcore.NewCore(getEncoder(), writer, zapcore.InfoLevel)
		l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

		gormCfg.Logger = logger.New(
			NewLogWriter(l),
			//log.New(file, "\r\n",log.LstdFlags | log.Lshortfile | log.LUTC), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
	}

	con, err = gorm.Open(mysql.Open(args), gormCfg)
	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, arg = %v the error is '%v'", args, err))
	}
	//设置表名前缀
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return defaultTableName[:]
	//}

	sqlDB, err := con.DB()
	if err != nil {
		panic(fmt.Sprintf("Got error when get con.DB, arg = %v the error is '%v'", args, err))
	}

	idle := maxCon
	if maxCon/3 >= 10 {
		idle = maxCon / 3
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(idle)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(2 * time.Hour)

	//若结构有变，则删除表重新创建
	//dropTable(con, arr...)
	//con.AutoMigrate(arr...) //若没有表，自动生成表
	return con
}

func (sqlDb *Db) Create() (error, interface{}) {
	//var err error
	result := sqlDb.DB.Create(sqlDb.Data)
	return result.Error, result.RowsAffected
}

func Create(db *gorm.DB, value interface{}) (error, interface{}) {
	//var err error
	result := db.Create(value)
	return result.Error, result.RowsAffected
}

func Save(db *gorm.DB, v interface{}) error {
	var err error
	err = db.Save(v).Error
	return err
}

// update语句 不允许使用orm特性
func getLogWriter(savepath string) zapcore.WriteSyncer {
	file, _ := os.OpenFile(savepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return zapcore.AddSync(file)
}

// 获取日志格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
