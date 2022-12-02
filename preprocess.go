package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
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
			go processFile(filename, s, wg)
		}
	}

	return nil
}

func processFile(filename string, s *Searcher, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open("./texts/" + filename)
	if err != nil {
		fmt.Println("Loading file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lowercaseFile, err := os.Create("data/lowercase---"+filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// soundexFile, err := os.Create("data/soundex---"+filename)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		// soundexLine := encodeSoundexLine(line)
		fmt.Fprintln(lowercaseFile, line)
		// fmt.Fprintln(soundexFile, soundexLine)
	}

	err = lowercaseFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// err = soundexFile.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

// func encodeSoundexLine(line string) string {
// 	soundexLine := ""
// 	words := strings.Split(line, " ")
// 	for _, word := range words {
// 		soundexLine += soundex.Code(word)
// 	}
// 	return soundexLine
// }