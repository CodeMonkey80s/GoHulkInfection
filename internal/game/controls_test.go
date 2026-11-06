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
		{'u', -1, -1}, {'i', 0, -1}, {'o', 1, -1},
		{'j', -1, 0}, {'l', 1, 0},
		{'m', -1, 1}, {',', 0, 1}, {'.', 1, 1},
	}
	for _, tc := range cases {
		dx, dy, ok := directionFromEvent(tcell.NewEventKey(tcell.KeyRune, tc.key, tcell.ModNone))
		if !ok || dx != tc.dx || dy != tc.dy {
			t.Fatalf("%q direction = %d,%d,%v, want %d,%d,true", tc.key, dx, dy, ok, tc.dx, tc.dy)
		}
	}
	for _, key := range []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'} {
		if _, _, ok := directionFromEvent(tcell.NewEventKey(tcell.KeyRune, key, tcell.ModNone)); ok {
			t.Fatalf("numpad key %q still produces a direction", key)
		}
	}
}

func TestActiveBindingsUseNewKeys(t *testing.T) {
	found := map[rune]bool{}
	for _, b := range ActiveKeyBindings {
		found[rune(b.Key)] = true
	}
	for _, key := range []rune{'u', 'i', 'o', 'j', 'k', 'l', 'm', ',', '.', 'e', 'v', '~', 'g', 'r', 'R'} {
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
