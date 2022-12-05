package main

import (
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

var defLevOpts = levenshtein.DefaultOptions

type LineResult struct {
	LineNum int
	Score int
	Act string
	Scene string
	SceneLine string
	Character string
	Line string
	Title string
}

// returns [[line number, levenshtein distance]] for a text
func levenshteinSearch(query [][]rune, lines []string, filename string) []LineResult {

	// levThreshold := int(math.Round(float64(len(query)) / 8.0))
	levThreshold := len(query)

	result := []LineResult{}

	title := lines[0]
	if (title == "lucrece" || title == "the phoenix and turtle" || title == "sonnets") {
		return result
	}

	for lineNum, line := range lines {
		words := strings.Split(line, ";")
		if lineNum > 0 && len(words) > 4 {

			score := 0
			start := 0
			for _, word := range words[4:] {
				trackScore := score
				for _, qs := range query[start:] {
					ln := []rune(word)
					target := qs
					source := ln

					distance := levenshtein.DistanceForStrings(target, source, defLevOpts)
					if (distance < 1) {
						score++
						if (len(target) > 3) {
							score++
						}
					} else if (distance < 2 && len(target) > 3) {
						score++
					}
				}
				if score > (trackScore+1) {
					start++
				}
			}

			if (score > levThreshold) {
				result = append(result, LineResult{
					LineNum: lineNum,
					Score: score,
					Act: words[0],
					Scene: words[1],
					SceneLine: words[2],
					Character: words[3],
				})
			}

		}
	}

	return result
}