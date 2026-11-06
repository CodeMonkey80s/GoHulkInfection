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

func (c *Critter) HasRanged() bool {
	w := c.Equip[EquipRanged]
	if w == nil {
		return false
	}
	d := w.Def()
	return d != nil && d.Kind == KindRanged
}

func (c *Critter) CurrentWeapon() *Item {
	if w := c.Equip[EquipMelee]; w != nil {
		return w
	}
	return nil
}

func (g *Game) attackSkill(c *Critter, mode int) int {
	skill := SkillUnarmed
	switch mode {
	case 1:
		skill = SkillMelee
	case 2:
		skill = SkillRanged
	case 3:
		skill = SkillThrowing
	}
	return c.Skill(skill)
}

// AttackValue computes the attacker's AV for a given attack mode.
func (g *Game) AttackValue(c *Critter, mode int) int {
	av := g.attackSkill(c, mode) + c.AGI
	switch mode {
	case 1:
		if w := c.Equip[EquipMelee]; w != nil {
			if wd := w.Def().Weapon(); wd != nil {
				av += wd.AV * 2
			}
		}
	case 2:
		if w := c.Equip[EquipRanged]; w != nil {
			if wd := w.Def().Weapon(); wd != nil {
				av += wd.AV * 2
			}
		}
	case 3:
		if w := c.Equip[EquipGrenade]; w != nil {
			if wd := w.Def().Weapon(); wd != nil {
				av += wd.AV * 2
			}
		}
	}
	switch c.Tactics {
	case TacticAggressive:
		av = av * 120 / 100
	case TacticDefensive:
		av = av * 80 / 100
	}
	if mode == 2 || mode == 3 {
		// Ranged attacks are harder when an enemy is adjacent.
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
	if w := c.Equip[EquipMelee]; w != nil {
		if wd := w.Def().Weapon(); wd != nil {
			dv += wd.DV
		}
	} else if w := c.Equip[EquipRanged]; w != nil {
		if wd := w.Def().Weapon(); wd != nil {
			dv += wd.DV
		}
	}
	if arm := c.Equip[EquipArmor]; arm != nil {
		if d := arm.Def(); d != nil {
			// Heavier armor slightly reduces dodging.
			dv -= d.Weight / 4000
		}
	}
	switch c.Tactics {
	case TacticAggressive:
		dv = dv * 80 / 100
	case TacticDefensive:
		dv = dv * 120 / 100
	}
	switch c.Burden() {
	case "light":
		dv = dv * 90 / 100
	case "heavy":
		dv = dv * 70 / 100
	case "overloaded":
		dv = dv * 50 / 100
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
	// Focusing helps land criticals rather than raw hits; a small bonus keeps
	// high-level fights from stalling.
	chance += attacker.Skill(SkillFocusing) / 20
	if chance < 5 {
		chance = 5
	}
	if chance > 95 {
		chance = 95
	}
	return chance
}

func (g *Game) PlayerAttack(target *Critter) {
	if target == nil || !target.IsAlive() {
		g.Msg("there is no target set...", 12)
		return
	}
	if g.Player == nil {
		return
	}
	if isAdjacent(g.Player, target) {
		g.PlayerMelee(target)
		return
	}
	if g.Player.HasRanged() {
		g.PlayerFire(target)
		return
	}
	g.Msg("the "+target.NameOrUnknown()+" is too far for melee attack...", 12)
}

func (g *Game) PlayerMelee(target *Critter) {
	p := g.Player
	mode := 1
	weapon := p.Equip[EquipMelee]
	if weapon == nil {
		mode = 0
	}
	chance := g.HitChance(p, target, mode)
	if !g.Chance(chance) {
		g.Msg("you attack and miss.", 7)
		if target.IsAlive() {
			g.MonsterRetaliate(target)
		}
		g.EndPlayerTurn()
		return
	}
	dmg := g.RollRange(1, 4)
	desc := "fists"
	if weapon != nil {
		if d := weapon.Def(); d != nil {
			desc = d.Name
			if wd := d.Weapon(); wd != nil {
				dmg = g.RollRange(wd.DmgMin, wd.DmgMax) + p.STR/3 + p.Skill(SkillMelee)/10
			}
		}
	} else {
		dmg = p.Skill(SkillUnarmed)/8 + p.STR/2 + g.Roll(4)
	}
	if p.Skill(SkillFocusing) > 0 && g.Chance(p.Skill(SkillFocusing)/4) {
		dmg = dmg * 3 / 2
		g.Msg("your hit was critical!", 14)
	}
	g.Msg(fmt.Sprintf("you attack with %s and hit.", desc), 7)
	g.DamageCritter(target, dmg, p, 0)
	g.EndPlayerTurn()
}

func (g *Game) PlayerFire(target *Critter) {
	p := g.Player
	w := p.Equip[EquipRanged]
	if w == nil {
		g.Msg("you do not have a ranged weapon ready.", 12)
		return
	}
	d := w.Def()
	if d == nil {
		return
	}
	if !w.Active {
		g.Msg("the "+d.Name+" needs to be activated first...", 12)
		return
	}
	if d.AmmoType < 49 && d.AmmoType >= 0 {
		if w.Charges <= 0 {
			g.Msg("your "+d.Name+" has empty clip.", 12)
			return
		}
	}
	dist := chebyshev(p.X, p.Y, target.X, target.Y)
	if dist > d.Range {
		g.Msg("the "+target.NameOrUnknown()+" is too far for melee attack...", 12)
		return
	}
	if !g.Map.LineOfSight(p.X, p.Y, target.X, target.Y) {
		g.Msg("you do not see the target.", 12)
		return
	}
	shots := d.Burst
	if shots < 1 {
		shots = 1
	}
	if d.AmmoType < 49 && shots > w.Charges {
		shots = w.Charges
	}
	if shots < 1 {
		g.Msg("your "+d.Name+" has empty clip.", 12)
		return
	}
	g.Msg("you fire with "+d.Name+".", 7)
	hits := 0
	for i := 0; i < shots; i++ {
		if d.AmmoType < 49 && w.Charges > 0 {
			w.Charges--
		}
		// A critter standing in the line of fire may absorb the shot.
		if t := g.LineOfFireTarget(p, target); t != nil && t != target {
			if g.Chance(g.Options.Get("optMechanicLineOfFireHitChance")) {
				g.Msg("your shot hits "+t.NameOrUnknown()+" in the line of fire!", 12)
				g.DamageCritter(t, d.Damage, p, 1)
				continue
			}
		}
		chance := g.HitChance(p, target, 2)
		if !g.Chance(chance) {
			continue
		}
		dmg := 1
		if wd := d.Weapon(); wd != nil {
			dmg = g.RollRange(wd.DmgMin, wd.DmgMax) + p.Skill(SkillRanged)/10
		}
		g.DamageCritter(target, dmg, p, 1)
		hits++
		if !target.IsAlive() {
			break
		}
	}
	if hits == 0 && target.IsAlive() {
		g.Msg("you miss.", 7)
	}
	if d.AmmoType < 49 {
		g.Msg(fmt.Sprintf("%d rounds left.", w.Charges), 8)
	}
	g.EndPlayerTurn()
}

func (g *Game) LineOfFireTarget(from, target *Critter) *Critter {
	line := g.Map.Line(from.X, from.Y, target.X, target.Y)
	for _, p := range line {
		if p[0] == target.X && p[1] == target.Y {
			return nil
		}
		if c := g.CritterAt(p[0], p[1]); c != nil && c != from && c != target {
			return c
		}
		if g.Map.BlocksSight(p[0], p[1]) {
			return nil
		}
	}
	return nil
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
	if c.Equip[EquipMelee] == nil {
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
	if w := c.Equip[EquipMelee]; w != nil {
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
	if arm := target.Equip[EquipArmor]; arm != nil {
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
			// Armor wears down.
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
		"you slay the " + target.NameOrUnknown() + " with mighty punch.",
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
		proto := 0
		if c.Class == 0 {
			proto = 0
		}
		g.Items = append(g.Items, &Item{Proto: proto, X: c.X, Y: c.Y})
	}
	if g.Options.Bool("optMapCreateLimbs") && g.Chance(45) {
		proto := 1 + g.Roll(4)
		g.Items = append(g.Items, &Item{Proto: proto, X: c.X, Y: c.Y})
	}
	if g.Chance(40) {
		g.Items = append(g.Items, &Item{Proto: 41 + g.Roll(5), X: c.X, Y: c.Y})
	}
	// Drop carried equipment.
	for _, it := range c.Inventory {
		it.X, it.Y = c.X, c.Y
		g.Items = append(g.Items, it)
	}
	for _, it := range c.Equip {
		if it != nil {
			it.X, it.Y = c.X, c.Y
			g.Items = append(g.Items, it)
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
	}
}

// Explode damages everything in a radius and leaves fire/acid behind.
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
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			nx, ny := x+dx, y+dy
			if !g.Map.InBounds(nx, ny) || chebyshev(x, y, nx, ny) > radius {
				continue
			}
			if g.Map.BlocksMove(nx, ny) {
				continue
			}
			switch kind {
			case 1: // incendiary
				g.Map.Fire[g.Map.idx(nx, ny)] = 15 + g.Roll(15)
			case 2: // acid
				g.Map.Acid[g.Map.idx(nx, ny)] = 15 + g.Roll(15)
			case 3: // smoke
				g.Map.Smoke[g.Map.idx(nx, ny)] = 20 + g.Roll(20)
			}
		}
	}
	if c := g.CritterAt(x, y); c != nil && c.IsAlive() {
		g.DamageCritter(c, damage, nil, 1)
	}
}

func (g *Game) ThrowGrenade(x, y int) {
	it := g.Player.Equip[EquipGrenade]
	if it == nil {
		for _, cand := range g.Player.Inventory {
			if d := cand.Def(); d != nil && d.Kind == KindGrenade {
				it = cand
				break
			}
		}
	}
	if it == nil {
		g.Msg("you do not have any throwable items ready.", 12)
		return
	}
	d := it.Def()
	if d == nil {
		return
	}
	// Remove from wherever it is.
	if it == g.Player.Equip[EquipGrenade] {
		g.Player.Equip[EquipGrenade] = nil
	} else {
		for i, cand := range g.Player.Inventory {
			if cand == it {
				g.Player.RemoveItemAt(i)
				break
			}
		}
	}
	g.Msg("you throw your "+d.Name+".", 7)
	// Grenade lands at target or first wall.
	lx, ly := g.Player.X, g.Player.Y
	line := g.Map.Line(g.Player.X, g.Player.Y, x, y)
	for _, p := range line {
		if g.Map.BlocksMove(p[0], p[1]) {
			break
		}
		lx, ly = p[0], p[1]
		if chebyshev(g.Player.X, g.Player.Y, p[0], p[1]) > d.Range {
			break
		}
	}
	kind := 0
	switch d.WeaponSub {
	case 15:
		kind = 1
	case 0:
		kind = 2
	}
	if d.Name == "smoke grenade" {
		kind = 3
	}
	if d.Name == "flashbang grenade" {
		for _, c := range g.Critters {
			if c.IsAlive() && chebyshev(c.X, c.Y, lx, ly) <= 3 {
				c.Alert = 0
			}
		}
		g.Msg("a blinding flash fills the area.", 15)
		g.EndPlayerTurn()
		return
	}
	radius := 2 + d.Radius/4
	dmg := 10
	if wd := d.Weapon(); wd != nil {
		dmg = g.RollRange(wd.DmgMin, wd.DmgMax)
	}
	g.Explode(lx, ly, radius, dmg+g.Player.Skill(SkillThrowing)/4, kind)
	g.EndPlayerTurn()
}

func (g *Game) Reload(c *Critter) {
	if c == nil || !c.Player {
		return
	}
	w := c.Equip[EquipRanged]
	if w == nil {
		g.Msg("there is nothing to reload.", 12)
		return
	}
	d := w.Def()
	if d == nil {
		return
	}
	if d.AmmoType >= 49 || d.MaxAmmo <= 0 {
		g.Msg("your "+d.Name+" does not need to be reloaded. the clip is full.", 8)
		return
	}
	if w.Charges >= d.MaxAmmo {
		g.Msg("your "+d.Name+" does not need to be reloaded. the clip is full.", 8)
		return
	}
	for i, it := range c.Inventory {
		ad := it.Def()
		if ad == nil || ad.Kind != KindAmmo || ad.AmmoType != d.AmmoType {
			continue
		}
		need := d.MaxAmmo - w.Charges
		take := it.Charges
		if take > need {
			take = need
		}
		w.Charges += take
		it.Charges -= take
		if it.Charges <= 0 {
			c.RemoveItemAt(i)
		}
		g.Msg("you reload your "+d.Name+".", 8)
		g.Msg(fmt.Sprintf("%d rounds left of %s.", w.Charges, ad.Name), 8)
		g.EndPlayerTurn()
		return
	}
	g.Msg("you do not have proper ammunition for your "+d.Name+".", 12)
}

func (g *Game) CheckAmmo() {
	w := g.Player.Equip[EquipRanged]
	if w == nil {
		g.Msg("you are not holding any weapon that uses ammo.", 8)
		return
	}
	d := w.Def()
	if d == nil {
		return
	}
	if d.AmmoType >= 49 {
		g.Msg("the "+d.Name+" uses an internal power source.", 8)
		return
	}
	g.Msg(fmt.Sprintf("you check the ammo for your %s, there are %d rounds left.", d.Name, w.Charges), 8)
	g.EndPlayerTurn()
}
