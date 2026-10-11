package main

import (
	"fmt"
	"os"
)

func main() {
	data, _ := os.Open("../test/GOGO.txt")

	stats, _ := statistics(data)

	fmt.Println(stats)

}
