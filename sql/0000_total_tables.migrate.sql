-- 插入用户数据
INSERT INTO `users` (`id`, `name`, `email`, `password`, `avatar_url`) VALUES
('user_1', '张三', 'zhangsan@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://example.com/avatar/1.jpg'),
('user_2', '李四', 'lisi@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://example.com/avatar/2.jpg'),
('user_3', '王五', 'wangwu@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://example.com/avatar/3.jpg'),
('user_4', '赵六', 'zhaoliu@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://example.com/avatar/4.jpg'),
('user_5', '孙七', 'sunqi@example.com', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://example.com/avatar/5.jpg');

-- 插入商品数据
INSERT INTO `products` (`id`, `title`, `description`, `price`, `status`, `image_url`, `sales`) VALUES
('P00001', '华为 Mate 50 Pro', '5G 手机，徕卡镜头，12GB 内存，256GB 存储', 5999.00, 1, 'https://example.com/products/1.jpg', 120),
('P00002', '小米 13 Ultra', '5G 手机，1 英寸大底相机，16GB 内存，512GB 存储', 4999.00, 1, 'https://example.com/products/2.jpg', 85),
('P00003', '苹果 iPhone 14 Pro Max', 'A16 芯片，ProMotion 屏幕，256GB 存储', 8999.00, 1, 'https://example.com/products/3.jpg', 200),
('P00004', 'vivo X90 Pro+', '5G 手机，蔡司镜头，12GB 内存，256GB 存储', 4799.00, 0, 'https://example.com/products/4.jpg', 90), -- 注意：此商品已下架
('P00005', 'OPPO Find X6 Pro', '5G 手机，哈苏镜头，16GB 内存，512GB 存储', 5499.00, 1, 'https://example.com/products/5.jpg', 110);

-- 插入购物车数据
INSERT INTO `cart_items` (`user_id`, `product_id`, `quantity`) VALUES
(1, 'P00001', 1),
(2, 'P00002', 3), -- 用户李四购买了3件小米 13 Ultra
(3, 'P00003', 1),
(4, 'P00005', 2), -- 用户赵六购买了2件OPPO Find X6 Pro
(5, 'P00001', 4); -- 用户孙七购买了4件华为 Mate 50 Pro

-- 插入订单数据
INSERT INTO `orders` (`id`, `user_id`, `item_id`, `total_amount`, `discount`, `final_amount`, `payment_id`, `payment_status`) VALUES
('O00001', 1, 1, 5999.00, 0.00, 5999.00, 'PAY00001', 2), -- 已支付
('O00002', 2, 2, 14997.00, 500.00, 14497.00, 'PAY00002', 2), -- 已支付
('O00003', 3, 3, 8999.00, 0.00, 8999.00, 'PAY00003', 1), -- 待支付
('O00004', 4, 4, 10998.00, 1000.00, 9998.00, 'PAY00004', 3), -- 支付超时
('O00005', 5, 5, 23996.00, 2000.00, 21996.00, 'PAY00005', 2); -- 已支付

-- 插入订单明细数据
INSERT INTO `order_items` (`id`, `product_id`, `quantity`, `price`) VALUES
(1, 'P00001', 1, 5999.00),
(2, 'P00002', 3, 4999.00),
(3, 'P00003', 1, 8999.00),
(4, 'P00005', 2, 5499.00),
(5, 'P00001', 4, 5999.00);

-- 插入购买历史数据
INSERT INTO `purchase_history` (`user_id`, `product_id`, `quantity`, `purchase_date`) VALUES
(1, 'P00001', 1, '2025-05-01 10:00:00'),
(2, 'P00002', 3, '2025-05-02 11:00:00'),
(3, 'P00003', 1, '2025-05-03 12:00:00'),
(4, 'P00005', 2, '2025-05-04 13:00:00'),
(5, 'P00001', 4, '2025-05-05 14:00:00');

-- 插入支付数据
INSERT INTO `payments` (`id`, `order_id`, `final_amount`, `payment_method`, `status`, `payment_gateway`, `payment_gateway_id`) VALUES
('PAY00001', 'O00001', 5999.00, '原力推', 2, 10, 'WX00001'),
('PAY00002', 'O00002', 14497.00, '支付宝', 2, 11, 'AL00002'),
('PAY00003', 'O00003', 8999.00, '支付宝', 1, 11, 'ST00003'), -- 待支付
('PAY00004', 'O00004', 9998.00, '微信支付', 3, 12, 'BT00004'), -- 支付超时
('PAY00005', 'O00005', 21996.00, '微信支付', 2, 12, 'WX00005');