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
-- @Chge 2025年5月9日14点48分 item_id int(11) -> varchar(32); 新增updated_at字段
-- @Chge 2025年5月13日10点09分 id varchar(16) -> varchar(32)
-- @TODO 增加source字段, 记录订单来源(如网站、移动端、API等), 方便分析不同渠道的销售情况。
CREATE TABLE orders (
    `id` varchar(32) NOT NULL COMMENT '订单唯一标识',
    `user_id` varchar(32) NOT NULL COMMENT '用户ID(关联用户表)',
    `item_id` varchar(32) NOT NULL COMMENT '订单明细ID(关联订单明细表)',
    `total_amount` decimal(10, 2) NOT NULL COMMENT '订单总金额',
    `discount` decimal(10, 2) DEFAULT 0.00 COMMENT '优惠券折扣金额',
    `final_amount` decimal(10, 2) NOT NULL COMMENT '最终支付金额(总金额 - 折扣)',
    `payment_id` varchar(255) NOT NULL DEFAULT '' COMMENT '支付ID(关联支付信息表)',
    `payment_status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '支付状态(0已创建, 1待支付, 2已支付, 3支付超时, 4支付失败, 5取消支付)',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    -- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_payment_status (payment_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

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
    `product_id` varchar(32) NOT NULL COMMENT '商品ID(关联商品表)',
    `quantity` int(8) NOT NULL COMMENT '购买数量',
    `price` decimal(10, 2) NOT NULL COMMENT '商品单价(记录下单时的价格)',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    -- FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    -- FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单明细表';
