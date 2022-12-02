package main

import (
	"math"
	"regexp"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

var defLevOpts = levenshtein.DefaultOptions

// returns [[line number, levenshtein distance]] for a text
func levenshteinSearch(query string, lines []string, filename string) [][2]int {
	qs := []rune(query)
	qsLen := len(qs)

	levThreshold := int(math.Round(float64(qsLen) / 8.0)) + 2

	// a very crude method for screening
	whitespace := regexp.MustCompile(`\s`).MatchString(query)

	result := [][2]int{}

	startCounting := 10000
	if (lines[0] == "lucrece") {
		startCounting = 9
	} else if (lines[0] == "the phoenix and turtle") {
		startCounting = 12
	} else if (lines[0] == "sonnets") {
		startCounting = 10
	}

	gotCharacters := false
	gotALineAfter := false

	for lineNum, line := range lines {
		if lineNum >= startCounting {
			ln := []rune(line)
			lnLen := len(ln)

			diff := lnLen - qsLen
			target := qs
			source := ln

			if diff == 0 {
				distance := levenshtein.DistanceForStrings(target, source, defLevOpts)
				if (distance < levThreshold) && whitespace {
					result = append(result, [2]int{ lineNum, distance })
				}

			} else {
				if diff < 0 {
					// the source is the bigger string
					source = qs
					target = ln
				}

				// proceed (skipping? todo) through the source
				targetLen := len(target)
				distance := 1000
				for i := 0; i < diff; i++ {
					s := source[i:targetLen+i]
					d := levenshtein.DistanceForStrings(target, s, defLevOpts)
					if (!whitespace && i<3 && d<2) {
						i = targetLen
					} else {
						if d < distance {
							distance = d
						}
					}
				}

				if (distance < levThreshold) {
					result = append(result, [2]int{ lineNum, distance })
				}
			}

		} else {
			if line == "characters in the play" {
				gotCharacters = true
			}
			if gotCharacters && (line == "\r" || line == "\n" || line == "") {
				gotALineAfter = true
			}
			if gotCharacters && gotALineAfter {
				startCounting = lineNum
			}
		}
	}

	return result
}