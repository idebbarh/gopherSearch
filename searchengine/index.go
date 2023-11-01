package searchengine

import (
	"fmt"
	"strings"
)

type FileData struct {
	// terms ferq in doc
	Terms TermsFrequency
	// doc title
	Title string
	// number of terms in the doc
	DocSize int
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

func indexHandler(curPath string) {
	files := getPathFiles(curPath)
	ftf := FilesTermsFrequency{}
	df := DocumentFrequency{Size: 0, Value: map[string]int{}}

	for _, f := range files {
		fmt.Printf("indexing %s....\n", f)
		fileContent, err := getFileContent(f)
		if err != nil {
			fmt.Printf("ERROR: Could not read file %s : %v", f, err)
			continue
		}

		var parsedFile string
		htmlParser(fileContent, &parsedFile)
		tf := getTermsFrequency(parsedFile)
		docTitle := getDocTitle(fileContent)

		ftf[f] = FileData{Terms: tf, Title: docTitle, DocSize: len(tf)}

		for t := range tf {
			_, ok := df.Value[t]
			if ok {
				df.Value[t] += 1
			} else {
				df.Value[t] = 1
			}

		}

		df.Size += 1
	}

	inMemoryData := InMemoryData{Ftf: ftf, Df: df}

	pathParts := strings.Split(curPath, "/")

	indexFileName := pathParts[len(pathParts)-1] + ".index.json"

	saveToJson(indexFileName, inMemoryData)
}
