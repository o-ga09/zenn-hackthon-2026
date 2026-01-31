-- +migrate Up
-- usersテーブル
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    uid VARCHAR(255) NOT NULL UNIQUE COMMENT 'FirebaseのUID',
    name VARCHAR(255) NOT NULL COMMENT '表示名',
    type VARCHAR(50) NOT NULL COMMENT 'ユーザータイプ: admin, tavinikkiy, tavinikkiy-agent',
    plan VARCHAR(50) NOT NULL COMMENT 'サブスクリプションプラン: free, premium',
    token_balance BIGINT NULL COMMENT 'トークン残高',
    is_public BOOLEAN NULL COMMENT 'プロフィールの公開設定',
    display_name VARCHAR(255) NULL COMMENT '表示名',
    bio TEXT NULL COMMENT '自己紹介',
    profile_image VARCHAR(512) NULL COMMENT 'プロフィール画像のURL',
    birth_day VARCHAR(50) NULL COMMENT '誕生日',
    gender VARCHAR(50) NULL COMMENT '性別',
    followers_count BIGINT NULL DEFAULT 0 COMMENT 'フォロワー数',
    following_count BIGINT NULL DEFAULT 0 COMMENT 'フォロー中数',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_uid (uid),
    INDEX idx_type (type),
    INDEX idx_plan (plan),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- vlogsテーブル
CREATE TABLE IF NOT EXISTS vlogs (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    video_id VARCHAR(255) NULL COMMENT '動画ID',
    video_url VARCHAR(512) NULL COMMENT '動画のURL',
    share_url VARCHAR(512) NULL COMMENT '動画の共有URL',
    duration DOUBLE NULL COMMENT '動画の長さ',
    thumbnail VARCHAR(512) NULL COMMENT '動画のサムネイル',
    status VARCHAR(50) NOT NULL DEFAULT 'pending' COMMENT '作成状況',
    error_message TEXT NULL COMMENT '失敗時のエラーメッセージ',
    progress DOUBLE NOT NULL DEFAULT 0 COMMENT '作成の進行状況',
    started_at TIMESTAMP NULL COMMENT '作成開始日時',
    completed_at TIMESTAMP NULL COMMENT '作成完了日時',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_vlogs_user_id FOREIGN KEY (create_user_id) REFERENCES users (id),
    INDEX idx_user_id (create_user_id),
    INDEX idx_video_id (video_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- paymentsテーブル
CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    user_id VARCHAR(255) NOT NULL COMMENT 'ユーザーID',
    type VARCHAR(50) NOT NULL COMMENT 'トークン購入、サブスクリプション',
    amount INT NOT NULL COMMENT '購入金額（円）',
    tokens_granted INT NOT NULL COMMENT '付与トークン数',
    status VARCHAR(50) NOT NULL COMMENT '保留中、完了、失敗',
    stripe_payment_id VARCHAR(255) NOT NULL COMMENT 'Stripeの支払いID',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    CONSTRAINT fk_payments_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_stripe_payment_id (stripe_payment_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- subscriptionsテーブル
CREATE TABLE IF NOT EXISTS subscriptions (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    user_id VARCHAR(255) NOT NULL COMMENT 'ユーザーID',
    plan VARCHAR(50) NOT NULL COMMENT 'プラン',
    status VARCHAR(50) NOT NULL COMMENT '契約状況',
    stripe_customer_id VARCHAR(255) NOT NULL COMMENT 'Stripeの顧客ID',
    stripe_subscription_id VARCHAR(255) NOT NULL COMMENT 'StripeのサブスクリプションID',
    current_period_end TIMESTAMP NOT NULL COMMENT '契約期間終了日時',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_subscriptions_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_stripe_customer_id (stripe_customer_id),
    INDEX idx_stripe_subscription_id (stripe_subscription_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- token_transactionsテーブル
CREATE TABLE IF NOT EXISTS token_transactions (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    user_id VARCHAR(255) NOT NULL COMMENT 'ユーザーID',
    type VARCHAR(50) NOT NULL COMMENT '購入、消費、ボーナス、返金',
    amount INT NOT NULL COMMENT 'トークン数（消費時はマイナス）',
    balance INT NOT NULL COMMENT '取引後の残高',
    description VARCHAR(512) NOT NULL COMMENT '動画生成、月額プラン付与など',
    metadata JSON NULL COMMENT 'メタデータ',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_token_transactions_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- media テーブル
CREATE TABLE IF NOT EXISTS media (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    url VARCHAR(512) NOT NULL COMMENT 'メディアのURL',
    size BIGINT NOT NULL COMMENT 'ファイルサイズ（バイト）',
    content_type VARCHAR(100) NOT NULL COMMENT 'MIMEタイプ',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- media_analyticsテーブル
CREATE TABLE IF NOT EXISTS media_analytics (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    file_id VARCHAR(255) NOT NULL COMMENT 'ファイルID',
    description TEXT NOT NULL COMMENT '全体的な説明',
    mood VARCHAR(100) NOT NULL COMMENT '雰囲気（楽しい、穏やか、など）',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT uc_file_id UNIQUE (file_id),
    CONSTRAINT fk_media_analytics_file_id FOREIGN KEY (file_id) REFERENCES media (id),
    INDEX idx_file_id (file_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- subtitle_segmentsテーブル
CREATE TABLE IF NOT EXISTS subtitle_segments (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    `index` INT NOT NULL,
    start VARCHAR(50) NOT NULL COMMENT '開始時間',
    end VARCHAR(50) NOT NULL COMMENT '終了時間',
    text TEXT NOT NULL COMMENT '表示テキスト',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_index (`index`),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- objectsテーブル
CREATE TABLE IF NOT EXISTS objects (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    media_analytics_id VARCHAR(255) NOT NULL COMMENT 'メディアアナリティクスID',
    name VARCHAR(255) NOT NULL COMMENT 'オブジェクト名',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_objects_media_analytics_id FOREIGN KEY (media_analytics_id) REFERENCES media_analytics (id),
    INDEX idx_media_analytics_id (media_analytics_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- landmarksテーブル
CREATE TABLE IF NOT EXISTS landmarks (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    media_analytics_id VARCHAR(255) NOT NULL COMMENT 'メディアアナリティクスID',
    name VARCHAR(255) NOT NULL COMMENT 'ランドマーク名',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_landmarks_media_analytics_id FOREIGN KEY (media_analytics_id) REFERENCES media_analytics (id),
    INDEX idx_media_analytics_id (media_analytics_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- activitiesテーブル
CREATE TABLE IF NOT EXISTS activities (
    id VARCHAR(255) PRIMARY KEY,
    version INT NOT NULL DEFAULT 0 COMMENT '楽観的ロックのバージョン',
    create_user_id VARCHAR(255) NULL COMMENT '作成者のユーザーID',
    update_user_id VARCHAR(255) NULL COMMENT '更新者のユーザーID',
    media_analytics_id VARCHAR(255) NOT NULL COMMENT 'メディアアナリティクスID',
    name VARCHAR(255) NOT NULL COMMENT 'アクティビティ名',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_activities_media_analytics_id FOREIGN KEY (media_analytics_id) REFERENCES media_analytics (id),
    INDEX idx_media_analytics_id (media_analytics_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;