# テスト計画とデプロイ手順

## 8. テスト計画

### 8.1 単体テスト
- 各ツールの入出力テスト
- エラーハンドリングテスト

### 8.2 統合テスト
- エージェントフローの全体テスト
- 各ツールの連携テスト

### 8.3 E2Eテスト
- ブラウザからの実際の操作テスト
- 複数ファイルのアップロードテスト
- 動画生成の完全フローテスト

### 8.4 パフォーマンステスト
- 同時接続テスト
- 大容量ファイルのテスト
- 長時間動画のテスト

## 9. デプロイ手順

### 9.1 環境変数（例）
VERTEX_AI_PROJECT=your-project-id
VERTEX_AI_LOCATION=asia-northeast1

GCS_BUCKET=your-vlog-bucket

VOICEVOX_ENDPOINT=https://voicevox-xxx.run.app

NEXT_PUBLIC_API_URL=https://backend-xxx.run.app

### 9.2 Cloud Build / デプロイ例
gcloud builds submit --config cloudbuild.yaml

cd frontend && gcloud run deploy vlog-frontend --source .
cd backend && gcloud run deploy vlog-backend --source .
