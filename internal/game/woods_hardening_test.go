package game

import (
	"strings"
	"testing"
)

func TestWoodsMapInvariantMatrix(t *testing.T) {
	seed := int64(17000)
	for theme := 0; theme < numWoodsMaps; theme++ {
		for entry := EdgeNorth; entry <= EdgeWest; entry++ {
			for exit := EdgeNorth; exit <= EdgeWest; exit++ {
				if entry == exit {
					continue
				}
				for repeat := 0; repeat < 4; repeat++ {
					g := NewGame(seed)
					g.SetupNewGame(0, "tester")
					g.generateWoodsMap(0, theme, entry, exit)
					seed++
					if problems := woodsMapProblems(g.Map, theme, entry, exit, g.Player.X, g.Player.Y); len(problems) > 0 {
						t.Fatalf("theme %d entry %d exit %d seed %d: %s", theme, entry, exit, seed-1, strings.Join(problems, "; "))
					}
				}
			}
		}
	}
}

func TestWoodsValidatorDetectsBrokenMaps(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(g *Game)
	}{
		{"extra border gap", func(g *Game) { g.Map.SetTile(0, 0, 2) }},
		{"missing portal", func(g *Game) {
			px, py := findPortal(g)
			g.Map.SetTile(px, py, 2)
		}},
		{"void tile", func(g *Game) { g.Map.SetTile(10, 10, 0) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGame(17500)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(0)
			tc.mutate(g)
			problems := woodsMapProblems(g.Map, g.WoodsRun.Themes[0], g.WoodsRun.Entries[0], g.WoodsRun.Exits[0], g.Player.X, g.Player.Y)
			if len(problems) == 0 {
				t.Fatalf("validator accepted a map with %s", tc.name)
			}
		})
	}
}

func TestWoodsRunMapsPassInvariants(t *testing.T) {
	for seed := int64(0); seed < 40; seed++ {
		g := NewGame(seed + 18000)
		g.SetupNewGame(0, "tester")
		for level := 0; level < numWoodsMaps; level++ {
			g.GenerateLevel(level)
			theme := g.WoodsRun.Themes[level]
			if problems := woodsMapProblems(g.Map, theme, g.WoodsRun.Entries[level], g.WoodsRun.Exits[level], g.Player.X, g.Player.Y); len(problems) > 0 {
				t.Fatalf("seed %d level %d: %s", seed, level, strings.Join(problems, "; "))
			}
		}
	}
}

func TestWoodsRunIsWalkableToEachPortal(t *testing.T) {
	for seed := int64(0); seed < 4; seed++ {
		g := NewGame(seed + 19500)
		g.SetupNewGame(0, "tester")
		g.State = StatePlaying
		g.Options["optGameUseAI"] = 0
		for level := 0; level < numWoodsMaps; level++ {
			g.Critters = nil
			px, py := findPortal(g)
			path := g.findPath(g.Player, px, py)
			if len(path) == 0 {
				t.Fatalf("seed %d level %d has no walkable path to the portal", seed, level)
			}
			for _, p := range path {
				dx, dy := p.x-g.Player.X, p.y-g.Player.Y
				beforeX, beforeY := g.Player.X, g.Player.Y
				g.MovePlayer(dx, dy)
				if g.Player.X == beforeX && g.Player.Y == beforeY {
					// A closed door opens on the first attempt; step through on the second.
					g.MovePlayer(dx, dy)
				}
				if g.Player.X != p.x || g.Player.Y != p.y {
					t.Fatalf("seed %d level %d could not walk from %d,%d to %d,%d", seed, level, beforeX, beforeY, p.x, p.y)
				}
			}
			if g.Player.X != px || g.Player.Y != py {
				t.Fatalf("seed %d level %d did not reach the portal", seed, level)
			}
			g.UsePortal()
		}
	}
}

func TestWoodsRunEndToEndViaPortal(t *testing.T) {
	for seed := int64(0); seed < 8; seed++ {
		g := NewGame(seed + 19000)
		g.SetupNewGame(0, "tester")
		for level := 0; level < numWoodsMaps; level++ {
			if g.Level != level {
				t.Fatalf("seed %d reached level %d, want %d", seed, g.Level, level)
			}
			if g.WoodsRun.Themes[level] != g.currentWoodsTheme() {
				t.Fatalf("seed %d level %d theme mismatch", seed, level)
			}
			entry := g.WoodsRun.Entries[level]
			gaps := edgeGapTiles(g.Map, entry)
			if len(gaps) != 1 {
				t.Fatalf("seed %d level %d entry edge %d has %d gaps", seed, level, entry, len(gaps))
			}
			if manhattan(g.Player.X, g.Player.Y, gaps[0][0], gaps[0][1]) > 4 {
				t.Fatalf("seed %d level %d player did not start near the entry edge", seed, level)
			}

			px, py := findPortal(g)
			if px < 0 {
				t.Fatalf("seed %d level %d has no portal", seed, level)
			}
			g.Player.X, g.Player.Y = px, py
			g.UpdateFOV()
			g.UsePortal()
			if level < numWoodsMaps-1 {
				if g.Level != level+1 || g.Theme != 0 {
					t.Fatalf("seed %d portal from level %d reached level %d theme %d", seed, level, g.Level, g.Theme)
				}
			} else if g.Level != 9 || g.Theme != 1 {
				t.Fatalf("seed %d Meteor Field portal reached level %d theme %d, want caves", seed, g.Level, g.Theme)
			}
		}
	}
}
