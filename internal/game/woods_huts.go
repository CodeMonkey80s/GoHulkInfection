package game

type woodsHut struct {
	x0, y0 int
	x1, y1 int
	doorX  int
	doorY  int
}

func (g *Game) placeWoodsHuts(m *Map, path [][2]int, startX, startY, exitX, exitY, minHuts, maxHuts int) []woodsHut {
	if maxHuts < minHuts {
		maxHuts = minHuts
	}
	want := minHuts + g.Roll(maxHuts-minHuts+1)
	huts := make([]woodsHut, 0, want)
	if len(path) < 12 {
		return huts
	}
	tryAt := func(i int) {
		p := path[i]
		sides := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		g.RNG.Shuffle(len(sides), func(a, b int) {
			sides[a], sides[b] = sides[b], sides[a]
		})
		for _, side := range sides {
			if len(huts) >= want {
				return
			}
			for offset := 1; offset <= 20; offset++ {
				if h, ok := g.tryPlaceHut(m, p, side, offset, startX, startY, exitX, exitY, huts); ok {
					huts = append(huts, h)
					return
				}
			}
		}
	}
	for i := 8; i < len(path)-8 && len(huts) < want; i += 4 + g.Roll(4) {
		tryAt(i)
	}
	for i := 4; i < len(path)-4 && len(huts) < want; i += 3 {
		tryAt(i)
	}
	for i := 2; i < len(path)-2 && len(huts) < want; i++ {
		tryAt(i)
	}
	return huts
}

func (g *Game) tryPlaceHut(m *Map, p, side [2]int, offset int, startX, startY, exitX, exitY int, huts []woodsHut) (woodsHut, bool) {
	w := 5 + g.Roll(3)
	h := 5 + g.Roll(2)
	var x0, y0, x1, y1, doorX, doorY int
	doorX, doorY = p[0]+side[0]*offset, p[1]+side[1]*offset
	switch {
	case side[1] == -1: // north of road
		x0 = doorX - w/2
		x1 = x0 + w - 1
		y1 = doorY
		y0 = y1 - h + 1
	case side[1] == 1: // south of road
		x0 = doorX - w/2
		x1 = x0 + w - 1
		y0 = doorY
		y1 = y0 + h - 1
	case side[0] == -1: // west of road
		y0 = doorY - h/2
		y1 = y0 + h - 1
		x1 = doorX
		x0 = x1 - w + 1
	default: // east of road
		y0 = doorY - h/2
		y1 = y0 + h - 1
		x0 = doorX
		x1 = x0 + w - 1
	}
	if x0 < 2 || y0 < 2 || x1 >= m.W-2 || y1 >= m.H-2 {
		return woodsHut{}, false
	}
	if abs(doorX-startX)+abs(doorY-startY) < 5 || abs(doorX-exitX)+abs(doorY-exitY) < 5 {
		return woodsHut{}, false
	}
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			t := m.Tile(x, y)
			if t == tileDirtRoad || t == tileWoodenBridge || t == portalTile(0) || t == tileRiverBank || isWoodsWater(t) || isWoodsStructure(t) {
				return woodsHut{}, false
			}
		}
	}
	for k := 1; k < offset; k++ {
		if isWoodsStructure(m.Tile(p[0]+side[0]*k, p[1]+side[1]*k)) {
			return woodsHut{}, false
		}
	}
	for _, other := range huts {
		if x0-2 <= other.x1 && x1+2 >= other.x0 && y0-2 <= other.y1 && y1+2 >= other.y0 {
			return woodsHut{}, false
		}
	}

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x == x0 || x == x1 || y == y0 || y == y1 {
				m.SetTile(x, y, tileWoodenWall)
			} else {
				m.SetTile(x, y, tileWoodenFloor)
			}
		}
	}
	m.SetTile(doorX, doorY, tileWoodenDoor)
	for k := 1; k < offset; k++ {
		g.setWoodsRoadTile(m, p[0]+side[0]*k, p[1]+side[1]*k)
	}
	hut := woodsHut{x0: x0, y0: y0, x1: x1, y1: y1, doorX: doorX, doorY: doorY}
	g.placeHutWindows(m, hut)
	g.furnishHut(m, hut)
	return hut, true
}

func (g *Game) placeHutWindows(m *Map, h woodsHut) {
	var spots [][2]int
	for x := h.x0 + 1; x < h.x1; x++ {
		if !(x == h.doorX && h.y0 == h.doorY) {
			spots = append(spots, [2]int{x, h.y0})
		}
		if !(x == h.doorX && h.y1 == h.doorY) {
			spots = append(spots, [2]int{x, h.y1})
		}
	}
	for y := h.y0 + 1; y < h.y1; y++ {
		if !(y == h.doorY && h.x0 == h.doorX) {
			spots = append(spots, [2]int{h.x0, y})
		}
		if !(y == h.doorY && h.x1 == h.doorX) {
			spots = append(spots, [2]int{h.x1, y})
		}
	}
	g.RNG.Shuffle(len(spots), func(a, b int) {
		spots[a], spots[b] = spots[b], spots[a]
	})
	count := 1 + g.Roll(2)
	for i := 0; i < count && i < len(spots); i++ {
		m.SetTile(spots[i][0], spots[i][1], tileWoodenWindow)
	}
}

func (g *Game) furnishHut(m *Map, h woodsHut) {
	entryX, entryY := h.doorX, h.doorY
	switch {
	case h.doorY == h.y0:
		entryY++
	case h.doorY == h.y1:
		entryY--
	case h.doorX == h.x0:
		entryX++
	default:
		entryX--
	}
	var spots [][2]int
	for y := h.y0 + 1; y < h.y1; y++ {
		for x := h.x0 + 1; x < h.x1; x++ {
			if x == entryX && y == entryY {
				continue
			}
			spots = append(spots, [2]int{x, y})
		}
	}
	g.RNG.Shuffle(len(spots), func(a, b int) {
		spots[a], spots[b] = spots[b], spots[a]
	})
	furniture := []uint8{tileTable, tileChair, tileBed, tileCrate, tileBarrel, tileFireplace, tileShelf}
	count := 2 + g.Roll(3)
	for i := 0; i < count && i < len(spots); i++ {
		m.SetTile(spots[i][0], spots[i][1], furniture[g.Roll(len(furniture))])
	}
}
