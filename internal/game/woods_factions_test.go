package game

import "testing"

func clearAround(g *Game, cx, cy, radius int) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			if g.Map.InBounds(x, y) {
				g.Map.SetTile(x, y, 2)
			}
		}
	}
}

func TestWoodsSpawnsVillagersAndAnimals(t *testing.T) {
	for level := 0; level < numWoodsMaps; level++ {
		for seed := int64(0); seed < 5; seed++ {
			g := NewGame(seed + int64(level)*7000)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(level)
			villagers, animals := 0, 0
			for _, c := range g.Critters {
				switch c.Align {
				case AlignVillager:
					villagers++
				case AlignNeutral:
					if c.IsAnimal {
						animals++
					}
				}
			}
			if villagers == 0 {
				t.Fatalf("level %d seed %d has no villagers", level, seed)
			}
			if animals == 0 {
				t.Fatalf("level %d seed %d has no neutral animals", level, seed)
			}
		}
	}
}

func TestZombieTargetingAndVillagerConversion(t *testing.T) {
	g := NewGame(7100)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	g.State = StatePlaying
	clearAround(g, g.Player.X, g.Player.Y, 8)

	zombie := g.SpawnCritter(0, g.Player.X+3, g.Player.Y)
	villager := g.SpawnCritter(villagerProto, g.Player.X+8, g.Player.Y)
	if zombie == nil || villager == nil {
		t.Fatal("spawn failed")
	}
	if target, _ := g.nearestHostile(zombie, 2); target != nil {
		t.Fatalf("zombie should ignore distant targets, got %s", target.NameOrUnknown())
	}
	if target, _ := g.nearestHostile(zombie, 6); target == nil {
		t.Fatal("zombie did not find a target within range")
	}

	g.Options["optMechanicBaseHitChance"] = 100
	villager.HP, villager.MaxHP = 999, 999
	for i := 0; i < 30 && villager.HP == 999; i++ {
		zombie.X, zombie.Y = g.Player.X+3, g.Player.Y
		villager.X, villager.Y = g.Player.X+4, g.Player.Y
		g.MonsterTurn(zombie)
	}
	if villager.HP == 999 {
		t.Fatal("zombie did not attack the adjacent villager")
	}

	villager.HP = 1
	g.DamageCritter(villager, 999, zombie, 0)
	if len(g.PendingRises) != 1 {
		t.Fatalf("pending rises = %d, want 1", len(g.PendingRises))
	}
	g.PendingRises[0].Turns = 1
	g.processPendingRises()
	found := false
	for _, c := range g.Critters {
		if c.Proto == 0 && chebyshev(c.X, c.Y, g.Player.X+4, g.Player.Y) <= 2 {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("villager did not rise as a zombie")
	}
}

func TestNeutralAnimalFleesInsteadOfAttacking(t *testing.T) {
	g := NewGame(7101)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	g.State = StatePlaying
	clearAround(g, g.Player.X, g.Player.Y, 8)

	rabbit := g.SpawnCritter(rabbitProto, g.Player.X+2, g.Player.Y)
	if rabbit == nil {
		t.Fatal("rabbit spawn failed")
	}
	threat, _ := g.nearestAnimalThreat(rabbit, 4)
	if threat != g.Player {
		t.Fatal("rabbit does not see the player as a threat")
	}
	hp := g.Player.HP
	for i := 0; i < 10; i++ {
		g.MonsterTurn(rabbit)
	}
	if g.Player.HP != hp {
		t.Fatal("neutral rabbit attacked the player")
	}
}
