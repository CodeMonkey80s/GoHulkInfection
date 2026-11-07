package game

import (
	"fmt"
	"strings"
)

// Monster attacks are explicit dice pools. Damage is always at least 1.
func monsterDamageDice(name string) (count, sides, bonus int) {
	switch strings.ToLower(name) {
	case "zombie":
		return 1, 6, 1
	case "toxic zombie":
		return 1, 6, 2
	case "zombie punk":
		return 1, 6, 1
	case "zombie bomber":
		return 2, 4, 1
	case "mutant zombie":
		return 2, 8, 2
	case "biotech nightmare":
		return 2, 8, 3
	case "abomination":
		return 3, 8, 4
	case "observer mainbrain":
		return 1, 8, 2
	case "observer drone":
		return 1, 6, 1
	case "tech-spider mainbrain":
		return 2, 6, 2
	case "tech-spider drone":
		return 1, 8, 1
	case "warp wolf":
		return 2, 6, 2
	case "meta fly":
		return 1, 4, 0
	case "meta worm":
		return 1, 8, 1
	case "wolf":
		return 1, 6, 0
	case "deer":
		return 1, 4, 0
	case "villager":
		return 1, 4, 1
	case "militia":
		return 1, 6, 1
	case "rabbit":
		return 1, 2, 0
	case "fox":
		return 1, 4, 0
	case "boar":
		return 1, 6, 1
	case "crow":
		return 1, 2, 0
	case "frog":
		return 1, 2, 0
	default:
		return 1, 4, 0
	}
}

func init() {
	for i := range CritterDefs {
		CritterDefs[i].DmgCount, CritterDefs[i].DmgSides, CritterDefs[i].DmgBonus = monsterDamageDice(CritterDefs[i].Name)
	}
}

func assignMonsterDamage(c *Critter) {
	if c == nil {
		return
	}
	c.DmgCount, c.DmgSides, c.DmgBonus = monsterDamageDice(c.NameOrUnknown())
}

func (c *Critter) DamageDescription() string {
	if c == nil {
		return "1d4"
	}
	count, sides, bonus := c.DmgCount, c.DmgSides, c.DmgBonus
	if count < 1 {
		count = 1
	}
	if sides < 2 {
		sides = 4
	}
	switch {
	case bonus > 0:
		return fmt.Sprintf("%dd%d+%d", count, sides, bonus)
	case bonus < 0:
		return fmt.Sprintf("%dd%d%d", count, sides, bonus)
	default:
		return fmt.Sprintf("%dd%d", count, sides)
	}
}

func (g *Game) rollMonsterDamage(c *Critter) int {
	if c == nil {
		return 1
	}
	count, sides, bonus := c.DmgCount, c.DmgSides, c.DmgBonus
	if count < 1 {
		count = 1
	}
	if sides < 2 {
		sides = 4
	}
	total := bonus
	for i := 0; i < count; i++ {
		total += 1 + g.Roll(sides)
	}
	if total < 1 {
		return 1
	}
	return total
}
