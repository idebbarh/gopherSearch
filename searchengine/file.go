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
