package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) ToggleFullscreen() {
	if ebiten.IsFullscreen() {
		// Turn fullscreen off
		ebiten.SetFullscreen(false)
	} else {
		// Go fullscreen
		ebiten.SetFullscreen(true)
	}
}

func (g *Game) SaveWindowGeometry() {
	// Update configuration and save information
	width, height := ebiten.WindowSize()
	g.Config.WindowSize = [2]int{width, height}

	x, y := ebiten.WindowPosition()
	g.Config.LastWindowPosition = [2]int{x, y}
}

func WindowBounds() image.Rectangle {
	x, y := ebiten.WindowSize()
	return image.Rect(0, 0, x, y)
}
