package searchengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

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
type Command struct {
	Subcommand string
	Path       string
}

const (
	NO_SUBCOMMAND = iota
	NO_PATH_TO_INDEX
	NO_FILE_TO_SERVE
	UNKOWN_SUBCOMMAND
	TOTAL_ERRORS
)

func PrintErrorToUser(errorType int) {
	assert(TOTAL_ERRORS == 4, "You are not handling all error types")
	switch errorType {
	case NO_SUBCOMMAND:
		fmt.Println("ERROR: You must provide a subcommand.")
		fmt.Println("Usage: program <subcommand> [arguments]")
		fmt.Println("Subcommands:")
		fmt.Println("  index  <path_to_files>  - Index the files.")
		fmt.Println("  serve <path_to_file>   - Serve the indexed files.")

	case NO_FILE_TO_SERVE:
		fmt.Println("ERROR: You must provide the path to the indexed file to serve.")
		fmt.Println("Usage: program serve <path_to_file>")

	case NO_PATH_TO_INDEX:
		fmt.Println("ERROR: you must provide a path to the file or directory to index.")
		fmt.Println("Usage: program index <path_to_file_or_folder>")
	case UNKOWN_SUBCOMMAND:
		fmt.Println("ERROR: Unknown subcommand")
		fmt.Println("Usage: program <subcommand> [arguments]")
		fmt.Println("Subcommands:")
		fmt.Println("  index  <path_to_files>  - Index the files.")
		fmt.Println("  serve <path_to_file>   - Serve the indexed files.")
	default:
		fmt.Println("ERROR: Unknown error")
	}

	os.Exit(1)
}

func PrintUsage() {
	fmt.Println("Usage: program <subcommand> [arguments]")
	fmt.Println("Subcommands:")
	fmt.Println("  index  <path_to_files>  - Index the files.")
	fmt.Println("  serve <path_to_file>   - Serve the indexed files.")
}

func (c Command) HandleCommand() {
	switch c.Subcommand {
	case "index":
		// indexFileName := getIndexFileNameFromPath(c.Path)
		// inMemoryData := InMemoryData{}
		// indexHandler(c.Path, &inMemoryData)
		// saveToJson(indexFileName, inMemoryData)
	case "serve":
		var inMemoryData InMemoryData
		var wg sync.WaitGroup
		wg.Add(1)

		indexFileName := getIndexFileNameFromPath(c.Path)

		if _, stateError := os.Stat(indexFileName); stateError == nil {
			fmt.Println("Looking for new or changed files to reindex...")
			loadedJsonFile, readFileErr := os.ReadFile(indexFileName)

			if readFileErr != nil {
				fmt.Println("ERROR: Failed to open json file")
				return
			}

			json.Unmarshal(loadedJsonFile, &inMemoryData)

		} else if errors.Is(stateError, os.ErrNotExist) {
			fmt.Println("First time Indexing...")
			ftf := FilesTermsFrequency{}
			df := DocumentFrequency{Size: 0, Value: map[string]int{}}
			inMemoryData = InMemoryData{Ftf: ftf, Df: df}
		} else {
			fmt.Printf("ERROR: file may or may not exist:%v\n", stateError)
			return
		}

		go indexHandler(c.Path, &inMemoryData, &wg)

		go serveHandler(&inMemoryData)

		wg.Wait()

		saveToJson(indexFileName, inMemoryData)

	default:
		PrintErrorToUser(UNKOWN_SUBCOMMAND)
	}
}
