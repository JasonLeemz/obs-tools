package tools

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/BurntSushi/toml"
	_ "github.com/BurntSushi/toml"
	"os"
	"strings"
)

var mapRtpmPrefix = map[string]string{
	"快手":   "kuaishou",
	"哔哩哔哩": "bilibili",
}

type RtmpData struct {
	VideoList string `json:"video_list,omitempty"`
	RtmpUrl   string `json:"rtmp_url,omitempty"`
	ShowTitle string `json:"show_title,omitempty"`
	Subtitle  string `json:"subtitle,omitempty"`
}

func showMenu() (*RtmpData, error) {
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
		{
			Name: "showtitle",
			Prompt: &survey.Select{
				Message: "是否在影片中加入文件名水印:",
				Options: []string{
					"是",
					"否",
				},
				Default: "是",
			},
		},
		//{
		//	Name: "subtitle",
		//	Prompt: &survey.Select{
		//		Message: "是否将字幕文件合成进入影片,默认检索同名文件名并且扩展名为.srt的文件:",
		//		Options: []string{
		//			"是",
		//			"否",
		//		},
		//		Default: "否",
		//	},
		//},
	}

	var answers = struct {
		RtmpType  string `survey:"rtmptype"`  // survey 会默认匹配首字母小写的name
		StreamKey string `survey:"streamkey"` // 或者你也可以用tag指定如何匹配
		VideoList string `survey:"videolist"` // 如果类型不一致，survey会尝试转换
		ShowTitle string `survey:"showtitle"` // 显示标题
		Subtitle  string `survey:"subtitle"`  // 字幕
	}{}

	// 执行提问
	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}

	rtmpType := strings.TrimSpace(answers.RtmpType)
	streamKey := strings.TrimSpace(answers.StreamKey)
	videoList := strings.TrimSpace(answers.VideoList)
	showTitle := strings.TrimSpace(answers.ShowTitle)
	subtitle := strings.TrimSpace(answers.Subtitle)

	fmt.Println("rtmpType:", rtmpType)
	fmt.Println("streamKey:", streamKey)
	fmt.Println("videoList:", videoList)

	if rtmpType == "" {
		return nil, errors.New("rtmpType不能为空")
	}
	if streamKey == "" {
		return nil, errors.New("streamKey不能为空")
	}
	if videoList == "" {
		return nil, errors.New("videoList不能为空")
	}

	rtmp := ""
	ok := false
	if rtmp, ok = mapRtpmPrefix[rtmpType]; !ok {
		return nil, errors.New("所选平台当前不支持")
	}
	rtmpUrl := rtmp + streamKey

	rtmpData := &RtmpData{
		VideoList: videoList,
		RtmpUrl:   rtmpUrl,
		ShowTitle: showTitle,
		Subtitle:  subtitle,
	}
	return rtmpData, nil
}

func switchPlatform() (string, error) {
	// the questions to ask
	var qs = []*survey.Question{
		{
			Name: "rtmptype",
			Prompt: &survey.Select{
				Message: "选择一个推流平台,回车键确认:",
				Options: []string{
					"快手",
					"哔哩哔哩",
					"全部",
				},
			},
		},
	}

	var answers = struct {
		RtmpType string `survey:"rtmptype"` // survey 会默认匹配首字母小写的name
	}{}

	// 执行提问
	err := survey.Ask(qs, &answers)
	if err != nil {
		return "", err
	}

	rtmpType := strings.TrimSpace(answers.RtmpType)
	fmt.Println("rtmpType:", rtmpType)

	if rtmpType == "" {
		return "", errors.New("rtmpType不能为空")
	}

	return rtmpType, nil
}

func Push() error {
	rtmpData := &RtmpData{}
	var err error

	// 读取配置文件，如果不存在，就显示交互式菜单
	rtmpType, err := switchPlatform()
	if err != nil {
		return err
	}
	rtdCfg, rtd, err := ReloadConfig(rtmpType)
	if rtdCfg.Enable == false {
		// 手动指定
		rtmpData, err = showMenu()
		if err != nil {
			return err
		}
	} else {
		// 从配置文件读取
		rtmpData = rtd
	}

	files, err := ListVideoFiles(rtmpData.VideoList)
	fmt.Println(files)
	if err != nil {
		return err
	}

	totalFiles := len(files)
	fmt.Println("共有", totalFiles, "个视频文件")
	movieName := ""
	for {
		curr := 0
		for _, file := range files {
			curr++
			fmt.Println(fmt.Sprintf("当前播放第%d个,文件名:%s", curr, file))
			if rtmpData.ShowTitle == "是" {
				movieName, _, err = ExtractFileNameInfo(file)
				if err != nil {
					return err
				}
			}

			err = pushStream(file, movieName, rtmpData.Subtitle, rtmpData.RtmpUrl)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pushStream(filePath, movieName, subtitle, rtmpUrl string) error {

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
		//ffmpeg -i input.mp4 -vf "drawtext=fontfile=simhei.ttf: text=技术是第一生产力:x=10:y=10:fontsize=24:fontcolor=white:shadowy=2" output.mp4
		cmdArguments = []string{
			"-re", "-i", filePath,
			"-c:v", "libx264",
			"-c:a", "copy",
			"-b:a", "192k",
			"-vf", "\"drawtext=fontfile=./resource/fonts/SourceHanSansCN-VF-2.otf: text=" + movieName + ":x=10:y=10:fontsize=10:fontcolor=white:shadowy=2\"",
			"-strict", "-2",
			"-f", "flv", rtmpUrl,
		}
	}

	//// 合成字幕
	//if subtitle == "是" && movieName != "" {
	//	subtitleFile := movieName + ".srt"
	//	cmdArguments = append(cmdArguments, "-i", subtitleFile)
	//}
	o, ok := ExecShell("ffmpeg", cmdArguments, "run.log")

	fmt.Println(o, ok)

	return nil
}

// RtmpConfig Rtmp推流配置
type RtmpConfig struct {
	Enable       bool   `toml:"enable"`
	RtmpUrl      string `toml:"rtmp_url"`
	VideoPath    string `toml:"video_path"`
	ShowTitle    bool   `toml:"show_title"`
	ShowSubtitle bool   `toml:"show_subtitle"`
}

func ReloadConfig(platform string) (*RtmpConfig, *RtmpData, error) {
	path := "./resource/config/rtmp.template.toml"
	if _, err := os.Stat(path); err != nil {
		return nil, nil, err
	}
	mRtmp := make(map[string]*RtmpConfig)
	_, err := toml.DecodeFile(path, &mRtmp)

	if err != nil {
		return nil, nil, err
	}

	cfg := &RtmpConfig{}
	ok := false
	if platform, ok = mapRtpmPrefix[platform]; !ok {
		return nil, nil, errors.New(fmt.Sprintf("[%s] is unsupport", platform))
	}

	if cfg, ok = mRtmp[platform]; !ok {
		return nil, nil, errors.New(fmt.Sprintf("[%s] is unsupport", platform))
	}
	rtd := &RtmpData{
		VideoList: cfg.VideoPath,
		RtmpUrl:   cfg.RtmpUrl,
	}

	// 标题水印
	if cfg.ShowTitle {
		rtd.ShowTitle = "是"
	}

	// 字幕
	if cfg.ShowSubtitle {
		rtd.Subtitle = "是"
	}
	return cfg, rtd, nil
}

func WatchCommand() {

}

func rePlay() {

}

func prev() {

}

func next() {

}
