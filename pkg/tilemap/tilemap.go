package tilemap

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zapturk/PixelDream-2D/pkg/display"
)

// Layer represents a single 2D grid plane of tile IDs.
type Layer struct {
	Name    string
	Visible bool
	Opacity float32
	Width   int
	Height  int
	tiles   []int // 1D array representing 2D grid: index = y * Width + x
}

// NewLayer creates an empty layer filled with TileTypeEmpty (0).
func NewLayer(name string, width, height int) *Layer {
	return &Layer{
		Name:    name,
		Visible: true,
		Opacity: 1.0,
		Width:   width,
		Height:  height,
		tiles:   make([]int, width*height),
	}
}

// GetTile returns the tile ID at grid coordinates (x, y). Returns 0 if out of bounds.
func (l *Layer) GetTile(x, y int) int {
	if x < 0 || x >= l.Width || y < 0 || y >= l.Height {
		return 0
	}
	return l.tiles[y*l.Width+x]
}

// SetTile sets the tile ID at grid coordinates (x, y).
func (l *Layer) SetTile(x, y, tileID int) {
	if x >= 0 && x < l.Width && y >= 0 && y < l.Height {
		l.tiles[y*l.Width+x] = tileID
	}
}

// Fill fills the entire layer or a rectangular region with tileID.
func (l *Layer) Fill(rect image.Rectangle, tileID int) {
	x1 := max(0, rect.Min.X)
	y1 := max(0, rect.Min.Y)
	x2 := min(l.Width, rect.Max.X)
	y2 := min(l.Height, rect.Max.Y)

	for y := y1; y < y2; y++ {
		rowOffset := y * l.Width
		for x := x1; x < x2; x++ {
			l.tiles[rowOffset+x] = tileID
		}
	}
}

// TileMap represents a multi-layered tile grid rendered via a Tileset.
type TileMap struct {
	Tileset    *Tileset
	Width      int // in tiles
	Height     int // in tiles
	TileWidth  int // in pixels
	TileHeight int // in pixels
	Layers     []*Layer
}

// NewTileMap creates a tilemap with the given grid size and tileset.
func NewTileMap(width, height int, tileset *Tileset) *TileMap {
	tw := 16
	th := 16
	if tileset != nil {
		tw = tileset.TileWidth
		th = tileset.TileHeight
	}

	return &TileMap{
		Tileset:    tileset,
		Width:      width,
		Height:     height,
		TileWidth:  tw,
		TileHeight: th,
		Layers:     make([]*Layer, 0),
	}
}

// AddLayer adds a new named layer to the tilemap.
func (tm *TileMap) AddLayer(name string) *Layer {
	layer := NewLayer(name, tm.Width, tm.Height)
	tm.Layers = append(tm.Layers, layer)
	return layer
}

// GetLayer returns a layer by name.
func (tm *TileMap) GetLayer(name string) *Layer {
	for _, l := range tm.Layers {
		if l.Name == name {
			return l
		}
	}
	return nil
}

// PixelWidth returns the full width of the tilemap in world pixels.
func (tm *TileMap) PixelWidth() int {
	return tm.Width * tm.TileWidth
}

// PixelHeight returns the full height of the tilemap in world pixels.
func (tm *TileMap) PixelHeight() int {
	return tm.Height * tm.TileHeight
}

// WorldToTile converts world coordinates into tile grid indices.
func (tm *TileMap) WorldToTile(worldX, worldY float64) (int, int) {
	tx := int(math.Floor(worldX / float64(tm.TileWidth)))
	ty := int(math.Floor(worldY / float64(tm.TileHeight)))
	return tx, ty
}

// TileToWorld converts tile grid coordinates into top-left world coordinates.
func (tm *TileMap) TileToWorld(tileX, tileY int) (float64, float64) {
	return float64(tileX * tm.TileWidth), float64(tileY * tm.TileHeight)
}

// Draw renders all visible layers with viewport culling based on the camera view.
func (tm *TileMap) Draw(canvas *display.Canvas, cam *display.Camera) {
	if tm.Tileset == nil || tm.Tileset.Texture == nil {
		return
	}

	// 1. Calculate camera visible bounds in world coordinates for culling
	var minTileX, minTileY, maxTileX, maxTileY int

	if cam != nil {
		// Visible world rectangle
		topLeftX, topLeftY := cam.ScreenToWorld(0, 0)
		botRightX, botRightY := cam.ScreenToWorld(float64(cam.ViewportWidth), float64(cam.ViewportHeight))

		minX := math.Min(topLeftX, botRightX)
		minY := math.Min(topLeftY, botRightY)
		maxX := math.Max(topLeftX, botRightX)
		maxY := math.Max(topLeftY, botRightY)

		// Expand by 1 tile buffer on all sides to avoid edge clipping
		minTileX = int(math.Floor(minX/float64(tm.TileWidth))) - 1
		minTileY = int(math.Floor(minY/float64(tm.TileHeight))) - 1
		maxTileX = int(math.Ceil(maxX/float64(tm.TileWidth))) + 1
		maxTileY = int(math.Ceil(maxY/float64(tm.TileHeight))) + 1
	} else {
		minTileX = 0
		minTileY = 0
		maxTileX = tm.Width
		maxTileY = tm.Height
	}

	// Clamp to tilemap bounds
	minTileX = max(0, minTileX)
	minTileY = max(0, minTileY)
	maxTileX = min(tm.Width, maxTileX)
	maxTileY = min(tm.Height, maxTileY)

	img := tm.Tileset.Texture.Image()

	var viewMat ebiten.GeoM
	if cam != nil {
		viewMat = cam.ViewMatrix()
	}

	// Draw each layer
	for _, layer := range tm.Layers {
		if !layer.Visible || layer.Opacity <= 0 {
			continue
		}

		for ty := minTileY; ty < maxTileY; ty++ {
			rowOffset := ty * layer.Width
			worldY := float64(ty * tm.TileHeight)

			for tx := minTileX; tx < maxTileX; tx++ {
				tileID := layer.tiles[rowOffset+tx]
				if tileID <= 0 {
					continue
				}

				srcRect := tm.Tileset.TileRect(tileID)
				if srcRect.Empty() {
					continue
				}

				worldX := float64(tx * tm.TileWidth)

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(worldX, worldY)
				if cam != nil {
					op.GeoM.Concat(viewMat)
				}

				if layer.Opacity < 1.0 {
					op.ColorScale.Scale(1, 1, 1, layer.Opacity)
				}
				op.Filter = ebiten.FilterNearest

				tileSub := img.SubImage(srcRect).(*ebiten.Image)
				canvas.Image().DrawImage(tileSub, op)
			}
		}
	}
}
