package searchengine

import (
	"encoding/json"
	"fmt"
	"os"
)

type TermsFrequency = map[string]int

func saveToJson(filename string, inMemoryData *InMemoryData) {
	jsonData, err := json.Marshal(*inMemoryData)
	if err != nil {
		fmt.Println("ERROR: could not convert data to json format")
		os.Exit(1)
	}

	err = os.WriteFile(filename, jsonData, 0666)

	if err != nil {
		fmt.Println("ERROR: could not save data to json file")
		os.Exit(1)
	}

	fmt.Println("Data saved to: ", filename)
}

func getTermsFrequency(fileContent string) TermsFrequency {
	tf := TermsFrequency{}
	for _, t := range lexer(fileContent) {
		_, ok := tf[t]
		if ok {
			tf[t] += 1
		} else {
			tf[t] = 1
		}

	}
	return tf
}

func getDocumentFrequency(df *DocumentFrequency, tf TermsFrequency) {
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

func removeDocumentFrequency(df *DocumentFrequency, ftf FilesTermsFrequency, f FileInfo) {
	fileData, isFileExist := ftf[f.filePath]

	if isFileExist {
		for term := range fileData.Terms {
			freq, isTermExit := df.Value[term] // this always should exist
			if isTermExit {
				if freq > 1 {
					df.Value[term] -= 1
				} else {
					delete(df.Value, term)
				}
			}
		}
	}
}
