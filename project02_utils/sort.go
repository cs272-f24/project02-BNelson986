package project02utils

type sort interface {
	Len([]results) int
	Less(r []results, i, j float64) bool
	Sort([]results) []results
	Swap(i, j results)
}

type results struct {
	URL   string
	Score float64
}

/*
Interface implementation for sorting the results of a search query based on the tf-idf score.
*/
func Len(r []results) int {
	return len(r)
}

func Less(r []results, i, j int) bool {
	return r[i].Score < r[j].Score
}

func Swap(r []results, i, j int) {
	r[i], r[j] = r[j], r[i]
}

func Sort(r []results) []results {
	for i := range r {
		for j := range r {
			// Sort by score (descending)
			if !Less(r, i, j) {
				Swap(r, i, j)
			}
			// If the scores are equal, sort by URL
			if r[i].Score == r[j].Score {
				if r[i].URL < r[j].URL {
					Swap(r, i, j)
				}
			}
		}
	}
	return r
}

/*
Function that sorts the results of a search query based on the tf-idf score.

Parameters:
- m: A pointer to a maps struct.
- query: A string representing the search query.

Returns:
- A slice of results structs representing the top 10 search results.
*/
func sortResults(m Maps, query string) []results {
	// Create a map to store the results
	docScores := []results{}

	// Calculate the term frequency-inverse document frequency for each document
	for doc := range m.invIndex[query] {
		// Calculate the tf-idf for the query in the document
		tfidf := tfIdf(m, query, doc)
		docScores = append(docScores, results{URL: doc, Score: tfidf})
	}

	// Sort the results by score
	Sort(docScores)

	// Create a map to store the top 10 results
	if len(docScores) > 10 {
		docScores = docScores[:10]
	}
	return docScores
}
