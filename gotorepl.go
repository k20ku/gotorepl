package gotorepl

import (
	"fmt"
	"io"

	liner "github.com/peterh/liner"
)

func Run() {
	cl := newContLiner()
	defer cl.Close()

	cl.InitDebugger()
	cl.Debugln("main: init contLiner")
	for {
		in, err := cl.Prompt("repl")
		if err != nil {
			if err == io.EOF {
				break
			} else if err == liner.ErrPromptAborted {
				continue
			}
			fmt.Printf("main: %+v\n", err)
		}

		if in == "" {
			continue
		}

		if err := cl.Reindent("rapl"); err != nil {
			// fmt.Fprintf(os.Stderr, "err: %s\n", err)
			cl.Clear()
			continue
		}

		if cl.CountDepth() < 0 {
			fmt.Printf("%s: %s", in, errUnmatchedBraces)
			cl.Clear()
			continue
		} else if cl.CountDepth() == 0 {
			fmt.Println(cl.buffer)
			cl.Accepted()
		} else {
			continue
		}

	}

}
