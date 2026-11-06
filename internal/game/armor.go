package game

// ArmorStats is the player-facing AP model for a body armour item.
type ArmorStats struct {
	MaxAP  int
	Absorb int // percentage of incoming damage absorbed while AP remains
	DV     int // flat DV modifier (normally negative for heavy armour)
}

func ArmorStatsFor(d *ItemDef) (ArmorStats, bool) {
	if d == nil || d.Kind != KindArmor {
		return ArmorStats{}, false
	}
	switch d.Name {
	case "jacket":
		return ArmorStats{MaxAP: 20, Absorb: 50, DV: -1}, true
	case "flak armor":
		return ArmorStats{MaxAP: 35, Absorb: 60, DV: -2}, true
	case "power jacket":
		return ArmorStats{MaxAP: 50, Absorb: 70, DV: -3}, true
	case "leather jacket":
		return ArmorStats{MaxAP: 25, Absorb: 55, DV: -1}, true
	case "hide vest":
		return ArmorStats{MaxAP: 20, Absorb: 50, DV: 0}, true
	case "padded coat":
		return ArmorStats{MaxAP: 30, Absorb: 55, DV: -1}, true
	case "chain shirt":
		return ArmorStats{MaxAP: 45, Absorb: 65, DV: -2}, true
	case "ranger cloak":
		return ArmorStats{MaxAP: 15, Absorb: 40, DV: 0}, true
	default:
		return ArmorStats{MaxAP: 10, Absorb: 30, DV: 0}, true
	}
}

// ensureArmorAP initializes a new armour item's AP exactly once.  Depleted
// armour keeps AP 0.
func ensureArmorAP(it *Item) {
	if it == nil {
		return
	}
	stats, ok := ArmorStatsFor(it.Def())
	if !ok {
		return
	}
	if !it.APInit {
		it.AP = stats.MaxAP
		it.APInit = true
	}
	if it.AP > stats.MaxAP {
		it.AP = stats.MaxAP
	}
}

func armorAP(it *Item) (cur, max int) {
	if it == nil {
		return 0, 0
	}
	stats, ok := ArmorStatsFor(it.Def())
	if !ok {
		return 0, 0
	}
	ensureArmorAP(it)
	return it.AP, stats.MaxAP
}

// playerArmorStats returns the currently worn body armour stats.
func (g *Game) playerArmorStats() (ArmorStats, bool) {
	if g.Player == nil || g.Player.Slots.Body == nil {
		return ArmorStats{}, false
	}
	return ArmorStatsFor(g.Player.Slots.Body.Def())
}
