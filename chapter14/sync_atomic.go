package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// notSafeAdd()
	// atomicAdd()
	syncCond()
}

// race conditionが発生する。結果は947だった
func notSafeAdd() {
	var wg sync.WaitGroup
	wg.Add(1000)
	var counter int64
	for i := 0; i < 1000; i++ {
		go func() {
			counter++
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

// lockをしなくても、複数のgoroutineから安全にアクセスできる
func atomicAdd() {
	var wg sync.WaitGroup
	wg.Add(1000)
	var counter int64
	for i := 0; i < 1000; i++ {
		go func() {
			atomic.AddInt64(&counter, 1)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

func syncCond() {
	var (
		mu      sync.Mutex
		cond    = sync.NewCond(&mu)
		started = 0 // ワーカー開始フラグ
		wg      sync.WaitGroup
	)

	worker := func(id int) {
		defer wg.Done()

		// started=0はワーカーが開始できない状態を表す
		// startedはworkerが実行された順序も表す
		// startedは複数のgoroutineからread/writeされるので、lockで保護する
		// lockが取れたworkerから処理を行える
		mu.Lock()
		for started <= 0 {
			// waitしてる間は、lockを解放してくれるみたい
			cond.Wait()
		}
		fmt.Printf("Worker=%d: started=%d\n", id, started)
		started++
		mu.Unlock()

		time.Sleep(time.Second)
	}

	// 複数のワーカーを準備
	numWorkers := 5
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i)
	}

	// 全てのworkerの初期化が終わるのを待つ的な設定
	time.Sleep(1 * time.Second)
	fmt.Println("Broadcasting start signal...")

	mu.Lock()
	started = 1      // ワーカー開始
	cond.Broadcast() // 全ワーカーに通知
	mu.Unlock()

	// すべてのワーカーが終了するのを待つ
	wg.Wait()
	fmt.Println("All workers completed")
}
