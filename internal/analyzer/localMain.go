package main

import "fmt"

func main() {
	files, _ := readDirectory("../test")
	for _, file := range files {
		fmt.Println(file)
	}
}
