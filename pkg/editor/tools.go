package editor

import (
	"image"

	"github.com/zapturk/PixelDream-2D/pkg/tilemap"
)

// ToolType represents the active editing tool.
type ToolType int

const (
	ToolPencil ToolType = iota
	ToolEraser
	ToolRect
	ToolBucket
	ToolPicker
)

func (t ToolType) String() string {
	switch t {
	case ToolPencil:
		return "Pencil"
	case ToolEraser:
		return "Eraser"
	case ToolRect:
		return "Rect Fill"
	case ToolBucket:
		return "Bucket"
	case ToolPicker:
		return "Picker"
	default:
		return "Unknown"
	}
}

// ApplyPencil places a single tile at (tx, ty).
func ApplyPencil(layer *tilemap.Layer, tx, ty, tileID int) {
	if layer != nil {
		layer.SetTile(tx, ty, tileID)
	}
}

// ApplyEraser clears the tile at (tx, ty) to 0.
func ApplyEraser(layer *tilemap.Layer, tx, ty int) {
	if layer != nil {
		layer.SetTile(tx, ty, 0)
	}
}

// ApplyRect fills the rectangle between (x1, y1) and (x2, y2) with tileID.
func ApplyRect(layer *tilemap.Layer, x1, y1, x2, y2, tileID int) {
	if layer == nil {
		return
	}
	minX := min(x1, x2)
	maxX := max(x1, x2)
	minY := min(y1, y2)
	maxY := max(y1, y2)

	layer.Fill(image.Rect(minX, minY, maxX+1, maxY+1), tileID)
}

// ApplyBucket performs a 4-way flood fill on the layer.
func ApplyBucket(layer *tilemap.Layer, startX, startY, newTileID int) {
	if layer == nil {
		return
	}
	if startX < 0 || startX >= layer.Width || startY < 0 || startY >= layer.Height {
		return
	}

	targetTileID := layer.GetTile(startX, startY)
	if targetTileID == newTileID {
		return
	}

	type point struct{ x, y int }
	queue := []point{{startX, startY}}

	visited := make([]bool, layer.Width*layer.Height)
	visited[startY*layer.Width+startX] = true

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		layer.SetTile(p.x, p.y, newTileID)

		// 4 neighbors: Left, Right, Up, Down
		dirs := []point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, d := range dirs {
			nx := p.x + d.x
			ny := p.y + d.y
			if nx >= 0 && nx < layer.Width && ny >= 0 && ny < layer.Height {
				idx := ny*layer.Width + nx
				if !visited[idx] && layer.GetTile(nx, ny) == targetTileID {
					visited[idx] = true
					queue = append(queue, point{nx, ny})
				}
			}
		}
	}
}

// ApplyPicker returns the tile ID at (tx, ty).
func ApplyPicker(layer *tilemap.Layer, tx, ty int) int {
	if layer == nil {
		return 0
	}
	return layer.GetTile(tx, ty)
}
