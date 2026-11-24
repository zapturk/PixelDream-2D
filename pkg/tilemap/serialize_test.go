package tilemap

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/zapturk/PixelDream-2D/pkg/gfx"
)

func TestTileMapJSONSerialization(t *testing.T) {
	tex := gfx.NewSolidColorTexture(32, 32, color.White)
	ts := NewTileset(tex, 16, 16)
	tm := NewTileMap(10, 10, ts)

	bg := tm.AddLayer("Background")
	bg.Fill(bg.BoundsRect(), 1)

	solid := tm.AddLayer("Solid")
	solid.SetTile(2, 3, 2)
	solid.SetTile(4, 5, 3)

	// Serialize
	jsonData, err := tm.ToJSON()
	if err != nil {
		t.Fatalf("failed to serialize: %v", err)
	}

	// Deserialize
	loadedMap, err := LoadFromJSON(jsonData, ts)
	if err != nil {
		t.Fatalf("failed to deserialize: %v", err)
	}

	if loadedMap.Width != 10 || loadedMap.Height != 10 {
		t.Fatalf("expected dimensions 10x10, got %dx%d", loadedMap.Width, loadedMap.Height)
	}
	if len(loadedMap.Layers) != 2 {
		t.Fatalf("expected 2 layers, got %d", len(loadedMap.Layers))
	}

	loadedSolid := loadedMap.GetLayer("Solid")
	if loadedSolid == nil {
		t.Fatalf("expected 'Solid' layer to exist")
	}
	if loadedSolid.GetTile(2, 3) != 2 || loadedSolid.GetTile(4, 5) != 3 {
		t.Fatalf("mismatched tile values in loaded solid layer")
	}

	// File test
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_map.json")
	if err := tm.SaveToFile(filePath); err != nil {
		t.Fatalf("failed to save to file: %v", err)
	}
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}

	fromFile, err := LoadFromFile(filePath, ts)
	if err != nil {
		t.Fatalf("failed to load from file: %v", err)
	}
	if fromFile.GetLayer("Solid").GetTile(2, 3) != 2 {
		t.Fatalf("failed to read tile from file-loaded map")
	}
}

// BoundsRect helper for Layer
func (l *Layer) BoundsRect() image.Rectangle {
	return image.Rect(0, 0, l.Width, l.Height)
}
