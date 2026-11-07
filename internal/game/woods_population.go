package game

// spawnWoodsPopulation keeps the compact maps at roughly 4-10 critters with a
// theme-appropriate mix of villagers, animals and undead.
func (g *Game) spawnWoodsPopulation(theme int) {
	switch theme {
	case WoodsChurch:
		for i := 0; i < 1+g.Roll(2); i++ {
			g.spawnWoodsCritterNearLandmark(0, 10, 7)
		}
		for i := 0; i < 1+g.Roll(2); i++ {
			g.spawnWoodsCritter(g.rollProto([]int{0, 2}), 7)
		}
		g.spawnWoodsCritter(crowProto, 6)
	case WoodsPeasantHouses:
		for i := 0; i < 2+g.Roll(3); i++ {
			if g.Roll(3) == 0 {
				g.spawnWoodsCritter(militiaProto, 6)
			} else {
				g.spawnWoodsCritter(villagerProto, 6)
			}
		}
		for i := 0; i < 1+g.Roll(3); i++ {
			g.spawnWoodsCritter(g.rollProto(neutralWoodsAnimals), 6)
		}
		for i := 0; i < 1+g.Roll(3); i++ {
			g.spawnWoodsCritter(g.rollProto([]int{0, 2}), 8)
		}
	case WoodsMillPond:
		for i := 0; i < 1+g.Roll(3); i++ {
			g.spawnWoodsCritter(g.rollProto(villagerProtos), 6)
		}
		for i := 0; i < 1+g.Roll(3); i++ {
			g.spawnWoodsCritter(g.rollProto([]int{rabbitProto, frogProto, crowProto, 15}), 6)
		}
		g.spawnWoodsCritter(g.rollProto([]int{0, 2}), 8)
	case WoodsCemetery:
		for i := 0; i < 2+g.Roll(2); i++ {
			g.spawnWoodsCritterNearLandmark(g.rollProto([]int{0, 1, 2, 4}), 11, 7)
		}
		for i := 0; i < 2+g.Roll(3); i++ {
			g.spawnWoodsCritter(g.rollProto([]int{0, 1, 2, 4}), 8)
		}
		g.spawnWoodsCritter(crowProto, 6)
	case WoodsMeteorField:
		for i := 0; i < 3+g.Roll(4); i++ {
			g.spawnWoodsCritter(g.rollProto([]int{1, 2, 4, 5, 13}), 7)
		}
		g.spawnWoodsCritterNearLandmark(g.rollProto([]int{1, 4, 5}), 9, 7)
	}
}

func (g *Game) rollProto(pool []int) int {
	if len(pool) == 0 {
		return -1
	}
	return pool[g.Roll(len(pool))]
}

func (g *Game) spawnWoodsCritter(proto, minDist int) *Critter {
	for i := 0; i < 300; i++ {
		x, y := g.Map.RandomFree(g.RNG)
		if isWoodsWater(g.Map.Tile(x, y)) || g.Map.Tile(x, y) == tileRiverBank || g.Map.Tile(x, y) == portalTile(0) {
			continue
		}
		if abs(x-g.Player.X)+abs(y-g.Player.Y) < minDist {
			continue
		}
		if g.CritterAt(x, y) != nil {
			continue
		}
		return g.SpawnCritter(proto, x, y)
	}
	return nil
}

func (g *Game) spawnWoodsCritterNearLandmark(proto, radius, minPlayerDist int) *Critter {
	for i := 0; i < 200; i++ {
		x := g.Map.LandmarkX - radius + g.Roll(radius*2+1)
		y := g.Map.LandmarkY - radius + g.Roll(radius*2+1)
		if !g.Map.IsPassable(x, y) || isWoodsWater(g.Map.Tile(x, y)) || g.Map.Tile(x, y) == portalTile(0) {
			continue
		}
		if abs(x-g.Player.X)+abs(y-g.Player.Y) < minPlayerDist {
			continue
		}
		if g.CritterAt(x, y) != nil {
			continue
		}
		return g.SpawnCritter(proto, x, y)
	}
	return nil
}

func (g *Game) spawnWoodsLoot(level, theme int) {
	pool := woodsLootPool(theme)
	if len(pool) == 0 {
		pool = filterAvailableItems(itemSpawnPool(0, level))
	}
	if len(pool) == 0 {
		return
	}
	for i := 0; i < 2+g.Roll(2); i++ {
		x, y := g.Map.RandomFree(g.RNG)
		if abs(x-g.Player.X)+abs(y-g.Player.Y) < 4 || g.Map.Tile(x, y) == portalTile(0) {
			continue
		}
		if g.ItemAt(x, y) != nil {
			continue
		}
		it := &Item{Proto: g.rollProto(pool), X: -1, Y: -1}
		if d := it.Def(); d != nil && d.Kind == KindMedical {
			it.Charges = d.Charges
		}
		g.PlaceGroundItemNear(it, x, y)
	}
}

func woodsLootPool(theme int) []int {
	var names []string
	switch theme {
	case WoodsChurch:
		names = []string{"holy water", "gold cross", "old tome", "bandage", "healing vial", "leather shard"}
	case WoodsPeasantHouses:
		names = []string{"apple", "bread", "cheese", "dried meat", "hatchet", "rope", "leather shard"}
	case WoodsMillPond:
		names = []string{"raw fish", "cooked fish", "fishing rod", "rope", "whetstone", "apple", "leather shard", "chain shard"}
	case WoodsCemetery:
		names = []string{"chain shard", "iron shard", "bone charm", "old tome", "shovel", "red potion"}
	case WoodsMeteorField:
		names = []string{"blue potion", "green potion", "bone charm", "silver chalice", "hatchet", "dried meat", "iron shard", "silver shard"}
	}
	ids := make([]int, 0, len(names))
	for _, name := range names {
		if proto := itemProtoByName(name); proto >= 0 {
			ids = append(ids, proto)
		}
	}
	return filterAvailableItems(ids)
}
