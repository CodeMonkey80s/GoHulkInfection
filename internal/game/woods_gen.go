package game

import "container/heap"

type roadNode struct {
	x, y int
	cost int
}

type roadHeap []roadNode

func (h roadHeap) Len() int            { return len(h) }
func (h roadHeap) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h roadHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *roadHeap) Push(x interface{}) { *h = append(*h, x.(roadNode)) }
func (h *roadHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

// generateWoodsLevel builds one compact 80x23 map from the run's theme order
// and snake direction sequence.
func (g *Game) generateWoodsLevel(level int) {
	index := woodsRunIndex(level)
	if !validWoodsRun(&g.WoodsRun) {
		g.WoodsRun = newWoodsRun(g.RNG)
	}
	theme := g.WoodsRun.Themes[index]
	entryEdge := g.WoodsRun.Entries[index]
	exitEdge := g.WoodsRun.Exits[index]

	m := NewMap(mapWidth, woodsMapHeight)
	m.Theme = 0
	m.Name = woodsThemeName(theme)
	for i := range m.Tiles {
		m.Tiles[i] = 2 // grass
	}
	g.buildWoodsTreeBorder(m)

	entryPos := g.woodsEdgePosition(m, entryEdge)
	exitPos := g.woodsEdgePosition(m, exitEdge)
	_, entryInside := woodsEdgePoints(m, entryEdge, entryPos)
	_, exitInside := woodsEdgePoints(m, exitEdge, exitPos)

	cx, cy := g.woodsLandmarkCenter(m, theme)
	anchor := g.buildWoodsThemeLandmark(m, theme, cx, cy, entryEdge)
	path := g.buildWoodsRoute(m, entryEdge, entryPos, exitEdge, exitPos, anchor)
	g.addWoodsThemeWater(m, theme, path, cx, cy)
	g.rebridgeWoodsPath(m, path)

	spec := woodThemeSpecFor(theme)
	g.placeWoodsHuts(m, path, entryInside[0], entryInside[1], exitInside[0], exitInside[1], spec.HutMin, spec.HutMax)
	g.placeWoodsLanterns(m, path)
	g.scatterWoodsProps(m, theme, path, entryInside, exitInside)

	if g.Player == nil {
		g.Player = &Critter{Proto: -1, Player: true, Align: AlignPlayer, Class: 4}
	}
	g.Player.X, g.Player.Y = m.FirstFreeNear(entryInside[0], entryInside[1])
	g.finishLevel(m, level, 0)
}

func (g *Game) buildWoodsTreeBorder(m *Map) {
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			if x >= 2 && y >= 2 && x < m.W-2 && y < m.H-2 {
				continue
			}
			switch {
			case g.Roll(100) < 18:
				m.SetTile(x, y, tileBramble)
			case g.Roll(100) < 22:
				m.SetTile(x, y, tileDeadTree)
			default:
				m.SetTile(x, y, tileDenseTree)
			}
		}
	}
}

func (g *Game) woodsEdgePosition(m *Map, edge int) int {
	switch edge {
	case EdgeNorth, EdgeSouth:
		return 10 + g.Roll(m.W-20)
	default:
		return 4 + g.Roll(m.H-8)
	}
}

// woodsEdgePoints returns the tile on the map edge and the first safe tile
// inside the two-tile tree border.
func woodsEdgePoints(m *Map, edge, pos int) (edgePoint, insidePoint [2]int) {
	switch edge {
	case EdgeNorth:
		return [2]int{pos, 0}, [2]int{pos, 3}
	case EdgeSouth:
		return [2]int{pos, m.H - 1}, [2]int{pos, m.H - 4}
	case EdgeEast:
		return [2]int{m.W - 1, pos}, [2]int{m.W - 4, pos}
	default:
		return [2]int{0, pos}, [2]int{3, pos}
	}
}

func (g *Game) buildWoodsRoute(m *Map, entryEdge, entryPos, exitEdge, exitPos int, anchor [2]int) [][2]int {
	entryEdgePoint, entryInside := woodsEdgePoints(m, entryEdge, entryPos)
	exitEdgePoint, exitInside := woodsEdgePoints(m, exitEdge, exitPos)
	if m.BlocksMove(anchor[0], anchor[1]) {
		m.SetTile(anchor[0], anchor[1], 2)
	}

	first := g.findWoodsRoad(m, entryInside[0], entryInside[1], anchor[0], anchor[1])
	second := g.findWoodsRoad(m, anchor[0], anchor[1], exitInside[0], exitInside[1])
	var path [][2]int
	path = append(path, first...)
	if len(second) > 1 {
		path = append(path, second[1:]...)
	}
	if len(first) == 0 || len(second) == 0 {
		path = g.directWoodsRoute(entryInside, anchor, exitInside)
	}
	for _, p := range path {
		g.setWoodsRoadTile(m, p[0], p[1])
	}
	path = append(path, g.carveWoodsEdgeRoad(m, entryEdgePoint, entryInside)...)
	path = append(path, g.carveWoodsEdgeRoad(m, exitEdgePoint, exitInside)...)
	m.SetTile(exitEdgePoint[0], exitEdgePoint[1], portalTile(0))
	return path
}

func (g *Game) carveWoodsEdgeRoad(m *Map, edgePoint, insidePoint [2]int) [][2]int {
	line := m.Line(edgePoint[0], edgePoint[1], insidePoint[0], insidePoint[1])
	path := make([][2]int, 0, len(line)+2)
	path = append(path, edgePoint)
	g.setWoodsRoadTile(m, edgePoint[0], edgePoint[1])
	for _, p := range line {
		path = append(path, p)
		g.setWoodsRoadTile(m, p[0], p[1])
	}
	path = append(path, insidePoint)
	g.setWoodsRoadTile(m, insidePoint[0], insidePoint[1])
	return path
}

func (g *Game) directWoodsRoute(points ...[2]int) [][2]int {
	var out [][2]int
	for i := 0; i < len(points); i++ {
		if i == 0 {
			out = append(out, points[i])
			continue
		}
		from := points[i-1]
		to := points[i]
		x, y := from[0], from[1]
		for x != to[0] {
			if x < to[0] {
				x++
			} else {
				x--
			}
			out = append(out, [2]int{x, y})
		}
		for y != to[1] {
			if y < to[1] {
				y++
			} else {
				y--
			}
			out = append(out, [2]int{x, y})
		}
	}
	return out
}

func (g *Game) rebridgeWoodsPath(m *Map, path [][2]int) {
	for _, p := range path {
		if isWoodsWater(m.Tile(p[0], p[1])) {
			m.SetTile(p[0], p[1], tileWoodenBridge)
		}
	}
}

func woodsRoadCost(t uint8) int {
	if isWoodsStructure(t) {
		return 80
	}
	switch t {
	case tileDeepWater:
		return 40
	case tileShallowWater, tileLilyPad:
		return 16
	case tileRiverBank:
		return 2
	case tileBramble, tileFallenLog, tileDeadTree, tileLanternPost, tileDenseTree:
		return 18
	case 3, 4: // tree, big tree
		return 4
	case tileBush, tileReeds:
		return 2
	default:
		return 1
	}
}

func (g *Game) findWoodsRoad(m *Map, sx, sy, tx, ty int) [][2]int {
	const inf = 1 << 30
	n := m.W * m.H
	dist := make([]int, n)
	prev := make([]int, n)
	for i := range dist {
		dist[i] = inf
		prev[i] = -1
	}
	start := m.idx(sx, sy)
	dist[start] = 0
	h := &roadHeap{{x: sx, y: sy, cost: 0}}
	heap.Init(h)
	for h.Len() > 0 {
		cur := heap.Pop(h).(roadNode)
		ci := m.idx(cur.x, cur.y)
		if cur.cost != dist[ci] {
			continue
		}
		if cur.x == tx && cur.y == ty {
			break
		}
		for _, d := range dirs4 {
			nx, ny := cur.x+d[0], cur.y+d[1]
			if nx < 2 || ny < 2 || nx >= m.W-2 || ny >= m.H-2 {
				continue
			}
			ni := m.idx(nx, ny)
			step := woodsRoadCost(m.Tile(nx, ny))
			nd := cur.cost + step
			if nd < dist[ni] {
				dist[ni] = nd
				prev[ni] = ci
				heap.Push(h, roadNode{x: nx, y: ny, cost: nd})
			}
		}
	}
	if prev[m.idx(tx, ty)] == -1 && (sx != tx || sy != ty) {
		return nil
	}
	var path [][2]int
	for at := m.idx(tx, ty); at != -1; at = prev[at] {
		path = append(path, [2]int{at % m.W, at / m.W})
		if at == start {
			break
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func (g *Game) placeWoodsLanterns(m *Map, path [][2]int) {
	if len(path) < 6 {
		return
	}
	placed := 0
	tryAt := func(p [2]int) bool {
		for radius := 1; radius <= 3; radius++ {
			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					if abs(dx) != radius && abs(dy) != radius {
						continue
					}
					x, y := p[0]+dx, p[1]+dy
					if x < 2 || y < 2 || x >= m.W-2 || y >= m.H-2 {
						continue
					}
					t := m.Tile(x, y)
					if t == 2 || t == tileBush || t == tileFlowers || t == tileMud {
						m.SetTile(x, y, tileLanternPost)
						return true
					}
				}
			}
		}
		return false
	}
	for i := 6; i < len(path)-2 && placed < 8; i += 5 + g.Roll(3) {
		if tryAt(path[i]) {
			placed++
		}
	}
	for i := 2; i < len(path)-1 && placed < 4; i++ {
		if tryAt(path[i]) {
			placed++
		}
	}
}

func (g *Game) scatterWoodsProps(m *Map, theme int, path [][2]int, entry, exit [2]int) {
	want := 6 + g.Roll(7)
	placed := 0
	for tries := 0; tries < 300 && placed < want; tries++ {
		x := 2 + g.Roll(m.W-4)
		y := 2 + g.Roll(m.H-4)
		if m.Tile(x, y) != 2 {
			continue
		}
		if abs(x-entry[0])+abs(y-entry[1]) < 4 || abs(x-exit[0])+abs(y-exit[1]) < 4 {
			continue
		}
		if abs(x-m.LandmarkX)+abs(y-m.LandmarkY) < 3 {
			continue
		}
		var tile uint8
		switch theme {
		case WoodsCemetery:
			switch r := g.Roll(100); {
			case r < 35:
				tile = tileGravestone
			case r < 65:
				tile = tileDeadTree
			case r < 85:
				tile = tileStump
			default:
				tile = tileBush
			}
		case WoodsMeteorField:
			switch r := g.Roll(100); {
			case r < 40:
				tile = tileBurnedGround
			case r < 65:
				tile = tileDeadTree
			case r < 85:
				tile = tileStump
			default:
				tile = tileMeteorRock
			}
		case WoodsMillPond:
			switch r := g.Roll(100); {
			case r < 30:
				tile = tileFlowers
			case r < 55:
				tile = tileStump
			case r < 75:
				tile = tileBush
			default:
				tile = tileReeds
			}
		default:
			switch r := g.Roll(100); {
			case r < 35:
				tile = tileFlowers
			case r < 60:
				tile = tileStump
			case r < 80:
				tile = tileBush
			default:
				tile = tileFallenLog
			}
		}
		m.SetTile(x, y, tile)
		placed++
	}
}

// addWoodsThemeWater paints the optional themed pond or river. The Mill's pond
// deliberately overlaps the road so the route crosses it on a bridge.
func (g *Game) addWoodsThemeWater(m *Map, theme int, path [][2]int, cx, cy int) {
	switch theme {
	case WoodsMillPond:
		order := make([]int, 0, len(path))
		for i := len(path) / 3; i < len(path)*2/3; i++ {
			order = append(order, i)
		}
		for i := 2; i < len(path)-2; i++ {
			order = append(order, i)
		}
		for _, idx := range order {
			if idx <= 1 || idx >= len(path)-2 {
				continue
			}
			p := path[idx]
			if abs(p[0]-cx)+abs(p[1]-cy) < 8 {
				continue
			}
			if g.paintWoodsPondIf(m, p[0], p[1], 4+g.Roll(2), 3+g.Roll(2), true) {
				return
			}
		}
		for _, idx := range order {
			if idx <= 1 || idx >= len(path)-2 {
				continue
			}
			p := path[idx]
			if g.paintWoodsPondIf(m, p[0], p[1], 4, 3, true) {
				return
			}
		}
	}
}

func (g *Game) paintWoodsPondIf(m *Map, cx, cy, rx, ry int, allowRoad bool) bool {
	if cx < 3 || cy < 3 || cx >= m.W-3 || cy >= m.H-3 {
		return false
	}
	x0, x1 := maxInt(2, cx-rx-2), minInt(m.W-3, cx+rx+2)
	y0, y1 := maxInt(2, cy-ry-2), minInt(m.H-3, cy+ry+2)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			t := m.Tile(x, y)
			if isWoodsStructure(t) || t == portalTile(0) {
				return false
			}
			if !allowRoad && (t == tileDirtRoad || t == tileWoodenBridge) {
				return false
			}
		}
	}
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
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
	return true
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
