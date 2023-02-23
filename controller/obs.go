package controller

import (
	"github.com/JasonLeemz/obs-tools/core/log"
	"github.com/JasonLeemz/obs-tools/tools"
	"net/http"
)

func PushStream(w http.ResponseWriter, r *http.Request) {
	//实例化日志类
	logger := log.InitLogger()

	err := tools.Push()
	logger.Info("推送结果:", err)

}
