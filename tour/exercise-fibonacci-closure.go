package main

import "fmt"

func fibonacci() func() int {
	current, next := 0, 1

	return func() int {
		fib := current
		current, next = next, current+next

		return fib
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
