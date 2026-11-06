package game

import "testing"

func findPortal(g *Game) (int, int) {
	for y := 0; y < g.Map.H; y++ {
		for x := 0; x < g.Map.W; x++ {
			switch g.Map.Tile(x, y) {
			case 23, 24, 59: // stairs, time portal, elevator
				return x, y
			}
		}
	}
	return -1, -1
}

func reachable(g *Game, sx, sy, tx, ty int) bool {
	if !g.Map.InBounds(tx, ty) {
		return false
	}
	seen := make([]bool, g.Map.W*g.Map.H)
	q := []int{g.Map.idx(sx, sy)}
	seen[g.Map.idx(sx, sy)] = true
	for head := 0; head < len(q); head++ {
		cur := q[head]
		cx, cy := cur%g.Map.W, cur/g.Map.W
		if cx == tx && cy == ty {
			return true
		}
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
	return false
}

func TestMapsHaveNoVoidAndReachPortal(t *testing.T) {
	levels := []int{0, 5, 9, 14, 18, 23, 27}
	for _, level := range levels {
		for seed := int64(0); seed < 8; seed++ {
			g := NewGame(seed + int64(level)*100)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(level)
			voids := 0
			for _, tile := range g.Map.Tiles {
				if tile == 0 || tile == 1 {
					voids++
				}
			}
			if voids > 0 {
				t.Fatalf("level %d seed %d has %d void tiles", level, seed, voids)
			}
			px, py := findPortal(g)
			if px < 0 {
				t.Fatalf("level %d seed %d has no portal", level, seed)
			}
			if !reachable(g, g.Player.X, g.Player.Y, px, py) {
				t.Fatalf("level %d seed %d portal at %d,%d not reachable from %d,%d", level, seed, px, py, g.Player.X, g.Player.Y)
			}
		}
	}
}
