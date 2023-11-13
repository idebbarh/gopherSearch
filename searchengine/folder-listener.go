package searchengine

import (
	"fmt"
	"io/fs"
	"os"
	"time"
)

func getEntriesInfo(entries []fs.DirEntry) []fs.FileInfo {
	entriesInfo := []fs.FileInfo{}
	for _, e := range entries {
		entrieInfo, err := e.Info()
		if err != nil {
			fmt.Printf("ERROR: Could not get info of %s : %v", e.Name(), err)
			os.Exit(1)
		}
		entriesInfo = append(entriesInfo, entrieInfo)
	}
	return entriesInfo
}

func FolderListener(folderPath string) {
	fi, err := os.Stat(folderPath)
	if err != nil {
		fmt.Printf("ERROR: Could not get info of %s : %v", folderPath, err)
		os.Exit(1)
	}

	mode := fi.Mode()

	if !mode.IsDir() {
		fmt.Printf("Error: could not listener to this path because its not a folder")
		os.Exit(1)
	}

	prevEntries, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("Error: could not get the entries of: %s: %s", folderPath, err)
		os.Exit(1)
	}

	prevFolderEntries := getEntriesInfo(prevEntries)

	fmt.Printf("listening on : %s\n", folderPath)

	for {
		newFolderEntries, err := os.ReadDir(folderPath)
		if err != nil {
			fmt.Printf("Error: could not get the entries of: %s: %s", folderPath, err)
			os.Exit(1)
		}

		if len(newFolderEntries) < len(prevFolderEntries) {
			fmt.Println("file deleted!!!!")
			prevFolderEntries = getEntriesInfo(newFolderEntries)
		} else if len(newFolderEntries) > len(prevFolderEntries) {
			fmt.Println("new file added!!!!")
			prevFolderEntries = getEntriesInfo(newFolderEntries)
		} else if len(newFolderEntries) == len(prevFolderEntries) {
			for i, curNewFileState := range newFolderEntries {
				curPrevFileInfo := prevFolderEntries[i]
				if err != nil {
					fmt.Printf("ERROR: Could not get info of %s : %v", curPrevFileInfo.Name(), err)
					os.Exit(1)
				}

				curNewFileInfo, err := curNewFileState.Info()
				if err != nil {
					fmt.Printf("ERROR: Could not get info of %s : %v", curNewFileState.Name(), err)
					os.Exit(1)
				}

				if !curPrevFileInfo.ModTime().Equal(curNewFileInfo.ModTime()) {
					fmt.Printf("file %s / %s is changed!!!!\n", folderPath, curPrevFileInfo.Name())
					prevFolderEntries = getEntriesInfo(newFolderEntries)
				}
			}
		}

		time.Sleep(1 * time.Second)

	}
}
