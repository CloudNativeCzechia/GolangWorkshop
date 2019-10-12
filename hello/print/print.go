package print

import "fmt"

// Print External function printing the stuff
func Print(word string) {
	fmt.Print(word)
}

// internal routine that is not exported
func internalRoutine() {
	// do something
}
