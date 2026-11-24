package physics

import (
	"math"
)

// AABB represents an Axis-Aligned Bounding Box.
type AABB struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// NewAABB creates an AABB.
func NewAABB(x, y, width, height float64) AABB {
	return AABB{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Left returns the left boundary.
func (b AABB) Left() float64 {
	return b.X
}

// Right returns the right boundary.
func (b AABB) Right() float64 {
	return b.X + b.Width
}

// Top returns the top boundary.
func (b AABB) Top() float64 {
	return b.Y
}

// Bottom returns the bottom boundary.
func (b AABB) Bottom() float64 {
	return b.Y + b.Height
}

// CenterX returns the horizontal center.
func (b AABB) CenterX() float64 {
	return b.X + b.Width/2.0
}

// CenterY returns the vertical center.
func (b AABB) CenterY() float64 {
	return b.Y + b.Height/2.0
}

// Overlaps returns true if this AABB overlaps with another AABB.
func (b AABB) Overlaps(other AABB) bool {
	return b.Left() < other.Right() &&
		b.Right() > other.Left() &&
		b.Top() < other.Bottom() &&
		b.Bottom() > other.Top()
}

// ContainsPoint returns true if the point (px, py) is inside the bounding box.
func (b AABB) ContainsPoint(px, py float64) bool {
	return px >= b.Left() && px <= b.Right() && py >= b.Top() && py <= b.Bottom()
}

// MinkovksiDiff computes Minkowski difference for continuous collision detection.
func (b AABB) MinkowskiDiff(other AABB) AABB {
	return AABB{
		X:      b.X - other.Right(),
		Y:      b.Y - other.Bottom(),
		Width:  b.Width + other.Width,
		Height: b.Height + other.Height,
	}
}

// Penetration returns the minimum penetration vector required to resolve overlap.
func (b AABB) Penetration(other AABB) (dx, dy float64) {
	overlapLeft := other.Right() - b.Left()
	overlapRight := b.Right() - other.Left()
	overlapTop := other.Bottom() - b.Top()
	overlapBottom := b.Bottom() - other.Top()

	minX := overlapLeft
	if overlapRight < minX {
		minX = -overlapRight
	}

	minY := overlapTop
	if overlapBottom < minY {
		minY = -overlapBottom
	}

	if math.Abs(minX) < math.Abs(minY) {
		return minX, 0
	}
	return 0, minY
}
