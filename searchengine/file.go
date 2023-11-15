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

func getFolderEntriesInfo(curPath string, entries []fs.DirEntry, entriesInfo FolderEntriesInfo) {
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
			getFolderEntriesInfo(newPath, newEntries, entriesInfo)
			continue
		}

		if entriesInfo[curPath+"/"+entryInfo.Name()] == nil {
			entriesInfo[curPath+"/"+entryInfo.Name()] = &Entry{}
		}

		entriesInfo[curPath+"/"+entryInfo.Name()].Info = entryInfo
	}
}

func getFolderEntries(curPath string, folderEntries FolderEntries) {
	curEntries, err := os.ReadDir(curPath)
	if err != nil {
		fmt.Printf("Error: could not get the entries of: %s: %s", curPath, err)
		os.Exit(1)
	}

	folderEntries[curPath] = curEntries

	for _, entry := range curEntries {
		entryInfo, err := entry.Info()
		if err != nil {
			fmt.Printf("ERROR: Could not get info of %s : %v", entry.Name(), err)
			os.Exit(1)
		}

		if entryInfo.IsDir() {
			getFolderEntries(curPath+"/"+entryInfo.Name(), folderEntries)
		}
	}
}

// TODO: replace the logs with a better way of notify the user if something changed
func folderListener(watchingPath string, folderEntriesInfo *FolderEntriesInfo, curEntries []fs.DirEntry, folderEntries FolderEntries) bool {
	if len(curEntries) < (*folderEntriesInfo)[watchingPath].Size {
		// TODO: print the created file of folder
		fmt.Printf("warning: folder or file was deleted inside %s\n", watchingPath)
		return true
	}

	if len(curEntries) > (*folderEntriesInfo)[watchingPath].Size {
		// TODO: print the deleted file or folder
		fmt.Printf("warning: folder or file was created inside %s\n", watchingPath)
		return true
	}

	if len(curEntries) == (*folderEntriesInfo)[watchingPath].Size {
		for _, curEntry := range curEntries {
			curEntryInfo, err := curEntry.Info()
			if err != nil {
				fmt.Printf("ERROR: Could not get info of %s : %v", curEntry.Name(), err)
				os.Exit(1)
			}

			curNewFileMode := curEntryInfo.Mode()

			if curNewFileMode.IsDir() {
				nextWatchingPath := watchingPath + "/" + curEntry.Name()
				isSomethingChange := folderListener(nextWatchingPath, folderEntriesInfo, folderEntries[nextWatchingPath], folderEntries)
				if isSomethingChange {
					return true
				}

			} else {
				// TODO: handle file name changed
				prevEntryInfo := (*folderEntriesInfo)[watchingPath+"/"+curEntryInfo.Name()]
				if !prevEntryInfo.Info.ModTime().Equal(curEntryInfo.ModTime()) {
					fmt.Printf("warning: %s/%s content is updated\n", watchingPath, prevEntryInfo.Info.Name())
					return true
				}
			}
		}
	}

	return false
}
