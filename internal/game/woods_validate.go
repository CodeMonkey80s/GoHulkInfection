package game

import "fmt"

// woodsMapProblems returns the invariants a generated compact woods map must
// satisfy. It is used by generator tests and by the exhaustive theme/edge matrix.
func woodsMapProblems(m *Map, theme, entryEdge, exitEdge, px, py int) []string {
	var problems []string
	if m == nil {
		return []string{"map is nil"}
	}
	if m.W != mapWidth || m.H != woodsMapHeight {
		problems = append(problems, fmt.Sprintf("map size is %dx%d, want %dx%d", m.W, m.H, mapWidth, woodsMapHeight))
	}
	if entryEdge == exitEdge {
		problems = append(problems, "entry and exit use the same edge")
	}

	voids := 0
	counts := map[uint8]int{}
	for _, tile := range m.Tiles {
		if tile == 0 || tile == 1 {
			voids++
		}
		counts[tile]++
	}
	if voids > 0 {
		problems = append(problems, fmt.Sprintf("%d void tiles", voids))
	}

	portalX, portalY, portals := -1, -1, 0
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			if m.Tile(x, y) == portalTile(0) {
				portalX, portalY, portals = x, y, portals+1
			}
			t := m.Tile(x, y)
			if (t == tileWoodenBridge || t == tileStoneBridge) && m.BlocksMove(x, y) {
				problems = append(problems, fmt.Sprintf("bridge at %d,%d blocks movement", x, y))
			}
		}
	}
	if portals != 1 {
		problems = append(problems, fmt.Sprintf("found %d portals, want 1", portals))
	} else if edgeAt(m, portalX, portalY) != exitEdge {
		problems = append(problems, fmt.Sprintf("portal at %d,%d is on edge %d, want %d", portalX, portalY, edgeAt(m, portalX, portalY), exitEdge))
	}

	for _, edge := range []int{EdgeNorth, EdgeSouth, EdgeEast, EdgeWest} {
		gaps := edgeGapTiles(m, edge)
		want := 0
		if edge == entryEdge || edge == exitEdge {
			want = 1
		}
		if len(gaps) != want {
			problems = append(problems, fmt.Sprintf("edge %d has %d passable edge tiles, want %d", edge, len(gaps), want))
			continue
		}
		if edge == entryEdge && len(gaps) == 1 {
			if d := manhattan(px, py, gaps[0][0], gaps[0][1]); d > 4 {
				problems = append(problems, fmt.Sprintf("player start %d,%d is %d tiles from the entry edge", px, py, d))
			}
		}
		if edge == exitEdge && len(gaps) == 1 && (gaps[0][0] != portalX || gaps[0][1] != portalY) {
			problems = append(problems, "exit edge gap is not the portal")
		}
	}

	if !m.IsPassable(px, py) {
		problems = append(problems, fmt.Sprintf("player start %d,%d is not passable", px, py))
	}
	if portalX < 0 || (px == portalX && py == portalY) {
		problems = append(problems, "player starts on the exit")
	}

	reachable := map[int]bool{}
	for _, index := range m.Reachable(px, py) {
		reachable[index] = true
	}
	roads := 0
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			t := m.Tile(x, y)
			if t == tileDirtRoad || t == tileWoodenBridge || t == tileStoneBridge || t == portalTile(0) {
				roads++
				if !reachable[m.idx(x, y)] {
					problems = append(problems, fmt.Sprintf("road tile %d,%d is not reachable", x, y))
				}
			}
		}
	}
	if roads < 20 {
		problems = append(problems, fmt.Sprintf("only %d road tiles", roads))
	}

	if !m.InBounds(m.LandmarkX, m.LandmarkY) {
		problems = append(problems, "landmark is outside the map")
	} else {
		identity := themeIdentityTile(theme)
		found := false
		roadNearLandmark := false
		for y := 0; y < m.H; y++ {
			for x := 0; x < m.W; x++ {
				if manhattan(x, y, m.LandmarkX, m.LandmarkY) > 10 {
					continue
				}
				if m.Tile(x, y) == identity {
					found = true
				}
				t := m.Tile(x, y)
				if (t == tileDirtRoad || t == tileWoodenBridge || t == tileStoneBridge) && reachable[m.idx(x, y)] {
					roadNearLandmark = true
				}
			}
		}
		if !found {
			problems = append(problems, fmt.Sprintf("no theme landmark tile %d near the landmark", identity))
		}
		if !roadNearLandmark {
			problems = append(problems, "landmark has no reachable road nearby")
		}
	}

	switch theme {
	case WoodsChurch:
		if counts[tileBellTower] != 1 || counts[tileChapelAltar] != 1 {
			problems = append(problems, fmt.Sprintf("church has %d bell towers and %d altars", counts[tileBellTower], counts[tileChapelAltar]))
		}
	case WoodsPeasantHouses:
		if counts[tileWell] != 1 {
			problems = append(problems, fmt.Sprintf("peasant houses has %d wells, want 1", counts[tileWell]))
		}
	case WoodsCemetery:
		if counts[tileCryptDoor] != 1 {
			problems = append(problems, fmt.Sprintf("cemetery has %d crypt doors, want 1", counts[tileCryptDoor]))
		}
	case WoodsMeteorField:
		if counts[tileMeteorRock] < 3 {
			problems = append(problems, fmt.Sprintf("meteor field has %d meteor rocks", counts[tileMeteorRock]))
		}
	}
	return problems
}

func themeIdentityTile(theme int) uint8 {
	switch theme {
	case WoodsChurch:
		return tileBellTower
	case WoodsPeasantHouses:
		return tileWell
	case WoodsMillPond:
		return tileMillStone
	case WoodsCemetery:
		return tileCryptDoor
	default:
		return tileMeteorRock
	}
}

func edgeAt(m *Map, x, y int) int {
	switch {
	case y == 0:
		return EdgeNorth
	case y == m.H-1:
		return EdgeSouth
	case x == m.W-1:
		return EdgeEast
	default:
		return EdgeWest
	}
}

func edgeGapTiles(m *Map, edge int) [][2]int {
	var out [][2]int
	switch edge {
	case EdgeNorth:
		for x := 0; x < m.W; x++ {
			if m.IsPassable(x, 0) {
				out = append(out, [2]int{x, 0})
			}
		}
	case EdgeSouth:
		for x := 0; x < m.W; x++ {
			if m.IsPassable(x, m.H-1) {
				out = append(out, [2]int{x, m.H - 1})
			}
		}
	case EdgeEast:
		for y := 0; y < m.H; y++ {
			if m.IsPassable(m.W-1, y) {
				out = append(out, [2]int{m.W - 1, y})
			}
		}
	default:
		for y := 0; y < m.H; y++ {
			if m.IsPassable(0, y) {
				out = append(out, [2]int{0, y})
			}
		}
	}
	return out
}
