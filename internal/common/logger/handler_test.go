package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"reflect"
	"testing"
	"testing/slogtest"
)

const (
	keyRepository = "repository"
	keyPRNumber   = "pr_number"
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
		ctx := InitLogContext(context.Background())
		SetAttr(ctx, slog.String(keyRepository, "r1"))

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
		ctx := InitLogContext(context.Background())
		callee := func(ctx context.Context) {
			SetAttr(ctx, slog.Int(keyPRNumber, 1))
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
		ctx := InitLogContext(context.Background())
		SetAttr(ctx, slog.Int(keyPRNumber, 1))
		SetAttr(ctx, slog.String(keyRepository, "r1"))
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
		ctx := InitLogContext(context.Background())
		SetAttr(ctx, slog.String(keyRepository, "r1"))

		// Act
		logger.InfoContext(ctx, message)

		// Assert
		got := decode(t, &buf)
		if got[keyRepository] != "r1" {
			t.Errorf("got %v, want %v", got[keyRepository], "r1")
		}
	})
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
