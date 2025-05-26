// Package text contains functions to process text for output on display.
package text

import "strings"

const _maxLineSymbols = 18 // max amount of symbols on one line

// Normalize splits text string by \n and display width.
func Normalize(msg string) []string {
	// var normalized []string

	// split by \n
	lines := strings.Split(msg, "\n")
	// copy(lines, normalized)

	// for _, line := range lines {

	// }

	return lines
}

// func splitByWidthStep(startIdx int, msgLines []string) (stopIdx int) {
// 	var stopIdx int

// 	for idx, line := range msgLines[startIdx:] {
// 		if len(line) > _maxLineSymbols {
// 			msgLines = append(msgLines[:idx], line[:_maxLineSymbols+1], line[_maxLineSymbols+1:])
// 			msgLines = append(msgLines, msgLines[idx:]...)
// 		}
// 	}
// }
