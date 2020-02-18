package notes

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/doctordesh/notes/logbook"
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

	for _, day := range days {
		s += fmt.Sprintf("%s\n\n", day.Date.Format("2006-01-02"))

		for _, entry := range day.Entries {
			s += fmt.Sprintf("%s - %s\n", entry.Time.Format("15:04:05"), entry.Note)
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
