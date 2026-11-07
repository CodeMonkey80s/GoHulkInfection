package game

import (
	"math/rand"
)

// Signature loot categories. At world creation these are shuffled and assigned
// one per woods level; the level landmark guarantees an item from its category.
var signatureCategoryItems = [][]string{
	{"saw", "shovel", "rope", "lantern", "torch", "lockpick", "fishing rod", "whetstone", "hatchet"},
	{"apple", "pear", "berries", "honey", "bread", "cheese", "dried meat", "mushroom", "cooked fish", "stew"},
	{"leather shard", "chain shard", "iron shard", "silver shard", "blessed shard", "bone charm"},
	{"fishing rod", "raw fish", "cooked fish", "rope", "whetstone"},
	{"hunting knife", "boar spear", "spear", "snare", "hatchet", "rope"},
	{"healing herb", "red potion", "green potion", "blue potion", "antidote", "holy water"},
	{"club", "spear", "sickle", "scythe", "pitchfork", "hatchet", "wood axe"},
	{"torch", "lantern", "flint", "oil flask", "holy fire flask", "incense flask"},
	{"silver chalice", "old tome", "bone charm", "gold cross", "holy water"},
}

func newWoodLootAssignment(rng *rand.Rand) [9]int {
	var cats [9]int
	for i := range cats {
		cats[i] = i
	}
	rng.Shuffle(len(cats), func(i, j int) {
		cats[i], cats[j] = cats[j], cats[i]
	})
	return cats
}

func validWoodLootAssignment(v []int) bool {
	if len(v) != 9 {
		return false
	}
	seen := [9]bool{}
	for _, cat := range v {
		if cat < 0 || cat >= len(seen) || seen[cat] {
			return false
		}
		seen[cat] = true
	}
	return true
}

func itemProtoByName(name string) int {
	for i := range ItemDefs {
		if ItemDefs[i].Name == name {
			return i
		}
	}
	return -1
}

func (g *Game) randomSignatureItem(category int) int {
	if category < 0 || category >= len(signatureCategoryItems) {
		return -1
	}
	names := signatureCategoryItems[category]
	start := g.Roll(len(names))
	for i := 0; i < len(names); i++ {
		proto := itemProtoByName(names[(start+i)%len(names)])
		if proto >= 0 && ItemAvailable(proto) {
			return proto
		}
	}
	return -1
}

func (g *Game) spawnSignatureLoot(level int) {
	if g.Map == nil || level < 0 || level >= len(g.WoodLoot) {
		return
	}
	category := g.WoodLoot[level]
	if proto := g.randomSignatureItem(category); proto >= 0 {
		it := &Item{Proto: proto, X: -1, Y: -1}
		if d := it.Def(); d != nil && d.Kind == KindMedical {
			it.Charges = d.Charges
		}
		g.placeLootNear(g.Map.LandmarkX, g.Map.LandmarkY, it)
	}
	// One extra signature item keeps the compact maps thematic without
	// flooding them with loot.
	for i := 0; i < 1; i++ {
		proto := g.randomSignatureItem(category)
		if proto < 0 {
			continue
		}
		x, y := g.Map.RandomFree(g.RNG)
		it := &Item{Proto: proto, X: -1, Y: -1}
		if d := it.Def(); d != nil && d.Kind == KindMedical {
			it.Charges = d.Charges
		}
		g.placeLootNear(x, y, it)
	}
}

func (g *Game) placeLootNear(x, y int, it *Item) bool {
	if it == nil || g.Map == nil {
		return false
	}
	for radius := 0; radius <= 10; radius++ {
		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				if radius > 0 && abs(dx) != radius && abs(dy) != radius {
					continue
				}
				nx, ny := x+dx, y+dy
				if !g.Map.InBounds(nx, ny) || g.Map.BlocksMove(nx, ny) {
					continue
				}
				t := g.Map.Tile(nx, ny)
				if isWoodsWater(t) || t == tileRiverBank || t == portalTile(0) {
					continue
				}
				if g.PlaceGroundItem(it, nx, ny) {
					return true
				}
			}
		}
	}
	return false
}

func init() {
	addWeapon := func(w WeaponDef) int {
		id := len(WeaponDefs)
		WeaponDefs = append(WeaponDefs, w)
		return id
	}
	addItem := func(d ItemDef) int {
		id := len(ItemDefs)
		ItemDefs = append(ItemDefs, d)
		return id
	}
	food := func(name, desc string, heal int, glyph byte, color uint8) {
		addItem(ItemDef{
			Name: name, Desc: desc, Kind: KindFood, Weight: 100,
			Charges: 1, Heal: heal, Glyph: glyph, Color: color,
		})
	}
	tool := func(name, desc string, glyph byte, color uint8) {
		addItem(ItemDef{
			Name: name, Desc: desc, Kind: KindTool, Weight: 800,
			Glyph: glyph, Color: color,
		})
	}
	weapon := func(name, desc string, sub int, glyph byte, color uint8) {
		addItem(ItemDef{
			Name: name, Desc: desc, Kind: KindMelee, Weight: 2000, Damage: 47,
			WeaponSub: sub, Burst: 1, Reload: 1, Range: 1, Glyph: glyph, Color: color,
		})
	}
	armor := func(name, desc string, glyph byte, color uint8) {
		addItem(ItemDef{
			Name: name, Desc: desc, Kind: KindArmor, Weight: 3000,
			ArmorMelee: 80, ArmorRanged: 80, ArmorAcid: 80, Glyph: glyph, Color: color,
		})
	}
	medical := func(name, desc string, heal int, glyph byte, color uint8) {
		addItem(ItemDef{
			Name: name, Desc: desc, Kind: KindMedical, Weight: 300,
			Charges: 1, Heal: heal, Glyph: glyph, Color: color,
		})
	}
	relic := func(name, desc string, glyph byte, color uint8) {
		addItem(ItemDef{
			Name: name, Desc: desc, Kind: KindRelic, Weight: 500,
			Glyph: glyph, Color: color,
		})
	}

	hatchetSub := addWeapon(WeaponDef{AV: 3, DV: 0, DmgMin: 2, DmgMax: 6})
	woodAxeSub := addWeapon(WeaponDef{AV: 4, DV: -1, DmgMin: 3, DmgMax: 9})
	pitchforkSub := addWeapon(WeaponDef{AV: 3, DV: 1, DmgMin: 2, DmgMax: 7})
	sickleSub := addWeapon(WeaponDef{AV: 4, DV: 1, DmgMin: 2, DmgMax: 6})
	scytheSub := addWeapon(WeaponDef{AV: 5, DV: -2, DmgMin: 3, DmgMax: 10})
	clubSub := addWeapon(WeaponDef{AV: 2, DV: 1, DmgMin: 2, DmgMax: 5})
	spearSub := addWeapon(WeaponDef{AV: 4, DV: 2, DmgMin: 2, DmgMax: 8})
	huntingKnifeSub := addWeapon(WeaponDef{AV: 5, DV: 3, DmgMin: 1, DmgMax: 5})
	boarSpearSub := addWeapon(WeaponDef{AV: 5, DV: 1, DmgMin: 3, DmgMax: 9})

	// Food.
	food("apple", "a bruised apple from an old orchard.", 4, '%', 4)
	food("pear", "a hard pear with a sweet core.", 4, '%', 14)
	food("berries", "a handful of dark woodland berries.", 3, '%', 13)
	food("honey", "a small jar of wild honey.", 6, '!', 14)
	food("bread", "a stale loaf of bread.", 6, '%', 6)
	food("cheese", "a wedge of hard cheese.", 5, '%', 14)
	food("dried meat", "a strip of smoked meat.", 8, '%', 4)
	food("mushroom", "a plain woodland mushroom.", 3, '%', 8)
	food("toxic mushroom", "a sickly mushroom with a bitter smell.", -5, '%', 10)
	food("raw fish", "a slippery raw fish.", 3, '%', 3)
	food("cooked fish", "a fish cooked over a fire.", 7, '%', 3)
	food("stew", "a bowl of thick woodland stew.", 10, '~', 6)

	// Tools.
	tool("saw", "a rusted saw for cutting wood.", '/', 7)
	tool("shovel", "a short-handled shovel.", '\\', 6)
	tool("rope", "a coil of sturdy rope.", '~', 6)
	tool("lantern", "a hooded lantern with a little oil left.", 'i', 14)
	tool("torch", "a pitch-soaked torch.", '|', 12)
	tool("lockpick", "a bent lockpick.", '-', 7)
	tool("fishing rod", "a simple fishing rod.", '/', 6)
	tool("whetstone", "a worn whetstone.", 'o', 8)
	tool("snare", "a wire snare for small game.", '&', 6)
	tool("flint", "a piece of flint and steel.", '*', 8)
	tool("oil flask", "a flask of lamp oil.", '!', 6)

	// Melee weapons.
	weapon("hatchet", "a woodcutter's hatchet.", hatchetSub, '/', 7)
	weapon("wood axe", "a heavy wood axe.", woodAxeSub, 'P', 7)
	weapon("pitchfork", "a farmer's pitchfork.", pitchforkSub, '/', 6)
	weapon("sickle", "a curved harvesting sickle.", sickleSub, ')', 7)
	weapon("scythe", "a long two-handed scythe.", scytheSub, '/', 7)
	weapon("club", "a heavy wooden club.", clubSub, '|', 6)
	weapon("spear", "a simple wooden spear.", spearSub, '/', 6)
	weapon("hunting knife", "a sharp hunting knife.", huntingKnifeSub, '/', 7)
	weapon("boar spear", "a broad-bladed boar spear.", boarSpearSub, '/', 6)

	// Armour.
	armor("leather jacket", "a supple leather jacket.", '[', 6)
	armor("hide vest", "a vest of thick animal hide.", '[', 6)
	armor("padded coat", "a padded cloth coat.", '[', 6)
	armor("chain shirt", "a heavy shirt of rusted chain.", '[', 7)
	armor("ranger cloak", "a weather-stained ranger cloak.", '[', 2)

	// Herbs and potions.
	medical("healing herb", "a bitter herb that numbs pain.", 5, '*', 10)
	medical("red potion", "a red potion that closes wounds.", 15, '!', 4)
	medical("green potion", "a green potion that steadies the body.", 8, '!', 10)
	medical("blue potion", "a blue potion that clears the mind.", 12, '!', 1)
	medical("antidote", "a sharp-smelling antidote.", 4, '!', 11)
	medical("holy water", "a small flask of blessed water.", 8, '!', 15)

	// Chapel relics.
	relic("silver chalice", "a tarnished silver chalice.", '$', 7)
	relic("old tome", "a water-stained tome of old prayers.", '?', 6)
	relic("bone charm", "a charm carved from yellowed bone.", '!', 7)
	relic("gold cross", "a small gold cross on a broken chain.", '+', 14)
}
