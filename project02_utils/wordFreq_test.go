package project02utils

import (
	"math"
	"testing"
)

func TestTermFreq(t *testing.T) {
	// Create local testing invIndex and WordCount maps
	m := Maps{
		invIndex: map[string]map[string]int{
			"term1": {"doc1": 3, "doc2": 2},
			"term2": {"doc1": 1},
		},
		wordCount: map[string]int{
			"doc1": 10,
			"doc2": 5,
		},
	}

	// Test cases
	tests := []struct {
		term     string
		doc      string
		expected float64
	}{
		{"term1", "doc1", 0.3},
		{"term1", "doc2", 0.4},
		{"term2", "doc1", 0.1},
	}

	for _, test := range tests {
		result := termFreq(m, test.term, test.doc)
		if result != test.expected {
			t.Errorf("termFreq(%v, %v, %v) = %v; want %v", m, test.term, test.doc, result, test.expected)
		}
	}
}

func TestIdf(t *testing.T) {
	// Create local testing invIndex and WordCount maps
	m := Maps{
		invIndex: map[string]map[string]int{
			"term1": {"doc1": 3, "doc2": 2},
			"term2": {"doc1": 1},
		},
		wordCount: map[string]int{
			"doc1": 10,
			"doc2": 5,
		},
	}

	// Test cases
	tests := []struct {
		term     string
		expected float64
	}{
		{"term1", math.Log10(2.0 / 2.0)},
		{"term2", math.Log10(2.0 / 1.0)},
	}

	for _, test := range tests {
		result := idf(m, test.term)
		if result != test.expected {
			t.Errorf("idf(%v, %v) = %v; want %v", m, test.term, result, test.expected)
		}
	}
}

func TestTfIdf(t *testing.T) {
	// Create local testing invIndex and WordCount maps
	m := Maps{
		invIndex: map[string]map[string]int{
			"term1": {"doc1": 3, "doc2": 2},
			"term2": {"doc1": 1},
		},
		wordCount: map[string]int{
			"doc1": 10,
			"doc2": 5,
		},
	}

	// Test cases
	tests := []struct {
		term     string
		doc      string
		expected float64
	}{
		{"term1", "doc1", math.Round(0.3*math.Log10(2.0/2.0)*10000) / 10000},
		{"term1", "doc2", math.Round(0.4*math.Log10(2.0/2.0)*10000) / 10000},
		{"term2", "doc1", math.Round(0.1*math.Log10(2.0/1.0)*10000) / 10000},
	}

	for _, test := range tests {
		result := tfIdf(m, test.term, test.doc)
		if result != test.expected {
			t.Errorf("tfidf(%v, %v, %v) = %v; want %v", m, test.term, test.doc, result, test.expected)
		}
	}
}
