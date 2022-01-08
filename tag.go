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
	"strings"

	"github.com/fatih/color"
)

func tag(i bool, l string, sd string, v bool) {

	tags := make(map[string]int)
	tagFiles := make(map[string][]string)
	var err error

	// get all the filenames

	fileList := []string{}
	fileList, err = getFileList(sd)
	check(err)

	// iterate over files and different tags and create the tags map

	for _, file := range fileList {
		f, _ := os.Open(file)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "tags:") {
				t := strings.TrimSpace(strings.Replace(scanner.Text(), "* tags:", "", -1))
				if len(t) != 0 {
					ts := strings.Split(t, ",")
					for _, tt := range ts {
						tt = strings.TrimSpace(tt)
						tags[tt]++
						tagFiles[tt] = append(tagFiles[tt], file)
					}
				}
			}
		}
	}

	// process index command

	if i {

		color.Green("Full index of the used tags and times used:")

		// sort numerically

		keys := make(map[string]int)

		for tag, _ := range tags {
			///keys[tag] = tags[tag]
			keys[tag] = tags[tag]
		}

		for _, nr := range rankByWordCount(keys) {
			fmt.Println(nr.Value, "time(s) used\t: ", nr.Key)
		}
	}

	// process list command

	if len(l) != 0 {
		if len(tagFiles[l]) != 0 {
			color.Green("The tag %s was used in the following files:", l)
			for a := range tagFiles[l] {
				color.Cyan(tagFiles[l][a])
				if v {
					fc, _ := ioutil.ReadFile(tagFiles[l][a])
					fmt.Println(string(fc))
				}
			}
		} else {
			color.Green("The tag %s is not present in your journal entries.", l)
		}
	}

	// return errors

	check(err)
}
