package drawer

import (
	"IvolgaOledManager/internal/pkg/text"
	"fmt"
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)

// lineToDraw is a one line of text used to output on display
// in DrawTextScreen method of drawer.
type lineToDraw struct {
	borders image.Rectangle
	text    *image1bit.VerticalLSB
}

// TextScreen is a ready-to-output slice of text lines.
type TextScreen []lineToDraw

// TextScreenBuilder is a builder for TextScreen.
type TextScreenBuilder struct {
	screenWidth  int
	screenHeight int
	linesAmount  int
}

// NewTextScreenBuilder returns new TextScreenBuilder.
func NewTextScreenBuilder(screenWidth, screenHeight int) *TextScreenBuilder {
	return &TextScreenBuilder{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// TextLine is a struct of text line.
type TextLine struct {
	// text of line
	Content string
	// relative part of the total height occupied by the text line
	RelativeSize int
}

// Build takes text lines and returns collected text screen to output.
func (b TextScreenBuilder) Build(lines ...TextLine) (*TextScreen, error) {
	b.linesAmount = len(lines)
	screen := make(TextScreen, 0, b.linesAmount)

	var allParts int
	for _, line := range lines {
		allParts += line.RelativeSize
	}
	onePartHeight := b.screenHeight / allParts

	var heightFrom, heightTo int
	var toDraw *lineToDraw
	var err error
	for _, line := range lines {
		heightFrom = heightTo
		heightTo += onePartHeight * line.RelativeSize

		toDraw, err = b.createLineToDraw(heightFrom, heightTo, line.Content)
		if err != nil {
			return nil, fmt.Errorf("create line to draw: %w", err)
		}
		screen = append(screen, *toDraw)
	}

	return &screen, nil
}

// createLineToDraw creates line to draw using given params.
func (b TextScreenBuilder) createLineToDraw(heightFrom, heightTo int,
	content string) (*lineToDraw, error) {

	height := heightTo - heightFrom
	// set up font for text line
	fontFace, err := text.NewRussianFont(float64(height))
	if err != nil {
		return nil, fmt.Errorf("font face: %w", err)
	}

	// image for final output
	img := image1bit.NewVerticalLSB(image.Rect(0, 0, b.screenWidth, height))
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
		borders: image.Rect(0, heightFrom, b.screenWidth, heightTo),
		text:    img,
	}, nil
}
