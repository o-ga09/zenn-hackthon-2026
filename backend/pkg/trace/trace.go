package trace

import (
	"context"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

// Span wraps Sentry span
type Span struct {
	sentrySpan *sentry.Span
	startTime  time.Time
}

// SetData sets data on Sentry span
func (s *Span) SetData(key string, value interface{}) {
	if s.sentrySpan != nil {
		s.sentrySpan.SetData(key, value)
	}
}

// StartSpan creates a new Sentry span with the given operation name and description
func StartSpan(ctx context.Context, operation, description string) (*Span, context.Context) {
	span := &Span{
		startTime: time.Now(),
	}

	// Start Sentry span
	if sentrySpan := sentry.SpanFromContext(ctx); sentrySpan != nil {
		childSpan := sentrySpan.StartChild(operation)
		childSpan.Description = description
		span.sentrySpan = childSpan
		ctx = childSpan.Context()
	}

	return span, ctx
}

// StartBusinessLogicSpan creates a span specifically for business logic operations
func StartBusinessLogicSpan(ctx context.Context, operation string) (*Span, context.Context) {
	return StartSpan(ctx, fmt.Sprintf("business.%s", operation), fmt.Sprintf("Execute business logic for %s", operation))
}

// StartValidationSpan creates a span specifically for validation operations
func StartValidationSpan(ctx context.Context, operation string) (*Span, context.Context) {
	return StartSpan(ctx, fmt.Sprintf("validation.%s", operation), fmt.Sprintf("Validate input for %s", operation))
}

// FinishSpan safely finishes Sentry span and records its duration
func FinishSpan(span *Span, err error) {
	if span == nil {
		return
	}

	duration := time.Since(span.startTime)

	// Finish Sentry span
	if span.sentrySpan != nil {
		if err != nil {
			span.sentrySpan.Status = sentry.SpanStatusInternalError
			span.sentrySpan.SetData("error", err.Error())
		} else {
			span.sentrySpan.Status = sentry.SpanStatusOK
		}
		span.sentrySpan.SetData("duration_ms", float64(duration.Microseconds())/1000.0)
		span.sentrySpan.Finish()
	}
}

// WithSpan is a helper function to wrap an operation with a span
func WithSpan(ctx context.Context, operation, description string, fn func(context.Context) error) error {
	span, ctx := StartSpan(ctx, operation, description)
	if span == nil {
		return fn(ctx)
	}

	err := fn(ctx)
	FinishSpan(span, err)
	return err
}
