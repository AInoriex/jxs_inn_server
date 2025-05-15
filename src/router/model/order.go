package model

import (
	"time"
)

const (
	OrderPaymentStatusCreate    int32 = 0 // 0 已创建
	OrderPaymentStatusToPay     int32 = 1 // 1 待支付
	OrderPaymentStatusPayed     int32 = 2 // 2 已支付
	OrderPaymentStatusTimeOut   int32 = 3 // 3 支付超时
	OrderPaymentStatusPayFail   int32 = 4 // 4 支付失败
	OrderPaymentStatusPayCancel int32 = 5 // 5 取消支付
)

/*
-- @Author AInoriex
-- @Desc 创建订单与商品关联表, 记录用户购买商品数量和订单金额, 外键关联用户ID和商品ID
-- @Desc 新增order_items表以支持一个订单与多个商品关联
-- @Desc 新增优惠券和订单状态设计
-- @Hint MySQL 8.0.16+ 版本才完全支持CHECK约束, 如果使用旧版本可能需要改用触发器实现。
-- @Hint 移除了所有列定义中的`CHECK`约束。如果需要确保`final_amount`的值等于`total_amount`-`discount`，可以在应用程序逻辑中进行验证，或者考虑使用触发器来实现这一逻辑。
-- @Chge 2025年5月5日16点24分 订单id INT -> varchar(16)
-- @Chge 2025年5月5日16点24分 取消外键users(id)
-- @Chge 2025年5月5日16点28分 新增item_id关联order_items表:订单商品信息
-- @Chge 2025年5月5日16点30分 新增payment_id关联payments表
-- @TODO 增加source字段, 记录订单来源(如网站、移动端、API等), 方便分析不同渠道的销售情况。
CREATE TABLE orders (
    `id` varchar(16) COMMENT '订单唯一标识',
    `user_id` varchar(32) NOT NULL COMMENT '用户ID(关联用户表)',
    `item_id` int(11) NOT NULL COMMENT '订单明细ID(关联订单明细表)',
    `total_amount` decimal(10, 2) NOT NULL COMMENT '订单总金额',
    `discount` decimal(10, 2) DEFAULT 0.00 COMMENT '优惠券折扣金额',
    `final_amount` decimal(10, 2) NOT NULL COMMENT '最终支付金额(总金额 - 折扣)',
    `payment_id` varchar(255) NOT NULL DEFAULT '' COMMENT '支付ID(关联支付信息表)',
    `payment_status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '支付状态(0已创建, 1待支付, 2已支付, 3支付超时, 4支付失败, 5取消支付)',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    -- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_payment_status (payment_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';
*/

type Order struct {
	Id            string    `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'订单唯一标识'"`
	UserId        string    `json:"user_id" gorm:"column:user_id;NOT NULL;comment:'用户ID(关联用户表)'"`
	ItemId        string    `json:"item_id" gorm:"column:item_id;NOT NULL;comment:'订单明细ID(关联订单明细表)'"`
	TotalAmount   float64   `json:"total_amount" gorm:"column:total_amount;NOT NULL;comment:'订单总金额'"`
	Discount      float64   `json:"discount" gorm:"column:discount;default:0.00;comment:'优惠券折扣金额'"`
	FinalAmount   float64   `json:"final_amount" gorm:"column:final_amount;NOT NULL;comment:'最终支付金额(总金额 - 折扣)'"`
	PaymentId     string    `json:"payment_id" gorm:"column:payment_id;NOT NULL;default:'';comment:'支付ID(关联支付信息表)'"`
	PaymentStatus int32     `json:"payment_status" gorm:"column:payment_status;NOT NULL;default:0;comment:'支付状态(0已创建, 1待支付, 2已支付, 3支付超时, 4支付失败, 5取消支付)'"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
}

func (t *Order) TableName() string {
	return "orders"
}

/*
-- @Author AInoriex
-- @Desc 用于记录每个订单中的商品信息(单价, 数量)
-- @Chge 2025年5月5日16点24分 取消外键orders(id)
-- @Chge 2025年5月5日16点24分 取消外键products(id)
-- @Chge 2025年5月5日16点28分 取消order_id, 在orders表用字段item_id关联此表
-- @Chge 2025年5月9日14点44分 取消id唯一键值束缚，同时id改为字符串类型
-- @Chge 2025年5月9日14点45分 新增created_at字段
-- @TODO 增加version字段(如果音乐作品有不同版本), 记录用户购买的商品版本。
CREATE TABLE order_items (
    `id` varchar(32) NOT NULL COMMENT '订单明细ID(同一订单下明细ID相同, 关联订单表)',
    -- `order_id` varchar(16) NOT NULL COMMENT '订单ID(关联订单表)',
    `product_id` varchar(16) NOT NULL COMMENT '商品ID(关联商品表)',
    `quantity` int(8) NOT NULL COMMENT '购买数量',
    `price` decimal(10, 2) NOT NULL COMMENT '商品单价(记录下单时的价格)',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    -- FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    -- FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单明细表';
*/

type OrderItem struct {
	Id        string    `json:"id" gorm:"column:id;NOT NULL;comment:'订单明细ID(同一订单下明细ID相同, 关联订单表)'"`
	ProductId string    `json:"product_id" gorm:"column:product_id;NOT NULL;comment:'商品ID(关联商品表)'"`
	Quantity  int32     `json:"quantity" gorm:"column:quantity;NOT NULL;comment:'购买数量'"`
	Price     float64   `json:"price" gorm:"column:price;NOT NULL;comment:'商品单价(记录下单时的价格)'"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;comment:'创建时间'"`
}

func (t *OrderItem) TableName() string {
	return "order_items"
}

// 创建订单请求参数
type CreateOrderReq struct {
	ItemList []struct { // 商品列表
		ProductId string `json:"product_id"`
		Quantity  int32  `json:"quantity"`
	} `json:"item_list"`
	PaymentMethod      string `json:"payment_method"`       // 支付方式: qrcode, bank, point
	PaymentGatewayType int32  `json:"payment_gateway_type"` // 支付网关: ylt, alipay, wechat
}
