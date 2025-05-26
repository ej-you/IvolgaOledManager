package ssd1306

import (
	"fmt"
	"image"
	_ "image/png" // init function for specific image codec
	"os"
)

// DisplayImage displays image from given path to the OLED display.
func (s *SSD1306) DisplayImage(imagePath string, x, y int) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("open image file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("decode image to bytes: %w", err)
	}
	if err := s.device.Draw(s.device.Bounds(), img, image.Pt(-x, -y)); err != nil {
		return fmt.Errorf("draw image: %w", err)
	}
	return nil
}
