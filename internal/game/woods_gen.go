package game

import (
	"container/heap"
	"math"
)

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

// generateWoodsLevel replaces the room-section layout for levels 0-8 with a
// first-pass outdoor terrain generator: grass, forest density, clearings,
// a main dirt road and lantern posts.
func (g *Game) generateWoodsLevel(level int) {
	m := NewMap(mapWidth, mapHeight)
	m.Theme = 0
	m.Name = LevelThemes[0]
	for i := range m.Tiles {
		m.Tiles[i] = 2 // grass
	}
	for x := 0; x < m.W; x++ {
		m.SetTile(x, 0, tileBramble)
		m.SetTile(x, m.H-1, tileBramble)
	}
	for y := 0; y < m.H; y++ {
		m.SetTile(0, y, tileBramble)
		m.SetTile(m.W-1, y, tileBramble)
	}

	startX, startY := 2, 6+g.Roll(m.H-12)
	exitX, exitY := g.woodsExit(m)

	clearings := make([][3]int, 0, 6)
	clearings = append(clearings,
		[3]int{startX, startY, 4},
		[3]int{exitX, exitY, 4},
	)
	for i := 0; i < 3+g.Roll(3); i++ {
		clearings = append(clearings, [3]int{
			8 + g.Roll(m.W-16),
			4 + g.Roll(m.H-8),
			3 + g.Roll(4),
		})
	}

	g.growWoodsForest(m, level, clearings)
	for _, c := range clearings {
		clearWoodsArea(m, c[0], c[1], c[2], 2)
	}

	path := g.findWoodsRoad(m, startX, startY, exitX, exitY)
	for _, p := range path {
		m.SetTile(p[0], p[1], tileDirtRoad)
	}
	m.SetTile(startX, startY, tileDirtRoad)
	m.SetTile(exitX, exitY, portalTile(0))
	g.placeWoodsLanterns(m, path)
	g.scatterWoodsProps(m, startX, startY, exitX, exitY)

	if g.Player == nil {
		g.Player = &Critter{Proto: -1, Player: true, Align: AlignPlayer, Class: 4}
	}
	g.Player.X, g.Player.Y = startX, startY
	g.finishLevel(m, level, 0)
}

func (g *Game) woodsExit(m *Map) (int, int) {
	switch g.Roll(3) {
	case 0: // east
		return m.W - 3, 6 + g.Roll(m.H-12)
	case 1: // north
		return m.W/2 + g.Roll(m.W/2-8), 2
	default: // south
		return m.W/2 + g.Roll(m.W/2-8), m.H - 3
	}
}

func (g *Game) growWoodsForest(m *Map, level int, clearings [][3]int) {
	const cell = 8
	gw := (m.W+cell-1)/cell + 1
	gh := (m.H+cell-1)/cell + 1
	grid := make([]float64, gw*gh)
	for i := range grid {
		grid[i] = g.RNG.Float64()
	}
	density := 0.52 + float64(level)*0.012
	for y := 1; y < m.H-1; y++ {
		for x := 1; x < m.W-1; x++ {
			n := sampleWoodsNoise(grid, gw, gh, float64(x)/cell, float64(y)/cell)
			n = n*0.85 + g.RNG.Float64()*0.15
			for _, c := range clearings {
				d := math.Hypot(float64(x-c[0]), float64(y-c[1]))
				if d < float64(c[2]) {
					n -= (1 - d/float64(c[2])) * 0.35
				}
			}
			switch {
			case n > density+0.22:
				m.SetTile(x, y, 4) // big tree
			case n > density+0.08:
				m.SetTile(x, y, 3) // tree
			case n > density-0.02:
				if g.Roll(100) < 20 {
					m.SetTile(x, y, tileBramble)
				} else {
					m.SetTile(x, y, tileBush)
				}
			}
		}
	}
}

func sampleWoodsNoise(grid []float64, gw, gh int, fx, fy float64) float64 {
	x0, y0 := int(fx), int(fy)
	tx, ty := fx-float64(x0), fy-float64(y0)
	if x0 < 0 {
		x0, tx = 0, 0
	}
	if y0 < 0 {
		y0, ty = 0, 0
	}
	if x0 >= gw-1 {
		x0, tx = gw-2, 1
	}
	if y0 >= gh-1 {
		y0, ty = gh-2, 1
	}
	tx = tx * tx * (3 - 2*tx)
	ty = ty * ty * (3 - 2*ty)
	a := grid[y0*gw+x0]
	b := grid[y0*gw+x0+1]
	c := grid[(y0+1)*gw+x0]
	d := grid[(y0+1)*gw+x0+1]
	top := a + (b-a)*tx
	bottom := c + (d-c)*tx
	return top + (bottom-top)*ty
}

func clearWoodsArea(m *Map, cx, cy, radius, edge int) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			if x < edge || y < edge || x >= m.W-edge || y >= m.H-edge {
				continue
			}
			if abs(x-cx)+abs(y-cy) <= radius+edge {
				m.SetTile(x, y, 2) // grass
			}
		}
	}
}

func woodsRoadCost(t uint8) int {
	switch t {
	case tileBramble, tileFallenLog, tileDeadTree, tileLanternPost:
		return 18
	case 3, 4: // tree, big tree
		return 3
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
		for _, d := range dirs8 {
			nx, ny := cur.x+d[0], cur.y+d[1]
			if nx < 1 || ny < 1 || nx >= m.W-1 || ny >= m.H-1 {
				continue
			}
			ni := m.idx(nx, ny)
			step := woodsRoadCost(m.Tile(nx, ny))
			if d[0] != 0 && d[1] != 0 {
				step = step * 3 / 2
			}
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
	if len(path) < 12 {
		return
	}
	for i := 6; i < len(path)-6; i += 6 + g.Roll(3) {
		p := path[i]
		for _, d := range [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
			x, y := p[0]+d[0], p[1]+d[1]
			if x < 1 || y < 1 || x >= m.W-1 || y >= m.H-1 {
				continue
			}
			if m.Tile(x, y) == 2 || m.Tile(x, y) == tileBush {
				m.SetTile(x, y, tileLanternPost)
				break
			}
		}
	}
}

func (g *Game) scatterWoodsProps(m *Map, startX, startY, exitX, exitY int) {
	for i := 0; i < 90; i++ {
		x := 2 + g.Roll(m.W-4)
		y := 2 + g.Roll(m.H-4)
		if m.Tile(x, y) != 2 {
			continue
		}
		if abs(x-startX)+abs(y-startY) < 5 || abs(x-exitX)+abs(y-exitY) < 5 {
			continue
		}
		switch r := g.Roll(100); {
		case r < 35:
			m.SetTile(x, y, tileFlowers)
		case r < 65:
			m.SetTile(x, y, tileStump)
		case r < 82:
			m.SetTile(x, y, tileFallenLog)
		case r < 92:
			m.SetTile(x, y, tileDeadTree)
		}
	}
}
