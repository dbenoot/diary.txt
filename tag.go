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
	"github.com/pmylund/sortutil"
	"os"
	"path/filepath"
	"strings"
)

func tag(i bool, l string, sd string) {

	tags := make(map[string]int)
	tagFiles := make(map[string][]string)
	var err error

	// get all the filenames

	fileList := []string{}
	err = filepath.Walk(sd, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	// filter the incorrect filenames out

	fileList = filterFile(fileList)

	// iterate over files and different tags and create the tags map

	for _, file := range fileList {
		f, _ := os.Open(file)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "tags:") == true {
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

	if i == true {
		color.Green("Full index of the used tags and times used:")

		// sort alphabetically

		var keys []string
		for k := range tags {
			keys = append(keys, k)
		}

		sortutil.CiAsc(keys)

		// report based on alphabetical slice

		for _, k := range keys {
			fmt.Println(k, "\t", tags[k])
		}
	}

	// process list command

	if len(l) != 0 {
		if len(tagFiles[l]) != 0 {
			color.Green("The tag %s was used in the following files:", l)
			for a, _ := range tagFiles[l] {
				fmt.Println(tagFiles[l][a])
			}
		} else {
			color.Green("The tag %s is not present in your journal entries.", l)
		}
	}

	// return errors

	if err != nil {
		fmt.Println(err)
	}
}
