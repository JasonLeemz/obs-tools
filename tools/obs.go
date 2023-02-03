package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/JasonLeemz/obs-tools/core/log"
	"os"
	"strings"
)

func Push() error {
	logger := log.InitLogger()

	// 加载推流配置
	rtmpConfig, err := ReloadConfig()
	logger.Sugar().Debug("rtmpConfig", rtmpConfig)
	sjson, _ := json.Marshal(rtmpConfig)
	logger.Sugar().Debug("rtmpConfig=", string(sjson))

	// 读取视频文件
	files, err := ListVideoFiles(rtmpConfig.VideoPath)
	totalFiles := len(files)
	logger.Sugar().Debugf("共有 %d 个视频文件,videoFiles=%v", totalFiles, files)
	if err != nil {
		return err
	}

	// 控制播放循环次数
	loopCount := rtmpConfig.LoopCount
	// 文件(视频标题)名称
	movieName := ""
	currentLoop := int32(0)
	for {
		// 按照阅读习惯，将0初始化为第1次
		currentLoop++

		// 按照读取的顺序，开始推流
		curr := 0
		for _, file := range files {
			curr++
			logger.Sugar().Infof("当前播放第%d次循环中的第%d个,文件名:%s", currentLoop, curr, file)
			if rtmpConfig.ShowTitle == true {
				movieName, _, err = ExtractFileNameInfo(file)
				if err != nil {
					logger.Sugar().Errorf(err.Error())
					return err
				}
			}

			output, err := pushStream(file, movieName, rtmpConfig)
			if err != nil {
				logger.Sugar().Errorw("pushStream", "output", output, "error", err.Error())
				return err
			}
		}

		// 判断循环次数
		if loopCount != 0 && currentLoop >= loopCount {
			msg := fmt.Sprintf("当前已经循环播放 [%d] 次，直播完成", currentLoop)
			logger.Sugar().Infof(msg)
			break
		}
	}

	return errors.New("播放完成")
}

func pushStream(filePath, movieName string, rtmpConfig *RtmpConfig) (string, error) {
	// 初始化logger
	logger := log.InitLogger()
	logger.Sugar().Debug("filePath=", filePath, "rtmpConfig=", rtmpConfig)

	movieName = strings.TrimSpace(movieName)
	cmdArguments := make([]string, 0)

	filePath = fmt.Sprintf("%s%s%s", "\"", filePath, "\"")
	rtmpUrl := fmt.Sprintf("%s%s%s", "\"", rtmpConfig.PushUrl, "\"")

	acodec := "copy"
	if rtmpConfig.FFMpegParams.ACodec != "" {
		acodec = rtmpConfig.FFMpegParams.ACodec
	}

	vcodec := "copy"
	if rtmpConfig.FFMpegParams.VCodec != "" {
		vcodec = rtmpConfig.FFMpegParams.VCodec
	}

	fontsize := int32(0)
	if rtmpConfig.FFMpegParams.FontSize != 0 {
		fontsize = rtmpConfig.FFMpegParams.FontSize
	} else if rtmpConfig.FFMpegParams.FontSize == 0 {
		// 如果为空 或者 为0，就不显示标题
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
		fsize := fmt.Sprintf("%d", fontsize)
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
			"-vf", "\"drawtext=fontfile=./resource/fonts/SourceHanSansCN-VF-2.otf: text=" + movieName + ":x=10:y=10:fontsize=" + fsize + ":fontcolor=white:shadowy=2\"",
			"-strict", "-2",
			"-f", "flv", rtmpUrl,
		}
	}

	logger.Sugar().Info("cmdArguments=", cmdArguments)
	out, err := ExecShell("", cmdArguments)

	// 打日志
	logger.Sugar().Debugw("current file complete", "outPut=", out, "error=", err)
	return out, err
}

// Stream 推流配置
type Stream map[string][]*RtmpConfig
type RtmpConfig struct {
	Platform     string        `toml:"platform"`
	Enable       bool          `toml:"enable"`
	RtmpServer   string        `toml:"rtmp_server"`
	StreamKey    string        `toml:"stream_key"`
	VideoPath    string        `toml:"video_path"`
	ShowTitle    bool          `toml:"show_title"`
	LoopCount    int32         `toml:"loop_count"`
	FFMpegParams *FFMpegParams `toml:"ffmpeg"`
	PushUrl      string
}

// FFMpegParams ffmpeg参数
type FFMpegParams struct {
	ACodec   string `toml:"acodec"`
	VCodec   string `toml:"vcodec"`
	FontSize int32  `toml:"font_size"`
}

func ReloadConfig() (*RtmpConfig, error) {

	path := "./config/rtmp.template.toml"
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	// 解析配置文件
	stream := Stream{}
	_, err := toml.DecodeFile(path, &stream)

	if err != nil {
		return nil, err
	}

	rtmpConfig := &RtmpConfig{}
	for _, config := range stream["stream"] {
		if config.Enable == false {
			continue
		}

		rtmpConfig = config
		rtmpConfig.PushUrl = config.RtmpServer + config.StreamKey // 拼接一下

		return rtmpConfig, nil
	}

	return nil, nil
}

func WatchCommand() {

}

func rePlay() {

}

func prev() {

}

func next() {

}
