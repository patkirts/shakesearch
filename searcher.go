package main

import (
	"fmt"
	"os"
	"strings"
)

type Searcher struct {
	TextNames    []string
	TextsMap     map[string][]string
	SearchMap map[string][]string
	// SoundexMap   map[string][]string
}

type SearchResults []LineResult

func (s *Searcher) Load() {
	s.TextsMap = make(map[string][]string)
	s.SearchMap = make(map[string][]string)
	// s.SoundexMap = make(map[string][]string)

	for _, filename := range s.TextNames {
		plain, err := os.ReadFile("./data/main---" + filename)
		if err != nil {
			fmt.Println("Loading file: %w", err)
		}
		s.TextsMap[filename] = strings.Split(string(plain), "\n")

		search, err := os.ReadFile("./data/search---" + filename)
		if err != nil {
			fmt.Println("Loading file: %w", err)
		}
		s.SearchMap[filename] = strings.Split(string(search), "\n")

		// sound, err := os.ReadFile("./data/soundex---" + filename)
		// if err != nil {
		// 	fmt.Println("Loading file: %w", err)
		// }
		// s.SoundexMap[filename] = strings.Split(string(sound), "\n")
	}
}

func (s *Searcher) Search(query [][]rune) SearchResults {
	// for t of texts
	// get line #s matching soundex
	// get line #s matching levenshtein
	// put 2 match lines ahead of 1 match lines
	// get the title and those lines (+ relevant before/after lines)
	// (relevance defined differently for plays, poems, sonnets)
	results := SearchResults{}
	for _, filename := range s.TextNames {
		levRes := levenshteinSearch(query, s.SearchMap[filename], filename)
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

func (s *Searcher) GetSearchResults(filename string, results []LineResult) SearchResults {
	text := s.TextsMap[filename]
	passages := SearchResults{}

	for _, res := range results {
		res.Title = text[0]
		res.Line = text[res.LineNum]
		passages = append(passages, res)
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

