package searchengine

import (
	"math"
)

func calcTF(tf TermsFrequency, t string, docSize int) float64 {
	n, ok := tf[t] // term freq on specific document

	if !ok {
		n = 0
	}

	N := docSize // number of terms in specific document

	return float64(n) / float64(N)
}

func calcIDF(totalDocs int, df map[string]int, t string) float64 {
	n := 0
	N := totalDocs

	v, ok := df[t]

	if ok {
		n = v
	} else {
		n = 1
	}

	return math.Log10(float64(N / n))
}
