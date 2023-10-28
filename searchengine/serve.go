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

type FileToServeInfo struct {
	FilePath string `josn:"filePath"`
}

func serveHandler(filePath string) {
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/search" {
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
			result := []ResultType{}
			response := ResponseType{}

			for _, path := range rankedDocs {
				result = append(result, ResultType{Path: path, Title: ftf[path].Title})
			}

			response.Result = result

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

			return
		} else if r.Method == http.MethodPost && r.URL.Path == "/file" {
			var fileToServeInfo FileToServeInfo
			if err := json.NewDecoder(r.Body).Decode(&fileToServeInfo); err != nil {
				http.Error(w, "Failed to decode request body", http.StatusInternalServerError)
			}

			fileToServePath := fileToServeInfo.FilePath
			http.ServeFile(w, r, fileToServePath)

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
