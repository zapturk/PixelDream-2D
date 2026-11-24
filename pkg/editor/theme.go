package editor

import (
	"image/color"
)

// Theme defines the pixel art aesthetic colors for the in-engine editor.
type Theme struct {
	BG             color.Color // Deep dark background
	PanelBG        color.Color // Slightly lighter container panel
	PanelHeader    color.Color // Panel title bars
	Border         color.Color // Panel and widget borders
	BorderActive   color.Color // Focused / selected border
	ButtonNormal   color.Color // Default button state
	ButtonHover    color.Color // Hovered button state
	ButtonActive   color.Color // Pressed or toggled active state
	Text           color.Color // Primary text
	TextDim        color.Color // Secondary/muted text
	Accent         color.Color // Highlight accent (cyan/blue)
	AccentSuccess  color.Color // Play/success accent (green)
	AccentDanger   color.Color // Delete/danger accent (red)
	GridLine       color.Color // Faint grid overlay line
}

// DefaultTheme provides a modern dark pixel-art theme.
var DefaultTheme = Theme{
	BG:            color.RGBA{R: 20, G: 22, B: 28, A: 255},
	PanelBG:       color.RGBA{R: 28, G: 32, B: 40, A: 255},
	PanelHeader:   color.RGBA{R: 38, G: 44, B: 56, A: 255},
	Border:        color.RGBA{R: 50, G: 58, B: 72, A: 255},
	BorderActive:  color.RGBA{R: 74, G: 144, B: 226, A: 255},
	ButtonNormal:  color.RGBA{R: 44, G: 52, B: 66, A: 255},
	ButtonHover:   color.RGBA{R: 58, G: 70, B: 88, A: 255},
	ButtonActive:  color.RGBA{R: 74, G: 144, B: 226, A: 255},
	Text:          color.RGBA{R: 240, G: 242, B: 245, A: 255},
	TextDim:       color.RGBA{R: 140, G: 150, B: 165, A: 255},
	Accent:        color.RGBA{R: 74, G: 144, B: 226, A: 255},
	AccentSuccess: color.RGBA{R: 46, G: 204, B: 113, A: 255},
	AccentDanger:  color.RGBA{R: 231, G: 76, B: 60, A: 255},
	GridLine:      color.RGBA{R: 255, G: 255, B: 255, A: 35},
}
