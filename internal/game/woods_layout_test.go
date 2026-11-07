package game

import (
	"path/filepath"
	"testing"
)

func TestWoodsFOVCoversWholeMap(t *testing.T) {
	for level := 0; level < numWoodsMaps; level++ {
		g := NewGame(int64(12000 + level))
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(level)
		g.UpdateFOV()
		for i := range g.Map.Visible {
			if !g.Map.Visible[i] || !g.Map.Seen[i] {
				t.Fatalf("level %d tile %d is not visible/seen under unlimited woods FOV", level, i)
			}
		}
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
