package main

import (
	"fmt"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

var defLevOpts = levenshtein.DefaultOptions

type LineResult struct {
	LineNum int
	Score int
	Act string
	Scene string
	Character string
	Line string
	Title string
}

// returns [[line number, levenshtein distance]] for a text
func levenshteinSearch(query [][]rune, lines []string, filename string) []LineResult {

	// levThreshold := int(math.Round(float64(len(query)) / 8.0))
	levThreshold := (len(query)-1)*2 + len(query)

	result := []LineResult{}

	title := lines[0]
	if (title == "lucrece" || title == "the phoenix and turtle" || title == "sonnets") {
		return result
	}

	for lineNum, line := range lines {
		words := strings.Split(line, ";")
		if lineNum > 0 && len(words) > 3 {

			score := 0
			start := 0
			for _, word := range words[3:] {
				for _, qs := range query[start:] {
					ln := []rune(word)
					target := qs
					source := ln

					distance := levenshtein.DistanceForStrings(target, source, defLevOpts)
					if (distance < 1) {
						score = score+2
						start++
						if score > 1 {
							score++
						}
					} else if (distance < 2) {
						score++
						start++
					}
				}
			}

			if (score > levThreshold) {
				fmt.Println(score, line)
				result = append(result, LineResult{
					LineNum: lineNum,
					Score: score,
					Act: words[0],
					Scene: words[1],
					Character: words[2],
				})
			}

		}
	}

	return result
}