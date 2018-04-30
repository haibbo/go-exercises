// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 73.

// Comma prints its argument numbers with a comma at each power of 1000.
//
// Example:
// 	$ go build gopl.io/ch3/comma
//	$ ./comma 1 12 123 1234 1234567890
// 	1
// 	12
// 	123
// 	1,234
// 	1,234,567,890
//
package main

import (
	"fmt"
	"os"
    "bytes"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma(os.Args[i]))
	}
}

//!+
// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
    var buf bytes.Buffer
    remainder := len(s) % 3
    if remainder != 0 {
        buf.WriteString(s[0:remainder] + ",")
    }
    for i, v := range s[remainder:] {
        if i % 3 == 0 && i != 0 {
            buf.WriteRune(',')
        }
        buf.WriteRune(v)
    }
	return buf.String()
}

//!-
