package logger

import (
	"context"
	"log/slog"
)

type Handler struct {
	innerHandler slog.Handler
}

func NewHandler(innerHandler slog.Handler) *Handler {
	return &Handler{
		innerHandler: innerHandler,
	}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic // slogのinterface仕様なので第2引数はポインタ型にできない
	contextAttributes := attributesFromContext(ctx)
	if len(contextAttributes) == 0 {
		return h.innerHandler.Handle(ctx, r)
	}

	// ログ呼び出しの引数とlogContextに同じキーがある場合は、引数の値を優先してlogContextの属性を出力しない
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

// WithGroup 以後のログでは、logContextの属性もトップレベルではなくグループ内に出力される。
// グループの処理はinnerHandlerに委譲しており、トップレベルに固定するにはグループを自前で保持する必要があり実装が複雑になる。
// 属性をまとめたい場合はWithGroupではなく、ログ呼び出しの引数にslog.Groupを使うようにしてください。
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		innerHandler: h.innerHandler.WithGroup(name),
	}
}
