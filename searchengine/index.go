package searchengine

import (
	"fmt"
	"strings"
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

type FilesTermsFrequency = map[string]FileData

type DocumentFrequency struct {
	// each term and number of documents appear in
	Value map[string]int
	// total document
	Size int
}

type InMemoryData struct {
	Ftf FilesTermsFrequency
	Df  DocumentFrequency
}

func indexHandler(curPath string, inMemoryData *InMemoryData) {
	files := getPathFiles(curPath)
	ftf := FilesTermsFrequency{}
	df := DocumentFrequency{Size: 0, Value: map[string]int{}}

	for _, f := range files {
		fmt.Printf("indexing %s....\n", f.filePath)
		fileContent, err := getFileContent(f.filePath)
		if err != nil {
			fmt.Printf("ERROR: Could not read file %s : %v", f.filePath, err)
			continue
		}

		v, ok := ftf[f.filePath]
		if ok {
			if v.LastUpdateTime.Equal(f.lastUpdateTime) {
				fmt.Printf("ignoring %s....\n", f.filePath)
				continue
			}
		}

		var parsedFile string
		htmlParser(fileContent, &parsedFile)

		removeDocumentFrequency(&df, ftf, f)

		tf := getTermsFrequency(parsedFile)
		docTitle := getDocTitle(fileContent)
		ftf[f.filePath] = FileData{Terms: tf, Title: docTitle, DocSize: len(tf), LastUpdateTime: f.lastUpdateTime}

		getDocumentFrequency(&df, tf)
	}

	inMemoryData.Ftf = ftf
	inMemoryData.Df = df

	pathParts := strings.Split(curPath, "/")

	indexFileName := pathParts[len(pathParts)-1] + ".index.json"

	saveToJson(indexFileName, inMemoryData)
}
