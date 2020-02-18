package logbook

import (
	"fmt"
	"time"
)

type Day struct {
	Date    time.Time `json:"date"`
	Entries []Entry   `json:"entries"`
}

func (d *Day) Add(e Entry) error {
	if sameDay(d.Date, e.Time) == false {
		return fmt.Errorf("Entry with time %v is not same day as day %v", e.Time, d.Date)
	}

	d.Entries = append(d.Entries, e)
	return nil
}
