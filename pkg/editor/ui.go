package editor

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// UIContext manages mouse input and drawing state for immediate-mode editor widgets.
type UIContext struct {
	MouseX       int
	MouseY       int
	MouseDown    bool
	MouseJustDown bool
	MouseJustUp   bool
	Theme        Theme

	// Internal state
	hotItem    string
	activeItem string
}

// NewUIContext creates a UIContext with the default theme.
func NewUIContext() *UIContext {
	return &UIContext{
		Theme: DefaultTheme,
	}
}

// Update reads current input for the UI frame.
func (ui *UIContext) Update(mouseX, mouseY int) {
	ui.MouseX = mouseX
	ui.MouseY = mouseY
	ui.MouseDown = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	ui.MouseJustDown = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	ui.MouseJustUp = inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

// DrawPanel draws a framed panel container.
func (ui *UIContext) DrawPanel(dst *ebiten.Image, rect image.Rectangle, title string) {
	x := float32(rect.Min.X)
	y := float32(rect.Min.Y)
	w := float32(rect.Dx())
	h := float32(rect.Dy())

	// Panel background
	vector.DrawFilledRect(dst, x, y, w, h, ui.Theme.PanelBG, false)

	// Header bar if title provided
	if title != "" {
		headerH := float32(18)
		vector.DrawFilledRect(dst, x, y, w, headerH, ui.Theme.PanelHeader, false)
		vector.StrokeLine(dst, x, y+headerH, x+w, y+headerH, 1.0, ui.Theme.Border, false)
		ebitenutil.DebugPrintAt(dst, title, int(x)+6, int(y)+3)
	}

	// Panel border
	vector.StrokeRect(dst, x, y, w, h, 1.0, ui.Theme.Border, false)
}

// Button renders an interactive button. Returns true on the frame it is clicked.
func (ui *UIContext) Button(dst *ebiten.Image, rect image.Rectangle, label string, active bool) bool {
	hovered := ui.MouseX >= rect.Min.X && ui.MouseX < rect.Max.X &&
		ui.MouseY >= rect.Min.Y && ui.MouseY < rect.Max.Y

	clicked := false
	bgCol := ui.Theme.ButtonNormal
	borderCol := ui.Theme.Border

	if active {
		bgCol = ui.Theme.ButtonActive
		borderCol = ui.Theme.BorderActive
	} else if hovered {
		bgCol = ui.Theme.ButtonHover
		if ui.MouseDown {
			bgCol = ui.Theme.ButtonActive
		}
		if ui.MouseJustUp {
			clicked = true
		}
	}

	x := float32(rect.Min.X)
	y := float32(rect.Min.Y)
	w := float32(rect.Dx())
	h := float32(rect.Dy())

	vector.DrawFilledRect(dst, x, y, w, h, bgCol, false)
	vector.StrokeRect(dst, x, y, w, h, 1.0, borderCol, false)

	// Center text roughly
	textX := rect.Min.X + (rect.Dx()-len(label)*6)/2
	textY := rect.Min.Y + (rect.Dy()-12)/2
	if textX < rect.Min.X+4 {
		textX = rect.Min.X + 4
	}
	ebitenutil.DebugPrintAt(dst, label, textX, textY)

	return clicked
}

// Checkbox renders a toggleable checkbox with label. Returns the new state.
func (ui *UIContext) Checkbox(dst *ebiten.Image, x, y int, label string, checked bool) bool {
	boxRect := image.Rect(x, y, x+12, y+12)
	hovered := ui.MouseX >= boxRect.Min.X && ui.MouseX < boxRect.Max.X &&
		ui.MouseY >= boxRect.Min.Y && ui.MouseY < boxRect.Max.Y

	clicked := hovered && ui.MouseJustUp
	if clicked {
		checked = !checked
	}

	bx := float32(x)
	by := float32(y)
	vector.DrawFilledRect(dst, bx, by, 12, 12, ui.Theme.ButtonNormal, false)
	vector.StrokeRect(dst, bx, by, 12, 12, 1.0, ui.Theme.Border, false)

	if checked {
		vector.DrawFilledRect(dst, bx+3, by+3, 6, 6, ui.Theme.Accent, false)
	}

	ebitenutil.DebugPrintAt(dst, label, x+16, y)
	return checked
}

// DrawGridOverlay renders a faint grid pattern over a rectangular canvas region.
func DrawGridOverlay(dst *ebiten.Image, startX, startY, width, height float64, cellSize float64, col color.Color) {
	if cellSize <= 0 {
		return
	}

	for x := startX; x <= startX+width; x += cellSize {
		fx := float32(x)
		vector.StrokeLine(dst, fx, float32(startY), fx, float32(startY+height), 1.0, col, false)
	}
	for y := startY; y <= startY+height; y += cellSize {
		fy := float32(y)
		vector.StrokeLine(dst, float32(startX), fy, float32(startX+width), fy, 1.0, col, false)
	}
}
