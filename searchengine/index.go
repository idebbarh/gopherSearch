package searchengine

import (
	"fmt"
	"sync"
	"time"
)

type FileData struct {
	// terms ferq in doc
	Terms TermsFrequency
	// doc title
	Title string
	// number of terms in the doc
	DocSize int
	// last update time
	LastUpdateTime time.Time
}

func indexHandler(curPath string, inMemoryData *InMemoryData, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, f := range getPathFiles(curPath) {
		v, ok := inMemoryData.Ftf[f.filePath]

		if ok {
			if v.LastUpdateTime.Equal(f.lastUpdateTime) {
				fmt.Printf("ignoring %s....\n", f.filePath)
				continue
			}
		}

		fmt.Printf("indexing %s....\n", f.filePath)

		fileContent, err := getFileContent(f.filePath)
		if err != nil {
			fmt.Printf("ERROR: Could not read file %s : %v", f.filePath, err)
			continue
		}

		var parsedFile string

		htmlParser(fileContent, &parsedFile)

		removeDocumentFrequency(&inMemoryData.Df, inMemoryData.Ftf, f)

		tf := getTermsFrequency(parsedFile)
		docTitle := getDocTitle(fileContent)
		inMemoryData.Ftf[f.filePath] = FileData{Terms: tf, Title: docTitle, DocSize: len(tf), LastUpdateTime: f.lastUpdateTime}

		getDocumentFrequency(&inMemoryData.Df, tf)
	}
}
