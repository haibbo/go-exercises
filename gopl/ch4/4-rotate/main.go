// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 86.

// Rev reverses a slice.
package main

import (
	"fmt"
)

func main() {
	//!+slice
	s := []int{0, 1, 2, 3, 4, 5}
	// Rotate s left by two positions.
    rotate(s[:], 9)
	fmt.Println(s) // "[2 3 4 5 0 1]"
	//!-slice

	// NOTE: ignoring potential errors from input.Err()
}

//!+rev
// reverse reverses a slice of ints in place.
func rotate(s []int, num int) {
    number := num % len(s)

    r := make([]int, number, len(s))
    copy(r, s[len(s)-number:])
    for i := 0; i< len(s) - number;i++ {
       r = append(r, s[i])
    }
    copy(s, r)
}

//!-rev
