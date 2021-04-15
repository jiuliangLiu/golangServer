/**
 * @Author:LJL
 * @Date: 2021/4/15 14:31
 */
package main

import (
	"fmt"
	"net"
	"strings"
)

func HandleConn(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	fmt.Println(addr, "connect successful")

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err =", err)
			return
		}
		fmt.Println("len(buf[:n]) =", len(buf[:n]))
		if "exit" == string(buf[:n-2]) {
			fmt.Println(addr, " exit")
			return
		}
		fmt.Println(addr, ": buf =", string(buf[:n]))
		conn.Write([]byte(strings.ToUpper(string(buf[:n]))))
	}

}

func main() {
	listner, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("err =", err)
		return
	}
	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("err =", err)
			return
		}
		//处理用户请求
		go HandleConn(conn)
	}
}
