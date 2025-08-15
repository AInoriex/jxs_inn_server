package model

import (
	"time"
)

const (
	ProductsPlayerFileTypeMp3 = ".mp3" // 文件类型 mp3
	ProductsPlayerFileTypeWav = ".wav" // 文件类型 wav

	ProductsPlayerPlayTypeHls = "hls" // 播放类型 hls

	ProductsPlayerStatusInit    = 0  // 文件播放状态 0:初始化
	ProductsPlayerStatusOk      = 1  // 文件播放状态 1:就绪
	ProductsPlayerStatusParsing = 2  // 文件播放状态 2:解析中
	ProductsPlayerStatusInvalid = -1 // 文件播放状态 -1:下架
	ProductsPlayerStatusError   = -2 // 文件播放状态 -2:异常
)

// 商品播放信息
type ProductsPlayer struct {
	Id        string    `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'播放id'"`
	ProductId string    `json:"product_id" gorm:"column:product_id;NOT NULL;comment:'商品id'"`
	Filename  string    `json:"filename" gorm:"column:filename;default:'';comment:'文件名'"`
	FileType  string    `json:"file_type" gorm:"column:file_type;default:'';comment:'文件类型'"`
	FileSize  int64     `json:"file_size" gorm:"column:file_size;default:0;comment:'文件大小（Byte字节）'"`
	Duration  int64     `json:"duration" gorm:"column:duration;default:0;comment:'文件时长'"`
	PlayType  string    `json:"play_type" gorm:"column:play_type;default:'';comment:'播放类型'"`
	PlayUrl   string    `json:"play_url" gorm:"column:play_url;default:'';comment:'播放地址'"`
	Status    int32     `json:"status" gorm:"column:status;default:0;comment:'文件状态'"`
	CreateAt  time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdateAt  time.Time `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:'更新时间'"`
}

func (t *ProductsPlayer) TableName() string {
	return "products_player"
}
