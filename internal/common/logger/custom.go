package logger

import (
	"context"
	"log/slog"
)

type CustomLogHandler struct {
	slog.Handler
}

func NewCustomLogHandler(handler slog.Handler) *CustomLogHandler {
	return &CustomLogHandler{
		Handler: handler,
	}
}

// Handle contextにセットされたLogContextMapからログ出力フィールドを追加する
func (h *CustomLogHandler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic // slogのinterface仕様なので第2引数はポインタ型にできない
	logContextMap, ok := LogContextMapFromContext(ctx)
	if !ok {
		return h.Handler.Handle(ctx, r)
	}

	r = r.Clone()
	for key, value := range logContextMap.Range {
		keyStr, ok2 := key.(string)
		if ok2 {
			r.AddAttrs(slog.Any(keyStr, value))
		}
	}

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

func (h *CustomLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomLogHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *CustomLogHandler) WithGroup(name string) slog.Handler {
	return &CustomLogHandler{
		Handler: h.Handler.WithGroup(name),
	}
}
