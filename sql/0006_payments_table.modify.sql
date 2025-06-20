-- 切换到eshop数据库
USE eshop;

-- @Author  AInoriex
-- @Des     purchase_history表新增订单ID字段
-- @Create  2025年6月20日14点03分
ALTER TABLE `eshop`.`purchase_history` 
ADD COLUMN `order_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '订单ID(关联订单表)' AFTER `purchased_at`;

-- @Author  AInoriex
-- @Des     添加订单ID后修复历史数据
-- @Create  2025年6月20日14点03分
UPDATE `purchase_history` AS ph
JOIN `payments` AS py ON ph.`payment_id` = py.`id`
SET ph.`order_id` = py.`order_id`;
