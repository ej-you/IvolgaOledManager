package display

const DefaultLinesAmount = 4 // amount of lines on screen (title and items/lines)

// TextLine represents a one text line on display screen.
type TextLine struct {
	// text
	Content string
	// relative part of the total display height occupied by the text line
	RelativeSize int
}

// NewDefaultTextLine returns a new instance of TextLine
// with default relative size (it is 1).
func NewDefaultTextLine(content string) TextLine {
	return TextLine{
		Content:      content,
		RelativeSize: 1,
	}
}
