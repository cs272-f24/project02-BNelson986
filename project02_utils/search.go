package project02utils

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
	"golang.org/x/net/html"
)

/*
Function that takes a byte slice representing the body of an HTML document and returns two slices:
one containing all the words found in the text nodes of the HTML document, and another containing
all the href values of the anchor tags in the HTML document.

Parameters:
- body: A byte slice representing the body of an HTML document.

Returns:
- A slice of strings containing all the words found in the text nodes of the HTML document.
- A slice of strings containing all the href values of the anchor tags in the HTML document.
*/
func extract(body []byte) ([]string, []string) {
	/*
		body =
		 	<!DOCTYPE html>
				<html>
					<head>
						<title>CS272 | Welcome</title>
					</head>
					<body>
						<p>Hello World!</p>
						<p>Welcome to <a href="https://cs272-f24.github.io/">CS272</a>!</p>
					</body>
				</html>
	*/

	words := []string{}
	hrefs := []string{}

	z, err := html.Parse(bytes.NewReader(body))

	if err != nil {
		log.Fatalf("Error parsing html. Error: %v", err)
	}

	// Recursively traverse the HTML document and extract words and href values
	var f func(*html.Node)

	f = func(n *html.Node) {
		// Extract href values from anchor tags
		if n.Type == html.ElementNode && n.Data == "a" {
			// Iterate through attributes of anchor tag and extract href values
			for _, a := range n.Attr {
				if a.Key == "href" {
					hrefs = append(hrefs, a.Val)
					break
				}
			}
		}
		// Extract words from text nodes
		if n.Type == html.TextNode {

			// Iterate through text node and split it into words
			// function performs split and cleaning before any further processing occurs
			for _, word := range strings.FieldsFunc(n.Data, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			}) {
				words = append(words, word)
			}
		}
		// Recursively traverse the children of the current node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	// Call the recursive function on the root node of the HTML document
	f(z)

	return words, hrefs
}

/*
Function that takes a host and an href value and returns a cleaned URL.
If the href value is an absolute URL, it is returned as is.

Parameters:
- host: A string representing the host of the URL.
- href: A string representing the href value to clean.

Returns:
- A string representing the cleaned URL.
*/
func clean(host string, href string) string {

	u, err := url.Parse(host)
	if err != nil {
		log.Printf("Error parsing URL: %v\n", err)

	}

	// If href is an absolute URL, return it as is
	if strings.HasPrefix(href, "http") {
		return href
	}

	// Separate parsing of href to comply with test case 4
	// Example href: "/?page=1"
	// Wrong output: "https://test.io/%3Fpage=1"
	// Correct output: "https://test.io/?page=1"
	hrefURL, err := url.Parse(href)

	if err != nil {
		log.Printf("Error parsing URL: %v\n", err)

	}

	return u.ResolveReference(hrefURL).String()
}

/*
Function that takes a URL and returns the body of the webpage at that URL.

Parameters:
- url: A string representing the URL of the webpage to download.

Returns:
- A byte slice representing the body of the webpage.
- An error if the webpage could not be downloaded.
*/
func download(url string) (body []byte, error error) {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		// Returns the body of the webpage
		return io.ReadAll(resp.Body)

	} else {
		// Bad response returns error code
		return nil, err
	}
}

/*
Function that takes a word and returns the stemmed version of the word.

Parameters:
- word: A string representing the word to stem.

Returns:
- A string representing the stemmed version of the word.
*/
func stem(word string) string {
	stemmed, err := snowball.Stem(word, "english", true)
	if err != nil {
		log.Fatalf("Error stemming word: %v", err)
	}
	return stemmed
}

/*
Function that takes a URL and crawls the webpage at that URL, extracting words and href values.
Utilizes download, extract, and clean functions.

Parameters:
- url: A string representing the URL of the webpage to crawl.
*/
func Crawl(m *Maps) map[string]map[string]int {

	// URL to start crawling from
	serverStart := "https://cs272-f24.github.io/top10/"

	// Initialize the needed structs
	downloadQueue := []string{}

	// Add the starting URL to the download queue
	downloadQueue = append(downloadQueue, serverStart)
	m.visited[serverStart] = struct{}{}

	// Extract the origin domain from the starting URL
	origin, err := url.Parse(serverStart)

	if err != nil {
		log.Printf("Error parsing URL: %v\n", err)
	}

	// Continue crawling while the download queue is not empty
	for len(downloadQueue) > 0 {
		// Pop the first URL from the download queue
		nextUp := downloadQueue[0]
		downloadQueue = downloadQueue[1:]

		newLink, err := url.Parse(nextUp)
		if err != nil {
			log.Printf("Error parsing URL: %v\n", err)
		}

		// Check if the URL is from the same domain as the origin domain
		if origin.Host != newLink.Host {
			continue
		}

		// Download the webpage at the given URL
		if body, err := download(nextUp); err == nil {
			// Extract words and href values from the downloaded webpage
			words, hrefs := extract(body)

			// Clean the href values and add them to the download queue
			for _, href := range hrefs {
				newURL := clean(nextUp, href)

				// Add the cleaned URL to the download queue if it has not been visited
				if _, ok := m.visited[newURL]; !ok {
					m.visited[newURL] = struct{}{}
					downloadQueue = append(downloadQueue, newURL)
				}
			}

			// Remove stop words from the extracted words
			words = m.removeStopWords(words)

			// Add the words to the inverted index
			for _, word := range words {

				stemmedWord := stem(word)
				// Check if the word is already in the inverted index
				// Make a new entry if it is not
				if m.invIndex[stemmedWord] == nil {
					m.invIndex[stemmedWord] = make(map[string]int)
				}
				// Increment the frequency of the word in the URL
				m.invIndex[stemmedWord][nextUp]++

				// Increment the total word count of the URL
				m.wordCount[nextUp]++
			}
		} else {
			log.Printf("Error downloading webpage: %v\n", err)
		}
	}

	return nil
}

/*
Function that takes a word and an inverted index and returns a map of URLs and their frequencies
where the word appears.

Parameters:
- word: A string representing the word to search for.
- invIndex: A map representing the inverted index.

Returns:
- A map of strings to ints showing frequencies of the word in URLs.
*/
func search(word string, invIndex map[string]map[string]int) map[string]int {
	// Stem the word
	stemmedWord, err := snowball.Stem(word, "english", true)
	if err != nil {
		log.Fatalf("Error stemming word: %v", err)
	}

	// Check if the word is in the inverted index
	if invIndex[stemmedWord] != nil {
		return invIndex[stemmedWord]
	}
	return nil
}
