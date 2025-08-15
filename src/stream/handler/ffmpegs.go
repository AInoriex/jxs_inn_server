package handler

import (
	// "bufio"
	// "fmt"
	"bytes"
	"context"
	"eshop_server/src/utils/log"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GenerateAudioM3u8 生成m3u8文件
// @param srcFile 源文件路径
// @param dstFile 目标文件路径
func GenerateAudioM3u8(srcFile string, dstFile string) error {
	log.Infof("generateM3U8 解析音频文件, srcFile: %s, dstFile: %s\n", srcFile, dstFile)
	// 判断srcFile文件是否存在
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		return err
	}

	// 使用ffmpeg将音频文件切分成.ts文件，并生成m3u8文件
	// 示例命令：ffmpeg -i src -codec: copy -start_number 0 -hls_time 30 -hls_list_size 0 -f hls ./segments/output.m3u8
	log.Infof("generateM3U8 执行ffmpeg命令: ffmpeg -i %s -codec: copy -start_number 0 -hls_time 30 -hls_list_size 0 -f hls %s\n", srcFile, dstFile)
	cmd := exec.Command(
		"ffmpeg",
		"-i", srcFile, // 输入文件
		"-codec:", "copy", // 使用原始编码
		"-start_number", "0", // 分片文件编号从0开始
		"-hls_time", "30", // 每个分片文件的时间长度为30秒
		"-hls_list_size", "0", // 不限制分片文件数量
		"-f", "hls", // 输出格式为m3u8
		dstFile,
	)
	return cmd.Run()
}

// GetMediaDuration 通过调用 ffprobe 来获取音视频文件的时长
// 函数接收一个字符串参数 filePath，表示音视频文件的路径
// 函数返回两个值：一个整数表示时长（秒），一个 error 表示可能发生的错误
func GetMediaDuration(filePath string) (int64, error) { 
    // 检查 ffprobe 是否在 PATH 中
    if _, err := exec.LookPath("ffprobe"); err != nil { 
        return 0, fmt.Errorf("GetMediaDuration ffprobe not found in PATH: %v", err) 
    }
    
    // 转换路径分隔符到 Linux 风格
    filePath = strings.ReplaceAll(filePath, "\\", "/")
    
    // 设置命令超时时间（例如 10 秒）
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) 
    defer cancel() 
    
	// 执行 ffprobe 命令
	// ffprobe -v 
    cmd := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath) 
    stdout, err := cmd.Output() 
    if err != nil { 
        if exiterr, ok := err.(*exec.ExitError); ok { 
            return 0, fmt.Errorf("GetMediaDuration ffprobe execution failed: %s, stderr: %s", exiterr.Error(), string(exiterr.Stderr)) 
        } 
        return 0, fmt.Errorf("GetMediaDuration failed to execute ffprobe: %v", err) 
    } 
    
    // 去除换行符
    stdout = bytes.TrimSpace(stdout) 
    
	// 类型转换
    duration, err := strconv.ParseFloat(string(stdout), 64) 
    if err != nil { 
        return 0, fmt.Errorf("GetMediaDuration failed to parse duration: %v, output: %s", err, string(stdout)) 
    } 
    return int64(duration), nil 
}
