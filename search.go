//   Copyright 2016 The diarytxt Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	// 	"text/scanner"
)

func searchString(file string, text string, tag string, v bool) (err error) {

	f, err := os.Open(file)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)

	tags := getTags(file)

	line := 1
	c := 0

	if len(tag) == 0 || stringInSlice(tag, tags) == true {
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), text) == true && c == 0 {
				color.Green("Search criteria present in journal entry " + file)
				color.Cyan("on line " + strconv.Itoa(line))
				if v == true {
					coloredText := strings.Replace(scanner.Text(), text, "%[1]s", -1)
					blue := color.New(color.Bold, color.FgBlue).SprintFunc()
					fmt.Printf(coloredText, blue(text))
					fmt.Print("\n\n")
					// fmt.Println("\n")
				}
				c++
			} else if strings.Contains(scanner.Text(), text) == true && c > 0 {
				color.Cyan("on line " + strconv.Itoa(line))
				if v == true {
					coloredText := strings.Replace(scanner.Text(), text, "%[1]s", -1)
					blue := color.New(color.Bold, color.FgBlue).SprintFunc()
					fmt.Printf(coloredText, blue(text))
					fmt.Print("\n\n")
					// fmt.Println("\n")
				}

			}
			line++
		}
	}
	// if err := scanner.Err(); err != nil {
	// 	// Handle the error
	// }

	return err
}

func getTags(file string) []string {

	f, _ := os.Open(file)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "tags:") == true {
			r := strings.Split(strings.Split(scanner.Text(), ":")[1], ",")
			var tags []string
			for _, str := range r {
				tags = append(tags, strings.TrimSpace(str))
			}
			return tags
		}
	}
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func search(location string, text string, tag string, v bool) (err error) {
	fileList := []string{}
	err = filepath.Walk(location, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	for _, file := range fileList {
		if strings.Contains(file, ".md") == true {
			searchString(file, text, tag, v)
		}
	}

	return err
}
