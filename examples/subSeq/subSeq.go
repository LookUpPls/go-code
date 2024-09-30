package subSeq

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// PrintList prints the entire linked list starting from the given node.
func PrintList(head *Node) {
	var sb strings.Builder
	sb.WriteString(" | ")

	for current := head; current != nil; current = current.next {
		sb.WriteString(strconv.Itoa(current.info))
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
	info int
	next *Node
	prev *Node
}

// String returns a string representation of the Node.
func (n *Node) String() string {
	return fmt.Sprintf("%d", n.info)
}

type SubSeq struct {
	headMap   map[int]*Node
	tailMap   map[int]*Node
	waitPoint int
	isPop     bool
}

func NewSubSeq(start int, isPop bool) *SubSeq {
	/**
	 * 已知一个消息流会不断吐出整数1-N，
	 * 但不一定按照顺序依次吐出
	 * 如果上次打印的序号为i，那么当i+1出现时
	 * 请打i+1及其之后接受过的并且连续的所有数
	 * 直到1-N全部接收并打印完
	 * 请设计这种接收并打印的结构
	 */
	ans := &SubSeq{}
	ans.headMap = make(map[int]*Node)
	ans.tailMap = make(map[int]*Node)
	ans.waitPoint = start
	ans.isPop = isPop
	return ans
}

func (this *SubSeq) put(num int, info int) *Node {
	if this.isPop && num < this.waitPoint {
		return nil
	}
	cur := &Node{info: info}
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
			// cur 删掉
			temp.prev = cur.prev
			cur.prev.next = temp
			//cur.prev = nil
			//cur.next = nil
		} else {
			cur.next = temp
			temp.prev = cur
		}
		//cur.next = this.headMap[num+1]
		delete(this.headMap, num+1)
	} else {
		this.tailMap[num] = cur
	}

	if this.isPop {
		if num == this.waitPoint {
			return this.pop()
		}
	}
	return nil
}

func (this *SubSeq) pop() *Node {
	node := this.headMap[this.waitPoint]
	delete(this.headMap, this.waitPoint)
	if node.next != nil {
		//fmt.Printf("%d ", node.next.info)
		this.waitPoint = node.next.info + 1
	} else {
		this.waitPoint++
	}
	delete(this.tailMap, this.waitPoint-1)
	return node
}

func (this *SubSeq) getAllSubSeq() string {
	var sb strings.Builder
	for _, node := range this.headMap {
		for current := node; current != nil; current = current.next {
			sb.WriteString(strconv.Itoa(current.info))
			if current.next != nil {
				sb.WriteString("-")
			}
		}
		sb.WriteString(" ")
	}
	return strings.TrimSuffix(sb.String(), " ")
}
