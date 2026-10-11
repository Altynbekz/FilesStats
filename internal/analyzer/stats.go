package main

import (
	"io"
	"sync"
)

type wordStats struct {
	mx    sync.Mutex
	words map[string]int
	total int
}

func wordCount(
	counter *wordStats,
	jobs <-chan string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for job := range jobs {
		counter.mx.Lock()
		counter.words[job]++
		counter.total++
		counter.mx.Unlock()
	}
}

func statistics(r io.Reader) (*wordStats, error) {
	data, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	text := tokenize(string(data))

	workers := 4
	jobs := make(chan string, workers*4)
	wg := &sync.WaitGroup{}

	counter := &wordStats{
		words: make(map[string]int),
	}

	for i := 1; i < workers; i++ {
		wg.Add(1)
		go wordCount(counter, jobs, wg)
	}

	for _, word := range text {
		jobs <- word
	}

	close(jobs)

	wg.Wait()
	return counter, nil

}
