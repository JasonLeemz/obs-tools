package controller

import (
	"github.com/JasonLeemz/obs-tools/core/log"
	"html/template"
	"net/http"
	"os"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	// 后台首页
	logger := log.InitLogger()

	dir, err := os.Getwd()
	if err != nil {
		logger.Error(err)
		return
	}

	t, err := template.ParseFiles(dir + "/static/index.html")
	if err != nil {
		logger.Error(err)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		//"company": data,
	})
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("ok:", nil)

}

func Kuaishou(w http.ResponseWriter, r *http.Request) {
	// 后台首页
	logger := log.InitLogger()

	dir, err := os.Getwd()
	if err != nil {
		logger.Error(err)
		return
	}

	t, err := template.ParseFiles(dir + "/static/kuaishou.html")
	if err != nil {
		logger.Error(err)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		//"company": data,
	})
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("ok:", nil)

}
