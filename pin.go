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
	"github.com/go-ini/ini"
	"github.com/pmylund/sortutil"
)

func pin(a string, r string, l string, i bool, ia bool, sd string, pins []string, archived_pins []string, cfgFile string, args []string) {

	var err error
	var cfg, _ = ini.LooseLoad(cfgFile)
	var p, ap string

	switch {
	// ADD PIN
	case len(a) > 0:
		if stringInSlice(a, cfg.Section("general").Key("archived_pins").Strings(",")) {
			fmt.Printf("Pin '%s' has been archived in the past and will be reactivated.\n", a)
		}
		ap = removeIniKey(cfg.Section("general").Key("archived_pins").Strings(","), a)
		p = appendToIniKey(cfg.Section("general").Key("pins").String(), a)

		_, _ = cfg.Section("general").NewKey("pins", p)
		_, _ = cfg.Section("general").NewKey("archived_pins", ap)
		err = cfg.SaveTo(cfgFile)
		check(err)

	// Remove a pin & TODO add it to a section archived_pins in the config file
	case len(r) > 0:

		if !stringInSlice(r, cfg.Section("general").Key("pins").Strings(",")) {
			fmt.Printf("Pin '%s' not defined.\n", r)
		} else {
			p = removeIniKey(cfg.Section("general").Key("pins").Strings(","), r)
			ap = appendToIniKey(cfg.Section("general").Key("archived_pins").String(), r)
		}

		_, _ = cfg.Section("general").NewKey("pins", p)
		_, _ = cfg.Section("general").NewKey("archived_pins", ap)
		err = cfg.SaveTo(cfgFile)
		check(err)

	// Index of pins
	case i:
		color.Green("Specified pins:")
		for i := range pins {
			fmt.Println(pins[i])
		}

		if len(archived_pins) > 0 {
			color.Red("Archived pins:")
			for i := range archived_pins {
				fmt.Println(archived_pins[i])
			}
		}
	// Indexall: list of pins and the contents
	case ia:
		var fileList []string
		fileList, err = getFileList(sd)
		check(err)

		color.Green("Full index of the specified pins and their unique values:")

		for i := range pins {

			var index []string

			for _, file := range fileList {
				fo, _ := os.Open(file)

				scanner := bufio.NewScanner(fo)

				for scanner.Scan() {
					if strings.Contains(scanner.Text(), "* "+pins[i]+":") {
						item := strings.TrimSpace(strings.Replace(scanner.Text(), "* "+pins[i]+":", "", -1))
						if len(item) > 0 {
							index = AppendIfMissing(index, item)
						}
					}
				}
			}

			color.Cyan("Entries for pin " + pins[i])
			if len(index) != 0 {
				sortutil.CiAsc(index)
				for j := range index {
					fmt.Printf("%v : %s \n", j+1, index[j])
				}
			}
		}

		color.Red("Full index of the archived pins and their unique values:")

		for i := range archived_pins {

			var index []string

			for _, file := range fileList {
				fo, _ := os.Open(file)

				scanner := bufio.NewScanner(fo)

				for scanner.Scan() {
					if strings.Contains(scanner.Text(), "* "+archived_pins[i]+":") {
						item := strings.TrimSpace(strings.Replace(scanner.Text(), "* "+archived_pins[i]+":", "", -1))
						if len(item) > 0 {
							index = AppendIfMissing(index, item)
						}
					}
				}
			}

			color.Cyan("Entries for pin " + archived_pins[i])
			if len(index) != 0 {
				sortutil.CiAsc(index)
				for j := range index {
					fmt.Printf("%v : %s \n", j+1, index[j])
				}
			}
		}

	// List function: shows all diary entries with a specific pin
	case len(l) > 0:
		var fileList []string
		fileList, err = getFileList(sd)

		var index []string
		var date []string
		var files []string

		for _, file := range fileList {

			var d string

			fo, _ := os.Open(file)

			scanner := bufio.NewScanner(fo)

			for scanner.Scan() {

				if strings.Contains(scanner.Text(), "* date:") {
					d = strings.TrimSpace(strings.Replace(scanner.Text(), "* date:", "", -1))
				}

				if strings.Contains(scanner.Text(), "* "+l+":") {
					item := strings.TrimSpace(strings.Replace(scanner.Text(), "* "+l+":", "", 1))
					if len(item) > 0 {
						index = append(index, item)
						date = append(date, d)
						files = append(files, file)
					}
				}
			}

		}

		if len(index) == 0 {

			color.Red("Pin '" + l + "' not present or never completed in journal entries.")

		} else {

			color.Green("Dated entries for pin '" + l + "':")

			for k := range index {
				fmt.Printf("%v \t %s \t %s \n", date[k], index[k], files[k])
			}
		}
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		fmt.Println("Use the help command 'diarytxt help' for help.")
		os.Exit(2)
	}

}

func appendToIniKey(ini string, newkey string) string {
	if len(ini) == 0 {
		ini = newkey
	} else {
		ini = ini + ", " + newkey
	}
	return ini
}

func removeIniKey(ini []string, key string) string {
	var p string
	for i, v := range ini {
		if v == key {
			ini = append(ini[:i], ini[i+1:]...)
			break
		}
	}

	// comma-separate the strings in the slice and export to a single string

	for j := range ini {
		if j == len(ini)-1 {
			p = p + ini[j]
		} else {
			p = p + ini[j] + ", "
		}
	}
	return p
}
