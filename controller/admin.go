package controller

import (
	"github.com/JasonLeemz/obs-tools/core/log"
	"net/http"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	// 后台首页
	logger := log.InitLogger()

	logger.Info("ok:", nil)

}
