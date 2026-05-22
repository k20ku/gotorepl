package gotorepl

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	liner "github.com/peterh/liner"
)

const (
	promptDefault  = "> "
	promptContinue = ". "
	indent         = "  "
)

type contLiner struct {
	*liner.State
	buffer string
	depth  int
	prompt string
	wc     io.WriteCloser
}

func newContLiner() *contLiner {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	return &contLiner{State: line}
}

func (cl *contLiner) promptStringWidth(str string, width int) string {
	var ps string
	if str != "" {
		ps = "(" + str + ") "
	}
	if cl.buffer != "" {
		return strings.Repeat(".", len(ps)) + promptContinue + strings.Repeat(indent, max(0, cl.depth+width))
		// return strings.Repeat(".", len(ps)) + promptContinue
	}

	return ps + promptDefault
}

func (cl *contLiner) promptString(str string) string {
	return cl.promptStringWidth(str, 0)
}

func (cl *contLiner) Prompt(str string) (string, error) {
	line, err := cl.State.Prompt(cl.promptString(str))

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
		fmt.Println("see you")
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
	cl.State.AppendHistory(strings.ReplaceAll(cl.buffer, "\n", " "))
	cl.Clear()
}

func (cl *contLiner) Clear() {
	cl.buffer = ""
	cl.depth = 0
	cl.prompt = ""
}

func (cl *contLiner) Close() {
	cl.Clear()
	cl.State.Close()
	if cl.wc != nil {
		cl.wc.Close()
	}
}

var errUnmatchedBraces = errors.New("unmatched braces")

func (cl *contLiner) Reindent(str string) error {
	oldDepth := cl.depth
	cl.depth = cl.CountDepth()

	if cl.depth < 0 {
		return errUnmatchedBraces
	}

	cl.Debugf("Reindent: cl.depth < oldDepth? cl.depth = %d, oldDepth=%d\n", oldDepth, cl.depth)

	lines := strings.Split(cl.buffer, "\n")
	if len(lines) > 1 {
		lastLine := lines[len(lines)-1]
		cl.Debugf("Reindent: len(lines) > 1. lastline = %s\n", lastLine)
		if cl.depth < oldDepth {
			cursorUp()
			fmt.Printf("\r%s%s", cl.promptString(str), lastLine)
			eraseInLine()
			fmt.Println()
		} else if hasPrefix(lastLine, RBRACE, RBRACKET, RPAREN) {
			cursorUp()
			fmt.Printf("\r%s%s", cl.promptStringWidth(str, -1), lastLine)
			eraseInLine()
			fmt.Println()
		}
	}

	return nil
}

func (cl *contLiner) CountDepth() int {
	l := NewLexer(cl.buffer)
	depth := 0
	for {
		switch l.NextToken() {
		case LPAREN, LBRACE, LBRACKET:
			depth++
		case RPAREN, RBRACE, RBRACKET:
			depth--
		case EOF:
			return depth
		}
	}
}

func (cl *contLiner) InitDebugger() {
	fmt.Println("gotorepl: init a debugger. wait a second...")
	time.Sleep(time.Second)
	conn, err := net.Dial("unix", "/tmp/gotorepl-debug.sock")
	if err != nil {
		fmt.Println("gotorepl: run 'make log' in another terminal to enable debugger")
		return
	}
	fmt.Println("gotorepl: debugger ready")
	cl.wc = conn
}

func (cl *contLiner) Debugf(format string, a ...any) {
	if cl.wc != nil {
		fmt.Fprintf(cl.wc, format, a...)
	}
}

func (cl *contLiner) Debugln(s string) {
	if cl.wc != nil {
		fmt.Fprintln(cl.wc, s)
	}
}
