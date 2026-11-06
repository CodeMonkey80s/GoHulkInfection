package game

import "fmt"

func chebyshev(x0, y0, x1, y1 int) int {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	if dx > dy {
		return dx
	}
	return dy
}

func isAdjacent(a, b *Critter) bool {
	return a != nil && b != nil && chebyshev(a.X, a.Y, b.X, b.Y) <= 1 && !(a.X == b.X && a.Y == b.Y)
}

func (c *Critter) CurrentWeapon() *Item {
	return c.HandItem()
}

func (g *Game) attackSkill(c *Critter, mode int) int {
	switch mode {
	case 1:
		return c.Skill(SkillMelee)
	case 3:
		return c.Skill(SkillThrowing)
	default:
		// Monsters may still use natural attacks.  The player never attacks
		// without a weapon.
		return c.Skill(SkillUnarmed)
	}
}

// AttackValue computes the attacker's AV for a given attack mode:
// 0 = natural/unarmed, 1 = melee, 3 = throwing.
func (g *Game) AttackValue(c *Critter, mode int) int {
	av := g.attackSkill(c, mode) + c.AGI
	switch mode {
	case 1:
		if w := c.HandItem(); w != nil {
			if wd := w.Def().Weapon(); wd != nil {
				av += wd.AV * 2
			}
		}
	case 3:
		if w := c.ThrowableItem(); w != nil {
			if wd := w.Def().Weapon(); wd != nil {
				av += wd.AV * 2
			}
		}
	}
	if mode == 3 {
		// Throwing is harder when an enemy is adjacent.
		for _, other := range g.Critters {
			if other.IsAlive() && other != c && chebyshev(c.X, c.Y, other.X, other.Y) <= 1 {
				av = av * g.Options.Get("optMechanicAVModEnemyIsNear") / 100
				break
			}
		}
	}
	if av < 1 {
		av = 1
	}
	return av
}

// DefenseValue computes the defender's DV.
func (g *Game) DefenseValue(c *Critter) int {
	dv := c.Skill(SkillDodging) + c.AGI/5
	if w := c.HandItem(); w != nil {
		if wd := w.Def().Weapon(); wd != nil {
			dv += wd.DV
		}
	}
	if arm := c.BodyItem(); arm != nil {
		if stats, ok := ArmorStatsFor(arm.Def()); ok {
			dv += stats.DV
		}
	}
	if dv < 1 {
		dv = 1
	}
	return dv
}

func (g *Game) HitChance(attacker, defender *Critter, mode int) int {
	av := g.AttackValue(attacker, mode)
	dv := g.DefenseValue(defender)
	base := g.Options.Get("optMechanicBaseHitChance")
	av1 := g.Options.Get("optMechanicAVMod1")
	av2 := g.Options.Get("optMechanicAVMod2")
	dv1 := g.Options.Get("optMechanicDVMod1")
	dv2 := g.Options.Get("optMechanicDVMod2")
	if av2 <= 0 {
		av2 = 200
	}
	if dv2 <= 0 {
		dv2 = 200
	}
	chance := base + av*av1/av2 - dv*dv1/dv2
	chance += attacker.Skill(SkillFocusing) / 20
	if chance < 5 {
		chance = 5
	}
	if chance > 95 {
		chance = 95
	}
	return chance
}

// PlayerDirectionalAttack attacks the adjacent critter in the chosen
// direction.  Movement into an enemy also triggers a bump attack.
func (g *Game) PlayerDirectionalAttack(dx, dy int) {
	p := g.Player
	if p == nil || !p.IsAlive() {
		return
	}
	if dx == 0 && dy == 0 {
		g.Msg("choose a direction to attack.", 8)
		return
	}
	c := g.CritterAt(p.X+dx, p.Y+dy)
	if c == nil || !c.IsAlive() {
		g.Msg("there is nothing to attack in that direction.", 8)
		return
	}
	g.PlayerMelee(c)
}

func (g *Game) PlayerMelee(target *Critter) {
	p := g.Player
	weapon := p.HandItem()
	if weapon == nil {
		g.Msg("you are defenceless without a weapon.", 12)
		return
	}
	chance := g.HitChance(p, target, 1)
	if !g.Chance(chance) {
		g.Msg("you attack and miss.", 7)
		if target.IsAlive() {
			g.MonsterRetaliate(target)
		}
		g.EndPlayerTurn()
		return
	}
	dmg := g.RollRange(1, 4)
	desc := "weapon"
	if d := weapon.Def(); d != nil {
		desc = d.Name
		if wd := d.Weapon(); wd != nil {
			dmg = g.RollRange(wd.DmgMin, wd.DmgMax) + p.STR/3 + p.Skill(SkillMelee)/10
		}
	}
	if p.Skill(SkillFocusing) > 0 && g.Chance(p.Skill(SkillFocusing)/4) {
		dmg = dmg * 3 / 2
		g.Msg("your hit was critical!", 14)
	}
	g.Msg(fmt.Sprintf("you attack with %s and hit.", desc), 7)
	g.DamageCritter(target, dmg, p, 0)
	g.EndPlayerTurn()
}

func (g *Game) MonsterRetaliate(c *Critter) {
	if c == nil || !isAdjacent(c, g.Player) {
		return
	}
	g.MonsterAttack(c, g.Player)
}

func (g *Game) MonsterAttack(c, target *Critter) {
	if c == nil || target == nil || !c.IsAlive() || !target.IsAlive() {
		return
	}
	if !isAdjacent(c, target) {
		return
	}
	mode := 1
	if c.HandItem() == nil {
		mode = 0
	}
	chance := g.HitChance(c, target, mode)
	if !g.Chance(chance) {
		if target.Player {
			g.Msg("the "+c.NameOrUnknown()+" attacks you and miss.", 7)
		}
		return
	}
	dmg := c.STR/2 + c.Skill(SkillUnarmed)/8 + g.Roll(4)
	if w := c.HandItem(); w != nil {
		if wd := w.Def().Weapon(); wd != nil {
			dmg += g.RollRange(wd.DmgMin, wd.DmgMax) / 2
		}
	}
	if dmg < 1 {
		dmg = 1
	}
	if target.Player {
		g.Msg("the "+c.NameOrUnknown()+" attacks you and hit.", 12)
	}
	g.DamageCritter(target, dmg, c, 0)
}

func (g *Game) DamageCritter(target *Critter, dmg int, source *Critter, dtype int) {
	if target == nil || !target.IsAlive() {
		return
	}
	if target.Align == AlignNeutral && source != nil && source.Player {
		target.Align = AlignEnemy
	}
	if target.Player && g.Options.Bool("optPlayerGodMode") {
		g.Msg("your god mode absorbs the damage.", 8)
		return
	}
	if target.Player {
		// AP absorbs a percentage of non-acid damage.  Acid bypasses armour.
		if dtype != 2 {
			if arm := target.Slots.Body; arm != nil {
				if stats, ok := ArmorStatsFor(arm.Def()); ok {
					ensureArmorAP(arm)
					if arm.AP > 0 {
						absorbed := (dmg*stats.Absorb + 99) / 100
						if absorbed > arm.AP {
							absorbed = arm.AP
						}
						arm.AP -= absorbed
						dmg -= absorbed
						if arm.AP == 0 {
							g.Msg("your armour is depleted!", 12)
						}
					}
				}
			}
		}
	} else if arm := target.Equip[EquipArmor]; arm != nil {
		if d := arm.Def(); d != nil {
			factor := d.ArmorMelee
			switch dtype {
			case 1:
				factor = d.ArmorRanged
			case 2:
				factor = d.ArmorAcid
			}
			if factor <= 0 {
				factor = 100
			}
			dmg = dmg * factor / 100
			if dmg < 1 {
				dmg = 1
			}
			if g.Options.Bool("optGameSendDamageToItems") && g.Chance(5) {
				arm.Charges--
				if arm.Charges < 0 {
					target.Equip[EquipArmor] = nil
					g.Msg("your "+d.Name+" is shattered!", 12)
				}
			}
		}
	}
	if target.Player {
		g.Msg(fmt.Sprintf("the attack does %d damage.", dmg), 12)
	}
	target.HP -= dmg
	target.Blood += dmg / 2
	if target.HP <= 0 {
		g.KillCritter(target, source)
	}
}

func (g *Game) KillCritter(target, killer *Critter) {
	if target == nil || target.Dead {
		return
	}
	target.Dead = true
	target.HP = 0
	if target.Player {
		g.Msg("you died.", 12)
		g.Msg("=== session ended ===", 7)
		g.State = StateGameOver
		return
	}
	deadMsgs := []string{
		"the " + target.NameOrUnknown() + " dies in pain.",
		"you cleave the " + target.NameOrUnknown() + " apart.",
		"you hit a vital spot, the " + target.NameOrUnknown() + " is wasted.",
		"the " + target.NameOrUnknown() + " dies.",
	}
	g.Msg(deadMsgs[g.Roll(len(deadMsgs))], 14)
	if killer != nil && killer.Player {
		xp := target.XP
		g.Player.XP += xp
		g.Msg(fmt.Sprintf("you gain %d experience points.", xp), 10)
		g.CheckLevelUp()
	}
	g.CreateRemains(target)
	g.RemoveCritter(target)
}

func (g *Game) CreateRemains(c *Critter) {
	if g.Options.Bool("optMapCreateBlood") {
		g.Map.Blood[g.Map.idx(c.X, c.Y)] = 1
		for i := 0; i < 2+g.Roll(4); i++ {
			x := c.X + g.RollRange(-1, 1)
			y := c.Y + g.RollRange(-1, 1)
			if g.Map.InBounds(x, y) {
				g.Map.Blood[g.Map.idx(x, y)] = 1
			}
		}
	}
	if g.Options.Bool("optMapCreateCorpses") {
		g.PlaceGroundItemNear(&Item{Proto: 0, X: -1, Y: -1}, c.X, c.Y)
	}
	if g.Options.Bool("optMapCreateLimbs") && g.Chance(45) {
		proto := 1 + g.Roll(4)
		g.PlaceGroundItemNear(&Item{Proto: proto, X: -1, Y: -1}, c.X, c.Y)
	}
	if g.Chance(15) {
		proto := 41 // HE grenade
		if g.Roll(2) == 0 {
			proto = 42 // molotov cocktail
		}
		g.PlaceGroundItemNear(&Item{Proto: proto, X: -1, Y: -1}, c.X, c.Y)
	}
	for _, it := range c.Inventory {
		if !ItemAvailable(it.Proto) {
			continue
		}
		g.PlaceGroundItemNear(it, c.X, c.Y)
	}
	for _, it := range c.Equip {
		if it != nil && ItemAvailable(it.Proto) {
			g.PlaceGroundItemNear(it, c.X, c.Y)
		}
	}
}

func (g *Game) CheckLevelUp() {
	p := g.Player
	for {
		need := 200 * p.Level
		if p.XP < need {
			return
		}
		p.XP -= need
		p.Level++
		p.MaxHP += 5
		p.HP = p.MaxHP
		pts := g.Options.Get("optPlayerFreePointsPerLevel")
		if pts < 1 {
			pts = 1
		}
		g.FreePoints += pts
		g.Msg(fmt.Sprintf("you are on %d level !", p.Level), 14)
		g.Msg(fmt.Sprintf("you have gained: %d skill points, go to player screen and increase your skills.", pts), 10)
		if p.Slots.BeltUnlocked < 2 && p.Level >= 3 {
			p.Slots.BeltUnlocked = 2
			g.Msg("you can now use belt slot 2 [d].", 10)
		}
		if p.Slots.BeltUnlocked < 3 && p.Level >= 6 {
			p.Slots.BeltUnlocked = 3
			g.Msg("you can now use belt slot 3 [e].", 10)
		}
	}
}

// Explode damages everything in a radius.  kind 1 leaves fire behind
// (molotov); kind 0 is a plain high-explosive blast.
func (g *Game) Explode(x, y, radius, damage, kind int) {
	if radius < 1 {
		radius = 1
	}
	g.Msg("boom!", 12)
	for _, c := range g.Critters {
		if !c.IsAlive() {
			continue
		}
		d := chebyshev(x, y, c.X, c.Y)
		if d > radius {
			continue
		}
		dmg := damage * (radius - d + 1) / (radius + 1)
		if dmg < 1 {
			dmg = 1
		}
		g.DamageCritter(c, dmg, nil, 1)
	}
	if kind == 1 {
		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				nx, ny := x+dx, y+dy
				if !g.Map.InBounds(nx, ny) || chebyshev(x, y, nx, ny) > radius {
					continue
				}
				if g.Map.BlocksMove(nx, ny) {
					continue
				}
				g.Map.Fire[g.Map.idx(nx, ny)] = 15 + g.Roll(15)
			}
		}
	}
	if c := g.CritterAt(x, y); c != nil && c.IsAlive() {
		g.DamageCritter(c, damage, nil, 1)
	}
}

func (g *Game) hasThrowable() bool {
	if it := g.Player.ThrowableItem(); it != nil && ItemIsThrowable(it.Proto) {
		return true
	}
	return false
}

// ThrowGrenade throws the selected HE grenade or molotov at a visible tile
// within the player's current FOV radius.
func (g *Game) ThrowGrenade(x, y int) {
	it := g.Player.ThrowableItem()
	if it == nil || !ItemIsThrowable(it.Proto) {
		it = nil
	}
	if it == nil {
		g.Msg("you do not have any throwable items ready.", 12)
		return
	}
	d := it.Def()
	if d == nil {
		return
	}
	if !g.Map.InBounds(x, y) || !g.Map.Visible[g.Map.idx(x, y)] {
		g.Msg("you can only throw at visible tiles.", 8)
		return
	}
	radius := g.FOVRadius()
	if chebyshev(g.Player.X, g.Player.Y, x, y) > radius {
		g.Msg("that is beyond your throwing range.", 8)
		return
	}
	g.Player.Slots.Hands = nil
	g.Msg("you throw your "+d.Name+".", 7)
	// The projectile stops at the first blocking tile before the target.
	lx, ly := g.Player.X, g.Player.Y
	for _, p := range g.Map.Line(g.Player.X, g.Player.Y, x, y) {
		if g.Map.BlocksMove(p[0], p[1]) {
			break
		}
		lx, ly = p[0], p[1]
	}
	// Throwing skill controls scatter; an adjacent enemy makes the throw
	// less accurate.
	accuracy := 60 + g.Player.Skill(SkillThrowing)/2
	for _, c := range g.Critters {
		if c.IsAlive() && isAdjacent(g.Player, c) {
			accuracy = accuracy * g.Options.Get("optMechanicAVModEnemyIsNear") / 100
			break
		}
	}
	if !g.Chance(accuracy) {
		nx := lx + g.RollRange(-1, 1)
		ny := ly + g.RollRange(-1, 1)
		if g.Map.InBounds(nx, ny) && !g.Map.BlocksMove(nx, ny) {
			lx, ly = nx, ny
			g.Msg("the throw lands off target.", 8)
		}
	}
	blastRadius := 2 + d.Radius/4
	dmg := 10
	if wd := d.Weapon(); wd != nil {
		dmg = g.RollRange(wd.DmgMin, wd.DmgMax)
	}
	kind := 0
	if d.Name == "molotov cocktail" {
		kind = 1
	}
	g.Explode(lx, ly, blastRadius, dmg+g.Player.Skill(SkillThrowing)/4, kind)
	g.EndPlayerTurn()
}
