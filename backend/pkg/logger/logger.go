package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"cloud.google.com/go/logging"

	"github.com/o-ga09/zenn-hackthon-2026/pkg/constant"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"

	"go.opentelemetry.io/otel/trace"
)

// cloud logging の Log level 定義
var (
	Severitydefault = slog.Level(logging.Default)
	SeverityInfo    = slog.Level(logging.Info)
	SeverityWarn    = slog.Level(logging.Warning)
	SeverityError   = slog.Level(logging.Error)
	SeverityNotice  = slog.Level(logging.Notice)
)

// traceId , spanId 追加
type traceHandler struct {
	slog.Handler
	projectID string
}

// traceHandler 実装
func (h *traceHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		trace := fmt.Sprintf("projects/%s/traces/%s", h.projectID, sc.TraceID().String())
		r.AddAttrs(slog.String("logging.googleapis.com/trace", trace),
			slog.String("logging.googleapis.com/spanId", sc.SpanID().String()))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttr(attrs []slog.Attr) slog.Handler {
	return &traceHandler{h.Handler.WithAttrs(attrs), h.projectID}
}

func (h *traceHandler) WithGroup(g string) slog.Handler {
	return h.Handler.WithGroup(g)
}

// logger 生成関数
func Logger(ctx context.Context) *slog.Logger {
	replacer := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.MessageKey {
			a.Key = "message"
		}

		if a.Key == slog.LevelKey {
			a.Key = "severity"
			a.Value = slog.StringValue(logging.Severity(a.Value.Any().(slog.Level)).String())
		}

		if a.Key == slog.SourceKey {
			a.Key = "logging.googleapis.com/sourceLocation"
		}

		return a
	}
	cfg := Ctx.GetCfgFromCtx(ctx)
	projectID := cfg.ProjectID
	h := traceHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replacer}), projectID}
	newh := h.WithAttr([]slog.Attr{
		slog.Group("logging.googleapis.com/labels",
			slog.String("app", "MH-API"),
			slog.String("env", cfg.Env),
		),
	})
	logger := slog.New(newh)
	slog.SetDefault(logger)
	return logger
}

func Info(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"requestId", Ctx.GetRequestID(ctx)}, args...)
	slog.Log(ctx, constant.SeverityInfo, msg, allArgs...)
}

func Error(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"requestId", Ctx.GetRequestID(ctx)}, args...)
	slog.Log(ctx, constant.SeverityError, msg, allArgs...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"requestId", Ctx.GetRequestID(ctx)}, args...)
	slog.Log(ctx, constant.SeverityWarn, msg, allArgs...)
}

func Notice(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"requestId", Ctx.GetRequestID(ctx)}, args...)
	slog.Log(ctx, constant.SeverityNotice, msg, allArgs...)
}
