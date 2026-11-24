package tilemap

import (
	"image"
	"image/color"
	"testing"

	"github.com/zapturk/PixelDream-2D/pkg/gfx"
)

func TestTilesetAndLayer(t *testing.T) {
	tex := gfx.NewSolidColorTexture(64, 64, color.White)
	ts := NewTileset(tex, 16, 16)

	if ts.Columns != 4 || ts.Rows != 4 || ts.TotalTiles() != 16 {
		t.Fatalf("expected 4x4=16 tiles, got %dx%d=%d", ts.Columns, ts.Rows, ts.TotalTiles())
	}

	// Tile 1 rect should be (0,0,16,16)
	r1 := ts.TileRect(1)
	if r1 != image.Rect(0, 0, 16, 16) {
		t.Fatalf("expected (0,0,16,16), got %v", r1)
	}

	// Tile 5 (first tile of second row) -> col 0, row 1 -> (0, 16, 16, 32)
	r5 := ts.TileRect(5)
	if r5 != image.Rect(0, 16, 16, 32) {
		t.Fatalf("expected (0,16,16,32), got %v", r5)
	}

	tm := NewTileMap(20, 15, ts)
	bg := tm.AddLayer("Background")
	solid := tm.AddLayer("Solid")

	solid.SetTile(5, 5, 1)
	if solid.GetTile(5, 5) != 1 {
		t.Fatalf("expected tile 1 at (5,5), got %d", solid.GetTile(5, 5))
	}

	if bg.GetTile(5, 5) != 0 {
		t.Fatalf("expected empty tile 0 in background at (5,5), got %d", bg.GetTile(5, 5))
	}
}

func TestCoordinateConversions(t *testing.T) {
	tex := gfx.NewSolidColorTexture(32, 32, color.White)
	ts := NewTileset(tex, 16, 16)
	tm := NewTileMap(10, 10, ts)

	// World (24, 35) -> Tile (1, 2)
	tx, ty := tm.WorldToTile(24, 35)
	if tx != 1 || ty != 2 {
		t.Fatalf("expected tile (1,2), got (%d,%d)", tx, ty)
	}

	wx, wy := tm.TileToWorld(2, 3)
	if wx != 32 || wy != 48 {
		t.Fatalf("expected world (32, 48), got (%f,%f)", wx, wy)
	}
}
