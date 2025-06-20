package model

import (
	"time"
)

const ()

/*
-- @Author  AInoriex
-- @Des     用于记录用户的购买历史&商品权限表
-- @Create  2025年5月12日17点15分
CREATE TABLE purchase_history (

	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '购买历史唯一标识',
	`user_id` varchar(32) NOT NULL COMMENT '用户ID(关联用户表)',
	`product_id` varchar(16) NOT NULL COMMENT '商品ID(关联商品表)',
	`quantity` int(8) NOT NULL COMMENT '购买数量',
	`payment_id` varchar(255) NOT NULL COMMENT '支付ID(关联支付表)',
	`purchased_at` datetime DEFAULT NULL COMMENT '支付时间',
	PRIMARY KEY (`id`),
	KEY `idx_user_id` (`user_id`),
	KEY `idx_user_product_id` (`user_id`,`product_id`)

) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购买历史订单表';
*/
type PurchaseHistory struct {
	Id          int64     `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'购买历史唯一标识'"`
	UserId      string    `json:"user_id" gorm:"column:user_id;default:NULL;comment:'用户ID(关联用户表)'"`
	ProductId   string    `json:"product_id" gorm:"column:product_id;default:NULL;comment:'商品ID(关联商品表)'"`
	Quantity    int32     `json:"quantity" gorm:"column:quantity;default:NULL;comment:'购买数量'"`
	OrderId     string    `json:"order_id" gorm:"column:order_id;default:NULL;comment:'订单ID(关联订单表)'"`
	PaymentId   string    `json:"payment_id" gorm:"column:payment_id;default:NULL;comment:'支付ID(关联支付表)'"`
	PurchasedAt time.Time `json:"purchased_at" gorm:"column:purchased_at;default:NULL;comment:'支付时间'"`
}

func (t *PurchaseHistory) TableName() string {
	return "purchase_history"
}

// 用户获取历史购买记录响应体
type GetUserPurchaseHistoryResp struct {
	Id                 int64     `json:"id"`
	OrderId            string    `json:"order_id"`
	ProductName        string    `json:"product_name"`
	FinalAmount        float64   `json:"final_amount"`
	Quantity           int32     `json:"quantity"`
	PurchaseStatus     int32     `json:"-"` // 不传递给前端
	PurchaseStatusDesc string    `json:"purchase_status_desc"`
	PurchaseDate       time.Time `json:"purchase_date"`
}
