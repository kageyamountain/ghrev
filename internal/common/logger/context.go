package logger

import (
	"context"
	"sync"
)

func NewLogContextMap() *sync.Map {
	return &sync.Map{}
}

type logContextKey struct{}

// WithLogContextMap LogContextMapをセットしたcontextを返す
func WithLogContextMap(ctx context.Context, logContextMap *sync.Map) context.Context {
	return context.WithValue(ctx, logContextKey{}, logContextMap)
}

// LogContextMapFromContext contextからLogContextMapを取り出す
func LogContextMapFromContext(ctx context.Context) (*sync.Map, bool) {
	logContextMap, ok := ctx.Value(logContextKey{}).(*sync.Map)
	return logContextMap, ok
}

type LogType string

const (
	LogTypeApp    LogType = "app_log"
	LogTypeAccess LogType = "access_log"
)

func ChangeLogType(ctx context.Context, logType LogType) {
	logContextMap, ok := LogContextMapFromContext(ctx)
	if ok {
		logContextMap.Store("log_type", logType)
	}
}
