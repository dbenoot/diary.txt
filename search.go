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
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func search(location string, text string, tag string, y string, m string, v bool) (err error) {
	// fileList := []string{}
	var fileList, tSlice, tagSlice, ySlice, mSlice []string

	fileList, err = getFileList(location)

	tagSlice = strings.Split(tag, " ")
	ySlice = strings.Split(y, " ")
	mSlice = strings.Split(m, " ")
	tSlice = strings.Split(text, " ")

	fileList = filterTag(fileList, tagSlice)
	fileList = filterYear(fileList, ySlice)
	fileList = filterMonth(fileList, mSlice)
	fileList = filterText(fileList, tSlice)

	report(fileList, tSlice, tagSlice, ySlice, mSlice, v)

	return err
}

func filterTag(f []string, t []string) []string {

	var fileList []string

	for _, file := range f {
		fo, _ := os.Open(file)

		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {
			for _, tt := range t {
				if strings.Contains(scanner.Text(), "* tags:") && strings.Contains(scanner.Text(), tt) {
					fileList = append(fileList, file)
				}
			}
		}
	}

	return fileList
}

func filterYear(f []string, y []string) []string {

	var fileList []string

	for _, file := range f {
		fo, _ := os.Open(file)

		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {

			if strings.Contains(scanner.Text(), "* date:") {
				year := getYear(scanner.Text())
				for _, yy := range y {
					if strings.Contains(year, yy) {
						fileList = append(fileList, file)
					}
				}
			}
		}
	}

	return fileList
}

func filterMonth(f []string, m []string) []string {

	var fileList []string

	for _, file := range f {
		fo, _ := os.Open(file)

		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "* date:") {
				month := getMonth(scanner.Text())
				for _, mm := range m {
					if strings.Contains(month, mm) {
						fileList = append(fileList, file)
					}
				}
			}
		}
	}

	return fileList
}

func filterText(f []string, t []string) []string {

	var fileList []string

	for _, file := range f {

		c := 0

		fo, _ := os.Open(file)

		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {
			for _, tt := range t {
				if strings.Contains(strings.ToUpper(scanner.Text()), strings.ToUpper(tt)) {
					c++
				}
			}
		}

		if c > 0 {
			fileList = append(fileList, file)
		}
	}

	return fileList
}

func report(f []string, text []string, tag []string, y []string, m []string, v bool) {
	if len(f) > 0 {
		for _, file := range f {

			var outputGrp, outputVerbose []string
			var output []int

			color.Green("Search criteria present in journal entry " + file)

			f, _ := os.Open(file)
			fc, _ := ioutil.ReadFile(file)

			scanner := bufio.NewScanner(f)

			line := 1

			for scanner.Scan() {

				// TAGS

				if len(tag) > 0 {
					if strings.Contains(scanner.Text(), "tags:") {
						for _, tt := range tag {
							if strings.Contains(strings.Split(scanner.Text(), ":")[1], tt) && len(tt) != 0 {
								color.Cyan("Journal entry is tagged with " + tt)
							}
						}
					}
				}

				// YEAR

				if len(y) > 0 {
					if strings.Contains(scanner.Text(), "date:") {
						year := getYear(scanner.Text())
						for _, yy := range y {
							if strings.Contains(year, yy) && len(yy) != 0 {
								color.Cyan("Journal entry was created in the year " + year)
							}
						}
					}
				}

				// MONTH

				if len(m) > 0 {
					if strings.Contains(scanner.Text(), "date:") {
						month := getMonth(scanner.Text())
						for _, mm := range m {
							if strings.Contains(month, mm) && len(mm) != 0 {
								color.Cyan("Journal entry was created in the month " + month)
							}
						}
					}
				}

				// TEXT
				// output is deferred to end of file processing as it will otherwise intersperse with tag, year and month output

				for _, tt := range text {
					if strings.Contains(strings.ToUpper(scanner.Text()), strings.ToUpper(tt)) && len(tt) > 0 {

						output = append(output, line)
						outputGrp = append(outputGrp, tt)

						if v {
							outputVerbose = append(outputVerbose, scanner.Text())
						}
					}
				}

				line++

			}

			// output for non-text searches with verbosity on

			if v && len(output) == 0 {
				fmt.Println(string(fc))
			} else {

				// output for text searches with verbosity on

				for i := 0; i < len(output); i++ {

					color.Cyan("Text " + outputGrp[i] + " is present on line " + strconv.Itoa(output[i]))

					if v {
						red := color.New(color.Bold, color.FgRed).SprintFunc()
						re := regexp.MustCompile("(?i)" + outputGrp[i])
						coloredText := re.ReplaceAllLiteralString(outputVerbose[i], "%[1]s")
						fmt.Printf(coloredText, red(strings.ToUpper(outputGrp[i])))
						fmt.Print("\n\n")
					}
				}
			}
		}
	} else {
		fmt.Println("Search query did not return any hits.")
	}
}
