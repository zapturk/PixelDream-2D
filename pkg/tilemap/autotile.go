package tilemap

// AutotileBitmask returns a 4-bit bitmask (0-15) indicating solid neighbors:
// North = 1, East = 2, South = 4, West = 8.
func Autotile4Bit(layer *Layer, x, y int, solidTileID func(int) bool) int {
	mask := 0
	// North
	if solidTileID(layer.GetTile(x, y-1)) {
		mask |= 1
	}
	// East
	if solidTileID(layer.GetTile(x+1, y)) {
		mask |= 2
	}
	// South
	if solidTileID(layer.GetTile(x, y+1)) {
		mask |= 4
	}
	// West
	if solidTileID(layer.GetTile(x-1, y)) {
		mask |= 8
	}
	return mask
}
