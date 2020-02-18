package notes

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/doctordesh/notes/logbook"
)

var (
	ErrEmptyNote = errors.New("empty note not allowed")
)

type Notes interface {
	Add(note string) error
	AddWithEditor(note string) error
	Print() error
	Store() error
}

type notes struct {
	output  io.Writer
	logbook logbook.Logbook
	printer Printer
	editor  Editor
}

func New(file string, p Printer, e Editor) (Notes, error) {
	n := notes{}
	n.logbook = logbook.New()
	n.printer = p
	n.editor = e

	if fileExists(file) {
		f, err := os.Open(file)
		if err != nil {
			return &n, fmt.Errorf("could not read file %s: %w", file, err)
		}

		_, err = io.Copy(n.logbook, f)
		if err != nil {
			return &n, fmt.Errorf("could not load from file: %w", err)
		}
	}

	// Remove file
	_ = os.Remove(file)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return &n, fmt.Errorf("could not create file %s: %w", file, err)
	}

	n.output = f

	return &n, nil
}

func (n *notes) Add(note string) error {
	note = strings.TrimSpace(note)
	if len(note) == 0 {
		return ErrEmptyNote
	}
	err := n.logbook.Add(note)
	if err != nil {
		return err
	}

	return nil
}

func (n *notes) Print() error {
	return n.printer.Print(n.logbook)
}

func (n *notes) AddWithEditor(note string) error {
	note, err := n.editor.Edit(note)
	if err != nil {
		return err
	}

	return n.Add(note)
}

func (n *notes) Store() error {
	_, err := io.Copy(n.output, n.logbook)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
