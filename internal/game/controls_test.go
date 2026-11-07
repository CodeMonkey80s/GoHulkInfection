package game

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestMovementKeysAndNoNumpad(t *testing.T) {
	cases := []struct {
		key rune
		dx  int
		dy  int
	}{
		{'h', -1, 0}, {'j', 0, 1}, {'k', 0, -1}, {'l', 1, 0},
	}
	for _, tc := range cases {
		dx, dy, ok := directionFromEvent(tcell.NewEventKey(tcell.KeyRune, tc.key, tcell.ModNone))
		if !ok || dx != tc.dx || dy != tc.dy {
			t.Fatalf("%q direction = %d,%d,%v, want %d,%d,true", tc.key, dx, dy, ok, tc.dx, tc.dy)
		}
	}
	for _, key := range []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'u', 'i', 'o', 'm', ',', '.'} {
		if _, _, ok := directionFromEvent(tcell.NewEventKey(tcell.KeyRune, key, tcell.ModNone)); ok {
			t.Fatalf("key %q still produces a movement direction", key)
		}
	}
}

func TestActiveBindingsUseNewKeys(t *testing.T) {
	found := map[rune]bool{}
	for _, b := range ActiveKeyBindings {
		found[rune(b.Key)] = true
	}
	for _, key := range []rune{'h', 'j', 'k', 'l', '.', 'e', 'v', '~', 'g', 'r', 'R'} {
		if !found[key] {
			t.Fatalf("active key bindings missing %q", key)
		}
	}
	for _, key := range []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'} {
		if found[key] {
			t.Fatalf("active key bindings still contain numpad key %q", key)
		}
	}
}

func TestContextualBash(t *testing.T) {
	g := NewGame(7200)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	g.State = StatePlaying
	clearAround(g, g.Player.X, g.Player.Y, 6)
	g.Player.STR = 50
	g.Options["optPlayerGodMode"] = 1

	x, y := g.Player.X+1, g.Player.Y
	g.Map.SetTile(x, y, tileBramble)
	if !g.tryBashTile(x, y) {
		t.Fatal("bramble was not recognised as bashable")
	}
	if g.Map.Tile(x, y) == tileBramble {
		t.Fatal("bramble survived a guaranteed bash")
	}
	g.Map.SetTile(x, y, 10) // stone wall
	if g.tryBashTile(x, y) {
		t.Fatal("stone wall should not be contextually bashable")
	}

	g.Map.SetTile(x, y, tileBramble)
	for i := 0; i < 10 && g.Map.Tile(x, y) == tileBramble; i++ {
		g.MovePlayer(1, 0)
	}
	if g.Map.Tile(x, y) == tileBramble {
		t.Fatal("walking into a bramble did not attempt a bash")
	}
}

func TestMenuJKNavigation(t *testing.T) {
	g := NewGame(7201)
	u := &UI{g: g}
	u.handleTitle(tcell.NewEventKey(tcell.KeyRune, 'j', tcell.ModNone))
	if u.titleSel != 1 {
		t.Fatalf("title j selection = %d, want 1", u.titleSel)
	}
	u.handleTitle(tcell.NewEventKey(tcell.KeyRune, 'k', tcell.ModNone))
	if u.titleSel != 0 {
		t.Fatalf("title k selection = %d, want 0", u.titleSel)
	}

	g.SetupNewGame(0, "tester")
	u.handleCharCreate(tcell.NewEventKey(tcell.KeyRune, 'j', tcell.ModNone))
	if u.classSel != 1 {
		t.Fatalf("class j selection = %d, want 1", u.classSel)
	}
	u.handleCharCreate(tcell.NewEventKey(tcell.KeyRune, 'k', tcell.ModNone))
	if u.classSel != 0 {
		t.Fatalf("class k selection = %d, want 0", u.classSel)
	}

	g.State = StateInventory
	u.slotSel = 0
	u.handleInventory(tcell.NewEventKey(tcell.KeyRune, 'j', tcell.ModNone))
	if u.slotSel != 1 {
		t.Fatalf("equipment j selection = %d, want 1", u.slotSel)
	}
	u.handleInventory(tcell.NewEventKey(tcell.KeyRune, 'k', tcell.ModNone))
	if u.slotSel != 0 {
		t.Fatalf("equipment k selection = %d, want 0", u.slotSel)
	}
}

func TestMonsterPathfindingIsCardinal(t *testing.T) {
	g := NewGame(7202)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	clearAround(g, g.Player.X, g.Player.Y, 12)
	zombie := g.SpawnCritter(0, g.Player.X+3, g.Player.Y+3)
	if zombie == nil {
		t.Fatal("zombie spawn failed")
	}
	dx, dy, ok := g.findStep(zombie, g.Player.X+8, g.Player.Y+6)
	if !ok {
		t.Fatal("no cardinal path found")
	}
	if dx != zombie.X && dy != zombie.Y {
		t.Fatalf("path step is diagonal: from %d,%d to %d,%d", zombie.X, zombie.Y, dx, dy)
	}
	if dx == zombie.X && dy == zombie.Y {
		t.Fatal("path step did not move")
	}
}

func TestBashDirectionKey(t *testing.T) {
	g := NewGame(7203)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	g.State = StatePlaying
	clearAround(g, g.Player.X, g.Player.Y, 6)
	g.Player.STR = 50
	x, y := g.Player.X+1, g.Player.Y
	g.Map.SetTile(x, y, tileBramble)
	g.PlayerBashDirection(1, 0)
	if g.Map.Tile(x, y) == tileBramble {
		t.Fatal("f + direction did not bash the bramble")
	}
}

func TestBumpAttack(t *testing.T) {
	g := NewGame(7204)
	g.SetupNewGame(0, "tester")
	g.Critters = nil
	g.State = StatePlaying
	clearAround(g, g.Player.X, g.Player.Y, 8)
	g.Options["optPlayerGodMode"] = 1
	g.Options["optMechanicBaseHitChance"] = 100
	monster := g.SpawnCritter(0, g.Player.X+1, g.Player.Y)
	monster.HP, monster.MaxHP = 999, 999
	for i := 0; i < 30 && monster.HP == 999; i++ {
		monster.X, monster.Y = g.Player.X+1, g.Player.Y
		g.MovePlayer(1, 0)
	}
	if monster.HP == 999 {
		t.Fatal("moving into a monster did not attack it")
	}
}

func TestSpaceClosesModals(t *testing.T) {
	states := []GameState{
		StateInventory, StateLook, StateThrow, StateMessages, StateOptions,
		StateSkills, StateHelp, StatePlayerScreen, StateQuitConfirm,
	}
	for _, state := range states {
		g := NewGame(7205)
		g.SetupNewGame(0, "tester")
		g.State = state
		u := &UI{g: g}
		u.handleKey(tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone))
		if g.State != StatePlaying {
			t.Fatalf("space did not close state %v (now %v)", state, g.State)
		}
	}
}
