package game

import (
	"fmt"
	"sort"
)

func (g *Game) EndPlayerTurn() {
	if g.State == StateGameOver || g.State == StateVictory {
		return
	}
	g.Turn++
	g.applyPlayerTileEffects()
	if g.Player.Poison > 0 {
		if g.Turn%3 == 0 {
			dmg := 1 + g.Roll(3)
			g.Player.HP -= dmg
			g.Player.Poison--
			g.Msg(fmt.Sprintf("the poison burns in your veins (%d).", dmg), 10)
			if g.Player.HP <= 0 {
				g.KillCritter(g.Player, nil)
			}
		}
	}
	if g.Turn%10 == 0 && g.Player.IsAlive() {
		heal := 1 + g.Player.Skill(SkillMetabolism)/25
		g.Player.HP += heal
		if g.Player.HP > g.Player.MaxHP {
			g.Player.HP = g.Player.MaxHP
		}
		g.Player.Blood = max(0, g.Player.Blood-1)
	}
	g.updateWeather()
	g.UpdateMonsters()
	g.decayTileEffects()
	g.UpdateFOV()
	if g.PlayerTarget() == nil {
		g.Target = -1
	}
	if !g.Player.IsAlive() {
		g.State = StateGameOver
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *Game) applyPlayerTileEffects() {
	if !g.Map.InBounds(g.Player.X, g.Player.Y) {
		return
	}
	i := g.Map.idx(g.Player.X, g.Player.Y)
	if g.Map.Fire[i] > 0 && g.Turn%2 == 0 {
		g.Player.HP -= 4
		g.Msg("you step into the burning flames.", 12)
	}
	if g.Map.Acid[i] > 0 && g.Turn%2 == 0 {
		g.Player.HP -= 6 // acid ignores armor
		g.Msg("you step into the pool of toxic acid.", 10)
	}
	if g.Map.Smoke[i] > 0 {
		g.Msg("there is smoke here.", 8)
	}
	if g.Player.HP <= 0 {
		g.KillCritter(g.Player, nil)
	}
}

func (g *Game) decayTileEffects() {
	for i := range g.Map.Fire {
		if g.Map.Fire[i] > 0 {
			g.Map.Fire[i]--
		}
		if g.Map.Acid[i] > 0 {
			g.Map.Acid[i]--
		}
		if g.Map.Smoke[i] > 0 {
			g.Map.Smoke[i]--
		}
	}
}

func (g *Game) updateWeather() {
	if g.Theme != 0 {
		g.Weather = WeatherClear
		return
	}
	if g.Weather != WeatherClear {
		g.WeatherLeft--
		if g.WeatherLeft <= 0 {
			g.Weather = WeatherClear
			g.Msg("the rain stops.", 8)
		}
		return
	}
	if g.Turn%30 == 0 && g.Chance(20) {
		g.Weather = WeatherRain
		g.WeatherLeft = 120 + g.Roll(240)
		g.Msg("it starts to rain.", 11)
	}
}

func (g *Game) MovePlayer(dx, dy int) {
	if g.State != StatePlaying || !g.Player.IsAlive() {
		return
	}
	nx, ny := g.Player.X+dx, g.Player.Y+dy
	if !g.Map.InBounds(nx, ny) {
		g.Msg("you cannot go there...", 8)
		return
	}
	if c := g.CritterAt(nx, ny); c != nil {
		if g.Options.Bool("optPlayerAttackOnMove") {
			g.SetTarget(c)
			g.PlayerAttack(c)
			return
		}
		g.Msg("the "+c.NameOrUnknown()+" blocks your way...", 8)
		return
	}
	d := g.Map.Def(nx, ny)
	if g.Options.Bool("optPlayerGhostMode") {
		g.Player.X, g.Player.Y = nx, ny
		g.UpdateFOV()
		g.EndPlayerTurn()
		return
	}
	if d.Type == 3 && d.BlockMove {
		if g.Options.Bool("optPlayerAutoOpenDoor") {
			g.OpenDoor(nx, ny)
			g.Msg("you open the "+d.Name+".", 8)
			g.EndPlayerTurn()
			return
		}
		g.Msg("you cannot go there; "+d.Name+" blocks your way...", 8)
		return
	}
	if d.BlockMove {
		g.Msg("you cannot go there; "+d.Name+" blocks your way...", 8)
		return
	}
	g.Player.X, g.Player.Y = nx, ny
	g.UpdateFOV()
	g.EndPlayerTurn()
}

func (g *Game) WaitTurn() {
	g.Msg("you wait and rest.", 8)
	g.EndPlayerTurn()
}

func (g *Game) PickUp() {
	if g.State != StatePlaying {
		return
	}
	items := g.ItemsAt(g.Player.X, g.Player.Y)
	if len(items) == 0 {
		g.Msg("there is nothing to pick up here.", 8)
		return
	}
	maxW := g.Player.MaxWeight()
	for _, it := range items {
		d := it.Def()
		if d == nil {
			continue
		}
		if g.Player.Weight()+d.Weight > maxW+maxW/4 {
			g.Msg("the "+d.Name+" weights too much.", 12)
			continue
		}
		g.RemoveItem(it)
		g.Player.AddItem(it)
		g.Msg("you picked up item: "+d.Name, 10)
	}
	g.EndPlayerTurn()
}

func (g *Game) DropItem(it *Item) {
	if it == nil {
		return
	}
	// Unequip if necessary.
	for i, e := range g.Player.Equip {
		if e == it {
			g.Player.Equip[i] = nil
		}
	}
	for i, cand := range g.Player.Inventory {
		if cand == it {
			g.Player.RemoveItemAt(i)
			break
		}
	}
	it.X, it.Y = g.Player.X, g.Player.Y
	g.Items = append(g.Items, it)
	g.Msg("you dropped item: "+it.Name(), 8)
	g.EndPlayerTurn()
}

func (g *Game) UseItem(it *Item) {
	if it == nil {
		return
	}
	d := it.Def()
	if d == nil {
		return
	}
	switch d.Kind {
	case KindMedical:
		if it.Charges <= 0 {
			it.Charges = d.Charges
		}
		heal := d.Heal
		if heal <= 0 {
			heal = 10
		}
		g.Player.HP += heal
		if g.Player.HP > g.Player.MaxHP {
			g.Player.HP = g.Player.MaxHP
		}
		it.Charges--
		g.Msg(fmt.Sprintf("you use %s and heal %d health.", d.Name, heal), 10)
		if it.Charges <= 0 {
			g.removeFromInventory(it)
		}
		g.EndPlayerTurn()
	case KindArmor:
		g.WearItem(it)
	case KindMelee, KindRanged, KindGrenade:
		g.WieldItem(it)
	case KindAmmo:
		if w := g.Player.Equip[EquipRanged]; w != nil {
			if wd := w.Def(); wd != nil && wd.AmmoType == d.AmmoType {
				g.Reload(g.Player)
				return
			}
		}
		g.Msg("you cannot use "+d.Name+" right now.", 8)
	default:
		g.Msg("you cannot use "+d.Name+".", 8)
	}
}

func (g *Game) removeFromInventory(it *Item) {
	for i, cand := range g.Player.Inventory {
		if cand == it {
			g.Player.RemoveItemAt(i)
			return
		}
	}
}

func (g *Game) WieldItem(it *Item) {
	d := it.Def()
	if d == nil {
		return
	}
	// Remove from inventory if present.
	g.removeFromInventory(it)
	switch d.Kind {
	case KindMelee:
		if old := g.Player.Equip[EquipMelee]; old != nil {
			g.Player.AddItem(old)
		}
		it.Active = true
		g.Player.Equip[EquipMelee] = it
		g.Msg("you are now using "+d.Name+".", 10)
	case KindRanged:
		if old := g.Player.Equip[EquipRanged]; old != nil {
			g.Player.AddItem(old)
		}
		it.Active = true
		if it.Charges == 0 && d.AmmoType < 49 {
			it.Charges = d.MaxAmmo
		}
		g.Player.Equip[EquipRanged] = it
		g.Msg("you are now using "+d.Name+".", 10)
	case KindArmor:
		if old := g.Player.Equip[EquipArmor]; old != nil {
			g.Player.AddItem(old)
			g.Msg("you are now taken off "+old.Name()+".", 8)
		}
		it.Active = true
		g.Player.Equip[EquipArmor] = it
		g.Msg("you are now wearing "+d.Name+".", 10)
	case KindGrenade:
		if old := g.Player.Equip[EquipGrenade]; old != nil {
			g.Player.AddItem(old)
		}
		g.Player.Equip[EquipGrenade] = it
		g.Msg("you are now using "+d.Name+".", 10)
	default:
		g.Player.AddItem(it)
	}
	g.EndPlayerTurn()
}

func (g *Game) WearItem(it *Item) { g.WieldItem(it) }

func (g *Game) ActivateWeapon() {
	w := g.Player.Equip[EquipRanged]
	if w == nil {
		w = g.Player.Equip[EquipMelee]
	}
	if w == nil {
		g.Msg("there is nothing to activate...", 8)
		return
	}
	d := w.Def()
	if d == nil {
		return
	}
	if w.Active {
		w.Active = false
		g.Msg("you deactivated your "+d.Name, 8)
	} else {
		w.Active = true
		g.Msg("you activated your "+d.Name, 10)
	}
	g.EndPlayerTurn()
}

func (g *Game) SwapWeapons() {
	a := g.Player.Equip[EquipMelee]
	b := g.Player.Equip[EquipRanged]
	if a == nil && b == nil {
		g.Msg("there is nothing to swap.", 8)
		return
	}
	if a != nil && b != nil && isAdjacent(g.Player, g.PlayerTargetOrNearest()) {
		g.Msg("you cannot swap weapons with an enemy nearby.", 12)
		return
	}
	g.Player.Equip[EquipMelee], g.Player.Equip[EquipRanged] = b, a
	g.Msg("you swap your weapons.", 8)
	g.EndPlayerTurn()
}

func (g *Game) CleanYourself() {
	g.Player.Blood = 0
	g.Msg("you clean up yourself and your hands.", 8)
	g.EndPlayerTurn()
}

func (g *Game) Search() {
	skill := g.Player.Skill(SkillSearching)
	g.Msg("searching...", 8)
	if g.Chance(20 + skill/2) {
		// Find a random nearby item if any; otherwise spawn a small find.
		x, y := g.Player.X, g.Player.Y
		found := false
		for dy := -4; dy <= 4 && !found; dy++ {
			for dx := -4; dx <= 4 && !found; dx++ {
				for _, it := range g.Items {
					if it.X == x+dx && it.Y == y+dy {
						g.Msg("you found something !", 10)
						found = true
					}
				}
			}
		}
		if !found {
			proto := 48
			if g.Chance(30) {
				proto = 46
			}
			g.Items = append(g.Items, &Item{Proto: proto, X: x, Y: y, Charges: ItemDefs[proto].Charges})
			g.Msg("you found something !", 10)
		}
	}
	g.EndPlayerTurn()
}

func (g *Game) Bash() {
	done := false
	for _, d := range dirs8 {
		x, y := g.Player.X+d[0], g.Player.Y+d[1]
		td := g.Map.Def(x, y)
		if td.Type == 3 || td.BlockMove {
			chance := 20 + g.Player.STR*2
			if g.Chance(chance) {
				g.Map.SetTile(x, y, floorTile(g.Theme))
				g.Msg("you smash the "+td.Name+"!", 12)
				done = true
			} else {
				g.Msg("you bash the "+td.Name+" but it holds.", 8)
				done = true
			}
			break
		}
	}
	if !done {
		g.Msg("there is nothing to destroy here.", 8)
		return
	}
	g.EndPlayerTurn()
}

func (g *Game) Operate() {
	// Stand on a door?
	if d := g.Map.Def(g.Player.X, g.Player.Y); d.Type == 3 {
		if d.BlockMove {
			g.OpenDoor(g.Player.X, g.Player.Y)
		} else {
			g.CloseDoor(g.Player.X, g.Player.Y)
		}
		g.EndPlayerTurn()
		return
	}
	for _, dir := range dirs8 {
		x, y := g.Player.X+dir[0], g.Player.Y+dir[1]
		if d := g.Map.Def(x, y); d.Type == 3 {
			if d.BlockMove {
				g.OpenDoor(x, y)
				g.Msg("the "+d.Name+" opens.", 8)
			} else {
				g.CloseDoor(x, y)
				g.Msg("the "+d.Name+" closes.", 8)
			}
			g.EndPlayerTurn()
			return
		}
	}
	g.Msg("there is nothing to operate here.", 8)
}

func (g *Game) UsePortal() {
	d := g.Map.Def(g.Player.X, g.Player.Y)
	if d.Type != 5 && d.Type != 3 && d.Type != 1 {
		g.Msg("there is nothing to enter here.", 8)
		return
	}
	// Portal/stairs/elevator tiles.
	name := d.Name
	if name != "time portal" && name != "stone stairs" && name != "elevator" {
		g.Msg("there is nothing to enter here.", 8)
		return
	}
	if g.Level >= 27 {
		g.Msg("you survive", 14)
		g.Msg("as you enter the time portal, you feel highly accelerated and years seem to be seconds ... after a while you were thrown out on the other side of the portal. you stand up slowly and look around, examining the cold, shiny walls of your cell ...", 10)
		g.State = StateVictory
		return
	}
	g.Msg("you enter the "+name+".", 10)
	g.ChangeLevel(g.Level + 1)
}

// ChangeLevelRelative jumps directly to the next/previous map.  It is a
// convenience/debug shortcut; the normal way is still to find the portal.
func (g *Game) ChangeLevelRelative(delta int) {
	if g.Player == nil || !g.Player.IsAlive() {
		return
	}
	level := g.Level + delta
	if level < 0 {
		level = 0
	}
	if level > 27 {
		level = 27
	}
	if level == g.Level {
		if delta > 0 {
			g.Msg("you are already on the last map.", 8)
		} else {
			g.Msg("you are already on the first map.", 8)
		}
		return
	}
	g.Msg(fmt.Sprintf("you travel directly to map level %d.", level), 10)
	g.ChangeLevel(level)
}

func (g *Game) ChangeLevel(level int) {
	// Remember player state.
	p := g.Player
	g.GenerateLevel(level)
	p.X, p.Y = g.Player.X, g.Player.Y
	// GenerateLevel created a fresh "player" placeholder; restore the real one.
	g.Player = p
	g.Player.X, g.Player.Y = p.X, p.Y
	g.UpdateFOV()
}

func (g *Game) ChangeTactics() {
	p := g.Player
	p.Tactics = (p.Tactics + 1) % 3
	names := []string{"defensive", "normal", "aggressive"}
	g.Msg("you changed tactic to "+names[p.Tactics], 8)
	g.EndPlayerTurn()
}

func (g *Game) PlayerTargetOrNearest() *Critter {
	if t := g.PlayerTarget(); t != nil {
		return t
	}
	return g.NearestTarget()
}

func (g *Game) NearestTarget() *Critter {
	var best *Critter
	bestD := 99999
	for _, c := range g.Critters {
		if !c.IsAlive() {
			continue
		}
		i := g.Map.idx(c.X, c.Y)
		if !g.Map.Visible[i] {
			continue
		}
		d := chebyshev(g.Player.X, g.Player.Y, c.X, c.Y)
		if d < bestD {
			bestD = d
			best = c
		}
	}
	return best
}

func (g *Game) BuildAimTargets() {
	type tc struct {
		c *Critter
		d int
	}
	var list []tc
	for _, c := range g.Critters {
		if !c.IsAlive() {
			continue
		}
		i := g.Map.idx(c.X, c.Y)
		if !g.Map.Visible[i] {
			continue
		}
		list = append(list, tc{c, chebyshev(g.Player.X, g.Player.Y, c.X, c.Y)})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].d < list[j].d })
	g.AimTargets = nil
	for _, t := range list {
		g.AimTargets = append(g.AimTargets, g.indexOfCritter(t.c))
	}
	if len(g.AimTargets) == 0 {
		g.Msg("there is no target set...", 8)
		return
	}
	g.AimPos = 0
	if t := g.PlayerTarget(); t != nil {
		for i, idx := range g.AimTargets {
			if idx >= 0 && idx < len(g.Critters) && g.Critters[idx] == t {
				g.AimPos = i
				break
			}
		}
	}
	g.SetTarget(g.Critters[g.AimTargets[g.AimPos]])
	g.Msg("new target selected, now targetting "+g.Critters[g.AimTargets[g.AimPos]].NameOrUnknown()+".", 8)
}

func (g *Game) indexOfCritter(c *Critter) int {
	for i, x := range g.Critters {
		if x == c {
			return i
		}
	}
	return -1
}

func (g *Game) NextAimTarget(dir int) {
	if len(g.AimTargets) == 0 {
		return
	}
	g.AimPos = (g.AimPos + dir + len(g.AimTargets)) % len(g.AimTargets)
	idx := g.AimTargets[g.AimPos]
	if idx >= 0 && idx < len(g.Critters) {
		g.SetTarget(g.Critters[idx])
		g.Msg("new target selected, now targetting "+g.Critters[idx].NameOrUnknown()+".", 8)
	}
}

func (g *Game) NextNearestTarget() {
	t := g.NearestTarget()
	if t == nil {
		g.Msg("trying to find nearest target... not found !", 8)
		return
	}
	g.SetTarget(t)
	g.Msg("new target selected, now targetting "+t.NameOrUnknown()+".", 8)
}

func (g *Game) ClearTarget() {
	g.Target = -1
	g.Msg("target cleared.", 8)
	if g.Options.Bool("optPlayerClearTargetResetTactics") {
		g.Player.Tactics = TacticNormal
		g.Msg(" and tactic set to normal", 8)
	}
}
