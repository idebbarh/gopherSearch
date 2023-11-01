package searchengine

import (
	"encoding/json"
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
		// loadedJsonFile, readFileErr := os.ReadFile(c.Path)
		//
		// if readFileErr != nil {
		// 	fmt.Println("ERROR: Failed to open json file")
		// 	return
		// }

		inMemoryData := InMemoryData{}

		// json.Unmarshal(loadedJsonFile, &inMemoryData)

		indexHandler(c.Path, &inMemoryData)
	case "serve":
		loadedJsonFile, readFileErr := os.ReadFile(c.Path)

		if readFileErr != nil {
			fmt.Println("ERROR: Failed to open json file")
			return
		}

		var inMemoryData InMemoryData

		json.Unmarshal(loadedJsonFile, &inMemoryData)

		serveHandler(&inMemoryData)
	default:
		PrintErrorToUser(UNKOWN_SUBCOMMAND)
	}
}
