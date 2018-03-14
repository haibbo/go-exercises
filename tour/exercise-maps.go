package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	wordc := make(map[string]int)
	for _, w := range strings.Fields(s) {
		wordc[w]++
	}
	return wordc
}

func main() {
	wc.Test(WordCount)
}
