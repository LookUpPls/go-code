package subSeq

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestBatchConfirm1(t *testing.T) {
	numbers := []int{2, 3, 1, 4, 9, 6, 8, 10, 7, 11}
	assert.Equal(t, "1-4 6-11", batchConfirmInner(t, numbers))
}
func TestBatchConfirm2(t *testing.T) {
	numbers := []int{2, 3, 1, 4, 9, 8, 10, 6, 11, 99, 1000, 1001, 1003}
	assert.Equal(t, "8-11 99 1000-1001 1003 6 1-4", batchConfirmInner(t, numbers))
}
func TestBatchConfirm(t *testing.T) {
	numbers := []int{2, 3, 1, 4, 5, 9, 6, 8, 10, 7, 11}
	assert.Equal(t, "1-11", batchConfirmInner(t, numbers))
}
func batchConfirmInner(t *testing.T, numbers []int) string {

	box := NewSubSeq(6, false)
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
	}
	PrintMemUsage()
	runtime.GC()
	PrintMemUsage()

	return box.getAllSubSeq()
	//todo: 有可能会添加重复
	//box.put(1, 1)
}
