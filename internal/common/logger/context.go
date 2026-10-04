package logger

import (
	"context"
	"log/slog"
	"slices"
	"sync"
)

type logContext struct {
	mutex      sync.RWMutex
	attributes []slog.Attr
}

type logContextKey struct{}

// InitLogContext 空のlogContextをセットしたcontextを返す。
func InitLogContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey{}, &logContext{})
}

// SetAttribute logContextに属性をセットする。同名キーは上書きする。
func SetAttribute(ctx context.Context, attribute slog.Attr) {
	logContext, ok := ctx.Value(logContextKey{}).(*logContext)
	if !ok {
		return
	}

	logContext.mutex.Lock()
	defer logContext.mutex.Unlock()

	for i := range logContext.attributes {
		if logContext.attributes[i].Key == attribute.Key {
			logContext.attributes[i] = attribute
			return
		}
	}
	logContext.attributes = append(logContext.attributes, attribute)
}

// ForkLogContext logContextを複製した新しいcontextを返す。
// 以降にセットする属性を、呼び出し元や並行する処理のログに出したくないときに使う。
func ForkLogContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey{}, &logContext{
		attributes: attributesFromContext(ctx),
	})
}

// attributesFromContext logContextが持つ属性のコピーを返す。
func attributesFromContext(ctx context.Context) []slog.Attr {
	logContext, ok := ctx.Value(logContextKey{}).(*logContext)
	if !ok {
		return nil
	}

	logContext.mutex.RLock()
	defer logContext.mutex.RUnlock()

	return slices.Clone(logContext.attributes)
}
