package common

import (
	"path/filepath"
)

// 检查fileName文件类型是否在指定列表extList中
func CheckFileTypes(fileName string, extList []string) bool {
	fileType := filepath.Ext(fileName)
	for _, ext := range extList {
		if fileType == ext {
			return true
		}
	}
	return false
}
