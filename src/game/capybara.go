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
	c.Sprite.MoveTo(float64(screen.Bounds().Dx()/2)-capyWidth/2, float64(screen.Bounds().Dy()/2)-capyHeight/2)

	op.GeoM.Translate(c.Sprite.X, c.Sprite.Y)

	screen.DrawImage(c.Sprite.Img, op)
}
