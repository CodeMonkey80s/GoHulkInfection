package game

import (
	"strings"
	"testing"
)

func TestOutOfThemeItemsAreRenamedOrUnavailable(t *testing.T) {
	for _, old := range []string{"he grenade", "incendiary grenade", "molotov cocktail", "flak armor", "power jacket", "plasteel stick", "power sword", "meta dagger", "vibra knife", "shock mace"} {
		if proto := itemProtoByName(old); proto >= 0 {
			t.Fatalf("old out-of-theme item %q still exists as proto %d", old, proto)
		}
	}
	for _, name := range []string{"holy water flask", "holy fire flask", "incense flask", "radiant flask", "padded jerkin", "silvered longsword", "consecrated dagger", "blessed knife", "censer mace", "healer's kit", "surgeon's kit"} {
		if proto := itemProtoByName(name); proto < 0 {
			t.Fatalf("themed item %q is missing", name)
		}
	}
}

func TestAvailableItemsFitTheme(t *testing.T) {
	banned := []string{"grenade", "gun", "pistol", "ammo", "clip", "rifle", "laser", "plasma", "flak", "power jacket", "plasteel", "meta ", "vibra", "shock ", "flashbang", "he grenade"}
	for proto, def := range ItemDefs {
		if !ItemAvailable(proto) {
			continue
		}
		name := strings.ToLower(def.Name)
		for _, word := range banned {
			if strings.Contains(name, word) {
				t.Fatalf("available item %d %q still uses out-of-theme term %q", proto, def.Name, word)
			}
		}
	}
}

func TestRangedWeaponsAmmoAndOldArmourAreUnavailable(t *testing.T) {
	for _, name := range []string{"lasgun", "personal pistol clip", "padded jerkin", "leather jacket"} {
		proto := itemProtoByName(name)
		if proto < 0 {
			continue
		}
		if ItemAvailable(proto) {
			t.Fatalf("%q should not be available in normal play", name)
		}
	}
	armour := itemProtoByName("padded jerkin")
	if armour < 0 || !CritterItemAvailable(armour) {
		t.Fatal("monsters should still be able to use armour prototypes")
	}
	if proto := itemProtoByName("holy fire flask"); proto < 0 || !ItemAvailable(proto) {
		t.Fatal("holy fire flask should be available")
	}
}
