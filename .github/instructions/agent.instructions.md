---
description: バックエンド開発のためのGoアプリケーションに関する指示
applyTo: "**/*.go,**/go.mod,**/go.sum"
---

# バックエンドアーキテクチャとコーディング規約

## アプリケーション概要

- Go 言語で実装される REST API サーバー
- **Web フレームワーク**: Gin
- **ORM**: Gorm
- **データベース**: TiDB
- **ユーザー認証**: Firebase Auth
- **AI エージェント**: Firebase Genkit
- **AI サービス**: Vertex AI (Gemma、Veo)
- **ホスティング**: Google Cloud Run
- マイクロサービスアーキテクチャを採用

## ディレクトリ構成

Go Standard Project Layout に厳密に従う：

```
backend/
├── cmd/           # メインアプリケーション
├── internal/      # プライベートアプリケーションコード
│   ├── handler/   # HTTPハンドラー (Gin)
│   ├── service/   # ビジネスロジック
│   ├── repository/ # データアクセス層 (Gorm)
│   ├── middleware/ # ミドルウェア
│   └── genkit/    # Firebase Genkit AIエージェント
├── pkg/           # 外部アプリケーションで使用可能なライブラリコード
└── api/           # OpenAPI/Swagger仕様、JSONスキーマファイル
```

## 技術スタック詳細

### Gin Web フレームワーク

- HTTP ルーティング
- ミドルウェア処理

# Firebase Genkit 開発ガイドライン

このドキュメントは、旅行振り返り動画生成機能で Firebase Genkit を使って AI エージェントを開発するためのガイドラインです。TiDB をデータソースとして一元管理し、Genkit を利用して画像解析・動画生成・チャット機能を実装する際の設計方針、実装パターン、セキュリティ、テストのベストプラクティスをまとめます。

## 目的

- TiDB を一次データストアとして利用し、ユーザー/旅行/画像/動画メタデータを管理する
- Firebase Genkit を使って AI タスク（画像解析、場所判定、動画生成、チャット）を実装する
- クライアント（React）→ Firebase（Auth/TiDB/Storage/Genkit）で完結するサーバレスアーキテクチャを推進する
- Cloud Functions を補助として長時間処理や外部連携、セキュリティ検査を行う

## 基本方針

1. TiDB を真のソースオブトゥルースとする。すべてのメタデータ（ユーザー、旅行、画像、動画、チャット履歴）は TiDB に保存する。
2. クライアントは Firebase SDK（Auth/TiDB/Storage）を直接利用して基本 CRUD を行う。権限は TiDB Security Rules で厳格に管理する。
3. 長時間処理（動画生成など）は TiDB ドキュメントのフィールド（例: videos/{videoId}.status）をトリガーにして Cloud Functions または Genkit のワークフローを起動する。
4. Genkit の呼び出しは Cloud Functions 内で行い、シークレット（API キー等）は Secret Manager で管理する。クライアントから直接 Genkit を叩くのは避ける（例外: クライアント向けの軽量なチャットは検討可）。

## データフロー（推奨パターン）

1. 画像アップロード

   - クライアントは Storage に画像をアップロードし、uploads/ または travels/{travelId}/images/{imageId} にメタデータを作成する。
   - Storage のアップロード完了イベントを Cloud Functions が受け取り、画像のサムネイル作成や初期分析ジョブ（Genkit）を開始する。

2. 画像解析 / メタデータ抽出

   - Cloud Functions が Genkit を呼び出し、解析結果を travels/{travelId}/images/{imageId}.analysisData に書き込む。
   - 解析で得た位置情報や日時、オブジェクトタグを元にシーン判定を行い、必要に応じて travel ドキュメントのサマリを更新する。

3. 動画生成

   - クライアントが videos コレクションに generate リクエスト（status: "requested"）を作成。
   - Cloud Functions がトリガーされ、Genkit/VertexAI を用いて動画生成ワークフローを開始する。
   - 生成中は videos/{videoId}.status を "generating" に更新し、進捗を随時更新する。
   - 生成完了時に videos ドキュメントを更新し、Storage に保存された動画の URL を書き込む。

4. チャット / 編集支援
   - Chat は travels/{travelId}/videos/{videoId}/chatHistory に記録。
   - ユーザーからのメッセージは Cloud Functions 経由で Genkit に渡し、応答と編集提案を作成して chatHistory に保存し、必要なら videos ドキュメントを更新する。

## TiDB スキーマ（抜粋）

- users/{userId}
  - email, name, createdAt, updatedAt
- travels/{travelId}
  - userId, title, description, startDate, endDate, status, createdAt, updatedAt
- travels/{travelId}/images/{imageId}
  - originalName, storagePath, url, size, mimeType, width, height, metadata, analysisData, createdAt
- travels/{travelId}/videos/{videoId}
  - title, storagePath, url, thumbnailUrl, duration, width, height, size, status, style, scenes, music, effects, shareUrl, isPublic, createdAt, updatedAt
- travels/{travelId}/videos/{videoId}/chatHistory/{chatId}
  - userId, message, response, suggestions, intent, createdAt

## Genkit 活用パターン

1. 画像解析

   - 入力: Storage の画像 URI またはバイナリ
   - 出力: オブジェクト検出、シーン分類、感情スコア、場所推定、EXIF 解析
   - 保存先: travels/{travelId}/images/{imageId}.analysisData

2. 位置情報の補完

   - 画像に GPS EXIF がない場合は、画像の内容（ランドマーク）から Genkit/外部 API で場所を推定
   - 推定結果は confidence とともに保存し、UI でユーザーが確認できるようにする

3. 動画生成

   - Genkit でテンプレートベースの短編動画（縦型）を生成。必要に応じて VertexAI の動画生成 API を呼ぶ
   - 生成処理は Cloud Functions で管理し、Genkit のワークフロー中に中間結果を Storage/TiDB に保存して進捗を可視化する

4. チャットエージェント
   - Genkit を使い、ユーザーの要求（例: 「この旅行の動画をもっと短く」）を解釈して具体的な VideoUpdate（music/style/order 等）を生成
   - 意図（intent）と信頼度を chatHistory に保存して、必要時にユーザーへ確認を促す

## セキュリティとシークレット管理

- API キーやサービスアカウントは Cloud Functions の環境変数または Secret Manager を使って管理する
- TiDB のルールは最小権限の原則に基づいて設計する（例: travels ドキュメントは owner のみ write 可）
- Storage への直接アップロードは認証済みユーザーに限定し、アップロード先パスとファイル名を厳格に制御する

## エラー / 再試行ポリシー

- Genkit 呼び出しは retryable な失敗（503 など）に対して指数バックオフで再試行する
- 永続的失敗は videos/{videoId}.error フィールドに記録し、運用用の通知を発行する
- Cloud Functions のタイムアウト設計は Genkit の想定処理時間に合わせる（大きすぎるとリソース浪費、小さすぎると失敗）

## ロギングと監視

- 重要なイベント（生成開始・完了・失敗）は TiDB にイベントログを残す
- Cloud Monitoring / Error Reporting / Trace と連携してエンドツーエンドの可観測性を確保する

## テスト戦略

- Firebase Emulator Suite を使ってローカルで TiDB / Auth / Functions / Storage の統合テストを実行する
- Genkit 呼び出しはユニットテストではモック化し、E2E でのみ実際の Genkit を使う（ステージング環境）
- シードデータと teardown スクリプトを CI に組み込み、再現性のあるテストを実行する

## 実装チェックリスト（開発者向け）

- [ ] TiDB のスキーマを設計し、主要コレクション/インデックスを定義した
- [ ] TiDB Security Rules を作成し、owner ベースの書き込み制御を実装した
- [ ] Cloud Functions の IAM と Secret Manager を設定した
- [ ] Storage のバケットとアップロードポリシーを設定した
- [ ] Genkit 呼び出しのラッパーを作成（retry / timeout / error handling を含む）
- [ ] ロギング・監視の基盤を設定（Error Reporting, Monitoring）
- [ ] Firebase Emulator Suite 用の seed/teardown スクリプトを用意した

## よくある実装注意点

- クライアントから直接大容量の動画生成ジョブを開始しない。必ず videos ドキュメントを作り、Cloud Functions 側で実行する。
- Genkit のレスポンスは必ず正規化してから TiDB に保存する（スキーマの変化に対応しやすくするため）。
- TiDB のクエリは課金に直結するため、必要なインデックスを事前に定義しておく。

---

このドキュメントはプロジェクトの初期指針です。実際の開発で得られた知見を元に随時更新してください。
