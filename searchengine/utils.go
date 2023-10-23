package searchengine

import "errors"

func assert(condition bool, message string) {
	if !condition {
		panic(errors.New("Assertion failed: " + message))
	}
}
