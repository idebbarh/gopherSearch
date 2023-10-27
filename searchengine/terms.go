package searchengine

import (
	"encoding/json"
	"fmt"
	"os"
)

type TermsFrequency = map[string]int

type FileData struct {
	Terms TermsFrequency
	Title string
}

type FilesTermsFrequency = map[string]FileData

func saveToJson(filename string, ftf FilesTermsFrequency) {
	jsonData, err := json.Marshal(ftf)
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
