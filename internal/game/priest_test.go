package game

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestPriestProfileGeneration(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		g := NewGame(seed)
		if g.PlayerName == "" {
			t.Fatalf("seed %d has no priest name", seed)
		}
		if len(g.KnownRites()) != 1 {
			t.Fatalf("seed %d has %d starting rites", seed, len(g.KnownRites()))
		}
		if g.PlayerRites[1] != RiteNone {
			t.Fatalf("seed %d starts with a second rite", seed)
		}
		if g.RiteCharges[0] != 1 {
			t.Fatalf("seed %d starting charge = %d", seed, g.RiteCharges[0])
		}
	}
}

func TestRiteChargesRefreshOnNewMap(t *testing.T) {
	g := NewGame(15100)
	g.SetupNewGame(0, "tester")
	g.State = StatePlaying
	g.PlayerRites[0] = RiteHealingWounds
	g.RiteCharges[0] = 1
	g.Player.HP = 1
	g.UseRite(0)
	if g.RiteCharges[0] != 0 {
		t.Fatalf("rite charge = %d after use", g.RiteCharges[0])
	}
	g.GenerateLevel(1)
	if g.RiteCharges[0] != 1 {
		t.Fatalf("rite charge = %d after map change, want 1", g.RiteCharges[0])
	}
}

func TestNoBaselineAutoHealing(t *testing.T) {
	g := NewGame(15200)
	g.SetupNewGame(0, "tester")
	g.State = StatePlaying
	g.Options["optGameUseAI"] = 0
	g.Player.HP = g.Player.MaxHP - 5
	for i := 0; i < 25; i++ {
		g.WaitTurn()
	}
	if g.Player.HP != g.Player.MaxHP-5 {
		t.Fatalf("baseline healing changed HP: %d -> %d", g.Player.MaxHP-5, g.Player.HP)
	}
}

func TestHealingWoundsRiteConsumesCharge(t *testing.T) {
	g := NewGame(15300)
	g.SetupNewGame(0, "tester")
	g.State = StatePlaying
	g.PlayerRites[0] = RiteHealingWounds
	g.RiteCharges[0] = 1
	g.Player.HP = 10
	g.UseRite(0)
	if g.Player.HP != 15 {
		t.Fatalf("Healing Wounds restored to %d, want 15", g.Player.HP)
	}
	if g.RiteCharges[0] != 0 {
		t.Fatalf("charge after use = %d", g.RiteCharges[0])
	}
	g.Player.HP = 10
	g.UseRite(0)
	if g.Player.HP != 10 {
		t.Fatal("rite healed again with no charges")
	}
}

func TestXPGivesNoLevels(t *testing.T) {
	g := NewGame(15350)
	g.SetupNewGame(0, "tester")
	g.Player.XP = 100000
	oldXP := g.Player.XP
	oldLevel := g.Player.Level
	target := g.SpawnCritter(0, g.Player.X+1, g.Player.Y)
	g.DamageCritter(target, 999, g.Player, 0)
	if g.Player.XP <= oldXP {
		t.Fatal("XP was not awarded")
	}
	if g.Player.Level != oldLevel {
		t.Fatalf("player level changed from %d to %d", oldLevel, g.Player.Level)
	}
}

func TestIntroStartsGameWithoutClassSelection(t *testing.T) {
	g := NewGame(15425)
	g.State = StateIntro
	u := &UI{g: g}
	u.handleIntro(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
	if u.g.State != StatePlaying {
		t.Fatalf("intro enter state = %v, want playing", u.g.State)
	}
	if u.g.Map == nil || u.g.Player == nil {
		t.Fatal("intro did not start the game")
	}
}

func TestIntroRerollRegeneratesEverything(t *testing.T) {
	g := NewGame(15450)
	g.State = StateIntro
	u := &UI{g: g}
	old := g
	u.handleIntro(tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone))
	if u.g == old {
		t.Fatal("reroll did not create a new game state")
	}
	if u.g.State != StateIntro {
		t.Fatalf("reroll state = %v, want intro", u.g.State)
	}
	if len(u.g.KnownRites()) != 1 {
		t.Fatal("rerolled priest has no starting rite")
	}
}

func TestBlessingRestoresShardAP(t *testing.T) {
	g := NewGame(15400)
	g.SetupNewGame(0, "tester")
	g.State = StatePlaying
	g.PlayerRites[0] = RiteBlessing
	g.RiteCharges[0] = 1
	shard := g.Player.Slots.Robe[0]
	ensureShardAP(shard)
	shard.AP = 1
	g.UseRite(0)
	if shard.AP != 3 {
		t.Fatalf("Blessing restored shard to %d AP, want 3", shard.AP)
	}
	if g.RiteCharges[0] != 0 {
		t.Fatal("Blessing did not consume its charge")
	}
}
