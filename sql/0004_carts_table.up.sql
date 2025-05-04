-- @Author AInoriex
-- @Desc 用于记录用户购物车信息
-- @TODO 如果音乐作品有不同的版本(如普通版、高清版、无损版), 可以在购物车表中增加version字段, 记录用户选择的商品版本。
CREATE TABLE cart_items (
    `id` INT AUTO_INCREMENT COMMENT '购物车项目唯一标识',
    `user_id` INT NOT NULL COMMENT '用户ID(关联用户表)',
    `product_id` INT NOT NULL COMMENT '商品ID(关联商品表)',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '购买数量',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_user_product (user_id, product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车表';