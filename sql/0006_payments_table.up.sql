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
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付信息表';