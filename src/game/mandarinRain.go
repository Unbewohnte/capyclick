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
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type MandarinRain struct {
	InProgress           bool
	MandarinBox          *Physical
	Mandarins            []*Physical
	Completed            bool
	mandarinCount        uint16
	mandarinInitialCount uint16
	mandarinsInBox       uint16
	boxFull              bool
	mandarinCountRange   [2]uint16
}

func NewMandarinRain(from uint16, to uint16) *MandarinRain {
	rain := MandarinRain{}
	rain.InProgress = false
	rain.mandarinInitialCount = uint16(rand.Int31n(int32(to-from)) + int32(from))
	rain.mandarinCountRange = [2]uint16{from, to}
	rain.mandarinCount = rain.mandarinInitialCount
	rain.mandarinsInBox = 0
	rain.boxFull = false
	rain.Completed = false

	rain.Mandarins = make([]*Physical, rain.mandarinInitialCount)
	for i := 0; i < int(rain.mandarinInitialCount); i++ {
		rain.Mandarins[i] = NewPhysical(NewSpriteFromFile("mandarin_orange.png"), 10.0)
	}

	rain.MandarinBox = NewPhysical(NewSpriteFromFile("mandarin_box_empty.png"), 5.5)

	return &rain
}

func (mr *MandarinRain) PhysicalAt(x int, y int) *Physical {
	for _, orange := range mr.Mandarins {
		if orange.Sprite.IsIn(x, y) {
			return orange
		}
	}

	if mr.MandarinBox.Sprite.IsIn(x, y) {
		return mr.MandarinBox
	}

	return nil
}

func (mr *MandarinRain) Run(game *Game) {
	if mr.InProgress {
		return
	}

	mr.InProgress = true

	// Move oranges to random positions on the top of the screen
	for _, orange := range mr.Mandarins {
		orange.Sprite.MoveTo(float64(rand.Int31n(int32(game.Screen.Bounds().Dx()-orange.Sprite.Img.Bounds().Dx()))), 10.0, game.Screen)
	}

	// Create mandarin box
	mr.MandarinBox.Sprite.MoveTo(
		float64(rand.Int31n(int32(game.Screen.Bounds().Dx()-mr.MandarinBox.Sprite.Img.Bounds().Dx()))),
		10.0, game.Screen,
	)
}

func (mr *MandarinRain) Update(game *Game) {
	// Oranges
	temp := mr.Mandarins[:0]
	for _, orange := range mr.Mandarins {
		orange.Acceleration.Vx = 0.0
		orange.Acceleration.Vy = 9.81 / orange.Mass

		orange.Velocity.Vx = orange.Velocity.Vx + orange.Acceleration.Vx*0.05
		orange.Velocity.Vy = orange.Velocity.Vy + orange.Acceleration.Vy*0.05

		oBounds := orange.Sprite.RealBounds()
		oX := orange.Sprite.X
		oY := orange.Sprite.Y

		// Constraints
		// Right
		if oX+float64(oBounds.Dx()) >= float64(game.Screen.Bounds().Dx()) {
			orange.Velocity.Vx = -orange.Velocity.Vx * 0.4
		}

		// Left
		if oX <= 0 {
			orange.Velocity.Vx = -orange.Velocity.Vx * 0.4
		}

		// Up
		if oY <= 0.0 {
			orange.Velocity.Vy = -orange.Velocity.Vy * 0.4
		}

		// Bottom
		if oY+float64(oBounds.Dy()) >= float64(game.Screen.Bounds().Dy()) {
			orange.Velocity.Vx = orange.Velocity.Vx * 0.4 // friction on the floor
			orange.Velocity.Vy = -orange.Velocity.Vy * 0.4
		}

		orange.Sprite.X += orange.Velocity.Vx
		orange.Sprite.Y += orange.Velocity.Vy

		// Move the orange
		orange.Sprite.MoveTo(orange.Sprite.X, orange.Sprite.Y, game.Screen)

		// Check whether it touches mandarin box
		if orange.InVicinity(mr.MandarinBox.Sprite.X, mr.MandarinBox.Sprite.Y, float64(mr.MandarinBox.Sprite.RealBounds().Dx())) {
			// Yes!
			mr.mandarinsInBox++
			mr.mandarinCount--
			game.PlaySound("orange_put")
		} else {
			// Do not include this orange in the next update (effectively, delete it)
			temp = append(temp, orange)
		}
	}
	mr.Mandarins = temp

	// Orange box
	mr.MandarinBox.Acceleration.Vx = 0.0
	mr.MandarinBox.Acceleration.Vy = 9.81 / mr.MandarinBox.Mass

	mr.MandarinBox.Velocity.Vx = mr.MandarinBox.Velocity.Vx + mr.MandarinBox.Acceleration.Vx*0.05
	mr.MandarinBox.Velocity.Vy = mr.MandarinBox.Velocity.Vy + mr.MandarinBox.Acceleration.Vy*0.05

	mBounds := mr.MandarinBox.Sprite.RealBounds()
	mX := mr.MandarinBox.Sprite.X
	mY := mr.MandarinBox.Sprite.Y

	// Right
	if mX+float64(mBounds.Dx()) >= float64(game.Screen.Bounds().Dx()) {
		mr.MandarinBox.Velocity.Vx = -mr.MandarinBox.Velocity.Vx * 0.3
	}

	// Left
	if mX <= 0 {
		mr.MandarinBox.Velocity.Vx = -mr.MandarinBox.Velocity.Vx * 0.3
	}

	// Up
	if mY <= 0.0 {
		mr.MandarinBox.Velocity.Vy = -mr.MandarinBox.Velocity.Vy * 0.3
	}

	// Bottom
	if mY+float64(mBounds.Dy()) >= float64(game.Screen.Bounds().Dy()) {
		mr.MandarinBox.Velocity.Vx = mr.MandarinBox.Velocity.Vx * 0.3 // friction on the floor
		mr.MandarinBox.Velocity.Vy = -mr.MandarinBox.Velocity.Vy * 0.3
	}

	mr.MandarinBox.Sprite.X += mr.MandarinBox.Velocity.Vx
	mr.MandarinBox.Sprite.Y += mr.MandarinBox.Velocity.Vy

	// Move box
	mr.MandarinBox.Sprite.MoveTo(mr.MandarinBox.Sprite.X, mr.MandarinBox.Sprite.Y, game.Screen)

	if mr.mandarinsInBox == mr.mandarinInitialCount && !mr.boxFull {
		// All oranges are in a box!
		mr.boxFull = true
		game.PlaySound("mandarin_box_full")
	}

	// If the box is full with mandarines and is near capybara - end mandarin rain and reward with points!
	if mr.boxFull && mr.MandarinBox.InVicinity(
		game.Capybara.Sprite.X+float64(game.Capybara.Sprite.RealBounds().Dx()/2),
		game.Capybara.Sprite.Y+float64(game.Capybara.Sprite.RealBounds().Dy()/2),
		float64(game.Screen.Bounds().Dx())/7) {
		// Give a reward and finish this mandarin rain!
		game.Save.Points += pointsForLevel(game.Save.Level+1) / 5
		game.PlaySound("mandarin_rain_completed")
		mr.InProgress = false
		mr.Completed = true
	}
}

func (mr *MandarinRain) Draw(screen *ebiten.Image) {
	if mr.InProgress {
		// Mandarin box
		if mr.mandarinsInBox < mr.mandarinInitialCount && mr.mandarinsInBox > 0 {
			mr.MandarinBox.Sprite.ChangeImageByName("mandarin_box_not_empty.png")
		} else if mr.mandarinsInBox == mr.mandarinInitialCount {
			mr.MandarinBox.Sprite.ChangeImageByName("mandarin_box_full.png")
		} else {
			mr.MandarinBox.Sprite.ChangeImageByName("mandarin_box_empty.png")
		}

		op := &ebiten.DrawImageOptions{}
		scale := float64(screen.Bounds().Dx()) / float64(mr.MandarinBox.Sprite.Img.Bounds().Dx()) / 6.0
		mr.MandarinBox.Sprite.Scale = scale // Save current scale for proper collision detection
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(mr.MandarinBox.Sprite.X, mr.MandarinBox.Sprite.Y)
		screen.DrawImage(mr.MandarinBox.Sprite.Img, op)

		// Oranges
		for _, orange := range mr.Mandarins {
			op = &ebiten.DrawImageOptions{}
			scale = float64(screen.Bounds().Dx()) / float64(orange.Sprite.Img.Bounds().Dx()) / 11.5
			orange.Sprite.Scale = scale // Save current scale for proper collision detection
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(orange.Sprite.X, orange.Sprite.Y)
			screen.DrawImage(orange.Sprite.Img, op)
		}
	}
}
