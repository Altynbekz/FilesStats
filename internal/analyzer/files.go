package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func checkExtension(
	jobs <-chan os.DirEntry,
	results chan<- string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for file := range jobs {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".txt" {
			results <- file.Name()
		}
	}
}

func readDirectory(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("The directory does not exist")
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("The directory does not have any content in it")
	}

	workers := 4
	wg := &sync.WaitGroup{}
	jobs := make(chan os.DirEntry, workers*2)
	result := make(chan string, workers*2)

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go checkExtension(jobs, result, wg)
	}
	go func() {
		for _, file := range files {
			jobs <- file
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(result)
	}()

	var content []string

	for filename := range result {
		content = append(content, filename)
	}

	return content, nil
}

func readFile(path string) (string, error) {
	file, err := os.ReadFile("../test/" + path)
	if err != nil {
		return "", fmt.Errorf("The error is %w", err)
	}

	return string(file), nil
}
