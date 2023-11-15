package searchengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"
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

type Entry struct {
	Size int
	Info fs.FileInfo
}

type FolderEntriesInfo = map[string]*Entry

type FolderEntries = map[string][]fs.DirEntry

const (
	NO_SUBCOMMAND = iota
	NO_PATH_TO_INDEX
	UNKOWN_SUBCOMMAND
	TOTAL_ERRORS
)

func PrintErrorToUser(errorType int) {
	assert(TOTAL_ERRORS == 3, "You are not handling all error types")
	switch errorType {
	case NO_SUBCOMMAND:
		fmt.Println("ERROR: You must provide a subcommand.")
		PrintUsage()
	case NO_PATH_TO_INDEX:
		fmt.Println("ERROR: you must provide a path to the file or directory to index.")
		PrintUsage()
	case UNKOWN_SUBCOMMAND:
		fmt.Println("ERROR: Unknown subcommand")
		PrintUsage()
	default:
		fmt.Println("ERROR: Unknown error")
	}

	os.Exit(1)
}

func PrintUsage() {
	fmt.Println("Usage: program <subcommand> [arguments]")
	fmt.Println("Subcommands:")
	fmt.Println("  serve <path_to_folder>   - Index and serve the folder.")
}

func (c Command) HandleCommand() {
	switch c.Subcommand {
	case "serve":
		var inMemoryData InMemoryData
		var indexingWG sync.WaitGroup
		var serverWG sync.WaitGroup

		indexingWG.Add(1)
		serverWG.Add(1)

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
			fmt.Printf("ERROR: file may or may not exist: %v\n", stateError)
			return
		}

		go func() {
			watchingPath, err := os.Getwd()
			if err != nil {
				return
			}
			watchingPath += "/" + "testListener"

			fi, err := os.Stat(watchingPath)
			if err != nil {
				fmt.Printf("ERROR: Could not get info of %s : %v", watchingPath, err)
				os.Exit(1)
			}

			mode := fi.Mode()

			if !mode.IsDir() {
				fmt.Printf("Error: could not listener to this path because its not a folder")
				return
			}

			prevEntries, err := os.ReadDir(watchingPath)
			if err != nil {
				fmt.Printf("Error: could not get the entries of: %s: %s", watchingPath, err)
				os.Exit(1)
			}

			fmt.Printf("listening on : %s\n", watchingPath)

			folderEntriesInfo := FolderEntriesInfo{}

			getFolderEntriesInfo(watchingPath, prevEntries, folderEntriesInfo)

			for {
				folderEntries := FolderEntries{}
				getFolderEntries(watchingPath, folderEntries)

				isSomethingChange := folderListener(watchingPath, &folderEntriesInfo, folderEntries[watchingPath], folderEntries)

				if isSomethingChange {
					getFolderEntriesInfo(watchingPath, folderEntries[watchingPath], folderEntriesInfo)
				}

				time.Sleep(1 * time.Second)
			}
		}()

		go indexHandler(c.Path, &inMemoryData, &indexingWG)

		go serveHandler(&inMemoryData, &serverWG)

		// wait for indexing to finish then save data to json
		indexingWG.Wait()

		saveToJson(indexFileName, inMemoryData)

		// wait for the server to finish then exit the program
		serverWG.Wait()
	default:
		PrintErrorToUser(UNKOWN_SUBCOMMAND)
	}
}
