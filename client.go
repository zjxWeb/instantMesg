package main

import (
	"net"
	"fmt"
	"flag"
)

type Client struct{
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int // 当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client{
	// 创建客户端对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 999,
	}

	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err!= nil {
		fmt.Println("连接服务器失败",err)
		return nil
	}
	client.conn = conn
	// 返回对象
	return client
}

func (client *Client) menu() bool{
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	}else{
		fmt.Println("输入错误，请重新输入")
		return false
	}
}

func(client *Client) Run(){
	for client.flag != 0{
		for client.menu() != true {
		}
		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			// 公聊天模式
			fmt.Println("公聊天模式")
			break
		case 2:
			// 私聊模式
			fmt.Println("私聊模式")
			break
		case 3:
			// 群聊模式
			fmt.Println("群聊模式")
			break
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1
func init(){
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "默认server ip 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8080, "默认server port 8080")
}
func main() {
	// 命令行解析
	flag.Parse()
	fmt.Println("ip:",serverIp)
	client := NewClient(serverIp, serverPort)
	if client == nil{
		fmt.Println("连接服务器失败")
		return
	}
	fmt.Println("连接服务器成功")

	// 启动客户端的业务
	client.Run()
}