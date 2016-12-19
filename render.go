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

// - read all files
//   filepath.walk??
// - sort files
//   sort.Strings(strs)
// - concatenate all the resulting strings
// - md -> html

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func render(location string) (err error) {

	// dir and file location

	ed := filepath.Join(location, "rendered")
	ef := filepath.Join(ed, "temp.md")

	//create dir and file

	if _, err = os.Stat(ed); os.IsNotExist(err) {
		fmt.Println(err, "ok, passed here")
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

	// concatenate contents
	var y string
	var m string
	buf := bytes.NewBuffer(nil)
	for _, file := range fileList {
		if strings.Contains(file, ".md") == true {
			// if strings.Split(file, "-")[0] != y {
			// 	y = strings.Split(file, "-")[0]
			// 	io.Copy(buf, ("#" + y)) // copy year string to file...
			// }
			f, _ := os.Open(file)
			io.Copy(buf, f)
			f.Close()
		}
	}
	err = ioutil.WriteFile(ef, buf.Bytes(), 0644)

	// defer exportFile.Close()
	return err

}
