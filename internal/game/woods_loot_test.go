package game

import (
	"path/filepath"
	"testing"
)

func itemInSignatureCategory(proto, category int) bool {
	if proto < 0 || proto >= len(ItemDefs) || category < 0 || category >= len(signatureCategoryItems) {
		return false
	}
	name := ItemDefs[proto].Name
	for _, want := range signatureCategoryItems[category] {
		if name == want {
			return true
		}
	}
	return false
}

func TestWoodLootAssignmentIsPermutation(t *testing.T) {
	g := NewGame(6001)
	g.SetupNewGame(0, "tester")
	if !validWoodLootAssignment(g.WoodLoot[:]) {
		t.Fatalf("invalid wood loot assignment: %v", g.WoodLoot)
	}
}

func TestSignatureCategoriesResolve(t *testing.T) {
	for category, names := range signatureCategoryItems {
		if len(names) == 0 {
			t.Fatalf("category %d is empty", category)
		}
		for _, name := range names {
			proto := itemProtoByName(name)
			if proto < 0 {
				t.Fatalf("category %d references unknown item %q", category, name)
			}
			if !ItemAvailable(proto) {
				t.Fatalf("category %d item %q is unavailable", category, name)
			}
		}
	}
}

func TestSignatureLootSpawnedAtLandmarks(t *testing.T) {
	for level := 0; level < numWoodsMaps; level++ {
		for seed := int64(0); seed < 6; seed++ {
			g := NewGame(seed + int64(level)*6000)
			g.SetupNewGame(0, "tester")
			g.GenerateLevel(level)
			category := g.WoodLoot[level]
			found, nearLandmark := false, false
			for _, it := range g.Items {
				if !itemInSignatureCategory(it.Proto, category) {
					continue
				}
				found = true
				if abs(it.X-g.Map.LandmarkX)+abs(it.Y-g.Map.LandmarkY) <= 12 {
					nearLandmark = true
				}
			}
			if !found {
				t.Fatalf("level %d seed %d has no signature category item", level, seed)
			}
			if !nearLandmark {
				t.Fatalf("level %d seed %d has no signature loot near its landmark", level, seed)
			}
		}
	}
}

func TestFoodAndNewArmour(t *testing.T) {
	g := NewGame(6002)
	g.SetupNewGame(0, "tester")
	apple := itemProtoByName("apple")
	toxic := itemProtoByName("toxic mushroom")
	if apple < 0 || toxic < 0 {
		t.Fatal("food prototypes missing")
	}
	g.Player.HP = 5
	appleItem := &Item{Proto: apple, Charges: 1}
	if !g.playerConsumeFood(appleItem) {
		t.Fatal("apple was not consumed")
	}
	if g.Player.HP <= 5 {
		t.Fatalf("apple did not heal: HP %d", g.Player.HP)
	}
	before := g.Player.HP
	toxicItem := &Item{Proto: toxic, Charges: 1}
	if !g.playerConsumeFood(toxicItem) {
		t.Fatal("toxic mushroom was not consumed")
	}
	if g.Player.HP >= before {
		t.Fatalf("toxic mushroom did not hurt: HP %d -> %d", before, g.Player.HP)
	}

	for _, name := range []string{"leather jacket", "hide vest", "padded coat", "chain shirt", "ranger cloak"} {
		proto := itemProtoByName(name)
		if proto < 0 {
			t.Fatalf("armour %q missing", name)
		}
		if stats, ok := ArmorStatsFor(&ItemDefs[proto]); !ok || stats.MaxAP <= 0 {
			t.Fatalf("armour %q has no stats", name)
		}
	}
}

func TestWoodLootSaveRoundtrip(t *testing.T) {
	g := NewGame(6003)
	g.SetupNewGame(0, "tester")
	path := filepath.Join(t.TempDir(), "loot.json")
	if err := g.SaveGame(path); err != nil {
		t.Fatal(err)
	}
	g2 := NewGame(6004)
	if err := g2.LoadGame(path); err != nil {
		t.Fatal(err)
	}
	if g2.WoodLoot != g.WoodLoot {
		t.Fatalf("wood loot assignment changed: %v vs %v", g2.WoodLoot, g.WoodLoot)
	}
}
