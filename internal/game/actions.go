package game

import "fmt"

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
			g.PlayerMelee(c)
			return
		}
		g.Msg("the "+c.NameOrUnknown()+" blocks your way...", 8)
		return
	}
	d := g.Map.Def(nx, ny)
	if g.Options.Bool("optPlayerGhostMode") {
		g.Player.X, g.Player.Y = nx, ny
		g.UpdateFOV()
		g.announceGroundItem()
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
	g.announceGroundItem()
	g.EndPlayerTurn()
}

func (g *Game) announceGroundItem() {
	if it := g.ItemAt(g.Player.X, g.Player.Y); it != nil {
		g.Msg(capitalizeFirst(it.Name())+" is lying here.", 10)
	}
}

func (g *Game) WaitTurn() {
	g.Msg("you wait and rest.", 8)
	g.EndPlayerTurn()
}

func (g *Game) PickUp() {
	if g.State != StatePlaying && g.State != StateInventory {
		return
	}
	it := g.ItemAt(g.Player.X, g.Player.Y)
	if it == nil {
		g.Msg("there is nothing to pick up here.", 8)
		return
	}
	d := it.Def()
	if d == nil {
		return
	}
	if !g.Player.AutoPlaceItem(it) {
		g.Msg("you have no free slot for "+d.Name+".", 8)
		return
	}
	g.RemoveItem(it)
	g.Msg("you picked up item: "+d.Name, 10)
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
