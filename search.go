package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	// 	"text/scanner"
)

func searchString(file string, text string) (err error) {

	f, err := os.Open(file)
	if err != nil {
		return err
	}

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			fmt.Println("Found in entry", file, "on line", line)
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		// Handle the error
	}

	return err
}

func walkFiles(location string, text string) (err error) {
	fileList := []string{}
	err = filepath.Walk(location, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	for _, file := range fileList {
		if strings.Contains(file, ".md") == true {
			searchString(file, text)
		}
	}

	return err
}
