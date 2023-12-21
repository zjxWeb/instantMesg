package main

import (
	"net"
	"fmt"
	"sync"
	"io"
	"time"
)

type Server struct {
	Ip string
	Port int
	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex
	
	// 消息广播的channel
	Message chan string
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessage(){
	for{
		msg := <-this.Message

		// 将msg发送给全部在线的user
		this.mapLock.RLock()
		for _,cli := range this.OnlineMap{
			cli.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

// 广播消息方法
func (this *Server) BroadCast(user *User, msg string){
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Hander(conn net.Conn){
	// ...当前链接业务
	// fmt.Println("连接建立成功")

	user := NewUser(conn,this)

	user.Offline()
	// // 用户上线了，将用户加入到onlineMap中
	// this.mapLock.Lock()
	// this.OnlineMap[user.Name] = user
	// this.mapLock.Unlock()
	// // 广播当前用户上线消息
	// this.BroadCast(user,"已上线")

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func(){
		buf := make([]byte, 4096)
		for{
			n, err := conn.Read(buf)
			if n == 0 {
				// this.BroadCast(user,"已下线")
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read error:", err)
				return
			}
			// 提取用户的消息，（去除 ‘\n’）
			msg := string(buf[:n-1])
			
			// 将得到的消息进行广播
			// this.BroadCast(user, msg)

			//用户针对msg进行消息处理
			user.DoMessage(msg)

			// 用户的任意消息，代表用户是一个活跃的
			isLive <- true
		}
	}()

	// 当前handle阻塞 
	for{
		select {
			case <- isLive:
				// 当前用户是或与的，应该重置定时器
				// 不做任何操作，为了激活select，更新定时器
			case <- time.After(time.Second * 300):
				// 已经超时
				// 将当前的User强制关闭
				user.SendMsg("你已经超时了，请重新登录")

				// 销毁用的资源
				close(user.C)

				// 关闭连接
				conn.Close()

				// 退出当前的Handleer
				return //runtime.Goexit()
		}
	}
}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err!= nil {
		fmt.Println("listen error:",err)
		return
	}
	// close listen socket
	defer listen.Close()
	//启动监听Message的goroutine
	go this.ListenMessage()

	for{
		// accept
		conn, err := listen.Accept()
		if err!= nil {
			fmt.Println("accept error:",err)
			continue
		}
		// do hander
		go this.Hander(conn)
	}

}
