// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 83.

// The sha256 command computes the SHA256 hash (an array) of a string.
package main

import "fmt"

//!+
import "crypto/sha256"

func PopCountByShifting(b1, b2 []byte) int {
	n := 0
    for k := 0; k < len(b1);k++ {
        for i := uint(0); i < 8; i++ {
            if b1[k]&(1<<i) != b2[k]&(1<<i) {
                n++
            }
        }
    }
	return n
}

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
    fmt.Printf("%d\n", PopCountByShifting(c1[:], c2[:]))
	// Output:
	// 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	// 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	// false
	// [32]uint8
}

//!-
