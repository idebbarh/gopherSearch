package searchengine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SearchQuery struct {
	query string
}

type FilesRank = map[string]float64

type ResultType struct {
	Path  string
	Title string
}

type ResponseType struct {
	Result         []ResultType
	IsCompleteData bool
}

type PaginationInfo struct {
	perRequest   int
	currentIndex int
}

func docsResponseHandler(result []ResultType, p *PaginationInfo, w http.ResponseWriter) {
	response := ResponseType{}

	paginationStart := p.currentIndex*p.perRequest - p.perRequest

	paginationEnd := paginationStart + p.perRequest

	if paginationStart >= len(result) {
		paginationStart = max(0, len(result)-1)
	}

	if paginationEnd >= len(result) {
		paginationEnd = len(result)
	}

	response.Result = result[paginationStart:paginationEnd]

	if len(response.Result) == 0 {
		// return empty slice instead of null
		response.Result = make([]ResultType, 0)
	}

	response.IsCompleteData = false

	if paginationEnd < len(result) {
		p.currentIndex += 1
	} else {
		response.IsCompleteData = true
	}

	jsonResponse, marshalErr := json.Marshal(response)

	if marshalErr != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, writeErr := w.Write(jsonResponse)

	if writeErr != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func serveHandler(inMemoryDataChan chan InMemoryData) {
	fmt.Println("serve")
	result := []ResultType{}
	paginationInfo := PaginationInfo{}
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case inMemoryData := <-inMemoryDataChan:
			if r.Method == http.MethodGet && r.URL.Path == "/search" {

				paginationInfo.perRequest = 10
				paginationInfo.currentIndex = 1

				searchQuery := lexer(strings.Join(r.URL.Query()["q"], " "))

				filesRank := FilesRank{}

				var termsIDFValue float64 = 0

				for _, term := range searchQuery {
					termsIDFValue += calcIDF(inMemoryData.Df.Size, inMemoryData.Df.Value, term)
				}

				for f, tf := range inMemoryData.Ftf {
					var rank float64 = 0
					for _, term := range searchQuery {
						termTFValue := calcTF(tf.Terms, term, tf.DocSize)
						rank += termsIDFValue * termTFValue
					}
					if rank == 0 {
						continue
					}

					filesRank[f] = rank
				}

				rankedDocs := rankDocs(filesRank)

				result = nil

				for _, path := range rankedDocs {
					result = append(result, ResultType{Path: path, Title: inMemoryData.Ftf[path].Title})
				}

				docsResponseHandler(result, &paginationInfo, w)
			} else if r.Method == http.MethodGet && r.URL.Path == "/nextSearch" {
				docsResponseHandler(result, &paginationInfo, w)
			} else if r.Method == http.MethodGet && r.URL.Path == "/file" {

				fileToServePath := r.URL.Query().Get("path")
				http.ServeFile(w, r, fileToServePath)

			} else {
				fs.ServeHTTP(w, r)
			}
		default:
			fmt.Println("No data received yet")
		}
	}))

	address := "localhost:8080"

	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("ERROR: Could not serve files: %v\n", err)
	}
	fmt.Println("Listening on :8080...")
}
