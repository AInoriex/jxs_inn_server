package log

import (
	"eshop_server/src/utils/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var logger *zap.Logger

// 初始化日志
func InitLogger() *zap.Logger {
	var module string = config.CommonConfig.AppName
	var logPath string = config.CommonConfig.Log.SavePath
	if err := os.MkdirAll(logPath, 0755); err != nil {
		panic(err)
	}

	now := time.Now()
	hook := &lumberjack.Logger{
		// 日志存储位置
		Filename:   fmt.Sprintf("%s/%04d%02d%02d-%02d/%s.log", logPath, now.Year(), now.Month(), now.Day(), now.Hour(), module),
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

	// 同时输出到控制台和文件
	fileWriter := zapcore.AddSync(hook)
	consoleWriter := zapcore.AddSync(os.Stdout)
	writer := zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)

	core := zapcore.NewCore(getEncoder(), writer, zapcore.InfoLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger.Info("InitLogger success", zap.String("log_path", hook.Filename))

	zap.ReplaceGlobals(logger)
	return logger
}

func Sync() {
	logger.Sync()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}


func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, fields...)
}

func Success(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func Fail(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

func getLogWriter(savepath string) zapcore.WriteSyncer {
	file, _ := os.Create(savepath)
	return zapcore.AddSync(file)
}

func Infof(msg string, args ...interface{}) {
	zap.S().Infof(msg, args...)
}

func Infow(msg string, args ...interface{}) {
	zap.S().Infow(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	zap.S().Warnf(msg, args...)
}

func Warnw(msg string, args ...interface{}) {
	zap.S().Warnw(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	zap.S().Errorf(msg, args...)
}

func Errorw(msg string, args ...interface{}) {
	zap.S().Errorw(msg, args...)
}

func Debugf(msg string, args ...interface{}) {
	zap.S().Debugf(msg, args...)
}

func Debugw(msg string, args ...interface{}) {
	zap.S().Debugw(msg, args...)
}

func Fatalw(msg string, args ...interface{}) {
	zap.S().Fatalw(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	zap.S().Fatalf(msg, args...)
}
