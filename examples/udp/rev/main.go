package main

import (
	"fmt"
	"net"
)

func main() {
	// 创建UDP地址结构
	addr, err := net.ResolveUDPAddr("udp", ":4000")
	// addr, err := net.ResolveUDPAddr("udp", "192.168.1.220:4000")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	// 监听UDP端口
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening on UDP port:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server is listening on port 4000...")

	// 接收数据的缓冲区
	buffer := make([]byte, 1024)

	for {
		// 读取数据
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading data from UDP:", err)
			continue
		}

		// 打印接收到的数据
		data := buffer[:n]
		for i := range data {
			fmt.Printf("%d ", data[i])
		}
		fmt.Printf("\n ")
		fmt.Printf("Received data from %s \n", remoteAddr.String())
	}
}
