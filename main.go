package main

import (
	"errors"
	"fmt"
	"os"
	"unicode"
)

type Command struct {
	subcommand string
	path       string
}

type TermsFrequency = map[string]int

type FilesTermsFrequency = map[string]TermsFrequency

const (
	NO_SUBCOMMAND = iota
	NO_PATH_TO_INDEX
	NO_FILE_TO_SERVE
	UNKOWN_SUBCOMMAND
	TOTAL_ERRORS
)

func assert(condition bool, message string) {
	if !condition {
		panic(errors.New("Assertion failed: " + message))
	}
}

func printHelpToUser(errorType int) {
	assert(TOTAL_ERRORS == 4, "You are not handling all error types")
	switch errorType {
	case NO_SUBCOMMAND:
		fmt.Println("ERROR: You must provide a subcommand.")
		fmt.Println("Usage: program <subcommand>")
		fmt.Println("Subcommands:")
		fmt.Println("  index  <path_to_files>  - Index the files.")
		fmt.Println("  serve <path_to_file>   - Serve the indexed files.")

	case NO_FILE_TO_SERVE:
		fmt.Println("ERROR: You must provide the path to the indexed file to serve.")
		fmt.Println("Usage: program serve <path_to_file>")

	case NO_PATH_TO_INDEX:
		fmt.Println("ERROR: you must provide a path to the file or directory to index.")
		fmt.Println("Usage: program index <path_to_file_or_folder>")
	case UNKOWN_SUBCOMMAND:
		fmt.Println("ERROR: Unknown subcommand")
		fmt.Println("Usage: program <subcommand>")
		fmt.Println("Subcommands:")
		fmt.Println("  index  <path_to_files>  - Index the files.")
		fmt.Println("  serve <path_to_file>   - Serve the indexed files.")
	default:
		fmt.Println("ERROR: Unknown error")
	}

	os.Exit(1)
}

func htmlParser(htmlContent string, parsedContent *string) {
	if len(htmlContent) == 0 {
		return
	}

	index := 0

	if string(htmlContent[index]) == "<" {
		for index < len(htmlContent) && string(htmlContent[index]) != ">" {
			index += 1
		}

		if index < len(htmlContent) {
			htmlContent = htmlContent[index+1:]
			htmlParser(htmlContent, parsedContent)
		}
		return
	}

	for index < len(htmlContent) && string(htmlContent[index]) != "<" {
		index += 1
	}

	*parsedContent += htmlContent[:index]

	if index < len(htmlContent) {
		htmlContent = htmlContent[index:]
		htmlParser(htmlContent, parsedContent)
	}
}

func lexer(content string) []string {
	if len(content) == 0 {
		return []string{}
	}

	var res []string
	index := 0

	r := rune(content[index])

	if unicode.IsSpace(r) {
		for index < len(content) && unicode.IsSpace(rune(content[index])) {
			index += 1
		}
		if index < len(content) {
			content = content[index:]
			res = append(res, lexer(content)...)
		}
	} else if unicode.IsLetter(r) {
		for index < len(content) && (unicode.IsLetter(rune(content[index])) || unicode.IsDigit(rune(content[index]))) {
			index += 1
		}

		res = append(res, content[:index])
		if index < len(content) {
			content = content[index:]
			res = append(res, lexer(content)...)
		}
	} else if unicode.IsDigit(r) {

		for index < len(content) && unicode.IsDigit(rune(content[index])) {
			index += 1
		}

		res = append(res, content[:index])
		if index < len(content) {
			content = content[index:]
			res = append(res, lexer(content)...)
		}
	} else {
		index += 1
		res = append(res, content[:index])
		if index < len(content) {
			content = content[index:]
			res = append(res, lexer(content)...)
		}
	}

	return res
}

func getFileContent(filePath string) (string, error) {
	fileContent, err := os.ReadFile(filePath)
	return string(fileContent), err
}

func getPathFiles(curPath string) []string {
	curFiles := []string{}
	fi, err := os.Stat(curPath)
	if err != nil {
		fmt.Printf("ERROR: Could not get info of %s : %v", curPath, err)
		os.Exit(1)
	}
	mode := fi.Mode()
	if mode.IsRegular() {
		return []string{curPath}
	} else if mode.IsDir() {
		entries, err := os.ReadDir(curPath)
		if err != nil {
			fmt.Printf("ERROR: Could not read dir %s : %v", curPath, err)
			os.Exit(1)
		}

		for _, l := range entries {
			curFiles = append(curFiles, getPathFiles(curPath+"/"+l.Name())...)
		}
	}
	return curFiles
}

func getTermsFrequency(fileContent string) TermsFrequency {
	tf := TermsFrequency{}
	for _, t := range lexer(fileContent) {
		_, ok := tf[t]
		if ok {
			tf[t] += 1
		} else {
			tf[t] = 1
		}
	}
	return tf
}

func indexHandler(curPath string) {
	files := getPathFiles(curPath)
	ftf := FilesTermsFrequency{}

	for _, f := range files {
		fileContent, err := getFileContent(f)
		if err != nil {
			fmt.Printf("ERROR: Could not read file %s : %v", f, err)
			continue
		}

		var parsedFile string
		htmlParser(fileContent, &parsedFile)
		tf := getTermsFrequency(parsedFile)
		ftf[f] = tf
	}

	for f, tf := range ftf {
		fmt.Println(len(tf))
		fmt.Println(f, "==>")
		for w, freq := range tf {
			fmt.Println("  ", w, "==>", freq)
		}
	}
}

func (c Command) handleCommand() {
	switch c.subcommand {
	case "index":
		fmt.Printf("indexing: %s\n", c.path)
		indexHandler(c.path)
	case "serve":
		fmt.Printf("serving: %s", c.path)
	default:
		printHelpToUser(UNKOWN_SUBCOMMAND)
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		printHelpToUser(NO_SUBCOMMAND)
	}

	args = args[1:]

	if len(args) < 2 {
		subcommand := args[0]
		if subcommand == "index" {
			printHelpToUser(NO_PATH_TO_INDEX)
		} else if subcommand == "serve" {
			printHelpToUser(NO_FILE_TO_SERVE)
		} else {
			printHelpToUser(UNKOWN_SUBCOMMAND)
		}
	}

	command := Command{
		subcommand: args[0],
		path:       args[1],
	}

	command.handleCommand()
}
