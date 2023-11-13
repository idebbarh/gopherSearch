package searchengine

import (
	"fmt"
	"os"
	"time"
)

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

	prevFolderEntries, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("Error: could not get the entries of: %s: %s", folderPath, err)
		os.Exit(1)
	}

	fmt.Printf("listening on : %s\n", folderPath)

	for {
		newFolderEntries, err := os.ReadDir(folderPath)
		if err != nil {
			fmt.Printf("Error: could not get the entries of: %s: %s", folderPath, err)
			os.Exit(1)
		}

		if len(newFolderEntries) < len(prevFolderEntries) {
			fmt.Println("file deleted!!!!")
		} else if len(newFolderEntries) > len(prevFolderEntries) {
			fmt.Println("new file added!!!!")
		} else if len(newFolderEntries) == len(prevFolderEntries) {
			for i, curNewFileState := range newFolderEntries {
				curPrevFileState := prevFolderEntries[i]
				curPrevFileInfo, err := curPrevFileState.Info()
				if err != nil {
					fmt.Printf("ERROR: Could not get info of %s : %v", folderPath+"/"+curPrevFileState.Name(), err)
					os.Exit(1)
				}

				curNewFileInfo, err := curNewFileState.Info()
				if err != nil {
					fmt.Printf("ERROR: Could not get info of %s : %v", folderPath+"/"+curNewFileState.Name(), err)
					os.Exit(1)
				}

				if !curPrevFileInfo.ModTime().Equal(curNewFileInfo.ModTime()) {
					fmt.Printf("file %s / %s is changed!!!!\n", folderPath, curPrevFileInfo.Name())
				}
			}
		}

		prevFolderEntries = newFolderEntries

		time.Sleep(2 * time.Second) // Add a delay between iterations

	}
}
