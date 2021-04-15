/**
 * @Author:LJL
 * @Date: 2021/4/15 15:13
 */
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("err =", err)
		return
	}
	defer conn.Close()
	str := make([]byte, 1024)
	go func() {
		for {
			n, err := os.Stdin.Read(str)
			if err != nil {
				fmt.Println("err =", err)
				return
			}
			conn.Write(str[:n])
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err =", err)
			return
		}
		fmt.Printf("buf[:%d]:%s\n", n, string(buf[:n]))
	}
}
