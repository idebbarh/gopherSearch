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

func calcIDF(ftf FilesTermsFrequency, t string) float64 {
	n := 0        // number of document that have this term
	N := len(ftf) // number of document
	for _, tf := range ftf {
		_, ok := tf.Terms[t]
		if ok {
			n += 1
		}
	}

	if n == 0 { // to avoid divition on zero if the term apeare in non document
		n += 1
	}

	return math.Log10(float64(N / n))
}
