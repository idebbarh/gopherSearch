package searchengine

import (
	"errors"
	"slices"
	"sort"
	"strings"
)

type Pair struct {
	key   string
	value float64
}

type Pairs []Pair

func (p Pairs) Len() int {
	return len(p)
}

func (p Pairs) Less(i, j int) bool {
	return p[i].value > p[j].value
}

func (p Pairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func rankDocs(f FilesRank) []string {
	docs := []string{}
	pairs := Pairs{}

	for file, rank := range f {
		pairs = append(pairs, Pair{key: file, value: rank})
	}

	sort.Sort(pairs)

	for _, p := range pairs {
		docs = append(docs, p.key)
	}

	return docs
}

func getIndexFileNameFromPath(path string) string {
	pathParts := strings.Split(path, "/")

	return pathParts[len(pathParts)-1] + ".index.json"
}

func assert(condition bool, message string) {
	if !condition {
		panic(errors.New("Assertion failed: " + message))
	}
}

func isPathContainsPath(parentPath string, childPath string) (bool, string) {
	var children []string
	for len(childPath) > 0 {
		if parentPath == childPath {
			if len(children) > 0 {
				slices.Reverse(children)
				return true, strings.Join(children, "/")
			} else {
				return true, ""
			}
		}
		paths := strings.Split(childPath, "/")
		children = append(children, paths[len(paths)-1])
		childPath = strings.Join(paths[:len(paths)-1], "/")
	}
	return false, ""
}
