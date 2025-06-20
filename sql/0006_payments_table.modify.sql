-- @Author  AInoriex
-- @Des     purchase_history表新增订单ID字段
-- @Create  2025年6月20日14点03分
ALTER TABLE `eshop`.`purchase_history` 
ADD COLUMN `order_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '订单ID(关联订单表)' AFTER `purchased_at`;
