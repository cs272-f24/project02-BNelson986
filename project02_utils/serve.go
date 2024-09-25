package project02utils

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
)

/*
Struct to contain the inverted index. And share between functions
*/
type Server struct {
	invIndex map[string]map[string]int
}

/*
Struct to contain requested search word.
*/
type SearchRequest struct {
	Query string
}

/*
Struct to contain the search query and the list of hits.
*/
type SearchResponse struct {
	Query   string
	Results []string
}

/*
Function to handle displaying the search query.

Receiver:
- s: The Server struct. Used to access invIndex.

Parameters:
- w: The http.ResponseWriter object.
- r: The http.Request object.
*/
func (s *Server) queryForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	qForm, err := os.ReadFile("project02_utils/static/queryForm.html")
	if err != nil {
		fmt.Fprintf(w, "Error reading query form: %v", err)
		return
	}

	w.Write(qForm)
}

/*
Function to handle the search query.

Receiver:
- s: The Server struct. Used to access invIndex.

Parameters:
- w: The http.ResponseWriter object.
- r: The http.Request object.
*/
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve queried word from page
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "No query provided", http.StatusBadRequest)
		return
	}

	results := search(query, s.invIndex)

	// Parse results with ranking into simple URLs
	urls := func(r map[string]int) []string {
		var urls []string
		for k := range r {
			urls = append(urls, k)
		}
		return urls
	}(results)

	// Create response struct
	resp := SearchResponse{
		Query:   query,
		Results: urls,
	}

	// Load search results template
	err := template.New("project02_utils/static/top10.html").Execute(w, resp)
	if err != nil {
		http.Error(w, "Error loading search results template", http.StatusInternalServerError)
		return
	}
}

/*
Function to serve the search page.
*/
func Serve(m Maps) {

	// Create a new server
	server := &Server{
		invIndex: m.invIndex,
	}

	http.HandleFunc("/", server.queryForm)
	http.HandleFunc("/top10/", server.searchHandler)
	go http.ListenAndServe("localhost:8080", nil)
}

/*func main() {
	index := newMaps()
	Serve(index.invIndex)

	for {
		time.Sleep(1 * time.Second)
	}
}*/
