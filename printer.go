package notes

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/doctordesh/notes/logbook"
	"github.com/eidolon/wordwrap"
)

type Printer interface {
	Print(logbook logbook.Logbook) error
}

type printer struct {
	cmd string
}

func NewPrinter(cmd string) Printer {
	return printer{cmd}
}

func (p printer) Print(l logbook.Logbook) error {
	days := l.Entries()

	for i := range days {
		for range days[i].Entries {
			sort.SliceStable(days[i].Entries, func(si, sj int) bool {
				return days[i].Entries[si].Time.Unix() > days[i].Entries[sj].Time.Unix()
			})
		}

		sort.SliceStable(days, func(si, sj int) bool {
			return days[si].Date.Unix() > days[sj].Date.Unix()
		})
	}

	s := ""
	wrapper := wordwrap.Wrapper(100, false)
	first := true
	for _, day := range days {
		if !first {
			// Divider between days
			s += "\n\n- - - - - - - \n"
		}
		first = false
		s += fmt.Sprintf("\n%s\n", day.Date.Format("2006-01-02 - Monday"))

		for _, entry := range day.Entries {
			timestamp := fmt.Sprintf("  %s - ", entry.Time.Format("15:04:05"))
			note := wrapper(entry.Note)
			s += wordwrap.Indent(note, timestamp, false) + "\n"
		}
	}

	buf := bytes.NewBufferString(s)

	cmd := exec.Command(p.cmd)
	cmd.Stdin = buf
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
