package logbook

import (
	"time"
)

type Logbook interface {
	Add(note string, tags []string) error
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Entries() []Day
}

func New() Logbook {
	return &logbook{Days: []Day{}}
}

type logbook struct {
	Days []Day `json:"days"`

	readBuffer  []byte
	readPointer int

	writeBuffer []byte
}

func (l *logbook) Entries() []Day {
	return l.Days
}

func (l *logbook) Add(note string, tags []string) error {
	e := Entry{Time: time.Now(), Note: note, Tags: tags}
	index := 0
	found := false

	for i := range l.Days {
		if sameDay(l.Days[i].Date, e.Time) {
			index = i
			found = true
		}
	}

	if found == false {
		l.Days = append(l.Days, Day{Date: time.Now(), Entries: []Entry{}})
		index = len(l.Days) - 1
	}

	return l.Days[index].Add(e)
}

func sameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
