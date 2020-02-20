package notes

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Tagger interface {
	ScanTags() []string
}

type tagger struct{}

func NewTagger() Tagger {
	return tagger{}
}

func (t tagger) ScanTags() []string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Tags?: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}
	return strings.Split(text, " ")
}
