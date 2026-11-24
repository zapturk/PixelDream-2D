package display

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScaleMode defines how the virtual canvas is scaled to fit the window.
type ScaleMode int

const (
	// ScaleModeInteger scales by the largest whole integer factor, ensuring razor-sharp square pixels.
	ScaleModeInteger ScaleMode = iota
	// ScaleModeFit scales smoothly to fit the window while maintaining aspect ratio.
	ScaleModeFit
	// ScaleModeStretch fills the entire window, ignoring aspect ratio.
	ScaleModeStretch
)

// Canvas represents the low-resolution virtual framebuffer for pixel-perfect rendering.
type Canvas struct {
	virtualWidth  int
	virtualHeight int
	scaleMode     ScaleMode
	clearColor    color.Color
	letterboxColor color.Color

	// Internal render target
	image *ebiten.Image

	// Calculated layout cache
	scaleX      float64
	scaleY      float64
	offsetX     float64
	offsetY     float64
	screenWidth int
	screenHeight int
}

// NewCanvas creates a new virtual canvas with the specified virtual resolution.
func NewCanvas(virtualWidth, virtualHeight int) *Canvas {
	if virtualWidth <= 0 {
		virtualWidth = 320
	}
	if virtualHeight <= 0 {
		virtualHeight = 180
	}

	return &Canvas{
		virtualWidth:   virtualWidth,
		virtualHeight:  virtualHeight,
		scaleMode:      ScaleModeInteger,
		clearColor:     color.RGBA{R: 20, G: 20, B: 24, A: 255},
		letterboxColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		image:          ebiten.NewImage(virtualWidth, virtualHeight),
	}
}

// VirtualWidth returns the width of the virtual canvas in pixels.
func (c *Canvas) VirtualWidth() int {
	return c.virtualWidth
}

// VirtualHeight returns the height of the virtual canvas in pixels.
func (c *Canvas) VirtualHeight() int {
	return c.virtualHeight
}

// SetScaleMode sets the scaling mode (Integer, Fit, or Stretch).
func (c *Canvas) SetScaleMode(mode ScaleMode) {
	c.scaleMode = mode
}

// ScaleMode returns the current scaling mode.
func (c *Canvas) ScaleMode() ScaleMode {
	return c.scaleMode
}

// SetClearColor sets the background color used to clear the canvas each frame.
func (c *Canvas) SetClearColor(col color.Color) {
	c.clearColor = col
}

// SetLetterboxColor sets the color of the bars outside the canvas.
func (c *Canvas) SetLetterboxColor(col color.Color) {
	c.letterboxColor = col
}

// Image returns the underlying ebiten.Image target for custom drawing.
func (c *Canvas) Image() *ebiten.Image {
	return c.image
}

// Clear clears the virtual canvas with the configured clear color.
func (c *Canvas) Clear() {
	c.image.Fill(c.clearColor)
}

// UpdateLayout calculates scale factors and centering offsets for the given screen dimensions.
func (c *Canvas) UpdateLayout(screenWidth, screenHeight int) {
	c.screenWidth = screenWidth
	c.screenHeight = screenHeight

	vw := float64(c.virtualWidth)
	vh := float64(c.virtualHeight)
	sw := float64(screenWidth)
	sh := float64(screenHeight)

	switch c.scaleMode {
	case ScaleModeInteger:
		factorX := math.Floor(sw / vw)
		factorY := math.Floor(sh / vh)
		factor := math.Min(factorX, factorY)
		if factor < 1 {
			factor = 1
		}
		c.scaleX = factor
		c.scaleY = factor
		c.offsetX = math.Floor((sw - vw*factor) / 2)
		c.offsetY = math.Floor((sh - vh*factor) / 2)

	case ScaleModeFit:
		factor := math.Min(sw/vw, sh/vh)
		c.scaleX = factor
		c.scaleY = factor
		c.offsetX = (sw - vw*factor) / 2
		c.offsetY = (sh - vh*factor) / 2

	case ScaleModeStretch:
		c.scaleX = sw / vw
		c.scaleY = sh / vh
		c.offsetX = 0
		c.offsetY = 0
	}
}

// RenderToScreen draws the virtual canvas onto the physical screen buffer with proper scaling and letterboxing.
func (c *Canvas) RenderToScreen(screen *ebiten.Image) {
	screen.Fill(c.letterboxColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.scaleX, c.scaleY)
	op.GeoM.Translate(c.offsetX, c.offsetY)
	// Nearest-neighbor filter preserves razor-sharp pixel edges
	op.Filter = ebiten.FilterNearest

	screen.DrawImage(c.image, op)
}

// ScreenToCanvas transforms a physical window screen coordinate into virtual canvas coordinates.
func (c *Canvas) ScreenToCanvas(screenX, screenY int) (int, int, bool) {
	if c.scaleX <= 0 || c.scaleY <= 0 {
		return 0, 0, false
	}

	relX := float64(screenX) - c.offsetX
	relY := float64(screenY) - c.offsetY

	canvasX := int(math.Floor(relX / c.scaleX))
	canvasY := int(math.Floor(relY / c.scaleY))

	inBounds := canvasX >= 0 && canvasX < c.virtualWidth && canvasY >= 0 && canvasY < c.virtualHeight
	return canvasX, canvasY, inBounds
}

// CanvasToScreen transforms a virtual canvas coordinate into physical screen coordinates.
func (c *Canvas) CanvasToScreen(canvasX, canvasY int) (int, int) {
	sx := int(math.Round(float64(canvasX)*c.scaleX + c.offsetX))
	sy := int(math.Round(float64(canvasY)*c.scaleY + c.offsetY))
	return sx, sy
}

// LayoutMetrics returns the current scaling factors and offsets.
func (c *Canvas) LayoutMetrics() (scaleX, scaleY, offsetX, offsetY float64) {
	return c.scaleX, c.scaleY, c.offsetX, c.offsetY
}

// Bounds returns the rectangle bounds of the virtual canvas.
func (c *Canvas) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.virtualWidth, c.virtualHeight)
}
