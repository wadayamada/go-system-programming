package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// first_start_first_out()
	// parallel_blocking_sync()
	// parallel_non_blocking_sync()
	parallel_blocking_async()
}

// チャネルのチャネルを使うことで、先に処理が始まったgoroutineの結果を先にチャネルで受け取る
func first_start_first_out() {
	task_count := 10
	ch := make(chan chan int, task_count)

	var wg sync.WaitGroup
	wg.Add(task_count)

	// goroutineが全て終わったらチャネルを閉じる
	go func() {
		wg.Wait()
		close(ch)
	}()

	// task=0から順にタスクを実行する
	for i := 0; i < task_count; i++ {
		task := make(chan int)
		ch <- task
		go task_run(i, task, &wg)
	}

	// goroutineなのでタスクの開始・終了の順序はバラバラだが、
	// チャネルのチャネルを使うことで、処理開始順で結果を受け取ることができる
	for task := range ch {
		fmt.Println("result ", <-task)
	}
}

func task_run(number int, task chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("start ", number)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("done", number)
	task <- number
}

func parallel_blocking_sync() {

	var wg sync.WaitGroup
	wg.Add(2)

	// 2つのgoroutineでカウント処理を並列に行う
	// 2つのカウント処理は非同期的に行われている
	counter1 := 0
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			counter1++
		}
	}()
	counter2 := 0
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			counter2++
		}
	}()

	// WaitGroupで2つのgoroutineが終わるのを同期的に待つ
	// ブロッキングで待つ
	wg.Wait()
	fmt.Println(counter1 + counter2)
}

func parallel_non_blocking_sync() {

	ch1 := make(chan int)
	ch2 := make(chan int)

	// 2つのgoroutineでカウント処理を並列に行う
	// 2つのカウント処理は非同期的に行われている
	counter1 := 0
	go func() {
		for i := 0; i < 10; i++ {
			counter1++
		}
		time.Sleep(5 * time.Second)
		ch1 <- counter1
	}()
	counter2 := 0
	go func() {
		for i := 0; i < 10; i++ {
			counter2++
		}
		time.Sleep(5 * time.Second)
		ch2 <- counter2
	}()

	// select defaultで、ノンブロッキング・同期で待つ
	// defaultあるから非同期だけど、if文で両方の結果を待ってるから、同期的と言える
	result := 0
	finish_counter1 := false
	finish_counter2 := false
	for {
		select {
		case counter1 = <-ch1:
			finish_counter1 = true
		case counter2 = <-ch2:
			finish_counter2 = true
		default:
			if finish_counter1 && finish_counter2 {
				result = counter1 + counter2
				fmt.Println("result", result)
				return
			}
			time.Sleep(1 * time.Second)
			fmt.Println("no goroutine finished")
		}
	}
}

func parallel_blocking_async() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// 2つのgoroutineでカウント処理を並列に行う
	// 2つのカウント処理は非同期的に行われている
	counter1 := 0
	go func() {
		for i := 0; i < 10; i++ {
			counter1++
		}
		time.Sleep(3 * time.Second)
		ch1 <- counter1
	}()
	counter2 := 0
	go func() {
		for i := 0; i < 10; i++ {
			counter2++
		}
		time.Sleep(5 * time.Second)
		ch2 <- counter2
	}()

	// selectで、ブロッキング・非同期で待つ
	// counter1, counter2を合流させてないので、この二つは非同期
	for {
		select {
		case counter1 = <-ch1:
			fmt.Println("counter1", counter1)
		case counter2 = <-ch2:
			fmt.Println("counter2", counter2)
			return
		}
	}
}
