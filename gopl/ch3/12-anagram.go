package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 3 {
		fmt.Printf("usage: %v string string", os.Args[0])
		os.Exit(1)
	}
	if anagram(os.Args[1], os.Args[2]) {
		fmt.Printf("%s and %s is anagram", os.Args[1], os.Args[2])
	} else {
		fmt.Printf("%s and %s is not anagram", os.Args[1], os.Args[2])
	}
}

func anagram(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for _, v := range s1 {
		if strings.Count(s1, string(v)) != strings.Count(s2, string(v)) {
			return false
		}
	}
	return true
}
