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
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
