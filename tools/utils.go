package tools

import (
	"errors"
	"fmt"
	"github.com/JasonLeemz/obs-tools/core/log"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

var videoFileType = map[string]bool{
	"wmv":  true,
	"asf":  true,
	"asx":  true,
	"rm":   true,
	"rmvb": true,
	"mpg":  true,
	"mpeg": true,
	"mpe":  true,
	"3gp":  true,
	"mov":  true,
	"mp4":  true,
	"m4v":  true,
	"avi":  true,
	"dat":  true,
	"mkv":  true,
	"flv":  true,
	"vob":  true,
}

// ListVideoFiles 返回所给路径中所有媒体文件
func ListVideoFiles(root string) ([]string, error) {

	files := make([]string, 0)
	// 如果是单个文件直接返回
	if !IsDir(root) {
		fmt.Println("当前文件[", root, "]不是文件夹")
		// 校验文件合法性
		_, t, err := ExtractFileNameInfo(root)
		if err != nil {
			return nil, err
		}
		if _, ok := videoFileType[t]; ok {
			files = append(files, root)
		}

		return files, nil
	}
	logger := log.InitLogger()
	logger.Sugar().Debugf("开始遍历[%s]", root)

	// 遍历文件夹
	return walkDir(root)
}

// 遍历文件夹
func walkDir(root string) ([]string, error) {
	_, files, err := FileForEachComplete(root)
	return files, err
}

func FileForEachComplete(fileFullPath string) ([]fs.FileInfo, []string, error) {
	files, err := ioutil.ReadDir(fileFullPath)
	if err != nil {
		return nil, nil, err
	}
	var myFile []fs.FileInfo
	var sFiles []string
	for _, file := range files {
		if file.IsDir() {
			path := strings.TrimSuffix(fileFullPath, "/") + "/" + file.Name()
			subFile, subFp, _ := FileForEachComplete(path)
			if len(subFile) > 0 {
				myFile = append(myFile, subFile...)
			}
			if len(subFp) > 0 {
				sFiles = append(sFiles, subFp...)
			}
		} else {
			fn := file.Name()

			// 如果文件名最前面第一个字符是"." 则跳过
			// 隐藏文件不记录
			rfn := []rune(fn)
			if string(rfn[0]) == "." {
				continue
			}

			strArr := strings.Split(fn, ".")
			l := len(strArr)
			if _, ok := videoFileType[strArr[l-1]]; ok {
				myFile = append(myFile, file)
				fp := strings.TrimSuffix(fileFullPath, "/") + "/" + fn
				sFiles = append(sFiles, fp)
			}

		}
	}
	return myFile, sFiles, nil
}

// ExtractFileNameInfo 提取文件的名字和扩展名
func ExtractFileNameInfo(path string) (string, string, error) {

	path = strings.TrimSpace(path)
	if path == "" {
		return "", "", errors.New("file path is empty")
	}

	strArr := strings.Split(path, "/")
	l := len(strArr)
	fileName := strArr[l-1]

	// 提取出文件名(移除后缀扩展名)
	fileNameCpy := []rune(fileName)
	fl := len(fileNameCpy)
	dotPosition := -1
	for i := fl - 1; i >= 0; i-- {
		if string(fileNameCpy[i]) == "." {
			dotPosition = i
			break
		}
	}

	movieName := ""
	movieType := ""
	if dotPosition == -1 {
		movieName = fileName
	} else {
		for i := 0; i < fl; i++ {
			if i < dotPosition {
				movieName += string(fileNameCpy[i])
			} else if i > dotPosition {
				movieType += string(fileNameCpy[i])
			}
		}
	}

	return movieName, movieType, nil
}
