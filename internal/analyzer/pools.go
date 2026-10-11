package main

import (
	"os"
	"sync"
)

func process(path string, top int) (*FileStats, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	name := fileName(path)
	words, total, err := countWord(file)
	words = sortTop(top, words)
	if err != nil {
		return nil, err
	}

	return &FileStats{
		fileName: name,
		words:    words,
		total:    total,
	}, nil
}

func worker(
	jobs <-chan string,
	top int,
	result chan<- FileResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for job := range jobs {
		stats, err := process(job, top)

		if err != nil {
			result <- FileResult{err: err}
			continue
		}

		result <- FileResult{
			fileStat: *stats,
			err:      err,
		}
	}

}

func work(top, workers int, dir string) ([]FileResult, error) {

	wg := &sync.WaitGroup{}
	jobs := make(chan string, workers*2)
	result := make(chan FileResult)
	files, err := readDirectory(dir)

	if err != nil {
		return []FileResult{}, err
	}

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(jobs, top, result, wg)
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

	content := []FileResult{}

	for file := range result {

		content = append(content, file)
	}

	return content, nil
}
