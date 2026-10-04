package logger

import (
	"context"
	"log/slog"
	"slices"
)

// Handler logContextを扱うためのslog.Handlerのラッパー。
//
// 制約:
// slogはlogContextの属性をログ呼び出しの引数と同じ扱いにするため、With・WithGroupは使用しないこと。
// WithGroup後はlogContextの属性がグループ内に入り、Withで付けた属性とlogContextの属性に同じキーがあると両方が出力されてしまう。
// WithではなくSetAttribute、WithGroupではなくslog.Groupを使用すること。
// この制約はforbidigoで、slogの書き方はsloglintで強制することを推奨する。
type Handler struct {
	innerHandler slog.Handler
}

func NewHandler(innerHandler slog.Handler) *Handler {
	return &Handler{
		innerHandler: innerHandler,
	}
}

// Handle logContextの属性とtraceの属性をログに追加して、innerHandlerに渡す。
func (h *Handler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic // slogのinterface仕様なので第2引数はポインタ型にできない
	logContextAttributes := logContextAttributesFromContext(ctx)
	traceAttributes := traceAttributesFromContext(ctx)
	contextAttributes := slices.Concat(logContextAttributes, traceAttributes)
	if len(contextAttributes) == 0 {
		return h.innerHandler.Handle(ctx, r)
	}

	// ログ呼び出しの引数とcontextの属性に同じキーがある場合は、引数の値を優先してcontextの属性を出力しない
	argAttributeKeys := make(map[string]struct{}, r.NumAttrs())
	for argAttribute := range r.Attrs {
		argAttributeKeys[argAttribute.Key] = struct{}{}
	}

	r = r.Clone()
	for _, contextAttribute := range contextAttributes {
		_, ok := argAttributeKeys[contextAttribute.Key]
		if ok {
			continue
		}
		r.AddAttrs(contextAttribute)
	}

	return h.innerHandler.Handle(ctx, r)
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.innerHandler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attributes []slog.Attr) slog.Handler {
	return &Handler{
		innerHandler: h.innerHandler.WithAttrs(attributes),
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		innerHandler: h.innerHandler.WithGroup(name),
	}
}
