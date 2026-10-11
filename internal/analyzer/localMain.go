package main

import (
	"fmt"
	"time"
)

func run() error {
	userConfig := loadConfig()

	result, err := work(userConfig.top, userConfig.workers, userConfig.dir)

	if err != nil {
		return fmt.Errorf("work failed: %w", err)
	}

	for _, item := range result {
		fmt.Println("----------------------------------------")
		fmt.Printf("| Name:  %s\n", item.fileStat.fileName)
		fmt.Printf("| Total: %d\n", item.fileStat.total)
		fmt.Println("|---------------------------------------")
		fmt.Println("| Stats")

		for word, count := range item.fileStat.words {
			fmt.Printf("| %-20s | %d\n", word, count)
		}

		fmt.Println("----------------------------------------")
	}

	return nil
}

func main() {
	start := time.Now()
	err := run()
	fmt.Println("Since it started ", time.Since(start))
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
