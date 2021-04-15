/**
 * @Author:LJL
 * @Date: 2021/4/15 16:16
 */
package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func sendFile(path string, conn net.Conn) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("os.open err", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024*4)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件读取完毕")
			} else {
				fmt.Println("文件读取失败:", err)
			}
			return
		}
		//发送读取的内容
		conn.Write(buf[:n])
	}

}

func main() {
	fmt.Println("请输入要发送的文件路径：")
	var path string
	fmt.Scan(&path)
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("err =", err)
		return
	}

	conn, err1 := net.Dial("tcp", "127.0.0.1:8000")
	if err1 != nil {
		fmt.Println("err1 =", err1)
		return
	}

	defer conn.Close()

	//发送文件名
	_, err3 := conn.Write([]byte(info.Name()))
	if err3 != nil {
		fmt.Println("err3 =", err3)
		return
	}

	//接收对方的回复
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("err =", err)
		return
	}
	if "ok" == string(buf[:n]) {
		//发送文件内容
		sendFile(path, conn)
	}

}
