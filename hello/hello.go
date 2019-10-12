package main

import (
	"fmt"
	"hello/print"
)

// Main entry point
// hhhhhhh
func main() {

	// variable assingment
	var variable int
	variable = 10

	// alternative assignment of the variable
	variable2 := 20

	fmt.Printf("Hello world for %d and %d", variable, variable2)
}

// INIT function for package
func init() {
	print.Print("hello init")
}
