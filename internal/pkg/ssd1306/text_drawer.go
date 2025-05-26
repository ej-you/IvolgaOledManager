package ssd1306

import (
	"fmt"
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)

const (
	_lineHeight = 16 // height for one text line on display
	_maxLines   = 4  // max lines amount for display
)

var _ TextDrawer = (*textDrawer)(nil)

type TextDrawer interface {
	AddLine(prefix, msg string)
	Draw() error
}

type textLine struct {
	borders image.Rectangle
	text    *image1bit.VerticalLSB
}

// TextDrawer implementation.
type textDrawer struct {
	ssd1306  *SSD1306
	fontFace *basicfont.Face
	lines    []textLine
}

func (s *SSD1306) NewTextDrawer() TextDrawer {
	// set up font for text
	fontFace := basicfont.Face7x13
	fontFace.Height = _lineHeight

	return &textDrawer{
		ssd1306:  s,
		fontFace: fontFace,
		lines:    make([]textLine, 0),
	}
}

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
		Dot:  fixed.P(0, d.fontFace.Ascent),
	}

	// draw text on line
	drawer.DrawString(prefix + msg)
	// create new textLine instance and append it to slice
	newTextLine := textLine{
		borders: image.Rect(0, lines*_lineHeight, _displayWidth, (lines+1)*_lineHeight),
		text:    img,
	}
	d.lines = append(d.lines, newTextLine)
}

func (d *textDrawer) Draw() error {
	var err error
	for _, line := range d.lines {
		if err = d.ssd1306.device.Draw(line.borders, line.text, image.Pt(0, 0)); err != nil {
			return fmt.Errorf("draw text line: %w", err)
		}
	}
	return nil
}
