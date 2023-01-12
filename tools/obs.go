package tools

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"strings"
)

var mapRtpmPrefix = map[string]string{
	"快手":       "rtmp://open-push.voip.yximgs.com/gifshow/",
	"Bilibili": "rtmp://live-push.bilivideo.com/live-bvc/",
}

func showMenu() (string, string, error) {
	// the questions to ask
	var qs = []*survey.Question{
		{
			Name: "rtmptype",
			Prompt: &survey.Select{
				Message: "选择一个推流平台,回车键确认:",
				Options: []string{
					"快手",
					"Bilibili",
				},
			},
		},
		{
			Name:   "streamkey",
			Prompt: &survey.Input{Message: "请输入直播码/串流密钥(Stream Key)"},
		},
		{
			Name:   "videolist",
			Prompt: &survey.Input{Message: "请输入需要推流的文件夹或文件路径"},
		},
	}

	var answers = struct {
		RtmpType  string `survey:"rtmptype"`  // survey 会默认匹配首字母小写的name
		StreamKey string `survey:"streamkey"` // 或者你也可以用tag指定如何匹配
		VideoList string `survey:"videolist"` // 如果类型不一致，survey会尝试转换
	}{}

	// 执行提问
	err := survey.Ask(qs, &answers)
	if err != nil {
		return "", "", err
	}

	rtmpType := strings.TrimSpace(answers.RtmpType)
	streamKey := strings.TrimSpace(answers.StreamKey)
	videoList := strings.TrimSpace(answers.VideoList)

	fmt.Println("rtmpType:", rtmpType)
	fmt.Println("streamKey:", streamKey)
	fmt.Println("videoList:", videoList)

	if rtmpType == "" {
		return "", "", errors.New("rtmpType不能为空")
	}
	if streamKey == "" {
		return "", "", errors.New("streamKey不能为空")
	}
	if videoList == "" {
		return "", "", errors.New("videoList不能为空")
	}

	rtmp := ""
	ok := false
	if rtmp, ok = mapRtpmPrefix[rtmpType]; !ok {
		return "", "", errors.New("所选平台当前不支持")
	}
	rtmpUrl := rtmp + streamKey
	return videoList, rtmpUrl, nil
}

func Push() error {
	videoList, rtmpUrl, err := showMenu()
	if err != nil {
		return err
	}
	files, err := ListVideoFiles(videoList)
	if err != nil {
		return err
	}

	totalFiles := len(files)
	fmt.Println("共有", totalFiles, "个视频文件")

	for {
		curr := 0
		for _, file := range files {
			curr++
			fmt.Println(fmt.Sprintf("当前播放第%d个,文件名:%s", curr, file))
			movieName, _, err := ExtractFileNameInfo(file)
			if err != nil {
				return err
			}
			err = pushStream(file, movieName, rtmpUrl)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pushStream(filePath, movieName, rtmpUrl string) error {

	movieName = strings.TrimSpace(movieName)
	//movieName = ""
	cmdArguments := make([]string, 0)

	filePath = fmt.Sprintf("%s%s%s", "\"", filePath, "\"")
	rtmpUrl = fmt.Sprintf("%s%s%s", "\"", rtmpUrl, "\"")
	if movieName == "" {
		//ffmpeg -re -i "mhls1.mp4" -c:v copy -c:a copy -b:a 192k -strict -2 -f flv "rtmp://live-push.bilivideo.com/live-bvc/?streamname=xxx"
		cmdArguments = []string{
			"-re", "-i", filePath,
			"-c:v", "copy",
			"-c:a", "copy",
			"-b:a", "192k",
			"-strict", "-2",
			"-f", "flv", rtmpUrl,
		}
	} else {
		//ffmpeg -i input.mp4 -vf "drawtext=fontfile=simhei.ttf: text=‘技术是第一生产力’:x=10:y=10:fontsize=24:fontcolor=white:shadowy=2" output.mp4
		cmdArguments = []string{
			"-re", "-i", filePath,
			"-c:v", "libx264",
			"-c:a", "copy",
			"-b:a", "192k",
			"-vf", "\"drawtext=fontfile=./resource/SourceHanSansCN-VF-2.otf: text=" + movieName + ":x=10:y=10:fontsize=10:fontcolor=white:shadowy=2\"",
			"-strict", "-2",
			"-f", "flv", rtmpUrl,
		}
	}

	o, ok := ExecShell("ffmpeg", cmdArguments, "run.log")

	fmt.Println(o, ok)

	return nil
}
