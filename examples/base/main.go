package main

import "fmt"

func main() {
	l := "2"
	go checkLink(l)
	fmt.Println(l)
	go checkLink(l)
	fmt.Println(l)
}
func checkLink(link string) {
	link = "1"
	fmt.Println("func " + link)
}
