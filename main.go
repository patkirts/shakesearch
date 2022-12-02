package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	// "github.com/umahmood/soundex"
)

func main() {
	var wg sync.WaitGroup

	searcher := Searcher{}
	texts, _ := ioutil.ReadDir("./texts")
	searcher.TextNames = make([]string, 0)

	if _, err := os.Stat("data"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("data", os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// TODO error handling
	err := processTexts(texts, &searcher, &wg)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()

	searcher.Load()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]

		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		fmt.Printf("Got a query: %s\n", query[0])

		qs := strings.ToLower(query[0])
		results := searcher.Search(qs)
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}
