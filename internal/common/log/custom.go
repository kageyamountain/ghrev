package log

import (
	"context"
	"log/slog"
	"sync"
)

type logContextKey string

const LogContextKey logContextKey = "LogContext"

type LogType string

const (
	LogTypeApp    LogType = "app_log"
	LogTypeAccess LogType = "access_log"
)

type CustomLogHandler struct {
	slog.Handler
}

func NewCustomLogHandler(handler slog.Handler) *CustomLogHandler {
	return &CustomLogHandler{
		Handler: handler,
	}
}

// Handle contextにセットされたLogContextからログ出力フィールドを追加する
func (h *CustomLogHandler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic // slogのinterface仕様なので第2引数はポインタ型にできない
	// contextからLogContextを取得
	logContext, ok := ctx.Value(LogContextKey).(*sync.Map)
	if !ok {
		return h.Handler.Handle(ctx, r)
	}

	// LogContextの全エントリをログ出力属性に追加
	logContext.Range(func(key, value interface{}) bool {
		keyStr, ok2 := key.(string)
		if ok2 {
			r.AddAttrs(slog.Attr{Key: keyStr, Value: slog.AnyValue(value)})
		}
		return true
	})

	// info: OpenTelemetryと連携する場合は有効化
	//// OpenTelemetry trace情報を追加
	// if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
	//	spanCtx := span.SpanContext()
	//	r.AddAttrs(
	//		slog.String("trace_id", spanCtx.TraceID().String()),
	//		slog.String("span_id", spanCtx.SpanID().String()),
	//	)
	//}

	return h.Handler.Handle(ctx, r)
}
