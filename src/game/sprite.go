/*
  	capyclick - Capybara clicker game
    Copyright (C) 2024  Kasianov Nikolai Alekseevich (Unbewohnte)

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package game

import (
	"Unbewohnte/capyclick/resources"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimationData struct {
	Squish              float64
	Theta               float64
	BounceDirectionFlag bool
}

// Drawable image structure
type Sprite struct {
	Img       *ebiten.Image
	X         float64
	Y         float64
	Animation AnimationData
	Scale     float64
	Dragged   bool
}

func NewSprite(img image.Image) *Sprite {
	return &Sprite{
		Img: ebiten.NewImageFromImage(img),
		X:   0.0,
		Y:   0.0,
		Animation: AnimationData{
			Squish:              0.0,
			Theta:               0.0,
			BounceDirectionFlag: false,
		},
		Scale:   1.0,
		Dragged: false,
	}
}

func NewSpriteFromFile(fileName string) *Sprite {
	return NewSprite(resources.ImageFromFile(fileName))
}

func (s *Sprite) ChangeImageByName(fileName string) {
	s.Img = ebiten.NewImageFromImage(resources.ImageFromFile(fileName))
}

// Returns how big the image is with applied scale factor
func (s *Sprite) RealBounds() image.Rectangle {
	bounds := s.Img.Bounds()
	realBounds := image.Rect(0, 0, bounds.Dx()*int(s.Scale), bounds.Dy()*int(s.Scale))
	return realBounds
}

func (s *Sprite) IsIn(x int, y int) bool {
	if x >= int(s.X) && x <= (int(s.X)+s.RealBounds().Dx()) &&
		y >= int(s.Y) && y <= (int(s.Y)+s.RealBounds().Dy()) {
		return true
	}

	return false
}

// Moves sprite to given positions. Respects window constraints
func (s *Sprite) MoveTo(x float64, y float64, screenBounds *ebiten.Image) {
	s.X = x
	s.Y = y
	// Constraints
	// Right
	if s.X+float64(s.RealBounds().Dx()) >= float64(screenBounds.Bounds().Dx()) {
		s.X = float64(screenBounds.Bounds().Dx()) - float64(s.RealBounds().Dx())
	}

	// Left
	if s.X <= 0 {
		s.X = 0
	}

	// Up
	if s.Y <= 0.0 {
		s.Y = 0.0
	}

	// Bottom
	if s.Y+float64(s.RealBounds().Dy()) >= float64(screenBounds.Bounds().Dy()) {
		s.Y = float64(screenBounds.Bounds().Dy()) - float64(s.RealBounds().Dy())
	}

}
