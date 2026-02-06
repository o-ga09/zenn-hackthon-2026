# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

## 開発コマンド

### アプリケーション実行
```bash
pnpm dev          # 開発サーバー起動 (http://localhost:3000)
pnpm build        # プロダクションビルド
pnpm start        # プロダクションサーバー起動
```

### コード品質
```bash
pnpm lint         # Oxlintによる型対応リンティング
pnpm type-check   # tsgoによるTypeScript型チェック
pnpm format       # Prettierフォーマットチェック
pnpm format:fix   # 自動フォーマット修正
```

### データベース管理
```bash
pnpm db:up        # データベースコンテナ起動
pnpm db:down      # データベースコンテナ停止・削除
pnpm db:migrate   # データベースマイグレーション実行
pnpm db:seed      # シードデータ投入
pnpm db:reset     # 完全リセット (down + up + migrate + seed)
```

### Prisma
```bash
pnpm prisma:generate  # スキーマからPrismaクライアント生成
pnpm prisma:studio    # Prisma Studio UI起動
```

## プロジェクトアーキテクチャ

### 技術スタック
- **フレームワーク**: Next.js 15 (App Router)
- **言語**: TypeScript (strict mode)
- **スタイリング**: Tailwind CSS v4 + shadcn/uiコンポーネント
- **状態管理**: React Context + TanStack Query
- **認証**: Firebase Auth + セッションクッキー
- **APIクライアント**: Axiosリクエストインターセプター
- **フォーム**: React Hook Form + Zodバリデーション

### ディレクトリ構造

```
frontend/
├── app/              # Next.js App Routerページ
│   ├── dashboard/    # 保護されたダッシュボードルート
│   ├── profile/      # ユーザープロフィールページ (動的ルート)
│   ├── upload/       # メディアアップロードフロー
│   └── videos/       # 動画ライブラリ
├── api/              # APIクライアント層 (app/apiとは別)
│   ├── client.ts     # 認証インターセプター付きAxiosインスタンス
│   ├── types.ts      # 共通API型定義
│   └── *.ts          # ドメイン別APIモジュール
├── components/       # 再利用可能なReactコンポーネント
│   ├── ui/           # shadcn/ui基本コンポーネント
│   ├── layout/       # レイアウトコンポーネント (ヘッダー、フッター等)
│   └── */            # 機能別コンポーネントフォルダ
├── context/          # React Contextプロバイダー
├── hooks/            # カスタムReactフック
├── lib/              # ユーティリティ関数・設定
└── types/            # TypeScript型定義
```

### パスエイリアス

プロジェクトは`tsconfig.json`で設定されたTypeScriptパスエイリアスを使用:
- `@/*` → ルートディレクトリ (例: `@/api/client`, `@/components/ui/button`)

### 認証アーキテクチャ

**Firebase認証 + バックエンドセッションクッキー**

1. **クライアント側フロー**:
   - Firebase Auth (`signInWithPopup`) でGoogleサインイン
   - Firebase IDトークンをバックエンドに送信してセッションクッキー作成
   - バックエンドがトークンを検証して`__session`クッキーを作成
   - クライアント側のFirebaseセッションは即座にクリア (`signOut`)

2. **リクエスト認証**:
   - APIクライアント (`api/client.ts`) がすべてのリクエストをインターセプト
   - `Authorization`ヘッダーにFirebase IDトークンを追加
   - `X-Tavinikkiy-User-Id`ヘッダーにユーザーIDを追加
   - クッキーベース認証のために資格情報を含める

3. **保護されたルート**:
   - ミドルウェア (`middleware.ts`) が`/dashboard`と`/profile`ルートを保護
   - Firebase Admin SDKを使用してサーバー側でセッションクッキーを検証
   - 未認証ユーザーはホームページにリダイレクト

4. **Auth Context** (`context/authContext.tsx`):
   - `user`, `loading`, `login()`, `logout()`, `refetchUser()` を提供
   - ユーザー状態はバックエンドの`/api/auth/user`エンドポイントと同期

### APIクライアントアーキテクチャ

**2つのAxiosインスタンス** (`api/client.ts`):

1. **`apiClient`**: 認証付きリクエスト用
   - Firebase IDトークンを自動的にヘッダーに追加
   - 資格情報 (クッキー) を含む
   - すべての保護されたエンドポイントで使用

2. **`noCredentialApiClient`**: 公開エンドポイント用
   - 認証ヘッダーなし
   - 資格情報なし

**APIモジュールパターン**:
- 各ドメインは独自のAPIファイルを持つ (例: `mediaApi.ts`, `user.ts`, `vlogAPi.ts`)
- 適切なクライアントを使用する型付き関数をエクスポート
- APIロジックを集中化し、エンドポイントを見つけやすくする

### SSEによるリアルタイム更新

アプリはServer-Sent Eventsをリアルタイム進捗更新に使用:

- **`useMediaSSE`フック** (`hooks/useMediaSSE.ts`): メディア分析進捗を監視
- 接続断時に自動再接続
- TanStack Queryと統合してキャッシュ無効化
- メディアアップロード/分析フローで使用

### コンポーネントパターン

**shadcn/ui統合**:
- `components/ui/`に基本UIコンポーネント
- `components.json`で設定済み、簡単に追加可能
- 新しいコンポーネントは`pnpm dlx shadcn add <component>`で追加

**機能別コンポーネント**:
- `_components/`フォルダを使ってルートと同じ場所に配置 (Next.js規約)
- コンポーネントがルートになるのを防ぐ
- 関連コードをまとめて管理

**フォーム処理**:
- React Hook Form + Zodスキーマバリデーション
- 例: `app/upload/_components/form-schema.ts`
- クライアントとサーバー間で一貫したバリデーション

### 状態管理戦略

1. **サーバー状態**: TanStack Query
   - キャッシング、バックグラウンド更新、楽観的更新
   - クエリキーは集中管理 (例: `mediaApi.ts`の`MEDIA_QUERY_KEYS`)

2. **UI状態**: React Context
   - `AuthContext`: グローバル認証状態
   - `NotificationContext`: アプリ全体の通知
   - 機能別コンテキスト (例: `UploadFormContext`)

3. **フォーム状態**: React Hook Form
   - バリデーション付きローカルフォーム状態

### スタイリング規約

- **Tailwind CSS v4**: ネイティブCSS機能を持つ最新バージョン
- **CSS変数**: `globals.css`でテーマカラーを定義
- **コンポーネントバリアント**: `class-variance-authority` (cva) を使用
- **アニメーション**: 複雑なアニメーションは`framer-motion`、シンプルなものはTailwindトランジション

### 環境変数

必須変数 (`.env.example`参照):
- `DATABASE_URL`: データベース接続文字列
- `NEXT_PUBLIC_API_BASE_URL`: バックエンドAPI URL (デフォルト: `http://localhost:8080`)
- Firebase設定変数 (`.env`にあるが`.env.example`にはない)

### Next.js設定

**画像最適化** (`next.config.ts`):
- リモートパターン設定:
  - Cloudflare R2ストレージ
  - Googleユーザー写真
  - Firebase Storage

**ミドルウェア**:
- 認証が必要なルートを保護
- Firebase Admin SDKでサーバー側トークン検証

### 型安全性

- **Strict TypeScript**: すべてのコンパイラオプションで厳密な型付けを強制
- **API型**: `api/types.ts`で集中管理
- **Zodスキーマ**: ランタイムバリデーションがTypeScript型と一致
- **型チェック**: コミット前に`pnpm type-check`を実行

### 主要なサードパーティ統合

- **Firebase**: 認証 (クライアント) とAdmin SDK (サーバー側検証)
- **Cloudflare R2**: メディアストレージ (画像と動画)
- **Prisma**: データベースORM (スキーマは`schemes/prisma/schema.prisma`)
- **TanStack Query**: サーバー状態管理とキャッシング
- **shadcn/ui**: Radix UIベースのコンポーネントライブラリ

## 開発ワークフロー

1. **開発開始**:
   ```bash
   pnpm db:up              # 必要に応じてデータベース起動
   pnpm dev                # 開発サーバー起動
   ```

2. **コミット前**:
   ```bash
   pnpm type-check         # TypeScript検証
   pnpm lint               # コード品質チェック
   pnpm format             # フォーマットチェック
   ```

3. **データベース変更**:
   - `schemes/prisma/schema.prisma`でPrismaスキーマを更新
   - `pnpm db:migrate`で変更を適用
   - `pnpm prisma:generate`でPrismaクライアントを更新

4. **新しいAPIエンドポイント追加**:
   - `api/types.ts`に型を追加
   - 関連するAPIモジュールを作成または更新 (例: `api/mediaApi.ts`)
   - 認証付きエンドポイントには`apiClient`、公開エンドポイントには`noCredentialApiClient`を使用
   - 適切なエラーハンドリング付きの型付き関数をエクスポート

5. **新しいUIコンポーネント追加**:
   - shadcn/uiコンポーネントの場合: `pnpm dlx shadcn add <component>`
   - カスタムコンポーネントの場合: `components/`またはルート固有の`_components/`に配置
   - propsにはTypeScriptインターフェースを使用
   - Tailwindスタイリング規約に従う
