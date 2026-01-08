package database

import (
	"testing"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUUIDPlugin(t *testing.T) {
	// テスト用のインメモリデータベースを作成
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// UUIDプラグインを登録
	err = db.Use(NewUUIDPlugin())
	assert.NoError(t, err)

	// テスト用のテーブルを作成
	err = db.AutoMigrate(&domain.User{})
	assert.NoError(t, err)

	t.Run("IDが空の場合、自動的にUUIDが付与される", func(t *testing.T) {
		user := &domain.User{
			BaseModel: domain.BaseModel{
				ID: "", // IDを空にする
			},
			UID:  "test-firebase-uid",
			Name: "Test User",
			Type: "tavinikkiy",
			Plan: "free",
		}

		// ユーザーを作成
		err := db.Create(user).Error
		assert.NoError(t, err)

		// IDが自動的に付与されていることを確認
		assert.NotEmpty(t, user.ID)
		assert.NotEqual(t, "", user.ID)

		// UUIDのフォーマットを確認（36文字）
		assert.Equal(t, 36, len(user.ID))
	})

	t.Run("IDが既に設定されている場合、上書きされない", func(t *testing.T) {
		existingID := "existing-uuid-123"
		user := &domain.User{
			BaseModel: domain.BaseModel{
				ID: existingID,
			},
			UID:  "test-firebase-uid-2",
			Name: "Test User 2",
			Type: "tavinikkiy",
			Plan: "free",
		}

		// ユーザーを作成
		err := db.Create(user).Error
		assert.NoError(t, err)

		// IDが上書きされていないことを確認
		assert.Equal(t, existingID, user.ID)
	})

	t.Run("複数レコード作成時、それぞれ異なるUUIDが付与される", func(t *testing.T) {
		users := []domain.User{
			{
				BaseModel: domain.BaseModel{ID: ""},
				UID:       "user1-firebase-uid",
				Name:      "User 1",
				Type:      "tavinikkiy",
				Plan:      "free",
			},
			{
				BaseModel: domain.BaseModel{ID: ""},
				UID:       "user2-firebase-uid",
				Name:      "User 2",
				Type:      "tavinikkiy",
				Plan:      "premium",
			},
			{
				BaseModel: domain.BaseModel{ID: ""},
				UID:       "user3-firebase-uid",
				Name:      "User 3",
				Type:      "tavinikkiy-agent",
				Plan:      "free",
			},
		}

		// 複数ユーザーを作成
		err := db.Create(&users).Error
		assert.NoError(t, err)

		// 全てのユーザーにUUIDが付与されていることを確認
		for i, user := range users {
			assert.NotEmpty(t, user.ID, "User %d should have an ID", i+1)
			assert.Equal(t, 36, len(user.ID), "User %d ID should be 36 characters", i+1)
		}

		// 全てのUUIDが異なることを確認
		idMap := make(map[string]bool)
		for _, user := range users {
			assert.False(t, idMap[user.ID], "ID should be unique: %s", user.ID)
			idMap[user.ID] = true
		}
	})
}

func TestUUIDPluginName(t *testing.T) {
	plugin := NewUUIDPlugin()
	assert.Equal(t, "UUIDPlugin", plugin.Name())
}
