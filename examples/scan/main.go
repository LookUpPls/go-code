package main

import (
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"
)

func main() {
	// 获取本机的IPv4地址和子网掩码
	localIPs := getLocalIPv4()

	var wg sync.WaitGroup
	for _, localIP := range localIPs {
		// 根据子网掩码计算扫描范围
		startIP, endIP := calculateScanRange(localIP)
		fmt.Printf("Scanning IP range: %s - %s\n", startIP, endIP)

		// 使用goroutine并发扫描
		for ip := startIP; !ip.Equal(endIP); incrementIP(ip) {
			wg.Add(1)
			go func(ip net.IP) {
				defer wg.Done()
				if isReachable(ip.String()) {
					fmt.Printf("Found device: %s\n", ip)
				}
			}(dupIP(ip))
		}
	}

	wg.Wait()
	fmt.Println("Scanning completed.")
}

// 获取本机的IPv4地址和子网掩码
func getLocalIPv4() []net.IP {
	var localIPs []net.IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting local IP addresses:", err)
		return localIPs
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIPs = append(localIPs, ipnet.IP)
			}
		}
	}

	return localIPs
}

// 根据IP地址和子网掩码计算扫描范围
func calculateScanRange(ip net.IP) (net.IP, net.IP) {
	mask := net.CIDRMask(24, 32)
	network := ip.Mask(mask)
	start := dupIP(network)
	start[3] = 1
	end := dupIP(network)
	end[3] = 254
	return start, end
}

func isReachable1(ip string) bool {
	//"ping"是要执行的命令。
	//"-c", "1"表示只发送一次ping请求。
	//"-W", "1"表示等待ping响应的超时时间为1秒。
	//ip是要ping的目标IP地址。
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()

	if err != nil {
		fmt.Println(ip, "is unreachable", err)
		return false
	}
	return true
}

// 检查IP地址是否可达, 这个还是可以的,只是扫描22端口
func isReachable(ip string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", ip+":22", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// 增加IP地址
func incrementIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

// 复制IP地址
func dupIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

//Found device: 192.168.31.1
//Found device: 192.168.31.8
//Found device: 192.168.194.10
//Found device: 192.168.31.161
