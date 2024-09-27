package project02utils

import (
	"fmt"
	"testing"
)

type testStruct struct {
	testName string
	m        *Maps
	query    string
	expected []results
}

/*
Helper function to compare expected and actual results.
*/
func resultCompare(a, b []results) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].URL != b[i].URL || a[i].Score != b[i].Score {
			return false
		}
	}
	return true
}

/*
TF-IDF Calculation:
    TF = doc(term) / totalWords(doc)
	IDF = log10(totalDocs / numDocs(term))
	TF-IDF = TF * IDF
*/

func TestSortResults(t *testing.T) {
	tests := []testStruct{
		// Basic test case with two documents
		{
			testName: "Basic test case with two documents",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"term1": {"doc1": 3},
					"term2": {"doc2": 4},
				},
				wordCount: map[string]int{
					"doc1": 10,
					"doc2": 20,
				},
			},
			query: "term1",
			expected: []results{
				{"doc1", 0.09030900},
			},
		},
		{
			testName: "Equal TF-IDF scores",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"unique1": {"doc1": 2},
					"unique2": {"doc2": 2},
				},
				wordCount: map[string]int{
					"doc1": 10,
					"doc2": 10,
				},
			},
			query: "unique1",
			expected: []results{
				{"doc1", 0.06020600},
			},
		},
		{
			testName: "More than 10 documents",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"term": {
						"doc1": 1, "doc2": 2, "doc3": 3, "doc4": 4, "doc5": 5,
						"doc6": 6, "doc7": 7, "doc8": 8, "doc9": 9, "doc10": 10,
					},
				},
				wordCount: map[string]int{
					"doc1": 10, "doc2": 10, "doc3": 10, "doc4": 10, "doc5": 10,
					"doc6": 10, "doc7": 10, "doc8": 10, "doc9": 10, "doc10": 10,
					"doc11": 10, "doc12": 10,
				},
			},
			query: "term",
			expected: []results{
				{"doc10", 0.07918125},
				{"doc9", 0.07126312},
				{"doc8", 0.063345},
				{"doc7", 0.05542687},
				{"doc6", 0.04750875},
				{"doc5", 0.03959062},
				{"doc4", 0.0316725},
				{"doc3", 0.02375437},
				{"doc2", 0.01583625},
				{"doc1", 0.00791812},
			},
		},
		{
			testName: "Zero TF-IDF scores",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"term": {"doc1": 2, "doc5": 1, "doc10": 3},
				},
				wordCount: map[string]int{
					"doc1": 20, "doc2": 30, "doc3": 25, "doc4": 15, "doc5": 50,
					"doc6": 40, "doc7": 35, "doc8": 45, "doc9": 55, "doc10": 60,
				},
			},
			query: "term",
			expected: []results{
				{"doc1", 0.05228787},
				{"doc10", 0.02614394},
				{"doc5", 0.01045757},
			},
		},
		{
			testName: "URL/Doc order preservation",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"common": {
						"doc1": 5, "doc2": 5, "doc3": 5, "doc4": 5, "doc5": 5,
						"doc6": 5, "doc7": 5, "doc8": 5, "doc9": 5, "doc10": 5,
					},
				},
				wordCount: map[string]int{
					"doc1": 100, "doc2": 100, "doc3": 100, "doc4": 100, "doc5": 100,
					"doc6": 100, "doc7": 100, "doc8": 100, "doc9": 100, "doc10": 100,
				},
			},
			query: "common",
			expected: []results{
				{"doc1", 0},
				{"doc10", 0},
				{"doc2", 0},
				{"doc3", 0},
				{"doc4", 0},
				{"doc5", 0},
				{"doc6", 0},
				{"doc7", 0},
				{"doc8", 0},
				{"doc9", 0},
			},
		},
		{
			testName: "Low TF-IDF scores",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"rare": {"doc3": 2},
				},
				wordCount: map[string]int{
					"doc1": 100, "doc2": 100, "doc3": 100, "doc4": 100, "doc5": 100,
				},
			},
			query: "rare",
			expected: []results{
				{"doc3", 0.0139794},
			},
		},
		{
			testName: "Sort by URL",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"term": {"doc1": 1, "doc2": 1, "doc3": 1, "doc4": 1},
				},
				wordCount: map[string]int{
					"doc1": 100, "doc2": 100, "doc3": 100, "doc4": 100,
				},
			},
			query: "term",
			expected: []results{
				{"doc1", 0},
				{"doc2", 0},
				{"doc3", 0},
				{"doc4", 0},
			},
		},
		{
			testName: "Empty query",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"term1": {"doc1": 2},
				},
				wordCount: map[string]int{
					"doc1": 50, "doc2": 50,
				},
			},
			query:    "",
			expected: []results{},
		},
	}

	for i, test := range tests {
		testName := fmt.Sprint("Test case ", i)
		t.Run(testName, func(t *testing.T) {
			actual := sortResults(*test.m, test.query)
			if !resultCompare(actual, test.expected) {
				t.Errorf("Test case %d failed: expected %v, got %v", i, test.expected, actual)
			}
		})
	}
}
