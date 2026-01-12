-- +migrate Up
-- usersテーブル
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT 'Optimistic locking version',
    create_user_id VARCHAR(255) NULL COMMENT 'Created by user ID',
    update_user_id VARCHAR(255) NULL COMMENT 'Updated by user ID',
    uid VARCHAR(255) NOT NULL UNIQUE COMMENT 'Firebase UID',
    name VARCHAR(255) NOT NULL COMMENT 'Display name',
    type VARCHAR(50) NOT NULL COMMENT 'User type: admin, tavinikkiy, tavinikkiy-agent',
    plan VARCHAR(50) NOT NULL COMMENT 'Subscription plan: free, premium',
    token_balance BIGINT NULL COMMENT 'Token balance',
    is_public BOOLEAN NULL COMMENT 'Is profile public',
    display_name VARCHAR(255) NULL COMMENT 'Display name',
    bio TEXT NULL COMMENT 'Biography',
    profile_image VARCHAR(512) NULL COMMENT 'Profile image URL',
    birth_day VARCHAR(50) NULL COMMENT 'Birthday',
    gender VARCHAR(50) NULL COMMENT 'Gender',
    followers_count BIGINT NULL DEFAULT 0 COMMENT 'Number of followers',
    following_count BIGINT NULL DEFAULT 0 COMMENT 'Number of following',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_uid (uid),
    INDEX idx_type (type),
    INDEX idx_plan (plan),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- vlogsテーブル
CREATE TABLE IF NOT EXISTS vlogs (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT 'Optimistic locking version',
    create_user_id VARCHAR(255) NULL COMMENT 'Created by user ID',
    update_user_id VARCHAR(255) NULL COMMENT 'Updated by user ID',
    video_id VARCHAR(255) NOT NULL,
    video_url VARCHAR(512) NOT NULL,
    share_url VARCHAR(512) NOT NULL,
    duration DOUBLE NOT NULL,
    thumbnail VARCHAR(512) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_video_id (video_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- paymentsテーブル
CREATE TABLE IF NOT EXISTS payments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT 'Optimistic locking version',
    create_user_id VARCHAR(255) NULL COMMENT 'Created by user ID',
    update_user_id VARCHAR(255) NULL COMMENT 'Updated by user ID',
    uid VARCHAR(255) NOT NULL COMMENT 'Firebase UID',
    type VARCHAR(50) NOT NULL COMMENT 'token_purchase, subscription',
    amount INT NOT NULL COMMENT '金額（円）',
    tokens_granted INT NOT NULL COMMENT '付与トークン数',
    status VARCHAR(50) NOT NULL COMMENT 'pending, completed, failed',
    stripe_payment_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    INDEX idx_uid (uid),
    INDEX idx_status (status),
    INDEX idx_stripe_payment_id (stripe_payment_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- subscriptionsテーブル
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT 'Optimistic locking version',
    create_user_id VARCHAR(255) NULL COMMENT 'Created by user ID',
    update_user_id VARCHAR(255) NULL COMMENT 'Updated by user ID',
    uid VARCHAR(255) NOT NULL COMMENT 'Firebase UID',
    plan VARCHAR(50) NOT NULL COMMENT 'monthly, yearly',
    status VARCHAR(50) NOT NULL COMMENT 'active, cancelled, expired',
    stripe_customer_id VARCHAR(255) NOT NULL,
    stripe_subscription_id VARCHAR(255) NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_uid (uid),
    INDEX idx_status (status),
    INDEX idx_stripe_customer_id (stripe_customer_id),
    INDEX idx_stripe_subscription_id (stripe_subscription_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- token_transactionsテーブル
CREATE TABLE IF NOT EXISTS token_transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT 'Optimistic locking version',
    create_user_id VARCHAR(255) NULL COMMENT 'Created by user ID',
    update_user_id VARCHAR(255) NULL COMMENT 'Updated by user ID',
    uid VARCHAR(255) NOT NULL COMMENT 'Firebase UID',
    type VARCHAR(50) NOT NULL COMMENT 'purchase, consumption, bonus, refund',
    amount INT NOT NULL COMMENT 'トークン数（消費時はマイナス）',
    balance INT NOT NULL COMMENT '取引後の残高',
    description VARCHAR(512) NOT NULL COMMENT '動画生成、月額プラン付与など',
    metadata JSON NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_uid (uid),
    INDEX idx_type (type),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- media_analyticsテーブル
CREATE TABLE IF NOT EXISTS media_analytics (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT 'Optimistic locking version',
    create_user_id VARCHAR(255) NULL COMMENT 'Created by user ID',
    update_user_id VARCHAR(255) NULL COMMENT 'Updated by user ID',
    file_id VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL COMMENT 'image or video',
    description TEXT NOT NULL COMMENT '全体的な説明',
    objects JSON NULL COMMENT '検出されたオブジェクト',
    landmarks JSON NULL COMMENT '観光地・ランドマーク',
    activities JSON NULL COMMENT 'アクティビティ',
    mood VARCHAR(100) NOT NULL COMMENT '雰囲気（楽しい、穏やか、など）',
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_file_id (file_id),
    INDEX idx_type (type),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- subtitle_segmentsテーブル
CREATE TABLE IF NOT EXISTS subtitle_segments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT 'Optimistic locking version',
    create_user_id VARCHAR(255) NULL COMMENT 'Created by user ID',
    update_user_id VARCHAR(255) NULL COMMENT 'Updated by user ID',
    `index` INT NOT NULL,
    start VARCHAR(50) NOT NULL COMMENT '00:00:01,000',
    end VARCHAR(50) NOT NULL COMMENT '00:00:04,000',
    text TEXT NOT NULL COMMENT '表示テキスト',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_index (`index`),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- media テーブル
CREATE TABLE IF NOT EXISTS media (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT 'Optimistic locking version',
    create_user_id VARCHAR(255) NULL COMMENT 'Created by user ID',
    update_user_id VARCHAR(255) NULL COMMENT 'Updated by user ID',
    url VARCHAR(512) NOT NULL COMMENT 'Media URL',
    type VARCHAR(50) NOT NULL COMMENT 'image, video, audio',
    size BIGINT NOT NULL COMMENT 'File size in bytes',
    content_type VARCHAR(100) NOT NULL COMMENT 'MIME type',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_type (type),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
