-- @Author AInoriex
-- @Desc 创建订单与商品关联表, 记录用户购买商品数量和订单金额, 外键关联用户ID和商品ID
-- @Desc 新增order_items表以支持一个订单与多个商品关联
-- @Desc 新增优惠券和订单状态设计
-- @Hint MySQL 8.0.16+ 版本才完全支持CHECK约束, 如果使用旧版本可能需要改用触发器实现。
-- @Hint 移除了所有列定义中的`CHECK`约束。如果需要确保`final_amount`的值等于`total_amount`-`discount`，可以在应用程序逻辑中进行验证，或者考虑使用触发器来实现这一逻辑。
-- @TODO 增加source字段, 记录订单来源(如网站、移动端、API等), 方便分析不同渠道的销售情况。
CREATE TABLE orders (
    `id` INT AUTO_INCREMENT COMMENT '订单唯一标识',
    `user_id` INT NOT NULL COMMENT '用户ID(关联用户表)',
    `total_amount` DECIMAL(10, 2) NOT NULL COMMENT '订单总金额',
    `discount` DECIMAL(10, 2) DEFAULT 0.00 COMMENT '优惠券折扣金额',
    `final_amount` DECIMAL(10, 2) NOT NULL COMMENT '最终支付金额(总金额 - 折扣)',
    `payment_status` INT NOT NULL DEFAULT '0' COMMENT '支付状态(0已创建, 1待支付, 2已支付, 3支付超时, 4支付失败, 5取消支付)',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_payment_status (payment_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';


-- @Author AInoriex
-- @Desc 用于记录每个订单中的商品信息(单价, 数量)
-- @TODO 增加version字段(如果音乐作品有不同版本), 记录用户购买的商品版本。
CREATE TABLE order_items (
    `id` INT AUTO_INCREMENT COMMENT '订单明细唯一标识',
    `order_id` INT NOT NULL COMMENT '订单ID(关联订单表)',
    `product_id` INT NOT NULL COMMENT '商品ID(关联商品表)',
    `quantity` INT NOT NULL COMMENT '购买数量',
    `price` DECIMAL(10, 2) NOT NULL COMMENT '商品单价(记录下单时的价格)',
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_order_id (order_id),
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单明细表';


-- @Author AInoriex
-- @Desc 用于记录用户购买历史。记录商品级别的信息, 便于分析用户对每个商品的购买行为分析和个性化推荐, 增强用户体验。
CREATE TABLE purchase_history (
    `id` INT AUTO_INCREMENT COMMENT '购买历史唯一标识',
    `user_id` INT NOT NULL COMMENT '用户ID(关联用户表)',
    `product_id` INT NOT NULL COMMENT '商品ID(关联商品表)',
    `purchase_date` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '购买日期',
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_purchase_date (purchase_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购买历史订单表';


-- @Author AInoriex
-- @Desc 创建订单与商品关联表，记录用户购买商品数量和订单金额，外键关联用户ID和商品ID
-- CREATE TABLE _orders (
--     id INT AUTO_INCREMENT COMMENT '订单唯一标识',
--     user_id INT NOT NULL COMMENT '用户ID（关联用户表）',
--     product_id INT NOT NULL COMMENT '商品ID（关联商品表）',
--     quantity INT NOT NULL COMMENT '购买数量',
--     total DECIMAL(10, 2) NOT NULL COMMENT '订单总金额',
--     purchase_date DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '购买日期',
--     PRIMARY KEY (id),
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
--     FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
--     INDEX idx_user_id (user_id),
--     INDEX idx_purchase_date (purchase_date)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- @Author AInoriex
-- @Desc 创建订单与商品关联表，记录用户购买商品数量和订单金额，外键关联用户ID和商品ID
-- @Desc 新增order_items表以支持一个订单与多个商品关联
-- CREATE TABLE __orders (
--     id INT AUTO_INCREMENT COMMENT '订单唯一标识',
--     user_id INT NOT NULL COMMENT '用户ID（关联用户表）',
--     total DECIMAL(10, 2) NOT NULL COMMENT '订单总金额',
--     purchase_date DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '购买日期',
--     PRIMARY KEY (id),
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE COMMENT '用户外键',
--     INDEX idx_user_id (user_id),
--     INDEX idx_purchase_date (purchase_date)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';
