package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestBetterPrintFiles1(t *testing.T) {
	files, err := os.ReadDir("C:\\WorkSpace\\_home\\")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		if file.Name() != "_home" && file.IsDir() {
			fmt.Print(file.Name() + "    ")
		}
	}
}
func TestBetterPrintFiles(t *testing.T) {
	files, err := os.ReadDir("C:\\WorkSpace\\_home\\")
	if err != nil {
		fmt.Println(err)
	}
	all := make([]string, len(files))
	index := 0
	maxLen := 0
	for _, file := range files {
		if file.Name() != "_home" && file.IsDir() {
			all[index] = file.Name()
			if ll := len(file.Name()); ll > maxLen {
				maxLen = ll
			}
			index++
		}
	}
	all = all[0:index]

	for _, each := range all {
		fmt.Printf("%-"+strconv.Itoa(maxLen)+"s", each+"\t")
	}
}
