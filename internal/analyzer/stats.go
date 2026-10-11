package main

import (
	"io"
)

type FileStats struct {
	fileName string
	words    map[string]int
	total    int
}

type FileResult struct {
	fileStat FileStats
	err      error
}

func countWord(r io.Reader) (map[string]int, int, error) {
	file, err := io.ReadAll(r)

	if err != nil {
		return make(map[string]int), 0, err
	}

	text := readFile(file)
	fileStats := map[string]int{}
	total := 0

	for _, word := range text {
		fileStats[word]++
		total++
	}

	return fileStats, total, nil

}
