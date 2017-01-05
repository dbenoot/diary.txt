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
	"github.com/go-ini/ini"
	"github.com/pmylund/sortutil"
	"os"
	"strings"
)

func pin(a string, r string, l string, i bool, ia bool, sd string, pins []string, cfgFile string, args []string) {

	var err error

	if len(a) > 0 {
		var cfg, _ = ini.LooseLoad(cfgFile)
		p := cfg.Section("general").Key("pins").String()
		if len(p) == 0 {
			p = a
		} else {
			p = p + ", " + a
		}
		_, _ = cfg.Section("general").NewKey("pins", p)
		err = cfg.SaveTo(cfgFile)
	}

	if len(r) > 0 {
		var p string

		if stringInSlice(r, pins) == false {
			fmt.Printf("Pin '%s' not defined.\n", r)
		} else {
			var cfg, _ = ini.LooseLoad(cfgFile)
			for i, v := range pins {
				if v == r {
					pins = append(pins[:i], pins[i+1:]...)
					break
				}
			}
			for j := range pins {
				if j == len(pins)-1 {
					p = p + pins[j]
				} else {
					p = p + pins[j] + ", "
				}
			}
			_, _ = cfg.Section("general").NewKey("pins", p)
			err = cfg.SaveTo(cfgFile)
		}
	}

	if i == true {
		color.Green("Specified pins:")
		for i := range pins {
			fmt.Println(pins[i])
		}
	}

	if ia == true {

		fileList := []string{}
		fileList, err = getFileList(sd)

		color.Green("Full index of the specified pins and their unique values:")

		for i := range pins {

			var index []string

			for _, file := range fileList {
				fo, _ := os.Open(file)

				scanner := bufio.NewScanner(fo)

				for scanner.Scan() {
					if strings.Contains(scanner.Text(), "* "+pins[i]+":") == true {
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
	}

	if len(l) > 0 {

		fileList := []string{}
		fileList, err = getFileList(sd) //function used in search.go

		var index []string
		var date []string
		var files []string

		for _, file := range fileList {

			var d string

			fo, _ := os.Open(file)

			scanner := bufio.NewScanner(fo)

			for scanner.Scan() {

				if strings.Contains(scanner.Text(), "* date:") == true {
					d = strings.TrimSpace(strings.Replace(scanner.Text(), "* date:", "", -1))
				}

				if strings.Contains(scanner.Text(), "* "+l+":") == true {
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

			color.Green("Pin '" + l + "' not present or never completed in journal entries.")

		} else {

			color.Green("Dated entries for pin '" + l + "':")

			sortutil.CiAsc(index)
			for k := range index {
				fmt.Printf("%v \t %s \t %s \n", date[k], index[k], files[k])
			}
		}
	}

	check(err)
}

func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
