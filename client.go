package main

import (
	"net"
	"fmt"
	"flag"
	"io"
	"os"
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

// 处理server回应的消息，直接现实标准输出即可
func (client *Client) DeslResponse() {
	// 永久阻塞监听  一旦有数据，就直接copy到stdout 标准输出上
	io.Copy(os.Stdout, client.conn)
	// 等价于下面的写法
	// for{
	// 	buf := make()
	// 	client.conn.Read(buf)
	// 	fmt.Println(buf)
	// }

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

//查询当前都有哪些用户在线
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err!= nil {
		fmt.Println("查询用户失败",err)
		return
	}
}

// 私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	client.SelectUsers()
	fmt.Println("请选择用户,exit退出")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println("请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			// 消息部位空则发送
			if len(chatMsg) != 0{
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_,err := client.conn.Write([]byte(sendMsg))
				if err!= nil {
					fmt.Println("发送消息失败",err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请选择用户,exit退出")
			fmt.Scanln(&remoteName)
		}
		client.SelectUsers()
		fmt.Println("请选择用户,exit退出")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string
	fmt.Println("请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)
	for chatMsg!= "exit" {
		// 发送消息给服务器

		// 消息部位空则发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err!= nil {
				fmt.Println("发送消息失败:",err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

// 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println("请输入新的用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename| " + client.Name + "\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err!= nil{
		fmt.Println("更新用户名失败")
		return false
	}
	return true
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
			client.PublicChat()
			break
		case 2:
			// 私聊模式
			fmt.Println("私聊模式")
			client.PrivateChat()
			break
		case 3:
			// 更新用户名
			client.UpdateName()
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

	// 单独卡其一个gorountine去处理server回执的消息
	go client.DeslResponse()

	fmt.Println("连接服务器成功")

	// 启动客户端的业务
	client.Run()
}