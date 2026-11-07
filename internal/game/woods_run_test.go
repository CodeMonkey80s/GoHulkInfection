package game

import "testing"

func TestWoodsRunThemeOrderAndDirections(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		g := NewGame(seed)
		run := g.WoodsRun
		if !validWoodsRun(&run) {
			t.Fatalf("seed %d produced an invalid run: %+v", seed, run)
		}
		if run.Themes[0] != WoodsChurch {
			t.Fatalf("seed %d first theme = %d, want Church", seed, run.Themes[0])
		}
		if run.Themes[4] != WoodsMeteorField {
			t.Fatalf("seed %d last theme = %d, want Meteor Field", seed, run.Themes[4])
		}
		seen := map[int]bool{}
		for _, theme := range run.Themes[1:4] {
			if theme < WoodsPeasantHouses || theme > WoodsCemetery || seen[theme] {
				t.Fatalf("seed %d middle themes are not a permutation: %v", seed, run.Themes)
			}
			seen[theme] = true
		}
		if run.Entries[0] != EdgeWest || run.Exits[0] != EdgeEast {
			t.Fatalf("seed %d Church is not west-to-east: %+v", seed, run)
		}
		if run.Entries[4] != EdgeWest || run.Exits[4] != EdgeEast {
			t.Fatalf("seed %d Meteor Field is not west-to-east: %+v", seed, run)
		}
		for i := 0; i < numWoodsMaps; i++ {
			if run.Entries[i] == run.Exits[i] {
				t.Fatalf("seed %d map %d immediately reverses: %+v", seed, i, run)
			}
			if i > 0 && run.Entries[i] != oppositeEdge(run.Exits[i-1]) {
				t.Fatalf("seed %d map %d entry is not opposite the previous exit: %+v", seed, i, run)
			}
		}
	}
}

func TestWoodsRunMiddleThemesAreShuffled(t *testing.T) {
	first := map[[3]int]bool{}
	for seed := int64(0); seed < 100; seed++ {
		g := NewGame(seed)
		key := [3]int{g.WoodsRun.Themes[1], g.WoodsRun.Themes[2], g.WoodsRun.Themes[3]}
		first[key] = true
	}
	if len(first) < 2 {
		t.Fatalf("middle theme order never varied: %v", first)
	}
}

func TestWoodsRunNoImmediateReversal(t *testing.T) {
	for seed := int64(0); seed < 200; seed++ {
		g := NewGame(seed)
		for i := 0; i < numWoodsMaps; i++ {
			if g.WoodsRun.Exits[i] == g.WoodsRun.Entries[i] {
				t.Fatalf("seed %d map %d exits through its entry edge", seed, i)
			}
			if i > 0 && g.WoodsRun.Entries[i] == g.WoodsRun.Exits[i-1] {
				t.Fatalf("seed %d map %d moves back through the previous exit", seed, i)
			}
		}
	}
}
