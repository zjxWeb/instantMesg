package main

import (
	"net"
	"fmt"
)

type Server struct {
	Ip string
	Port int
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
	}
	return server
}
func (this *Server) Hander(conn net.Conn){
	// ...当前链接业务
	fmt.Println("连接建立成功")
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
