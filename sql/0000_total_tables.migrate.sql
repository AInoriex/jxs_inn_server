-- 用户表测试数据
INSERT INTO users (name, email, password, avatar_url) VALUES 
('江明月', 'jmy@qq.com', '123456', 'https://example.com/avatars/zhangsan.jpg'),
('匿名用户', 'test_user@example.com', '123456', 'https://example.com/avatars/lisi.jpg'),
('王五', 'wangwu@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'https://example.com/avatars/wangwu.jpg');

-- 商品表测试数据
INSERT INTO products (title, description, price, status, image_url, sales) VALUES
('音声作品A', '这是一个非常棒的音声作品A', 29.99, 1, 'https://example.com/products/a.jpg', 100),
('音声作品B', '这是一个非常棒的音声作品B', 39.99, 1, 'https://example.com/products/b.jpg', 50),
('音声作品C', '这是一个非常棒的音声作品C', 19.99, 0, 'https://example.com/products/c.jpg', 200);

-- 购物车表测试数据
INSERT INTO cart_items (user_id, product_id, quantity) VALUES
(1, 1, 2),
(1, 2, 1),
(2, 1, 1);

-- 订单表测试数据
INSERT INTO orders (user_id, total_amount, discount, final_amount, payment_status) VALUES
(1, 99.97, 10.00, 89.97, 2),
(2, 29.99, 0.00, 29.99, 1),
(3, 59.98, 5.00, 54.98, 0);

-- 订单明细表测试数据
INSERT INTO order_items (order_id, product_id, quantity, price) VALUES
(1, 1, 2, 29.99),
(1, 2, 1, 39.99),
(2, 1, 1, 29.99);

-- 购买历史表测试数据
INSERT INTO purchase_history (user_id, product_id) VALUES
(1, 1),
(1, 2),
(2, 1);

-- 支付表测试数据
INSERT INTO payments (id, order_id, final_amount, payment_method, status, payment_gateway, payment_gateway_id) VALUES
('PAY20230601001', 1, 89.97, '支付宝', 2, 'Alipay', 'ALI123456789'),
('PAY20230601002', 2, 29.99, '微信支付', 1, 'WeChatPay', 'WX123456789'),
('PAY20230601003', 3, 54.98, '银行卡', 0, 'UnionPay', 'UP123456789');
