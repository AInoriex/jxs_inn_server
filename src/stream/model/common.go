package model

import (
	"os"
	// "path/filepath"
)

const (
	StreamFileUploadPath string = "uploads" // 音频上传路径
	StreamFileSegmentPath string = "segments" // 音频切片路径
)

// 初始化流媒体服务路径
func InitStreamingPaths() (err error) {
	if err = os.MkdirAll(StreamFileUploadPath, os.ModePerm); err != nil {
		return err
	}
	if err = os.MkdirAll(StreamFileSegmentPath, os.ModePerm); err != nil {
		return err
	}
	return nil
}