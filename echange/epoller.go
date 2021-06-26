package echange

import (
	"errors"
	"fmt"
	"sync"
	"syscall"
)

type (
	EpollerI interface {
		//Add 新增
		Add(conn *Socket) (int, error)
		//Remove 移除连接
		Remove(conn *Socket) error
		//Wait 等待连接
		Wait() ([]*Socket, error)
	}

	epoller struct {
		fd          int
		connections sync.Map
	}
)

var (
	epoll                  = epoller{}
	invalidConnectionError = errors.New("invalid socket ,connection is CLOSED")
)

func EpollerServer() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	var err error
	if err = epoll.InitEpoll(); err != nil {
		panic(err)
	}
	fmt.Println("rlimit", rLimit.Cur, " 建立epoll服务器")
	//go start()
}

func GetSocketPusher() *epoller {
	return &epoll
}
