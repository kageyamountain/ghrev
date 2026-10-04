package logger

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"testing"
)

type otherContextKey struct{}

func TestSetAttribute(t *testing.T) {
	t.Parallel()

	t.Run("正常系: 派生させたcontextでセットした場合、親contextからも参照できること", func(t *testing.T) {
		t.Parallel()

		// Arrange
		parent := InitLogContext(context.Background())
		child := context.WithValue(parent, otherContextKey{}, "x")
		want := []slog.Attr{slog.Int(keyPRNumber, 1)}

		// Act
		SetAttribute(child, slog.Int(keyPRNumber, 1))

		// Assert
		got := logContextAttributesFromContext(parent)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("正常系: 同じキーをセットした場合、順序を保ったまま値が上書きされること", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctx := InitLogContext(context.Background())
		SetAttribute(ctx, slog.String(keyRepository, "r1"))
		SetAttribute(ctx, slog.Int(keyPRNumber, 1))
		want := []slog.Attr{slog.String(keyRepository, "r2"), slog.Int(keyPRNumber, 1)}

		// Act
		SetAttribute(ctx, slog.String(keyRepository, "r2"))

		// Assert
		got := logContextAttributesFromContext(ctx)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("正常系: logContextが未設定の場合、何もセットされないこと", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctx := context.Background()

		// Act
		SetAttribute(ctx, slog.Int(keyPRNumber, 1))

		// Assert
		got := logContextAttributesFromContext(ctx)
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})
}

func TestForkLogContext(t *testing.T) {
	t.Parallel()

	t.Run("正常系: ForkLogContextした場合、親の属性を引き継ぎ、子でセットした属性は親に反映されないこと", func(t *testing.T) {
		t.Parallel()

		// Arrange
		parent := InitLogContext(context.Background())
		SetAttribute(parent, slog.String(keyRepository, "r1"))
		wantParent := []slog.Attr{slog.String(keyRepository, "r1")}
		wantChild := []slog.Attr{slog.String(keyRepository, "r1"), slog.Int(keyPRNumber, 1)}

		// Act
		child := ForkLogContext(parent)
		SetAttribute(child, slog.Int(keyPRNumber, 1))

		// Assert
		gotParent := logContextAttributesFromContext(parent)
		if !reflect.DeepEqual(gotParent, wantParent) {
			t.Errorf("parent: got %v, want %v", gotParent, wantParent)
		}
		gotChild := logContextAttributesFromContext(child)
		if !reflect.DeepEqual(gotChild, wantChild) {
			t.Errorf("child: got %v, want %v", gotChild, wantChild)
		}
	})

	t.Run("正常系: goroutineごとにForkLogContextした場合、それぞれのcontextが自身の属性だけを持つこと", func(t *testing.T) {
		t.Parallel()

		// Arrange
		const fanOut = 10
		parent := InitLogContext(context.Background())
		SetAttribute(parent, slog.String(keyRepository, "r1"))
		got := make([][]slog.Attr, fanOut)

		// Act
		var wg sync.WaitGroup
		for i := range fanOut {
			wg.Go(func() {
				ctx := ForkLogContext(parent)
				SetAttribute(ctx, slog.Int(keyPRNumber, i))
				got[i] = logContextAttributesFromContext(ctx)
			})
		}
		wg.Wait()

		// Assert
		for i := range fanOut {
			want := []slog.Attr{slog.String(keyRepository, "r1"), slog.Int(keyPRNumber, i)}
			if !reflect.DeepEqual(got[i], want) {
				t.Errorf("index %d: got %v, want %v", i, got[i], want)
			}
		}
	})
}
