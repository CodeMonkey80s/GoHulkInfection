package game

import (
	"fmt"
	"strings"
)

const RiteNone = -1

const (
	RiteBlessing = iota
	RiteHealingWounds
	RiteShout
	RiteBanish
	RiteImmobilize
	RiteSmite
	RiteSanctuary
)

type RiteDef struct {
	Name string
	Desc string
}

var RiteDefs = []RiteDef{
	{Name: "Blessing", Desc: "restore 2 AP to the most damaged shard."},
	{Name: "Healing Wounds", Desc: "restore 5 HP."},
	{Name: "Shout", Desc: "drive nearby enemies away for 3 turns."},
	{Name: "Banish", Desc: "repel undead; weak undead are removed."},
	{Name: "Immobilize", Desc: "hold one enemy in place for 3 turns."},
	{Name: "Smite", Desc: "strike one visible enemy with holy damage."},
	{Name: "Sanctuary", Desc: "enemies lose interest for 3 turns."},
}

var PriestNames = []string{
	"Father Aleksy", "Father Barnaba", "Father Cyryl", "Father Damian",
	"Father Eliasz", "Father Feliks", "Father Grzegorz", "Father Hieronim",
	"Father Ignacy", "Father Jakub", "Father Kazimierz", "Father Leon",
	"Father Mateusz", "Father Mikołaj", "Father Piotr", "Father Rafał",
	"Father Szymon", "Father Tomasz", "Father Urban", "Father Wiktor",
}

func (g *Game) RollPriestProfile() {
	g.PlayerName = PriestNames[g.Roll(len(PriestNames))]
	g.PlayerRites = [MaxRites]int{RiteNone, RiteNone}
	g.RiteCharges = [MaxRites]int{}
	g.PlayerRites[0] = g.Roll(len(RiteDefs))
	g.RiteCharges[0] = 1
}

func (g *Game) KnownRites() []int {
	var out []int
	for _, rite := range g.PlayerRites {
		if rite >= 0 && rite < len(RiteDefs) {
			out = append(out, rite)
		}
	}
	return out
}

func (g *Game) RefreshRiteCharges() {
	for i, rite := range g.PlayerRites {
		if rite >= 0 && rite < len(RiteDefs) {
			g.RiteCharges[i] = 1
		}
	}
}

func riteName(rite int) string {
	if rite < 0 || rite >= len(RiteDefs) {
		return "none"
	}
	return RiteDefs[rite].Name
}

func (g *Game) RiteSummary() string {
	parts := make([]string, 0, MaxRites)
	for i, rite := range g.PlayerRites {
		if rite < 0 || rite >= len(RiteDefs) {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s %d", riteName(rite), g.RiteCharges[i]))
	}
	if len(parts) == 0 {
		return "Rite: none"
	}
	return "Rite: " + strings.Join(parts, "  ")
}

func (g *Game) priestStartingItem(name string) *Item {
	proto := itemProtoByName(name)
	if proto < 0 {
		return nil
	}
	it := startingItem(proto)
	if it != nil && it.Def() != nil && it.Def().Kind == KindShard {
		ensureShardAP(it)
	}
	return it
}

func init() {
	addItem := func(d ItemDef) {
		ItemDefs = append(ItemDefs, d)
	}
	addItem(ItemDef{
		Name: "healing vial", Desc: "a small vial of bitter healing tonic.",
		Kind: KindMedical, Weight: 100, Charges: 1, Heal: 5, Glyph: '!', Color: 10,
	})
	shard := func(name, desc string, maxAP int, glyph byte, color uint8) {
		addItem(ItemDef{
			Name: name, Desc: desc, Kind: KindShard, Weight: 200,
			Charges: maxAP, Glyph: glyph, Color: color,
		})
	}
	shard("leather shard", "a hardened leather plate for the robe.", 3, '=', 6)
	shard("chain shard", "a small linked-mail plate for the robe.", 5, '=', 7)
	shard("iron shard", "a heavy iron plate for the robe.", 7, '=', 8)
	shard("silver shard", "a silver plate that turns aside claws.", 9, '=', 15)
	shard("blessed shard", "a blessed plate warm to the touch.", 12, '=', 14)
	addItem(ItemDef{
		Name: "plain hat", Desc: "a simple weather-stained hat.", Kind: KindHead,
		Weight: 100, Glyph: '^', Color: 6,
	})
}
