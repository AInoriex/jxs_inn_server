package model

import (
	"time"
)

const (
	ProductStatusOff int32 = 0 //下架
	ProductStatusOn  int32 = 1 //上架
)

/*
-- @Author AInoriex
-- @Desc 用于记录商品基本信息
-- @TODO 补充特定商品的属性信息(音声格式, 音声时长……)
-- @Chge 2025年5月6日11点06分 id int(11) -> varchar(16)
-- @Chge 2025年5月9日17点46分 新增字段external_id, external_link
CREATE TABLE `products` (

	`id` varchar(16) NOT NULL COMMENT '商品唯一标识',
	`title` varchar(100) NOT NULL COMMENT '商品标题',
	`description` text COMMENT '商品描述',
	`price` decimal(10,2) NOT NULL COMMENT '商品价格',
	`status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '商品状态(0下架, 1上架)',
	`image_url` varchar(255) DEFAULT NULL COMMENT '商品图片URL',
	`sales` int(11) DEFAULT '0' COMMENT '商品销量',
	`external_id` varchar(64) DEFAULT NULL COMMENT '外部商品ID',
	`external_link` varchar(64) DEFAULT NULL COMMENT '外部商品链接',
	`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	PRIMARY KEY (`id`),
	KEY `idx_title` (`title`),
	KEY `idx_price` (`price`),
	KEY `idx_status` (`status`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';
*/
type Products struct {
	Id           string    `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'商品唯一标识'"`
	Title        string    `json:"title" gorm:"column:title;default:NULL;comment:'商品标题'"`
	Description  string    `json:"description" gorm:"column:description;default:NULL;comment:'商品描述'"`
	Price        float64   `json:"price" gorm:"column:price;default:NULL;comment:'商品价格'"`
	Status       int32     `json:"status" gorm:"column:status;default:1;comment:'商品状态(0下架, 1上架)'"`
	ImageUrl     string    `json:"image_url" gorm:"column:image_url;default:NULL;comment:'商品图片URL'"`
	Sales        int64     `json:"sales" gorm:"column:sales;default:0;comment:'商品销量'"`
	ExternalId   string    `json:"external_id" gorm:"column:external_id;default:NULL;comment:'外部商品ID'"`
	ExternalLink string    `json:"external_link" gorm:"column:external_link;default:NULL;comment:'外部商品链接'"`
	CreateTime   time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
}

func (t *Products) TableName() string {
	return "products"
}
