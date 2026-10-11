package main

import "sort"

type wordCount struct {
	word  string
	count int
}

func sortTop(top int, words map[string]int) map[string]int {
	capacity := words
	result := make([]wordCount, 0, len(capacity))

	for key, value := range capacity {
		result = append(result, wordCount{
			word:  key,
			count: value,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].count == result[j].count {
			return result[i].word < result[j].word
		}
		return result[i].count > result[j].count
	})

	newMap := map[string]int{}

	if top > len(result) {
		top = len(result)
	}

	for _, item := range result[:top] {
		newMap[item.word] = item.count
	}
	return newMap
}
