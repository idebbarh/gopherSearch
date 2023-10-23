package searchengine

import (
	"testing"
)

type TestCases []TestCase

type TestCase struct {
	input    string
	expected []string
}

func isSameSlices(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func TestHtmlParserAndLexer(t *testing.T) {
	testCases := TestCases{{input: `
    <!doctype html>
    <head>
      <title></title>
    </head>
    <html>
      <body>
        <p>
          The quick brown fox jumps over the lazy dog. The dog barks, and the fox
          runs away.
        </p>
      </body>
    </html>
        `, expected: lexer("The quick brown fox jumps over the lazy dog. The dog barks, and the fox runs away.")}}

	for _, test := range testCases {
		var parsedContent string
		htmlParser(test.input, &parsedContent)
		result := lexer(parsedContent)
		expected := test.expected
		if !isSameSlices(result, expected) {
			t.Errorf("expected:\n====> %v\nresult:\n===> %v", result, expected)
		}
	}
}
