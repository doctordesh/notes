package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/doctordesh/notes"
)

type config struct {
	FilePath string `json:"file_path"`
}

func getConfig(configFilePath string) (config, error) {
	var c config
	var err error

	f, err := os.Open(configFilePath)
	if err != nil {
		return c, fmt.Errorf("ERROR: Config file ~/.notes-config.json is missing")
	}

	configFile, err := ioutil.ReadAll(f)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(configFile, &c)
	if err != nil {
		return c, fmt.Errorf("ERROR: ~/.notes-config.json is not well formatted JSON (%s)", err)
	}

	if c.FilePath == "" {
		return c, fmt.Errorf("ERROR: Key 'file_path' missing from ~/.notes-config.json config file")
	}

	return c, nil
}

func main() {
	var editorFlag bool
	var listFlag bool
	var notesConfigFile string
	var manualFlag bool

	flag.BoolVar(&editorFlag, "e", false, "Open editor with give string as message")
	flag.BoolVar(&listFlag, "l", false, "Lists the notebook in a 'less' styled output")
	flag.StringVar(&notesConfigFile, "c", "/Users/emil/.notes-config.json", "Config file to use")
	flag.BoolVar(&manualFlag, "m", false, "Opens $EDITOR with the data file for manual editing")

	flag.Parse()

	c, err := getConfig(notesConfigFile)
	if err != nil {
		fmt.Println(err)
	}

	commandPrinter := os.Getenv("PAGER")
	if commandPrinter == "" {
		commandPrinter = "less"
	}

	commandEditor := os.Getenv("EDITOR")
	if commandEditor == "" {
		commandEditor = "nano"
	}

	n, err := notes.New(
		c.FilePath,
		notes.NewPrinter(commandPrinter),
		notes.NewEditor(commandEditor),
		notes.NewTagger(),
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Call 'Notes'
	if listFlag {
		err = n.Print()
	} else if editorFlag {
		err = n.AddWithEditor(strings.Join(flag.Args(), " "))
	} else if manualFlag {
		err = n.ManualEdit()
	} else {
		err = n.Add(strings.Join(flag.Args(), " "))
	}

	if err != nil {
		if err == notes.ErrEmptyNote {
			flag.Usage()
		} else {
			fmt.Printf("Error: %s", err)
			return
		}
	}
}
