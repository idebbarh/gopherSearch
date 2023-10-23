package searchengine

import (
	"strings"
	"unicode"
)

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

		res = append(res, strings.ToUpper(content[:index]))
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
