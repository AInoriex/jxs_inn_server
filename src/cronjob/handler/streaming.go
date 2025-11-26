package handler

import (
	router_dao "eshop_server/src/router/dao"
	router_model "eshop_server/src/router/model"
	"eshop_server/src/stream/model"
	"eshop_server/src/utils/log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// @Title		定时任务清理文件
// @Description	清理管理后台上传的迷路音频文件
func CleanStreamingLostFiles() {
	var err error
	var totalDeleted int
	startTime := time.Now()
	log.Infof("=-=-=-= CleanStreamingLostFiles 开始执行定时任务清理文件, 检索上传目录: %s, 检索切片目录: %s =-=-=-=", model.StreamFileUploadPath, model.StreamFileSegmentPath)
	defer func() {
		log.Infof("=-=-=-= CleanStreamingLostFiles 定时任务清理文件执行完成，清理了%d个文件, 耗时: %v秒 =-=-=-=", totalDeleted, time.Since(startTime).Seconds())
	}()

	// 获取上传目录下所有音频文件
	uploadFiles, err := getDirectorySteamingFiles(model.StreamFileUploadPath)
	if err != nil {
		log.Errorf("获取上传目录下所有音频文件失败: %v", err)
		return
	}

	// 获取切片目录下所有音频文件
	segmentFiles, err := getDirectorySteamingFiles(model.StreamFileSegmentPath)
	if err != nil {
		log.Errorf("获取切片目录下所有音频文件失败: %v", err)
		return
	}

	// 聚合所有文件, 按player_id分组
	fileMap := make(map[string][]string)

	// 处理上传文件
	for _, filename := range uploadFiles {
		playerID := extractPlayerIDFromFilename(filename)
		fullPath := filepath.Join(model.StreamFileUploadPath, filename)
		fileMap[playerID] = append(fileMap[playerID], fullPath)
	}

	// 处理切片文件
	for _, filename := range segmentFiles {
		playerID := extractPlayerIDFromFilename(filename)
		fullPath := filepath.Join(model.StreamFileSegmentPath, filename)
		fileMap[playerID] = append(fileMap[playerID], fullPath)
	}

	// 查询数据库所有products_player记录
	players, err := router_dao.GetAllProductsPlayer()
	if err != nil {
		log.Errorf("查询数据库所有products_player记录失败: %v", err)
		return
	}

	// 创建已存在的player_id集合
	existingPlayerIDs := make(map[string]bool)
	for _, player := range players {
		existingPlayerIDs[player.Id] = true
	}

	// 遍历fileMap, 清理迷路文件
	for playerID, files := range fileMap {
		// 如果player_id不存在于数据库中, 则移动到 model.StreamFileRecyclePath 目录
		if !existingPlayerIDs[playerID] {
			log.Infof("检索到迷路文件, player_id: %s, 文件数量: %d, 文件路径: %v", playerID, len(files), files)
			for _, file := range files {
				// err := os.Remove(file)
				recyclePath := filepath.Join(model.StreamFileRecyclePath, filepath.Base(file))
				if err := os.Rename(file, recyclePath); err != nil {
					log.Errorf("迁移迷路文件失败: %v, 文件路径: %s", err, file)
				} else {
					log.Infof("文件已迁移到回收站: %s -> %s", file, recyclePath)
					totalDeleted++
				}
			}
		}
	}
}

// 判断文件是否为流媒体文件
func isStreamingFile(filename string) bool {
	// 文件以mp3/wav/m3u8/ts结尾
	format_list := append(router_model.ProductPlayerSupportFileTypeList, ".m3u8", ".ts")
	for _, format := range format_list {
		if strings.HasSuffix(filename, format) {
			return true
		}
	}
	return false
}

// 从文件路径提取player_id
// ./uploads/2c6443db0a3c4c18aa3eeb4c8775.mp3 -> 2c6443db0a3c4c18aa3eeb4c8775
func extractPlayerIDFromFilename(filename string) (playerID string) {
	subStr := strings.TrimSuffix(filename, filepath.Ext(filename))
	subStr = strings.TrimPrefix(subStr, ".")
	playerID = filepath.Base(subStr)
	return playerID
}

// 获取目录下所有音频文件
func getDirectorySteamingFiles(dir string) (files []string, err error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		if isStreamingFile(entry.Name()) {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}
