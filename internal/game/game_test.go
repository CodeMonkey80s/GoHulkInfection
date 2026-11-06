package game

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupAndFOV(t *testing.T) {
	g := NewGame(1)
	g.SetupNewGame(0, "tester")
	if g.Map == nil || g.Map.W != mapWidth || g.Map.H != mapHeight {
		t.Fatalf("bad map size")
	}
	if !g.Map.IsPassable(g.Player.X, g.Player.Y) {
		t.Fatalf("player start is not passable")
	}
	g.UpdateFOV()
	vis := 0
	for _, v := range g.Map.Visible {
		if v {
			vis++
		}
	}
	if vis < 3 {
		t.Fatalf("fov too small: %d", vis)
	}
	if len(g.Critters) == 0 {
		t.Fatalf("no critters spawned")
	}
}

func TestCombatKill(t *testing.T) {
	g := NewGame(2)
	g.SetupNewGame(1, "tester")
	c := g.SpawnCritter(14, g.Player.X+1, g.Player.Y)
	if c == nil {
		t.Fatal("spawn failed")
	}
	c.HP = 1
	c.MaxHP = 1
	g.Options["optPlayerGodMode"] = 1
	oldXP := g.Player.XP
	g.Options["optMechanicBaseHitChance"] = 95
	for i := 0; i < 50 && c.IsAlive(); i++ {
		c.X, c.Y = g.Player.X+1, g.Player.Y
		c.HP = 1
		g.PlayerMelee(c)
	}
	if g.Player.XP <= oldXP {
		t.Fatalf("no xp gained: %d -> %d", oldXP, g.Player.XP)
	}
}

func TestSaveLoadRoundtrip(t *testing.T) {
	g := NewGame(3)
	g.SetupNewGame(2, "tester")
	dir := t.TempDir()
	path := filepath.Join(dir, "save.json")
	if err := g.SaveGame(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	g2 := NewGame(4)
	if err := g2.LoadGame(path); err != nil {
		t.Fatal(err)
	}
	if g2.Player.Name != g.Player.Name || g2.Level != g.Level {
		t.Fatalf("roundtrip mismatch")
	}
	if g2.Player.X != g.Player.X || g2.Player.Y != g.Player.Y {
		t.Fatalf("position mismatch: %d,%d vs %d,%d", g.Player.X, g.Player.Y, g2.Player.X, g2.Player.Y)
	}
	if g2.Map.W != g.Map.W || g2.Map.H != g.Map.H {
		t.Fatalf("map size mismatch: %dx%d vs %dx%d", g.Map.W, g.Map.H, g2.Map.W, g2.Map.H)
	}
	if len(g2.Map.Tiles) != len(g.Map.Tiles) {
		t.Fatalf("map tile count mismatch")
	}
}

func TestChangeLevelRelative(t *testing.T) {
	g := NewGame(9)
	g.SetupNewGame(0, "tester")
	g.ChangeLevelRelative(-1)
	if g.Level != 0 {
		t.Fatalf("previous from 0 moved to %d", g.Level)
	}
	g.ChangeLevelRelative(1)
	if g.Level != 1 {
		t.Fatalf("next from 0 moved to %d", g.Level)
	}
	g.ChangeLevelRelative(-1)
	if g.Level != 0 {
		t.Fatalf("previous from 1 moved to %d", g.Level)
	}
	for i := 0; i < 30; i++ {
		g.ChangeLevelRelative(1)
	}
	if g.Level != 27 {
		t.Fatalf("next should clamp at 27, got %d", g.Level)
	}
	g.ChangeLevel(9)
	if g.Theme != 1 {
		t.Fatalf("level 9 should be caves, theme %d", g.Theme)
	}
}

func TestKillCreatesCorpse(t *testing.T) {
	g := NewGame(11)
	g.SetupNewGame(0, "tester")
	if !g.Options.Bool("optMapCreateCorpses") {
		t.Fatal("corpse option should be enabled by default")
	}
	x, y := g.Player.X+2, g.Player.Y
	c := g.SpawnCritter(0, x, y)
	if c == nil {
		t.Fatal("spawn failed")
	}
	g.DamageCritter(c, 999, g.Player, 0)
	if g.CritterAt(x, y) != nil {
		t.Fatal("dead critter still occupies the tile")
	}
	found := false
	for _, it := range g.Items {
		if it.X == x && it.Y == y && it.Proto == 0 {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("corpse item was not left on the map")
	}
}

func TestPlayerIsDefencelessWithoutWeapon(t *testing.T) {
	g := NewGame(21)
	g.SetupNewGame(0, "tester")
	g.Player.Slots.Hands = nil
	c := g.SpawnCritter(0, g.Player.X+1, g.Player.Y)
	c.HP = 999
	c.MaxHP = 999
	before := c.HP
	g.PlayerDirectionalAttack(1, 0)
	if c.HP != before {
		t.Fatalf("unarmed player dealt damage: %d -> %d", before, c.HP)
	}
}

func TestThrowLimitAndAmmoFreeThrowables(t *testing.T) {
	g := NewGame(22)
	g.SetupNewGame(0, "tester")
	g.Options["optPlayerGodMode"] = 1
	g.Player.Slots.Hands = &Item{Proto: 41, X: -1, Y: -1}
	radius := g.FOVRadius()

	// A visible but out-of-range tile must be rejected and keep the grenade.
	farX := g.Player.X + radius + 2
	farY := g.Player.Y
	if g.Map.InBounds(farX, farY) {
		g.Map.Visible[g.Map.idx(farX, farY)] = true
		g.ThrowGrenade(farX, farY)
		if g.Player.Slots.Hands == nil {
			t.Fatal("out-of-range throw consumed the grenade")
		}
	}

	// A visible adjacent tile should accept the throw.
	nearX, nearY := g.Player.X+1, g.Player.Y
	if !g.Map.InBounds(nearX, nearY) {
		t.Fatal("bad test map")
	}
	g.Map.Visible[g.Map.idx(nearX, nearY)] = true
	g.ThrowGrenade(nearX, nearY)
	if g.Player.Slots.Hands != nil {
		t.Fatal("in-range throw did not consume the grenade")
	}
}

func TestSaveMigrationDropsRangedSlot(t *testing.T) {
	sc := saveCritter{Equip: []saveItem{
		{Proto: 16}, // melee
		{Proto: 22}, // old ranged
		{Proto: 5},  // armor
		{Proto: 41}, // throwable
		{Proto: 46}, // tool/medical slot in old layout
	}}
	c, removed := critterFromSave(sc, 1)
	if removed != 1 {
		t.Fatalf("expected 1 removed item, got %d", removed)
	}
	want := []int{16, 5, 41, 46}
	for i, proto := range want {
		if c.Equip[i] == nil || c.Equip[i].Proto != proto {
			t.Fatalf("equip slot %d = %v, want proto %d", i, c.Equip[i], proto)
		}
	}
}

func TestSpawnPoolsUseAvailableItemsOnly(t *testing.T) {
	for _, level := range []int{0, 5, 10, 20, 27} {
		for _, theme := range []int{0, 1, 2, 3} {
			for _, proto := range itemSpawnPool(theme, level) {
				if !ItemAvailable(proto) {
					t.Fatalf("spawn pool contains unavailable item %d", proto)
				}
			}
		}
	}
}

func TestActiveScope(t *testing.T) {
	if len(ActiveSkills) != 6 {
		t.Fatalf("expected 6 active skills, got %d", len(ActiveSkills))
	}
	if len(ActiveKeyBindings) == 0 {
		t.Fatal("no active key bindings")
	}
	for _, d := range ActiveOptionDefs() {
		if inactiveOptions[d.Name] {
			t.Fatalf("inactive option leaked into active list: %s", d.Name)
		}
	}
}

func TestPlayerSlotsStartAndBeltUnlock(t *testing.T) {
	g := NewGame(31)
	g.SetupNewGame(0, "tester")
	if g.Player.Slots.Hands == nil || g.Player.Slots.Body == nil {
		t.Fatal("starting hands/body missing")
	}
	if g.Player.Slots.BeltUnlocked != 2 {
		t.Fatalf("starting belt unlock = %d, want 2", g.Player.Slots.BeltUnlocked)
	}
	if g.Player.Slots.Belt[0] == nil || g.Player.Slots.Belt[1] == nil {
		t.Fatal("starting belt items missing")
	}
	if g.Player.Slots.Belt[2] != nil {
		t.Fatal("third belt slot should start empty and locked")
	}
	g.Player.Level = 2
	g.Player.XP = 400
	g.CheckLevelUp()
	if g.Player.Slots.BeltUnlocked != 2 {
		t.Fatalf("belt unlock before level 6 = %d, want 2", g.Player.Slots.BeltUnlocked)
	}
	g.Player.Level = 5
	g.Player.XP = 1000
	g.CheckLevelUp()
	if g.Player.Slots.BeltUnlocked != 3 {
		t.Fatalf("belt unlock after level 6 = %d", g.Player.Slots.BeltUnlocked)
	}
}

func TestPickupAutoPlaceAndOneItemPerTile(t *testing.T) {
	g := NewGame(32)
	g.SetupNewGame(0, "tester")
	g.State = StatePlaying
	g.Player.Slots.Hands = nil
	g.Player.Slots.Body = nil
	g.Player.Slots.Belt = [3]*Item{}
	g.Player.Slots.BeltUnlocked = 1

	it := &Item{Proto: 8, X: g.Player.X, Y: g.Player.Y}
	g.Items = append(g.Items, it)
	g.PickUp()
	if g.Player.Slots.Hands == nil || g.Player.Slots.Hands.Proto != 8 {
		t.Fatal("item did not auto-place into Hands")
	}
	if g.ItemAt(g.Player.X, g.Player.Y) != nil {
		t.Fatal("picked-up item remained on the ground")
	}

	// One item per tile: refusing to drop when occupied.
	dropped := &Item{Proto: 47, X: -1, Y: -1, Charges: 1}
	g.Player.Slots.Belt[0] = dropped
	occupied := &Item{Proto: 48, X: g.Player.X, Y: g.Player.Y, Charges: 5}
	g.Items = append(g.Items, occupied)
	g.DropSlot(SlotBelt1)
	if g.Player.Slots.Belt[0] != dropped {
		t.Fatal("drop should have been refused on an occupied tile")
	}
}

func TestPlayerSlotSaveRoundtrip(t *testing.T) {
	g := NewGame(33)
	g.SetupNewGame(1, "tester")
	g.Player.Slots.BeltUnlocked = 3
	g.Player.Slots.Belt[1] = &Item{Proto: 41, X: -1, Y: -1}
	g.Player.Slots.Belt[2] = &Item{Proto: 47, X: -1, Y: -1, Charges: 1}
	ensureArmorAP(g.Player.Slots.Body)
	g.Player.Slots.Body.AP = 7
	path := filepath.Join(t.TempDir(), "slots.json")
	if err := g.SaveGame(path); err != nil {
		t.Fatal(err)
	}
	g2 := NewGame(34)
	if err := g2.LoadGame(path); err != nil {
		t.Fatal(err)
	}
	if g2.Player.Slots.Hands == nil || g2.Player.Slots.Body == nil {
		t.Fatal("slots not restored")
	}
	if g2.Player.Slots.BeltUnlocked != 3 {
		t.Fatalf("belt unlock lost: %d", g2.Player.Slots.BeltUnlocked)
	}
	if g2.Player.Slots.Belt[1] == nil || g2.Player.Slots.Belt[1].Proto != 41 {
		t.Fatal("belt 2 item lost")
	}
	if g2.Player.Slots.Body == nil || g2.Player.Slots.Body.AP != 7 {
		t.Fatal("armour AP not restored")
	}
}

func TestPlayerSlotMigration(t *testing.T) {
	sc := saveCritter{Player: true, Equip: []saveItem{
		{Proto: 16}, {Proto: 5}, {Proto: 41}, {Proto: 46},
	}, Inv: []saveItem{{Proto: 47, Charges: 1}}}
	c, removed := critterFromSave(sc, 2)
	if removed != 0 {
		t.Fatalf("unexpected removals: %d", removed)
	}
	if c.Slots.Hands == nil || c.Slots.Hands.Proto != 16 {
		t.Fatal("hands migration failed")
	}
	if c.Slots.Body == nil || c.Slots.Body.Proto != 5 {
		t.Fatal("body migration failed")
	}
	ensureArmorAP(c.Slots.Body)
	if c.Slots.Body.AP != 20 {
		t.Fatalf("migrated body AP = %d, want 20", c.Slots.Body.AP)
	}
	want := []int{41, 46, 47}
	for i, proto := range want {
		if c.Slots.Belt[i] == nil || c.Slots.Belt[i].Proto != proto {
			t.Fatalf("belt %d = %v, want proto %d", i, c.Slots.Belt[i], proto)
		}
	}
}

func TestArmorAPAbsorption(t *testing.T) {
	g := NewGame(41)
	g.SetupNewGame(0, "tester")
	arm := g.Player.Slots.Body
	if arm == nil {
		t.Fatal("starting armour missing")
	}
	ensureArmorAP(arm)
	if arm.AP != 20 {
		t.Fatalf("starting AP = %d, want 20", arm.AP)
	}
	hp := g.Player.HP
	g.DamageCritter(g.Player, 20, nil, 0)
	if arm.AP != 10 {
		t.Fatalf("AP after 20 damage = %d, want 10", arm.AP)
	}
	if g.Player.HP != hp-10 {
		t.Fatalf("HP after absorbed damage = %d, want %d", g.Player.HP, hp-10)
	}

	// Acid bypasses AP completely.
	apBefore := arm.AP
	hp = g.Player.HP
	g.DamageCritter(g.Player, 10, nil, 2)
	if arm.AP != apBefore {
		t.Fatalf("acid changed AP: %d -> %d", apBefore, arm.AP)
	}
	if g.Player.HP != hp-10 {
		t.Fatalf("acid HP damage = %d, want 10", hp-g.Player.HP)
	}
}

func TestArmorVialRestoresAP(t *testing.T) {
	g := NewGame(42)
	g.SetupNewGame(0, "tester")
	g.Player.Slots.BeltUnlocked = 1
	g.Player.Slots.Belt[0] = &Item{Proto: ArmourVialProto, X: -1, Y: -1, Charges: 1}
	ensureArmorAP(g.Player.Slots.Body)
	g.Player.Slots.Body.AP = 5
	g.UseSlot(SlotBelt1)
	cur, _ := armorAP(g.Player.Slots.Body)
	if cur != 15 {
		t.Fatalf("AP after vial = %d, want 15", cur)
	}
	if g.Player.Slots.Belt[0] != nil {
		t.Fatal("empty vial remained in belt")
	}
}

func TestDirectionalAttack(t *testing.T) {
	g := NewGame(51)
	g.SetupNewGame(0, "tester")
	east := g.SpawnCritter(0, g.Player.X+1, g.Player.Y)
	west := g.SpawnCritter(0, g.Player.X-1, g.Player.Y)
	east.HP, east.MaxHP = 999, 999
	west.HP, west.MaxHP = 999, 999
	g.Options["optMechanicBaseHitChance"] = 95
	g.PlayerDirectionalAttack(1, 0)
	if east.HP == 999 {
		t.Fatal("east critter was not attacked")
	}
	if west.HP != 999 {
		t.Fatal("west critter was attacked instead")
	}
}
