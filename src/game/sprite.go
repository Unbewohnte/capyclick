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
func (s *Sprite) MoveTo(x float64, y float64) {
	s.X = x
	s.Y = y
	screenBounds := WindowBounds()

	// Constraints
	// Right
	if s.X+float64(s.RealBounds().Dx()) >= float64(screenBounds.Dx()) {
		s.X = float64(screenBounds.Dx()) - float64(s.RealBounds().Dx())
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
	if s.Y+float64(s.RealBounds().Dy()) >= float64(screenBounds.Dy()) {
		s.Y = float64(screenBounds.Dy()) - float64(s.RealBounds().Dy())
	}

}
