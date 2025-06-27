package handler

import (
	router_handler "eshop_server/src/router/handler"
	"eshop_server/src/stream/model"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uuid"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadStreamingFile(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Infof("UploadStreamingFile 获取到文件, filename: %s, filesize: %vMB", file.Filename, file.Size/1024/1024)

	// 设置文件名（UUID）
	file_id := uuid.GetUuid()
	localFilename := file_id + filepath.Ext(file.Filename)
	
	// 保存上传的文件
	srcFile := filepath.Join(model.StreamFileUploadPath, localFilename)
	if err = c.SaveUploadedFile(file, srcFile); err != nil {
		log.Errorf("UploadStreamingFile 保存文件失败, filename: %s , error: %s", file.Filename, err.Error())
		router_handler.Fail(c, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Code, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Detail)
		return
	}

	// 生成m3u8文件和.ts分片文件
	// dstFile := filepath.Join(model.StreamFileSegmentPath, strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))+".m3u8") // 源文件名.m3u8
	dstFile := filepath.Join(model.StreamFileSegmentPath, file_id+".m3u8") // {file_id}.m3u8
	err = generateM3U8(srcFile, dstFile)
	if err != nil {
		log.Errorf("UploadStreamingFile 生成m3u8文件失败, filename: %s, error: %s", file.Filename, err.Error())
		router_handler.Fail(c, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Code, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Detail)
		return
	}

	// 返回成功响应
	dataMap["filename"] = localFilename
	dataMap["id"] = file_id
	dataMap["ext"] = strings.Split(localFilename, ".")[1]
	router_handler.Success(c, dataMap)
}

func generateM3U8(srcFile string, dstFile string) error {
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

// 提供m3u8文件和.ts分片文件的下载播放
func StreamingPlayer(c *gin.Context) {
	// var err error
	// dataMap := make(map[string]interface{})

	// 请求参数校验
	filename := c.Param("filename")
	if filename == "" {
		router_handler.FailWithFileNotFound(c)
		return
	}
	log.Infof("StreamingPlayer 请求参数, filename:%s", filename)

	// 判断文件是否存在
	filePath := filepath.Join(model.StreamFileSegmentPath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		router_handler.FailWithFileNotFound(c)
		return
	}

	// TODO 文件请求成功，日志记录，新增缓存
	// cache key: IP-用户-m3u8文件
	
	c.File(filePath)	
}
