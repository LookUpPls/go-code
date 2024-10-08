package subSeq

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestBatchMaxInner(t *testing.T) {
	testBatchMaxInner(t, 50, true)
}
func TestBatchMaxInner1(t *testing.T) {
	testBatchMaxInner(t, 500, true)
}
func TestBatchMaxInner2(t *testing.T) {
	testBatchMaxInner(t, 5000, true)
}
func TestBatchMaxInner3(t *testing.T) {
	testBatchMaxInner(t, 50000, true)
}
func TestBatchMaxInner4(t *testing.T) {
	testBatchMaxInner(t, 500000, true)
}

func testBatchMaxInner(t *testing.T, n int, isPop bool) {
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

	box := NewSubSeq(6, isPop)
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
		}
		if false || i == len(numbers)-1 {
			fmt.Print("head ")
			for _, node := range box.headMap {
				PrintList(node)
			}
			fmt.Print("\ntail ")
			for _, node := range box.tailMap {
				PrintList(node)
			}
			fmt.Println()
			fmt.Println()
		}
		if isPop && i == len(numbers)-1 {
			assert.Equal(t, n, pop.next.info)
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
