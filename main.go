package main

import (
	"fmt"
	"github.com/JasonLeemz/obs-tools/controller"
	"github.com/JasonLeemz/obs-tools/core/log"
	"io"
	"net/http"
)

func main() {
	//实例化日志类
	logger := log.InitLogger()

	// 设置访问路由
	http.HandleFunc("/", startServer)
	http.HandleFunc("/push", controller.PushStream)
	http.HandleFunc("/admin", controller.Admin)
	http.HandleFunc("/index", controller.Admin)
	// 设置监听的端口
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		logger.Fatal("ListenAndServe: ", err)
		panic(err)
	}

	defer logger.Sync() // 将 buffer 中的日志写到文件中
}

func startServer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）

	fmt.Fprintf(w, "service start...") //这个写入到w的是输出到客户端的 也可以用下面的 io.WriteString对象

	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	//query := r.URL.Query()
	//var uid string // 初始化定义变量
	//if r.Method == "GET" {
	//	uid = r.FormValue("uid")
	//} else if r.Method == "POST" {
	//	uid = r.PostFormValue("uid")
	//}
	io.WriteString(w, "service start...")
}

//
//func switchPlatform() (string, error) {
//	// the questions to ask
//	var qs = []*survey.Question{
//		{
//			Name: "rtmptype",
//			Prompt: &survey.Select{
//				Message: "选择一个推流平台,回车键确认:",
//				Options: []string{
//					"快手",
//					"哔哩哔哩",
//					//"全部",
//				},
//			},
//		},
//	}
//
//	var answers = struct {
//		RtmpType string `survey:"rtmptype"` // survey 会默认匹配首字母小写的name
//	}{}
//
//	// 执行提问
//	err := survey.Ask(qs, &answers)
//	if err != nil {
//		return "", err
//	}
//
//	rtmpType := strings.TrimSpace(answers.RtmpType)
//	fmt.Println("rtmpType:", rtmpType)
//
//	if rtmpType == "" {
//		return "", errors.New("rtmpType不能为空")
//	}
//
//	return rtmpType, nil
//}

//
//func showMenu() (*RtmpData, error) {
//	// the questions to ask
//	var qs = []*survey.Question{
//		{
//			Name:   "rtmpurl",
//			Prompt: &survey.Input{Message: "请输入推流地址(服务器地址+串流密钥)"},
//		},
//		{
//			Name:   "videolist",
//			Prompt: &survey.Input{Message: "请输入需要推流的文件夹或文件路径"},
//		},
//		{
//			Name: "showtitle",
//			Prompt: &survey.Select{
//				Message: "是否在影片中加入文件名水印:",
//				Options: []string{
//					"是",
//					"否",
//				},
//				Default: "是",
//			},
//		},
//	}
//
//	var answers = struct {
//		RtmpUrl   string `survey:"rtmpurl"`   // 或者你也可以用tag指定如何匹配
//		VideoList string `survey:"videolist"` // 如果类型不一致，survey会尝试转换
//		ShowTitle string `survey:"showtitle"` // 显示标题
//	}{}
//
//	// 执行提问
//	err := survey.Ask(qs, &answers)
//	if err != nil {
//		return nil, err
//	}
//
//	rtmpUrl := strings.TrimSpace(answers.RtmpUrl)
//	videoList := strings.TrimSpace(answers.VideoList)
//	showTitle := strings.TrimSpace(answers.ShowTitle)
//
//	fmt.Println("streamKey:", rtmpUrl)
//	fmt.Println("videoList:", videoList)
//
//	if rtmpUrl == "" {
//		return nil, errors.New("streamKey不能为空")
//	}
//	if videoList == "" {
//		return nil, errors.New("videoList不能为空")
//	}
//
//	rtmpData := &RtmpData{
//		VideoList: videoList,
//		RtmpUrl:   rtmpUrl,
//		ShowTitle: showTitle,
//	}
//	return rtmpData, nil
//}
