package searchengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

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
		inMemoryDataChan := make(chan InMemoryData)
		var inMemoryData InMemoryData
		indexFileName := getIndexFileNameFromPath(c.Path)
		go func() {
			for {
				select {
				case inMemoryData := <-inMemoryDataChan:
					// Handle received inMemoryData
					fmt.Println("Received updated inMemoryData:", inMemoryData)
					// Process the received data, if needed
				default:
					// Do other tasks or just wait for data
				}
			}
		}()

		if _, stateError := os.Stat(indexFileName); stateError == nil {
			fmt.Println("Reindexing...")
			loadedJsonFile, readFileErr := os.ReadFile(indexFileName)

			if readFileErr != nil {
				fmt.Println("ERROR: Failed to open json file")
				return
			}

			json.Unmarshal(loadedJsonFile, &inMemoryData)

			inMemoryDataChan <- inMemoryData

		} else if errors.Is(stateError, os.ErrNotExist) {
			ftf := FilesTermsFrequency{}
			df := DocumentFrequency{Size: 0, Value: map[string]int{}}
			inMemoryData = InMemoryData{Ftf: ftf, Df: df}
			fmt.Println("Indexing...")
		} else {
			fmt.Printf("ERROR: file may or may not exist:%v\n", stateError)
			return
		}

		go indexHandler(c.Path, inMemoryData, inMemoryDataChan)

		serveHandler(inMemoryDataChan)

		// go indexHandler(c.Path, inMemoryDataChan)

		// updatedInMemoryData := <-inMemoryDataChan
		//
		// saveToJson(indexFileName, updatedInMemoryData)

		// for {
		// 	updatedInMemoryData := <-inMemoryDataChan
		// 	fmt.Println(updatedInMemoryData.Df.Size)
		// }

	default:
		PrintErrorToUser(UNKOWN_SUBCOMMAND)
	}
}
