-- 切换到eshop数据库
USE eshop;

-- @Author  AInoriex
-- @Des     用户表新增字段：角色权限、最后登录时间、状态、账户锁定时间
-- @Create  2025年6月26日13点58分
ALTER TABLE users
ADD COLUMN `roles` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'user' COMMENT '用户角色权限（admin:管理员, user:普通用户，逗号分隔）' AFTER `updated_at`,
ADD COLUMN `last_login` datetime DEFAULT NULL COMMENT '最后登录时间' AFTER `roles`,
ADD COLUMN `status` tinyint(1) DEFAULT '1' COMMENT '用户状态（1:正常, 0:禁用）' AFTER `last_login`,
ADD COLUMN `banned_at` datetime DEFAULT NULL COMMENT '账户锁定时间' AFTER `status`;
