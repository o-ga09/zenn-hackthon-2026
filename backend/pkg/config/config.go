package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/caarlos0/env/v6"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"google.golang.org/api/option"
	"google.golang.org/genai"
)

type Env string
type GenAIConfig string

const CtxEnvKey Env = "env"
const CtxGenAIKey GenAIConfig = "genai_config"

type Config struct {
	Env                       string `env:"ENV" envDefault:"dev"`
	Port                      string `env:"PORT" envDefault:"8080"`
	Database_url              string `env:"DATABASE_URL_FOR_AGENT_TAVINIKKIY" envDefault:""`
	Sentry_DSN                string `env:"SENTRY_DSN" envDefault:""`
	ProjectID                 string `env:"PROJECT_ID" envDefault:"tavinikkiy"`
	CLOUDFLARE_R2_ACCOUNT_ID  string `env:"CLOUDFLARE_R2_ACCOUNT_ID" envDefault:""`
	CLOUDFLARE_R2_ACCESSKEY   string `env:"CLOUDFLARE_R2_ACCESSKEY" envDefault:""`
	CLOUDFLARE_R2_SECRETKEY   string `env:"CLOUDFLARE_R2_SECRETKEY" envDefault:""`
	CLOUDFLARE_R2_BUCKET_NAME string `env:"CLOUDFLARE_R2_BUCKET_NAME" envDefault:"tavinikkiy-local"`
	CLOUDFLARE_R2_PUBLIC_URL  string `env:"CLOUDFLARE_R2_PUBLIC_URL" envDefault:"http://localhost:4566"`
	COOKIE_DOMAIN             string `env:"COOKIE_DOMAIN" envDefault:"localhost"`
	BASE_URL                  string `env:"BASE_URL" envDefault:"http://localhost:3000"`
	GCS_TEMP_BUCKET           string `env:"GCS_TEMP_BUCKET" envDefault:"tavinikkiy-temp"`
	GCS_LOCATION              string `env:"GCS_LOCATION" envDefault:"us-central1"`
	CLOUD_TASKS_QUEUE_NAME    string `env:"CLOUD_TASKS_QUEUE_NAME" envDefault:"tavinikkiy-agent-queue"`
	CLOUD_TASKS_LOCATION      string `env:"CLOUD_TASKS_LOCATION" envDefault:"asia-northeast1"`
	SERVICE_ACCOUNT_EMAIL     string `env:"SERVICE_ACCOUNT_EMAIL" envDefault:""`
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

func InitGenAI(ctx context.Context) context.Context {
	cfg := GetCtxEnv(ctx)
	// Initialize Genkit with the Vertex AI plugin
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.VertexAI{ProjectID: cfg.ProjectID, Location: cfg.GCS_LOCATION}),
		genkit.WithDefaultModel("vertexai/gemini-2.5-flash"),
	)
	return context.WithValue(ctx, CtxGenAIKey, g)
}

func GetGenkitCtx(ctx context.Context) *genkit.Genkit {
	var g *genkit.Genkit
	var ok bool
	if g, ok = ctx.Value(CtxGenAIKey).(*genkit.Genkit); !ok {
		log.Fatal("genkit not found")
	}
	return g
}

// GCSクライアントとGenAIクライアントのシングルトン
var (
	gcsClient     *storage.Client
	gcsClientOnce sync.Once
	gcsClientErr  error

	genaiClient     *genai.Client
	genaiClientOnce sync.Once
	genaiClientErr  error
)

// GetGCSClient はGCSクライアントを取得します（シングルトン）
func GetGCSClient(ctx context.Context) (*storage.Client, error) {
	gcsClientOnce.Do(func() {
		gcsClient, gcsClientErr = storage.NewClient(ctx)
		if gcsClientErr != nil {
			gcsClientErr = fmt.Errorf("failed to create GCS client: %w", gcsClientErr)
		}
	})
	return gcsClient, gcsClientErr
}

// GetGenAIClient はGoogle Gen AIクライアントを取得します（シングルトン）
func GetGenAIClient(ctx context.Context) (*genai.Client, error) {
	genaiClientOnce.Do(func() {
		cfg := GetCtxEnv(ctx)
		genaiClient, genaiClientErr = genai.NewClient(ctx, &genai.ClientConfig{
			Project:  cfg.ProjectID,
			Location: cfg.GCS_LOCATION,
			Backend:  genai.BackendVertexAI,
		})
		if genaiClientErr != nil {
			genaiClientErr = fmt.Errorf("failed to create GenAI client: %w", genaiClientErr)
		}
	})
	return genaiClient, genaiClientErr
}
