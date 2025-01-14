package main

import (
	"fmt"
	"syscall"
)

func main() {
	kq, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}
	fd, err := syscall.Open("./test.txt", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	// ファイルの変更を通知する
	ev := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_VNODE,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE | syscall.EV_ONESHOT,
		Fflags: syscall.NOTE_DELETE | syscall.NOTE_WRITE,
		Data:   0,
		Udata:  nil,
	}
	for {
		events := make([]syscall.Kevent_t, 10)
		// ファイルの変更があるまで、ここでブロッキング
		// ディレクトリを監視して、任意のファイルが作成されるのを待って、作成されたら通知を受け取って、書き込みとかもできる
		nev, err := syscall.Kevent(kq, []syscall.Kevent_t{ev}, events, nil)
		if err != nil {
			panic(err)
		}
		for i := 0; i < nev; i++ {
			fmt.Printf("Event: %+v", events[i])
		}
	}
}
