package tilemap

import (
	"encoding/json"
	"os"
)

// MapDataJSON is the serializable structure for a tilemap.
type MapDataJSON struct {
	Version    int             `json:"version"`
	Width      int             `json:"width"`
	Height     int             `json:"height"`
	TileWidth  int             `json:"tile_width"`
	TileHeight int             `json:"tile_height"`
	Layers     []LayerDataJSON `json:"layers"`
}

// LayerDataJSON is the serializable structure for a layer.
type LayerDataJSON struct {
	Name    string  `json:"name"`
	Visible bool    `json:"visible"`
	Opacity float32 `json:"opacity"`
	Tiles   []int   `json:"tiles"`
}

// ToJSON serializes the TileMap to JSON bytes.
func (tm *TileMap) ToJSON() ([]byte, error) {
	data := MapDataJSON{
		Version:    1,
		Width:      tm.Width,
		Height:     tm.Height,
		TileWidth:  tm.TileWidth,
		TileHeight: tm.TileHeight,
		Layers:     make([]LayerDataJSON, len(tm.Layers)),
	}

	for i, layer := range tm.Layers {
		// Copy tiles slice
		tilesCopy := make([]int, len(layer.tiles))
		copy(tilesCopy, layer.tiles)

		data.Layers[i] = LayerDataJSON{
			Name:    layer.Name,
			Visible: layer.Visible,
			Opacity: layer.Opacity,
			Tiles:   tilesCopy,
		}
	}

	return json.MarshalIndent(data, "", "  ")
}

// SaveToFile writes the tilemap to a JSON file.
func (tm *TileMap) SaveToFile(filePath string) error {
	bytes, err := tm.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, bytes, 0644)
}

// LoadFromJSON parses JSON bytes and populates the tilemap with layers.
func LoadFromJSON(data []byte, tileset *Tileset) (*TileMap, error) {
	var mapData MapDataJSON
	if err := json.Unmarshal(data, &mapData); err != nil {
		return nil, err
	}

	tm := NewTileMap(mapData.Width, mapData.Height, tileset)
	if mapData.TileWidth > 0 {
		tm.TileWidth = mapData.TileWidth
	}
	if mapData.TileHeight > 0 {
		tm.TileHeight = mapData.TileHeight
	}

	for _, lData := range mapData.Layers {
		layer := tm.AddLayer(lData.Name)
		layer.Visible = lData.Visible
		if lData.Opacity > 0 {
			layer.Opacity = lData.Opacity
		}
		if len(lData.Tiles) == tm.Width*tm.Height {
			copy(layer.tiles, lData.Tiles)
		}
	}

	return tm, nil
}

// LoadFromFile reads and reconstructs a TileMap from a JSON file.
func LoadFromFile(filePath string, tileset *Tileset) (*TileMap, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return LoadFromJSON(bytes, tileset)
}
