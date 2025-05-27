// Package text contains functions for working with text strings and format it.
package text

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

// Normalize create lines slice of msg with lineLen max line length.
// Used to display it on device as text lines.
func Normalize(msg string, lineLen int) []string {
	linesString := text.WrapSoft(msg, lineLen)
	// lines := make([]string, 0)

	return strings.Split(linesString, "\n")
}
