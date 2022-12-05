package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// These functions create the initial files for us
// that Searcher Loads() into data arrays for us

func processTexts(texts []fs.FileInfo, s *Searcher, wg *sync.WaitGroup) error {
	fmt.Println("Processing each of Shakespeare's texts for searchability...")

	for _, text := range texts {
		if !text.IsDir() {
			wg.Add(1)
			filename := text.Name()
			s.TextNames = append(s.TextNames, filename)
			// go processFile(filename, s, wg)
			go processFile(filename, s, wg)
		}
	}

	return nil
}

// types of lines to distinguish
// PLAYS
// - title
// - other meta
// - eliminate non A-Za-z rows
// - Characters in the Play &ff which stops at next blank line
// - then play text begins
// - Do we have a Prologue or Act 1?
// - Act/Scene #
// - stage direction line(s): text btwn []
// - CHARACTER NAME-only lines
// - CHARACTER NAME-starting lines
//
// information we want to retrieve, per line
// - Act, Scene, Speaker, line number
// - prior line if not first line in scene
// - next line if not the last line
// - index(es) range(s) that match
//
// Options for play-search: TODO
// - Filter by specific characters
// - Limit works to be searched
// - We want to exclude stage directions from search by default?
//
// Open questions: TODO
// - character parsing (can I ask the server whether a word is a character's name?)
//
// Things we want: TODO
// - "fetch more" of the text matching the particular lines
//
// SONNETS TOOD
// want to return Sonnet # and specific lineNum

// ranking lines...

var actRegex = regexp.MustCompile("A[C|c][T|t] [1-5]")
var sceneRegex = regexp.MustCompile("S[C|c][E|e][N|n][E|e] [1-9]")
var eqRegex = regexp.MustCompile("=")
var charLineRegex = regexp.MustCompile("^[A-Z][A-Z][A-Z ]* ?[A-Z][A-Z]")
var inlineStageRegex = regexp.MustCompile(`\[([A-Za-z ,\.;:{}])*\]`)
var startStageRegex = regexp.MustCompile(`^\[([A-Za-z ,\.;:{}])*[^\]]$`)
var endStageRegex = regexp.MustCompile(`^([A-Za-z ,\.;:{}])*\]$`)
var puncRegex = regexp.MustCompile(`[.,";:!?()]`)

type Line struct {
	Act int
	Scene int
	Character string
	Words []string
}

func processFile(filename string, s *Searcher, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open("./texts/" + filename)
	if err != nil {
		fmt.Println("Loading file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	mainFile, err := os.Create("data/main---"+filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	searchFile, err := os.Create("data/search---"+filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// soundexFile, err := os.Create("data/soundex---"+filename)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	readingCharacters := false
	throughCharacters := false
	stageDirections := false
	act := 0 // prologues implicitly tracked
	scene := 0
	// stageInstructions := false
	// lines := make([]Line, 1)
	title := ""
	tryline := ""
	character := ""

	for  scanner.Scan() {
		line := scanner.Text()
		words := make([]string, 0)

		if title == "" {
			title = line
			fmt.Fprintln(mainFile, title)
			fmt.Fprintln(searchFile, strings.ToLower(strings.TrimSpace(title)))
		} else if line == "Characters in the Play" {
			readingCharacters = true
		} else if !readingCharacters {
			// we're not doing anything with the metadata yet

		} else if !throughCharacters {
			if line == "\r" || line == "" {
				throughCharacters = true
			} else {
				// potential TODO parse characters amidst
			}

		} else if line == "" || line == "\r" || eqRegex.MatchString(line) {
			// best way to check for empty string in go? TODO

		} else if stageDirections {
			if endStageRegex.MatchString(line) {
				stageDirections = false
			}

		} else if startStageRegex.MatchString(line) {
			stageDirections = true

		} else if actRegex.MatchString(line) {
			n, err := strconv.ParseInt(strings.Split(line, " ")[1], 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			act = int(n)

		} else if sceneRegex.MatchString(line) {
			n, err := strconv.ParseInt(strings.Split(line, " ")[1], 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			scene = int(n)

		} else if charLineRegex.MatchString(line) {
			nameLen := charLineRegex.FindStringIndex(line)[1]
			character = line[0:nameLen]
			remainder := line[nameLen:]
			if len(remainder) > 2 && remainder[0:2] == ", " {
				remainder = remainder[2:]
			}
			if remainder != "" {
				tryline = remainder
			}
		} else {
			tryline = line
		}

		if tryline != "" {
			if (inlineStageRegex.MatchString(tryline)) {
				idx := inlineStageRegex.FindStringIndex(tryline)
				tryline = tryline[0:idx[0]] + tryline[idx[1]:]
			}
			if tryline != "" {
				for _, word := range strings.Split(tryline, " ") {
					if (word != "") {
						words = append(words, strings.ToLower(puncRegex.ReplaceAllString(strings.TrimSpace(word), "")))
					}
				}
				// fmt.Printf("%d;%d;%s;%s\n", act, scene, character, strings.Join(words, ";"))
				fmt.Fprintf(searchFile, "%d;%d;%s;%s\n", act, scene, character, strings.Join(words, ";"))
				fmt.Fprintln(mainFile, line)
			}
			tryline = ""
		}

		// soundexLine := encodeSoundexLine(line)
		// fmt.Fprintln(soundexFile, soundexLine)
	}

	err = mainFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = searchFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// func encodeSoundexLine(line string) string {
// 	soundexLine := ""
// 	words := strings.Split(line, " ")
// 	for _, word := range words {
// 		soundexLine += soundex.Code(word)
// 	}
// 	return soundexLine
// }