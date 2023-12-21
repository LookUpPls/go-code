package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	testProcess(t, 50)
}
func TestProcess1(t *testing.T) {
	testProcess(t, 500)
}
func TestProcess2(t *testing.T) {
	testProcess(t, 5000)
}
func TestProcess3(t *testing.T) {
	testProcess(t, 50000)
}
func TestProcess4(t *testing.T) {
	testProcess(t, 500000)
}

func testProcess(t *testing.T, n int) {
	rand.Seed(time.Now().UnixNano())
	// 生成连续的数字序列
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

	box := NewMessageBox()
	PrintMemUsage()
	for i, number := range numbers {
		pop := box.put(number, number)

		if false {
			fmt.Printf("put %d", number)
			if pop != nil {
				fmt.Printf("      max %d ", pop.info)
				if pop.next != nil {
					fmt.Printf("> %d ", pop.next.info)
				}
			}

			fmt.Println()
			if false {
				fmt.Print("head ")
				for _, node := range box.headMap {
					PrintList(node)
				}
				fmt.Print("\ntail ")
				for _, node := range box.tailMap {
					PrintList(node)
				}
				fmt.Println()
			}
			fmt.Println()
		}
		if i == len(numbers) {
			assert.Equal(t, strconv.Itoa(n), pop.next.info)
			assert.Equal(t, 0, len(box.headMap))
			assert.Equal(t, 0, len(box.tailMap))
		}
	}
	PrintMemUsage()
	runtime.GC()
	PrintMemUsage()

	//todo: 有可能会添加重复
	box.put(1, 1)
}
