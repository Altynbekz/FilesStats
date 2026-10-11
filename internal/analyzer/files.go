package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func fileName(file string) string {
	name := strings.Split(file, "/")
	return name[len(name)-1]
}

func readDirectory(path string) ([]string, error) {
	var content []string
	err := filepath.WalkDir(path, func(currentPath string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.Type().IsRegular() && filepath.Ext(entry.Name()) == ".txt" {
			content = append(content, currentPath)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not walk directory: %w", err)
	}

	return content, nil
}

func readFile(input []byte) []string {
	return tokenize(string(input))
}
