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
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	// define variables

	t := time.Now()
	tStr := t.Format("2006-01-02T1504")
	tStrTitle := t.Format("02 January 2006")

	// define the flag command options

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	titleCreateFlag := createCommand.String("title", tStrTitle, "Title your diary entry. Default is today's date.")
	dateCreateFlag := createCommand.String("date", tStr, "Specify the date for your diary entry. Default is today.")

	searchCommand := flag.NewFlagSet("search", flag.ExitOnError)
	verboseSearchFlag := searchCommand.Bool("v", false, "Search verbosity. Default is false.")
	textSearchFlag := searchCommand.String("text", "", "Search text. Default is empty.")

	// define command switch

	switch os.Args[1] {
	case "create":
		createCommand.Parse(os.Args[2:])
	case "render":
		fmt.Println("Rendering!")
	case "search":
		searchCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	// Parse create command

	if createCommand.Parsed() {
		createEntry(*titleCreateFlag, *dateCreateFlag)
	}

	if searchCommand.Parsed() {
		search(".", *textSearchFlag, *verboseSearchFlag)
	}
}
