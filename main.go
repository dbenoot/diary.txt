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
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {

	// define variables
	switch os.Args[1] {
	case "create":
		createEntry()
	case "render":
		fmt.Println("Rendering!")
	case "search":
		fmt.Println(os.Args[2:])
		//search(os.Args[2:])
	case "archive":
		//archiveCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

func createEntry() {

	//Check if the subdir for this year and month already exists. If not, create it.
	t := time.Now()
	dir := filepath.Join(strconv.Itoa(t.Year()), strconv.Itoa(int(t.Month())))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
		// if err2 != nil {
		// 	return err2
		// }
	}

	// Create the markdown file as YYYYMMDD-HHMM_title(if specified).md

	os.Create(filepath.Join(dir, t.Format("2006-01-02T1504")+"_"+os.Args[2]+".md"))
	fmt.Println("Created ", filepath.Join(dir, t.Format("2006-01-02T1504")+"_"+os.Args[2]+".md"))
}
