-- シードデータ: users テーブル
INSERT INTO users (id, uid, name, type, plan, token_balance, created_at, updated_at)
VALUES
    ('1', 'user1', 'User One', 'admin', 'premium', 10000, NOW(), NOW()),
    ('2', 'user2', 'User Two', 'tavinikkiy', 'free', 0, NOW(), NOW()),
    ('3', 'user3', 'User Three', 'tavinikkiy-agent', 'premium', 5000, NOW(), NOW());