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

	"github.com/pmylund/sortutil"

	//"io/ioutil"
	"os"
	//"path/filepath"
	"strings"
	"time"
)

func shame(wd string) {
	var fileList []string

	fileList, _ = getFileList(wd)
	//check(err)

	lastDate := time.Now()

	sortutil.CiAsc(fileList)

	lastfile := fileList[len(fileList)-1]

	scanfile, _ := os.Open(lastfile)
	scanner := bufio.NewScanner(scanfile)

	for scanner.Scan() {

		if strings.Contains(scanner.Text(), "* date: ") {
			lastDate, _ = time.Parse("20060102T1504", strings.TrimSpace(strings.Split(scanner.Text(), ":")[1]))
			//fmt.Println(lastDate)
		}
	}

	fmt.Println("Last entry on:", lastDate)

	diffDate := time.Since(lastDate)

	fmt.Println("Which was", diffDate.String(), "ago.")
}
