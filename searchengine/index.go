package searchengine

import (
	"fmt"
)

func indexHandler(curPath string) {
	files := getPathFiles(curPath)
	ftf := FilesTermsFrequency{}
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
		ftf[f] = tf
	}

	saveToJson("index.json", ftf)
}
