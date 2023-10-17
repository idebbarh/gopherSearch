package main

import (
	"errors"
	"fmt"
	"os"
)

func assert(condition bool, message string) {
	if !condition {
		panic(errors.New("Assertion failed: " + message))
	}
}

type Command struct {
	subcommand string
	path       string
}

const (
	NO_SUBCOMMAND     = iota
	NO_PATH_TO_INDEX  = iota
	NO_FILE_TO_SERVE  = iota
	UNKOWN_SUBCOMMAND = iota
	TOTAL_ERRORS      = iota
)

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

func handleCommands(c Command) {
	switch c.subcommand {
	case "index":
		fmt.Printf("indexing: %s", c.path)
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

	handleCommands(command)
}
