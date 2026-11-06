package game

import "math"

type woodLevelSpec struct {
	Name   string
	HutMin int
	HutMax int
}

var woodLevelSpecs = []woodLevelSpec{
	{Name: "The Woodcutter's Clearing", HutMin: 2, HutMax: 6},
	{Name: "The Old Orchard", HutMin: 2, HutMax: 6},
	{Name: "The Forgotten Cemetery", HutMin: 2, HutMax: 6},
	{Name: "The River Mill", HutMin: 2, HutMax: 6},
	{Name: "The Poacher's Camp", HutMin: 2, HutMax: 6},
	{Name: "The Druid's Grove", HutMin: 2, HutMax: 6},
	{Name: "The Abandoned Village", HutMin: 5, HutMax: 8},
	{Name: "The Charcoal Burner's Clearing", HutMin: 2, HutMax: 6},
	{Name: "The Ruined Chapel", HutMin: 2, HutMax: 4},
}

func woodLevelSpecFor(level int) woodLevelSpec {
	if level < 0 || level >= len(woodLevelSpecs) {
		return woodLevelSpecs[0]
	}
	return woodLevelSpecs[level]
}

func isWoodsStructure(t uint8) bool {
	switch t {
	case 9, 10, 12, 13, 14, 15, 16, 17, 18, 20, 21, 22, 23, 24,
		tileTable, tileChair, tileBed, tileCrate, tileBarrel, tileFireplace, tileShelf,
		tileOrchardTree, tileBeehive, tileCiderPress, tileGravestone, tileIronFence,
		tileCryptDoor, tileMillStone, tileTent, tileCampfire, tileCage, tileDryingRack,
		tileMushroomRing, tileDruidAltar, tileWell, tileCart, tileMarketStall,
		tileHaystack, tileKiln, tileChapelAltar, tileBellTower, tileWoodpile:
		return true
	}
	return false
}

func (g *Game) placeWoodsLandmark(m *Map, level int, path [][2]int, startX, startY, exitX, exitY int) {
	switch level {
	case 0:
		g.buildWoodcutterLandmark(m, path)
	case 1:
		g.buildOrchardLandmark(m, path)
	case 2:
		g.buildCemeteryLandmark(m, path)
	case 3:
		g.buildRiverMillLandmark(m, path)
	case 4:
		g.buildPoacherCampLandmark(m, path)
	case 5:
		g.buildDruidGroveLandmark(m, path)
	case 6:
		g.buildVillageLandmark(m, path)
	case 7:
		g.buildCharcoalLandmark(m, path)
	case 8:
		g.buildRuinedChapelAtExit(m, path, exitX, exitY)
	}
}

func landmarkAreaOK(m *Map, x0, y0, x1, y1 int) (bool, bool) {
	hasWater := false
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			t := m.Tile(x, y)
			if t == tileDirtRoad || t == tileWoodenBridge || t == portalTile(0) || isWoodsStructure(t) {
				return false, false
			}
			if t == tileRiverBank || isWoodsWater(t) {
				hasWater = true
			}
		}
	}
	return true, hasWater
}

func (g *Game) reserveLandmarkArea(m *Map, path [][2]int, w, h int) (int, int, [2]int, bool) {
	mid := len(path) / 2
	var fallbackX, fallbackY int
	var fallbackSide [2]int
	fallbackOK := false
	for off := 0; off < len(path)/2; off++ {
		for _, idx := range []int{mid + off, mid - off} {
			if idx < 6 || idx >= len(path)-6 {
				continue
			}
			p := path[idx]
			sides := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
			g.RNG.Shuffle(len(sides), func(a, b int) {
				sides[a], sides[b] = sides[b], sides[a]
			})
			for _, side := range sides {
				cx := p[0] + side[0]*(w/2+2)
				cy := p[1] + side[1]*(h/2+2)
				x0, x1 := cx-w/2-1, cx+w/2+1
				y0, y1 := cy-h/2-1, cy+h/2+1
				if x0 < 2 || y0 < 2 || x1 >= m.W-2 || y1 >= m.H-2 {
					continue
				}
				ok, hasWater := landmarkAreaOK(m, x0, y0, x1, y1)
				if !ok {
					continue
				}
				if !hasWater {
					clearWoodsArea(m, cx, cy, w/2+1, 0)
					return cx, cy, side, true
				}
				if !fallbackOK {
					fallbackX, fallbackY, fallbackSide, fallbackOK = cx, cy, side, true
				}
			}
		}
	}
	if fallbackOK {
		clearWoodsArea(m, fallbackX, fallbackY, w/2+1, 0)
		return fallbackX, fallbackY, fallbackSide, true
	}

	// Last resort: scan the whole map for a dry area away from the road.
	for y0 := 2; y0 <= m.H-h-2; y0++ {
		for x0 := 2; x0 <= m.W-w-2; x0++ {
			ok, hasWater := landmarkAreaOK(m, x0, y0, x0+w+1, y0+h+1)
			if !ok || hasWater {
				continue
			}
			cx, cy := x0+w/2, y0+h/2
			clearWoodsArea(m, cx, cy, w/2+1, 0)
			return cx, cy, [2]int{0, 1}, true
		}
	}
	return 0, 0, [2]int{}, false
}

func buildLandmarkBuilding(m *Map, cx, cy, w, h int, side [2]int, wall, floor, door uint8) {
	x0, x1 := cx-w/2, cx-w/2+w-1
	y0, y1 := cy-h/2, cy-h/2+h-1
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x == x0 || x == x1 || y == y0 || y == y1 {
				m.SetTile(x, y, wall)
			} else {
				m.SetTile(x, y, floor)
			}
		}
	}
	var dx, dy int
	switch {
	case side[1] == -1:
		dx, dy = cx, y1
	case side[1] == 1:
		dx, dy = cx, y0
	case side[0] == -1:
		dx, dy = x1, cy
	default:
		dx, dy = x0, cy
	}
	m.SetTile(dx, dy, door)
}

func (g *Game) buildWoodcutterLandmark(m *Map, path [][2]int) {
	cx, cy, side, ok := g.reserveLandmarkArea(m, path, 9, 7)
	if !ok {
		return
	}
	buildLandmarkBuilding(m, cx, cy, 5, 5, side, tileWoodenWall, tileWoodenFloor, tileWoodenDoor)
	m.SetTile(cx+3, cy, tileWoodpile)
	for i := 0; i < 4; i++ {
		x, y := cx-3+g.Roll(7), cy-2+g.Roll(5)
		if m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileWoodpile)
		}
	}
	for i := 0; i < 7; i++ {
		x, y := cx-4+g.Roll(9), cy-3+g.Roll(7)
		if m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileStump)
		}
	}
}

func (g *Game) buildOrchardLandmark(m *Map, path [][2]int) {
	cx, cy, _, ok := g.reserveLandmarkArea(m, path, 11, 9)
	if !ok {
		return
	}
	for y := cy - 3; y <= cy+3; y += 2 {
		for x := cx - 4; x <= cx+4; x += 2 {
			if m.Tile(x, y) == 2 {
				m.SetTile(x, y, tileOrchardTree)
			}
		}
	}
	m.SetTile(cx, cy, tileCiderPress)
	for _, p := range [][2]int{{cx - 4, cy - 3}, {cx + 4, cy - 3}, {cx - 4, cy + 3}, {cx + 4, cy + 3}} {
		m.SetTile(p[0], p[1], tileBeehive)
	}
}

func (g *Game) buildCemeteryLandmark(m *Map, path [][2]int) {
	cx, cy, side, ok := g.reserveLandmarkArea(m, path, 13, 11)
	if !ok {
		return
	}
	x0, x1 := cx-6, cx+6
	y0, y1 := cy-5, cy+5
	for x := x0; x <= x1; x++ {
		m.SetTile(x, y0, tileIronFence)
		m.SetTile(x, y1, tileIronFence)
	}
	for y := y0; y <= y1; y++ {
		m.SetTile(x0, y, tileIronFence)
		m.SetTile(x1, y, tileIronFence)
	}
	gateX, gateY := cx, y0
	switch {
	case side[1] == -1:
		gateX, gateY = cx, y1
	case side[1] == 1:
		gateX, gateY = cx, y0
	case side[0] == -1:
		gateX, gateY = x1, cy
	default:
		gateX, gateY = x0, cy
	}
	m.SetTile(gateX, gateY, 2)
	if gateY == y0 {
		m.SetTile(gateX, y0+1, 2)
	} else if gateY == y1 {
		m.SetTile(gateX, y1-1, 2)
	} else if gateX == x0 {
		m.SetTile(x0+1, gateY, 2)
	} else {
		m.SetTile(x1-1, gateY, 2)
	}
	for y := y0 + 2; y <= y1-2; y += 2 {
		for x := x0 + 2; x <= x1-2; x += 2 {
			if x == gateX || m.Tile(x, y) != 2 {
				continue
			}
			m.SetTile(x, y, tileGravestone)
		}
	}
	buildLandmarkBuilding(m, cx-3, cy-2, 5, 5, side, 10, 9, tileCryptDoor)
	for i := 0; i < 4; i++ {
		x, y := cx-7+g.Roll(15), cy-6+g.Roll(13)
		if m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileDeadTree)
		}
	}
}

func (g *Game) buildRiverMillLandmark(m *Map, path [][2]int) {
	cx, cy, side, ok := g.reserveLandmarkArea(m, path, 9, 7)
	if !ok {
		return
	}
	buildLandmarkBuilding(m, cx, cy, 6, 5, side, tileWoodenWall, tileWoodenFloor, tileWoodenDoor)
	for _, p := range [][2]int{{cx - 4, cy}, {cx + 4, cy}, {cx, cy - 3}, {cx, cy + 3}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileMillStone)
			break
		}
	}
}

func (g *Game) buildPoacherCampLandmark(m *Map, path [][2]int) {
	cx, cy, _, ok := g.reserveLandmarkArea(m, path, 11, 9)
	if !ok {
		return
	}
	m.SetTile(cx, cy, tileCampfire)
	for _, p := range [][2]int{{cx - 3, cy - 2}, {cx + 3, cy - 2}, {cx, cy + 3}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileTent)
		}
	}
	for _, p := range [][2]int{{cx - 3, cy + 2}, {cx + 3, cy + 2}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileCage)
		}
	}
	for _, p := range [][2]int{{cx - 1, cy - 3}, {cx + 1, cy - 3}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileDryingRack)
		}
	}
}

func (g *Game) buildDruidGroveLandmark(m *Map, path [][2]int) {
	cx, cy, _, ok := g.reserveLandmarkArea(m, path, 13, 11)
	if !ok {
		return
	}
	for i := 0; i < 8; i++ {
		a := float64(i) * math.Pi / 4
		x := cx + int(math.Round(math.Cos(a)*4))
		y := cy + int(math.Round(math.Sin(a)*3))
		m.SetTile(x, y, 12)
	}
	for i := 0; i < 8; i++ {
		a := float64(i) * math.Pi / 4
		x := cx + int(math.Round(math.Cos(a)*2))
		y := cy + int(math.Round(math.Sin(a)*2))
		if m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileMushroomRing)
		}
	}
	m.SetTile(cx, cy, tileDruidAltar)
	for _, p := range [][2]int{{cx - 5, cy - 4}, {cx + 5, cy - 4}, {cx - 5, cy + 4}, {cx + 5, cy + 4}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], 4) // giant tree
		}
	}
}

func (g *Game) buildVillageLandmark(m *Map, path [][2]int) {
	cx, cy, _, ok := g.reserveLandmarkArea(m, path, 13, 11)
	if !ok {
		return
	}
	m.SetTile(cx, cy, tileWell)
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
	for _, p := range [][2]int{{cx - 5, cy}, {cx + 5, cy}} {
		if m.Tile(p[0], p[1]) == 2 {
			m.SetTile(p[0], p[1], tileHaystack)
		}
	}
}

func (g *Game) buildCharcoalLandmark(m *Map, path [][2]int) {
	cx, cy, _, ok := g.reserveLandmarkArea(m, path, 11, 9)
	if !ok {
		return
	}
	for i := 0; i < 3; i++ {
		x := cx - 3 + i*3
		if m.Tile(x, cy) == 2 {
			m.SetTile(x, cy, tileKiln)
			if m.InBounds(x+1, cy) && m.Tile(x+1, cy) == 2 {
				m.Fire[m.idx(x+1, cy)] = 20 + g.Roll(20)
			}
			if m.InBounds(x, cy-1) && m.Tile(x, cy-1) == 2 {
				m.Smoke[m.idx(x, cy-1)] = 20 + g.Roll(20)
			}
		}
	}
	for i := 0; i < 6; i++ {
		x, y := cx-5+g.Roll(11), cy-4+g.Roll(9)
		if m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileDeadTree)
		}
	}
}

func (g *Game) buildRuinedChapelAtExit(m *Map, path [][2]int, exitX, exitY int) {
	w, h := 11, 9
	x0, x1 := exitX-w/2, exitX-w/2+w-1
	y0, y1 := exitY-h/2, exitY-h/2+h-1
	if x0 < 2 {
		x0, x1 = 2, 2+w-1
	}
	if x1 >= m.W-2 {
		x1, x0 = m.W-3, m.W-3-w+1
	}
	if y0 < 2 {
		y0, y1 = 2, 2+h-1
	}
	if y1 >= m.H-2 {
		y1, y0 = m.H-3, m.H-3-h+1
	}
	openSide := [2]int{0, 1}
	if len(path) >= 2 {
		prev := path[len(path)-2]
		dx, dy := exitX-prev[0], exitY-prev[1]
		if abs(dx) > abs(dy) {
			if dx > 0 {
				openSide = [2]int{-1, 0}
			} else {
				openSide = [2]int{1, 0}
			}
		} else if dy > 0 {
			openSide = [2]int{0, -1}
		} else {
			openSide = [2]int{0, 1}
		}
	}
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if m.Tile(x, y) == tileDirtRoad || m.Tile(x, y) == portalTile(0) {
				continue
			}
			wall := x == x0 || x == x1 || y == y0 || y == y1
			if wall {
				if (openSide[0] == 0 && openSide[1] == -1 && y == y0) ||
					(openSide[0] == 0 && openSide[1] == 1 && y == y1) ||
					(openSide[0] == -1 && openSide[1] == 0 && x == x0) ||
					(openSide[0] == 1 && openSide[1] == 0 && x == x1) {
					wall = false
				}
			}
			if wall {
				m.SetTile(x, y, 10)
			} else {
				m.SetTile(x, y, 9)
			}
		}
	}
	m.SetTile(exitX, exitY, portalTile(0))
	altarX, altarY := exitX, exitY
	switch {
	case openSide[1] == -1:
		altarY = y1 - 1
	case openSide[1] == 1:
		altarY = y0 + 1
	case openSide[0] == -1:
		altarX = x1 - 1
	default:
		altarX = x0 + 1
	}
	m.SetTile(altarX, altarY, tileChapelAltar)
	m.SetTile(x0+1, y0+1, tileBellTower)
	for i := 0; i < 8; i++ {
		x, y := x0-2+g.Roll(w+4), y0-2+g.Roll(h+4)
		if m.Tile(x, y) == 2 {
			m.SetTile(x, y, tileGravestone)
		}
	}
}
