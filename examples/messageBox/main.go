package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// PrintList prints the entire linked list starting from the given node.
func PrintList(head *Node) {
	var sb strings.Builder
	sb.WriteString(" | ")

	for current := head; current != nil; current = current.next {
		sb.WriteString(current.info)
		if current.next != nil {
			sb.WriteString(" > ")
		}
	}

	fmt.Print(sb.String())
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 打印内存使用情况
	fmt.Printf("Alloc = %v KB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v KB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v KB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024
}

type Node struct {
	info string
	next *Node
	prev *Node
}

// String returns a string representation of the Node.
func (n *Node) String() string {
	return fmt.Sprintf("%s", n.info)
}

type MessageBox struct {
	headMap   map[int]*Node
	tailMap   map[int]*Node
	waitPoint int
	min       int
	minNode   *Node
}

func NewMessageBox() *MessageBox {
	ans := &MessageBox{}
	ans.headMap = make(map[int]*Node)
	ans.tailMap = make(map[int]*Node)
	ans.waitPoint = 1
	ans.min = 2<<32 - 1
	return ans
}

// 消息的编号，info消息的内容, 消息一定从1开始
func (this *MessageBox) put(num int, info string) {
	if num < 1 {
		return
	}
	cur := &Node{info: info}
	if num < this.min {
		this.min = num
		this.minNode = cur
	}
	// num~num
	// 建立了num~num这个连续区间的头和尾
	// 查询有没有某个连续区间以num-1结尾
	var temp *Node
	if tail, ok := this.tailMap[num-1]; ok {
		if tail.prev != nil {
			temp = tail.prev
			tail.prev = nil
		} else {
			temp = tail
		}
		cur.prev = temp
		temp.next = cur
		//this.tailMap[num-1].next = cur
		delete(this.tailMap, num-1)
	} else {
		this.headMap[num] = cur
	}
	// 查询有没有某个连续区间以num+1开头的
	if prev, ok := this.headMap[num+1]; ok {
		if prev.next != nil {
			temp = prev.next
			prev.next = nil
		} else {
			temp = prev
		}
		if cur.prev != nil {
			cur.prev.next = temp
			temp.prev = cur.prev
		} else {
			cur.next = temp
			temp.prev = cur
		}
		//cur.next = this.headMap[num+1]
		delete(this.headMap, num+1)
	} else {
		this.tailMap[num] = cur
	}

	if false {
		fmt.Printf("put %d max %s\n", num, this.minNode.next)
		fmt.Printf("")
		fmt.Print("head ")
		for _, node := range this.headMap {
			PrintList(node)
		}
		fmt.Print("\ntail ")
		for _, node := range this.tailMap {
			PrintList(node)
		}
		fmt.Println()
		fmt.Println()
	}
	//if num == this.waitPoint {
	//	this.pop()
	//}
}

func (this *MessageBox) pop() {
	fmt.Print("\nprint ")
	node := this.headMap[this.waitPoint]
	delete(this.headMap, this.waitPoint)
	for node != nil {
		fmt.Print(node.info + " ")
		node = node.next
		this.waitPoint++
	}
	delete(this.tailMap, this.waitPoint-1)
	fmt.Println()
}

func main() {
	/**
	 * 已知一个消息流会不断吐出整数1-N，
	 * 但不一定按照顺序依次吐出
	 * 如果上次打印的序号为i，那么当i+1出现时
	 * 请打i+1及其之后接受过的并且连续的所有数
	 * 直到1-N全部接收并打印完
	 * 请设计这种接收并打印的结构
	 */
	box := NewMessageBox()
	// 1....
	rand.Seed(time.Now().UnixNano())
	// 生成连续的数字序列
	n := 700000 // 比如生成1到10的数字
	numbers := make([]int, n)
	for i := range numbers {
		numbers[i] = i + 1
	}
	// 洗牌
	rand.Shuffle(n, func(i, j int) {
		if rand.Int() > 1000 {
			numbers[i], numbers[j] = numbers[j], numbers[i]
		}
	})
	//numbers = []int{2, 3, 1, 4, 5, 9, 6, 8, 10, 7, 11}

	PrintMemUsage()
	for _, number := range numbers {
		box.put(number, strconv.Itoa(number))
	}
	PrintMemUsage()
	runtime.GC()
	PrintMemUsage()
	box.put(1, "1")
}
