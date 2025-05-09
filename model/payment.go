package model

import (
	"time"
)

const (
	PaymentMethodQrcode string = "qrcode" // 支付方法:扫码支付
	PaymentMethodBank   string = "bank"   // 支付方法:银行卡支付
	PaymentMethodPoint  string = "point"  // 支付方法:积分兑换

	PaymentStatusCreate    int32 = 0 // 0 支付状态:已创建
	PaymentStatusToPay     int32 = 1 // 1 支付状态:待支付
	PaymentStatusPayed     int32 = 2 // 2 支付状态:已支付
	PaymentStatusTimeOut   int32 = 3 // 3 支付状态:支付超时
	PaymentStatusPayFail   int32 = 4 // 4 支付状态:支付失败
	PaymentStatusPayCancel int32 = 5 // 5 支付状态:取消支付

	PaymentGatewayTypeYlt    int32 = 10 // 10 支付类别:原力通
	PaymentGatewayTypeAlipay int32 = 11 // 11 支付类别:支付宝
	PaymentGatewayTypeWechat int32 = 12 // 12 支付类别:微信
)

/*
-- @Author AInoriex
-- @Desc 用于记录支付渠道的支付结果。不使用触发器强制同步更新orders.status的状态。
-- @Hint 移除了所有列定义中的`CHECK`约束。如果需要确保`final_amount`的值等于`total_amount`-`discount`，可以在应用程序逻辑中进行验证，或者考虑使用触发器来实现这一逻辑。
-- @Chge 2025年5月5日16点24分 取消外键orders(id)
-- @Chge 2025年5月9日17点34分 新增字段agent
-- @Chge 2025年5月9日17点49分 调整字段名gateway->gateway_type
CREATE TABLE payments (
    `id` varchar(255) NOT NULL COMMENT '支付唯一标识',
    `order_id` varchar(16) NOT NULL COMMENT '订单ID(关联订单表)',
    `final_amount` decimal(10, 2) NOT NULL COMMENT '最终支付金额',
    `method` varchar(255) NOT NULL COMMENT '支付方式(如扫码，积分，银行转账等)',
    `status` tinyint(3) NOT NULL DEFAULT 0 COMMENT '支付状态(0已创建, 1待支付, 2已支付, 3支付超时, 4支付失败, 5取消支付)',
    `gateway_type` tinyint(3) NOT NULL DEFAULT 0 COMMENT '支付网关(10ylt, 11zfb, 12wx)',
    `gateway_id` varchar(255) NOT NULL DEFAULT '' COMMENT '支付网关订单ID',
    `agent` varchar(16) NULL COMMENT '支付代理人',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付信息表';
*/

type Payment struct {
	ID          string    `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'支付唯一标识'"`
	OrderID     string    `json:"order_id" gorm:"column:order_id;NOT NULL;comment:'订单ID(关联订单表)'"`
	FinalAmount float64   `json:"final_amount" gorm:"column:final_amount;NOT NULL;comment:'最终支付金额'"`
	Method      string    `json:"method" gorm:"column:method;NOT NULL;comment:'支付方式(如信用卡、银行转账等)'"`
	Status      int32     `json:"status" gorm:"column:status;NOT NULL;default:0;comment:'支付状态(0已创建, 1待支付, 2已支付, 3支付超时, 4支付失败, 5取消支付)'"`
	GatewayType int32     `json:"gateway_type" gorm:"column:gateway_type;NOT NULL;default:0;comment:'支付网关(10ylt, 11zfb, 12wx)'"`
	GatewayID   string    `json:"gateway_id" gorm:"column:gateway_id;NOT NULL;default:'';comment:'支付网关ID(来自支付网关)'"`
	Agent       string    `json:"agent" gorm:"column:agent;comment:'支付代理人'"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;default:NULL ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'"`
}

func (t *Payment) TableName() string {
	return "payments"
}
