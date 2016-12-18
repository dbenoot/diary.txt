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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func createEntry(wd string, title string, t string, p []string) {

	// Derive the year and month

	y := strings.Split(t, "-")[0]
	m := strings.Split(t, "-")[1]

	//Check if the subdir for this year and month already exists. If not, create it.

	dir := filepath.Join(wd, y, m)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}

	filename := t + "_" + title + ".md"
	file := filepath.Join(dir, filename)

	// Create the markdown file as YYYYMMDD-HHMM_title(if specified).md
	if _, err := os.Stat(file); err == nil {
		fmt.Println("File already exists.")
	} else {
		os.Create(file)
		fmt.Println("Created ", file)

		// Add title and tags to file
		// Entry title has md headertag ### as the year and month will get # and ## respectively during the render

		content := "### " + title + "\n\ndate: " + t + "\ntags: "

		for i := 0; i < len(p); i++ {
			content = content + "\n" + p[i] + ":"
		}

		err := ioutil.WriteFile(file, []byte(content), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
