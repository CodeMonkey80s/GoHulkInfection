package game

import (
	"path/filepath"
	"testing"
)

func TestSavePreservesPortalAndPosition(t *testing.T) {
	g := NewGame(7)
	g.SetupNewGame(0, "t")
	px, py := -1, -1
	for y := 0; y < g.Map.H && px < 0; y++ {
		for x := 0; x < g.Map.W; x++ {
			if g.Map.Tile(x, y) == 23 {
				px, py = x, y
				break
			}
		}
	}
	if px < 0 {
		t.Fatal("no portal generated")
	}
	g.Player.X, g.Player.Y = px, py
	g.UpdateFOV()
	path := filepath.Join(t.TempDir(), "s.json")
	if err := g.SaveGame(path); err != nil {
		t.Fatal(err)
	}
	g2 := NewGame(8)
	if err := g2.LoadGame(path); err != nil {
		t.Fatal(err)
	}
	if g2.Player.X != px || g2.Player.Y != py {
		t.Fatalf("position lost: %d,%d want %d,%d", g2.Player.X, g2.Player.Y, px, py)
	}
	if g2.Map.Tile(px, py) != 23 {
		t.Fatalf("portal tile lost: %d", g2.Map.Tile(px, py))
	}
}
