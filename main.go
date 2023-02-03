package main

import (
	"fmt"
	"github.com/JasonLeemz/obs-tools/core/log"
	"github.com/JasonLeemz/obs-tools/tools"
	"go.uber.org/zap"
)

func main() {
	//实例化日志类
	logger := log.InitLogger()

	err := tools.Push()
	errStr := fmt.Sprintf("%v", err)
	logger.Info("推送结果", zap.String("Push", errStr))

}
