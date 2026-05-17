package progress

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const dotInterval = 1 * time.Second

type StopFunc func()

// Start は label を stdout に表示した後、dotInterval ごとにドットを追加して
// 処理中であることを伝える。処理完了後に返り値の StopFunc を呼ぶこと。StopFunc は
// 内部の goroutine の終了を待ち、最後に改行を出力する。
func Start(label string) StopFunc {
	fmt.Print(label)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(dotInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Print(".")
			}
		}
	}()
	return func() {
		cancel()
		wg.Wait()
		fmt.Println()
	}
}
