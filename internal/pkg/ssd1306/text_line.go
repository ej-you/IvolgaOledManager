package ssd1306

import (
	"fmt"
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)

// func newTextLine(prefix, msg string, line int) *textLine {
// 	// image for final output
// 	img := image1bit.NewVerticalLSB(image.Rect(0, 0, _displayWidth, _lineHeight))

// 	// set up font for text
// 	face := basicfont.Face7x13
// 	face.Height = _lineHeight

// 	// init drawer to output text as image on display
// 	drawer := font.Drawer{
// 		Dst:  img,
// 		Src:  image.White,
// 		Face: face,
// 		Dot:  fixed.P(0, face.Ascent),
// 	}

// 	lineRect := image.Rect(0, line*_lineHeight, _displayWidth, (line+1)*_lineHeight)
// 	// draw text on line
// 	drawer.DrawString(prefix + msg)
// }

// DisplayText displays text line.
func (s SSD1306) DisplayTextLine(prefix, msg string, line int) error {
	// image for final output
	img := image1bit.NewVerticalLSB(image.Rect(0, 0, _displayWidth, _lineHeight))

	// set up font for text
	face := basicfont.Face7x13
	face.Height = _lineHeight

	// init drawer to output text as image on display
	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(0, face.Ascent),
	}

	lineRect := image.Rect(0, line*_lineHeight, _displayWidth, (line+1)*_lineHeight)
	// draw every text line on a new line
	// drawer.Dot = fixed.P(0, 0+face.Ascent+line*face.Height)
	drawer.DrawString(prefix + msg)
	// for idx, line := range text.Normalize(msg) {
	// }
	if err := s.device.Draw(lineRect, img, image.Pt(0, 0)); err != nil {
		return fmt.Errorf("draw text: %w", err)
	}
	return nil
}
