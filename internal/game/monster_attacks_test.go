package game

import (
	"strings"
	"testing"
)

func TestEveryMonsterHasDamageDice(t *testing.T) {
	g := NewGame(17500)
	g.SetupNewGame(0, "tester")
	for proto := range CritterDefs {
		c := g.SpawnCritter(proto, 1, 1)
		if c == nil {
			t.Fatalf("critter proto %d did not spawn", proto)
		}
		if c.DmgCount < 1 || c.DmgSides < 2 {
			t.Fatalf("critter %q has no damage dice: %s", c.Name, c.DamageDescription())
		}
	}
}

func TestMonsterDamageUsesDiceAndAlwaysAtLeastOne(t *testing.T) {
	g := NewGame(17501)
	g.SetupNewGame(0, "tester")
	zombie := g.SpawnCritter(0, g.Player.X+2, g.Player.Y)
	if got := zombie.DamageDescription(); got != "1d6+1" {
		t.Fatalf("zombie damage = %q, want 1d6+1", got)
	}
	for i := 0; i < 200; i++ {
		if dmg := g.rollMonsterDamage(zombie); dmg < 2 || dmg > 7 {
			t.Fatalf("zombie damage roll out of range: %d", dmg)
		}
	}
	zombie.DmgCount, zombie.DmgSides, zombie.DmgBonus = 1, 6, -20
	if dmg := g.rollMonsterDamage(zombie); dmg != 1 {
		t.Fatalf("negative damage was not clamped: %d", dmg)
	}
}

func TestMonstersHaveNoArmour(t *testing.T) {
	g := NewGame(17502)
	g.SetupNewGame(0, "tester")
	for i := 0; i < 20; i++ {
		c := g.SpawnCritter(0, g.Player.X+2+i, g.Player.Y)
		if c == nil {
			t.Fatal("spawn failed")
		}
		if c.Equip[EquipArmor] != nil {
			t.Fatal("monster spawned with armour")
		}
	}
	target := g.SpawnCritter(0, g.Player.X+1, g.Player.Y)
	target.Equip[EquipArmor] = &Item{Proto: itemProtoByName("padded jerkin")}
	target.HP, target.MaxHP = 100, 100
	g.DamageCritter(target, 10, g.Player, 0)
	if target.HP != 90 {
		t.Fatalf("monster armour reduced damage: HP %d, want 90", target.HP)
	}
}

func TestCombatSetsAndClearsCurrentTarget(t *testing.T) {
	g := NewGame(17504)
	g.SetupNewGame(0, "tester")
	target := g.SpawnCritter(0, g.Player.X+1, g.Player.Y)
	g.DamageCritter(target, 1, g.Player, 0)
	if g.CurrentTarget != target {
		t.Fatal("attacking did not select the current target")
	}
	target.HP = 1
	g.DamageCritter(target, 10, g.Player, 0)
	if g.CurrentTarget != nil {
		t.Fatal("current target was not cleared after death")
	}
}

func TestHUDShowsCurrentTarget(t *testing.T) {
	g := NewGame(17503)
	g.SetupNewGame(0, "tester")
	zombie := g.SpawnCritter(0, g.Player.X+1, g.Player.Y)
	g.CurrentTarget = zombie
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.drawStatus(hudRow)
	s.Show()
	rows := strings.Split(screenText(s), "\n")
	hud := rows[hudRow]
	for _, want := range []string{"zombie", "1d6+1", "HP ", "z"} {
		if !strings.Contains(hud, want) {
			t.Fatalf("HUD missing target %q: %q", want, hud)
		}
	}
}
