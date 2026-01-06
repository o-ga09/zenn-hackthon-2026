package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"github.com/caarlos0/env/v6"
	"google.golang.org/api/option"
)

type Env string

const CtxEnvKey Env = "env"

type Config struct {
	Env                       string `env:"ENV" envDefault:"dev"`
	Port                      string `env:"PORT" envDefault:"8080"`
	Database_url              string `env:"DATABASE_URL" envDefult:""`
	Sentry_DSN                string `env:"SENTRY_DSN" envDefult:""`
	ProjectID                 string `env:"PROJECT_ID" envDefult:""`
	CLOUDFLARE_R2_ACCOUNT_ID  string `env:"CLOUDFLARE_R2_ACCOUNT_ID" envDefult:""`
	CLOUDFLARE_R2_ACCESSKEY   string `env:"CLOUDFLARE_R2_ACCESSKEY" envDefult:""`
	CLOUDFLARE_R2_SECRETKEY   string `env:"CLOUDFLARE_R2_SECRETKEY" envDefult:""`
	CLOUDFLARE_R2_BUCKET_NAME string `env:"CLOUDFLARE_R2_BUCKET_NAME" envDefult:""`
	COOKIE_DOMAIN             string `env:"COOKIE_DOMAIN" envDefault:"localhost"`
}

func New(ctx context.Context) (context.Context, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return context.WithValue(ctx, CtxEnvKey, cfg), nil
}

var (
	firebaseApp     *firebase.App
	firebaseAppOnce sync.Once
	firebaseAppErr  error
)

// GetFirebaseApp はFirebaseアプリケーションインスタンスを取得します（シングルトン）
func GetFirebaseApp(ctx context.Context) (*firebase.App, error) {
	firebaseAppOnce.Do(func() {
		var credentialsPath string
		env := GetCtxEnv(ctx)

		if env.Env == "local" {
			// 環境変数からサービスアカウントのパスを取得
			credentialsPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
			opt := option.WithCredentialsFile(credentialsPath)
			firebaseApp, firebaseAppErr = firebase.NewApp(ctx, nil, opt)
		} else {
			// Cloud Run環境: デフォルトの認証情報を使用
			firebaseApp, firebaseAppErr = firebase.NewApp(ctx, nil)
		}

		if firebaseAppErr != nil {
			firebaseAppErr = fmt.Errorf("failed to initialize firebase app: %w", firebaseAppErr)
		}
	})

	return firebaseApp, firebaseAppErr
}

func GetCtxEnv(ctx context.Context) *Config {
	var cfg *Config
	var ok bool
	if cfg, ok = ctx.Value(CtxEnvKey).(*Config); !ok {
		log.Fatal("config not found")
	}
	return cfg
}
