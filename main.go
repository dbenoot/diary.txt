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
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/go-ini/ini"
)

func main() {

	// define variables

	t := time.Now()
	tStr := t.Format("2006-01-02T1504")
	tStrTitle := t.Format("02 January 2006")

	// define Diary

	type Diary struct {
		wd   string
		pins []string
	}

	// Check that the settings directory (sd) exists and if not create a preliminary config file
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	sd := filepath.Join(usr.HomeDir, ".diarytxt")
	cfgFile := filepath.Join(sd, "config.ini")

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		_ = os.MkdirAll(sd, 0755)
		os.Create(cfgFile)
	}

	var cfg, _ = ini.LooseLoad(cfgFile)

	var diary = Diary{
		wd:   cfg.Section("general").Key("home").String(),
		pins: cfg.Section("general").Key("pins").Strings(","),
	}

	if len(diary.wd) == 0 {
		fmt.Println("Home dir not set, please add to config.ini")
	}

	// define the flag command options

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	titleCreateFlag := createCommand.String("title", tStrTitle, "Title your diary entry. Default is today's date.")
	dateCreateFlag := createCommand.String("date", tStr, "Specify the date for your diary entry. Default is today.")
	textCreateFlag := createCommand.String("text", "", "Add text to the diary entry. Especially useful for short notes, for larger notes an editor is best used. Default is empty.")
	tagCreateFlag := createCommand.String("tag", "", "Add tags (comma-separated) to journal entry. Can also be added using editor. Default is empty.")
	pinCreateFlag := createCommand.Bool("pin", true, "Specify if the pins should be present. Notation example: -pin=false (include equal sign). Default is true.")

	searchCommand := flag.NewFlagSet("search", flag.ExitOnError)
	verboseSearchFlag := searchCommand.Bool("v", false, "Set the output verbosity. Default is false.")
	textSearchFlag := searchCommand.String("text", "", "Search text. Default is empty.")
	tagSearchFlag := searchCommand.String("tag", "", "Search for entries with a specific tag. Default is empty.")
	yearSearchFlag := searchCommand.String("year", "", "Search for entries with a specific year. Default is empty.")
	monthSearchFlag := searchCommand.String("month", "", "Search for entries with a specific month. Default is empty.")

	renderCommand := flag.NewFlagSet("render", flag.ExitOnError)
	tagRenderFlag := renderCommand.String("tag", "", "Render journal entries with a specific tag. Default is empty.")
	yearRenderFlag := renderCommand.String("year", "", "Render journal entries from a specific year. Default is empty.")
	monthRenderFlag := renderCommand.String("month", "", "Render journal entries from a specific month. Default is empty. Format is 2 digit numeric.")

	setCommand := flag.NewFlagSet("set", flag.ExitOnError)
	// wdSetFlag := setCommand.String("home", "~/diary", "Set the home directory. The default is ~/diary")
	// addPinSetFlag := setCommand.String("add-pin", "", "Add a pinned item. A pinned item is an item that will be created in all new journal entries. E.g. weight, book reading,...")
	// removePinSetFlag := setCommand.String("remove-pin", "", "Remove a pinned item.")

	// What to show when no arguments are defined

	if len(os.Args) == 1 {
		fmt.Println("Please provide secondary command.")
		fmt.Println("")
		fmt.Println("The following commands can be issued:")
		fmt.Println("")
		fmt.Println("create          Creates a new journal entry")
		fmt.Println("  -title        Title your diary entry. Default is today's date.")
		fmt.Println("  -date         Specify the date for your diary entry. Format should be yyyy-mm-ddThhmm (e.g. 2006-01-02T1504). Default is today.")
		fmt.Println("  -text         Add text to the diary entry. Especially useful for short notes, for larger notes an editor is best used. Default is empty.")
		fmt.Println("  -tag          Add tags (comma-separated) to journal entry. Can also be added using editor. Default is empty.")
		fmt.Println("  -pin          Specify if the pins should be present. Notation example: -pin=false (include equal sign). Default is true.")
		fmt.Println("")
		fmt.Println("render          Renders your diary entries to a single markdown and html document located in the rendered subfolder in your diary home directory.")
		fmt.Println("  -tag          Render journal entries with a specific tag. Default is empty.")
		fmt.Println("  -year         Render journal entries from a specific year. Default is empty.")
		fmt.Println("  -month        Render journal entries from a specific month. Default is empty. Format is 2 digit numeric.")
		fmt.Println("")
		fmt.Println("search          Search your journal entries")
		fmt.Println("  -tag          Search for entries with a specific tag. Default is empty.")
		fmt.Println("  -text         Search text. Default is empty.")
		fmt.Println("  -v            Set the output verbosity. Default is false.")
	} else {

		// define command switch

		switch os.Args[1] {
		case "create":
			createCommand.Parse(os.Args[2:])
		case "render":
			renderCommand.Parse(os.Args[2:])
		case "search":
			searchCommand.Parse(os.Args[2:])
		case "set":
			setCommand.Parse(os.Args[2:])
		default:
			fmt.Printf("%q is not valid command.\n", os.Args[1])
			os.Exit(2)
		}
	}

	// Parse commands

	if createCommand.Parsed() {
		createEntry(diary.wd, *titleCreateFlag, *dateCreateFlag, *tagCreateFlag, *pinCreateFlag, diary.pins, *textCreateFlag)
	}

	if searchCommand.Parsed() {
		search(diary.wd, *textSearchFlag, *tagSearchFlag, *yearSearchFlag, *monthSearchFlag, *verboseSearchFlag)
	}

	if renderCommand.Parsed() {
		render(diary.wd, *tagRenderFlag, *yearRenderFlag, *monthRenderFlag)
	}
}
