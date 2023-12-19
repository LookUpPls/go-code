package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	// 这个是独立在for之外的， 不管超不超时都是1s
	tick := time.Tick(time.Second)
	for {
		select {
		case <-c:
			fmt.Println("random 01")
		case <-tick:
			fmt.Println("tick")
		//case <-time.After(3000 * time.Millisecond):
		// 这个每次进for循环会重新计时
		case <-time.After(400 * time.Millisecond):
			fmt.Println("timeout")
			//default:  //有default后 time.after就一次不会触发
			//fmt.Println("default")
		}
	}
}
