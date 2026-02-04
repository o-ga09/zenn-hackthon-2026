# 楽観ロックスキップ機能

## 概要

一括更新などの特定のケースで、楽観的排他制御（Optimistic Lock）のバージョンチェックをスキップする機能を提供します。

## 背景

通常の更新では、楽観ロックプラグインが以下の処理を行います：

1. **更新前（beforeUpdate）**: 現在のバージョン番号をWHERE条件に追加し、バージョンをインクリメント
2. **更新後（afterUpdate）**: 影響行数が0の場合、楽観ロックエラーを発生

しかし、以下のケースでは楽観ロックチェックが不要または不適切な場合があります：

- **一括更新**: 複数のレコードを同時に更新する場合
- **バッチ処理**: 定期的なバックグラウンド処理での更新
- **システム管理操作**: 管理者による強制的な状態変更

## 実装方法

### 1. コンテキストキーの定義

`pkg/context/context.go` に楽観ロックスキップ用のキーを定義しています：

```go
type SkipOptimisticLockKey string
const SkipOptimisticLock SkipOptimisticLockKey = "skipOptimisticLock"
```

### 2. ヘルパー関数

コンテキストに楽観ロックスキップフラグを設定するヘルパー関数を提供：

```go
// WithSkipOptimisticLock は楽観ロックチェックをスキップするコンテキストを返す
func WithSkipOptimisticLock(ctx context.Context) context.Context {
    return context.WithValue(ctx, SkipOptimisticLock, true)
}
```

### 3. プラグインでのチェック

`OptimisticLockPlugin` の `beforeUpdate` と `afterUpdate` メソッドで、コンテキストからフラグをチェック：

```go
func (p *OptimisticLockPlugin) beforeUpdate(db *gorm.DB) {
    // コンテキストから楽観ロックスキップフラグをチェック
    if skip, ok := db.Statement.Context.Value(pkgctx.SkipOptimisticLock).(bool); ok && skip {
        return
    }
    // ... 通常の楽観ロックチェック処理
}

func (p *OptimisticLockPlugin) afterUpdate(db *gorm.DB) {
    // コンテキストから楽観ロックスキップフラグをチェック
    if skip, ok := db.Statement.Context.Value(pkgctx.SkipOptimisticLock).(bool); ok && skip {
        return
    }
    // ... 通常の影響行数チェック処理
}
```

## 使用例

### 通知の一括既読機能

`NotificationRepository.MarkAllAsRead` メソッドでの使用例：

```go
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, notification *domain.Notification) (int64, error) {
    // 楽観ロックチェックをスキップするコンテキストを設定
    ctxWithSkip := Ctx.WithSkipOptimisticLock(ctx)
    
    result := Ctx.GetDB(ctxWithSkip).Model(notification).
        Updates(map[string]interface{}{
            "read": true,
        })

    if result.Error != nil {
        return 0, errors.Wrap(ctx, result.Error)
    }

    return result.RowsAffected, nil
}
```

### その他の使用例

#### バッチ処理での一括更新

```go
func (r *SomeRepository) BatchUpdate(ctx context.Context, ids []string, status string) error {
    ctxWithSkip := Ctx.WithSkipOptimisticLock(ctx)
    
    return Ctx.GetDB(ctxWithSkip).
        Model(&domain.SomeModel{}).
        Where("id IN ?", ids).
        Update("status", status).Error
}
```

#### 管理者による強制更新

```go
func (r *UserRepository) AdminForceUpdate(ctx context.Context, user *domain.User) error {
    ctxWithSkip := Ctx.WithSkipOptimisticLock(ctx)
    
    return Ctx.GetDB(ctxWithSkip).
        Model(user).
        Updates(user).Error
}
```

## 注意事項

### いつ使うべきか

✅ **使用すべき場合:**
- 複数のレコードを一括で更新する場合
- 並行更新の競合が発生しない保証がある場合
- システムによる自動処理での更新
- 管理者による強制的な状態変更

❌ **使用すべきでない場合:**
- ユーザーによる個別レコードの更新
- 並行更新の可能性がある場合
- データの整合性が重要な場合

### セキュリティ考慮事項

- この機能を使用する場合、並行更新による競合が発生しないことを確認してください
- 一括更新の権限チェックを適切に実装してください
- 監査ログに楽観ロックスキップを記録することを検討してください

## まとめ

この機能により、一括更新などの特定のケースで楽観ロックチェックを柔軟にスキップできるようになりました。ただし、データの整合性を保つため、慎重に使用してください。
