package logbook

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	l := logbook{}
	err := l.Add("Some string")
	if err != nil {
		t.Error(err)
	}

	if len(l.Days) != 1 {
		t.Error("Entry not added, days are still 0")
	}

	if len(l.Days[0].Entries) != 1 {
		t.Error("Entry not added, entries are still 0")
	}
}

func TestAddMultiple(t *testing.T) {
	l := logbook{}
	err := l.Add("Some string 1")
	if err != nil {
		t.Error(err)
	}

	err = l.Add("Some string 2")
	if err != nil {
		t.Error(err)
	}

	if len(l.Days) != 1 {
		t.Errorf("days are wrong, got %d, expected 1", len(l.Days))
	}

	if len(l.Days[0].Entries) != 2 {
		t.Errorf("entries are wrong, got %d, expected 2", len(l.Days[0].Entries))
	}
}

func TestReader(t *testing.T) {
	l := logbook{}

	l.Add("Went to ForMAX to fix some things")
	l.Add("Balder weekly meeting")

	path, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	p := path + "/test_notes.json"

	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	n, err := io.Copy(f, &l)
	if err != nil {
		t.Error(err)
	}

	if n == 0 {
		t.Error("Nothing copied")
	}

	_ = os.Remove(p)
}

func TestWriter(t *testing.T) {
	json_string := []byte(`{"days":[{"date":"2020-02-17T16:34:56.761244+01:00","entries":[{"time":"2020-02-17T16:34:56.761244+01:00","note":"Lorem ipsum"}]}]}`)
	l := logbook{}

	data := bytes.NewBuffer(json_string)

	for {
		b, err := data.ReadByte()
		if err != nil {
			break
		}

		buf := []byte{b}
		_, err = l.Write(buf)
		if err != nil {
			panic(err)
		}
	}

	if len(l.Days) != 1 {
		t.Fatal("No days added")
	}

	if len(l.Days[0].Entries) != 1 {
		t.Error("No entries added")
	}

	expected := "Lorem ipsum"
	if l.Days[0].Entries[0].Note != expected {
		t.Errorf("Unexpected data, got %s expected %s", l.Days[0].Entries[0].Note, expected)
	}
}

func TestOverwriteWithWriter(t *testing.T) {
	data := []byte(`{"days":[{"date":"2020-02-17T16:34:56.761244+01:00","entries":[{"time":"2020-02-17T16:34:56.761244+01:00","note":"Lorem ipsum"}]}]}`)
	l := logbook{}
	l.Add("testing")

	_, err := io.Copy(&l, bytes.NewBuffer(data))

	if err == nil {
		t.Error("expected error. should not allow overwrite")
	}

	if l.Days[0].Entries[0].Note != "testing" {
		t.Error("should not allow overwrite")
	}
}
