package editor

import (
	"testing"

	"github.com/zapturk/PixelDream-2D/pkg/tilemap"
)

func TestEditorTools(t *testing.T) {
	layer := tilemap.NewLayer("TestLayer", 10, 10)

	// Pencil
	ApplyPencil(layer, 2, 2, 5)
	if layer.GetTile(2, 2) != 5 {
		t.Fatalf("expected pencil to set tile 5 at (2,2), got %d", layer.GetTile(2, 2))
	}

	// Eraser
	ApplyEraser(layer, 2, 2)
	if layer.GetTile(2, 2) != 0 {
		t.Fatalf("expected eraser to set tile 0 at (2,2), got %d", layer.GetTile(2, 2))
	}

	// Rect Fill (from 1,1 to 3,3) -> 3x3 = 9 tiles
	ApplyRect(layer, 1, 1, 3, 3, 7)
	for y := 1; y <= 3; y++ {
		for x := 1; x <= 3; x++ {
			if layer.GetTile(x, y) != 7 {
				t.Fatalf("expected tile 7 at (%d,%d), got %d", x, y, layer.GetTile(x, y))
			}
		}
	}

	// Picker
	picked := ApplyPicker(layer, 2, 2)
	if picked != 7 {
		t.Fatalf("expected picker to return 7, got %d", picked)
	}

	// Bucket fill
	ApplyBucket(layer, 2, 2, 9)
	for y := 1; y <= 3; y++ {
		for x := 1; x <= 3; x++ {
			if layer.GetTile(x, y) != 9 {
				t.Fatalf("expected bucket filled tile 9 at (%d,%d), got %d", x, y, layer.GetTile(x, y))
			}
		}
	}
	// Surrounding tile (0, 0) should still be 0
	if layer.GetTile(0, 0) != 0 {
		t.Fatalf("expected (0,0) to remain untouched 0, got %d", layer.GetTile(0, 0))
	}
}
