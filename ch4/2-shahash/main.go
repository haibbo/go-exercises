// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 83.

// The sha256 command computes the SHA256 hash (an array) of a string.
package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"os"
)

func main() {
	htype := "HASH256"
	if len(os.Args[:]) == 3 {
		htype = os.Args[1]
	}
	//fmt.Printf("%d %d %d\n", len(os.Args[1:]))
	if htype == "384" {
		fmt.Printf("Hash %s: %x\n", htype, sha512.Sum384([]byte(os.Args[2])))
	} else if htype == "512" {
		fmt.Printf("Hash %s: %x", htype, sha512.Sum512([]byte(os.Args[2])))
	} else {
		fmt.Printf("Hash %s: %x", htype, sha256.Sum256([]byte(os.Args[1])))
	}
	// Output:
	// 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	// 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	// false
	// [32]uint8
}

//!-
