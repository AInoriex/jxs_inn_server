-- 切换到eshop数据库
USE eshop;

-- 插入用户数据
INSERT INTO `users` (`id`, `name`, `email`, `password`, `avatar_url`) VALUES
('user_1', '小O', 'xiaoo@163.com', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92', 'https://example.com/avatar/2.jpg'),
('user_2', '小A', 'xiaoa@163.com', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92', 'https://example.com/avatar/3.jpg'),
('user_3', '小B', 'xiaob@163.com', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92', 'https://example.com/avatar/4.jpg');

-- 插入商品数据
INSERT INTO `products` (`id`, `title`, `description`, `price`, `status`, `image_url`, `sales`, `external_id`, `external_link`) VALUES
('Ad00000', '露天电影', '由莫妮卡主演的爱情题材电影。荒无人烟的郊外突然放起露天电影,这件事逐渐作为怪谈流传开来。', 25.00, 1, 'https://act-webstatic.mihoyo.com/puzzle/zzz/pz_zCT32XGuUS/resource/puzzle/2025/07/02/ed00325e8e7d9bcdd775419af2e79d12_4960607706852427232.png?x-oss-process=image/format,webp/quality,Q_90', 200, '001', 'https://zzz.mihoyo.com/?utm_source=ooppgw'),
('Ad00001', '大身材小危机', '鲍勃一觉醒来发现自己变成了躺在废弃站点里的破旧OO。为了找回自己的身体,鲍勃踏上了啼笑皆非的冒险之旅。', 19.99, 1, 'https://act-webstatic.mihoyo.com/puzzle/zzz/pz_zCT32XGuUS/resource/puzzle/2025/07/02/e3dfe89db7bfeeac49e63a5eae796fd8_6023131665309002104.png?x-oss-process=image/format,webp/quality,Q_90', 120, '002', 'https://zzz.mihoyo.com/?utm_source=ooppgw'),
('Ad00002', '流浪的足迹', '根据真实故事改编,在空洞中流浪的小狗,在好心人的帮助下度过危机,逐渐成长为与空洞调查员并肩战斗的明星救助犬。', 14.50, 1, 'https://act-webstatic.mihoyo.com/puzzle/zzz/pz_zCT32XGuUS/resource/puzzle/2025/07/02/6ea4ccc02a22af5e8cf2951bd3635758_4763043475781278964.png?x-oss-process=image/format,webp/quality,Q_90', 85, '003', 'https://zzz.mihoyo.com/?utm_source=ooppgw');

-- 插入购物车数据
INSERT INTO `cart_items` (`user_id`, `product_id`, `quantity`) VALUES
('user_1', 'Ad00000', 1),
('user_2', 'Ad00001', 3),
('user_3', 'Ad00002', 1),

-- 插入订单数据
INSERT INTO `orders` (`id`, `user_id`, `item_id`, `total_amount`, `discount`, `final_amount`, `payment_id`, `payment_status`) VALUES
('Order00001', 'user_1', 'OrderItem00001', 0.50, 0.00, 0.50, 'PAY00001', 2), -- 已支付
('Order00002', 'user_2', 'OrderItem00002', 114.00, 0.00, 114.00, 'PAY00002', 2), -- 已支付
('Order00003', 'user_3', 'OrderItem00003', 38.00, 0.00, 38.00, 'PAY00003', 1), -- 待支付

-- 插入订单明细数据
INSERT INTO `order_items` (`id`, `product_id`, `quantity`, `price`) VALUES
('OrderItem00001', 'Ad00000', 1, 0.50),
('OrderItem00002', 'Ad00001', 3, 114.00),
('OrderItem00003', 'Ad00002', 1, 38.00),
('OrderItem00004', 'Ad00000', 4, 2.00);

-- 插入支付数据
INSERT INTO `payments` (`id`, `order_id`, `final_amount`, `method`, `status`, `gateway_type`, `gateway_id`, `agent`, `purchased_at`) VALUES
('PAY00001', 'Order00001', 0.50, 'wx', 2, 10, 'WX00001', '', '2025-05-01 10:00:00'), -- 已支付
('PAY00002', 'Order00002', 114.00, 'zfb', 2, 11, 'AL00002', '', '2025-05-02 11:00:00'), -- 已支付
('PAY00003', 'Order00003', 38.00, 'zfb', 1, 11, 'ST00003', '', NULL); -- 待支付

-- 插入购买历史数据
INSERT INTO `purchase_history` (`user_id`, `product_id`, `quantity`, `payment_id`, `purchased_at`) VALUES
('user_1', 'Ad00000', 1, 'PAY00001', '2025-05-01 10:00:00'),
('user_2', 'Ad00001', 3, 'PAY00002', '2025-05-02 11:00:00');

