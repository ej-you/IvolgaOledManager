package ssd1306

import (
	"fmt"
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"periph.io/x/devices/v3/ssd1306/image1bit"

	"sschmc/internal/pkg/text"
)

const _lineHeight = 15 // height for one text line on display

// DisplayText displays text.
func (s SSD1306) DisplayText(msg string, x, y int) error {
	// image for final output
	img := image1bit.NewVerticalLSB(image.Rect(0, 0, 128, 64))

	// init drawer to output text as image on display
	face := basicfont.Face7x13
	face.Height = _lineHeight

	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
	}

	// draw every text line on a new line
	for idx, line := range text.Normalize(msg) {
		drawer.Dot = fixed.P(0, 0+face.Ascent+idx*face.Height)
		drawer.DrawString(line)
	}
	if err := s.device.Draw(s.device.Bounds(), img, image.Pt(-x, -y)); err != nil {
		return fmt.Errorf("draw text: %w", err)
	}
	return nil
}
