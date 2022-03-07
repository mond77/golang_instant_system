package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"src/GolangStudy/golang_instant_system/GolangStudy/golang_instant_system/redis"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient() *Client {
	client := &Client{
		flag: 99,
	}
	return client

}
func (client *Client) Dial(serverIp string, serverPort int) {
	client.ServerIp = serverIp
	client.ServerPort = serverPort
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)

	}
	client.conn = conn

	client.ConnSetName()

}

func (client *Client) DealResponse() {
	//永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 2 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>请输入合法数字<<<")
		return false

	}

}

func (client *Client) ConnSetName() {

	sendmsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendmsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)

	}
}
func (client *Client) PublicChat() {
	var chatmsg string
	fmt.Println(">>>请输入聊天内容,exit退出")
	fmt.Scanln(&chatmsg)

	for chatmsg != "exit" {
		if len(chatmsg) != 0 {
			sendmsg := chatmsg + "\n"
			_, err := client.conn.Write([]byte(sendmsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break

			}
		}
		chatmsg = ""
		fmt.Println(">>>公聊模式,exit退出")
		fmt.Scanln(&chatmsg)

	}

}
func (client *Client) SelectUsers() {
	sendmsg := "who\n"
	_, err := client.conn.Write([]byte(sendmsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return

	}
}
func (client *Client) PrivateChat() {
	var remotename string
	var chatmsg string
	client.SelectUsers()
	fmt.Println(">>>请输入聊天对象,exit退出")
	fmt.Scanln(&remotename)

	for remotename != "exit" {
		fmt.Println("请输入消息内容,exit退出:")
		fmt.Scanln(&chatmsg)

		for chatmsg != "exit" {
			if len(chatmsg) != 0 {
				sendmsg := "to|" + remotename + "|" + chatmsg + "\n\n"
				_, err := client.conn.Write([]byte(sendmsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break

				}
			}
			chatmsg = ""
			fmt.Println(">>>请输入消息内容,exit退出")
			fmt.Scanln(&chatmsg)
		}
		client.SelectUsers()
		fmt.Println(">>>请输入聊天对象,exit退出")
		fmt.Scanln(&remotename)
	}

}

func (client *Client) verify() bool {
	fmt.Println("请输入" + client.Name + "的密码")
	var psw string
	for {
		fmt.Scanln(&psw)
		if psw == "exit" {
			return false
		}
		ok := redis.VerifyUser(client.Name, psw)
		if ok {
			break
		}
	}
	return true

}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			client.PublicChat()

			break
		case 2:
			client.PrivateChat()

			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口")
}
func (client *Client) SetName() {
	fmt.Println(">>>请输入用户名:")
	fmt.Scanln(&client.Name)
}
func main() {
	//命令行解析

	flag.Parse()
	client := NewClient()
	client.SetName()
	ok := client.verify()
	if !ok {
		fmt.Println("verify failure")
		return
	}
	client.Dial(serverIp, serverPort)
	//处理server回执消息
	go client.DealResponse()
	fmt.Println(">>>链接服务器成功...")
	//启动客户端业务
	//验证用户是否注册
	client.Run()

}
