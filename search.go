package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
	// 	"text/scanner"
)

func searchString(file string, text string, v bool) (err error) {

	f, err := os.Open(file)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)

	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			fmt.Println("Found in entry", file, "on line", line)
			if v == true {
				coloredText := strings.Replace(scanner.Text(), text, "%[1]s", -1)
				blue := color.New(color.Bold, color.FgBlue).SprintFunc()
				fmt.Printf(coloredText, blue(text))
				fmt.Println("\n")
			}
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		// Handle the error
	}

	return err
}

func search(location string, text string, v bool) (err error) {
	fileList := []string{}
	err = filepath.Walk(location, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	for _, file := range fileList {
		if strings.Contains(file, ".md") == true {
			searchString(file, text, v)
		}
	}

	return err
}
