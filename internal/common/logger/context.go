package logger

import (
	"context"
	"sync"
)

type LogContext struct {
	mutex      sync.RWMutex
	attributes map[string]any
}

func NewLogContext() *LogContext {
	return &LogContext{attributes: map[string]any{}}
}

// Set 属性を追加・更新する
func (l *LogContext) Set(key string, value any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.attributes[key] = value
}

// ForEach 全属性を列挙する。読み出し中は書き込みがブロックされ、列挙はatomic snapshotとなる
func (l *LogContext) ForEach(fn func(key string, value any)) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for key, value := range l.attributes {
		fn(key, value)
	}
}

type logContextKey struct{}

// WithLogContext LogContextをセットしたcontextを返す
func WithLogContext(ctx context.Context, logContext *LogContext) context.Context {
	return context.WithValue(ctx, logContextKey{}, logContext)
}

// LogContextFromContext contextからLogContextを取り出す
func LogContextFromContext(ctx context.Context) (*LogContext, bool) {
	logContext, ok := ctx.Value(logContextKey{}).(*LogContext)
	return logContext, ok
}

// SetLogContextAttribute contextに紐づくLogContextに属性を追加する。
// LogContextが未設定の場合は何もしない
func SetLogContextAttribute(ctx context.Context, key string, value any) {
	logContext, ok := LogContextFromContext(ctx)
	if ok {
		logContext.Set(key, value)
	}
}
