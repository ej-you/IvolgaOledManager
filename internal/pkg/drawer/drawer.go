// Package drawer provides interface Drawer to draw images and text.
package drawer

type Drawer interface {
	DrawImage(imagePath string, x, y int) error
	DrawTextScreen(screen *TextScreen) error
}
