package pixeleditor

import (
	"image/color"
)

// Tool defines the active pixel drawing tool.
type Tool int

const (
	ToolPencil Tool = iota
	ToolEraser
	ToolPicker
	ToolBucket
	ToolLine
	ToolRect
)

func (t Tool) String() string {
	switch t {
	case ToolPencil:
		return "Pencil"
	case ToolEraser:
		return "Eraser"
	case ToolPicker:
		return "Eyedropper"
	case ToolBucket:
		return "Bucket Fill"
	case ToolLine:
		return "Line"
	case ToolRect:
		return "Rectangle"
	default:
		return "Unknown"
	}
}

// BresenhamLine executes callback for every pixel along the line from (x0, y0) to (x1, y1).
func BresenhamLine(x0, y0, x1, y1 int, plot func(x, y int)) {
	dx := absInt(x1 - x0)
	dy := -absInt(y1 - y0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	err := dx + dy

	for {
		plot(x0, y0)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// FloodFill fills connected pixels of matching color with newColor.
func FloodFill(doc *Document, startX, startY int, newColor color.RGBA) {
	if doc == nil || startX < 0 || startX >= doc.Width || startY < 0 || startY >= doc.Height {
		return
	}

	targetColor := doc.GetPixel(startX, startY)
	if targetColor == newColor {
		return
	}

	type pt struct{ x, y int }
	queue := []pt{{startX, startY}}

	visited := make([]bool, doc.Width*doc.Height)
	visited[startY*doc.Width+startX] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		doc.SetPixel(curr.x, curr.y, newColor)

		dirs := []pt{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, d := range dirs {
			nx := curr.x + d.x
			ny := curr.y + d.y
			if nx >= 0 && nx < doc.Width && ny >= 0 && ny < doc.Height {
				idx := ny*doc.Width + nx
				if !visited[idx] && doc.GetPixel(nx, ny) == targetColor {
					visited[idx] = true
					queue = append(queue, pt{nx, ny})
				}
			}
		}
	}
}

// DrawStrokeLine draws a line of pixels between two points using Bresenham algorithm.
func DrawStrokeLine(doc *Document, x0, y0, x1, y1 int, col color.RGBA) {
	BresenhamLine(x0, y0, x1, y1, func(x, y int) {
		doc.SetPixel(x, y, col)
	})
}

// DrawRectStroke draws an outlined or filled rectangle.
func DrawRectStroke(doc *Document, x0, y0, x1, y1 int, col color.RGBA, filled bool) {
	minX := min(x0, x1)
	maxX := max(x0, x1)
	minY := min(y0, y1)
	maxY := max(y0, y1)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if filled || x == minX || x == maxX || y == minY || y == maxY {
				doc.SetPixel(x, y, col)
			}
		}
	}
}
