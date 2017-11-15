package main

import (
	"fmt"
	"strings"
)

type DictionaryItem struct {
	Transcript string    `json:"transcript,omitempty"`
	Meanings   []Meaning `json:"meanings,omitempty"`
}

type Meaning struct {
	Translation string   `json:"translation"`
	References  []string `json:"references"`
}

func (meaning *Meaning) String() string {
	return fmt.Sprintf(
		`%s (%s)`,
		meaning.Translation,
		strings.Join(meaning.References, `, `),
	)
}
