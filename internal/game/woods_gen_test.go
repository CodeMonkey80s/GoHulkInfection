package game

import "testing"

func TestWoodsGeneratorTerrain(t *testing.T) {
	for level := 0; level < 9; level++ {
		for seed := int64(0); seed < 5; seed++ {
			g := NewGame(seed + int64(level)*1000)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(level)

			if g.Map.W != mapWidth || g.Map.H != mapHeight {
				t.Fatalf("level %d seed %d map = %dx%d", level, seed, g.Map.W, g.Map.H)
			}
			roads, lanterns, bushes, trees, brambles := 0, 0, 0, 0, 0
			for _, tile := range g.Map.Tiles {
				switch tile {
				case tileDirtRoad:
					roads++
				case tileLanternPost:
					lanterns++
				case tileBush:
					bushes++
				case 3, 4:
					trees++
				case tileBramble:
					brambles++
				}
			}
			if roads < 20 {
				t.Fatalf("level %d seed %d has only %d road tiles", level, seed, roads)
			}
			if lanterns == 0 {
				t.Fatalf("level %d seed %d has no lantern posts", level, seed)
			}
			if bushes == 0 || trees == 0 || brambles == 0 {
				t.Fatalf("level %d seed %d missing vegetation: bushes=%d trees=%d brambles=%d", level, seed, bushes, trees, brambles)
			}
			if !g.Map.IsPassable(g.Player.X, g.Player.Y) {
				t.Fatalf("level %d seed %d start is not passable", level, seed)
			}
			px, py := findPortal(g)
			if px < 0 {
				t.Fatalf("level %d seed %d has no portal", level, seed)
			}
			if !reachable(g, g.Player.X, g.Player.Y, px, py) {
				t.Fatalf("level %d seed %d portal is not reachable", level, seed)
			}
		}
	}
}

func TestWoodsGeneratorIsDeterministic(t *testing.T) {
	build := func() *Game {
		g := NewGame(4242)
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(4)
		return g
	}
	a, b := build(), build()
	if len(a.Map.Tiles) != len(b.Map.Tiles) {
		t.Fatal("map sizes differ")
	}
	for i := range a.Map.Tiles {
		if a.Map.Tiles[i] != b.Map.Tiles[i] {
			t.Fatalf("tile %d differs: %d vs %d", i, a.Map.Tiles[i], b.Map.Tiles[i])
		}
	}
	if a.Player.X != b.Player.X || a.Player.Y != b.Player.Y {
		t.Fatal("starting positions differ")
	}
}
