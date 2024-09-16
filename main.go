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
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

func main() {

	// define Diary

	type Diary struct {
		wd       string
		pins     []string
		copyPins bool
	}

	// define variables

	t := time.Now()
	tStr := t.Format("20060102T1504")
	tStrTitle := t.Format("02 January 2006")
	var localCfgFile string

	diary := Diary{}

	// Setup the diary home directory

	if len(os.Args) > 1 && os.Args[1] == "setup" {
		saveWorkDir(strings.Join(os.Args[2:], " "))
	}

	// check that the local config directory (sd) exists and if not create a preliminary config file

	diary.wd = setWorkDir()

	// Create and load diary specific config file (only the pins at this moment)
	// This makes it possible to sync diary specific settings when using cloud sync, as this sync file is in the diary folder.
	// Splitting up the config files gives the ability to run the executable from anywhere and not only in the diary folder and sanely keep diary specific settings when syncing.

	diary.pins, localCfgFile, diary.copyPins = setPins(diary.wd)

	// make sure the necessary paths are present

	createDirs(diary.wd)

	// define the flag command options

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	titleCreateFlag := createCommand.String("title", tStrTitle, "Title your diary entry. Default is today's date.")
	dateCreateFlag := createCommand.String("date", tStr, "Specify the date for your diary entry. Default is today.")
	textCreateFlag := createCommand.String("text", "", "Add text to the diary entry. Especially useful for short notes, for larger notes an editor is best used. Default is empty.")
	tagCreateFlag := createCommand.String("tag", "", "Add tags (comma-separated) to journal entry. Can also be added using editor. Default is empty.")
	pinCreateFlag := createCommand.Bool("pin", true, "Specify if the pins should be present. Notation example: -pin=false (include equal sign). Default is true.")
	pinCopyFlag := createCommand.Bool("copypin", diary.copyPins, "Copy pins contents from your most recent journal entry. Default is set in local config file.")

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

	tagCommand := flag.NewFlagSet("tag", flag.ExitOnError)
	indexTagFlag := tagCommand.Bool("index", false, "Shows all tags.")
	listTagFlag := tagCommand.String("list", "", "Shows all items for a specific tag.")
	verboseTagFlag := tagCommand.Bool("v", false, "Set the output verbosity. Default is false.")

	pinCommand := flag.NewFlagSet("pin", flag.ExitOnError)
	addPinFlag := pinCommand.String("add", "", "Add a pin. Default is empty.")
	removePinFlag := pinCommand.String("remove", "", "Remove a pin. Default is empty.")
	indexPinFlag := pinCommand.Bool("index", false, "Shows all prespecified pins.")
	indexAllPinFlag := pinCommand.Bool("indexall", false, "Shows all prespecifed pins with their distinct answers.")
	listPinFlag := pinCommand.String("list", "", "Shows all items for a specific pin.")

	// What to show when no arguments are defined

	if len(os.Args) == 1 {

		fmt.Println("Please provide arguments. Type 'diarytxt help' for more information.")

	} else {

		// define command switch

		switch os.Args[1] {
		case "create":
			createCommand.Parse(os.Args[2:])
		case "render":
			renderCommand.Parse(os.Args[2:])
		case "search":
			searchCommand.Parse(os.Args[2:])
		case "pin":
			pinCommand.Parse(os.Args[2:])
		case "tag":
			tagCommand.Parse(os.Args[2:])
		case "stat":
			statistics(diary.wd, tStr)
		case "shame":
			shame(diary.wd)
		case "help":
			printHelp()
		case "autotag":
			autotag(os.Args[2:], diary.wd)
		case "setup":
			fmt.Println("Setup done.")
		default:
			fmt.Printf("%q is not valid command.\n", os.Args[1])
			fmt.Println("Use the help command 'diarytxt help' for help.")
			os.Exit(2)
		}
	}

	// Parse commands

	if createCommand.Parsed() {
		createEntry(diary.wd, *titleCreateFlag, *dateCreateFlag, *tagCreateFlag, *pinCreateFlag, *pinCopyFlag, diary.pins, *textCreateFlag)
	}

	if searchCommand.Parsed() {
		search(diary.wd, *textSearchFlag, *tagSearchFlag, *yearSearchFlag, *monthSearchFlag, *verboseSearchFlag)
	}

	if renderCommand.Parsed() {
		render(diary.wd, *tagRenderFlag, *yearRenderFlag, *monthRenderFlag)
	}

	if tagCommand.Parsed() {
		tag(*indexTagFlag, *listTagFlag, diary.wd, *verboseTagFlag)
	}

	if pinCommand.Parsed() {
		pin(*addPinFlag, *removePinFlag, *listPinFlag, *indexPinFlag, *indexAllPinFlag, diary.wd, diary.pins, localCfgFile, os.Args)
	}

}

func setWorkDir() string {
	usr, err := user.Current()
	check(err)

	sd := filepath.Join(usr.HomeDir, ".config", "diarytxt")
	cfgFile := filepath.Join(sd, "config.ini")

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		_ = os.MkdirAll(sd, 0755)
		os.Create(cfgFile)
		var cfg, _ = ini.LooseLoad(cfgFile)
		_, _ = cfg.Section("general").NewKey("home", "")
		err = cfg.SaveTo(cfgFile)
		check(err)
	}
	check(err)

	var cfg, _ = ini.LooseLoad(cfgFile)

	if len(cfg.Section("general").Key("home").String()) == 0 {
		fmt.Printf("Home directory not set, please add to config.ini. Config file is located here: %s \n", sd)
		os.Exit(2)
	}

	return cfg.Section("general").Key("home").String()
}

func saveWorkDir(workdir string) {

	usr, err := user.Current()
	check(err)

	sd := filepath.Join(usr.HomeDir, ".config", "diarytxt")
	cfgFile := filepath.Join(sd, "config.ini")

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		_ = os.MkdirAll(sd, 0755)
		os.Create(cfgFile)
		var cfg, _ = ini.LooseLoad(cfgFile)
		_, _ = cfg.Section("general").NewKey("home", workdir)
		err = cfg.SaveTo(cfgFile)
		fmt.Println("HERE")
		check(err)
	}
}

func setPins(wd string) ([]string, string, bool) {
	settingsDir := filepath.Join(wd, "settings")
	localCfgFile := filepath.Join(settingsDir, "local_config.ini")

	if _, err := os.Stat(localCfgFile); os.IsNotExist(err) {
		_ = os.MkdirAll(settingsDir, 0755)
		os.Create(localCfgFile)
		var localCfg, _ = ini.LooseLoad(localCfgFile)
		_, _ = localCfg.Section("general").NewKey("pins", "")
		_, _ = localCfg.Section("general").NewKey("copy_pin_content", "true")
		err = localCfg.SaveTo(localCfgFile)
		check(err)
	}

	var localCfg, _ = ini.LooseLoad(localCfgFile)

	copyPins, _ := localCfg.Section("general").Key("copy_pin_content").Bool()

	return localCfg.Section("general").Key("pins").Strings(","), localCfgFile, copyPins
}

func createDirs(wd string) {
	renderdir := filepath.Join(wd, "rendered")
	if _, err := os.Stat(renderdir); os.IsNotExist(err) {
		_ = os.MkdirAll(renderdir, 0755)
	}

	logdir := filepath.Join(wd, "logs")
	if _, err := os.Stat(logdir); os.IsNotExist(err) {
		_ = os.MkdirAll(logdir, 0755)
	}

	filedir := filepath.Join(wd, "files")
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		_ = os.MkdirAll(filedir, 0755)
	}
}

func printHelp() {
	fmt.Println("Please provide secondary command.")
	fmt.Println("")
	fmt.Println("The following commands can be issued:")
	fmt.Println("")
	fmt.Println("create          Creates a new journal entry")
	fmt.Println("  -title        Title your diary entry. Default is today's date.")
	fmt.Println("  -date         Specify the date for your diary entry. Format should be yyyymmddThhmm (e.g. 20060102T1504). Default is today.")
	fmt.Println("  -text         Add text to the diary entry. Especially useful for short notes, for larger notes an editor is best used. Default is empty.")
	fmt.Println("  -tag          Add tags (comma-separated) to journal entry. Can also be added using editor. Default is empty.")
	fmt.Println("  -pin          Specify if the pins should be present. Notation example: -pin=false (include equal sign). Default is true.")
	fmt.Println("  -copypin      Copies the pin content from the last written journal entry.")
	fmt.Println("")
	fmt.Println("render          Renders your diary entries to a single markdown and html document located in the rendered folder in your diary home directory.")
	fmt.Println("  -tag          Render journal entries with a specific tag. Default is empty.")
	fmt.Println("  -year         Render journal entries from a specific year. Default is empty.")
	fmt.Println("  -month        Render journal entries from a specific month. Default is empty. Format is 2 digit numeric.")
	fmt.Println("")
	fmt.Println("search          Search your journal entries")
	fmt.Println("  -tag          Search for entries with a specific tag. Default is empty.")
	fmt.Println("  -text         Search text. Default is empty.")
	fmt.Println("  -year         Search journal items for a specific year. Default is empty.")
	fmt.Println("  -month        Search journal items for a specific month. Default is empty.")
	fmt.Println("  -v            Set the output verbosity. Default is false.")
	fmt.Println("")
	fmt.Println("tag")
	fmt.Println("  -index        Shows all tags.")
	fmt.Println("  -list         Shows all journal entries for a specific tag.")
	fmt.Println("  -v            Set the output verbosity. Default is false.")
	fmt.Println("")
	fmt.Println("autotag [date]  Checks the content of the text for words that have been tagged before and adds these to the tags list.")
	fmt.Println("")
	fmt.Println("pin             Administrate the journal pins.")
	fmt.Println("  -add          Add a pin. Default is empty.")
	fmt.Println("  -remove       Remove a pin. Default is empty.")
	fmt.Println("  -index        Shows all prespecified pins.")
	fmt.Println("  -indexall     Shows all prespecifed pins with their distinct answers.")
	fmt.Println("  -list         Shows all items for a specific pin.")
	fmt.Println("")
	fmt.Println("stat            Some basic statistics about your diary. Output is saved in the logs folder in your diary home directory.")
	fmt.Println("")
	fmt.Println("setup           Setup the diary home directory.")
	fmt.Println("")
}
