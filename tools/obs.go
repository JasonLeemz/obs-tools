package tools

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/BurntSushi/toml"
	"github.com/JasonLeemz/obs-tools/core/log"
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
	LoopCount int32  `json:"loop_count,omitempty"`
}

func showMenu() (*RtmpData, error) {
	// the questions to ask
	var qs = []*survey.Question{
		{
			Name:   "rtmpurl",
			Prompt: &survey.Input{Message: "请输入推流地址(服务器地址+串流密钥)"},
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
	}

	var answers = struct {
		RtmpUrl   string `survey:"rtmpurl"`   // 或者你也可以用tag指定如何匹配
		VideoList string `survey:"videolist"` // 如果类型不一致，survey会尝试转换
		ShowTitle string `survey:"showtitle"` // 显示标题
	}{}

	// 执行提问
	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}

	rtmpUrl := strings.TrimSpace(answers.RtmpUrl)
	videoList := strings.TrimSpace(answers.VideoList)
	showTitle := strings.TrimSpace(answers.ShowTitle)

	fmt.Println("streamKey:", rtmpUrl)
	fmt.Println("videoList:", videoList)

	if rtmpUrl == "" {
		return nil, errors.New("streamKey不能为空")
	}
	if videoList == "" {
		return nil, errors.New("videoList不能为空")
	}

	rtmpData := &RtmpData{
		VideoList: videoList,
		RtmpUrl:   rtmpUrl,
		ShowTitle: showTitle,
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
					//"全部",
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
	rtmpCfg, rtd, err := ReloadConfig(rtmpType)
	if rtmpCfg.Enable == false {
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

	logger := log.InitLogger()
	logger.Sugar().Infof("共有 %d 个视频文件", totalFiles)

	// 读取配置，控制循环次数
	loopCount := rtmpData.LoopCount
	movieName := ""
	for {
		curr := int32(0)
		for _, file := range files {
			curr++
			logger.Sugar().Infof("当前播放第%d个,文件名:%s", curr, file)
			if rtmpData.ShowTitle == "是" {
				movieName, _, err = ExtractFileNameInfo(file)
				if err != nil {
					logger.Sugar().Errorf(err.Error())
					return err
				}
			}

			output, err := pushStream(file, movieName, rtmpData, rtmpCfg)
			if err != nil {
				logger.Sugar().Errorw("pushStream", "output", output, "error", err.Error())
				return err
			}
		}

		// 判断循环次数
		if loopCount != 0 && curr >= loopCount {
			msg := fmt.Sprintf("当前已经循环播放 [%d] 次，直播完成", curr)
			logger.Sugar().Infof(msg)
			break
		}
	}

	return errors.New("播放完成")
}

func pushStream(filePath, movieName string, rtmpData *RtmpData, rtmpConfig *RtmpConfig) (string, error) {
	// 初始化logger
	logger := log.InitLogger()
	logger.Sugar().Debug("filePath=", filePath, "rtmpData=", rtmpData, "rtmpConfig=", rtmpConfig)

	movieName = strings.TrimSpace(movieName)
	cmdArguments := make([]string, 0)

	filePath = fmt.Sprintf("%s%s%s", "\"", filePath, "\"")
	rtmpUrl := fmt.Sprintf("%s%s%s", "\"", rtmpData.RtmpUrl, "\"")

	acodec := "copy"
	if rtmpConfig.FFMpegParams.ACodec != "" {
		acodec = rtmpConfig.FFMpegParams.ACodec
	}

	vcodec := "copy"
	if rtmpConfig.FFMpegParams.VCodec != "" {
		vcodec = rtmpConfig.FFMpegParams.VCodec
	}

	fontsize := ""
	if rtmpConfig.FFMpegParams.FontSize != "" {
		fontsize = rtmpConfig.FFMpegParams.FontSize
	} else if rtmpConfig.FFMpegParams.FontSize == "" || rtmpConfig.FFMpegParams.FontSize == "0" {
		// 如果为空 或者 为0，就不显示标题
		// TODO 这里设置0好像不生效
		movieName = ""
	}

	if movieName == "" {
		//ffmpeg -re -i "mhls1.mp4" -c:v copy -c:a copy -b:a 192k  -strict -2 -f flv "rtmp://live-push.bilivideo.com/live-bvc/?streamname=xxx"
		cmdArguments = []string{
			"ffmpeg",
			"-re", "-i", filePath,
			"-preset", "ultrafast",
			//"-c:v", "copy",
			"-c:v", vcodec,
			"-c:a", acodec,
			"-b:a", "92k",
			"-b:v", "1500k",
			"-g", "60",
			"-strict", "-2",
			"-f", "flv", rtmpUrl,
		}
	} else {
		//ffmpeg -i input.mp4 -vf "drawtext=fontfile=simhei.ttf: text=技术是第一生产力:x=10:y=10:fontsize=24:fontcolor=white:shadowy=2" output.mp4
		cmdArguments = []string{
			"ffmpeg",
			"-re", "-i", filePath,
			"-preset", "ultrafast",
			//"-c:v", "libx264",
			"-c:v", vcodec,
			"-c:a", acodec,
			"-b:a", "92k",
			"-b:v", "1500k",
			"-g", "60",
			"-vf", "\"drawtext=fontfile=./resource/fonts/SourceHanSansCN-VF-2.otf: text=" + movieName + ":x=10:y=10:fontsize=" + fontsize + ":fontcolor=white:shadowy=2\"",
			"-strict", "-2",
			"-f", "flv", rtmpUrl,
		}
	}

	//// 合成字幕
	//if subtitle == "是" && movieName != "" {
	//	subtitleFile := movieName + ".srt"
	//	cmdArguments = append(cmdArguments, "-i", subtitleFile)
	//}

	logger.Sugar().Info("construct=", cmdArguments)
	out, err := ExecShell("", cmdArguments)

	// 打日志
	logger.Sugar().Debugw("current file complete", "outPut=", out, "error=", err)
	return out, err
}

// RtmpConfig Rtmp推流配置
type RtmpConfig struct {
	Enable       bool          `toml:"enable"`
	RtmpUrl      string        `toml:"rtmp_url"`
	VideoPath    string        `toml:"video_path"`
	ShowTitle    bool          `toml:"show_title"`
	ShowSubtitle bool          `toml:"show_subtitle"`
	Loop         int32         `toml:"loop"`
	FFMpegParams *FFMpegParams `toml:"ffmpeg"`
}

// FFMpegParams ffmpeg参数
type FFMpegParams struct {
	ACodec   string `toml:"acodec"`
	VCodec   string `toml:"vcodec"`
	FontSize string `toml:"font_size"`
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
