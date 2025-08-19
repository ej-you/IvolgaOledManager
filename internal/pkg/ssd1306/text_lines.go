package ssd1306

import (
	"IvolgaOledManager/internal/pkg/text"
	"fmt"
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)

// TextLine is a one text line.
type TextLine struct {
	// text
	Content string
	// relative part of the total display height occupied by the text line
	RelativeSize int
}

// lineToDraw is a one line of text used to output on display.
type lineToDraw struct {
	borders image.Rectangle
	text    *image1bit.VerticalLSB
}

// DisplayTextLines displays given text lines.
func (s *SSD1306) DisplayTextLines(lines ...TextLine) error {
	preparedLines, err := s.processTextLines(lines)
	if err != nil {
		return fmt.Errorf("prepare text lines: %w", err)
	}

	for _, line := range preparedLines {
		if err = s.device.Draw(line.borders, line.text, image.Pt(0, 0)); err != nil {
			return fmt.Errorf("draw text line: %w", err)
		}
	}
	return nil
}

// processTextLines prepare text lines for output on display.
func (s *SSD1306) processTextLines(lines []TextLine) ([]lineToDraw, error) {
	linesAmount := len(lines)
	textLines := make([]lineToDraw, 0, linesAmount)

	// calc proportions
	var allParts int
	for _, line := range lines {
		allParts += line.RelativeSize
	}
	onePartHeight := s.screenHeight / allParts

	var heightFrom, heightTo int
	var toDraw *lineToDraw
	var err error
	// prepare lines one by one
	for _, line := range lines {
		heightFrom = heightTo
		heightTo += onePartHeight * line.RelativeSize

		toDraw, err = s.createLineToDraw(heightFrom, heightTo, line.Content)
		if err != nil {
			return nil, fmt.Errorf("create line to draw: %w", err)
		}
		textLines = append(textLines, *toDraw)
	}
	return textLines, nil
}

// createLineToDraw creates line to draw using given params.
func (s *SSD1306) createLineToDraw(heightFrom, heightTo int,
	content string) (*lineToDraw, error) {

	height := heightTo - heightFrom
	// set up font for text line
	fontFace, err := text.NewRussianFont(float64(height))
	if err != nil {
		return nil, fmt.Errorf("font face: %w", err)
	}

	// image for final output
	img := image1bit.NewVerticalLSB(image.Rect(0, 0, s.screenWidth, height))
	// set up drawer
	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: fontFace,
		Dot:  fixed.P(0, fontFace.Metrics().Ascent.Ceil()),
	}

	// draw text on line
	drawer.DrawString(content)
	// create new textLine instance and append it to slice
	return &lineToDraw{
		borders: image.Rect(0, heightFrom, s.screenWidth, heightTo),
		text:    img,
	}, nil
}
