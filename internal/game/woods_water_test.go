package game

import "testing"

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
			if g.WoodsRun.Themes[level] == WoodsChurch && water == 0 {
				t.Fatalf("seed %d Church map has no pond or rain pool", seed)
			}
			totalBridges += bridges
		}
	}
	if totalBridges == 0 {
		t.Fatal("no woods level generated a bridge")
	}
}
