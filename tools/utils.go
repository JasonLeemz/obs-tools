package tools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

// ListVideoFiles 返回所给路径中所有媒体文件
func ListVideoFiles(root string) ([]string, error) {

	fileType := map[string]bool{
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
	files := make([]string, 0)
	// 如果是单个文件直接返回
	if !IsDir(root) {
		fmt.Println("当前文件[", root, "]不是文件夹")
		// 校验文件合法性
		_, t, err := ExtractFileNameInfo(root)
		if err != nil {
			return nil, err
		}
		if _, ok := fileType[t]; ok {
			files = append(files, root)
		}

		return files, nil
	}

	fmt.Println("开始遍历[", root, "]")
	// 遍历文件夹
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		strArr := strings.Split(path, ".")
		fmt.Println(path, strArr)
		l := len(strArr)
		if l < 2 {
			return nil
		}
		fmt.Println(strArr[l-2])
		if _, ok := fileType[strArr[l-2]]; ok {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// ExtractFileNameInfo 提取文件的名字和扩展名
func ExtractFileNameInfo(path string) (string, string, error) {

	strArr := strings.Split(path, "/")
	l := len(strArr)
	if l < 1 {
		return "", "", errors.New(fmt.Sprintf("file path error:%s", path))
	}

	fileName := strArr[l-1]
	fArr := strings.Split(fileName, ".")
	if len(fArr) < 2 {
		return "", "", errors.New(fmt.Sprintf("fileName error:%s", fileName))
	}
	movieName := fArr[0]
	movieType := fArr[1]
	fmt.Println(movieName)
	return movieName, movieType, nil
}
