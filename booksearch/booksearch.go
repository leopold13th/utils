package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	files := make([]string, 10)
	searchWords := []string{}
	for _, arg := range args {
		searchWords = append(searchWords, arg)
	}

	err := filepath.Walk("/home/shaman/Documents/GDrive/leopold13th",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			// fmt.Println(info.Size(), path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		var hits int
		for _, word := range searchWords {
			if strings.Contains(file, word) {
				hits++
			}
		}
		if hits >= len(searchWords) {
			fmt.Println(hits, file)
		}
	}
}
