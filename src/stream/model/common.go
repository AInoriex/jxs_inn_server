package model

import (
	"os"
	"path/filepath"
)

const (
	StreamFileUploadPath string = "uploads" // 音频上传路径
	StreamFileSegmentPath string = "segments" // 音频切片路径
)

// 初始化路径
func InitStreamingPaths() error {
	if err := os.MkdirAll(filepath.Dir(StreamFileUploadPath), 0750); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(StreamFileSegmentPath), 0750); err != nil {
		return err
	}
	return nil
}