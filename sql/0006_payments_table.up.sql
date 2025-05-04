-- @Author AInoriex
-- @Desc 用于记录支付渠道的支付结果。不使用触发器强制同步更新orders.status的状态。
-- @Hint 移除了所有列定义中的`CHECK`约束。如果需要确保`final_amount`的值等于`total_amount`-`discount`，可以在应用程序逻辑中进行验证，或者考虑使用触发器来实现这一逻辑。
CREATE TABLE payments (
    `id` VARCHAR(255) NOT NULL COMMENT '支付唯一标识',
    `order_id` INT NOT NULL COMMENT '订单ID(关联订单表)',
    `final_amount` DECIMAL(10, 2) NOT NULL COMMENT '最终支付金额',
    `payment_method` VARCHAR(255) NOT NULL COMMENT '支付方式(如信用卡、银行转账等)',
    `status` INT NOT NULL DEFAULT 0 COMMENT '支付状态(0已创建, 1待支付, 2已支付, 3支付超时, 4支付失败, 5取消支付)',
    `payment_gateway` VARCHAR(255) COMMENT '支付网关(如Stripe、PayPal等)',
    `payment_gateway_id` VARCHAR(255) COMMENT '支付网关ID(来自支付网关)',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
    -- FOREIGN KEY (status) REFERENCES orders(payment_status),
    -- TRIGGER after_payment_status_update
    --     AFTER UPDATE ON payments
    --     FOR EACH ROW
    --     BEGIN
    --         UPDATE orders SET payment_status = NEW.status 
    --         WHERE id = NEW.order_id;
    --     END;
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付订单关联支付方式表';