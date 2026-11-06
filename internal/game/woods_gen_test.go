package game

import (
	"strings"
	"testing"
)

func TestWoodsGeneratorTerrain(t *testing.T) {
	for level := 0; level < 9; level++ {
		for seed := int64(0); seed < 10; seed++ {
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
			if px == g.Player.X && py == g.Player.Y {
				t.Fatalf("level %d seed %d spawned the player on the exit", level, seed)
			}
		}
	}
}

func TestWoodsGeneratorBuildsHuts(t *testing.T) {
	for level := 0; level < 9; level++ {
		for seed := int64(0); seed < 10; seed++ {
			g := NewGame(seed + int64(level)*2000)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(level)
			doors, windows, furniture := 0, 0, 0
			for y := 0; y < g.Map.H; y++ {
				for x := 0; x < g.Map.W; x++ {
					switch g.Map.Tile(x, y) {
					case tileWoodenDoor:
						doors++
						if !reachable(g, g.Player.X, g.Player.Y, x, y) {
							t.Fatalf("level %d seed %d hut door at %d,%d is unreachable", level, seed, x, y)
						}
						roadside := false
						for _, d := range [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
							t := g.Map.Tile(x+d[0], y+d[1])
							if t == tileDirtRoad || t == tileWoodenBridge {
								roadside = true
								break
							}
						}
						if !roadside {
							t.Fatalf("level %d seed %d hut door at %d,%d is not connected to a road", level, seed, x, y)
						}
					case tileWoodenWindow:
						windows++
					case tileTable, tileChair, tileBed, tileCrate, tileBarrel, tileFireplace, tileShelf:
						furniture++
					}
				}
			}
			spec := woodLevelSpecFor(level)
			maxDoors := spec.HutMax
			if level == 0 || level == 3 {
				maxDoors++ // landmark shed/mill house also has a wooden door
			}
			if doors < spec.HutMin || doors > maxDoors {
				t.Fatalf("level %d seed %d has %d hut doors, want %d-%d", level, seed, doors, spec.HutMin, maxDoors)
			}
			if windows == 0 {
				t.Fatalf("level %d seed %d has no hut windows", level, seed)
			}
			if furniture == 0 {
				t.Fatalf("level %d seed %d has no hut furniture", level, seed)
			}
		}
	}
}

func TestWoodsLandmarks(t *testing.T) {
	expected := [][]uint8{
		{tileWoodpile},
		{tileOrchardTree, tileBeehive, tileCiderPress},
		{tileGravestone, tileIronFence, tileCryptDoor},
		{tileMillStone},
		{tileTent, tileCampfire, tileCage, tileDryingRack},
		{tileMushroomRing, tileDruidAltar, 12},
		{tileWell, tileMarketStall, tileCart, tileHaystack},
		{tileKiln},
		{tileChapelAltar, tileBellTower, tileGravestone},
	}
	for level := 0; level < 9; level++ {
		for seed := int64(0); seed < 5; seed++ {
			g := NewGame(seed + int64(level)*5000)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(level)
			spec := woodLevelSpecFor(level)
			if g.Map.Name != spec.Name {
				t.Fatalf("level %d name = %q, want %q", level, g.Map.Name, spec.Name)
			}
			counts := map[uint8]int{}
			for _, tile := range g.Map.Tiles {
				counts[tile]++
			}
			for _, tile := range expected[level] {
				if counts[tile] == 0 {
					t.Fatalf("level %d seed %d missing landmark tile %d", level, seed, tile)
				}
			}
		}
	}
}

func TestWoodsExitModes(t *testing.T) {
	seen := map[int]int{}
	for level := 0; level < 8; level++ {
		for seed := int64(0); seed < 30; seed++ {
			g := NewGame(seed + int64(level)*8000)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(level)
			seen[g.Map.ExitMode]++
		}
	}
	for _, mode := range []int{ExitEdge, ExitCenter, ExitLandmark} {
		if seen[mode] == 0 {
			t.Fatalf("exit mode %d never appeared: %v", mode, seen)
		}
	}
	for seed := int64(0); seed < 20; seed++ {
		g := NewGame(seed + 90000)
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(8)
		if g.Map.ExitMode != ExitLandmark {
			t.Fatalf("level 8 exit mode = %d, want landmark", g.Map.ExitMode)
		}
	}
}

func TestAllWoodsLevelsRender(t *testing.T) {
	for level := 0; level < 9; level++ {
		g := NewGame(int64(91000 + level))
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(level)
		g.State = StatePlaying
		g.UpdateFOV()
		u, s := newTestUI(g, ColorModeTrueColor)
		u.drawGame()
		s.Show()
		text := screenText(s)
		if !strings.Contains(text, "@") {
			t.Fatalf("level %d did not render the player", level)
		}
		s.Fini()
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
