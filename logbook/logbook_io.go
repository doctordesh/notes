package logbook

import (
	"encoding/json"
	"fmt"
	"io"
)

func (l *logbook) Read(p []byte) (n int, err error) {
	if l.readBuffer == nil {
		b, err := json.MarshalIndent(l, "", "    ")
		if err != nil {
			return 0, err
		}

		l.readBuffer = b
	}

	if l.readPointer >= len(l.readBuffer) {
		return 0, io.EOF
	}

	x := len(l.readBuffer) - l.readPointer
	n, bound := 0, 0
	if x >= len(p) {
		bound = len(p)
	} else if x <= len(p) {
		bound = x
	}

	copy(p, l.readBuffer[l.readPointer:l.readPointer+bound])

	l.readPointer += bound

	return bound, nil
}

func (l *logbook) Write(p []byte) (n int, err error) {
	if len(l.Days) > 0 {
		return 0, fmt.Errorf("not allowed to overwrite Logbook")
	}

	if l.writeBuffer == nil {
		l.writeBuffer = []byte{}
	}

	l.writeBuffer = append(l.writeBuffer, p...)

	if len(l.writeBuffer) > 0 && l.writeBuffer[0] != 0x7B /* the character '{' */ {
		return 0, fmt.Errorf("feeding Logbook with invalid JSON")
	}

	err = json.Unmarshal(l.writeBuffer, l)
	if err == nil { // Done
		l.writeBuffer = nil
	}

	return len(p), nil
}
