package game

import "testing"

func TestWoodsRiverIsEdgeToEdge(t *testing.T) {
	g := NewGame(5150)
	m := NewMap(mapWidth, mapHeight)
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

func TestWoodsGeneratorWaterAndBridges(t *testing.T) {
	totalBridges := 0
	for level := 0; level < 9; level++ {
		for seed := int64(0); seed < 10; seed++ {
			g := NewGame(seed + int64(level)*3000)
			g.SetupNewGame(0, "tester")
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
						adjacentWater := false
						for _, d := range [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
							if isWoodsWater(g.Map.Tile(x+d[0], y+d[1])) {
								adjacentWater = true
								break
							}
						}
						if !adjacentWater {
							t.Fatalf("level %d seed %d bridge at %d,%d has no adjacent water", level, seed, x, y)
						}
					}
				}
			}
			if water == 0 {
				t.Fatalf("level %d seed %d has no water", level, seed)
			}
			if banks == 0 {
				t.Fatalf("level %d seed %d has no river banks", level, seed)
			}
			totalBridges += bridges
		}
	}
	if totalBridges == 0 {
		t.Fatal("no woods level generated a bridge")
	}
}
