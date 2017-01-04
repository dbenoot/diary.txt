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

	var r = regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{4}_['\\w,\\s-\\.]*\\.md")
	logdir := filepath.Join(wd, "logs")
	renderdir := filepath.Join(wd, "rendered")

	for _, file := range f {
		if r.MatchString(file) && strings.Contains(file, renderdir) == false { //&& strings.Contains(file, ".md") {
			fo = append(fo, file)
		} else {
			fi, _ := os.Stat(file)
			if fi.Mode().IsRegular() == true && strings.Contains(file, renderdir) == false && strings.Contains(file, logdir) == false {
				fmt.Printf("File was not included in the filterlist %s. Please check filterFile function. \n", file)
			}
		}
	}

	return fo
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
