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
	m *Maps
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
	Results []results
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

	// Search for the query in the inverted index and get the results
	results := search(query, s.m)

	// Create response struct
	resp := SearchResponse{
		Query:   query,
		Results: results,
	}

	// Load search results template
	tmpl, err := template.ParseFiles("project02_utils/static/top10.html")
	if err != nil {
		http.Error(w, "Error loading search results template", http.StatusInternalServerError)
		return
	}

	// Execute template with response struct
	err = tmpl.Execute(w, resp)
}

/*
Function to serve the search page.
*/
func Serve(m *Maps) {

	// Create a new server
	server := &Server{
		m: m,
	}

	http.HandleFunc("/", server.queryForm)
	http.HandleFunc("/results", server.searchHandler)
	go http.ListenAndServe("localhost:8080", nil)
}

/*func main() {
	index := newMaps()
	Serve(index.invIndex)

	for {
		time.Sleep(1 * time.Second)
	}
}*/
