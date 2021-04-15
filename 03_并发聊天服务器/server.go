/**
 * @Author:LJL
 * @Date: 2021/4/15 17:41
 */
package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	C    chan string //用于发送数据
	Name string      //用户名
	Addr string      //地址
}

//保存在线用户，key:cliAddr
var onlineMap map[string]Client

//存储message的channel
var message = make(chan string)

func WriteMsgToClient(cli Client, conn net.Conn) {
	for msg := range cli.C {
		conn.Write([]byte(msg + "\n"))
	}
}

func SendMessage() {
	//给map分配空间
	onlineMap = make(map[string]Client)
	for {
		//没有消息时会阻塞
		msg := <-message
		//遍历onlineMap，给每个成员发送message
		for _, cli := range onlineMap {
			cli.C <- msg
		}
	}
}

func MakeMsg(cli Client, msg string) (buf string) {
	buf = "[" + cli.Addr + "] " + cli.Name + " : " + msg
	return
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	//获取客户端的地址
	cliAddr := conn.RemoteAddr().String()
	//创建client结构体，用户名和结构体都为cliAddr
	cli := Client{make(chan string), cliAddr, cliAddr}
	//将结构体添加到map
	onlineMap[cliAddr] = cli

	//新开一个协程，专门给当前客户端发送信息
	go WriteMsgToClient(cli, conn)

	//广播该节点已在线
	message <- MakeMsg(cli, "login")
	//提示，我是谁
	cli.C <- MakeMsg(cli, "I am here")

	isQuit := make(chan bool)  //是否主动退出
	hasData := make(chan bool) //是否有数据

	//该协程用于接收用户发送过来的数据
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				isQuit <- true
				fmt.Println("conn.Read() err =", err)
				return
			}
			msg := string(buf[:n-1])
			//fmt.Printf("msg ==== %s", msg)
			if len(msg) == 3 && msg == "who" {
				//遍历map，给当前用户发送所有成员
				conn.Write([]byte("user list:\n"))
				for _, tmp := range onlineMap {
					msg := tmp.Addr + ":" + tmp.Name + "\n"
					conn.Write([]byte(msg))
				}
			} else if msg == "" {
				continue
			} else if len(msg) >= 8 && msg[:6] == "rename" {
				//修改用户名
				name := strings.Split(msg, "|")[1]
				cli.Name = name
				onlineMap[cliAddr] = cli
				conn.Write([]byte("rename success\n"))
			} else {
				//转发此内容
				message <- MakeMsg(cli, msg)
			}
			hasData <- true
		}

	}()

	select {
	case <-isQuit:
		//从map移除当前用户
		delete(onlineMap, cliAddr)
		message <- MakeMsg(cli, "login out")
		return
	case <-hasData:
	case <-time.After(60 * time.Second):
		delete(onlineMap, cliAddr)
		message <- MakeMsg(cli, "time out leave out")
		return
	}
}

func main() {
	//监听TCP
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("net.Listen() error", err)
		return
	}
	//main结束后关闭listener
	defer listener.Close()

	//该协程转发消息给onlineMap中的成员
	go SendMessage()

	//循环等待TCP连接
	for {
		//没有连接时阻塞等待
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept() error", err)
			continue
		}
		//处理连接
		go HandleConn(conn)
	}

}
