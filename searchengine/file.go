package searchengine

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	filePath       string
	lastUpdateTime time.Time
}

func getFileContent(filePath string) (string, error) {
	fileContent, err := os.ReadFile(filePath)
	return string(fileContent), err
}

func getPathFiles(curPath string) []FileInfo {
	curFiles := []FileInfo{}
	fi, err := os.Stat(curPath)
	if err != nil {
		fmt.Printf("ERROR: Could not get info of %s : %v", curPath, err)
		os.Exit(1)
	}
	mode := fi.Mode()

	if mode.IsRegular() {
		res := []FileInfo{}

		pathParts := strings.Split(curPath, "/")
		lastPart := pathParts[len(pathParts)-1]
		lastPartParts := strings.Split(lastPart, ".")

		if len(lastPartParts) >= 2 && lastPartParts[len(lastPartParts)-1] == "html" {
			res = append(res, FileInfo{filePath: curPath, lastUpdateTime: fi.ModTime()})
		}

		return res
	} else if mode.IsDir() {
		entries, err := os.ReadDir(curPath)
		if err != nil {
			fmt.Printf("ERROR: Could not read dir %s : %v", curPath, err)
			os.Exit(1)
		}

		for _, l := range entries {
			curFiles = append(curFiles, getPathFiles(curPath+"/"+l.Name())...)
		}
	} else {
		return []FileInfo{}
	}
	return curFiles
}

func getFolderEntriesInfo(curPath string, entriesInfo FolderEntriesInfo) {
	entries, err := os.ReadDir(curPath)
	if err != nil {
		fmt.Printf("Error: could not get the entries of: %s: %s", curPath, err)
		os.Exit(1)
	}

	if entriesInfo[curPath] == nil {
		entriesInfo[curPath] = &DirInfo{}
	}

	entriesInfo[curPath].Entries = []string{}
	entriesInfo[curPath].isDir = true

	for _, e := range entries {
		entryInfo, err := e.Info()
		if err != nil {
			fmt.Printf("ERROR: Could not get info of %s : %v", e.Name(), err)
			os.Exit(1)
		}

		entryPath := curPath + "/" + entryInfo.Name()

		entriesInfo[curPath].Entries = append(entriesInfo[curPath].Entries, entryPath)

		if entryInfo.IsDir() {
			getFolderEntriesInfo(entryPath, entriesInfo)
			continue
		}

		if entriesInfo[entryPath] == nil {
			entriesInfo[entryPath] = &DirInfo{}
		}

		entriesInfo[entryPath].ModTime = entryInfo.ModTime()
		entriesInfo[entryPath].isDir = false
	}
}

// TODO: replace the logs with a better way of notify the user if something changed
func folderListener(watchingPath string, prevFolderEntriesInfo FolderEntriesInfo, curFolderEntriesInfo FolderEntriesInfo) bool {
	if len(curFolderEntriesInfo[watchingPath].Entries) < len(prevFolderEntriesInfo[watchingPath].Entries) {
		// TODO: print the created file of folder
		fmt.Printf("warning: folder or file was deleted inside %s\n", watchingPath)
		return true
	}

	if len(curFolderEntriesInfo[watchingPath].Entries) > len(prevFolderEntriesInfo[watchingPath].Entries) {
		// TODO: print the deleted file or folder
		fmt.Printf("warning: folder or file was created inside %s\n", watchingPath)
		return true
	}

	if len(curFolderEntriesInfo[watchingPath].Entries) == len(prevFolderEntriesInfo[watchingPath].Entries) {
		for _, curEntryPath := range curFolderEntriesInfo[watchingPath].Entries {
			isCurEntryDir := curFolderEntriesInfo[curEntryPath].isDir
			if isCurEntryDir {
				isSomethingChange := folderListener(curEntryPath, prevFolderEntriesInfo, curFolderEntriesInfo)
				if isSomethingChange {
					return true
				}
			} else {
				curEntryInfo := curFolderEntriesInfo[curEntryPath]
				prevEntryInfo, ok := prevFolderEntriesInfo[curEntryPath]
				if !ok || prevEntryInfo.ModTime.Second() != curEntryInfo.ModTime.Second() {
					// TODO: tell what changed the content or the file name
					fmt.Printf("warning: %s content is updated\n", curEntryPath)
					return true
				}
			}
		}
	}

	return false
}
