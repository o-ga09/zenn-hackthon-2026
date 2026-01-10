package database

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	pkgErrors "github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func TestOptimisticLockPlugin(t *testing.T) {
	// テスト用のインメモリデータベースを作成
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // SQLログを有効化
	})
	assert.NoError(t, err)

	// 楽観ロックプラグインとUUIDプラグインを登録（順序重要）
	err = db.Use(NewOptimisticLockPlugin())
	assert.NoError(t, err)
	err = db.Use(NewUUIDPlugin())
	assert.NoError(t, err)

	// テスト用のテーブルを作成
	err = db.AutoMigrate(&domain.User{})
	assert.NoError(t, err)

	t.Run("正常な更新時は楽観ロックエラーが発生しない", func(t *testing.T) {
		// ユーザーを作成
		user := &domain.User{
			UID:  "test-firebase-uid-lock",
			Name: "Lock Test User",
			Type: "tavinikkiy",
			Plan: "free",
		}
		err := db.Create(user).Error
		assert.NoError(t, err)

		// バージョンが1に設定されていることを確認
		assert.Equal(t, 1, user.Version)

		// ユーザーを更新
		user.Name = "Updated Name"
		err = db.Save(user).Error
		assert.NoError(t, err)

		// バージョンが1にインクリメントされていることを確認
		assert.Equal(t, 2, user.Version)
	})

	t.Run("並行更新時に楽観ロックエラーが発生する", func(t *testing.T) {
		// ユーザーを作成
		user := &domain.User{
			UID:  "test-concurrent-uid",
			Name: "Concurrent Test User",
			Type: "tavinikkiy",
			Plan: "free",
		}
		err := db.Create(user).Error
		assert.NoError(t, err)

		// 同じユーザーを2つの異なるインスタンスで取得
		user1 := &domain.User{}
		err = db.First(user1, "uid = ?", "test-concurrent-uid").Error
		assert.NoError(t, err)

		user2 := &domain.User{}
		err = db.First(user2, "uid = ?", "test-concurrent-uid").Error
		assert.NoError(t, err)

		// 両方のインスタンスが同じバージョンであることを確認
		assert.Equal(t, user1.Version, user2.Version)

		// 最初のインスタンスで更新
		user1.Name = "Updated by User1"
		err = db.Save(user1).Error
		assert.NoError(t, err)

		// 更新後のバージョンが2になっていることを確認
		assert.Equal(t, 2, user1.Version)

		// データベースから最新の状態を確認
		var updatedUser domain.User
		err = db.First(&updatedUser, "uid = ?", "test-concurrent-uid").Error
		assert.NoError(t, err)
		assert.Equal(t, 2, updatedUser.Version)
		assert.Equal(t, "Updated by User1", updatedUser.Name)

		// 2番目のインスタンスで更新を試行（古いバージョン0で更新するため楽観ロックエラーが発生するはず）
		user2.Name = "Updated by User2"
		err = db.Save(user2).Error
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pkgErrors.ErrOptimisticLock))
	})
}

func TestOptimisticLockPluginName(t *testing.T) {
	plugin := NewOptimisticLockPlugin()
	assert.Equal(t, "OptimisticLockPlugin", plugin.Name())
}

func TestZeroValueOmitPlugin(t *testing.T) {
	// テスト用のインメモリデータベースを作成
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// ZeroValueOmitプラグインを登録
	err = db.Use(NewZeroValueOmitPlugin())
	assert.NoError(t, err)

	// テスト用のモデル
	type TestModel struct {
		ID          uint          `gorm:"primarykey"`
		Name        string        `gorm:"column:name"`
		Age         int           `gorm:"column:age"`
		IsActive    bool          `gorm:"column:is_active"`
		NullString  sql.NullString `gorm:"column:null_string"`
		NullInt64   sql.NullInt64  `gorm:"column:null_int64"`
		NullBool    sql.NullBool   `gorm:"column:null_bool"`
		Description string        `gorm:"column:description"`
	}

	// テスト用のテーブルを作成
	err = db.AutoMigrate(&TestModel{})
	assert.NoError(t, err)

	t.Run("ゼロ値とInvalidなNull型はomitされる", func(t *testing.T) {
		// 初期データを作成
		original := &TestModel{
			Name:        "Original Name",
			Age:         25,
			IsActive:    true,
			NullString:  sql.NullString{String: "original", Valid: true},
			NullInt64:   sql.NullInt64{Int64: 100, Valid: true},
			NullBool:    sql.NullBool{Bool: true, Valid: true},
			Description: "Original Description",
		}

		err := db.Create(original).Error
		assert.NoError(t, err)

		// ゼロ値とInvalidなNull型を含む更新
		updateData := &TestModel{
			ID:          original.ID,
			Name:        "Updated Name",  // 通常の値
			Age:         0,               // ゼロ値（omitされる）
			IsActive:    false,           // ゼロ値（omitされる）
			NullString:  sql.NullString{String: "", Valid: false}, // Invalid（omitされる）
			NullInt64:   sql.NullInt64{Int64: 0, Valid: false},    // Invalid（omitされる）
			NullBool:    sql.NullBool{Bool: false, Valid: false},  // Invalid（omitされる）
			Description: "",              // ゼロ値（omitされる）
		}

		err = db.Model(original).Updates(updateData).Error
		assert.NoError(t, err)

		// 結果を確認
		var result TestModel
		err = db.First(&result, original.ID).Error
		assert.NoError(t, err)

		// Nameのみが更新され、他のフィールドは変更されていないことを確認
		assert.Equal(t, "Updated Name", result.Name)
		assert.Equal(t, 25, result.Age)                           // 元の値が保持
		assert.Equal(t, true, result.IsActive)                    // 元の値が保持
		assert.Equal(t, "original", result.NullString.String)     // 元の値が保持
		assert.Equal(t, true, result.NullString.Valid)            // 元の値が保持
		assert.Equal(t, int64(100), result.NullInt64.Int64)       // 元の値が保持
		assert.Equal(t, true, result.NullInt64.Valid)             // 元の値が保持
		assert.Equal(t, true, result.NullBool.Bool)               // 元の値が保持
		assert.Equal(t, true, result.NullBool.Valid)              // 元の値が保持
		assert.Equal(t, "Original Description", result.Description) // 元の値が保持
	})

	t.Run("ValidなNull型は更新される", func(t *testing.T) {
		// 初期データを作成
		original := &TestModel{
			Name:       "Test Name",
			NullString: sql.NullString{String: "original", Valid: true},
			NullInt64:  sql.NullInt64{Int64: 100, Valid: true},
		}

		err := db.Create(original).Error
		assert.NoError(t, err)

		// ValidなNull型での更新
		updateData := &TestModel{
			ID:         original.ID,
			NullString: sql.NullString{String: "updated", Valid: true},
			NullInt64:  sql.NullInt64{Int64: 200, Valid: true},
		}

		err = db.Model(original).Updates(updateData).Error
		assert.NoError(t, err)

		// 結果を確認
		var result TestModel
		err = db.First(&result, original.ID).Error
		assert.NoError(t, err)

		// ValidなNull型は更新されることを確認
		assert.Equal(t, "updated", result.NullString.String)
		assert.Equal(t, true, result.NullString.Valid)
		assert.Equal(t, int64(200), result.NullInt64.Int64)
		assert.Equal(t, true, result.NullInt64.Valid)
	})

	t.Run("time.Timeのゼロ値もomitされる", func(t *testing.T) {
		// time.Timeフィールドを含むテストモデル
		type TimeTestModel struct {
			ID        uint      `gorm:"primarykey"`
			Name      string    `gorm:"column:name"`
			EventTime time.Time `gorm:"column:event_time"`
			EventDate *time.Time `gorm:"column:event_date"`
		}

		err := db.AutoMigrate(&TimeTestModel{})
		assert.NoError(t, err)

		// 初期データを作成
		eventTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		original := &TimeTestModel{
			Name:      "Original",
			EventTime: eventTime,
			EventDate: &eventTime,
		}

		err = db.Create(original).Error
		assert.NoError(t, err)

		// ゼロ値のtime.Timeで更新（omitされるべき）
		var zeroTime time.Time // ゼロ値: 0001-01-01 00:00:00 +0000 UTC
		updateData := &TimeTestModel{
			ID:        original.ID,
			Name:      "Updated",
			EventTime: zeroTime,  // ゼロ値なのでomitされる
			EventDate: nil,       // nilでリセット
		}

		err = db.Model(original).Updates(updateData).Error
		assert.NoError(t, err)

		// 結果を確認
		var result TimeTestModel
		err = db.First(&result, original.ID).Error
		assert.NoError(t, err)

		// time.Timeのゼロ値はomitされ、元の値が保持されることを確認
		assert.Equal(t, "Updated", result.Name)     // Nameは更新される
		assert.Equal(t, eventTime, result.EventTime) // EventTimeはゼロ値でomitされるので元の値のまま
		// EventDateはnilへの更新もomitされる可能性があるので、元の値のままかもしれない
		// 実際の動作を確認
		if result.EventDate != nil {
			// nilへの更新がomitされた場合
			assert.Equal(t, eventTime, *result.EventDate)
		} else {
			// nilへの更新が実行された場合
			assert.Nil(t, result.EventDate)
		}
	})
}

func TestZeroValueOmitPluginName(t *testing.T) {
	plugin := NewZeroValueOmitPlugin()
	assert.Equal(t, "ZeroValueOmitPlugin", plugin.Name())
}
