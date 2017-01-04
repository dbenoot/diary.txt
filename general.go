package main

import (
	"os"
	"path/filepath"
	"strings"
)

func getFileList(wd string) ([]string, error) {

	fileList := []string{}

	err := filepath.Walk(wd, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	fileList = filterFile(fileList)

	return fileList, err
}

func filterFile(f []string) []string {

	var fo []string

	for _, file := range f {
		if strings.Contains(file, "rendered_diary") == false && strings.Contains(file, ".md") {
			fo = append(fo, file)
		}
	}

	return fo
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
