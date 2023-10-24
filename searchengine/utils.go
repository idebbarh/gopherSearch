package searchengine

import (
	"errors"
	"fmt"
	"sort"
)

type Item struct {
	key   string
	value float64
}

type ItemSlice []Item

func (s ItemSlice) Len() int {
	return len(s)
}

func (s ItemSlice) Less(i, j int) bool {
	return s[i].value > s[j].value
}

func (s ItemSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sortMap(f FilesRank, n int) FilesRank {
	if n < 0 {
		n = len(f)
	}

	pairsSlice := ItemSlice{}

	for file, rank := range f {
		pairsSlice = append(pairsSlice, Item{key: file, value: rank})
	}

	sort.Sort(pairsSlice)

	sortedFilesRank := make(map[string]float64)

	count := 0

	for _, pair := range pairsSlice {
		if count == n {
			break
		}
		count += 1
		fmt.Println(pair.value)
		sortedFilesRank[pair.key] = pair.value
	}

	return sortedFilesRank
}

func assert(condition bool, message string) {
	if !condition {
		panic(errors.New("Assertion failed: " + message))
	}
}
