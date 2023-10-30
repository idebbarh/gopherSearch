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
	Result []ResultType
}

type PaginationInfo struct {
	perRequest   int
	currentIndex int
}

func (response *ResponseType) setResponseResult(result []ResultType, p *PaginationInfo, w http.ResponseWriter) {
	paginationStart := p.currentIndex*p.perRequest - p.perRequest

	paginationEnd := paginationStart + p.perRequest

	response.Result = result[paginationStart:paginationEnd]

	p.currentIndex += 1

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
	fs := http.FileServer(http.Dir("./static"))

	result := []ResultType{}
	paginationInfo := PaginationInfo{perRequest: 5, currentIndex: 1}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/search" {
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

			response := ResponseType{}

			response.setResponseResult(result, &paginationInfo, w)

		} else if r.Method == http.MethodGet && r.URL.Path == "/file" {
			fileToServePath := r.URL.Query().Get("path")
			http.ServeFile(w, r, fileToServePath)

		} else if r.Method == http.MethodGet && r.URL.Path == "nextSearch" {
			response := ResponseType{}
			response.setResponseResult(result, &paginationInfo, w)

		} else {
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
