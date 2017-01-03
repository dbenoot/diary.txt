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
	"bytes"
	// "fmt"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func subselectTag(f []string, t string) []string {

	list := []string{}
	for i := 0; i < len(f); i++ {
		file, _ := os.Open(f[i])
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "tags:") == true {
				r := strings.Split(scanner.Text(), ":")[1]
				if strings.Contains(r, t) {
					list = append(list, f[i])
				}
			}
		}
	}

	return list
}

func subselectYear(f []string, y string) []string {

	list := []string{}

	for i := 0; i < len(f); i++ {
		file, _ := os.Open(f[i])
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "date:") == true {
				r := strings.Split(strings.Split(scanner.Text(), ":")[1], "-")[0]
				if strings.Contains(r, y) {
					list = append(list, f[i])
				}
			}
		}
	}

	return list
}

func subselectMonth(f []string, m string) []string {
	list := []string{}

	for i := 0; i < len(f); i++ {
		file, _ := os.Open(f[i])
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "date:") == true {
				r := strings.Split(strings.Split(scanner.Text(), ":")[1], "-")[1]
				if strings.Contains(r, m) {
					list = append(list, f[i])
				}
			}
		}
	}

	return list
}

func render(location string, tag string, year string, month string) (err error) {

	// dir and file location

	ed := filepath.Join(location, "rendered")
	ef := filepath.Join(ed, "rendered_diary.md")
	hf := filepath.Join(ed, "rendered_diary.html")

	//create dir and file

	if _, err = os.Stat(ed); os.IsNotExist(err) {
		_ = os.MkdirAll(ed, 0755)
	}

	os.Create(ef)

	// create fileList of all .md files

	fileList := []string{}

	err = filepath.Walk(location, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	// sort the fileList alphabetically

	sort.Strings(fileList)

	// subselect fileList on tag, year and month

	if len(tag) > 0 {
		fileList = subselectTag(fileList, tag)
	}

	if len(year) > 0 {
		fileList = subselectYear(fileList, year)
	}

	if len(month) > 0 {
		fileList = subselectMonth(fileList, month)
	}

	// concatenate contents
	var y string
	var m int
	buf := bytes.NewBuffer(nil)
	for _, file := range fileList {
		if strings.Contains(file, ".md") == true && strings.Contains(file, "rendered_diary.md") == false { //Change this a more robust file structure checker: only if file structure is correct, include the file
			if strings.Split(strings.Split(file, "/")[len(strings.Split(file, "/"))-1], "-")[0] != y {
				y = strings.Split(strings.Split(file, "/")[len(strings.Split(file, "/"))-1], "-")[0]
				b := bytes.NewBufferString("# " + y + "\n")
				io.Copy(buf, b)
				m = 0
			}

			if strings.Split(strings.Split(file, "/")[len(strings.Split(file, "/"))-1], "-")[1] != strconv.Itoa(m) {
				m, err = strconv.Atoi(strings.Split(strings.Split(file, "/")[len(strings.Split(file, "/"))-1], "-")[1])
				months := [...]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
				b := bytes.NewBufferString("## " + months[m-1] + "\n")
				io.Copy(buf, b)
			}

			nl := bytes.NewBufferString("\n\n")
			f, _ := os.Open(file)
			io.Copy(buf, f)
			io.Copy(buf, nl)
			f.Close()
		}
	}
	err = ioutil.WriteFile(ef, buf.Bytes(), 0644)

	// Create HTML

	output := blackfriday.MarkdownBasic(buf.Bytes())

	r := bytes.NewReader(output)

	html := bytes.NewBufferString("<html>\n<head>\n<link href=\"https://fonts.googleapis.com/css?family=Fira+Sans|Lobster\" rel=\"stylesheet\">\n\t<style>\n\t\thtml {\n\t\t\tbackground: #2b292e;\n\t\t\t}\n\t\tbody {\n\t\t\twidth: 60%;\n\t\t\tmargin-left: auto;\n\t\t\tmargin-right: auto;\n\t\t\tmargin-top: 0;\n\t\t\tpadding: 0 10px 50px 10px;\n\t\t\tbox-shadow: 0 0 10px white;\n\t\t\tfont-family: 'Fira Sans', sans-serif;\n\t\t\tfont-size: 1em;\n\t\t\tbackground: whitesmoke;\n\t\t}\n\t\th1 {\n\t\t\tcolor: orange;\n\t\t\tfont-family: 'Lobster', cursive;\n\t\t\tfont-size: 100;\n\t\t\ttext-align: center;\n\t\t}\n\t\th2 {\n\t\t\tcolor: #336699;\n\t\t\tfont-family: 'Lobster', cursive;\n\t\t\tfont-size: 50;\n\t\t\ttext-align: center;\n\t\t}\n\t\th3 {\n\t\t\ttext-align: center;\n\t\t\t}\n\t\tul {\n\t\t\tfont-size: 0.7em;\n\t\t}\n\t\tul {; text-align: center;}\n\t\tul li { list-style: none; display: inline}\n\t\tul li:after { content: \" \\00b7\"; }\n\t\tul li:last-child:after { content: none; }\n\t\tp {\n\t\t\tpadding-left: 5%;\n\t\t\tpadding-right: 5%;\n\t\t\tline-height: 1.8;\n\t\t\tcolor: #51565C;\n\t\t}\n\t</style>\n</head>\n<body>\n")
	htmlEnd := bytes.NewBufferString("\n</body>\n</html>")

	io.Copy(html, r)
	io.Copy(html, htmlEnd)

	err = ioutil.WriteFile(hf, html.Bytes(), 0644)

	// defer exportFile.Close()
	return err

}
