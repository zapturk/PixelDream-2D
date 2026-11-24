package tilemap

import (
	"image"

	"github.com/zapturk/PixelDream-2D/pkg/gfx"
)

// TileType defines the collision or behavior class of a tile.
type TileType int

const (
	// TileTypeEmpty represents no tile / air.
	TileTypeEmpty TileType = 0
	// TileTypeSolid is fully blocking on all sides.
	TileTypeSolid TileType = 1
	// TileTypeOneWay allows passing through from bottom/sides, but solid when landing from above.
	TileTypeOneWay TileType = 2
	// TileTypeSlopeLeft is a 45-degree slope rising to the right (ground goes from low left to high right).
	TileTypeSlopeRight TileType = 3
	// TileTypeSlopeLeft is a 45-degree slope rising to the left (ground goes from low right to high left).
	TileTypeSlopeLeft TileType = 4
	// TileTypeHazard hurts the player on contact (e.g. spikes, lava).
	TileTypeHazard TileType = 5
)

// TileDef defines the visual and collision attributes of a tile in a tileset.
type TileDef struct {
	ID        int
	Rect      image.Rectangle
	Type      TileType
	Tag       string
	UserData  any
}

// Tileset manages a texture atlas divided into uniform grid tiles.
type Tileset struct {
	Texture    *gfx.Texture
	TileWidth  int
	TileHeight int
	Columns    int
	Rows       int
	tileDefs   map[int]*TileDef
}

// NewTileset creates a tileset slicing a texture into grid tiles of tileWidth x tileHeight.
func NewTileset(texture *gfx.Texture, tileWidth, tileHeight int) *Tileset {
	if tileWidth <= 0 {
		tileWidth = 16
	}
	if tileHeight <= 0 {
		tileHeight = 16
	}

	cols := 0
	rows := 0
	if texture != nil {
		cols = texture.Width() / tileWidth
		rows = texture.Height() / tileHeight
	}

	ts := &Tileset{
		Texture:    texture,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
		Columns:    cols,
		Rows:       rows,
		tileDefs:   make(map[int]*TileDef),
	}

	// Auto-populate default tile definitions based on grid
	total := cols * rows
	for id := 1; id <= total; id++ {
		col := (id - 1) % cols
		row := (id - 1) / cols
		rect := image.Rect(col*tileWidth, row*tileHeight, (col+1)*tileWidth, (row+1)*tileHeight)
		ts.tileDefs[id] = &TileDef{
			ID:   id,
			Rect: rect,
			Type: TileTypeSolid, // Default to solid, can be customized
		}
	}

	return ts
}

// SetTileType configures the collision type for a specific tile ID.
func (ts *Tileset) SetTileType(tileID int, tType TileType) {
	def, ok := ts.tileDefs[tileID]
	if !ok {
		def = &TileDef{ID: tileID}
		ts.tileDefs[tileID] = def
	}
	def.Type = tType
}

// GetTileDef returns the definition for a tile ID.
func (ts *Tileset) GetTileDef(tileID int) *TileDef {
	return ts.tileDefs[tileID]
}

// TileRect returns the pixel source rectangle for a tile ID.
func (ts *Tileset) TileRect(tileID int) image.Rectangle {
	if def, ok := ts.tileDefs[tileID]; ok {
		return def.Rect
	}
	return image.Rectangle{}
}

// TotalTiles returns the total number of sliced tiles available in the tileset.
func (ts *Tileset) TotalTiles() int {
	return ts.Columns * ts.Rows
}
