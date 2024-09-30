package main

import (
	"fmt"
)

func parseData(data []byte) {
	if len(data) < 35 {
		fmt.Println("Invalid data length")
		return
	}

	if data[0] != 0x68 {
		fmt.Println("Invalid header")
		return
	}

	length := int(data[1]) + int(data[2])<<8
	if len(data) != length+3 {
		fmt.Println("Data length mismatch")
		return
	}

	targetAddr := data[3]
	sourceAddr := data[4]

	fmt.Printf("Target address: 0x%02X\n", targetAddr)
	fmt.Printf("Source address: 0x%02X\n", sourceAddr)

	// Parse YX
	parseYX(data[5:11])

	// Parse YC
	parseYC(data[11:35])

	//powerAdjType := data[35]
	//powerAdjParam := int(data[36]) + int(data[37])<<8
	//fmt.Printf("Current power adjustment command type: 0x%02X\n", powerAdjType)
	//fmt.Printf("Current power adjustment parameter: %d\n", powerAdjParam)
}

func parseYX(data []byte) {
	yx0 := data[0]
	workStatusInt := yx0 & 0x03
	var workStatus string
	switch workStatusInt {
	case 0:
		workStatus = "待机"
	case 1:
		workStatus = "工作"
	case 2:
		workStatus = "充电完成"
	case 3:
		workStatus = "充电暂停"
	}
	fmt.Printf("工作状态: %s\n", workStatus)
	//totalFault := (yx0 >> 2) & 0x01
	//fmt.Printf("Total fault: 0x%02X\n", totalFault)
	//totalAlarm := (yx0 >> 3) & 0x01
	//fmt.Printf("Total alarm: 0x%02X\n", totalAlarm)

	connectStatus := (data[4]) & 0x10
	if connectStatus == 0 {
		fmt.Printf("车辆连接状态: 已连接\n")
	} else if connectStatus == 1 {
		fmt.Printf("车辆连接状态: 未连接\n")
	}
}

func parseYC(data []byte) {
	outputVoltage := int(data[0]) + int(data[1])<<8
	fmt.Printf("充电输出电压: %.1f V\n", float64(outputVoltage)/10)

	outputCurrent := int(data[2]) + int(data[3])<<8
	fmt.Printf("充电输出电流: %.1f A\n", float64(outputCurrent)/10)

	outputStatus := int(data[22]) + int(data[23])<<8
	switch outputStatus {
	case 0:
		fmt.Printf("充放电状态: 待机\n")
	case 1:
		fmt.Printf("充放电状态: 充电\n")
	case 2:
		fmt.Printf("充放电状态: 放电\n")
	default:
		fmt.Printf("充放电状态: 未知\n")
	}

	soc := int(data[4]) + int(data[5])<<8
	fmt.Printf("SOC: %d%%\n", soc)
}

func main() {
	// Example data from the document
	data := []byte{0x68, 0x23, 0x00, 0xdc, 0x67, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0x00, 0x00, 0xa0, 0x0f, 0x00, 0x00, 0x32, 0x00, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa5, 0x04, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	//Data: 682300dc6700000000002c0000a00f000032003200000000008a01000000000400e001000000
	parseData(data)
}
