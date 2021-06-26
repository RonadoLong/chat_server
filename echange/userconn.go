package echange

import (
	"container/list"
	"fmt"
	"sync"
)

//UserConnection 用户连接管理
type UserConnection struct {
	connections sync.Map // map[string]*list.List //用户的socket连接
}

type connectionListObject struct {
	FD int //文件描述符
	// Platform string //平台
	UserID string //
	Conn   *Socket
}

var (
	//ConnectionPool 连接池
	ConnectionPool = UserConnection{
		connections: sync.Map{},
	}
)

//Add 添加新的连接
func (u *UserConnection) Add(fd int, socket *Socket) {
	var connObj = connectionListObject{
		FD:     fd,
		UserID: socket.UserID,
		Conn:   socket,
	}
	var connList = u.GetConnections(socket.UserID)
	fmt.Println("connList, ", connList)
	if connList == nil {
		connList = list.New()
	} else {
		for item := connList.Front(); item != nil; item = item.Next() {
			if obj, ok := item.Value.(connectionListObject); ok {
				var goawayPeaceful = socket.UserID == obj.Conn.UserID
				if !goawayPeaceful {
					obj.Conn.DestroyConnection()
					fmt.Println(socket.UserID, "new device get conn")
					break
				}
			}
		}
	}
	// 用户新连接，更新设备ID
	connList.PushFront(connObj)
	u.connections.Store(socket.UserID, connList)
	go socket.StartHealthCheck()
}

//remove 通过文件描述符移除用户连接
func (u *UserConnection) remove(conn *Socket, fd int) *list.List {
	if conn == nil {
		return nil
	}
	var connList = u.GetConnections(conn.UserID)
	if connList == nil {
		return nil
	}
	for item := connList.Front(); item != nil; item = item.Next() {
		if obj, ok := item.Value.(connectionListObject); ok {
			if obj.UserID == conn.UserID {
				connList.Remove(item)
				// 用户断开连接，更新登录状态
				break
			}
		}
	}
	return connList
}

func (u *UserConnection) GetConnections(userID string) *list.List {
	var targetList *list.List
	if v, ok := u.connections.Load(userID); ok {
		if l, ok := v.(*list.List); ok {
			targetList = l
		}
	}
	return targetList
}
