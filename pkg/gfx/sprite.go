package gfx

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zapturk/PixelDream-2D/pkg/display"
)

// OriginPreset defines standard pivot points for sprites.
type OriginPreset int

const (
	OriginTopLeft OriginPreset = iota
	OriginCenter
	OriginBottomCenter
	OriginBottomLeft
	OriginCustom
)

// Sprite represents a 2D renderable entity on the pixel canvas.
type Sprite struct {
	Texture *Texture
	// Source rectangle within the texture (zero rect means entire texture)
	SourceRect image.Rectangle

	// World or screen position
	X float64
	Y float64

	// Scale multipliers (1.0 = normal)
	ScaleX float64
	ScaleY float64

	// Rotation in radians
	Rotation float64

	// Pivot calculation mode
	OriginPreset OriginPreset
	// Custom origin/pivot coordinates (used when OriginPreset == OriginCustom)
	OriginX float64
	OriginY float64

	// Horizontal and vertical flip
	FlipX bool
	FlipY bool

	// Color tint and opacity
	ColorScale color.Color

	// Pixel snapping eliminates sub-pixel texture sampling blur
	PixelSnap bool
}

// NewSprite creates a sprite initialized with the given texture.
func NewSprite(texture *Texture) *Sprite {
	var rect image.Rectangle
	if texture != nil {
		rect = image.Rect(0, 0, texture.Width(), texture.Height())
	}

	return &Sprite{
		Texture:      texture,
		SourceRect:   rect,
		ScaleX:       1.0,
		ScaleY:       1.0,
		OriginPreset: OriginTopLeft,
		ColorScale:   color.White,
		PixelSnap:    true,
	}
}

// SetOriginPreset sets the pivot point to automatically follow a preset based on active frame size.
func (s *Sprite) SetOriginPreset(preset OriginPreset) {
	s.OriginPreset = preset
}

// SetCustomOrigin sets a manual pixel pivot point relative to frame top-left.
func (s *Sprite) SetCustomOrigin(ox, oy float64) {
	s.OriginPreset = OriginCustom
	s.OriginX = ox
	s.OriginY = oy
}

// Origin returns the effective pivot point in pixels relative to the active frame's top-left.
func (s *Sprite) Origin() (float64, float64) {
	w := float64(s.FrameWidth())
	h := float64(s.FrameHeight())

	switch s.OriginPreset {
	case OriginTopLeft:
		return 0, 0
	case OriginCenter:
		return w / 2.0, h / 2.0
	case OriginBottomCenter:
		return w / 2.0, h
	case OriginBottomLeft:
		return 0, h
	default:
		return s.OriginX, s.OriginY
	}
}

// FrameWidth returns the width of the active frame rectangle.
func (s *Sprite) FrameWidth() int {
	if !s.SourceRect.Empty() {
		return s.SourceRect.Dx()
	}
	if s.Texture != nil {
		return s.Texture.Width()
	}
	return 0
}

// FrameHeight returns the height of the active frame rectangle.
func (s *Sprite) FrameHeight() int {
	if !s.SourceRect.Empty() {
		return s.SourceRect.Dy()
	}
	if s.Texture != nil {
		return s.Texture.Height()
	}
	return 0
}

// Draw renders the sprite onto a Canvas, optionally transformed by a Camera view matrix.
func (s *Sprite) Draw(canvas *display.Canvas, cam *display.Camera) {
	if s.Texture == nil || s.Texture.Image() == nil {
		return
	}

	img := s.Texture.Image()
	if !s.SourceRect.Empty() {
		img = img.SubImage(s.SourceRect).(*ebiten.Image)
	}

	op := &ebiten.DrawImageOptions{}

	// 1. Pivot offset (move origin to 0,0)
	ox, oy := s.Origin()
	op.GeoM.Translate(-ox, -oy)

	// 2. Flip
	sx := s.ScaleX
	sy := s.ScaleY
	if s.FlipX {
		sx = -sx
	}
	if s.FlipY {
		sy = -sy
	}
	op.GeoM.Scale(sx, sy)

	// 3. Rotation
	if s.Rotation != 0 {
		op.GeoM.Rotate(s.Rotation)
	}

	// 4. Translation to world position
	posX := s.X
	posY := s.Y
	if s.PixelSnap {
		posX = math.Round(posX)
		posY = math.Round(posY)
	}
	op.GeoM.Translate(posX, posY)

	// 5. Apply camera matrix if provided
	if cam != nil {
		op.GeoM.Concat(cam.ViewMatrix())
	}

	// 6. Color tint
	if s.ColorScale != nil && s.ColorScale != color.White {
		r, g, b, a := s.ColorScale.RGBA()
		op.ColorScale.Scale(
			float32(r)/65535.0,
			float32(g)/65535.0,
			float32(b)/65535.0,
			float32(a)/65535.0,
		)
	}

	op.Filter = ebiten.FilterNearest
	canvas.Image().DrawImage(img, op)
}
