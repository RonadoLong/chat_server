//+build linux

package echange

import (
	"fmt"
	"reflect"
	"sync"
	"syscall"
)
import "golang.org/x/sys/unix"

func (e *epoller) InitEpoll() error {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		fmt.Println(err)
		return err
	}
	e.fd = fd
	e.connections = sync.Map{}
	return nil
}

func (e *epoller) Add(conn *Socket) (int, error) {
	fd := getSocketFd(conn)
	fmt.Println("===============111", fd, e.fd)
	if fd == -1 {
		return 0, invalidConnectionError
	}
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		fmt.Println(err)
		return e.fd, err
	}
	conn.FD = fd
	fmt.Println("===============1112")
	e.connections.Store(fd, conn)
	fmt.Println("===============113")
	ConnectionPool.Add(fd, conn)
	fmt.Println("===============114")
	return fd, nil
}

func (e epoller) Remove(conn *Socket) error {
	fd := websocketFD(conn)
	if fd == -1 {
		return invalidConnectionError
	}
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.connections.Delete(fd)

	ConnectionPool.remove(conn, fd)
	return nil
}

func (e epoller) Wait() ([]*Socket, error) {
	panic("implement me")
}

func getSocketFd(conn *Socket) int {
	if conn == nil || conn.Conn == nil {
		return  -1
	}
	var (
		connReflectValue = reflect.ValueOf(conn.Conn)
	)
	if !connReflectValue.IsValid() {
		return -1
	}
	var connIndirect = reflect.Indirect(connReflectValue)
	if !connIndirect.IsValid() {
		return -1
	}
	tcpConn := connIndirect.FieldByName("conn")
	if !tcpConn.IsValid() {
		return -1
	}
	fdVal := tcpConn.FieldByName("fd")
	if !fdVal.IsValid() {
		return -1
	}
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}


