package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	liner "github.com/peterh/liner"
)

var (
	history_fname = filepath.Join(os.TempDir(), ".gotorepl_history")
	names         = []string{"alice", "bob", "carol", "duke"}
)

const (
	promptDefault  = ">> "
	promptContinue = ".. "
	indent         = "	"
)

type contLiner struct {
	*liner.State
	buffer string
	depth  int
}

func newContLiner() *contLiner {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	return &contLiner{State: line}
}

func (cl *contLiner) promptString() string {
	// >> foo {
	// .. <indent>
	if cl.buffer != "" {
		return promptContinue + strings.Repeat(indent, cl.depth)
	}

	return promptDefault
}

func (cl *contLiner) Prompt(prompt string) (string, error) {
	line, err := cl.State.Prompt(cl.promptString())

	switch err {
	case nil:
		if cl.buffer != "" {
			cl.buffer = cl.buffer + "\n" + line
		} else {
			cl.buffer = line
		}
	case io.EOF:
		// when ^D
		if cl.buffer != "" {
			// cancel line continuation if in continuation
			cl.Accepted()
		}
		fmt.Println()
		// else do nothing
	case liner.ErrPromptAborted:
		if cl.buffer != "" {
			cl.Accepted()
		} else {
			fmt.Println("(^D to quit)")
		}
	}

	return cl.buffer, err
}

func (cl *contLiner) Accepted() {
	// cl.State.AppendHistory(cl.buffer)
	cl.Clear()
}

func (cl *contLiner) Clear() {
	cl.buffer = ""
	cl.depth = 0
}

func (cl *contLiner) Close() {
	cl.Clear()
	cl.State.Close()
}

var errUnmatchedBraces = errors.New("unmatched braces")

func (cl *contLiner) ReIndent() error {
	oldDepth := cl.depth
	cl.depth = cl.countDepth(cl.buffer)

	if cl.depth < 0 {
		return errUnmatchedBraces
	}

	if oldDepth < cl.countDepth(cl.buffer) {
		lines := strings.Split(cl.buffer, "\n")
		if len(lines) > 1 {
			lastLine := lines[len(lines)-1]
			cursorUp()
			fmt.Printf("\r%s%s", cl.promptString(), lastLine) // ..
			eraseInLine()
			fmt.Println()
		}
	}

	return nil
}

func (cl *contLiner) countDepth(src string) int {
	depth := 0
	for _, r := range src {
		switch r {
		case '{', '(':
			depth++
		case '}', ')':
			depth--
		}
	}

	return depth
}

func cursorUp() {
	defer time.Sleep(1 * time.Second)
	fmt.Print("\x1b[1A")
}

func cursorDown() {
	defer time.Sleep(1 * time.Second)
	fmt.Print("\x1b[1B")
}

func eraseInLine() {
	defer time.Sleep(1 * time.Second)
	fmt.Print("\x1b[0K")
}

func cursorToBeginThenDownBy(n int) {
	defer time.Sleep(1 * time.Second)
	fmt.Printf("\x1b[%dE", n)
}

func main() {
	defer println()

	// contl := newContLiner()
	// defer func() {
	// 	contl.Close()
	// 	println("Thank you. Goodbye!")
	// }()
	for {
		fmt.Println("row1: ")
		fmt.Println("row2: Good Morning!")
		fmt.Println("row3: ")

		cursorUp()
		cursorUp()
		fmt.Print("row2: Good Night!")
		time.Sleep(1 * time.Second)
		eraseInLine()
		cursorToBeginThenDownBy(2)
		println()
	}
}
