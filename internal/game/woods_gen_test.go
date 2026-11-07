package game

import (
	"strings"
	"testing"
)

func TestWoodsMapsAreCompactThemedAndConnected(t *testing.T) {
	for seed := int64(0); seed < 12; seed++ {
		g := NewGame(seed)
		g.SetupNewGame(0, "tester")
		for level := 0; level < numWoodsMaps; level++ {
			g.GenerateLevel(level)
			if g.Map.W != mapWidth || g.Map.H != woodsMapHeight {
				t.Fatalf("seed %d level %d map = %dx%d, want %dx%d", seed, level, g.Map.W, g.Map.H, mapWidth, woodsMapHeight)
			}
			theme := g.WoodsRun.Themes[level]
			if g.Map.Name != woodsThemeName(theme) {
				t.Fatalf("seed %d level %d name = %q, want %q", seed, level, g.Map.Name, woodsThemeName(theme))
			}
			voids := 0
			roads, lanterns, trees := 0, 0, 0
			for _, tile := range g.Map.Tiles {
				switch tile {
				case 0, 1:
					voids++
				case tileDirtRoad, tileWoodenBridge, tileStoneBridge:
					roads++
				case tileLanternPost:
					lanterns++
				case 3, 4, tileDenseTree:
					trees++
				}
			}
			if voids > 0 {
				t.Fatalf("seed %d level %d has %d void tiles", seed, level, voids)
			}
			if roads < 20 {
				t.Fatalf("seed %d level %d has only %d road tiles", seed, level, roads)
			}
			if lanterns == 0 {
				t.Fatalf("seed %d level %d has no lantern posts", seed, level)
			}
			if trees == 0 {
				t.Fatalf("seed %d level %d has no border trees", seed, level)
			}
			if !g.Map.IsPassable(g.Player.X, g.Player.Y) {
				t.Fatalf("seed %d level %d start is not passable", seed, level)
			}
			px, py := findPortal(g)
			if px < 0 {
				t.Fatalf("seed %d level %d has no portal", seed, level)
			}
			if portalEdge(g, px, py) != g.WoodsRun.Exits[level] {
				t.Fatalf("seed %d level %d portal at %d,%d is not on exit edge %d", seed, level, px, py, g.WoodsRun.Exits[level])
			}
			if !reachable(g, g.Player.X, g.Player.Y, px, py) {
				t.Fatalf("seed %d level %d portal is not reachable", seed, level)
			}
			if px == g.Player.X && py == g.Player.Y {
				t.Fatalf("seed %d level %d spawned the player on the exit", seed, level)
			}
			checkWoodsBorderGaps(t, g, level)
		}
	}
}

func portalEdge(g *Game, x, y int) int {
	switch {
	case y == 0:
		return EdgeNorth
	case y == g.Map.H-1:
		return EdgeSouth
	case x == g.Map.W-1:
		return EdgeEast
	default:
		return EdgeWest
	}
}

func checkWoodsBorderGaps(t *testing.T, g *Game, level int) {
	t.Helper()
	cells := map[int]int{}
	for x := 0; x < g.Map.W; x++ {
		if g.Map.IsPassable(x, 0) {
			cells[EdgeNorth]++
		}
		if g.Map.IsPassable(x, g.Map.H-1) {
			cells[EdgeSouth]++
		}
	}
	for y := 0; y < g.Map.H; y++ {
		if g.Map.IsPassable(0, y) {
			cells[EdgeWest]++
		}
		if g.Map.IsPassable(g.Map.W-1, y) {
			cells[EdgeEast]++
		}
	}
	entry, exit := g.WoodsRun.Entries[level], g.WoodsRun.Exits[level]
	for _, edge := range []int{EdgeNorth, EdgeSouth, EdgeEast, EdgeWest} {
		want := 0
		if edge == entry || edge == exit {
			want = 1
		}
		if cells[edge] != want {
			t.Fatalf("level %d edge %d has %d passable border tiles, want %d", level, edge, cells[edge], want)
		}
	}
}

func TestWoodsGeneratorBuildsHuts(t *testing.T) {
	for seed := int64(0); seed < 8; seed++ {
		g := NewGame(seed + 2000)
		g.SetupNewGame(0, "tester")
		for level := 0; level < numWoodsMaps; level++ {
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
							if t == tileDirtRoad || t == tileWoodenBridge || t == tileStoneBridge {
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
			if g.WoodsRun.Themes[level] == WoodsPeasantHouses {
				if doors < 2 {
					t.Fatalf("level %d seed %d has only %d cottage doors", level, seed, doors)
				}
				if windows == 0 || furniture == 0 {
					t.Fatalf("level %d seed %d cottage is missing windows/furniture", level, seed)
				}
			}
		}
	}
}

func TestWoodsLandmarksIdentifyEachTheme(t *testing.T) {
	want := map[int][]uint8{
		WoodsChurch:        {tileBellTower, tileChapelAltar, tileGravestone},
		WoodsPeasantHouses: {tileWell, tileMarketStall, tileCart, tileHaystack},
		WoodsMillPond:      {tileMillStone, tileBarrel, tileCart},
		WoodsCemetery:      {tileIronFence, tileCryptDoor, tileGravestone},
		WoodsMeteorField:   {tileMeteorRock, tileBurnedGround, tileDeadTree},
	}
	for seed := int64(0); seed < 8; seed++ {
		g := NewGame(seed + 5000)
		g.SetupNewGame(0, "tester")
		for level := 0; level < numWoodsMaps; level++ {
			g.GenerateLevel(level)
			theme := g.WoodsRun.Themes[level]
			counts := map[uint8]int{}
			for _, tile := range g.Map.Tiles {
				counts[tile]++
			}
			for _, tile := range want[theme] {
				if counts[tile] == 0 {
					t.Fatalf("seed %d level %d theme %d missing landmark tile %d", seed, level, theme, tile)
				}
			}
			switch theme {
			case WoodsChurch:
				if counts[tileBellTower] != 1 || counts[tileChapelAltar] != 1 {
					t.Fatalf("seed %d Church has %d bell towers and %d altars", seed, counts[tileBellTower], counts[tileChapelAltar])
				}
			case WoodsPeasantHouses:
				if counts[tileWell] != 1 {
					t.Fatalf("seed %d Peasant Houses has %d wells", seed, counts[tileWell])
				}
			case WoodsCemetery:
				if counts[tileCryptDoor] != 1 {
					t.Fatalf("seed %d Cemetery has %d crypt doors", seed, counts[tileCryptDoor])
				}
			case WoodsMeteorField:
				if counts[tileMeteorRock] < 3 {
					t.Fatalf("seed %d Meteor Field has only %d meteor rocks", seed, counts[tileMeteorRock])
				}
			}
			centralTile := map[int]uint8{
				WoodsChurch:        tileBellTower,
				WoodsPeasantHouses: tileWell,
				WoodsMillPond:      tileMillStone,
				WoodsCemetery:      tileCryptDoor,
				WoodsMeteorField:   tileMeteorRock,
			}[theme]
			foundCentral := false
			for y := 0; y < g.Map.H && !foundCentral; y++ {
				for x := 0; x < g.Map.W; x++ {
					if g.Map.Tile(x, y) == centralTile && abs(x-g.Map.LandmarkX)+abs(y-g.Map.LandmarkY) <= 10 {
						foundCentral = true
						break
					}
				}
			}
			if !foundCentral {
				t.Fatalf("seed %d level %d theme %d has no central identity tile", seed, level, theme)
			}
			if g.Map.LandmarkX < g.Map.W/4 || g.Map.LandmarkX > g.Map.W*3/4 {
				t.Fatalf("seed %d level %d landmark is not central: x=%d", seed, level, g.Map.LandmarkX)
			}
			if g.Map.LandmarkY < 2 || g.Map.LandmarkY >= g.Map.H-2 {
				t.Fatalf("seed %d level %d landmark y out of bounds: %d", seed, level, g.Map.LandmarkY)
			}
		}
	}
}

func TestAllWoodsLevelsRender(t *testing.T) {
	for level := 0; level < numWoodsMaps; level++ {
		g := NewGame(int64(91000 + level))
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(level)
		g.State = StatePlaying
		g.UpdateFOV()
		for _, mode := range []ColorMode{ColorModeTrueColor, ColorMode16} {
			u, s := newTestUI(g, mode)
			u.drawGame()
			s.Show()
			text := screenText(s)
			if !strings.Contains(text, "@") {
				t.Fatalf("level %d mode %s did not render the player", level, mode)
			}
			rows := strings.Split(text, "\n")
			if len(rows) != screenH+1 {
				t.Fatalf("level %d mode %s rendered %d rows, want %d", level, mode, len(rows)-1, screenH)
			}
			s.Fini()
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
