# 008_メディア分析結果の表示・編集機能

## 概要

素材メディア一覧の各画像をクリックすると分析結果が表示され、分析結果が気に入らない場合は編集して保存できる機能を実装する。

## 目的

- ユーザーがAIによる分析結果を確認できるようにする
- 分析結果が不正確な場合、ユーザーが手動で修正できるようにする
- 修正した内容を楽観的UI更新パターンで即座に反映し、UXを向上させる

## 機能仕様

### 1. 分析結果の表示

#### 表示対象
- **Description**: メディアの全体的な説明（最大500文字）
- **Mood**: 雰囲気（プリセット: 楽しい/穏やか/エキサイティング/ロマンチック）
- **Objects**: 検出されたオブジェクト（タグ形式）
- **Landmarks**: ランドマーク（タグ形式）
- **Activities**: アクティビティ（タグ形式）

#### UI要件
- メディア一覧で分析完了（`status === 'completed'`）のメディアのみクリック可能
- 分析中（`status === 'pending' || 'uploading'`）のメディアはクリック無効、ツールチップで「分析中です」と表示
- 分析失敗（`status === 'failed'`）のメディアはクリック無効、ツールチップで「分析失敗」と表示
- 右クリックまたは「詳細」ボタンクリックでダイアログを開く

### 2. 分析結果の編集

#### 編集可能項目
1. **Description**: テキストエリアで編集、最大500文字制限
2. **Mood**: セレクトボックスで選択（プリセット: 楽しい/穏やか/エキサイティング/ロマンチック）
3. **Objects**: TagInputコンポーネントでタグ追加・削除
4. **Landmarks**: TagInputコンポーネントでタグ追加・削除
5. **Activities**: TagInputコンポーネントでタグ追加・削除

#### バリデーション
- **Description**: 最大500文字
- **Mood**: プリセット値のいずれか（必須ではない）
- **Tags (Objects/Landmarks/Activities)**: 各タグの長さは最大50文字

### 3. 保存処理

#### 楽観的UI更新パターン
1. **保存ボタンクリック**: 即座にローカルキャッシュを更新し、UIに反映
2. **API呼び出し**: バックグラウンドでPUT `/api/media/:id/analytics`を呼び出し
3. **成功時**: `toast.success('タグを更新しました')` を表示し、ダイアログを閉じる
4. **失敗時**: ローカルキャッシュをロールバックし、`toast.error('タグの更新に失敗しました')` を表示

## アーキテクチャ

### バックエンド

#### API エンドポイント

##### 分析結果取得
- **Method**: `GET`
- **Path**: `/api/media/:id/analytics`
- **Response**:
```json
{
  "file_id": "string",
  "description": "string",
  "mood": "string",
  "objects": ["string"],
  "landmarks": ["string"],
  "activities": ["string"]
}
```

##### 分析結果更新
- **Method**: `PUT`
- **Path**: `/api/media/:id/analytics`
- **Request**:
```json
{
  "description": "string (optional, max 500)",
  "mood": "string (optional)",
  "objects": ["string (optional, max 50 each)"],
  "landmarks": ["string (optional, max 50 each)"],
  "activities": ["string (optional, max 50 each)"]
}
```
- **Response**: 更新後の分析結果（GET と同じ形式）

#### データモデル

```go
type MediaAnalytics struct {
    BaseModel
    FileID      string           // メディアID
    Description string           // 全体的な説明
    Mood        string           // 雰囲気
    Objects     []DetectedObject // 検出オブジェクト
    Landmarks   []Landmark       // ランドマーク
    Activities  []Activity       // アクティビティ
}

type DetectedObject struct {
    BaseModel
    MediaAnalyticsID string
    Name             string
}

type Landmark struct {
    BaseModel
    MediaAnalyticsID string
    Name             string
}

type Activity struct {
    BaseModel
    MediaAnalyticsID string
    Name             string
}
```

#### 更新処理
- トランザクション内で既存の関連データ（Objects, Landmarks, Activities）を削除
- 新しいデータを挿入
- MediaAnalytics本体を更新（Description, Mood）

### フロントエンド

#### コンポーネント構成

```
PhotoUpload.tsx (既存)
├── MediaAnalyticsDialog.tsx (新規)
│   ├── Dialog (shadcn/ui)
│   ├── Textarea (Description編集)
│   ├── Select (Mood選択)
│   └── TagInput (Objects/Landmarks/Activities編集)
└── Tooltip (分析中/失敗時の表示)
```

#### 状態管理
- **TanStack Query**: サーバー状態管理
  - `useGetMediaAnalytics(fileId)`: 分析結果取得
  - `useUpdateMediaAnalytics(fileId)`: 分析結果更新（楽観的更新）
- **useState**: ローカル状態管理
  - `selectedMedia`: クリックされたメディア
  - `isDialogOpen`: ダイアログの開閉状態

#### TagInputコンポーネント
- Enterキーでタグ追加
- Xボタンでタグ削除
- Backspaceキーで最後のタグを削除
- trim処理による重複チェック（大文字小文字やひらがな/カタカナの正規化は行わない）

## UX フロー

1. **メディア一覧表示**
   - 分析完了メディア: クリック可能、ホバー時に「詳細」ボタン表示
   - 分析中メディア: クリック無効、オーバーレイで「分析中...」表示、ツールチップで「分析中です」
   - 分析失敗メディア: クリック無効、バッジで「分析失敗」表示、ツールチップで「分析に失敗しました」

2. **ダイアログ表示**
   - 右クリックまたは「詳細」ボタンクリックでダイアログを開く
   - ローディング状態: スピナーを表示
   - エラー状態: エラーメッセージを表示

3. **編集**
   - Description: テキストエリアで編集、文字数カウンター表示
   - Mood: セレクトボックスで選択
   - Tags: TagInputでタグ追加・削除

4. **保存**
   - 保存ボタンクリック → 即座にUIを更新（楽観的更新）
   - API呼び出し成功 → トースト通知「タグを更新しました」、ダイアログを閉じる
   - API呼び出し失敗 → ロールバック、トースト通知「タグの更新に失敗しました」

5. **キャンセル**
   - キャンセルボタンクリック → 元の値に戻す、ダイアログを閉じる

## テスト項目

### バックエンド
- [ ] GET `/api/media/:id/analytics` で分析結果を取得できる
- [ ] PUT `/api/media/:id/analytics` で分析結果を更新できる
- [ ] 存在しないメディアIDでエラーが返される
- [ ] バリデーションエラーが適切に処理される
- [ ] トランザクション内で関連データが正しく更新される

### フロントエンド
- [ ] 分析完了メディアをクリックしてダイアログが開く
- [ ] 分析中メディアはクリック無効で、ツールチップが表示される
- [ ] 分析失敗メディアはクリック無効で、ツールチップが表示される
- [ ] TagInputでタグの追加・削除ができる
- [ ] Moodセレクトボックスでプリセット値を選択できる
- [ ] Descriptionが500文字制限される
- [ ] 保存ボタンクリック時に楽観的更新が行われる
- [ ] API呼び出し成功時にトースト通知が表示される
- [ ] API呼び出し失敗時にロールバックが行われる

## 非機能要件

### パフォーマンス
- 楽観的UI更新により、保存時の体感速度を向上
- キャッシュ戦略: 5分間staleTime、10分間gcTime

### アクセシビリティ
- ツールチップでクリック不可の理由を明示
- キーボード操作対応（Enter/Backspace）
- aria属性による適切なラベリング

### エラーハンドリング
- API呼び出し失敗時の自動ロールバック
- トースト通知によるエラーフィードバック
- ローディング状態とエラー状態の明示

## 実装ファイル一覧

### バックエンド
- `backend/internal/domain/media_analytics.go`: リポジトリインターフェースにUpdateメソッド追加
- `backend/internal/infra/database/mysql/media_analytics.go`: Update実装追加
- `backend/internal/handler/request/media_analytics.go`: リクエストDTO追加
- `backend/internal/handler/response/media_analytics.go`: レスポンスDTO追加
- `backend/internal/handler/media.go`: GetAnalytics, UpdateAnalyticsメソッド追加
- `backend/internal/server/router.go`: ルーティング追加
- `backend/internal/server/server.go`: ImageServer初期化にanalyticsRepo追加

### フロントエンド
- `frontend/components/ui/tag-input.tsx`: TagInputコンポーネント新規作成
- `frontend/components/ui/select.tsx`: shadcn/uiからインストール
- `frontend/components/ui/tooltip.tsx`: shadcn/uiからインストール
- `frontend/api/types.ts`: MediaAnalyticsResponse, UpdateMediaAnalyticsRequest追加
- `frontend/api/mediaApi.ts`: useGetMediaAnalytics, useUpdateMediaAnalytics追加
- `frontend/app/upload/_components/MediaAnalyticsDialog.tsx`: ダイアログコンポーネント新規作成
- `frontend/app/upload/_components/PhotoUpload.tsx`: クリックハンドラーとツールチップ追加

## 今後の拡張案

1. **一括編集機能**: 複数メディアの分析結果を一括編集
2. **プリセット管理**: Moodのプリセット値をユーザーがカスタマイズ可能に
3. **AI再分析**: 分析失敗メディアの再分析機能
4. **履歴管理**: 分析結果の編集履歴を保存・参照
5. **タグサジェスト**: 既存タグをもとに入力補完機能を追加
