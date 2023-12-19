package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://www.baidu.com",
		"http://www.google.com",
		"http://www.bing.com",
	}
	c := make(chan string)

	for _, l := range links {
		go checkLink(l, c)
	}
	//for {
	//	go checkLink(<-c, c)
	//}
	for l := range c {
		//go func() { //闭包的副作用会导致l变量被覆盖。
		//	time.Sleep(1 * time.Second)
		//	go checkLink(l, c)
		//}()
		go sleepAndCheck(l, c)
	}
}
func sleepAndCheck(l string, c chan string) {
	time.Sleep(1 * time.Second)
	go checkLink(l, c)
}
func checkLink(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println("error " + link + "  " + err.Error())
		return
	} else {
		fmt.Println("success " + link)
	}
	c <- link
}
