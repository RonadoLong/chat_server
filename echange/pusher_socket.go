package echange

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Socket struct {
	net.Conn          `json:"-"`  //连接
	FD                int         `json:"fd"`                  //文件描述符
	UserID            string      `json:"user_id"`             //用户ID
	DeviceID          string      `json:"device_id,omitempty"` //用户设备ID，仅移动端有效
	Platform          string      `json:"platform"`            //连接所属平台
	LastHeartbeatTime time.Time   `json:"last_heartbeat_time"`
	heartBeatChan     chan []byte //心跳包通道
	lock              sync.Mutex
}

//NewSocket 新建socket对象
func NewSocket(userID string, conn net.Conn) *Socket {
	return &Socket{
		Conn:          conn,
		UserID:        userID,
		heartBeatChan: make(chan []byte, 10),
	}
}

// offline
func (conn *Socket) DestroyConnection() {
	if conn == nil || conn.Conn == nil {
		return
	}
	conn.lock.Lock()
	defer conn.lock.Unlock()
	conn.FD = -1
	// wsutil.WriteServerMessage(conn.Conn, ws.OpText, []byte("goaway"))
	if _, err := conn.Conn.Write(ws.CompiledCloseNormalClosure); err != nil {
		fmt.Println(err)
		return
	}
	conn.Conn.Close()
	_ = epoll.Remove(conn)
	conn.Conn = nil
}

func (conn *Socket) StartHealthCheck() {
	log.Println("===========================3333")
	go conn.ReadConnectionData()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		if conn == nil || conn.Conn == nil {
			return
		}
		select {
		case b := <-conn.heartBeatChan:
			if b == nil {
				break
			}
			conn.LastHeartbeatTime = time.Now()
			//fmt.Println("新websocket消息：", conn.UserID, conn.Platform, string(b))
		case <-time.After(time.Second * 15):
			log.Println("读取心跳包超时（15秒），主动断开连接")
			conn.DestroyConnection()
		}
	}
}

//ReadConnectionData 读取数据
func (conn *Socket) ReadConnectionData() {
	log.Println("read data")
	defer func() {
		if err := recover(); err != nil {
			// 	Msg("读取心跳包发生异常")
			fmt.Println("read ping err ", err)
		}
	}()
	for {
		if conn == nil || conn.Conn == nil {
			return
		}
		if msg, _, err := wsutil.ReadClientData(conn.Conn); err != nil {
			conn.DestroyConnection()
			return
		} else {
			if conn == nil || conn.Conn == nil {
				return
			}
			if len(msg) != 0 && string(msg) == "q" {
				//客户端主动退出
				fmt.Println("client quit conn")
				conn.DestroyConnection()
				return
			}
			wsutil.WriteServerMessage(conn.Conn, ws.OpText, []byte("0"))
			conn.heartBeatChan <- msg
		}
	}
}
