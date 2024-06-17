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

import "math"

type Vec2f struct {
	Vx float64
	Vy float64
}

func newVec2f(x float64, y float64) Vec2f {
	return Vec2f{
		Vx: x,
		Vy: y,
	}
}

type Physical struct {
	Sprite       *Sprite
	Velocity     Vec2f
	Acceleration Vec2f
	Mass         float64
}

func NewPhysical(sprite *Sprite, mass float64) *Physical {
	return &Physical{
		Sprite:       sprite,
		Velocity:     newVec2f(0.0, 0.0),
		Acceleration: newVec2f(0.0, 0.0),
		Mass:         10.0,
	}
}

// Returns true if x and y coordinates are in the radius of the physical sprite
func (ph *Physical) InVicinity(x float64, y float64, radius float64) bool {
	distance := math.Sqrt(
		math.Pow(ph.Sprite.X-x, 2.0) + math.Pow(ph.Sprite.Y-y, 2.0),
	)

	return distance <= radius
}
