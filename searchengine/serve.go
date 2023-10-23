package searchengine

import (
	"fmt"
	"net/http"
)

func serveHandler(filePath string) {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	address := "localhost:8080"
	fmt.Println("Listening on :8080...")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("ERROR: Could not serve files: %v\n", err)
	}
}
