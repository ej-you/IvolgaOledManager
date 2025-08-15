package drawer

import (
	"fmt"
	"image"
	_ "image/png" // init function for specific image codec
	"os"

	"github.com/pkg/errors"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)

var _ Drawer = (*SSD1306)(nil)

// SSD1306 is a Drawer implementation. It draws data on oled display.
type SSD1306 struct {
	ScreenWidth  int
	ScreenHeight int

	device    *ssd1306.Dev
	busCloser func() error
}

// NewDrawerSSD1306 returns new SSD1306 drawer.
func NewDrawerSSD1306(bus string, screenWidth, screenHeight int) (*SSD1306, error) {
	busCloser, err := i2creg.Open(bus)
	if err != nil {
		return nil, fmt.Errorf("open bus: %w", err)
	}
	device, err := ssd1306.NewI2C(busCloser, &ssd1306.DefaultOpts)
	if err != nil {
		return nil, fmt.Errorf("open device: %w", err)
	}

	instance := &SSD1306{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,

		device:    device,
		busCloser: busCloser.Close,
	}
	if err := instance.DisplayClear(); err != nil {
		return nil, fmt.Errorf("clear device on startup: %w", err)
	}

	return instance, nil
}

// DisplayClear clears display.
func (s *SSD1306) DisplayClear() error {
	black := image1bit.NewVerticalLSB(image.Rect(0, 0, s.ScreenWidth, s.ScreenHeight))
	if err := s.device.Draw(s.device.Bounds(), black, image.Pt(0, 0)); err != nil {
		return fmt.Errorf("clear display: %w", err)
	}
	return nil
}

// Close clears display (turn off all pixels) and closes its bus.
func (s *SSD1306) Close() error {
	if err := s.device.Halt(); err != nil {
		return fmt.Errorf("clear ssd1306 oled: %w", err)
	}
	return errors.Wrap(s.busCloser(), "close ssd1306 oled")
}

// DrawImage draws image from given path on the display.
func (s *SSD1306) DrawImage(imagePath string, x, y int) error { //nolint:varnamelen // obviously
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

// DrawTextScreen draws given text screen on the display.
func (s *SSD1306) DrawTextScreen(screen *TextScreen) error {
	var err error
	for _, line := range *screen {
		if err = s.device.Draw(line.borders, line.text, image.Pt(0, 0)); err != nil {
			return fmt.Errorf("draw text line: %w", err)
		}
	}
	return nil
}
