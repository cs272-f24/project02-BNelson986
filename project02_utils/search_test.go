package project02utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"
	"time"
)

var testDir string = "TestExtract_Cases/"

func TestExtract(t *testing.T) {
	tests := []struct {
		/*
			Fields:
				body: The HTML document to extract words and href values from
				words: The expected words extracted from the HTML document
				hrefs: The expected href values extracted from the HTML document
		*/
		testFile     string
		words, hrefs []string
	}{
		// Test Case 1
		{
			testDir + "case_1.html",
			[]string{
				"Test", "Case", "Test", "Case", "This", "is", "a", "sample",
				"paragraph", "Content", "Section", "This", "is", "the", "content",
				"section", "Example", "Link", "Footer", "section",
			},
			[]string{
				"https://example.com",
			},
		},
		// Test Case 2
		{
			testDir + "case_2.html",
			[]string{},
			[]string{},
		},
		// Test Case 3
		{
			testDir + "case_3.html",
			[]string{
				"Test", "Extract", "Function", "Test", "Case", "3", "Example",
				"Link", "GitHub", "Link", "Stack", "Overflow", "Link",
			},
			[]string{
				"https://example.com",
				"https://github.com",
				"https://stackoverflow.com",
			},
		},
		{
			testDir + "case_4.html",
			[]string{
				"CS272", "Welcome", "Hello",
				"World", "Welcome", "to", "CS272",
			},
			[]string{
				"https://cs272-f24.github.io/",
			},
		},
	}

	for _, test := range tests {
		if body, err := os.ReadFile(test.testFile); err == nil {
			words, hrefs := extract(body)
			if !slices.Equal(words, test.words) {
				t.Errorf("Words do not match expected: \n%v \n!=\n%v", words, test.words)
			}
			if !slices.Equal(hrefs, test.hrefs) {
				t.Errorf("Hrefs do not match expected: \n%v \n!=\n%v", hrefs, test.hrefs)
			}
		} else {
			t.Errorf("Error reading test file: %v\nError: %v", test.testFile, err)
		}
	}
}

func TestCleanHref(t *testing.T) {
	/*
		Fields:
			host: The host of the website
			hrefs: The href values to clean
			expected: The expected cleaned href values
	*/
	tests := []struct {
		host     string
		hrefs    []string
		expected []string
	}{
		{ // TEST CASE 1
			"https://cs272-f24.github.io/",
			[]string{
				"/",
				"/help/",
				"/syllabus/",
				"https://gobyexample.com/",
			},
			[]string{
				"https://cs272-f24.github.io/",
				"https://cs272-f24.github.io/help/",
				"https://cs272-f24.github.io/syllabus/",
				"https://gobyexample.com/",
			},
		},
		{ // TEST CASE 2
			"https://example.com/",
			[]string{
				"/",
				"/help/faq/",
				"https://example.com/about/",
				"/contact/",
				"",
			},
			[]string{
				"https://example.com/",
				"https://example.com/help/faq/",
				"https://example.com/about/",
				"https://example.com/contact/",
				"https://example.com/",
			},
		},
		{ // TEST CASE 3
			"https://test.io/",
			[]string{
				"/",
				"/products/recent/",
				"about/",
				"/company/law_suits/",
				"/?page=1",
			},
			[]string{
				"https://test.io/",
				"https://test.io/products/recent/",
				"https://test.io/about/",
				"https://test.io/company/law_suits/",
				"https://test.io/?page=1",
			},
		},
	}
	{
		// Execute all test cases
		for _, test := range tests {
			for i, href := range test.hrefs {
				cleaned := clean(test.host, href)
				if cleaned != test.expected[i] {
					t.Errorf("Expected: %v\nGot: %v", test.expected[i], cleaned)
				}
			}
		}
	}
}

func TestStem(t *testing.T) {
	/*
		Fields:
			word: The word to stem
			expected: The expected stemmed word
	*/
	tests := []struct {
		word     string
		expected string
	}{
		{ // TEST CASE 1
			"running",
			"run",
		},
		{ // TEST CASE 2
			"testing",
			"test",
		},
		{ // TEST CASE 3
			"proposition",
			"proposit",
		},
		{ // TEST CASE 4
			"make",
			"make",
		},
		{ // TEST CASE 5
			"Mastery",
			"masteri",
		},
	}
	{
		// Execute all test cases
		for _, test := range tests {
			stemmed := stem(test.word)
			if stemmed != test.expected {
				t.Errorf("Expected: %v\nGot: %v", test.expected, stemmed)
			}
		}
	}
}

func TestDownload(t *testing.T) {
	// Test Case 1: Bad URL
	t.Run("Test Case 1: Bad URL", func(t *testing.T) {
		if _, err := download("https://bad.domain"); err != nil {
			t.Logf("SUCCESS :)\tError downloading bad URL: %v", err)
		} else {
			t.Errorf("FAIL :(\tExpected error downloading bad URL")
		}
	})

	// Test Case 2: Network Timeout
	t.Run("Test Case 2: Network Timeout", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use sleep to simulate network timeout
			time.Sleep(10 * time.Second)

			w.Write([]byte(`<html><body><p>Test Case</p></body></html>`))
		})
		server := httptest.NewServer(handler)
		defer server.Close()

		_, err := download(server.URL)
		if err == nil || strings.Contains(err.Error(), "ReadTimeout") {
			t.Logf("SUCCESS :)\tNetwork timeout error: %v", err)
		} else {
			t.Errorf("FAIL :(\tExpected network timeout error")
		}
	})

	// Test Case 3a: Successful Download
	t.Run("Test Case 3a: Successful Download", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(
				`<html><body><a href="https://example.com">Example Link</a>
				<p>Test Case</p></body></html>`))
		})
		server := httptest.NewServer(handler)
		defer server.Close()

		// Download the webpage and compare the body with the expected body
		if body, err := download(server.URL); err == nil {
			expectedBody := []byte(
				`<html><body><a href="https://example.com">Example Link</a>
				<p>Test Case</p></body></html>`)
			if bytes.Equal(body, expectedBody) {
				t.Logf("SUCCESS :)\tDownloaded: %v", string(body))
			} else {
				t.Errorf("FAIL :(\tExpected: %v\nGot: %v", string(expectedBody), string(body))
			}
		} else {
			t.Errorf("FAIL :(\tError downloading: %v", err)
		}
	})

	// Test Case 3b: Successful Download
	t.Run("Test Case 3b: Successful Download", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html><body><p>This is test case 2</p></body></html>`))
		})
		server := httptest.NewServer(handler)
		defer server.Close()

		// Download the webpage and compare the body with the expected body
		if body, err := download(server.URL); err == nil {
			expectedBody := []byte(`<html><body><p>This is test case 2</p></body></html>`)
			if bytes.Equal(body, expectedBody) {
				t.Logf("SUCCESS :)\tDownloaded: %v", string(body))
			} else {
				t.Errorf("FAIL :(\tExpected: %v\nGot: %v", string(expectedBody), string(body))
			}
		} else {
			t.Errorf("FAIL :(\tError downloading: %v", err)
		}
	})

	// Test Case 4: Success and Extract
	t.Run("Test Case 4: Success and Extract", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(
				`<html><body><a href="https://example.com">Example Link</a>
				<p>Test Case</p></body></html>`))
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		expectedWords := []string{"Example", "Link", "Test", "Case"}
		expectedHrefs := []string{"https://example.com"}

		// Download the webpage and extract words and href values
		// Compare the extracted values with the expected values
		if body, err := download(server.URL); err == nil {
			words, hrefs := extract(body)
			if slices.Equal(words, expectedWords) &&
				slices.Equal(hrefs, expectedHrefs) {
				t.Logf("SUCCESS :)\tExtracted words: %v Extracted hrefs: %v", words, hrefs)
			} else {
				t.Errorf("FAIL :(\nExpected words: %v\nGot: %v", expectedWords, words)
				t.Errorf("Expected hrefs: %v\nGot: %v", expectedHrefs, hrefs)
			}
		} else {
			t.Errorf("FAIL :(\tError downloading: %v", err)
		}

	})
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		name           string
		word           string
		m              *Maps
		expectedResult []string
	}{
		{
			name: "Word found in index",
			word: "Test",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"test": {
						"https://example.com/page1": 3,
						"https://example.com/page2": 1,
					},
				},
			},
			expectedResult: []string{
				"https://example.com/page1",
				"https://example.com/page2",
			},
		},
		{
			name: "Word not found in index",
			word: "Nonexistent",
			m: &Maps{
				invIndex: map[string]map[string]int{
					"test": {
						"https://example.com/page1": 3,
					},
				},
			},
			expectedResult: nil,
		},
		{
			name:           "Empty inverted index",
			word:           "test",
			m:              &Maps{invIndex: map[string]map[string]int{}},
			expectedResult: nil,
		},
	}

	// Execute all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run the search function and compare the result with the expected result
			result := search(tc.word, tc.m)
			// Same as 'resultCompare' function in sort_test.go
			if !func(a, b []string) bool {
				if len(a) != len(b) {
					return false
				}
				for i := range a {
					if a[i] != b[i] {
						return false
					}
				}
				return true
			}(result, tc.expectedResult) {
				t.Errorf("expected %v, but got %v", tc.expectedResult, result)
			}
		})
	}
}
