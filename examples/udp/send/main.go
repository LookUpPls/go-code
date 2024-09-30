package main

import (
	"fmt"
	"net"
)

func main() {
	// 创建UDP地址结构
	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.1.101:4000")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	// 连接UDP服务器
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error dialing UDP server:", err)
		return
	}
	defer conn.Close()

	// 发送数据
	for true {
		message := "Hello, UDP server!"
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending data to UDP server:", err)
			return
		}
		fmt.Println("Data sent to UDP server:", message)
	}

}
