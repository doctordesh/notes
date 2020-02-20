package logbook

import "time"

type Entry struct {
	Time time.Time `json:"time"`
	Note string    `json:"note"`
	Tags []string  `json:"tags"`
}
