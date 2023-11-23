package searchengine

import (
	"encoding/json"
	"errors"
	"fmt"
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

type DirInfo struct {
	Entries []string
	ModTime time.Time
	isDir   bool
}

type FolderEntriesInfo = map[string]*DirInfo

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
		var serverWG sync.WaitGroup

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

		events := goWatch(c.Path)

		go func() {
			for {
				select {
				case event := <-events:
					switch true {
					case event.Types.Write:
						file := event.Info.WriteInfo.Name
						fmt.Printf("edited file name: %s\n", file)
						fileInfo := make(FilesInfo)
						getPathFiles(file, fileInfo)
						go indexHandler(fileInfo, &inMemoryData, indexFileName, false)
					case event.Types.Create:
						file := event.Info.CreateInfo.Name
						fmt.Printf("created file name: %s\n", file)
						fileInfo := make(FilesInfo)
						getPathFiles(file, fileInfo)
						go indexHandler(fileInfo, &inMemoryData, indexFileName, false)
					case event.Types.Delete:
						file := event.Info.DeleteInfo.Name
						fmt.Printf("deleted file name: %s\n", file)
						fmt.Printf("removing %s from the cache...\n", file)
						delete(inMemoryData.Ftf, file)
						saveToJson(indexFileName, inMemoryData)
					case event.Types.Rename:
						fmt.Printf("prevname: %s, new name: %s\n", event.Info.RenameInfo.PrevName, event.Info.RenameInfo.NewName)
					}
				}
			}
		}()

		filesInfo := make(FilesInfo)
		getPathFiles(c.Path, filesInfo)
		go indexHandler(filesInfo, &inMemoryData, indexFileName, true)

		go serveHandler(&inMemoryData, &serverWG)
		// wait for the server to finish then exit the program
		serverWG.Wait()
	default:
		PrintErrorToUser(UNKOWN_SUBCOMMAND)
	}
}
