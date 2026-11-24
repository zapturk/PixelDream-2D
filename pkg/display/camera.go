package display

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera manages viewport position, tracking, zoom, screen shake, and pixel-snapping.
type Camera struct {
	// World position of the camera center
	X float64
	Y float64

	// Viewport dimensions (usually matching Canvas virtual dimensions)
	ViewportWidth  int
	ViewportHeight int

	// Zoom factor (1.0 = normal)
	Zoom float64

	// Pixel snapping eliminates sub-pixel jitter
	PixelSnap bool

	// Bounds limits (optional)
	hasBounds bool
	minX      float64
	minY      float64
	maxX      float64
	maxY      float64

	// Screen shake parameters
	shakeIntensity float64
	shakeDuration  float64
	shakeOffsetX   float64
	shakeOffsetY   float64
}

// NewCamera creates a camera centered on the given viewport size.
func NewCamera(viewportWidth, viewportHeight int) *Camera {
	return &Camera{
		ViewportWidth:  viewportWidth,
		ViewportHeight: viewportHeight,
		Zoom:           1.0,
		PixelSnap:      true,
	}
}

// SetBounds restricts camera movement within the specified rectangular world boundaries.
func (c *Camera) SetBounds(minX, minY, maxX, maxY float64) {
	c.hasBounds = true
	c.minX = minX
	c.minY = minY
	c.maxX = maxX
	c.maxY = maxY
}

// ClearBounds removes boundary restrictions.
func (c *Camera) ClearBounds() {
	c.hasBounds = false
}

// Follow smoothly interpolates the camera towards the target position.
func (c *Camera) Follow(targetX, targetY, speed float64) {
	if speed >= 1.0 {
		c.X = targetX
		c.Y = targetY
	} else if speed > 0 {
		c.X += (targetX - c.X) * speed
		c.Y += (targetY - c.Y) * speed
	}
	c.clampToBounds()
}

// Shake triggers screen shake with intensity (max offset in pixels) and duration (in seconds).
func (c *Camera) Shake(intensity, duration float64) {
	c.shakeIntensity = intensity
	c.shakeDuration = duration
}

// Update advances camera state such as screen shake decay.
func (c *Camera) Update(dt float64) {
	if c.shakeDuration > 0 {
		c.shakeDuration -= dt
		if c.shakeDuration <= 0 {
			c.shakeDuration = 0
			c.shakeOffsetX = 0
			c.shakeOffsetY = 0
		} else {
			c.shakeOffsetX = (rand.Float64()*2 - 1) * c.shakeIntensity
			c.shakeOffsetY = (rand.Float64()*2 - 1) * c.shakeIntensity
		}
	}
}

func (c *Camera) clampToBounds() {
	if !c.hasBounds {
		return
	}
	halfW := float64(c.ViewportWidth) / (2.0 * c.Zoom)
	halfH := float64(c.ViewportHeight) / (2.0 * c.Zoom)

	if c.maxX-c.minX > float64(c.ViewportWidth)/c.Zoom {
		if c.X-halfW < c.minX {
			c.X = c.minX + halfW
		}
		if c.X+halfW > c.maxX {
			c.X = c.maxX - halfW
		}
	}
	if c.maxY-c.minY > float64(c.ViewportHeight)/c.Zoom {
		if c.Y-halfH < c.minY {
			c.Y = c.minY + halfH
		}
		if c.Y+halfH > c.maxY {
			c.Y = c.maxY - halfH
		}
	}
}

// ViewMatrix computes the transformation matrix that maps world coordinates into viewport canvas coordinates.
func (c *Camera) ViewMatrix() ebiten.GeoM {
	camX := c.X + c.shakeOffsetX
	camY := c.Y + c.shakeOffsetY

	if c.PixelSnap {
		camX = math.Round(camX)
		camY = math.Round(camY)
	}

	halfVW := float64(c.ViewportWidth) / 2.0
	halfVH := float64(c.ViewportHeight) / 2.0

	var m ebiten.GeoM
	// Translate world coordinates so camera center is at origin (0,0)
	m.Translate(-camX, -camY)
	// Apply zoom scale
	m.Scale(c.Zoom, c.Zoom)
	// Translate origin to viewport center
	m.Translate(halfVW, halfVH)

	return m
}

// WorldToScreen transforms world coordinates into virtual canvas coordinates.
func (c *Camera) WorldToScreen(worldX, worldY float64) (float64, float64) {
	m := c.ViewMatrix()
	return m.Apply(worldX, worldY)
}

// ScreenToWorld transforms virtual canvas coordinates into world coordinates.
func (c *Camera) ScreenToWorld(screenX, screenY float64) (float64, float64) {
	m := c.ViewMatrix()
	m.Invert()
	return m.Apply(screenX, screenY)
}
