package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestBetterPrin111tFiles1(t *testing.T) {

	numbers := []int{1, 5, 3, 4, 2, 8, 7, 6}

	// 输出结果
	fmt.Println(numbers)

	s := succession{}
	s.init(5)
	for _, number := range numbers {
		fmt.Printf("put %d\n", number)
		s.put(number, int64(number))
	}

	fmt.Printf("%d  %d \n", s.bp, s.len)
	fmt.Printf("%v", s.body.buf)

	assert.Equal(t, 0, s.len)
	assert.Equal(t, true, s.body.isEmpty)
}
func TestBetterPrin111tFi111111les1(t *testing.T) {

	numbers := []int{94, 78, 6, 67, 80, 47, 27, 32, 3, 69, 88, 75, 46, 8, 56, 13, 17, 79, 85, 50, 97, 30, 21, 83, 89, 65, 41, 2, 73, 22, 37, 91, 15, 40, 54, 23, 48, 42, 60, 90, 16, 93, 64, 11, 58, 35, 61, 38, 5, 25, 57, 82, 18, 28, 96, 34, 4, 9, 43, 77, 14, 68, 12, 92, 39, 24, 81, 63, 99, 49, 29, 87, 98, 7, 26, 71, 19, 20, 36, 72, 66, 70, 51, 76, 53, 74, 44, 59, 52, 1, 95, 55, 86, 31, 100, 33, 84, 45, 62, 10}

	// 输出结果
	fmt.Println(numbers)

	s := succession{}
	s.init(5)
	for _, number := range numbers {
		fmt.Printf("%d ", number)
		s.put(number, int64(number))
	}
	fmt.Printf("\n%d  %d", s.bp, s.len)
	fmt.Printf("%v", s.body.buf)
	assert.Equal(t, 0, s.len)
	//assert.Equal(t, len(numbers)-1, s.bp)
	assert.Equal(t, true, s.body.isEmpty)
}

func TestBetterPrintFiles1(t *testing.T) {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成连续的数字序列
	n := 10000 // 比如生成1到10的数字
	numbers := make([]int, n)
	for i := range numbers {
		numbers[i] = i + 1
	}

	// 洗牌
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	// 输出结果
	fmt.Println(numbers)

	s := succession{}
	s.init(5)
	for _, number := range numbers {
		fmt.Printf("%d ", number)
		s.put(number, int64(number))
	}
	fmt.Printf("\nsuccess  bp:%d  len:%d bufLen:%d\n", s.bp, s.len, s.body.size)
	fmt.Printf("%v\n\n", s.body.buf)
	assert.Equal(t, 0, s.len)
	//assert.Equal(t, len(numbers)-1, s.bp)
	assert.Equal(t, true, s.body.isEmpty)

}
func TestBetterPrintFiles(t *testing.T) {
	r := NewRingSlice(30)
	diff := 0
	r.WriteSingleDiff(diff, 1)
	diff++
	r.WriteSingleDiff(diff, 2)
	diff += 2
	r.WriteSingleDiff(diff, 4)
	fmt.Println(r.buf)
}
func TestBetterPrintFil1es1(t *testing.T) {
	r := NewRingSlice(5)
	r.WriteSingle(11)
	r.WriteSingle(12)
	r.WriteSingle(13)
	r.WriteSingle(14)
	r.WriteSingle(15)
	fmt.Println(r.buf)
	r.Read(make([]int64, 4))
	fmt.Println(r.buf)
	r.WriteSingleDiff(0, 1)
	fmt.Println(r.buf)
	r.WriteSingleDiff(0, 2)
	fmt.Println(r.buf)
	r.WriteSingleDiff(1, 4)
	fmt.Println(r.buf)
}
func TestBetterPrintFil111es1(t *testing.T) {
	r := NewRingSlice(5)
	r.WriteSingle(11)
	r.WriteSingle(12)
	r.WriteSingle(13)
	r.WriteSingle(14)
	r.WriteSingle(15)
	fmt.Println(r.buf)
	r.Read(make([]int64, 4))
	fmt.Println(r.buf)
	r.WriteSingleDiff(0, 1)
	fmt.Println(r.buf)
	r.WriteSingleDiff(0, 2)
	fmt.Println(r.buf)
	r.WriteSingleDiff(1, 4)
	fmt.Println(r.buf)
	r.WriteSingleDiff(1, 6)
	fmt.Println(r.buf)

	diff, _ := r.ReadOneDiff(0)
	assert.Equal(t, int64(15), diff, "read 0 error")
	diff, _ = r.ReadOneDiff(1)
	assert.Equal(t, int64(1), diff, "read 0 error")
	diff, _ = r.ReadOneDiff(2)
	assert.Equal(t, int64(2), diff, "read 0 error")
	diff, _ = r.ReadOneDiff(4)
	assert.Equal(t, int64(4), diff, "read 0 error")
	diff, _ = r.ReadOneDiff(6)
	assert.Equal(t, int64(6), diff, "read 0 error")
}

func TestBetterPri1111ntFil111es1(t *testing.T) {
	r := NewRingSlice(5)
	diff := 0
	r.callback = func(d int) {
		diff += -d
	}
	r.WriteSingle(11)
	r.WriteSingle(12)
	r.WriteSingle(13)
	r.WriteSingle(14)
	r.WriteSingle(15)
	fmt.Println(r.buf)
	r.Read(make([]int64, 4))
	fmt.Println(r.buf)
	r.WriteSingle(16)
	r.WriteSingle(17)
	fmt.Println(r.buf)

	result, _ := r.ReadOneDiff(0 - diff)
	assert.Equal(t, int64(15), result, "read 0 error")
	result, _ = r.ReadOneDiff(1 - diff)
	assert.Equal(t, int64(16), result, "read 0 error")
	result, _ = r.ReadOneDiff(2 - diff)
	assert.Equal(t, int64(17), result, "read 0 error")
	result, _ = r.ReadOneDiff(4 - diff)
	assert.Equal(t, int64(0), result, "read 0 error")
}

func TestBetterPrii1111ntFil111es1(t *testing.T) {
	r := NewRingSlice(5)
	diff := 0
	r.callback = func(d int) {
		diff += -d
	}
	r.WriteSingle(11)
	r.WriteSingle(12)
	r.WriteSingle(13)
	r.WriteSingle(14)
	r.WriteSingle(15)
	fmt.Println(r.buf)
	r.Read(make([]int64, 4))
	fmt.Println(r.buf)
	r.WriteSingle(16)
	r.WriteSingle(17)
	fmt.Println(r.buf)

	result, _ := r.ReadOneDiff(0 - diff)
	assert.Equal(t, int64(15), result, "read 0 error")
	result, _ = r.ReadOneDiff(1 - diff)
	assert.Equal(t, int64(16), result, "read 0 error")
	result, _ = r.ReadOneDiff(2 - diff)
	assert.Equal(t, int64(17), result, "read 0 error")
	result, _ = r.ReadOneDiff(4 - diff)
	assert.Equal(t, int64(0), result, "read 0 error")
}
