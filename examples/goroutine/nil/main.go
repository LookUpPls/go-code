package main

import (
	"fmt"
	"strconv"
)

func main() {
	a := make(chan int)
	b := make(chan int)
	// 这个是独立在for之外的， 不管超不超时都是1s
	go func() {
		for i := 0; i < 200; i++ {
			fmt.Println("i" + strconv.Itoa(i))
			select {
			case a <- 1:
				//a = nil
			case b <- 2:
				//b = nil
			}
		}
	}()
	fmt.Println(<-a)
	fmt.Println(<-b)
}
