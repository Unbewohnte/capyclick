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
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type StrokeSource interface {
	Position() (int, int)
	IsJustReleased() bool
}

type MouseStrokeSource struct{}

func (m *MouseStrokeSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (m *MouseStrokeSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

type TouchStrokeSource struct {
	ID ebiten.TouchID
}

func (t *TouchStrokeSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

func (t *TouchStrokeSource) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}

type Stroke struct {
	source   StrokeSource
	offsetX  float64
	offsetY  float64
	physical *Physical
}

func NewStroke(source StrokeSource, physical *Physical) *Stroke {
	physical.Sprite.Dragged = true
	x, y := source.Position()
	return &Stroke{
		source:   source,
		offsetX:  float64(x) - physical.Sprite.X,
		offsetY:  float64(y) - physical.Sprite.Y,
		physical: physical,
	}
}

func (s *Stroke) Update(game *Game) {
	if !s.physical.Sprite.Dragged {
		return
	}

	// xp, yp := s.source.Position()
	// difference := newVec2f(
	// 	(float64(xp-s.physical.Sprite.RealBounds().Dx()/2)-s.physical.Sprite.X)*3.5,
	// 	(float64(yp-s.physical.Sprite.RealBounds().Dy()/2)-s.physical.Sprite.Y)*3.5,
	// )

	// s.physical.Acceleration.Vx = difference.Vx / s.physical.Mass
	// s.physical.Acceleration.Vy = difference.Vy / s.physical.Mass

	if s.source.IsJustReleased() {
		s.physical.Sprite.Dragged = false
		return
	}

	ix, iy := s.source.Position()
	x := float64(ix) - s.offsetX
	y := float64(iy) - s.offsetY
	s.physical.Sprite.MoveTo(x, y, game.Screen)
}

func (s *Stroke) Physical() *Physical {
	return s.physical
}
