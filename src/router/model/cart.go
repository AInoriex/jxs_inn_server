package model

import (
	"time"
)

const ()

/*
-- @Author AInoriex
-- @Desc 用于记录用户购物车信息
-- @Chge 2025年5月5日16点38分 取消外键users(id), products(id)
-- @Chge 2025年5月6日14点46分 新增唯一键值user_id, product_id
-- @TODO 如果音乐作品有不同的版本(如普通版、高清版、无损版), 可以在购物车表中增加version字段, 记录用户选择的商品版本。
CREATE TABLE cart_items (
    `id` int(11) AUTO_INCREMENT COMMENT '购物车项目唯一标识',
    `user_id` varchar(32) NOT NULL COMMENT '用户ID(关联用户表)',
    `product_id` varchar(16) NOT NULL COMMENT '商品ID(关联商品表)',
    `quantity` int(5) NOT NULL DEFAULT 1 COMMENT '购买数量',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY `user_product` (`user_id`,`product_id`), -- 用户与商品绑定唯一键值
    -- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    -- FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_user_product (user_id, product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车表';
*/

type CartItem struct {
	Id        int64     `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;"`                //购物车物品唯一标识
	UserId    string    `json:"user_id" gorm:"column:user_id;NOT NULL;"`                        //用户ID(关联用户表)
	ProductId string    `json:"product_id" gorm:"column:product_id;NOT NULL;"`                  //商品ID(关联商品表)
	Quantity  int32     `json:"quantity" gorm:"column:quantity;NOT NULL;"`                      //购买数量
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;"` //创建时间
}

func (t *CartItem) TableName() string {
	return "cart_items"
}

// 获取购物车列表单个商品响应结构体
type GetCartListItemResponse struct {
	Id       string  `json:"id"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
	Image    string  `json:"image"`
}

// 创建购物车请求结构体
type CreateCartItemReq struct {
	UserId    string `json:"user_id"`    //用户ID(关联用户表)
	ProductId string `json:"product_id"` //商品ID(关联商品表)
	Quantity  int32  `json:"quantity"`   //购买数量
}

// 移除购物车物品请求结构体
type RemoveCartItemReq struct {
	ProductId string `json:"product_id"` //商品ID
}
