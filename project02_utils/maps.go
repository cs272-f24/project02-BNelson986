package project02utils

import (
	"bufio"
	"log"
	"os"
)

/*
m represents a data structure for managing an inverted index.

Fields:
  - invIndex: A map where the key is a word (string) and the value is another map. The inner map's key is
    a URL (string) and the value is the count of occurrences (int) of the word in that document.
  - visited: A map where the key is a URL (string) and the value is an empty struct,
    representing the collection of visited documents.
  - wordCount: A map where the key is a URL (string) and the value is the total count of
    words (int) in that document.
  - stopWords: A map where the key is a word (string) and the value is an empty struct, representing the
    collection of stop words.
*/
type Maps struct {
	invIndex  map[string]map[string]int
	visited   map[string]struct{}
	wordCount map[string]int
	stopWords map[string]struct{}
}

/*
Function that creates and returns a new instance of the maps struct.
Gathers all stopwords at initialization.

Returns:
- A pointer to a new instance of the maps struct containing four initialized maps:
  - invIndex: a map where keys are words and values are maps of string URLs to integers.
  - visited: a map where keys are string URLs and values are empty structs.
  - wordCount: a map where keys are string URLs and values number of words.
  - stopWords: a map where keys are words and values are empty structs.
*/
func NewMaps() *Maps {
	m := &Maps{
		invIndex:  make(map[string]map[string]int),
		visited:   make(map[string]struct{}),
		wordCount: make(map[string]int),
		stopWords: make(map[string]struct{}),
	}
	getStopWords(m)

	return m
}

/*
Function that reads a list of stop words from a file and adds them to the stopWords map in the maps struct.

Parameters:
- m: A pointer to a maps struct.
*/
func getStopWords(m *Maps) {
	fileName := "project02_utils/stopwords-en.txt"

	if file, err := os.Open(fileName); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			m.stopWords[scanner.Text()] = struct{}{}
		}
	} else {
		log.Fatalf("Error opening file: %v", err)
	}
}

func (m Maps) removeStopWords(words []string) []string {
	var cleanedWords []string
	for _, word := range words {
		if _, ok := m.stopWords[word]; !ok {
			cleanedWords = append(cleanedWords, word)
		}
	}
	return cleanedWords
}
