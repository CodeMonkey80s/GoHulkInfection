package game

import (
	"fmt"
	"strings"
)

func (g *Game) UseRite(index int) {
	if g.Player == nil || g.State != StatePlaying {
		return
	}
	if index < 0 || index >= MaxRites {
		return
	}
	rite := g.PlayerRites[index]
	if rite < 0 || rite >= len(RiteDefs) {
		g.Msg("you have no rite to use.", 8)
		return
	}
	if g.RiteCharges[index] <= 0 {
		g.Msg(fmt.Sprintf("the rite of %s has no charges left on this map.", riteName(rite)), 8)
		return
	}
	used := false
	switch rite {
	case RiteBlessing:
		used = g.blessShards(2)
		if !used {
			g.Msg("your shards are already fully restored.", 8)
		}
	case RiteHealingWounds:
		if g.Player.HP >= g.Player.MaxHP {
			g.Msg("you are already at full health.", 8)
			return
		}
		g.Player.HP += 5
		if g.Player.HP > g.Player.MaxHP {
			g.Player.HP = g.Player.MaxHP
		}
		g.Msg("your wounds close with a prayer.", 10)
		used = true
	case RiteShout:
		targets := 0
		for _, c := range g.Critters {
			if c == nil || !c.IsAlive() || c.Align != AlignEnemy {
				continue
			}
			if _, ok := g.canSeeCritter(g.Player, c, 6); !ok {
				continue
			}
			c.Flee = 3
			c.Path = nil
			targets++
		}
		if targets > 0 {
			g.Msg("your shout drives the dead back.", 14)
			used = true
		} else {
			g.Msg("there is nothing nearby to hear your shout.", 8)
		}
	case RiteBanish:
		target := g.nearestRiteTarget(8)
		if target == nil {
			g.Msg("there is no enemy in sight to banish.", 8)
			return
		}
		if !isUndeadCritter(target) {
			g.Msg("the rite only affects the undead.", 8)
			return
		}
		if target.MaxHP <= 18 {
			g.Msg("the "+target.NameOrUnknown()+" is banished from this place.", 14)
			target.Dead = true
			target.HP = 0
			g.RemoveCritter(target)
		} else {
			target.Flee = 5
			target.Path = nil
			g.Msg("the "+target.NameOrUnknown()+" recoils from your prayer.", 14)
		}
		used = true
	case RiteImmobilize:
		target := g.nearestRiteTarget(8)
		if target == nil {
			g.Msg("there is no enemy in sight to hold.", 8)
			return
		}
		target.Immobile = 3
		target.Path = nil
		g.Msg("the "+target.NameOrUnknown()+" is held in place.", 14)
		used = true
	case RiteSmite:
		target := g.nearestRiteTarget(8)
		if target == nil {
			g.Msg("there is no enemy in sight to smite.", 8)
			return
		}
		dmg := 6
		if isUndeadCritter(target) {
			dmg += 2
		}
		g.Msg("you call down a smiting light.", 14)
		g.DamageCritter(target, dmg, g.Player, 0)
		used = true
	case RiteSanctuary:
		g.PeaceTurns = 3
		for _, c := range g.Critters {
			if c != nil && c.Align == AlignEnemy {
				c.Alert = 0
				c.Path = nil
			}
		}
		g.Msg("a quiet sanctuary settles over the road.", 10)
		used = true
	}
	if used {
		g.RiteCharges[index]--
		g.EndPlayerTurn()
	}
}

func (g *Game) nearestRiteTarget(maxDist int) *Critter {
	var best *Critter
	bestDist := maxDist + 1
	for _, c := range g.Critters {
		if c == nil || !c.IsAlive() || c.Align != AlignEnemy {
			continue
		}
		dist, ok := g.canSeeCritter(g.Player, c, maxDist)
		if !ok || dist >= bestDist {
			continue
		}
		best, bestDist = c, dist
	}
	return best
}

func isUndeadCritter(c *Critter) bool {
	if c == nil {
		return false
	}
	name := strings.ToLower(c.NameOrUnknown())
	return strings.Contains(name, "zombie") ||
		strings.Contains(name, "ghoul") ||
		strings.Contains(name, "wraith") ||
		strings.Contains(name, "undead")
}
