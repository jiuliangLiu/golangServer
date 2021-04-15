/**
 * @Author:LJL
 * @Date: 2021/4/15 16:40
 */
package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func RecvFile(fileName string, conn net.Conn) {
	//新建文件
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("os.Create()失败：", err)
		return
	}

	//接收对方发过来的文件内容
	buf := make([]byte, 4*1024)
	for {
		n, err := conn.Read(buf) //读取发送的文件内容
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件接收完毕")
			} else {
				fmt.Println("conn.Read()失败：", err)
			}
			return
		}
		//将读取的文件内容写入文件
		f.Write(buf[:n])
	}
}

func main() {
	//监听
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("net.Listen失败：", err)
		return
	}

	//阻塞等待
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("listener.Accept()失败：", err)
		return
	}

	defer conn.Close()
	//读取内容
	buf := make([]byte, 1024)
	n, err := conn.Read(buf) //读取发送的文件名
	if err != nil {
		fmt.Println("conn.Read()失败：", err)
		return
	}

	//获取文件名称
	fileName := string(buf[:n])

	//回复ok
	conn.Write([]byte("ok"))

	//接收文件内容
	RecvFile(fileName, conn)
}
