package text

import (
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const _fontDPI = 63 // dots per inch resolution for font

// NewRussianFont returns font face with Cyrillic support.
func NewRussianFont(height float64) (font.Face, error) {
	// load font with Cyrillic support
	ttf, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("parse ttf: %w", err)
	}
	// crearte font face
	fontFace, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    height,
		DPI:     _fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return fontFace, nil
}
