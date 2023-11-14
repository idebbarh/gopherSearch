package searchengine

import (
	"fmt"
	"io/fs"
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

func getEntriesInfo(curPath string, entries []fs.DirEntry) EntriesInfo {
	entriesInfo := EntriesInfo{}

	if entriesInfo[curPath] == nil {
		entriesInfo[curPath] = &Entry{}
	}

	entriesInfo[curPath].Size = len(entries)

	for _, e := range entries {
		entryInfo, err := e.Info()
		if err != nil {
			fmt.Printf("ERROR: Could not get info of %s : %v", e.Name(), err)
			os.Exit(1)
		}
		if entryInfo.IsDir() {
			newPath := curPath + "/" + entryInfo.Name()
			newEntries, err := os.ReadDir(newPath)
			if err != nil {
				fmt.Printf("Error: could not get the entries of: %s: %s", newPath, err)
				os.Exit(1)
			}
			getEntriesInfo(newPath, newEntries)
			continue
		}
		entriesInfo[curPath+"/"+entryInfo.Name()].Info = entryInfo
	}
	return entriesInfo
}

func folderListener(watchingPath string, prevFolderEntries *EntriesInfo) {
	newFolderEntries, err := os.ReadDir(watchingPath)
	if err != nil {
		fmt.Printf("Error: could not get the entries of: %s: %s", watchingPath, err)
		os.Exit(1)
	}

	if len(newFolderEntries) < (*prevFolderEntries)[watchingPath].Size {
		fmt.Printf("warning: folder or file was deleted inside %s\n", watchingPath)
		*prevFolderEntries = getEntriesInfo(watchingPath, newFolderEntries)
	} else if len(newFolderEntries) > (*prevFolderEntries)[watchingPath].Size {
		fmt.Printf("warning: folder or file was created inside %s\n", watchingPath)
		*prevFolderEntries = getEntriesInfo(watchingPath, newFolderEntries)
	} else if len(newFolderEntries) == (*prevFolderEntries)[watchingPath].Size {
		for _, curNewFileState := range newFolderEntries {
			curNewFileInfo, err := curNewFileState.Info()
			if err != nil {
				fmt.Printf("ERROR: Could not get info of %s : %v", curNewFileState.Name(), err)
				os.Exit(1)
			}

			curNewFileMode := curNewFileInfo.Mode()

			if curNewFileMode.IsDir() {
				folderListener(watchingPath+"/"+curNewFileInfo.Name(), prevFolderEntries)
				continue
			}

			curPrevFileInfo := (*prevFolderEntries)[watchingPath+"/"+curNewFileInfo.Name()]

			if !curPrevFileInfo.Info.ModTime().Equal(curNewFileInfo.ModTime()) {
				fmt.Printf("warning: %s/%s content is updated\n", watchingPath, curPrevFileInfo.Info.Name())
				*prevFolderEntries = getEntriesInfo(watchingPath, newFolderEntries)
			}
		}
	}
}
