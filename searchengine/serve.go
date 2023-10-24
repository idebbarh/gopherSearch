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

			fmt.Println(termsIDFValue)
			for f, tf := range ftf {
				var rank float64 = 0
				for _, term := range searchQuery {
					termTFValue := calcTF(tf, term)
					rank += termsIDFValue * termTFValue
				}
				filesRank[f] = rank
			}

			filesRank = sortMap(filesRank, 10)

			for key, value := range filesRank {
				fmt.Printf("%s\t%f\n", key, value)
			}

			jsonResponse, marshalErr := json.Marshal(filesRank)

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
		}

		fs.ServeHTTP(w, r)
	}))

	address := "localhost:8080"
	fmt.Println("Listening on :8080...")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("ERROR: Could not serve files: %v\n", err)
	}
}
