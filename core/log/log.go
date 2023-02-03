package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//声明日志类全局变量
var sugarLogger *zap.SugaredLogger

//日志类初始化方法
func InitLogger() *zap.SugaredLogger {

	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	zapLogger := zap.New(core, zap.AddCaller())
	sugarLogger = zapLogger.Sugar()

	return sugarLogger
}

//日志记录地址
func getLogWriter() zapcore.WriteSyncer {
	//定义日志文件名，设置权限，当日志文件不存在时创建文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/obs.log",
		MaxSize:    1,  // 切割大小
		MaxBackups: 5,  // 保留最大数量
		MaxAge:     30, // 保留最大天数
		Compress:   false,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)

}

//日志编码方式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// TODO 生产环境 和 开发环境
	//return zapcore.NewJSONEncoder(encoderConfig)
	return zapcore.NewConsoleEncoder(encoderConfig)
}
