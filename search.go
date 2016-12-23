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
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func search(location string, text string, tag string, y string, m string, v bool) (err error) {
	fileList := []string{}
	err = filepath.Walk(location, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	fileList = filterFile(fileList)
	fileList = filterTag(fileList, tag)
	fileList = filterYear(fileList, y)
	fileList = filterMonth(fileList, m)
	fileList = filterText(fileList, text)

	report(fileList, text, tag, y, m, v)

	return err
}

func filterFile(f []string) []string {

	var fo []string

	for _, file := range f {
		if strings.Contains(file, "rendered_diary") == false {
			fo = append(fo, file)
		}
	}

	return fo
}

func filterTag(f []string, t string) []string {

	var fileList []string

	for _, file := range f {
		fo, _ := os.Open(file)

		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "tags:") == true && strings.Contains(scanner.Text(), t) == true {
				fileList = append(fileList, file)
			}
		}
	}

	return fileList
}

func filterYear(f []string, y string) []string {

	var fileList []string

	for _, file := range f {
		fo, _ := os.Open(file)

		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "date:") == true {
				year := strings.Split(strings.Split(scanner.Text(), ":")[1], "-")[0]
				if strings.Contains(year, y) == true {
					fileList = append(fileList, file)
				}
			}
		}
	}

	return fileList
}

func filterMonth(f []string, m string) []string {

	var fileList []string

	for _, file := range f {
		fo, _ := os.Open(file)

		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "date:") == true {
				month := strings.Split(strings.Split(scanner.Text(), ":")[1], "-")[1]
				if strings.Contains(month, m) == true {
					fileList = append(fileList, file)
				}
			}
		}
	}

	return fileList
}

func filterText(f []string, t string) []string {

	var fileList []string

	for _, file := range f {
		fo, _ := os.Open(file)
		c := 0
		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {
			if strings.Contains(strings.ToUpper(scanner.Text()), strings.ToUpper(t)) == true {
				c++
			}
		}
		if c > 0 {
			fileList = append(fileList, file)
		}
	}

	return fileList
}

func report(f []string, text string, tag string, y string, m string, v bool) {
	if len(f) > 0 {
		for _, file := range f {

			var output []string
			var outputVerbose []string
			var caseText []string
			var outputCaseText []([]string)

			color.Green("Search criteria present in journal entry " + file)

			f, _ := os.Open(file)
			fc, _ := ioutil.ReadFile(file)
			// if err != nil {
			// 	return err
			// }

			scanner := bufio.NewScanner(f)

			// tags := getTags(file)

			line := 1
			// c := 0

			for scanner.Scan() {

				// var output []string
				// var outputVerbose []string

				// TAGS

				if len(tag) > 0 {
					if strings.Contains(scanner.Text(), "tags:") == true {
						if strings.Contains(strings.Split(scanner.Text(), ":")[1], tag) == true {
							color.Cyan("Journal entry is tagged with " + tag)
						}
					}
				}

				// YEAR

				if len(y) > 0 {
					if strings.Contains(scanner.Text(), "date:") == true {
						year := strings.Split(strings.Split(scanner.Text(), ":")[1], "-")[0]
						if strings.Contains(year, y) == true {
							color.Cyan("Journal entry was created in the year " + year)
						}
					}
				}

				// MONTH

				if len(m) > 0 {
					if strings.Contains(scanner.Text(), "date:") == true {
						month := strings.Split(strings.Split(scanner.Text(), ":")[1], "-")[1]
						if strings.Contains(month, m) == true {
							color.Cyan("Journal entry was created in the month " + month)
						}
					}
				}

				// TEXT
				// output is deferred to end of file processing as it will otherwise intersperse with tag, year and month output

				if strings.Contains(strings.ToUpper(scanner.Text()), strings.ToUpper(text)) == true && len(text) > 0 {

					output = append(output, "Text is present on line "+strconv.Itoa(line))

					if v == true {
						re := regexp.MustCompile("(?i)" + text)
						caseText = re.FindAllString(scanner.Text(), -1)
						coloredText := re.ReplaceAllLiteralString(scanner.Text(), "%[1]s")
						outputVerbose = append(outputVerbose, coloredText)
						outputCaseText = append(outputCaseText, caseText)
					}

				}

				line++

			}

			if v == true && len(output) == 0 {
				fmt.Println(string(fc))
			} else {
				for i := 0; i < len(output); i++ {
					color.Cyan(output[i])
					if v == true && len(text) > 0 {
						red := color.New(color.Bold, color.FgRed).SprintFunc()
						fmt.Printf(outputVerbose[i], red(strings.ToUpper(text)))
						fmt.Print("\n\n")

						// fmt.Println(blue(fmt.Print(outputCaseText[i][0]), ","))
					}
				}
			}
		}
	} else {
		fmt.Println("Search query did not return any hits.")
	}
}
