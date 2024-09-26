package project02utils

import (
	"math"
)

/*
Function that calculates the term frequency of a term in a document.

Parameters:
- m: A pointer to a maps struct.
- term: A string representing the term to calculate the term frequency for.
- doc: A string representing the document to calculate the term frequency in.

Returns:
- A float64 representing the term frequency of the term in the document.
*/
func termFreq(m Maps, term, doc string) float64 {
	// Get the total number of words in the document
	termInDoc := float64(m.invIndex[term][doc])
	totalWords := float64(m.wordCount[doc])

	// Calculate the term frequency
	return termInDoc / totalWords
}

/*
Function that calculates the inverse document frequency of a term in a collection of documents.

log10(totalDocs / numDocs)

Parameters:
- m: A pointer to a maps struct.
- term: A string representing the term to calculate the inverse document frequency for.

Returns:
- A float64 representing the inverse document frequency of the term in the collection of documents.
*/
func idf(m Maps, term string) float64 {
	// Get the total number of documents
	totalDocs := float64(len(m.wordCount))

	// Get the number of documents containing the term
	numDocs := float64(len(m.invIndex[term]))

	// Calculate the inverse document frequency
	return math.Log10(totalDocs / numDocs)
}

/*
Function that calculates the term frequency-inverse document frequency of a term in a document.

Parameters:
- m: A pointer to a maps struct.
- term: A string representing the term to calculate the tf-idf for.
- doc: A string representing the document to calculate the tf-idf in.

Returns:
- A float64 representing the term frequency-inverse document frequency of the term in the document.
*/
func tfIdf(m Maps, term, doc string) float64 {
	score := termFreq(m, term, doc) * idf(m, term)

	// Truncate the score to 4 decimal places
	return math.Round(score*10000) / 10000
}
