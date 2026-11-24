package physics

import (
	"math"
)

// KinematicBody represents an entity with platformer/kinematic physics interacting with a tile grid.
type KinematicBody struct {
	Box AABB

	VelocityX float64
	VelocityY float64

	Gravity      float64
	MaxFallSpeed float64

	OnGround  bool
	OnWall    bool
	OnCeiling bool
}

// NewKinematicBody creates a new kinematic body.
func NewKinematicBody(x, y, width, height float64) *KinematicBody {
	return &KinematicBody{
		Box:          NewAABB(x, y, width, height),
		Gravity:      800.0, // pixels/sec^2
		MaxFallSpeed: 400.0, // pixels/sec
	}
}

// SetPosition sets the position of the body's bounding box.
func (kb *KinematicBody) SetPosition(x, y float64) {
	kb.Box.X = x
	kb.Box.Y = y
}

// Position returns top-left coordinates.
func (kb *KinematicBody) Position() (float64, float64) {
	return kb.Box.X, kb.Box.Y
}

// Center returns the center coordinates.
func (kb *KinematicBody) Center() (float64, float64) {
	return kb.Box.CenterX(), kb.Box.CenterY()
}

// ApplyGravity applies gravity acceleration for the given dt.
func (kb *KinematicBody) ApplyGravity(dt float64) {
	kb.VelocityY += kb.Gravity * dt
	if kb.VelocityY > kb.MaxFallSpeed {
		kb.VelocityY = kb.MaxFallSpeed
	}
}

// MoveAndSlide moves the body by Velocity * dt, resolving collisions against a solid tile grid.
func (kb *KinematicBody) MoveAndSlide(dt float64, isSolid func(tileX, tileY int) bool, tileSize int) {
	ts := float64(tileSize)
	if ts <= 0 {
		ts = 16.0
	}

	kb.OnWall = false
	kb.OnGround = false
	kb.OnCeiling = false

	// 1. Horizontal Movement & Collision Resolution
	deltaX := kb.VelocityX * dt
	if deltaX != 0 {
		kb.Box.X += deltaX

		// Find potential overlapping tile range
		minTX := int(math.Floor(kb.Box.Left() / ts))
		maxTX := int(math.Floor((kb.Box.Right() - 0.001) / ts))
		minTY := int(math.Floor(kb.Box.Top() / ts))
		maxTY := int(math.Floor((kb.Box.Bottom() - 0.001) / ts))

		for ty := minTY; ty <= maxTY; ty++ {
			for tx := minTX; tx <= maxTX; tx++ {
				if isSolid(tx, ty) {
					tileBox := NewAABB(float64(tx)*ts, float64(ty)*ts, ts, ts)
					if kb.Box.Overlaps(tileBox) {
						if deltaX > 0 {
							// Moving right, hit left edge of tile
							kb.Box.X = tileBox.Left() - kb.Box.Width
							kb.VelocityX = 0
							kb.OnWall = true
						} else if deltaX < 0 {
							// Moving left, hit right edge of tile
							kb.Box.X = tileBox.Right()
							kb.VelocityX = 0
							kb.OnWall = true
						}
					}
				}
			}
		}
	}

	// 2. Vertical Movement & Collision Resolution
	deltaY := kb.VelocityY * dt
	if deltaY != 0 {
		kb.Box.Y += deltaY

		minTX := int(math.Floor(kb.Box.Left() / ts))
		maxTX := int(math.Floor((kb.Box.Right() - 0.001) / ts))
		minTY := int(math.Floor(kb.Box.Top() / ts))
		maxTY := int(math.Floor((kb.Box.Bottom() - 0.001) / ts))

		for ty := minTY; ty <= maxTY; ty++ {
			for tx := minTX; tx <= maxTX; tx++ {
				if isSolid(tx, ty) {
					tileBox := NewAABB(float64(tx)*ts, float64(ty)*ts, ts, ts)
					if kb.Box.Overlaps(tileBox) {
						if deltaY > 0 {
							// Falling down, hit top of ground tile
							kb.Box.Y = tileBox.Top() - kb.Box.Height
							kb.VelocityY = 0
							kb.OnGround = true
						} else if deltaY < 0 {
							// Jumping up, hit bottom of ceiling tile
							kb.Box.Y = tileBox.Bottom()
							kb.VelocityY = 0
							kb.OnCeiling = true
						}
					}
				}
			}
		}
	} else {
		// Even if not moving vertically, check if grounded right on top of a tile
		checkY := kb.Box.Bottom() + 0.1
		minTX := int(math.Floor(kb.Box.Left() / ts))
		maxTX := int(math.Floor((kb.Box.Right() - 0.001) / ts))
		ty := int(math.Floor(checkY / ts))

		for tx := minTX; tx <= maxTX; tx++ {
			if isSolid(tx, ty) {
				kb.OnGround = true
				break
			}
		}
	}
}
