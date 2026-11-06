package game

const (
	tileShallowWater = 7
	tileDeepWater    = 8
)

// generateWoodsWater adds either a meandering edge-to-edge river or a pond.
func (g *Game) generateWoodsWater(m *Map) {
	if g.Roll(100) < 60 {
		g.generateWoodsRiver(m)
	} else {
		g.generateWoodsPond(m)
	}
}

func (g *Game) generateWoodsRiver(m *Map) {
	if g.Roll(2) == 0 {
		center := 12 + g.Roll(m.W-24)
		drift := 0
		for y := 0; y < m.H; y++ {
			if y%4 == 0 {
				drift += g.Roll(3) - 1
				drift = clampInt(drift, -5, 5)
			}
			g.paintRiverSlice(m, center+drift, y, true)
		}
		return
	}
	center := 8 + g.Roll(m.H-16)
	drift := 0
	for x := 0; x < m.W; x++ {
		if x%4 == 0 {
			drift += g.Roll(3) - 1
			drift = clampInt(drift, -4, 4)
		}
		g.paintRiverSlice(m, x, center+drift, false)
	}
}

func (g *Game) paintRiverSlice(m *Map, cx, cy int, vertical bool) {
	radius := 1 + g.Roll(2)
	for off := -radius - 2; off <= radius+2; off++ {
		x, y := cx, cy
		if vertical {
			x = cx + off
		} else {
			y = cy + off
		}
		if x < 0 || y < 0 || x >= m.W || y >= m.H {
			continue
		}
		switch {
		case abs(off) <= radius:
			m.SetTile(x, y, tileDeepWater)
		case abs(off) <= radius+1:
			if g.Roll(100) < 12 {
				m.SetTile(x, y, tileLilyPad)
			} else {
				m.SetTile(x, y, tileShallowWater)
			}
		case abs(off) <= radius+2:
			if g.Roll(100) < 22 {
				m.SetTile(x, y, tileReeds)
			} else {
				m.SetTile(x, y, tileRiverBank)
			}
		}
	}
}

func (g *Game) generateWoodsPond(m *Map) {
	cx := 16 + g.Roll(m.W-32)
	cy := 8 + g.Roll(m.H-16)
	rx := 4 + g.Roll(4)
	ry := 3 + g.Roll(3)
	for y := cy - ry - 3; y <= cy+ry+3; y++ {
		for x := cx - rx - 3; x <= cx+rx+3; x++ {
			if x < 1 || y < 1 || x >= m.W-1 || y >= m.H-1 {
				continue
			}
			dx := float64(x-cx) / float64(rx)
			dy := float64(y-cy) / float64(ry)
			d := dx*dx + dy*dy
			d += (g.RNG.Float64() - 0.5) * 0.35
			switch {
			case d < 0.55:
				m.SetTile(x, y, tileDeepWater)
			case d < 1.0:
				if g.Roll(100) < 8 {
					m.SetTile(x, y, tileLilyPad)
				} else {
					m.SetTile(x, y, tileShallowWater)
				}
			case d < 1.35:
				if g.Roll(100) < 20 {
					m.SetTile(x, y, tileReeds)
				} else {
					m.SetTile(x, y, tileRiverBank)
				}
			}
		}
	}
}

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
