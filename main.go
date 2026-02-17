package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	liner "github.com/peterh/liner"
)

var (
	history_fname = filepath.Join(os.TempDir(), ".gotorepl_history")
	names         = []string{"alice", "bob", "carol", "duke"}
)

const PROMPT = ">> "

func main() {

	line := liner.NewLiner()
	defer func() {
		if hfd, err := os.Create(history_fname); err == nil { // truncate history_file if it exists
			line.WriteHistory(hfd)
			hfd.Close()
		} else {
			fmt.Println("Error writing history file ", err)
		}
		line.Close()
		println("Thank you. Goodbye!")
	}()

	// setting
	line.SetMultiLineMode(true)
	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (c []string) {
		for _, nm := range names {
			if strings.HasPrefix(nm, strings.ToLower(line)) {
				c = append(c, nm)
			}
		}
		return
	})

	if hfd, err := os.Open(history_fname); err == nil {
		line.ReadHistory(hfd)
		hfd.Close()
	}

	if user, err := user.Current(); err == nil {
		fmt.Printf("Hello %s! ", user.Username)
	} else {
		fmt.Print("Hello! ")
	}

	fmt.Println("This is the Monkey programming language!")
	fmt.Println("Feel free to type in commands")

	for {
		if name, err := line.Prompt(PROMPT); err == nil {

			if strings.HasPrefix(name, "/exit") {
				break
			}
			fmt.Println("Got:", name)
			line.AppendHistory(name)

		} else if err == liner.ErrPromptAborted {
			continue

		} else {
			fmt.Println("Error reading line: ", err)
		}
	}

}
