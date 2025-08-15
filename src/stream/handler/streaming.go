package handler

import (
	"bytes"
	router_dao "eshop_server/src/router/dao"
	router_handler "eshop_server/src/router/handler"
	router_model "eshop_server/src/router/model"
	"eshop_server/src/stream/model"
	"eshop_server/src/utils/common"
	"eshop_server/src/utils/config"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 上传文件
func UploadStreamingFile(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})

	// ~~设置文件大小限制~~
	// 放在 Nginx 限制
	// c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8<<20) // 8 MiB 

	// 从请求中获取product_id 
	product_id := c.PostForm("product_id") 

	// 获取上传的文件 
	file, err := c.FormFile("file") 
	if err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload failed"}) 
		return 
	} 

	// 验证文件真实类型 
	// fileHeader, _ := file.Open() 
	// defer fileHeader.Close() 
	// buffer := make([]byte, 12) 
	// _, err = fileHeader.Read(buffer) 
	// if err != nil { 
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "file read failed"}) 
	// 	return 
	// } 
	// // 魔数校验示例，需根据实际情况完善 
	// isValid := validateFileType(buffer) 
	// if !isValid { 
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "file type must be mp3 or wav"}) 
	// 	return 
	// }
	if !common.CheckFileTypes(file.Filename, []string{".mp3", ".wav"}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type must be mp3 or wav"}) 
		return 
	}
	fileType := filepath.Ext(file.Filename)
	log.Infof("UploadStreamingFile 获取到文件, filename: %s, filesize: %vMB", file.Filename, file.Size/1024/1024)

	// 创建数据库记录
	player := &router_model.ProductsPlayer{ 
		Id:        uuid.GetUuid(), 
		ProductId: product_id, 
		Filename:  strings.TrimSuffix(file.Filename, fileType), // 只保留文件名不包含后缀 
		FileType:  fileType, 
		FileSize:  file.Size, 
		Status:    router_model.ProductsPlayerStatusInit, 
		Duration:  0, 
		PlayType:  "", 
		PlayUrl:   "", 
	} 
	_, err = router_dao.CreateProductsPlayer(player)
	if err != nil  || player == nil { 
		log.Errorf("UploadStreamingFile 创建ProductsPlayer记录失败, filename: %s , error: %s", file.Filename, err.Error()) 
		router_handler.Fail(c, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Code, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Detail) 
		return 
	} 

	// 修改 defer 函数，通过闭包捕获 err 变量地址
	defer func() { 
		if err != nil { 
			player.Status = router_model.ProductsPlayerStatusError // 假设存在失败状态
			// 更新数据库记录
			if _, updateErr := router_dao.UpdateProductsPlayerByField(player, []string{"status"}); updateErr != nil { 
				log.Errorf("UploadStreamingFile 更新ProductsPlayer记录失败, player: %+v , error: %s", player, updateErr.Error()) 
			} 
			router_handler.Fail(c, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Code, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Detail)
			return
		} 
	}() 

	// 重新设置文件名
	newFilename := player.Id + fileType
	
	// 保存上传的文件
	srcFile := filepath.Join(model.StreamFileUploadPath, newFilename)
	if err = c.SaveUploadedFile(file, srcFile); err != nil {
		log.Errorf("UploadStreamingFile 保存文件失败, filename: %s , error: %s", file.Filename, err.Error())
		return
	}
	
	// 更新数据库
	player.Filename = player.Id
	player.Duration, err = GetMediaDuration(srcFile)
	if err != nil {
		log.Errorf("UploadStreamingFile 获取文件时长失败, srcFile: %s, error: %s", srcFile, err.Error())
		return
	}
	player.Status = router_model.ProductsPlayerStatusParsing
	if _, err = router_dao.UpdateProductsPlayerByField(player, []string{"filename", "duration", "status"}); err != nil {
		log.Errorf("UploadStreamingFile 更新ProductsPlayer记录失败, m: %+v , error: %s", player, err.Error())
		return
	}

	// 生成m3u8文件和.ts分片文件
	m3u8File := filepath.Join(model.StreamFileSegmentPath, player.Id+".m3u8") // {file_id}.m3u8
	err = GenerateAudioM3u8(srcFile, m3u8File)
	if err != nil {
		log.Errorf("UploadStreamingFile 生成m3u8文件失败, srcFile: %s, m3u8File: %s, error: %s", srcFile, m3u8File, err.Error())
		return
	}

	// 更新数据库
	player.PlayType = router_model.ProductsPlayerPlayTypeHls
	player.PlayUrl = config.StreamConfig.Host + player.Id + ".m3u8"
	player.Status = router_model.ProductsPlayerStatusOk
	if _, err = router_dao.UpdateProductsPlayerByField(player, []string{"play_type", "play_url", "status"}); err != nil {
		log.Errorf("UploadStreamingFile 更新ProductsPlayer记录失败, m: %+v , error: %s", player, err.Error())
		return
	}

	// 返回成功响应
	// dataMap["filename"] = newFilename
	// dataMap["id"] = product_id
	// dataMap["ext"] = strings.Split(newFilename, ".")[1]
	dataMap["result"] = player
	router_handler.Success(c, dataMap)
}

// 验证文件类型的魔数校验函数(暂不使用)
func validateFileType(buffer []byte) bool {
	println("DEBUG validateFileType buffer: ", buffer)
    // mp3 文件魔数 
    if bytes.HasPrefix(buffer, []byte{0x49, 0x44, 0x33}) { 
        return true 
    } 
    // wav 文件魔数 
    if bytes.HasPrefix(buffer, []byte{0x52, 0x49, 0x46, 0x46}) { 
        return true 
    } 
	// m4a 文件魔数
	if bytes.HasPrefix(buffer, []byte{0x4D, 0x41, 0x44, 0x46}) {
		return true
	}
	// mpeg 文件魔数
	if bytes.HasPrefix(buffer, []byte{0x4D, 0x50, 0x45, 0x47}) {
		return true
	}
    return false 
}

// 提供m3u8文件和.ts分片文件的下载播放
func StreamingPlayer(c *gin.Context) {
	// var err error
	// dataMap := make(map[string]interface{})

	// 请求参数校验
	filename := c.Param("filename")
	if filename == "" {
		log.Errorf("StreamingPlayer 请求参数filename为空")
		router_handler.FailWithFileNotFound(c)
		return
	}

	// 检查请求文件类型是否为m3u8或ts
	if !common.CheckFileTypes(filename, []string{".m3u8", ".ts"}) {
		log.Errorf("StreamingPlayer 请求参数格式错误, filename:%s", filename)
		router_handler.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":格式错误")

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
