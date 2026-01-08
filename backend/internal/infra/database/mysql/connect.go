package mysql

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/o-ga09/zenn-hackthon-2026/internal/infra/database"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
)

const (
	maxRetries      = 5
	retryInterval   = 2 * time.Second
	maxIdleConns    = 10
	maxOpenConns    = 100
	connMaxLifetime = 1 * time.Hour
)

func Connect(ctx context.Context) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	env := config.GetCtxEnv(ctx)
	dsn := env.Database_url
	logger := database.NewSentryLogger()
	// リトライ処理
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: false,
			},
			Logger: logger,
		})
		if err != nil {
			if i == maxRetries-1 {
				return nil, fmt.Errorf("failed to open database after %d retries: %w", maxRetries, err)
			}
			time.Sleep(retryInterval)
			continue
		}

		// 接続成功
		break
	}

	// SQLDBインスタンスを取得
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// コネクションプール設定
	sqlDB.SetMaxIdleConns(10)           // アイドル状態の最大接続数
	sqlDB.SetMaxOpenConns(100)          // 最大接続数
	sqlDB.SetConnMaxLifetime(time.Hour) // 接続の最大生存期間

	// UUID自動付与プラグインを登録
	if err := db.Use(database.NewUUIDPlugin()); err != nil {
		return nil, fmt.Errorf("failed to register UUID plugin: %w", err)
	}

	return db, nil
}
