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
	"os"
	"strings"

	"github.com/fatih/color"
)

func getTags(sd string) (map[string]int, map[string][]string) {
	// get all the filenames
	tags := make(map[string]int)
	tagFiles := make(map[string][]string)
	var err error

	var fileList []string
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
	return tags, tagFiles
}

func tag(i bool, l string, sd string, v bool) {

	var err error

	tags, tagFiles := getTags(sd)

	// process index command

	if i {

		color.Green("Full index of the used tags and times used:")

		// sort numerically

		keys := make(map[string]int)

		for tag := range tags {
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
					fc, _ := os.ReadFile(tagFiles[l][a])
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

func autotag(vars []string, wd string) {

	var err error
	var tagList []string

	var fileList []string
	fileList, err = getFileList(wd)
	check(err)

	tags, _ := getTags(wd)

	// iterate over the files and select the file with the correct search string
	// afterwards iterate over the tags and check if they are present in the list of tags.

	for _, f := range fileList {
		if strings.Contains(f, vars[0]) {
			content, err := os.ReadFile(f)
			check(err)

			cntStr := string(content)

			for t := range tags {
				if strings.Contains(strings.ToLower(cntStr), strings.ToLower(t)) {
					tagList = AppendIfMissing(tagList, strings.ToLower(t))
				}
			}

			// Replace the empty tags field with the completed tags field

			newContent := strings.Replace(cntStr, "* tags:", "* tags: "+stringSliceToString(tagList), 1)

			// Write the file

			err = os.WriteFile(f, []byte(newContent), 0644)
			check(err)
		}
	}
}

func stringSliceToString(slice []string) string {
	var nt string

	// Iterate over the slice and add to string, adding a comma in between iterations

	for _, tl := range slice {
		nt = nt + tl + ", "
	}

	// Remove the last trailing comma

	nt = strings.TrimSuffix(nt, ", ")

	return nt
}
