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
