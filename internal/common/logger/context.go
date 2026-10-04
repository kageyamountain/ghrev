package logger

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

type logAttributes struct {
	mutex      sync.RWMutex
	attributes []slog.Attr
}

type logAttributesKey struct{}

// InitAttributes 空のログ属性を持つcontextを返す。
func InitAttributes(ctx context.Context) context.Context {
	return context.WithValue(ctx, logAttributesKey{}, &logAttributes{})
}

// ForkAttributes ログ属性を複製した新しいcontextを返す。
func ForkAttributes(ctx context.Context) context.Context {
	return context.WithValue(ctx, logAttributesKey{}, &logAttributes{
		attributes: logAttributesFromContext(ctx),
	})
}

// SetAttribute ログ属性をセットする。同名キーは上書きする。
func SetAttribute(ctx context.Context, attribute slog.Attr) {
	logAttributes, ok := ctx.Value(logAttributesKey{}).(*logAttributes)
	if !ok {
		return
	}

	logAttributes.mutex.Lock()
	defer logAttributes.mutex.Unlock()

	for i := range logAttributes.attributes {
		if logAttributes.attributes[i].Key == attribute.Key {
			logAttributes.attributes[i] = attribute
			return
		}
	}
	logAttributes.attributes = append(logAttributes.attributes, attribute)
}

// logAttributesFromContext ログ属性のコピーを返す。
func logAttributesFromContext(ctx context.Context) []slog.Attr {
	logAttributes, ok := ctx.Value(logAttributesKey{}).(*logAttributes)
	if !ok {
		return nil
	}

	logAttributes.mutex.RLock()
	defer logAttributes.mutex.RUnlock()

	return slices.Clone(logAttributes.attributes)
}

// traceAttributesFromContext contextに有効なspanがあれば、そのtrace_idとspan_idを属性として返す。
func traceAttributesFromContext(ctx context.Context) []slog.Attr {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return nil
	}

	return []slog.Attr{
		slog.String("trace_id", spanContext.TraceID().String()),
		slog.String("span_id", spanContext.SpanID().String()),
	}
}
