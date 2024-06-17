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

type Capybara struct {
	Sprite *Sprite
}

func NewCapybara(sprite *Sprite) *Capybara {
	return &Capybara{
		Sprite: sprite,
	}
}

func (c *Capybara) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) != 0 {
		c.Sprite.Animation.Squish += 0.5
	}

	// Capybara Animation
	capyAniData := &c.Sprite.Animation
	if capyAniData.Theta >= 0.03 {
		capyAniData.BounceDirectionFlag = false
	} else if capyAniData.Theta <= -0.03 {
		capyAniData.BounceDirectionFlag = true
	}

	if capyAniData.Squish >= 0 {
		capyAniData.Squish -= 0.05
	}

	if capyAniData.BounceDirectionFlag {
		capyAniData.Theta += 0.001
	} else {
		capyAniData.Theta -= 0.001
	}
}

func (c *Capybara) Draw(screen *ebiten.Image, level uint32) {
	// Capybara
	switch level {
	case 1:
		c.Sprite.ChangeImageByName("capybara_1.png")
	case 2:
		c.Sprite.ChangeImageByName("capybara_2.png")
	case 3:
		c.Sprite.ChangeImageByName("capybara_3.png")
	default:
		c.Sprite.ChangeImageByName("capybara_3.png")
	}

	op := &ebiten.DrawImageOptions{}
	capybaraBounds := c.Sprite.Img.Bounds()
	scale := float64(screen.Bounds().Dx()) / float64(capybaraBounds.Dx()) / 2.5
	c.Sprite.Scale = scale
	op.GeoM.Scale(
		scale+c.Sprite.Animation.Squish,
		scale-c.Sprite.Animation.Squish,
	)
	op.GeoM.Rotate(c.Sprite.Animation.Theta)

	capyWidth := float64(c.Sprite.RealBounds().Dx())
	capyHeight := float64(c.Sprite.RealBounds().Dy())
	c.Sprite.MoveTo(float64(screen.Bounds().Dx()/2)-capyWidth/2, float64(screen.Bounds().Dy()/2)-capyHeight/2, screen)

	op.GeoM.Translate(c.Sprite.X, c.Sprite.Y)

	screen.DrawImage(c.Sprite.Img, op)
}
