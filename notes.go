package notes

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	ManualEdit() error
}

type notes struct {
	filename string
	logbook  logbook.Logbook
	printer  Printer
	editor   Editor
	tagger   Tagger

	sizeOfOriginal int64
}

func New(filename string, p Printer, e Editor, t Tagger) (Notes, error) {
	n := notes{}
	n.logbook = logbook.New()
	n.printer = p
	n.editor = e
	n.tagger = t
	n.filename = filename

	if fileExists(n.filename) {
		f, err := os.Open(n.filename)
		if err != nil {
			return &n, fmt.Errorf("could not read file %s: %w", n.filename, err)
		}

		w, err := io.Copy(n.logbook, f)
		if err != nil {
			return &n, fmt.Errorf("could not load from file: %w", err)
		}

		n.sizeOfOriginal = w
	}

	return &n, nil
}

func (n *notes) Add(note string) error {
	note = strings.TrimSpace(note)
	if len(note) == 0 {
		return ErrEmptyNote
	}
	tags := n.tagger.ScanTags()
	err := n.logbook.Add(note, tags)
	if err != nil {
		return err
	}

	err = n.store()
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

	err = n.Add(note)
	if err != nil {
		return err
	}

	return nil
}

func (n *notes) ManualEdit() error {
	return n.editor.ManualEdit(n.filename)
}

func (n *notes) store() error {
	b, err := ioutil.ReadAll(n.logbook)
	if err != nil {
		return err
	}

	if int64(len(b)) <= n.sizeOfOriginal {
		panic("New data is smaller than the data in the file. Wont save because of potential data loss")
	}

	err = ioutil.WriteFile(n.filename, b, 0644)
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
