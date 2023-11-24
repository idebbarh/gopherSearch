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
