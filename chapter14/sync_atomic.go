package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	notSafeAdd()
	atomicAdd()
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
