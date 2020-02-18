package logbook

import (
	"testing"
	"time"
)

func TestAddEntry(t *testing.T) {
	d := Day{time.Now(), []Entry{}}
	e := Entry{time.Now(), "My Message"}

	err := d.Add(e)
	if err != nil {
		t.Error(err)
	}

	if len(d.Entries) != 1 {
		t.Error("Didn't add")
	}
}

func TestAddEntryBadDate(t *testing.T) {
	d := Day{time.Now(), []Entry{}}
	e := Entry{time.Date(2008, 8, 24, 0, 0, 0, 0, time.UTC), "My Message"}

	err := d.Add(e)
	if err == nil {
		t.Error("Expected error")
	}

	if len(d.Entries) != 0 {
		t.Error("Added anyway")
	}
}
