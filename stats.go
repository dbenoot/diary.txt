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
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func statistics(wd string, t string) {

	// variables

	var entries, twc, y, year, years, m, month, months, tags, uniquetags int
	var err error
	var f []string
	var superstring string
	var pintags = regexp.MustCompile(`\*\s\w+:`)
	var wc map[string]int
	tagsMap := make(map[string]int)

	of := filepath.Join(wd, "logs", "statistics_"+t+".txt")

	fo, err := os.Create(of)
	check(err)
	defer fo.Close()

	// populate variables

	f, err = getFileList(wd)

	// get amount of journal entries

	entries = len(f)

	// extract other variables from within the files

	for _, file := range f {

		fo, _ := os.Open(file)

		scanner := bufio.NewScanner(fo)

		for scanner.Scan() {

			// add text in the superstring

			if !pintags.MatchString(scanner.Text()) {
				superstring = superstring + scanner.Text() + "\n"
			}

			// date operations

			if !strings.Contains(scanner.Text(), "* date:") {

				// get the amount of years

				y, err = strconv.Atoi(getYear(scanner.Text()))
				check(err)

				if year != y {
					years++
				}

				year = y

				// get the amount of months

				m, err = strconv.Atoi(getMonth(scanner.Text()))

				if month != m {
					months++
				}

				month = m

			}

			// tag operations: total amount of tags

			if strings.Contains(scanner.Text(), "* tags:") {
				t := strings.TrimSpace(strings.Replace(scanner.Text(), "* tags:", "", -1))
				if len(t) != 0 {
					ts := strings.Split(t, ",")
					tags = tags + len(ts)

					for _, tt := range ts {
						tt = strings.TrimSpace(tt)
						if tagsMap[tt] == 0 {
							uniquetags++
						}
						tagsMap[tt]++
					}

				}
			}

		}
	}

	wc, twc = wordcount(superstring)

	// output

	fmt.Printf("Statistics log created: %s \n", of)

	fo.WriteString("Amount of years writing journal:   " + strconv.Itoa(years) + "\n")
	fo.WriteString("Amount of months writing journal:  " + strconv.Itoa(months) + "\n")
	fo.WriteString("Total amount of journal entries:   " + strconv.Itoa(entries) + "\n")
	fo.WriteString("Averge entries per year:           " + strconv.Itoa(entries/years) + "\n")
	fo.WriteString("Average entries per month:         " + strconv.Itoa(entries/months) + "\n")
	fo.WriteString("Amount of unique tags used:        " + strconv.Itoa(uniquetags) + "\n")
	fo.WriteString("Total amount of tags used:         " + strconv.Itoa(tags) + "\n")
	fo.WriteString("Average amount of tags per entry:  " + strconv.FormatFloat(float64(tags)/float64(entries), 'f', 3, 64) + "\n")
	fo.WriteString("Total word count:                  " + strconv.Itoa(twc) + "\n\n")

	fo.WriteString("Tab-separated list of amount of times a word is used: \n\n")
	fo.WriteString("word\t#\tword density\n")
	rankWC := rankByWordCount(wc)
	for _, wcp := range rankWC {
		density := float64(wcp.Value) / float64(twc) * 100
		fo.WriteString(wcp.Key + "\t" + strconv.Itoa(wcp.Value) + "\t" + strconv.FormatFloat(density, 'f', 3, 64) + "%" + "\n")
	}

	check(err)

}

func wordcount(s string) (map[string]int, int) {

	totalWordCount := 0

	unwantedSymbols := [...]string{":", "-", "_", "(", ")", "'", "*", ",", ".", ":", "###", "!", ";", "‘", "'", ")", "?", "(", "\"", "?)", ":"}

	substrs := strings.Fields(s)

	for a, str := range substrs {

		substrs[a] = strings.ToLower(substrs[a])

		// remove unicode symbols

		for _, letter := range str {
			if unicode.IsSymbol(letter) {
				substrs[a] = strings.Replace(str, string(letter), "", -1)
			}

		}

		// remove unwanted symbols

		for _, us := range unwantedSymbols {
			if strings.ContainsAny(substrs[a], us) {
				substrs[a] = strings.Replace(substrs[a], us, "", -1)
			}
		}

		totalWordCount++
	}

	// key,value  ==> word,count

	counts := make(map[string]int)

	for _, word := range substrs {
		if len(word) > 0 {
			counts[word]++
		}
	}

	return counts, totalWordCount
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
