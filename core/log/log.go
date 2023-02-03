package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

//声明日志类全局变量
var logger *zap.Logger

//日志类初始化方法
func InitLogger() *zap.Logger {
	////需要将日志写入什么地方
	//writeSyncer := getLogWriter()
	////日志编码方式
	//encoder := getEncoder()
	////传入参数进行实例化，最后一个参数为日志记录级别，当日志级别小于当前级别时，不写入日志文件
	//core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	//logger = zap.New(core, zap.AddCaller())

	var err error
	logger, err = zap.NewDevelopment() // 忽略了错误
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // 将 buffer 中的日志写到文件中

	//使用日志方法
	//logger.Debug("调试日志", zap.String("code", "200"))
	//logger.Info("详情日志", zap.String("code", "200"))
	//logger.Error("错误记录", zap.String("code", "200"))
	//logger.Fatal("失败日志", zap.String("code", "200"))
	////logger.Panic("系统错误日志",zap.String("code","200"))

	return logger
}

//日志记录地址
func getLogWriter() zapcore.WriteSyncer {
	//定义日志文件名，设置权限，当日志文件不存在时创建文件
	file, err := os.OpenFile("./log/run.log", os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Printf("打开文件资源失败 err:%v", err)
	}
	return zapcore.AddSync(file)
}

//日志编码方式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
