package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func BatchConvert(root, outPath string) error {

	ls, err := ListVideoFiles(root)
	if err != nil {
		return err
	}
	for i, file := range ls {
		fmt.Println("current work:[", i, "]")
		strArr := strings.Split(file, "/")
		l := len(strArr)
		fileName := strArr[l-1]
		strArr = strings.Split(fileName, ".")

		run(file, outPath, strArr[0])
	}

	return nil
}

func run(inPath, outPath, fileName string) {
	//cmdArguments := []string{"-i", "551.mp4", "-max_muxing_queue_size", "1024", "-b:v", "400k", "-crf", "25", "-s", "820*676", "506.mp4"}
	//ffmpeg -i input_filename.avi -c:v copy -c:a copy -y output_filename.mp4
	//ffmpeg -i file_example_AVI_1280_1_5MG.avi -c:a copy -c:v vp9 -b:v 100K outputVP9.mp4
	cmdArguments := []string{"-i", inPath, "-vcodec", "h264", "-acodec", "aac", "-threads", "10", "-preset", "ultrafast", "-r", "22", "-vf", "scale=1920:-2", outPath + fileName + "_h264.mp4"}
	//cmdArguments := []string{"-i", inPath, "-c:a", "copy", "-c:v", "h264", "-b:v", "100K", outPath + fileName + "_h264.mp4"}
	cmd := exec.Command("ffmpeg", cmdArguments...)
	fmt.Println(cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		fmt.Println("===")
	}
	fmt.Println("err---")
}
