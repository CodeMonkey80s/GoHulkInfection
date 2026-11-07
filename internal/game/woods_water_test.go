package game

import "testing"

func TestWoodsRiverIsEdgeToEdge(t *testing.T) {
	g := NewGame(5150)
	m := NewMap(mapWidth, woodsMapHeight)
	for i := range m.Tiles {
		m.Tiles[i] = 2
	}
	g.generateWoodsRiver(m)

	top, bottom, left, right := 0, 0, 0, 0
	for x := 0; x < m.W; x++ {
		if isWoodsWater(m.Tile(x, 0)) {
			top++
		}
		if isWoodsWater(m.Tile(x, m.H-1)) {
			bottom++
		}
	}
	for y := 0; y < m.H; y++ {
		if isWoodsWater(m.Tile(0, y)) {
			left++
		}
		if isWoodsWater(m.Tile(m.W-1, y)) {
			right++
		}
	}
	if !((top > 0 && bottom > 0) || (left > 0 && right > 0)) {
		t.Fatalf("river does not reach opposite edges: top=%d bottom=%d left=%d right=%d", top, bottom, left, right)
	}
}

func TestWoodsThemeWaterAndBridges(t *testing.T) {
	totalBridges := 0
	for seed := int64(0); seed < 12; seed++ {
		g := NewGame(seed + 3000)
		g.SetupNewGame(0, "tester")
		for level := 0; level < numWoodsMaps; level++ {
			g.GenerateLevel(level)
			water, banks, bridges := 0, 0, 0
			for y := 0; y < g.Map.H; y++ {
				for x := 0; x < g.Map.W; x++ {
					tile := g.Map.Tile(x, y)
					switch {
					case isWoodsWater(tile):
						water++
					case tile == tileRiverBank:
						banks++
					case tile == tileWoodenBridge:
						bridges++
						if g.Map.BlocksMove(x, y) {
							t.Fatalf("level %d seed %d bridge at %d,%d is impassable", level, seed, x, y)
						}
					}
				}
			}
			if water > 0 && banks == 0 {
				t.Fatalf("level %d seed %d has water without banks", level, seed)
			}
			if g.WoodsRun.Themes[level] == WoodsMillPond {
				if water == 0 {
					t.Fatalf("seed %d Mill map has no water", seed)
				}
				if bridges == 0 {
					t.Fatalf("seed %d Mill map has no road bridge", seed)
				}
			}
			totalBridges += bridges
		}
	}
	if totalBridges == 0 {
		t.Fatal("no woods level generated a bridge")
	}
}
