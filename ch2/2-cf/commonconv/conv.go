// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 41.

//!+

package commonconv


import "fmt"

type Celsius float64
type Fahrenheit float64

type Pound float64
type Kilogram float64

type Inch float64
type Mile float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

func (c Celsius) String() string    { return fmt.Sprintf("%g°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%g°F", f) }

func (p Pound) String() string    { return fmt.Sprintf("%glb", p) }
func (k Kilogram) String() string { return fmt.Sprintf("%gkg", k) }

func (i Inch) String() string    { return fmt.Sprintf("%g″", i) }
func (m Mile) String() string { return fmt.Sprintf("%gm", m) }

// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

// PToK converts a Pound to kilogram.
func PToK(p Pound) Kilogram { return Kilogram(p * 0.45359237) }

// KtoP converts a Kilogram to Pound.
func KToP(k Kilogram) Pound { return Pound(k / 0.45359237) }

// ItoM converts a Inch to Mile.
func IToM(i Inch) Mile { return Mile(i * 0.0254) }

// MToI converts a Mile to Inch.
func MToI(m Mile) Inch { return Inch(m / 0.0254) }
//!-
