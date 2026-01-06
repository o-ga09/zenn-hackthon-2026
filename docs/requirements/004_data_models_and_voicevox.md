# データモデル & ずんだモン（VOICEVOX）仕様

## 3.3 データモデル

### 3.3.1 MediaAnalysis
```go
type MediaAnalysis struct {
    FileID      string   `json:"file_id"`
    Type        string   `json:"type"`         // "image" or "video"
    Description string   `json:"description"`  // 全体的な説明
    Objects     []string `json:"objects"`      // 検出されたオブジェクト
    Landmarks   []string `json:"landmarks"`    // 観光地・ランドマーク
    Activities  []string `json:"activities"`   // アクティビティ
    Mood        string   `json:"mood"`         // 雰囲気（楽しい、穏やか、など）
    Timestamp   time.Time `json:"timestamp"`
}
```

### 3.3.2 SubtitleSegment
```go
type SubtitleSegment struct {
    Index int     `json:"index"`
    Start string  `json:"start"`  // "00:00:01,000"
    End   string  `json:"end"`    // "00:00:04,000"
    Text  string  `json:"text"`   // 表示テキスト
}
```

### 3.3.3 VlogOutput
```go
type VlogOutput struct {
    VideoID   string  `json:"video_id"`
    VideoURL  string  `json:"video_url"`
    ShareURL  string  `json:"share_url"`
    Duration  float64 `json:"duration"`
    Thumbnail string  `json:"thumbnail"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 3.3.4 User
```go
type User struct {
    UID           string    `TiDB:"uid"`
    Email         string    `TiDB:"email"`
    DisplayName   string    `TiDB:"display_name"`
    Plan          string    `TiDB:"plan"`           // "free", "monthly", "yearly", "payg"
    TokenBalance  int       `TiDB:"token_balance"`  // 残トークン数
    CreatedAt     time.Time `TiDB:"created_at"`
    UpdatedAt     time.Time `TiDB:"updated_at"`
}
```

### 3.3.5 Subscription
```go
type Subscription struct {
    UID              string    `TiDB:"uid"`
    Plan             string    `TiDB:"plan"`           // "monthly", "yearly"
    Status           string    `TiDB:"status"`         // "active", "cancelled", "expired"
    StripeCustomerID string    `TiDB:"stripe_customer_id"`
    StripeSubID      string    `TiDB:"stripe_subscription_id"`
    CurrentPeriodEnd time.Time `TiDB:"current_period_end"`
    CreatedAt        time.Time `TiDB:"created_at"`
    UpdatedAt        time.Time `TiDB:"updated_at"`
}
```

### 3.3.6 TokenTransaction
```go
type TokenTransaction struct {
    ID          string    `TiDB:"id"`
    UID         string    `TiDB:"uid"`
    Type        string    `TiDB:"type"`         // "purchase", "consumption", "bonus", "refund"
    Amount      int       `TiDB:"amount"`       // トークン数（消費時はマイナス）
    Balance     int       `TiDB:"balance"`      // 取引後の残高
    Description string    `TiDB:"description"`  // "動画生成", "月額プラン付与"など
    Metadata    map[string]interface{} `TiDB:"metadata"`
    CreatedAt   time.Time `TiDB:"created_at"`
}
```

### 3.3.7 Payment
```go
type Payment struct {
    ID              string    `TiDB:"id"`
    UID             string    `TiDB:"uid"`
    Type            string    `TiDB:"type"`           // "token_purchase", "subscription"
    Amount          int       `TiDB:"amount"`         // 金額（円）
    TokensGranted   int       `TiDB:"tokens_granted"` // 付与トークン数
    Status          string    `TiDB:"status"`         // "pending", "completed", "failed"
    StripePaymentID string    `TiDB:"stripe_payment_id"`
    CreatedAt       time.Time `TiDB:"created_at"`
    CompletedAt     *time.Time `TiDB:"completed_at,omitempty"`
}
```

---

## 3.4 ずんだモンキャラクター仕様

### 3.4.1 口調
- 語尾: 「〜なのだ」「〜だよ」「〜なんだ」
- 一人称: 「ボク」
- 性格: 元気、明るい、好奇心旺盛

### 3.4.2 スクリプト例
わぁ！富士山がとってもきれいなのだ！
春の桜と青空が最高なのだよ！
こんな素敵な景色を見られて、ボク幸せなんだ〜！
次はどこに行こうかな？楽しみなのだ！

### 3.4.3 VOICEVOX設定
- 話者ID: 3 (ずんだもん)
- 速度: 1.0（標準）
- 音高: 0.0（標準）
- 抑揚: 1.0（標準）
