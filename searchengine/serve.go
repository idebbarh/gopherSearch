package searchengine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func serveHandler(filePath string) {
	result := []ResultType{}
	paginationInfo := PaginationInfo{}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/search" {

			paginationInfo.perRequest = 10
			paginationInfo.currentIndex = 1

			searchQuery := lexer(strings.Join(r.URL.Query()["q"], " "))
			loadedJsonFile, readFileErr := os.ReadFile(filePath)

			if readFileErr != nil {
				http.Error(w, "Failed to open json file", http.StatusInternalServerError)
				return
			}

			var ftf FilesTermsFrequency

			filesRank := FilesRank{}

			json.Unmarshal(loadedJsonFile, &ftf)

			var termsIDFValue float64 = 0

			for _, term := range searchQuery {
				termsIDFValue += calcIDF(ftf, term)
			}

			for f, tf := range ftf {
				var rank float64 = 0
				for _, term := range searchQuery {
					termTFValue := calcTF(tf.Terms, term)
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
				result = append(result, ResultType{Path: path, Title: ftf[path].Title})
			}

			docsResponseHandler(result, &paginationInfo, w)
		} else if r.Method == http.MethodGet && r.URL.Path == "/nextSearch" {
			docsResponseHandler(result, &paginationInfo, w)
		} else if r.Method == http.MethodGet && r.URL.Path == "/file" {
			fileToServePath := r.URL.Query().Get("path")
			http.ServeFile(w, r, fileToServePath)
		} else {
			fs := http.FileServer(http.Dir("./static"))
			fs.ServeHTTP(w, r)

		}
	}))

	address := "localhost:8080"
	fmt.Println("Listening on :8080...")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("ERROR: Could not serve files: %v\n", err)
	}
}
