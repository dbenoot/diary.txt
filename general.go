//   Copyright 2022 The diarytxt Authors
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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func getFileList(wd string) ([]string, error) {

	var fileList []string

	err := filepath.Walk(wd, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	fileList = filterFile(fileList, wd)

	return fileList, err
}

func filterFile(f []string, wd string) []string {

	var fo []string

	// use regexp [0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{4}_['\\w,\\s-\\.]*\\.md
	//
	// meaning
	// [0-9]{4} : 4 digits
	// - : the character -
	// [0-9]{2} : 2 digits
	// - : the character -
	// [0-9]{2} : 2 digits
	// T : the character T
	// [0-9]{4} : 4 digits
	// _ : the character _
	// ['\\w,\\s-\\.]* : text consisting of  ', all word symbols, all whitespace symbols, dashes, dots
	// \\.md : .md
	//
	// Remark the double escape for w, s and . -> otherwise the string parser complains (and '' didn't work...)

	var r = regexp.MustCompile("[0-9]{4}[0-9]{2}[0-9]{2}T[0-9]{4}_['\\w,\\s-\\.]*\\.md")
	logdir := filepath.Join(wd, "logs")
	renderdir := filepath.Join(wd, "rendered")
	settingsdir := filepath.Join(wd, "settings")

	for _, file := range f {
		if r.MatchString(file) && !strings.Contains(file, renderdir) { //&& strings.Contains(file, ".md") {
			fo = append(fo, file)
		} else {
			fi, _ := os.Stat(file)
			if fi.Mode().IsRegular() && !strings.Contains(file, renderdir) && !strings.Contains(file, logdir) && !strings.Contains(file, settingsdir) {
				fmt.Printf("File was not included in the filterlist %s. Please check filterFile function. \n", file)
			}
		}
	}

	return fo
}

func checkDate(date string) bool {
	var r = regexp.MustCompile("[0-9]{4}[0-9]{2}[0-9]{2}T[0-9]{4}")

	d, _ := strconv.Atoi(getDay(date))
	m, _ := strconv.Atoi(getMonth(date))
	y, _ := strconv.Atoi(getYear(date))

	if r.MatchString(date) && date == r.FindString(date) && checkDay(d, m, y) && checkMonth(m) {
		return true
	}

	return false
}

func getYear(date string) string {
	var r = regexp.MustCompile("[0-9]{4}[0-9]{2}[0-9]{2}T[0-9]{4}")
	date = r.FindString(date)
	return date[0:4]
}

func getMonth(date string) string {
	var r = regexp.MustCompile("[0-9]{4}[0-9]{2}[0-9]{2}T[0-9]{4}")
	date = r.FindString(date)
	return date[4:6]
}

func getDay(date string) string {
	var r = regexp.MustCompile("[0-9]{4}[0-9]{2}[0-9]{2}T[0-9]{4}")
	date = r.FindString(date)
	return date[6:8]
}

func checkDay(day int, month int, year int) bool {

	thirtyone := []int{1, 3, 5, 7, 8, 10, 12}
	thirty := []int{4, 6, 9, 11}
	feb := []int{2}

	if contains(thirtyone, month) {
		if day > 0 && day < 32 {
			return true
		}
	}

	if contains(thirty, month) {
		if day > 0 && day < 31 {
			return true
		}
	}

	if contains(feb, month) {
		if year%4 == 0 && year%100 != 0 || year%400 == 0 {
			if day > 0 && day < 30 {
				return true
			}
		} else {
			if day > 0 && day < 29 {
				return true
			}
		}
	}
	return false
}

func checkMonth(m int) bool {
	if m > 0 && m < 13 {
		return true
	}
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
