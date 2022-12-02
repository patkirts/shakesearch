package main

import (
	"fmt"
	"os"
	"strings"
)

type Searcher struct {
	TextNames    []string
	TextsMap     map[string][]string
	LowercaseMap map[string][]string
	// SoundexMap   map[string][]string
}

type SearchResult struct {
	Title string
	Lines []string
	Rank int
}

type SearchResults []SearchResult

func (s *Searcher) Load() {
	s.TextsMap = make(map[string][]string)
	s.LowercaseMap = make(map[string][]string)
	// s.SoundexMap = make(map[string][]string)

	for _, filename := range s.TextNames {
		plain, err := os.ReadFile("./texts/" + filename)
		if err != nil {
			fmt.Println("Loading file: %w", err)
		}
		s.TextsMap[filename] = strings.Split(string(plain), "\n")

		lower, err := os.ReadFile("./data/lowercase---" + filename)
		if err != nil {
			fmt.Println("Loading file: %w", err)
		}
		s.LowercaseMap[filename] = strings.Split(string(lower), "\n")

		// sound, err := os.ReadFile("./data/soundex---" + filename)
		// if err != nil {
		// 	fmt.Println("Loading file: %w", err)
		// }
		// s.SoundexMap[filename] = strings.Split(string(sound), "\n")
	}
}

func (s *Searcher) Search(query string) SearchResults {
	// for t of texts
	// get line #s matching soundex
	// get line #s matching levenshtein
	// put 2 match lines ahead of 1 match lines
	// get the title and those lines (+ relevant before/after lines)
	// (relevance defined differently for plays, poems, sonnets)
	results := SearchResults{}
	for _, filename := range s.TextNames {
		levRes := levenshteinSearch(query, s.LowercaseMap[filename], filename)
		if len(levRes) > 0 {
			passages := s.GetSearchResults(filename, levRes)
			results = append(results, passages...)
		}
		// soundexSearch(query, s.SoundexMap[filename])
	}
	return results
	// return sort.SliceStable(results, func(i, j int) bool {
	// 	return results[i].Rank < results[j].Rank
	// })
}

func (s *Searcher) GetSearchResults(filename string, results [][2]int) SearchResults {
	text := s.TextsMap[filename]
	passages := SearchResults{}
	resLen := len(results)

	for i := 0; i < resLen; i++{
		lines := []string{}
		r := results[i]
		lineNum := r[0]
		rank := r[1]
		// if lineNum > 0 {
		// 	// prev := TODO think about this
		// }
		lines = append(lines, text[lineNum])

		passages = append(passages, SearchResult{text[0], lines, rank})
	}

	return passages
}

// TODO (own file)
// func soundexSearch(query string, lines []string) {
// 	const soundexQS = soundex.Encode()
// 	for _, l := range lines {
// 		if

// 	}
// }

