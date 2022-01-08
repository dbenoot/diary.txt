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
	"path/filepath"
	"strings"

	"github.com/pmylund/sortutil"
)

func createEntry(wd string, title string, t string, tag string, pb bool, cp bool, p []string, c string) {

	// check that t is in the correct format

	if !checkDate(t) {
		fmt.Println("Specified datetime format is incorrect.")
		os.Exit(1)
	}

	// Derive the year and month

	y := getYear(t)
	m := getMonth(t)

	//Check if the subdir for this year and month already exists. If not, create it.

	dir := filepath.Join(wd, y, m)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}

	// create filename

	filename := t + "_" + title + ".md"
	file := filepath.Join(dir, filename)

	// Create the markdown file as YYYYMMDDTHHMM_title(if specified).md
	if _, err := os.Stat(file); err == nil {

		fmt.Println("A journal entry for this date and time already exists.")

	} else {

		os.Create(file)
		fmt.Println("Created ", file)

		// Add title and tags to file
		// Entry title has md headertag ### as the year and month will get # and ## respectively during the render

		content := "### " + title + "\n\n* date: " + t + "\n* tags: " + tag

		// copy pins content from previous entry, if applicable

		if cp {

			fileList := []string{}

			fileList, err = getFileList(wd)
			check(err)

			sortutil.CiAsc(fileList)

			lastfile := fileList[len(fileList)-2]

			fmt.Printf("Pin content copied from %s \n", lastfile)
			scanfile, _ := os.Open(lastfile)
			scanner := bufio.NewScanner(scanfile)

			var added []string

			for scanner.Scan() {

				for _, a := range p {

					if strings.Contains(scanner.Text(), "* "+a+":") {
						content = content + "\n" + scanner.Text()
						added = append(added, a)
					}

				}
			}
			tocreate := difference(added, p)
			if len(tocreate) != 0 {
				for _, add := range tocreate {
					content = content + "\n" + "* " + add + ":"
				}
			}
		}

		// add pins if applicable

		if pb && len(p) > 0 && !cp {
			for i := 0; i < len(p); i++ {
				content = content + "\n* " + p[i] + ":"
			}
		}

		// add content in case it is included in the command line

		if len(c) > 0 {
			content = content + "\n\n" + c
		}

		// write the file

		err := ioutil.WriteFile(file, []byte(content), 0644)
		check(err)

	}
}
