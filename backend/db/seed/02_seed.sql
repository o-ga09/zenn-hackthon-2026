-- シードデータ: users テーブル
INSERT INTO users (id, uid, name, type, plan, token_balance, created_at, updated_at)
VALUES
    ('1', 'user1', 'User One', 'admin', 'premium', 10000, NOW(), NOW()),
    ('2', 'user2', 'User Two', 'tavinikkiy', 'free', 0, NOW(), NOW()),
    ('3', 'user3', 'User Three', 'tavinikkiy-agent', 'premium', 5000, NOW(), NOW());

-- vlogテーブ
INSERT INTO vlogs (id,video_id,video_url, thumbnail,duration,share_url, created_at, updated_at)
VALUES
    ('1', 'vid1', 'https://pub-0a072bc79aa54f28b971e7bd751566a4.r2.dev/%E7%94%BB%E9%9D%A2%E5%8F%8E%E9%8C%B2%202026-01-12%2015.06.54.mov', 'https://example.com/thumbnails/vlog1.jpg', 30, 'https://example.com/share/vlog1', NOW(), NOW()),
    ('2', 'vid2', 'https://www.youtube.com/shorts/qhiEsT9reHw', 'https://example.com/thumbnails/vlog2.jpg', 45, 'https://example.com/share/vlog2', NOW(), NOW()),
    ('3', 'vid3', 'https://www.youtube.com/watch?v=TMstofGVMBM', 'https://example.com/thumbnails/vlog3.jpg', 120, 'https://example.com/share/vlog3', NOW(), NOW());