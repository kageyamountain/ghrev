package logger

import (
	"context"
	"log/slog"
	"slices"
	"sync"
)

type logContext struct {
	mutex sync.RWMutex
	attrs []slog.Attr
}

type logContextKey struct{}

// InitLogContext 空のlogContextをセットしたcontextを返す。
func InitLogContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey{}, &logContext{})
}

// SetAttr logContextに属性をセットする。同名キーは上書きする。
func SetAttr(ctx context.Context, attr slog.Attr) {
	logContext, ok := ctx.Value(logContextKey{}).(*logContext)
	if !ok {
		return
	}

	logContext.mutex.Lock()
	defer logContext.mutex.Unlock()

	for i := range logContext.attrs {
		if logContext.attrs[i].Key == attr.Key {
			logContext.attrs[i] = attr
			return
		}
	}
	logContext.attrs = append(logContext.attrs, attr)
}

// ForkLogContext logContextを複製した新しいcontextを返す。
// 以降にセットする属性を、呼び出し元や並行する処理のログに出したくないときに使う。
func ForkLogContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey{}, &logContext{
		attrs: attrsFromContext(ctx),
	})
}

// attrsFromContext logContextが持つ属性のコピーを返す。
func attrsFromContext(ctx context.Context) []slog.Attr {
	logContext, ok := ctx.Value(logContextKey{}).(*logContext)
	if !ok {
		return nil
	}

	logContext.mutex.RLock()
	defer logContext.mutex.RUnlock()

	return slices.Clone(logContext.attrs)
}
