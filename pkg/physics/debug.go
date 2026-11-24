package physics

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/zapturk/PixelDream-2D/pkg/display"
)

// DrawAABB draws an outline and optional translucent fill for an AABB transformed into canvas space.
func DrawAABB(canvas *display.Canvas, cam *display.Camera, box AABB, strokeColor color.Color, fillColor color.Color) {
	if canvas == nil || canvas.Image() == nil {
		return
	}

	var screenX, screenY float64
	zoom := 1.0

	if cam != nil {
		screenX, screenY = cam.WorldToScreen(box.X, box.Y)
		zoom = cam.Zoom
	} else {
		screenX, screenY = box.X, box.Y
	}

	sx := float32(math.Round(screenX))
	sy := float32(math.Round(screenY))
	w := float32(math.Round(box.Width * zoom))
	h := float32(math.Round(box.Height * zoom))

	if fillColor != nil {
		vector.DrawFilledRect(canvas.Image(), sx, sy, w, h, fillColor, false)
	}
	if strokeColor != nil {
		vector.StrokeRect(canvas.Image(), sx, sy, w, h, 1.0, strokeColor, false)
	}
}

// DrawTileColliders draws debug outlines for all solid tiles within the specified tile grid range.
func DrawTileColliders(
	canvas *display.Canvas,
	cam *display.Camera,
	isSolid func(tx, ty int) bool,
	tileSize int,
	minTX, minTY, maxTX, maxTY int,
	strokeColor color.Color,
	fillColor color.Color,
) {
	if isSolid == nil {
		return
	}
	ts := float64(tileSize)
	for ty := minTY; ty < maxTY; ty++ {
		for tx := minTX; tx < maxTX; tx++ {
			if isSolid(tx, ty) {
				box := NewAABB(float64(tx)*ts, float64(ty)*ts, ts, ts)
				DrawAABB(canvas, cam, box, strokeColor, fillColor)
			}
		}
	}
}
