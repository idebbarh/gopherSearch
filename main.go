package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

type TermsFrequency = map[string]int

type FilesTermsFrequency = map[string]TermsFrequency

const (
	NO_SUBCOMMAND = iota
	NO_PATH_TO_INDEX
	NO_FILE_TO_SERVE
	UNKOWN_SUBCOMMAND
	TOTAL_ERRORS
)

func simpleHtmlParser(htmlContent string, parsedContent *string) {
	if len(htmlContent) == 0 {
		return
	}

	count := 0

	if string(htmlContent[count]) == "<" {
		for count < len(htmlContent) && string(htmlContent[count]) != ">" {
			count += 1
		}

		if count < len(htmlContent) {
			htmlContent = htmlContent[count+1:]
			simpleHtmlParser(htmlContent, parsedContent)
		}
		return
	}

	for count < len(htmlContent) && string(htmlContent[count]) != "<" {
		count += 1
	}

	*parsedContent += htmlContent[:count]

	if count < len(htmlContent) {
		htmlContent = htmlContent[count:]
		simpleHtmlParser(htmlContent, parsedContent)
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

func getFileContent(filePath string) []string {
	return strings.Split("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.", " ")
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

func getTermsFrequency(terms []string) TermsFrequency {
	tf := TermsFrequency{}
	for _, t := range terms {
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
		fileContent := getFileContent(f)
		tf := getTermsFrequency(fileContent)
		ftf[f] = tf
	}

	fmt.Println(len(ftf))
	for v, k := range ftf {
		fmt.Println(v, k)
	}
}

func (c Command) handleCommand() {
	switch c.subcommand {
	case "index":
		fmt.Printf("indexing: %s", c.path)
		indexHandler(c.path)
	case "serve":
		fmt.Printf("serving: %s", c.path)
	default:
		printHelpToUser(UNKOWN_SUBCOMMAND)
	}
}

func main() {
	simpleHtml := `<!DOCTYPE html>
<html>
<head>
    <title>Complex HTML Example</title>
</head>
<body>
    <header>
        <h1>Welcome to our website!</h1>
        <nav>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/about">About</a></li>
                <li><a href="/services">Services</a></li>
                <li><a href="/contact">Contact</a></li>
            </ul>
        </nav>
    </header>
    <main>
        <section id="about">
            <h2>About Us</h2>
            <p>
                We are a dedicated team of professionals with a mission to provide high-quality services to our clients.
            </p>
        </section>
        <section id="services">
            <h2>Our Services</h2>
            <ul>
                <li>Web Development</li>
                <li>Mobile App Development</li>
                <li>Digital Marketing</li>
            </ul>
        </section>
        <section id="contact">
            <h2>Contact Us</h2>
            <address>
                Email: <a href="mailto:contact@example.com">contact@example.com</a><br>
                Phone: <a href="tel:+1234567890">123-456-7890</a>
            </address>
        </section>
    </main>
    <footer>
        &copy; 2023 Company Name
    </footer>
</body>
</html>`
	var parsedContent string
	simpleHtmlParser(simpleHtml, &parsedContent)

	fmt.Println(parsedContent)

	// args := os.Args
	//
	// if len(args) < 2 {
	// 	printHelpToUser(NO_SUBCOMMAND)
	// }
	//
	// args = args[1:]
	//
	// if len(args) < 2 {
	// 	subcommand := args[0]
	// 	if subcommand == "index" {
	// 		printHelpToUser(NO_PATH_TO_INDEX)
	// 	} else if subcommand == "serve" {
	// 		printHelpToUser(NO_FILE_TO_SERVE)
	// 	} else {
	// 		printHelpToUser(UNKOWN_SUBCOMMAND)
	// 	}
	// }
	//
	// command := Command{
	// 	subcommand: args[0],
	// 	path:       args[1],
	// }
	//
	// command.handleCommand()
}
