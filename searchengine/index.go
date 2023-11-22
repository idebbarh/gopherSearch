package searchengine

import (
	"fmt"
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

func indexHandler(filesInfo []FileInfo, inMemoryData *InMemoryData, indexFileName string) {
	for _, f := range filesInfo {
		// check if the file already indexed.
		v, ok := inMemoryData.Ftf[f.filePath]

		// if the file already index.
		if ok {
			// check if file modified or not ,if its ignore it,
			// else redindex it.
			if v.LastUpdateTime.Equal(f.lastUpdateTime) {
				fmt.Printf("ignoring %s....\n", f.filePath)
				continue
			}
		}

		fmt.Printf("indexing %s....\n", f.filePath)

		// read the file
		fileContent, err := getFileContent(f.filePath)
		if err != nil {
			fmt.Printf("ERROR: Could not read file %s : %v", f.filePath, err)
			continue
		}

		var parsedFile string

		// parse the html
		htmlParser(fileContent, &parsedFile)

		// reset the df, because will recalc the df, of the file and we dont want to calc the terms again.
		removeDocumentFrequency(&inMemoryData.Df, inMemoryData.Ftf, f)

		tf := getTermsFrequency(parsedFile)
		docTitle := getDocTitle(fileContent)
		inMemoryData.Ftf[f.filePath] = FileData{Terms: tf, Title: docTitle, DocSize: len(tf), LastUpdateTime: f.lastUpdateTime}

		getDocumentFrequency(&inMemoryData.Df, tf)
	}

	saveToJson(indexFileName, *inMemoryData)
}
