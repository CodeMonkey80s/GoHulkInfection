package game

import (
	"path/filepath"
	"testing"
)

func TestWoodsFOVHasUnlimitedRadiusButBlockersOcclude(t *testing.T) {
	blockers := []struct {
		name string
		tile uint8
	}{
		{"dense tree", tileDenseTree},
		{"bramble", tileBramble},
		{"bush", tileBush},
		{"stone wall", 10},
		{"closed wooden door", tileWoodenDoor},
	}
	for _, tc := range blockers {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGame(12000)
			g.SetupNewGame(0, "tester")
			g.Map = NewMap(mapWidth, woodsMapHeight)
			g.Theme = 0
			for i := range g.Map.Tiles {
				g.Map.Tiles[i] = 2
			}
			g.Player.X, g.Player.Y = 10, 11
			g.Map.SetTile(25, 11, tc.tile)
			g.UpdateFOV()
			if !g.Map.Visible[g.Map.idx(24, 11)] {
				t.Fatalf("tile before %s is not visible", tc.name)
			}
			if g.Map.Visible[g.Map.idx(26, 11)] {
				t.Fatalf("woods FOV sees straight through %s", tc.name)
			}
		})
	}

	// The radius is unlimited: a distant tile with clear line of sight is
	// visible even though the map is much wider than the normal FOV radius.
	g := NewGame(12001)
	g.SetupNewGame(0, "tester")
	g.Map = NewMap(mapWidth, woodsMapHeight)
	g.Theme = 0
	for i := range g.Map.Tiles {
		g.Map.Tiles[i] = 2
	}
	g.Player.X, g.Player.Y = 10, 11
	g.Map.SetTile(25, 11, tileDenseTree)
	g.UpdateFOV()
	if !g.Map.Visible[g.Map.idx(79, 2)] {
		t.Fatal("unlimited woods FOV was culled by a radius limit")
	}
	if g.Map.Visible[g.Map.idx(79, 11)] {
		t.Fatal("woods FOV sees through a dense tree")
	}
}

func TestWoodsRoadNetworkIsCardinallyReachable(t *testing.T) {
	for seed := int64(0); seed < 10; seed++ {
		g := NewGame(seed + 13000)
		g.SetupNewGame(0, "tester")
		for level := 0; level < numWoodsMaps; level++ {
			g.GenerateLevel(level)
			seen := make([]bool, g.Map.W*g.Map.H)
			start := g.Map.idx(g.Player.X, g.Player.Y)
			seen[start] = true
			q := []int{start}
			for head := 0; head < len(q); head++ {
				cur := q[head]
				cx, cy := cur%g.Map.W, cur/g.Map.W
				for _, d := range dirs4 {
					nx, ny := cx+d[0], cy+d[1]
					if !g.Map.InBounds(nx, ny) {
						continue
					}
					ni := g.Map.idx(nx, ny)
					if seen[ni] {
						continue
					}
					if g.Map.BlocksMove(nx, ny) && g.Map.Def(nx, ny).Type != 3 {
						continue
					}
					seen[ni] = true
					q = append(q, ni)
				}
			}
			px, py := findPortal(g)
			for y := 0; y < g.Map.H; y++ {
				for x := 0; x < g.Map.W; x++ {
					tile := g.Map.Tile(x, y)
					if tile == tileDirtRoad || tile == tileWoodenBridge || tile == tileStoneBridge || (x == px && y == py) {
						if !seen[g.Map.idx(x, y)] {
							t.Fatalf("seed %d level %d road tile %d,%d is not cardinally reachable", seed, level, x, y)
						}
					}
				}
			}
		}
	}
}

func TestWoodsLootBudgetStaysCompact(t *testing.T) {
	for level := 0; level < numWoodsMaps; level++ {
		for seed := int64(0); seed < 10; seed++ {
			g := NewGame(seed + int64(level)*14000)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(level)
			if len(g.Items) == 0 {
				t.Fatalf("level %d seed %d has no loot", level, seed)
			}
			if len(g.Items) > 12 {
				t.Fatalf("level %d seed %d has %d items, compact budget exceeded", level, seed, len(g.Items))
			}
		}
	}
}

func TestWoodsRunSaveRoundtrip(t *testing.T) {
	g := NewGame(15000)
	g.SetupNewGame(0, "tester")
	g.GenerateLevel(2)
	path := filepath.Join(t.TempDir(), "woods.json")
	if err := g.SaveGame(path); err != nil {
		t.Fatal(err)
	}
	g2 := NewGame(15001)
	if err := g2.LoadGame(path); err != nil {
		t.Fatal(err)
	}
	if g2.WoodsRun != g.WoodsRun {
		t.Fatalf("woods run changed across save/load: %+v vs %+v", g2.WoodsRun, g.WoodsRun)
	}
	if g2.Level != g.Level || g2.Map.Name != g.Map.Name {
		t.Fatalf("loaded woods map changed: level %d %q vs %d %q", g2.Level, g2.Map.Name, g.Level, g.Map.Name)
	}
}
