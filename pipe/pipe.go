package main

import (
	// "syscall"
	"unsafe"

	// "io"

	// "os"
	"log"

	// "unix"
	"golang.org/x/sys/unix"
)

type workers struct {
	id   []int
	name []string
}

func main() {

	pipefds := make([]int, 2)
	unix.Pipe(pipefds)

	rfd := pipefds[0]
	wfd := pipefds[1]
	defer unix.Close(rfd)
	defer unix.Close(wfd)

	var network workers
	var vm workers
	var qa workers
	network.id = []int{1, 2, 3}
	network.name = []string{"a", "b", "c"}
	vm.id = []int{4, 5, 6}
	vm.name = []string{"aa", "bb", "cc"}
	qa.id = []int{7, 8, 9}
	qa.name = []string{"aaa", "bbb", "ccc"}

	chrw := make(chan bool)
	chwr := make(chan bool)

	log.Println("rfd, wfd: ", rfd, wfd)

	go func() {

		epfd, err := unix.EpollCreate1(0)
		if err != nil {
			log.Fatalf("epoll creation error:%v", err)
		}
		// log.Println("epfd: ", epfd)
		defer unix.Close(epfd)

		err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, rfd, &unix.EpollEvent{
			Events: unix.EPOLLIN | unix.EPOLLET,
			Fd:     int32(rfd),
		})
		if err != nil {
			log.Fatalf("epoll ctrl error:%v", err)
		}

		chrw <- true
		<-chwr

		for i := 0; i < 3; i++ {

			var events []unix.EpollEvent
			n, err := unix.EpollWait(epfd, events[:], -1)
			log.Println("Returned events:", events)
			if n < 0 {
				if err == unix.EAGAIN || err == unix.EINTR {
					continue
				}
				log.Fatalf("epoll wait error:%v", err)
			} else {
				var outval workers
				for _, ev := range events[:n] {
					if int(ev.Fd) == rfd {
						unix.Read(rfd, (*(*[unsafe.Sizeof(outval)]byte)(unsafe.Pointer(&outval)))[:])
						log.Println("out: ", outval)
					}
				}

			}

		}
		chrw <- true

	}()

	<-chrw

	byteval := (*(*[unsafe.Sizeof(network)]byte)(unsafe.Pointer(&network)))[:]
	unix.Write(wfd, byteval)

	byteval = (*(*[unsafe.Sizeof(vm)]byte)(unsafe.Pointer(&vm)))[:]
	unix.Write(wfd, byteval)

	byteval = (*(*[unsafe.Sizeof(qa)]byte)(unsafe.Pointer(&qa)))[:]
	unix.Write(wfd, byteval)

	chwr <- true

	<-chrw

}
