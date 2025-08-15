-- 切换到eshop数据库
USE eshop;

-- @Author AInoriex
-- @Desc 用于记录商品基本信息
-- @TODO 补充特定商品的属性信息(音声格式, 音声时长……)
-- @Chge 2025年5月6日11点06分 id int(11) -> varchar(16)
-- @Chge 2025年5月9日17点46分 新增字段external_id, external_link
-- @Chge 2025年5月13日10点08分 id varchar(16) -> varchar(32)
CREATE TABLE `products` (
  `id` varchar(32) NOT NULL COMMENT '商品唯一标识',
  `title` varchar(100) NOT NULL COMMENT '商品标题',
  `description` text COMMENT '商品描述',
  `price` decimal(10,2) NOT NULL COMMENT '商品价格',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '商品状态(0下架, 1上架)',
  `image_url` varchar(255) DEFAULT NULL COMMENT '商品图片URL',
  `sales` int(11) DEFAULT '0' COMMENT '商品销量',
  `external_id` varchar(64) DEFAULT NULL COMMENT '外部商品ID',
  `external_link` varchar(64) DEFAULT NULL COMMENT '外部商品链接',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_title` (`title`),
  KEY `idx_price` (`price`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- @Author AInoriex
-- @Desc 新增字段source_type, 用于记录商品来源类别
-- @Chge 2025年8月11日18点26分 新增字段source_type
ALTER TABLE `eshop`.`products` 
ADD COLUMN `source_type` int(128) NULL DEFAULT 0 COMMENT '来源类别' AFTER `created_at`;


-- @Author AInoriex
-- @Desc 商品播放信息
-- @Chge 2025年8月11日18点26分 创建表products_player
CREATE TABLE `products_player` (
  `id` varchar(32) NOT NULL COMMENT '播放id',
  `product_id` varchar(32) NOT NULL COMMENT '商品id',
  `filename` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_type` varchar(16) NOT NULL DEFAULT '' COMMENT '文件类型',
  `play_type` varchar(16) NOT NULL DEFAULT '' COMMENT '播放类型',
  `file_size` int(11) NOT NULL DEFAULT '0' COMMENT '文件大小（Byte字节）',
  `duration` int(11) NOT NULL DEFAULT '0' COMMENT '文件时长',
  `play_url` varchar(255) NOT NULL DEFAULT '' COMMENT '播放地址',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品播放信息表';

