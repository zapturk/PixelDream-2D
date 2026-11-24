package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/zapturk/PixelDream-2D/pkg/display"
)

// Binding represents an input trigger (keyboard key, mouse button, etc.).
type Binding struct {
	Keys         []ebiten.Key
	MouseButtons []ebiten.MouseButton
}

// ActionMap manages mappings from logical action names to input triggers.
type ActionMap struct {
	bindings map[string]*Binding
}

// NewActionMap creates an empty ActionMap.
func NewActionMap() *ActionMap {
	return &ActionMap{
		bindings: make(map[string]*Binding),
	}
}

// BindKey associates an action with one or more keyboard keys.
func (am *ActionMap) BindKey(action string, keys ...ebiten.Key) {
	b, ok := am.bindings[action]
	if !ok {
		b = &Binding{}
		am.bindings[action] = b
	}
	b.Keys = append(b.Keys, keys...)
}

// BindMouseButton associates an action with one or more mouse buttons.
func (am *ActionMap) BindMouseButton(action string, buttons ...ebiten.MouseButton) {
	b, ok := am.bindings[action]
	if !ok {
		b = &Binding{}
		am.bindings[action] = b
	}
	b.MouseButtons = append(b.MouseButtons, buttons...)
}

// IsPressed returns true if any bound input for the action is currently held down.
func (am *ActionMap) IsPressed(action string) bool {
	b, ok := am.bindings[action]
	if !ok {
		return false
	}
	for _, k := range b.Keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	for _, mb := range b.MouseButtons {
		if ebiten.IsMouseButtonPressed(mb) {
			return true
		}
	}
	return false
}

// IsJustPressed returns true if any bound input was pressed in the current frame.
func (am *ActionMap) IsJustPressed(action string) bool {
	b, ok := am.bindings[action]
	if !ok {
		return false
	}
	for _, k := range b.Keys {
		if inpututil.IsKeyJustPressed(k) {
			return true
		}
	}
	for _, mb := range b.MouseButtons {
		if inpututil.IsMouseButtonJustPressed(mb) {
			return true
		}
	}
	return false
}

// IsJustReleased returns true if any bound input was released in the current frame.
func (am *ActionMap) IsJustReleased(action string) bool {
	b, ok := am.bindings[action]
	if !ok {
		return false
	}
	for _, k := range b.Keys {
		if inpututil.IsKeyJustReleased(k) {
			return true
		}
	}
	for _, mb := range b.MouseButtons {
		if inpututil.IsMouseButtonJustReleased(mb) {
			return true
		}
	}
	return false
}

// Axis2D computes a normalized directional vector (-1.0 to 1.0) for four directional actions.
func (am *ActionMap) Axis2D(leftAction, rightAction, upAction, downAction string) (dx, dy float64) {
	if am.IsPressed(leftAction) {
		dx -= 1.0
	}
	if am.IsPressed(rightAction) {
		dx += 1.0
	}
	if am.IsPressed(upAction) {
		dy -= 1.0
	}
	if am.IsPressed(downAction) {
		dy += 1.0
	}

	// Normalize diagonals so diagonal movement doesn't exceed 1.0 speed
	if dx != 0 && dy != 0 {
		invSqrt2 := 0.7071067811865475
		dx *= invSqrt2
		dy *= invSqrt2
	}
	return dx, dy
}

// CursorCanvasPosition returns the current mouse position in virtual canvas space.
func CursorCanvasPosition(canvas *display.Canvas) (int, int, bool) {
	mx, my := ebiten.CursorPosition()
	return canvas.ScreenToCanvas(mx, my)
}
