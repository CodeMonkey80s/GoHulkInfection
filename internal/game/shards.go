package game

import "fmt"

func shardMaxAP(it *Item) int {
	if it == nil {
		return 0
	}
	d := it.Def()
	if d == nil || d.Kind != KindShard {
		return 0
	}
	if d.Charges < 1 {
		return 1
	}
	return d.Charges
}

func ensureShardAP(it *Item) {
	maxAP := shardMaxAP(it)
	if maxAP == 0 {
		return
	}
	if !it.APInit {
		it.AP = maxAP
		it.APInit = true
	}
	if it.AP > maxAP {
		it.AP = maxAP
	}
	if it.AP < 0 {
		it.AP = 0
	}
}

func (c *Critter) RobeSlot(i int) *Item {
	if c == nil || i < 0 || i >= c.Slots.RobeSlots || i >= MaxRobeSlots {
		return nil
	}
	return c.Slots.Robe[i]
}

func (c *Critter) SetRobeSlot(i int, it *Item) {
	if c == nil || i < 0 || i >= c.Slots.RobeSlots || i >= MaxRobeSlots {
		return
	}
	c.Slots.Robe[i] = it
}

func (c *Critter) FirstFreeRobeSlot() int {
	if c == nil {
		return -1
	}
	for i := 0; i < c.Slots.RobeSlots && i < MaxRobeSlots; i++ {
		if c.Slots.Robe[i] == nil {
			return i
		}
	}
	return -1
}

func (g *Game) PlayerAP() (cur, max int) {
	if g.Player == nil {
		return 0, 0
	}
	if g.Player.Slots.Body != nil {
		return armorAP(g.Player.Slots.Body)
	}
	for i := 0; i < g.Player.Slots.RobeSlots && i < MaxRobeSlots; i++ {
		shard := g.Player.Slots.Robe[i]
		if shard == nil || shard.Def() == nil || shard.Def().Kind != KindShard {
			continue
		}
		ensureShardAP(shard)
		cur += shard.AP
		max += shardMaxAP(shard)
	}
	return cur, max
}

// damagePlayerShards drains equipped shards in slot order and returns the
// damage left for HP. Empty shards break and free their slot immediately.
func (g *Game) damagePlayerShards(dmg int) int {
	if g.Player == nil || dmg <= 0 {
		return dmg
	}
	for i := 0; i < g.Player.Slots.RobeSlots && i < MaxRobeSlots && dmg > 0; i++ {
		shard := g.Player.Slots.Robe[i]
		if shard == nil || shard.Def() == nil || shard.Def().Kind != KindShard {
			continue
		}
		ensureShardAP(shard)
		if shard.AP <= 0 {
			g.Player.Slots.Robe[i] = nil
			continue
		}
		absorbed := dmg
		if absorbed > shard.AP {
			absorbed = shard.AP
		}
		shard.AP -= absorbed
		dmg -= absorbed
		if shard.AP <= 0 {
			g.Msg("your "+shard.Name()+" breaks and falls away.", 12)
			g.Player.Slots.Robe[i] = nil
		}
	}
	return dmg
}

func (g *Game) blessShards(amount int) bool {
	if g.Player == nil || amount <= 0 {
		return false
	}
	bestIndex, bestMissing := -1, 0
	for i := 0; i < g.Player.Slots.RobeSlots && i < MaxRobeSlots; i++ {
		shard := g.Player.Slots.Robe[i]
		if shard == nil || shard.Def() == nil || shard.Def().Kind != KindShard {
			continue
		}
		ensureShardAP(shard)
		missing := shardMaxAP(shard) - shard.AP
		if missing > bestMissing {
			bestIndex, bestMissing = i, missing
		}
	}
	if bestIndex < 0 {
		return false
	}
	shard := g.Player.Slots.Robe[bestIndex]
	shard.AP += amount
	if shard.AP > shardMaxAP(shard) {
		shard.AP = shardMaxAP(shard)
	}
	g.Msg(fmt.Sprintf("your %s is blessed (%d/%d AP).", shard.Name(), shard.AP, shardMaxAP(shard)), 10)
	return true
}
