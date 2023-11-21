package searchengine

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ChangeType = int

type EventsType struct {
	Write  bool
	Create bool
	Delete bool
	Rename bool
}

const (
	WRITE ChangeType = iota
	CREATE
	DELETE
	RENAME
	NOCHANGE
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

	entriesInfo[curPath] = &DirInfo{}
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

		entriesInfo[entryPath] = &DirInfo{}
		entriesInfo[entryPath].ModTime = entryInfo.ModTime()
		entriesInfo[entryPath].isDir = false
	}
}

func listener(watchingPath string, events chan EventsType) {
	fmt.Printf("listening on : %s\n", watchingPath)
	prevFolderEntriesInfo := FolderEntriesInfo{}
	getFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)
	for {
		curFolderEntriesInfo := FolderEntriesInfo{}
		getFolderEntriesInfo(watchingPath, curFolderEntriesInfo)

		isSomethingChange, changeType := entriesScanner(watchingPath, prevFolderEntriesInfo, curFolderEntriesInfo)

		if isSomethingChange {
			prevFolderEntriesInfo = make(FolderEntriesInfo)
			getFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)
			events <- EventsType{Write: changeType == WRITE, Create: changeType == CREATE, Delete: changeType == DELETE, Rename: changeType == RENAME}
		}

		time.Sleep(1 * time.Second)
	}
}

func goWatch(watchingPath string) chan EventsType {
	fi, err := os.Stat(watchingPath)
	if err != nil {
		fmt.Printf("ERROR: Could not get info of %s : %v", watchingPath, err)
		os.Exit(1)
	}

	mode := fi.Mode()

	if !mode.IsDir() {
		fmt.Printf("Error: could not listener to this path because its not a folder")
		os.Exit(1)
	}
	events := make(chan EventsType)

	go listener(watchingPath, events)

	return events
}

func entriesScanner(watchingPath string, prevFolderEntriesInfo FolderEntriesInfo, curFolderEntriesInfo FolderEntriesInfo) (bool, ChangeType) {
	prevWatchingPathInfo, ok := prevFolderEntriesInfo[watchingPath]
	if !ok {
		for path := range prevFolderEntriesInfo {
			_, ok := curFolderEntriesInfo[path]
			if !ok {
				fmt.Printf("warning: folder name changed from %s -> %s \n", path, watchingPath)
				break
			}
		}

		return true, RENAME
	}
	curWatchingPathInfo := curFolderEntriesInfo[watchingPath]
	if len(curWatchingPathInfo.Entries) < len(prevWatchingPathInfo.Entries) {
		for path := range prevFolderEntriesInfo {
			_, ok := curFolderEntriesInfo[path]
			if !ok {
				fmt.Printf("warning: %s is deleted from %s\n", path, watchingPath)
				break
			}
		}
		return true, DELETE
	}

	if len(curWatchingPathInfo.Entries) > len(prevWatchingPathInfo.Entries) {
		for path := range curFolderEntriesInfo {
			_, ok := prevFolderEntriesInfo[path]
			if !ok {
				fmt.Printf("warning: %s is created inside %s\n", path, watchingPath)
				break
			}
		}
		return true, CREATE
	}

	if len(curWatchingPathInfo.Entries) == len(prevWatchingPathInfo.Entries) {
		for _, curEntryPath := range curWatchingPathInfo.Entries {
			curEntryInfo := curFolderEntriesInfo[curEntryPath]
			if curEntryInfo.isDir {
				isSomethingChange, changeType := entriesScanner(curEntryPath, prevFolderEntriesInfo, curFolderEntriesInfo)
				if isSomethingChange {
					return isSomethingChange, changeType
				}
			} else {
				prevEntryInfo, ok := prevFolderEntriesInfo[curEntryPath]
				if !ok || prevEntryInfo.ModTime.Second() != curEntryInfo.ModTime.Second() {
					if ok {
						fmt.Printf("warning: %s content is updated\n", curEntryPath)
						return true, WRITE
					} else {
						for path := range prevFolderEntriesInfo {
							_, ok := curFolderEntriesInfo[path]
							if !ok {
								fmt.Printf("warning: file name changed from %s -> %s \n", path, curEntryPath)
								break
							}
						}
						return true, RENAME
					}
				}
			}
		}
	}

	return false, NOCHANGE
}
