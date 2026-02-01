package retry

import (
	"context"
	"time"
)

// Config はリトライ動作の設定
type Config struct {
	MaxRetries     int           // 最大リトライ回数
	InitialBackoff time.Duration // 初期バックオフ時間
	MaxBackoff     time.Duration // 最大バックオフ時間
	Multiplier     float64       // バックオフ乗数
}

// DefaultConfig はデフォルトのリトライ設定
// - 最大3回リトライ
// - 初期バックオフ1秒
// - 最大バックオフ10秒
// - 指数バックオフ（2倍）
var DefaultConfig = Config{
	MaxRetries:     3,
	InitialBackoff: 1 * time.Second,
	MaxBackoff:     10 * time.Second,
	Multiplier:     2.0,
}

// Do は指定された関数を実行し、失敗時にExponential Backoffでリトライする
// ctx: コンテキスト（キャンセル対応）
// cfg: リトライ設定
// fn: 実行する関数
func Do(ctx context.Context, cfg Config, fn func() error) error {
	var lastErr error
	backoff := cfg.InitialBackoff

	for i := 0; i <= cfg.MaxRetries; i++ {
		// 2回目以降はバックオフを待つ
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}

			// 次回のバックオフ時間を計算
			backoff = time.Duration(float64(backoff) * cfg.Multiplier)
			if backoff > cfg.MaxBackoff {
				backoff = cfg.MaxBackoff
			}
		}

		// 関数実行
		if err := fn(); err != nil {
			lastErr = err
			continue
		}

		// 成功
		return nil
	}

	// 全てのリトライが失敗
	return lastErr
}
