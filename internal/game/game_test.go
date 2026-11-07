package game

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupAndFOV(t *testing.T) {
	g := NewGame(1)
	g.SetupNewGame(0, "tester")
	if g.Map == nil {
		t.Fatal("no map generated")
	}
	if g.Map.W != mapWidth || g.Map.H != woodsMapHeight {
		t.Fatalf("bad map size: %dx%d", g.Map.W, g.Map.H)
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

func TestDeathDeletesSave(t *testing.T) {
	g := NewGame(35)
	g.SetupNewGame(0, "tester")
	path := filepath.Join(t.TempDir(), "death.json")
	g.SavePath = path
	if err := g.SaveGame(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("save was not created")
	}
	g.DamageCritter(g.Player, 9999, nil, 0)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("save still exists after death: %v", err)
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
	g.Map.SetTile(x, y, 2) // keep the corpse test on ordinary grass
	for _, it := range append([]*Item(nil), g.Items...) {
		if it.X == x && it.Y == y {
			g.RemoveItem(it)
		}
	}
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
	if removed != 2 {
		t.Fatalf("expected 2 removed items, got %d", removed)
	}
	if c.Equip[0] == nil || c.Equip[0].Proto != 16 {
		t.Fatalf("equip slot 0 = %v, want proto 16", c.Equip[0])
	}
	if c.Equip[1] != nil {
		t.Fatal("old body armour should not be migrated")
	}
	if c.Equip[2] == nil || c.Equip[2].Proto != 41 {
		t.Fatalf("equip slot 2 = %v, want proto 41", c.Equip[2])
	}
	if c.Equip[3] == nil || c.Equip[3].Proto != 46 {
		t.Fatalf("equip slot 3 = %v, want proto 46", c.Equip[3])
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

func TestPriestStartsWithRobeShardVialAndRite(t *testing.T) {
	g := NewGame(31)
	g.SetupNewGame(0, "tester")
	if g.Player.Slots.Hands == nil || g.Player.Slots.Hands.Name() != "sickle" {
		t.Fatal("starting weapon missing")
	}
	if g.Player.Slots.RobeSlots != 3 {
		t.Fatalf("robe slots = %d, want 3", g.Player.Slots.RobeSlots)
	}
	if g.Player.Slots.Robe[0] == nil || g.Player.Slots.Robe[0].Def().Kind != KindShard {
		t.Fatal("starting shard missing")
	}
	if g.Player.Slots.Robe[1] == nil || g.Player.Slots.Robe[1].Def().Kind != KindMedical {
		t.Fatal("starting vial missing")
	}
	if g.Player.Slots.Robe[2] != nil {
		t.Fatal("third robe slot should start empty")
	}
	if len(g.KnownRites()) != 1 {
		t.Fatalf("known rites = %d, want 1", len(g.KnownRites()))
	}
	if g.RiteCharges[0] != 1 {
		t.Fatalf("starting rite charge = %d, want 1", g.RiteCharges[0])
	}
	if g.Player.Level != 1 {
		t.Fatalf("player level = %d, want 1", g.Player.Level)
	}
}

func TestPickupAutoPlaceAndOneItemPerTile(t *testing.T) {
	g := NewGame(32)
	g.SetupNewGame(0, "tester")
	g.State = StatePlaying
	g.Player.Slots.Hands = nil
	g.Player.Slots.Robe = [MaxRobeSlots]*Item{}

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
	g.Player.Slots.Robe[0] = dropped
	occupied := &Item{Proto: 48, X: g.Player.X, Y: g.Player.Y, Charges: 5}
	g.Items = append(g.Items, occupied)
	g.DropSlot(SlotRobe1)
	if g.Player.Slots.Robe[0] != dropped {
		t.Fatal("drop should have been refused on an occupied tile")
	}
}

func TestPlayerSlotSaveRoundtrip(t *testing.T) {
	g := NewGame(33)
	g.SetupNewGame(1, "tester")
	g.Player.Slots.Head = &Item{Proto: itemProtoByName("plain hat"), X: -1, Y: -1}
	g.Player.Slots.Robe[0] = &Item{Proto: itemProtoByName("chain shard"), X: -1, Y: -1}
	ensureShardAP(g.Player.Slots.Robe[0])
	g.Player.Slots.Robe[1] = &Item{Proto: itemProtoByName("healing vial"), X: -1, Y: -1, Charges: 1}
	g.Player.Slots.Robe[2] = &Item{Proto: 47, X: -1, Y: -1, Charges: 1}
	g.Player.Slots.RobeSlots = 3
	g.PlayerRites[1] = RiteSmite
	g.RiteCharges[1] = 1
	path := filepath.Join(t.TempDir(), "slots.json")
	if err := g.SaveGame(path); err != nil {
		t.Fatal(err)
	}
	g2 := NewGame(34)
	if err := g2.LoadGame(path); err != nil {
		t.Fatal(err)
	}
	if g2.Player.Slots.Hands == nil || g2.Player.Slots.Head == nil {
		t.Fatal("slots not restored")
	}
	if g2.Player.Slots.RobeSlots != 3 {
		t.Fatalf("robe slots lost: %d", g2.Player.Slots.RobeSlots)
	}
	if g2.Player.Slots.Robe[0] == nil || g2.Player.Slots.Robe[0].Name() != "chain shard" {
		t.Fatal("shard not restored")
	}
	if g2.Player.Slots.Robe[0].AP != 5 {
		t.Fatalf("shard AP not restored: %d", g2.Player.Slots.Robe[0].AP)
	}
	if g2.Player.Slots.Head.Name() != "plain hat" {
		t.Fatal("head item not restored")
	}
	if g2.PlayerRites[1] != RiteSmite || g2.RiteCharges[1] != 1 {
		t.Fatal("second rite not restored")
	}
}

func TestPlayerSlotMigrationDropsOldBodyArmour(t *testing.T) {
	sc := saveCritter{Player: true, Equip: []saveItem{
		{Proto: 16}, {Proto: 5}, {Proto: 41}, {Proto: 46},
	}, Inv: []saveItem{{Proto: 47, Charges: 1}}}
	c, _ := critterFromSave(sc, 2)
	if c.Slots.Hands == nil || c.Slots.Hands.Proto != 16 {
		t.Fatal("hands migration failed")
	}
	if c.Slots.Body != nil {
		t.Fatal("old body armour should not be migrated")
	}
	want := []int{41, 46, 47}
	for i, proto := range want {
		if c.Slots.Belt[i] == nil || c.Slots.Belt[i].Proto != proto {
			t.Fatalf("belt %d = %v, want proto %d", i, c.Slots.Belt[i], proto)
		}
	}
}

func TestShardsDrainSequentiallyAndBreak(t *testing.T) {
	g := NewGame(41)
	g.SetupNewGame(0, "tester")
	g.Player.Slots.Robe[0] = &Item{Proto: itemProtoByName("leather shard"), X: -1, Y: -1}
	g.Player.Slots.Robe[1] = &Item{Proto: itemProtoByName("chain shard"), X: -1, Y: -1}
	ensureShardAP(g.Player.Slots.Robe[0])
	ensureShardAP(g.Player.Slots.Robe[1])
	hp := g.Player.HP

	g.DamageCritter(g.Player, 6, nil, 0)
	if g.Player.Slots.Robe[0] != nil {
		t.Fatal("first shard should have broken")
	}
	if g.Player.Slots.Robe[1] == nil || g.Player.Slots.Robe[1].AP != 2 {
		t.Fatal("second shard did not absorb the remaining damage")
	}
	if g.Player.HP != hp {
		t.Fatalf("HP changed while shards remained: %d -> %d", hp, g.Player.HP)
	}

	g.DamageCritter(g.Player, 3, nil, 0)
	if g.Player.Slots.Robe[1] != nil {
		t.Fatal("second shard should have broken")
	}
	if g.Player.HP != hp-1 {
		t.Fatalf("excess damage did not reach HP: %d -> %d", hp, g.Player.HP)
	}
}

func TestHealingVialRestoresHP(t *testing.T) {
	g := NewGame(42)
	g.SetupNewGame(0, "tester")
	g.Player.HP = 10
	vial := g.Player.Slots.Robe[1]
	g.UseSlot(SlotRobe2)
	if g.Player.HP != 15 {
		t.Fatalf("vial healed to %d, want 15", g.Player.HP)
	}
	if g.Player.Slots.Robe[1] != nil || vial == nil {
		t.Fatal("empty vial was not removed from the robe")
	}
}

func TestDirectionalAttack(t *testing.T) {
	g := NewGame(51)
	g.SetupNewGame(0, "tester")
	east := g.SpawnCritter(0, g.Player.X+1, g.Player.Y)
	west := g.SpawnCritter(0, g.Player.X-1, g.Player.Y)
	east.HP, east.MaxHP = 999, 999
	west.HP, west.MaxHP = 999, 999
	g.Options["optMechanicBaseHitChance"] = 100
	g.Options["optPlayerGodMode"] = 1
	for i := 0; i < 50 && east.HP == 999; i++ {
		east.X, east.Y = g.Player.X+1, g.Player.Y
		west.X, west.Y = g.Player.X-1, g.Player.Y
		g.PlayerDirectionalAttack(1, 0)
	}
	if east.HP == 999 {
		t.Fatal("east critter was not attacked")
	}
	if west.HP != 999 {
		t.Fatal("west critter was attacked instead")
	}
}
