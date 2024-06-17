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
