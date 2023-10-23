package main

import (
	"os"

	se "github.com/idebbarh/GopherSearch/searchengine"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		se.PrintErrorToUser(se.NO_SUBCOMMAND)
	}

	args = args[1:]

	if (args[0]) == "help" {
		se.PrintUsage()
		os.Exit(0)
	}

	if len(args) < 2 {
		subcommand := args[0]
		if subcommand == "index" {
			se.PrintErrorToUser(se.NO_PATH_TO_INDEX)
		} else if subcommand == "serve" {
			se.PrintErrorToUser(se.NO_FILE_TO_SERVE)
		} else {
			se.PrintErrorToUser(se.UNKOWN_SUBCOMMAND)
		}
	}

	command := se.Command{
		Subcommand: args[0],
		Path:       args[1],
	}

	command.HandleCommand()
}
