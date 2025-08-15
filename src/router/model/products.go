package model

import (
	"time"
)

const (
	ProductStatusOff   int32 = 0 //下架
	ProductStatusOn    int32 = 1 //上架
	ProductStatusAudit int32 = 2 //待审核

	ProductSourceTypeDefault int32 = 0 // 默认
	ProductSourceTypeYlt     int32 = 1 // ylt

	ProductImageUrlDefault string = "https://ucarecdn.com/28285bd2-bfa6-46aa-af19-24e00ea396a9/-/preview/1000x562/" //默认商品图片链接
)

/*
-- @Author AInoriex
-- @Desc 用于记录商品基本信息
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
	`source_type` int(128) NULL DEFAULT 0 COMMENT '来源类别',
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
	SourceType   int32     `json:"source_type" gorm:"column:source_type;default:0;comment:'来源类别'"`
	ExternalId   string    `json:"external_id" gorm:"column:external_id;default:NULL;comment:'外部商品ID'"`
	ExternalLink string    `json:"external_link" gorm:"column:external_link;default:NULL;comment:'外部商品链接'"`
	CreateAt     time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
}

func (t *Products) TableName() string {
	return "products"
}

// @Title	创建商品请求体
// @Author  AInoriex (2025/08/11 14:20)
type CreateProductReq struct {
	Products
	PP []ProductsPlayer `json:"player_list"`
}

// @Title	用户查看商品列表格式化
// @Author  AInoriex  (2025/06/26 16:30)
type ProductUserView struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"image_url"`
	Sales       int64   `json:"sales"`
}

func (m *Products) UserViewFormat() (resList *ProductUserView) {
	resList = &ProductUserView{
		Id:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Price:       m.Price,
		ImageUrl:    m.ImageUrl,
		Sales:       m.Sales,
	}
	return
}
