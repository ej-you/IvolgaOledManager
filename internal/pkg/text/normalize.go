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
	return strings.Split(linesString, "\n")
}

// StringAlignCenter aligns given string in the center.
// It panics if width is more than str.
func StringAlignCenter(str string, width int) string {
	spacesLeft := int(float64(width-len(str)) / 2) //nolint:mnd // divide in half
	spacesRight := width - (spacesLeft + len(str))
	return strings.Repeat(" ", spacesLeft) + str + strings.Repeat(" ", spacesRight)
}
