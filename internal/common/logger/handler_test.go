package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"reflect"
	"testing"
	"testing/slogtest"

	"go.opentelemetry.io/otel/trace"
)

const (
	keyRepository = "repository"
	keyPRNumber   = "pr_number"
	keyTraceID    = "trace_id"
	keySpanID     = "span_id"
	message       = "msg"
)

func TestHandler_slogtest(t *testing.T) {
	t.Parallel()

	var buf *bytes.Buffer
	newHandler := func(*testing.T) slog.Handler {
		buf = &bytes.Buffer{}
		return NewHandler(slog.NewJSONHandler(buf, nil))
	}
	result := func(t *testing.T) map[string]any {
		t.Helper()
		return decode(t, buf)
	}

	slogtest.Run(t, newHandler, result)
}

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	t.Run("正常系: contextに属性がある場合、ログに出力されること", func(t *testing.T) {
		t.Parallel()

		// Arrange
		var buf bytes.Buffer
		logger := slog.New(NewHandler(slog.NewJSONHandler(&buf, nil)))
		ctx := InitAttributes(context.Background())
		SetAttribute(ctx, slog.String(keyRepository, "r1"))

		// Act
		logger.InfoContext(ctx, message)

		// Assert
		got := decode(t, &buf)
		if got[keyRepository] != "r1" {
			t.Errorf("got %v, want %v", got[keyRepository], "r1")
		}
	})

	t.Run("正常系: 呼び出し先でセットした属性の場合、呼び出し元のログにも出力されること", func(t *testing.T) {
		t.Parallel()

		// Arrange
		var buf bytes.Buffer
		logger := slog.New(NewHandler(slog.NewJSONHandler(&buf, nil)))
		ctx := InitAttributes(context.Background())
		callee := func(ctx context.Context) {
			SetAttribute(ctx, slog.Int(keyPRNumber, 1))
		}

		// Act
		callee(ctx)
		logger.InfoContext(ctx, message)

		// Assert
		got := decode(t, &buf)
		if got[keyPRNumber] != float64(1) {
			t.Errorf("got %v, want %v", got[keyPRNumber], 1)
		}
	})

	t.Run("正常系: ログ呼び出しの引数とcontextが同じキーを持つ場合、引数の値だけが出力されること", func(t *testing.T) {
		t.Parallel()

		// Arrange
		var buf bytes.Buffer
		logger := slog.New(NewHandler(slog.NewJSONHandler(&buf, nil)))
		ctx := InitAttributes(context.Background())
		SetAttribute(ctx, slog.Int(keyPRNumber, 1))
		SetAttribute(ctx, slog.String(keyRepository, "r1"))
		want := map[string]any{keyPRNumber: float64(2), keyRepository: "r1"}

		// Act
		logger.InfoContext(ctx, message, slog.Int(keyPRNumber, 2))

		// Assert
		got := decode(t, &buf)
		for _, key := range []string{slog.TimeKey, slog.LevelKey, slog.MessageKey} {
			delete(got, key)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
		count := bytes.Count(buf.Bytes(), []byte(`"`+keyPRNumber+`"`))
		if count != 1 {
			t.Errorf("key %q appears %d times, want 1", keyPRNumber, count)
		}
	})

	t.Run("正常系: Withで派生させたloggerの場合も、contextの属性が出力されること", func(t *testing.T) {
		t.Parallel()

		// Arrange
		var buf bytes.Buffer
		logger := slog.New(NewHandler(slog.NewJSONHandler(&buf, nil))).With(slog.Int(keyPRNumber, 1))
		ctx := InitAttributes(context.Background())
		SetAttribute(ctx, slog.String(keyRepository, "r1"))

		// Act
		logger.InfoContext(ctx, message)

		// Assert
		got := decode(t, &buf)
		if got[keyRepository] != "r1" {
			t.Errorf("got %v, want %v", got[keyRepository], "r1")
		}
	})
}

func TestHandler_Handle_trace(t *testing.T) {
	t.Parallel()

	const (
		traceID = "0102030405060708090a0b0c0d0e0f10"
		spanID  = "0102030405060708"
	)
	spanContextConfig := trace.SpanContextConfig{
		TraceID:    trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		SpanID:     trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
		TraceFlags: trace.FlagsSampled,
	}

	tests := []struct {
		name    string
		ctx     context.Context
		argAttr []any
		want    map[string]any
	}{
		{
			name: "正常系: contextに有効なspanがある場合、trace_idとspan_idが出力されること",
			ctx:  trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(spanContextConfig)),
			want: map[string]any{keyTraceID: traceID, keySpanID: spanID},
		},
		{
			name: "正常系: ログ属性とspanの両方がある場合、両方の属性が出力されること",
			ctx: func() context.Context {
				ctx := InitAttributes(context.Background())
				SetAttribute(ctx, slog.String(keyRepository, "r1"))
				return trace.ContextWithSpanContext(ctx, trace.NewSpanContext(spanContextConfig))
			}(),
			want: map[string]any{keyRepository: "r1", keyTraceID: traceID, keySpanID: spanID},
		},
		{
			name: "正常系: ログ呼び出しの引数とtrace_idが同じキーの場合、引数の値が出力されること",
			ctx:  trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(spanContextConfig)),
			argAttr: []any{
				slog.String(keyTraceID, "explicit"),
			},
			want: map[string]any{keyTraceID: "explicit", keySpanID: spanID},
		},
		{
			name: "正常系: contextにspanがない場合、trace_idとspan_idは出力されないこと",
			ctx:  context.Background(),
			want: map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			var buf bytes.Buffer
			logger := slog.New(NewHandler(slog.NewJSONHandler(&buf, nil)))

			// Act
			logger.InfoContext(tt.ctx, message, tt.argAttr...)

			// Assert
			got := decode(t, &buf)
			for _, key := range []string{slog.TimeKey, slog.LevelKey, slog.MessageKey} {
				delete(got, key)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func decode(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	var got map[string]any
	err := json.Unmarshal(buf.Bytes(), &got)
	if err != nil {
		t.Fatalf("failed to unmarshal log line. line: %s, err: %v", buf.Bytes(), err)
	}
	return got
}
