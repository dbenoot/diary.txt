package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getFileList(wd string) ([]string, error) {

	fileList := []string{}

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
		if r.MatchString(file) && strings.Contains(file, renderdir) == false { //&& strings.Contains(file, ".md") {
			fo = append(fo, file)
		} else {
			fi, _ := os.Stat(file)
			if fi.Mode().IsRegular() == true && strings.Contains(file, renderdir) == false && strings.Contains(file, logdir) == false && strings.Contains(file, settingsdir) == false {
				fmt.Printf("File was not included in the filterlist %s. Please check filterFile function. \n", file)
			}
		}
	}

	return fo
}

func getYear(y string) string {

	var r = regexp.MustCompile("[0-9]{4}[0-9]{2}[0-9]{2}T[0-9]{4}")

	y = r.FindString(y)

	return y[0:4]
}

func getMonth(y string) string {

	var r = regexp.MustCompile("[0-9]{4}[0-9]{2}[0-9]{2}T[0-9]{4}")

	y = r.FindString(y)
	return y[4:6]
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
