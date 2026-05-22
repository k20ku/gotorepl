package gotorepl

import (
	"fmt"
)

func cursorUp() {
	fmt.Print("\x1b[1A")
}

func cursorDown() {
	fmt.Print("\x1b[1B")
}

func eraseInLine() {
	fmt.Print("\x1b[0K")
}

func cursorToBeginThenDownBy(n int) {
	fmt.Printf("\x1b[%dE", n)
}
