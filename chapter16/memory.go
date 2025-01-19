package main

import (
	"fmt"
	"sync"
)

func main() {
	array_slice_map()
	sync_pool()
	memory_arena()
}

func array_slice_map() {
	// 配列は固定長
	array := [3]int{1, 2, 3}
	fmt.Println(array, len(array), cap(array))
	// 配列からスライスを作成
	slice1 := array[:]
	fmt.Println(slice1, len(slice1), cap(slice1))
	// サイズは2だけど、容量は3
	slice2 := array[0:2]
	fmt.Println(slice2, len(slice2), cap(slice2))
	// スライスを作成
	// 裏に配列も作成される
	slice3 := []int{1, 2, 3}
	fmt.Println(slice3, len(slice3), cap(slice3))
	// スライスの容量を指定して作成
	slice4 := make([]int, 3, 4)
	fmt.Println(slice4, len(slice4), cap(slice4))
	// 容量に余裕があるので、appendの時に新しい配列を作成しない
	slice4 = append(slice4, 4)
	fmt.Println(slice4, len(slice4), cap(slice4))
	// 容量が足りないので、新しい配列を作成して、そっちにコピーする
	// mallocが呼ばれる。できるだけ避けたい
	slice4 = append(slice4, 5)
	fmt.Println(slice4, len(slice4), cap(slice4))

	// mapの作成
	m := map[string]int{
		"apple":  100,
		"banana": 200,
	}
	fmt.Println(m)
	// 作成時に容量を指定
	// スライスと同様に容量が足りないとmallocが呼ばれるので、意識が必要
	m2 := make(map[string]int, 10)
	fmt.Println(m2)
}

func sync_pool() {
	// sync.Poolはgoroutine間で共有できる
	// 使い終わったら、GCされる
	pool := sync.Pool{
		New: func() interface{} {
			return 0
		},
	}
	// poolに追加
	pool.Put(1)
	// poolから取得
	fmt.Println(pool.Get())
	// 新規作成
	fmt.Println(pool.Get())
}
