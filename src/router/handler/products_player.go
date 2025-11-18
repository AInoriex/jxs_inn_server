package handler

import (
	"encoding/json"
	"eshop_server/src/common/api"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	"eshop_server/src/utils/common"
	"eshop_server/src/utils/config"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uuid"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Title		 上传流媒体文件
// @Description  上传流媒体文件到服务器
// @Response     json
// @Router       /v1/eshop_api/admin/player/upload_streaming_file [post]
func UploadStreamingFile(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})

	// 从请求中获取product_id
	product_id := c.PostForm("product_id")
	if product_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "necessary parameter is missing"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload failed"})
		return
	}

	// 检查文件类型
	if !common.CheckFileTypes(file.Filename, []string{".mp3", ".wav"}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type must be mp3 or wav"})
		return
	}
	fileType := filepath.Ext(file.Filename)
	log.Infof("UploadStreamingFile 获取到文件, filename: %s, filesize: %vMB", file.Filename, file.Size/1024/1024)

	// 创建数据库记录
	player := &model.ProductsPlayer{
		Id:        uuid.GetUuid(),
		ProductId: product_id,
		Filename:  strings.TrimSuffix(file.Filename, fileType), // 只保留文件名不包含后缀
		FileType:  fileType,
		Status:    model.ProductsPlayerStatusParsing,
	}
	_, err = dao.CreateProductsPlayer(player)
	if err != nil || player == nil {
		log.Errorf("UploadStreamingFile 创建ProductsPlayer记录失败, filename: %s , error: %s", file.Filename, err.Error())
		api.Fail(c, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Code, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Detail)
		return
	}

	// 修改 defer 函数，通过闭包捕获 err 变量地址
	defer func() {
		if err != nil {
			// 更新ProductsPlayer记录为错误状态
			player.Status = model.ProductsPlayerStatusError
			if _, updateErr := dao.UpdateProductsPlayerByField(player, []string{"status"}); updateErr != nil {
				log.Errorf("UploadStreamingFile 更新ProductsPlayer记录失败, player: %+v , error: %s", player, updateErr.Error())
			}
			api.Fail(c, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Code, uerrors.Parse(uerrors.ErrorStreamFileUploadFailed.Error()).Detail)
			return
		}
	}()

	// 调用stream服务上传文件获取player信息
	request_url := config.StreamConfig.Host + "/v1/steaming/internal_upload_streaming_file"
	streamPlayerInfo, err := RequestInternalUploadStreamingFile(request_url, c)
	if err != nil {
		log.Errorf("UploadStreamingFile 调用stream服务上传文件接口失败, error: %s", err.Error())
		return
	}

	// 更新player信息
	// player.PlayType = model.ProductsPlayerPlayTypeHls
	// player.PlayUrl = config.StreamConfig.Host + player.Id + ".m3u8"
	player.Duration = streamPlayerInfo.Duration
	player.FileSize = streamPlayerInfo.FileSize
	player.PlayType = streamPlayerInfo.PlayType
	player.PlayUrl = streamPlayerInfo.PlayUrl
	player.Status = model.ProductsPlayerStatusOk
	product_player_update_fields := []string{"duration", "file_size", "play_type", "play_url", "status"}
	if _, err = dao.UpdateProductsPlayerByField(player, product_player_update_fields); err != nil {
		log.Errorf("UploadStreamingFile 更新ProductsPlayer记录失败, m: %+v , error: %s", player, err.Error())
		return
	}

	// 返回成功响应
	dataMap["result"] = player
	api.Success(c, dataMap)
}

// 转发请求，请求流媒体服务器上传文件获取文件meta信息
func RequestInternalUploadStreamingFile(url string, c *gin.Context) (m model.ProductsPlayer, err error) {
	log.Infof("RequestInternalUploadStreamingFile request_url: %s", url)
	// 发送POST请求
	resp, err := http.Post(url, "multipart/form-data", c.Request.Body)
	if err != nil {
		log.Errorf("RequestInternalUploadStreamingFile 发送POST请求失败, error: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("RequestInternalUploadStreamingFile 读取响应体失败, error: %s", err.Error())
		return
	}
	// 解析响应体
	if err = json.Unmarshal(body, &m); err != nil {
		log.Errorf("RequestInternalUploadStreamingFile 解析响应体失败, error: %s", err.Error())
		return
	}
	return
}
