package notes

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type Editor interface {
	Edit(note string) (string, error)
	ManualEdit(filename string) error
}

type editor struct {
	cmd string
}

func NewEditor(cmd string) Editor {
	return editor{cmd}
}

func (e editor) Edit(note string) (string, error) {
	file, err := ioutil.TempFile("", "notes.*.tmp")
	if err != nil {
		return note, fmt.Errorf("could not create temporary file %w", err)
	}

	filename := file.Name()

	defer os.Remove(filename)

	// Write the note to the file
	err = ioutil.WriteFile(filename, []byte(note), 0644)
	if err != nil {
		return note, fmt.Errorf("could not write to temporary file %w", err)
	}

	// Close file, since the editor will open it
	err = file.Close()
	if err != nil {
		return note, err
	}

	// Prepare command
	cmd := exec.Command(e.cmd, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return note, fmt.Errorf("could not edit temporary file %w", err)
	}

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return note, err
	}

	note = string(fileContent)
	return note, nil
}

func (e editor) ManualEdit(filename string) error {
	cmd := exec.Command(e.cmd, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not manually edit file %w", err)
	}

	return nil
}
