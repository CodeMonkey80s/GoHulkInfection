package game

const (
	tileShallowWater = 7
	tileDeepWater    = 8
)

func isWoodsWater(t uint8) bool {
	switch t {
	case tileDeepWater, tileShallowWater, tileLilyPad:
		return true
	}
	return false
}

func (g *Game) setWoodsRoadTile(m *Map, x, y int) {
	if isWoodsWater(m.Tile(x, y)) {
		m.SetTile(x, y, tileWoodenBridge)
		return
	}
	m.SetTile(x, y, tileDirtRoad)
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
