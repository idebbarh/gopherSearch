package searchengine

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ChangeType = int

const (
	WRITE ChangeType = iota
	CREATE
	DELETE
	RENAME
	NOCHANGE
)

type EventsInfo struct {
	WriteInfo  struct{ Name string }
	CreateInfo struct{ Name string }
	DeleteInfo struct{ Name string }
	RenameInfo struct {
		PrevName string
		NewName  string
	}
}

type EventsType struct {
	Write  bool
	Create bool
	Delete bool
	Rename bool
}
type Event struct {
	Types EventsType
	Info  *EventsInfo
}
type FileInfo struct {
	filePath       string
	lastUpdateTime time.Time
}

type FilesInfo = map[string]FileInfo

func getFileContent(filePath string) (string, error) {
	fileContent, err := os.ReadFile(filePath)
	return string(fileContent), err
}

func getPathFiles(curPath string, filesInfo FilesInfo) {
	fi, err := os.Stat(curPath)
	if err != nil {
		fmt.Printf("ERROR: Could not get info of %s : %v\n", curPath, err)
		os.Exit(1)
	}
	mode := fi.Mode()

	if mode.IsRegular() {

		pathParts := strings.Split(curPath, "/")
		lastPart := pathParts[len(pathParts)-1]
		lastPartParts := strings.Split(lastPart, ".")

		if len(lastPartParts) >= 2 && lastPartParts[len(lastPartParts)-1] == "html" {
			filesInfo[curPath] = FileInfo{filePath: curPath, lastUpdateTime: fi.ModTime()}
		}

	} else if mode.IsDir() {
		entries, err := os.ReadDir(curPath)
		if err != nil {
			fmt.Printf("ERROR: Could not read dir %s : %v", curPath, err)
			os.Exit(1)
		}

		for _, l := range entries {
			getPathFiles(curPath+"/"+l.Name(), filesInfo)
		}
	}
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

func listener(watchingPath string, events chan Event) {
	fmt.Printf("listening on : %s\n", watchingPath)
	prevFolderEntriesInfo := FolderEntriesInfo{}
	getFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)

	for {
		curFolderEntriesInfo := FolderEntriesInfo{}
		getFolderEntriesInfo(watchingPath, curFolderEntriesInfo)

		isSomethingChange, changeType, eventInfo := entriesScanner(watchingPath, prevFolderEntriesInfo, curFolderEntriesInfo)

		if isSomethingChange {
			prevFolderEntriesInfo = make(FolderEntriesInfo)
			getFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)
			events <- Event{Types: EventsType{Write: changeType == WRITE, Create: changeType == CREATE, Delete: changeType == DELETE, Rename: changeType == RENAME}, Info: &eventInfo}
		}

		time.Sleep(1 * time.Second)
	}
}

func goWatch(watchingPath string) chan Event {
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
	events := make(chan Event)

	go listener(watchingPath, events)

	return events
}

func entriesScanner(watchingPath string, prevFolderEntriesInfo FolderEntriesInfo, curFolderEntriesInfo FolderEntriesInfo) (bool, ChangeType, EventsInfo) {
	prevWatchingPathInfo, ok := prevFolderEntriesInfo[watchingPath]
	eventInfo := EventsInfo{}
	if !ok {
		for path := range prevFolderEntriesInfo {
			_, ok := curFolderEntriesInfo[path]
			if !ok {
				eventInfo.RenameInfo.PrevName = path
				eventInfo.RenameInfo.NewName = watchingPath
				break
			}
		}

		return true, RENAME, eventInfo
	}
	curWatchingPathInfo := curFolderEntriesInfo[watchingPath]
	if len(curWatchingPathInfo.Entries) < len(prevWatchingPathInfo.Entries) {
		for path := range prevFolderEntriesInfo {
			_, ok := curFolderEntriesInfo[path]
			if !ok {
				eventInfo.DeleteInfo.Name = path
				break
			}
		}
		return true, DELETE, eventInfo
	}

	if len(curWatchingPathInfo.Entries) > len(prevWatchingPathInfo.Entries) {
		for path := range curFolderEntriesInfo {
			_, ok := prevFolderEntriesInfo[path]
			if !ok {
				eventInfo.CreateInfo.Name = path
				break
			}
		}
		return true, CREATE, eventInfo
	}

	if len(curWatchingPathInfo.Entries) == len(prevWatchingPathInfo.Entries) {
		for _, curEntryPath := range curWatchingPathInfo.Entries {
			curEntryInfo := curFolderEntriesInfo[curEntryPath]
			if curEntryInfo.isDir {
				isSomethingChange, changeType, eventInfo := entriesScanner(curEntryPath, prevFolderEntriesInfo, curFolderEntriesInfo)
				if isSomethingChange {
					return isSomethingChange, changeType, eventInfo
				}
			} else {
				prevEntryInfo, ok := prevFolderEntriesInfo[curEntryPath]
				if !ok || prevEntryInfo.ModTime.Second() != curEntryInfo.ModTime.Second() {
					if ok {
						eventInfo.WriteInfo.Name = curEntryPath
						return true, WRITE, eventInfo
					} else {
						for path := range prevFolderEntriesInfo {
							_, ok := curFolderEntriesInfo[path]
							if !ok {
								eventInfo.RenameInfo.PrevName = path
								eventInfo.RenameInfo.NewName = curEntryPath
								break
							}
						}
						return true, RENAME, eventInfo
					}
				}
			}
		}
	}

	return false, NOCHANGE, eventInfo
}
