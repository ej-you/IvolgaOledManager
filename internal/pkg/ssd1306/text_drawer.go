package ssd1306

import (
	"fmt"
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"periph.io/x/devices/v3/ssd1306/image1bit"

	"IvolgaOledManager/internal/pkg/text"
)

const (
	_fontAscent = 11   // for correct text output
	_lineHeight = 16.0 // height for one text line on display
	_maxLines   = 4    // max lines amount for display
)

var _ TextDrawer = (*textDrawer)(nil)

// Interface for draw text lines.
type TextDrawer interface {
	AddLine(prefix, msg string)
	FillEmpty()
	Draw() error
}

// textLine is one text line to output on display.
type textLine struct {
	borders image.Rectangle
	text    *image1bit.VerticalLSB
}

// TextDrawer implementation.
type textDrawer struct {
	ssd1306  *SSD1306
	fontFace font.Face
	lines    []textLine
}

// NewTextDrawer returns TextDrawer for text lines output.
func (s *SSD1306) NewTextDrawer() (TextDrawer, error) {
	// set up font for text
	fontFace, err := text.NewRussianFont(_lineHeight)
	if err != nil {
		return nil, fmt.Errorf("font face: %w", err)
	}

	return &textDrawer{
		ssd1306:  s,
		fontFace: fontFace,
		lines:    make([]textLine, 0, _maxLines),
	}, nil
}

// AddLine add new line with prefix and text message to slice to draw.
func (d *textDrawer) AddLine(prefix, msg string) {
	// lines amount
	lines := len(d.lines)
	if lines == _maxLines {
		return
	}

	// image for final output
	img := image1bit.NewVerticalLSB(image.Rect(0, 0, _displayWidth, _lineHeight))
	// set up drawer
	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: d.fontFace,
		Dot:  fixed.P(0, _fontAscent),
	}

	// draw text on line
	drawer.DrawString(prefix)
	drawer.DrawString(msg)
	// create new textLine instance and append it to slice
	newTextLine := textLine{
		borders: image.Rect(0, lines*_lineHeight, _displayWidth, (lines+1)*_lineHeight),
		text:    img,
	}
	d.lines = append(d.lines, newTextLine)
}

// FillEmpty fills empty lines in black color.
func (d *textDrawer) FillEmpty() {
	// lines amount
	lines := len(d.lines)
	if lines == _maxLines {
		return
	}
	for range _maxLines - lines {
		d.AddLine("", "")
	}
}

// Draw draws every text line in lines slice.
func (d *textDrawer) Draw() error {
	var err error
	for _, line := range d.lines {
		if err = d.ssd1306.device.Draw(line.borders, line.text, image.Pt(0, 0)); err != nil {
			return fmt.Errorf("draw text line: %w", err)
		}
	}
	return nil
}
