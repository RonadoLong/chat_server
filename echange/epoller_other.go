//+build !linux

package echange

import (
	"reflect"
	"sync"
)

func (e *epoller) InitEpoll() error {
	e.fd = 34
	e.connections = sync.Map{}
	return nil
}

func (e *epoller) Add(conn *Socket) (int, error) {
	fd := websocketFD(conn)
	ConnectionPool.Add(fd, conn)
	return fd, nil
}

func (e epoller) Remove(conn *Socket) error {
	fd := websocketFD(conn)
	e.connections.Delete(fd)
	ConnectionPool.remove(conn, fd)
	return nil
}

func (e epoller) Wait() ([]*Socket, error) {
	return nil, nil
}

func websocketFD(conn *Socket) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn.Conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func (e *epoller) GetConnection(fd int) *Socket {
	return &Socket{}
}
