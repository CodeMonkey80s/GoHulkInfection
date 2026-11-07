package game

import "math"

type woodsThemeSpec struct {
	HutMin int
	HutMax int
}

var woodsThemeSpecs = [numWoodsMaps]woodsThemeSpec{
	WoodsChurch:        {HutMin: 1, HutMax: 2},
	WoodsPeasantHouses: {HutMin: 3, HutMax: 4},
	WoodsMillPond:      {HutMin: 0, HutMax: 1},
	WoodsCemetery:      {HutMin: 0, HutMax: 1},
	WoodsMeteorField:   {HutMin: 0, HutMax: 0},
}

func woodThemeSpecFor(theme int) woodsThemeSpec {
	if theme < 0 || theme >= len(woodsThemeSpecs) {
		return woodsThemeSpecs[0]
	}
	return woodsThemeSpecs[theme]
}

func isWoodsStructure(t uint8) bool {
	switch t {
	case 9, 10, 12, 13, 14, 15, 16, 17, 18, 20, 21, 22, 23, 24,
		tileTable, tileChair, tileBed, tileCrate, tileBarrel, tileFireplace, tileShelf,
		tileOrchardTree, tileBeehive, tileCiderPress, tileGravestone, tileIronFence,
		tileCryptDoor, tileMillStone, tileTent, tileCampfire, tileCage, tileDryingRack,
		tileMushroomRing, tileDruidAltar, tileWell, tileCart, tileMarketStall,
		tileHaystack, tileKiln, tileChapelAltar, tileBellTower, tileWoodpile,
		tileMeteorRock:
		return true
	}
	return false
}

func (g *Game) woodsLandmarkCenter(m *Map, theme int) (int, int) {
	cx := 34 + g.Roll(13)
	cy := 9 + g.Roll(5)
	if theme == WoodsCemetery {
		cy = 11
	}
	return clampInt(cx, 26, 53), clampInt(cy, 8, 13)
}

func sideTowardEdge(edge int) [2]int {
	switch edge {
	case EdgeNorth:
		return [2]int{0, -1}
	case EdgeSouth:
		return [2]int{0, 1}
	case EdgeEast:
		return [2]int{1, 0}
	default:
		return [2]int{-1, 0}
	}
}

func clearWoodsRect(m *Map, x0, y0, x1, y1 int) {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x < 2 || y < 2 || x >= m.W-2 || y >= m.H-2 {
				continue
			}
			m.SetTile(x, y, 2)
		}
	}
}

func buildWoodsBuilding(m *Map, cx, cy, w, h int, wall, floor, door uint8, doorSide [2]int) [2]int {
	x0 := cx - w/2
	y0 := cy - h/2
	x1 := x0 + w - 1
	y1 := y0 + h - 1
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x == x0 || x == x1 || y == y0 || y == y1 {
				m.SetTile(x, y, wall)
			} else {
				m.SetTile(x, y, floor)
			}
		}
	}
	doorX, doorY := cx, cy
	switch doorSide {
	case [2]int{0, -1}:
		doorX, doorY = cx, y0
	case [2]int{0, 1}:
		doorX, doorY = cx, y1
	case [2]int{-1, 0}:
		doorX, doorY = x0, cy
	default:
		doorX, doorY = x1, cy
	}
	m.SetTile(doorX, doorY, door)
	return [2]int{doorX + doorSide[0], doorY + doorSide[1]}
}

func (g *Game) buildWoodsThemeLandmark(m *Map, theme, cx, cy, entryEdge int) [2]int {
	switch theme {
	case WoodsChurch:
		return g.buildChurchLandmark(m, cx, cy, entryEdge)
	case WoodsPeasantHouses:
		return g.buildPeasantHousesLandmark(m, cx, cy, entryEdge)
	case WoodsMillPond:
		return g.buildMillLandmark(m, cx, cy, entryEdge)
	case WoodsCemetery:
		return g.buildCemeteryLandmark(m, cx, cy, entryEdge)
	default:
		return g.buildMeteorLandmark(m, cx, cy, entryEdge)
	}
}

func (g *Game) buildChurchLandmark(m *Map, cx, cy, entryEdge int) [2]int {
	clearWoodsRect(m, cx-7, cy-5, cx+7, cy+5)
	side := sideTowardEdge(entryEdge)
	anchor := buildWoodsBuilding(m, cx, cy, 9, 7, 10, 9, tileWoodenDoor, side)
	m.SetTile(cx-3, cy-2, tileBellTower)
	altarX := cx - side[0]*3
	altarY := cy - side[1]*2
	m.SetTile(altarX, altarY, tileChapelAltar)
	for i := 0; i < 8; i++ {
		x, y := cx-7+g.Roll(15), cy-5+g.Roll(11)
		if m.Tile(x, y) == 2 && !(x == anchor[0] && y == anchor[1]) {
			m.SetTile(x, y, tileGravestone)
		}
	}
	for i := 0; i < 2; i++ {
		x, y := cx-7+g.Roll(15), cy-5+g.Roll(11)
		if m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileDeadTree)
		}
	}
	m.LandmarkX, m.LandmarkY = cx, cy
	return anchor
}

func (g *Game) buildPeasantHousesLandmark(m *Map, cx, cy, entryEdge int) [2]int {
	clearWoodsRect(m, cx-8, cy-5, cx+8, cy+5)
	m.SetTile(cx, cy, tileWell)
	side := sideTowardEdge(entryEdge)
	anchor := [2]int{cx + side[0], cy + side[1]}
	if m.Tile(anchor[0], anchor[1]) != 2 {
		m.SetTile(anchor[0], anchor[1], 2)
	}
	for _, p := range [][2]int{{cx - 3, cy - 2}, {cx + 3, cy - 2}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileMarketStall)
		}
	}
	for _, p := range [][2]int{{cx - 3, cy + 3}, {cx + 3, cy + 3}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileCart)
		}
	}
	for _, p := range [][2]int{{cx - 6, cy}, {cx + 6, cy}, {cx, cy - 4}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileHaystack)
		}
	}
	for i := 0; i < 6; i++ {
		x, y := cx-7+g.Roll(15), cy-4+g.Roll(9)
		if m.Tile(x, y) == 2 {
			if g.Roll(2) == 0 {
				m.SetTile(x, y, tileFlowers)
			} else {
				m.SetTile(x, y, tileBush)
			}
		}
	}
	m.LandmarkX, m.LandmarkY = cx, cy
	return anchor
}

func (g *Game) buildMillLandmark(m *Map, cx, cy, entryEdge int) [2]int {
	clearWoodsRect(m, cx-8, cy-5, cx+8, cy+5)
	side := sideTowardEdge(entryEdge)
	anchor := buildWoodsBuilding(m, cx, cy, 8, 5, tileWoodenWall, tileWoodenFloor, tileWoodenDoor, side)
	for _, p := range [][2]int{{cx + 4, cy - 2}, {cx - 4, cy + 2}, {cx + 4, cy + 2}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileMillStone)
		}
	}
	for _, p := range [][2]int{{cx - 5, cy - 2}, {cx + 5, cy}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileBarrel)
		}
	}
	if m.Tile(cx-5, cy+2) == 2 {
		m.SetTile(cx-5, cy+2, tileCart)
	}
	if m.Tile(cx+2, cy-3) == 2 {
		m.SetTile(cx+2, cy-3, tileCrate)
	}
	m.LandmarkX, m.LandmarkY = cx, cy
	return anchor
}

func (g *Game) buildCemeteryLandmark(m *Map, cx, cy, entryEdge int) [2]int {
	w, h := 15, 11
	x0, x1 := cx-w/2, cx-w/2+w-1
	y0, y1 := cy-h/2, cy-h/2+h-1
	clearWoodsRect(m, x0-2, y0-2, x1+2, y1+2)
	for x := x0; x <= x1; x++ {
		m.SetTile(x, y0, tileIronFence)
		m.SetTile(x, y1, tileIronFence)
	}
	for y := y0; y <= y1; y++ {
		m.SetTile(x0, y, tileIronFence)
		m.SetTile(x1, y, tileIronFence)
	}
	side := sideTowardEdge(entryEdge)
	var gateX, gateY int
	switch side {
	case [2]int{0, -1}:
		gateX, gateY = cx, y0
	case [2]int{0, 1}:
		gateX, gateY = cx, y1
	case [2]int{-1, 0}:
		gateX, gateY = x0, cy
	default:
		gateX, gateY = x1, cy
	}
	m.SetTile(gateX, gateY, 2)
	anchor := [2]int{gateX + side[0], gateY + side[1]}

	cryptCX, cryptCY := cx, cy
	buildWoodsBuilding(m, cryptCX, cryptCY, 7, 5, 10, 9, tileCryptDoor, side)
	for y := y0 + 2; y <= y1-2; y += 2 {
		for x := x0 + 2; x <= x1-2; x += 2 {
			if m.Tile(x, y) != 2 {
				continue
			}
			if abs(x-gateX)+abs(y-gateY) < 3 {
				continue
			}
			m.SetTile(x, y, tileGravestone)
		}
	}
	for i := 0; i < 6; i++ {
		x, y := x0-3+g.Roll(w+6), y0-3+g.Roll(h+6)
		if m.InBounds(x, y) && m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileDeadTree)
		}
	}
	m.LandmarkX, m.LandmarkY = cryptCX, cryptCY
	return anchor
}

func (g *Game) buildMeteorLandmark(m *Map, cx, cy, entryEdge int) [2]int {
	rx, ry := 9, 4
	clearWoodsRect(m, cx-rx-2, cy-ry-2, cx+rx+2, cy+ry+2)
	for y := cy - ry; y <= cy+ry; y++ {
		for x := cx - rx; x <= cx+rx; x++ {
			dx := float64(x-cx) / float64(rx)
			dy := float64(y-cy) / float64(ry)
			d := dx*dx + dy*dy
			if d <= 0.25 {
				m.SetTile(x, y, tileMud)
			} else if d <= 1.0 {
				m.SetTile(x, y, tileBurnedGround)
			}
		}
	}
	for i := 0; i < 9; i++ {
		a := float64(i) * 2 * math.Pi / 9
		x := cx + int(math.Round(math.Cos(a)*float64(rx)*0.78))
		y := cy + int(math.Round(math.Sin(a)*float64(ry)*0.78))
		if m.InBounds(x, y) {
			m.SetTile(x, y, tileMeteorRock)
		}
	}
	for i := 0; i < 7; i++ {
		x, y := cx-rx-3+g.Roll(2*rx+7), cy-ry-3+g.Roll(2*ry+7)
		if m.InBounds(x, y) && m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileDeadTree)
		}
	}
	side := sideTowardEdge(entryEdge)
	anchor := [2]int{cx + side[0]*(rx+1), cy + side[1]*(ry+1)}
	if !m.InBounds(anchor[0], anchor[1]) || m.Tile(anchor[0], anchor[1]) != 2 {
		anchor = [2]int{cx + rx + 1, cy}
	}
	m.LandmarkX, m.LandmarkY = cx, cy
	return anchor
}
