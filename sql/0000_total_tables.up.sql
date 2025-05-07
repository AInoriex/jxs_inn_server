-- @Author AInoriex
-- @Desc 目前只支持邮箱一种方式登陆
-- @Chge 2025年5月6日11点06分 id int(11) -> varchar(32)
-- @TODO 用户角色：如果未来有管理员、普通用户等不同角色, 可以增加一个role字段, 用于区分用户权限。
-- @TODO 联系方式：除了邮箱, 可以增加手机号字段, 方便用户接收验证码、订单通知等信息。
-- @TODO 账户锁定机制：可以增加login_attempts字段记录登录失败次数, 当连续多次登录失败时, 暂时锁定账户, 防止暴力破解。
-- @TTODO 会员信息：如果计划推出会员制度, 可以增加会员等级、会员积分等字段。
-- @TTODO 登录方式：除了邮箱登录, 可以考虑支持社交媒体账号登录(如微信、QQ、微博等), 增加social_login_id字段存储第三方登录的唯一标识。
CREATE TABLE `users` (
  `id` varchar(32) NOT NULL COMMENT '用户唯一标识',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户姓名',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户邮箱',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户密码(强加密算法存储, 如bcrypt、scrypt等)',
  `avatar_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户头像URL',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';


-- @Author AInoriex
-- @Desc 用于记录商品基本信息
-- @TODO 补充特定商品的属性信息(音声格式, 音声时长……)
-- @Chge id int(11) -> varchar(16)
CREATE TABLE `products` (
  `id` varchar(16) NOT NULL COMMENT '商品唯一标识',
  `title` varchar(100) NOT NULL COMMENT '商品标题',
  `description` text COMMENT '商品描述',
  `price` decimal(10,2) NOT NULL COMMENT '商品价格',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '商品状态(0下架, 1上架)',
  `image_url` varchar(255) DEFAULT NULL COMMENT '商品图片URL',
  `sales` int(11) DEFAULT '0' COMMENT '商品销量',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_title` (`title`),
  KEY `idx_price` (`price`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';


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
    `id` varchar(16) NOT NULL COMMENT '订单唯一标识',
    `user_id` int(11) NOT NULL COMMENT '用户ID(关联用户表)',
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

-- @Author AInoriex
-- @Desc 用于记录每个订单中的商品信息(单价, 数量)
-- @Chge 2025年5月5日16点24分 取消外键orders(id)
-- @Chge 2025年5月5日16点24分 取消外键products(id)
-- @Chge 2025年5月5日16点28分 取消order_id, 在orders表用字段item_id关联此表
-- @TODO 增加version字段(如果音乐作品有不同版本), 记录用户购买的商品版本。
CREATE TABLE order_items (
    `id` int(11) AUTO_INCREMENT COMMENT '订单明细唯一标识',
    -- `order_id` varchar(16) NOT NULL COMMENT '订单ID(关联订单表)',
    `product_id` varchar(16) NOT NULL COMMENT '商品ID(关联商品表)',
    `quantity` int(8) NOT NULL COMMENT '购买数量',
    `price` decimal(10, 2) NOT NULL COMMENT '商品单价(记录下单时的价格)',
    PRIMARY KEY (id),
    -- FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    -- FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单明细表';

-- @Author AInoriex
-- @Desc 用于记录用户购买历史。记录商品级别的信息, 便于分析用户对每个商品的购买行为分析和个性化推荐, 增强用户体验。
CREATE TABLE purchase_history (
    `id` int(11) AUTO_INCREMENT COMMENT '购买历史唯一标识',
    `user_id` varchar(32) NOT NULL COMMENT '用户ID(关联用户表)',
    `product_id` varchar(16) NOT NULL COMMENT '商品ID(关联商品表)',
    `quantity` int(8) NOT NULL COMMENT '购买数量',
    `purchase_date` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '购买日期',
    PRIMARY KEY (id),
    INDEX idx_user_id (user_id),
    INDEX idx_purchase_date (purchase_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购买历史订单表';


-- @Author AInoriex
-- @Desc 用于记录支付渠道的支付结果。不使用触发器强制同步更新orders.status的状态。
-- @Hint 移除了所有列定义中的`CHECK`约束。如果需要确保`final_amount`的值等于`total_amount`-`discount`，可以在应用程序逻辑中进行验证，或者考虑使用触发器来实现这一逻辑。
-- @Chge 2025年5月5日16点24分 取消外键orders(id)
CREATE TABLE payments (
    `id` varchar(255) NOT NULL COMMENT '支付唯一标识',
    `order_id` varchar(16) NOT NULL COMMENT '订单ID(关联订单表)',
    `final_amount` decimal(10, 2) NOT NULL COMMENT '最终支付金额',
    `payment_method` varchar(255) NOT NULL COMMENT '支付方式(如信用卡、银行转账等)',
    `status` tinyint(3) NOT NULL DEFAULT 0 COMMENT '支付状态(0已创建, 1待支付, 2已支付, 3支付超时, 4支付失败, 5取消支付)',
    `payment_gateway` tinyint(3) COMMENT '支付网关(10ylt, 11zfb, 12wx)',
    `payment_gateway_id` varchar(255) COMMENT '支付网关ID(来自支付网关)',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id)
    -- FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
    -- FOREIGN KEY (status) REFERENCES orders(payment_status),
    -- TRIGGER after_payment_status_update
    --     AFTER UPDATE ON payments
    --     FOR EACH ROW
    --     BEGIN
    --         UPDATE orders SET payment_status = NEW.status 
    --         WHERE id = NEW.order_id;
    --     END;
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付信息表';