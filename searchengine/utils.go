package searchengine

import (
	"errors"
	"fmt"
	"sort"
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
		fmt.Println(p.value)
		docs = append(docs, p.key)
	}

	return docs
}

func assert(condition bool, message string) {
	if !condition {
		panic(errors.New("Assertion failed: " + message))
	}
}
