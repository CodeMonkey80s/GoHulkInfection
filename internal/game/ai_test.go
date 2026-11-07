package game

import "testing"

func TestPathfindingAroundWallAndWater(t *testing.T) {
	g := NewGame(7300)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	clearAround(g, g.Player.X, g.Player.Y, 14)
	zombie := g.SpawnCritter(0, g.Player.X, g.Player.Y)
	if zombie == nil {
		t.Fatal("zombie spawn failed")
	}
	tx, ty := g.Player.X+4, g.Player.Y
	for y := g.Player.Y - 2; y <= g.Player.Y+2; y++ {
		g.Map.SetTile(g.Player.X+2, y, 10) // stone wall
	}
	g.Map.SetTile(g.Player.X+1, g.Player.Y+1, tileDeepWater)
	path := g.findPath(zombie, tx, ty)
	if len(path) == 0 {
		t.Fatal("no path around wall")
	}
	if path[len(path)-1].x != tx || path[len(path)-1].y != ty {
		t.Fatal("path does not end at the target")
	}
	for _, p := range path {
		if g.Map.BlocksMove(p.x, p.y) || g.Map.Tile(p.x, p.y) == tileDeepWater {
			t.Fatalf("path crosses blocked tile %d,%d", p.x, p.y)
		}
	}
}

func TestPathfindingAvoidsFireWhenPossible(t *testing.T) {
	g := NewGame(7301)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	clearAround(g, g.Player.X, g.Player.Y, 10)
	zombie := g.SpawnCritter(0, g.Player.X, g.Player.Y)
	tx, ty := g.Player.X+3, g.Player.Y
	fireX, fireY := g.Player.X+1, g.Player.Y
	g.Map.Fire[g.Map.idx(fireX, fireY)] = 20
	path := g.findPath(zombie, tx, ty)
	if len(path) == 0 {
		t.Fatal("no path around fire")
	}
	for _, p := range path {
		if p.x == fireX && p.y == fireY {
			t.Fatal("path stepped into fire despite a safe route")
		}
	}
}

func TestPathfindingOpensDoors(t *testing.T) {
	g := NewGame(7302)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	clearAround(g, g.Player.X, g.Player.Y, 8)
	zombie := g.SpawnCritter(0, g.Player.X, g.Player.Y)
	doorX, doorY := g.Player.X+1, g.Player.Y
	g.Map.SetTile(doorX, doorY, tileWoodenDoor)
	g.stepToward(zombie, g.Player.X+2, g.Player.Y)
	if g.Map.Tile(doorX, doorY) == tileWoodenDoor {
		t.Fatal("monster did not open the door")
	}
}
