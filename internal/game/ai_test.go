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

func TestAIHostileTargetingRequiresLineOfSight(t *testing.T) {
	blockers := []struct {
		name string
		tile uint8
	}{
		{"dense tree", tileDenseTree},
		{"bramble", tileBramble},
		{"bush", tileBush},
		{"stone wall", 10},
		{"closed door", tileWoodenDoor},
	}
	for _, tc := range blockers {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGame(7303)
			g.SetupNewGame(0, "tester")
			g.Critters = nil
			g.State = StatePlaying
			clearAround(g, g.Player.X, g.Player.Y, 8)
			zombie := g.SpawnCritter(0, g.Player.X+2, g.Player.Y)
			wallX, wallY := g.Player.X+1, g.Player.Y
			g.Map.SetTile(wallX, wallY, tc.tile)
			if target, _ := g.nearestHostile(zombie, 4); target != nil {
				t.Fatalf("zombie saw a target through %s", tc.name)
			}
			g.Map.SetTile(wallX, wallY, 2)
			if target, _ := g.nearestHostile(zombie, 4); target == nil {
				t.Fatalf("zombie failed to see a clear target after removing %s", tc.name)
			}
		})
	}
}

func TestAIAnimalThreatRequiresLineOfSight(t *testing.T) {
	g := NewGame(7304)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	clearAround(g, g.Player.X, g.Player.Y, 8)
	rabbit := g.SpawnCritter(rabbitProto, g.Player.X+2, g.Player.Y)
	g.Map.SetTile(g.Player.X+1, g.Player.Y, tileDenseTree)
	if threat, _ := g.nearestAnimalThreat(rabbit, 4); threat != nil {
		t.Fatal("animal sensed a threat through a dense tree")
	}
	g.Map.SetTile(g.Player.X+1, g.Player.Y, 2)
	if threat, _ := g.nearestAnimalThreat(rabbit, 4); threat != g.Player {
		t.Fatal("animal did not sense a clear threat")
	}
}

func TestMonstersDoNotAttackDiagonally(t *testing.T) {
	g := NewGame(7306)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	g.State = StatePlaying
	clearAround(g, g.Player.X, g.Player.Y, 4)
	zombie := g.SpawnCritter(0, g.Player.X+1, g.Player.Y+1)
	hp := g.Player.HP
	for i := 0; i < 3; i++ {
		zombie.X, zombie.Y = g.Player.X+1, g.Player.Y+1
		g.MonsterTurn(zombie)
	}
	if g.Player.HP != hp {
		t.Fatalf("monster attacked diagonally: HP %d -> %d", hp, g.Player.HP)
	}
}

func TestPlayerDoesNotAttackDiagonally(t *testing.T) {
	g := NewGame(7307)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	g.State = StatePlaying
	clearAround(g, g.Player.X, g.Player.Y, 4)
	target := g.SpawnCritter(0, g.Player.X+1, g.Player.Y+1)
	hp := target.HP
	g.PlayerDirectionalAttack(1, 1)
	if target.HP != hp {
		t.Fatal("player attacked diagonally")
	}
	g.PlayerMelee(target)
	if target.HP != hp {
		t.Fatal("PlayerMelee allowed a diagonal attack")
	}
}

func TestObserverUsesLineOfSight(t *testing.T) {
	g := NewGame(7305)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	clearAround(g, g.Player.X, g.Player.Y, 8)
	observer := g.SpawnCritter(8, g.Player.X+2, g.Player.Y) // observer drone
	g.Map.SetTile(g.Player.X+1, g.Player.Y, 10)
	if _, canSee := g.canSeeCritter(observer, g.Player, 14); canSee {
		t.Fatal("observer saw the player through a stone wall")
	}
	g.Map.SetTile(g.Player.X+1, g.Player.Y, 2)
	if _, canSee := g.canSeeCritter(observer, g.Player, 14); !canSee {
		t.Fatal("observer did not see the player through clear ground")
	}
}
