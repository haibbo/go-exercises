// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 43.
//!+

// Cf converts its numeric argument to Celsius and Fahrenheit.
package main

import (
	"fmt"
	"os"
	"strconv"
    "bufio"
    "strings"
	"./commonconv"
)

func main() {
    var arg string
    var t float64
    var err error
    if len(os.Args) == 1 {
        inputReader := bufio.NewReader(os.Stdin)
        arg,_ = inputReader.ReadString('\n')
    } else {
        arg = os.Args[1:][0]
    }
    t, err = strconv.ParseFloat(strings.TrimSuffix(arg, "\n"), 64)
    if err != nil {
        fmt.Fprintf(os.Stderr, "cf: %v\n", err)
        os.Exit(1)
    }
    f := commonconv.Fahrenheit(t)
    c := commonconv.Celsius(t)
    p, k := commonconv.Pound(t), commonconv.Kilogram(t)
    i, m := commonconv.Inch(t), commonconv.Mile(t)
    fmt.Printf("%s = %s, %s = %s\n",
        f, commonconv.FToC(f), c, commonconv.CToF(c))
    fmt.Printf("%s = %s, %s = %s\n",
        k, commonconv.KToP(k), p, commonconv.PToK(p))
    fmt.Printf("%s = %s, %s = %s\n",
        i, commonconv.IToM(i), m, commonconv.MToI(m))
}

//!-
