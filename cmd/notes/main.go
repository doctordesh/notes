package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/doctordesh/notes"
)

func main() {
	var editorFlag bool
	var listFlag bool

	flag.BoolVar(&editorFlag, "e", false, "Open editor with give string as message")
	flag.BoolVar(&listFlag, "l", false, "Lists the notebook in a 'less' styled output")

	flag.Parse()

	command := os.Getenv("PAGER")
	if command == "" {
		command = "less"
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	n, err := notes.New("/Users/emil/.notes", notes.NewPrinter(command), notes.NewEditor(editor))
	if err != nil {
		fmt.Printf("could not add note: %w", err)
		return
	}

	if listFlag {
		err = n.Print()
	} else if editorFlag {
		err = n.AddWithEditor(strings.Join(os.Args[2:], " "))
	} else {
		err = n.Add(strings.Join(os.Args[1:], " "))
	}

	if err != nil {
		if err == notes.ErrEmptyNote {
			flag.Usage()
		} else {
			fmt.Printf("Error: %s", err)
			return
		}
	}

	err = n.Store()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

}
