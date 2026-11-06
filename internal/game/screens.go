package game

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func (u *UI) groundItems() []*Item {
	if u.g.Map == nil || u.g.Player == nil {
		return nil
	}
	return u.g.ItemsAt(u.g.Player.X, u.g.Player.Y)
}

func (u *UI) activeList() []*Item {
	if u.invTab == 0 {
		return u.g.Player.Inventory
	}
	return u.goundItems()
}

func (u *UI) goundItems() []*Item { return u.groundItems() }

func (u *UI) clampInvSel() {
	n := len(u.activeList())
	if n == 0 {
		u.invSel = 0
		return
	}
	if u.invSel >= n {
		u.invSel = n - 1
	}
	if u.invSel < 0 {
		u.invSel = 0
	}
}

func (u *UI) drawInventory() {
	g := u.g
	u.clear(0)
	u.putStr(centerText("[ i ] inventory", screenW), 0, "[ i ] inventory", 2, 0)
	inv := g.Player.Inventory
	ground := u.groundItems()
	invTitle := fmt.Sprintf("[ i ] items in inventory... (%d)", len(inv))
	groundTitle := fmt.Sprintf("[ . ] items lying here... (%d)", len(ground))
	u.putStr(1, 2, invTitle, 10, 0)
	u.putStr(41, 2, groundTitle, 10, 0)
	u.putStr(1, 3, strings.Repeat("~", 37), 2, 0)
	u.putStr(41, 3, strings.Repeat("~", 37), 2, 0)
	start := 0
	if u.invTab == 0 && u.invSel >= 16 {
		start = u.invSel - 15
	}
	for i := start; i < len(inv) && i < start+16; i++ {
		it := inv[i]
		y := 4 + (i - start)
		mark := "  "
		fg := uint8(11)
		bg := uint8(0)
		if u.invTab == 0 && i == u.invSel {
			mark = "> "
			fg = 0
			bg = 7
		}
		line := mark + it.Name()
		w := it.Def()
		if w != nil {
			line += strings.Repeat(".", max(1, 30-len(line))) + fmt.Sprintf("%dg", w.Weight)
		}
		u.putStr(1, y, fmt.Sprintf("%-38s", line), fg, bg)
	}
	gstart := 0
	if u.invTab == 1 && u.invSel >= 16 {
		gstart = u.invSel - 15
	}
	for i := gstart; i < len(ground) && i < gstart+16; i++ {
		it := ground[i]
		y := 4 + (i - gstart)
		mark := "  "
		fg := uint8(11)
		bg := uint8(0)
		if u.invTab == 1 && i == u.invSel {
			mark = "> "
			fg = 0
			bg = 7
		}
		line := mark + it.Name()
		if w := it.Def(); w != nil {
			line += strings.Repeat(".", max(1, 30-len(line))) + fmt.Sprintf("%dg", w.Weight)
		}
		u.putStr(41, y, fmt.Sprintf("%-38s", line), fg, bg)
	}
	if len(inv) == 0 {
		u.putStr(1, 4, "you are not carrying anything...", 8, 0)
	}
	if len(ground) == 0 {
		u.putStr(41, 4, "there is nothing lying here...", 8, 0)
	}
	// Equipment panel.
	u.putStr(1, 21, "[ # ] equipment...", 10, 0)
	u.putStr(1, 22, strings.Repeat("~", 37), 2, 0)
	slots := []string{"melee", "armor", "throwable", "tool"}
	for i, name := range slots {
		val := "none"
		if it := g.Player.Equip[i]; it != nil {
			val = it.Name()
			if it.Charges > 0 {
				val += fmt.Sprintf(" (%d)", it.Charges)
			}
		}
		u.putStr(3, 23+i, fmt.Sprintf("%-8s: %s", name, val), 11, 0)
	}
	// Status.
	u.putStr(1, 29, "[ & ] status...", 10, 0)
	u.putStr(1, 30, strings.Repeat("~", 37), 2, 0)
	u.putStr(3, 31, fmt.Sprintf("current weight: %dg", g.Player.Weight()), 11, 0)
	u.putStr(3, 32, fmt.Sprintf("maximum weight: %dg", g.Player.MaxWeight()), 11, 0)
	u.putStr(3, 33, fmt.Sprintf("burden: %s", g.Player.Burden()), 11, 0)
	// Description panel.
	u.putStr(41, 21, "[ ? ] description...", 10, 0)
	u.putStr(41, 22, strings.Repeat("~", 37), 2, 0)
	var sel *Item
	list := u.activeList()
	if len(list) > 0 && u.invSel < len(list) {
		sel = list[u.invSel]
	}
	if sel != nil {
		d := sel.Def()
		if d != nil {
			lines := wrapText(d.Desc, 37)
			for i, l := range lines {
				if i > 2 {
					break
				}
				u.putStr(41, 23+i, l, 11, 0)
			}
			y := 27
			u.putStr(41, y, "item data:", 8, 0)
			y++
			u.putStr(41, y, fmt.Sprintf("weight: %dg", d.Weight), 11, 0)
			y++
			if d.Kind == KindMelee || d.Kind == KindGrenade {
				wd := d.Weapon()
				if wd != nil {
					u.putStr(41, y, fmt.Sprintf("dmg: %d-%d  av: %d  dv: %d", wd.DmgMin, wd.DmgMax, wd.AV, wd.DV), 11, 0)
				} else {
					u.putStr(41, y, "dmg: unknown", 11, 0)
				}
				y++
			}
			if d.Kind == KindArmor {
				u.putStr(41, y, fmt.Sprintf("protection: melee %d%% acid %d%%", d.ArmorMelee, d.ArmorAcid), 11, 0)
				y++
			}
			if d.Kind == KindGrenade {
				u.putStr(41, y, "throwable: HE blast or fire", 11, 0)
				y++
			}
			if d.Kind == KindMedical {
				u.putStr(41, y, fmt.Sprintf("heal: %d  charges: %d", d.Heal, sel.Charges), 11, 0)
				y++
			}
		}
	}
	u.fill(0, 48, screenW, 2, ' ', 8, 0)
	u.putStr(1, 48, "z quit screen   tab pick up / drop item", 8, 0)
	u.putStr(1, 49, "enter wield / wear / use   d drop   <> select", 8, 0)
}

func (u *UI) handleInventory(e *tcell.EventKey) {
	g := u.g
	list := u.activeList()
	switch e.Key() {
	case tcell.KeyUp:
		u.invSel--
		u.clampInvSel()
	case tcell.KeyDown:
		u.invSel++
		u.clampInvSel()
	case tcell.KeyTab:
		u.invTab = 1 - u.invTab
		u.invSel = 0
	case tcell.KeyEscape:
		g.State = StatePlaying
	case tcell.KeyEnter:
		if len(list) == 0 || u.invSel >= len(list) {
			return
		}
		it := list[u.invSel]
		if u.invTab == 1 {
			g.RemoveItem(it)
			g.Player.AddItem(it)
			g.Msg("you picked up item: "+it.Name(), 10)
			u.invSel = 0
			u.clampInvSel()
		} else {
			g.UseItem(it)
			if g.Options.Bool("optInvQuitAfterItemUse") {
				g.State = StatePlaying
			}
		}
	default:
		switch e.Rune() {
		case 'z', 'q', 'i':
			g.State = StatePlaying
		case 'd':
			if u.invTab == 0 && len(list) > 0 && u.invSel < len(list) {
				g.DropItem(list[u.invSel])
				u.clampInvSel()
			}
		case 'j':
			u.invSel++
			u.clampInvSel()
		case 'k':
			u.invSel--
			u.clampInvSel()
		}
	}
}

func (u *UI) drawPlayerScreen() {
	g := u.g
	p := g.Player
	u.clear(0)
	u.putStr(centerText("[ @ ] character information", screenW), 0, "[ @ ] character information", 2, 0)
	u.putStr(1, 2, "[   ] attributes...", 10, 0)
	u.putStr(41, 2, "[ = ] skills...", 10, 0)
	u.putStr(1, 3, strings.Repeat("~", 37), 2, 0)
	u.putStr(41, 3, strings.Repeat("~", 37), 2, 0)
	attrs := []string{
		fmt.Sprintf("class: %s", ClassDefs[g.PlayerClass].Name),
		fmt.Sprintf("name: %s", p.Name),
		fmt.Sprintf("strength: %d", p.STR),
		fmt.Sprintf("toughness: %d", p.TGH),
		fmt.Sprintf("agility: %d", p.AGI),
		fmt.Sprintf("intelligence: %d", p.INT),
		fmt.Sprintf("health: %d/%d", p.HP, p.MaxHP),
		fmt.Sprintf("weight: %d/%d", p.Weight(), p.MaxWeight()),
		fmt.Sprintf("gold: %d", p.Gold),
	}
	for i, s := range attrs {
		u.putStr(3, 5+i, s, 11, 0)
	}
	u.putStr(41, 5, fmt.Sprintf("level: %d", p.Level), 11, 0)
	u.putStr(41, 6, fmt.Sprintf("experience: %d", p.XP), 11, 0)
	u.putStr(41, 7, fmt.Sprintf("need to gain: %d", 200*p.Level), 11, 0)
	u.putStr(41, 8, fmt.Sprintf("skill points: %d", g.FreePoints), 11, 0)
	u.putStr(41, 10, "[ %% ] experience...", 10, 0)
	u.putStr(41, 11, strings.Repeat("~", 37), 2, 0)
	for i, skill := range ActiveSkills {
		y := 13 + i
		val := p.Skill(skill)
		u.putStr(43, y, fmt.Sprintf("%-16s %3d  %s", SkillNames[skill], val, skillLevelName(val)), 11, 0)
	}
	u.fill(0, 48, screenW, 2, ' ', 8, 0)
	u.putStr(1, 48, "z quit screen", 8, 0)
	u.putStr(1, 49, "x skill levels", 8, 0)
}

func skillLevelName(v int) string {
	idx := 0
	switch {
	case v >= 98:
		idx = 11
	case v >= 94:
		idx = 10
	case v >= 88:
		idx = 9
	case v >= 80:
		idx = 8
	case v >= 70:
		idx = 7
	case v >= 60:
		idx = 6
	case v >= 50:
		idx = 5
	case v >= 40:
		idx = 4
	case v >= 30:
		idx = 3
	case v >= 20:
		idx = 2
	case v >= 10:
		idx = 1
	}
	if idx < len(SkillLevelNames) {
		return SkillLevelNames[idx]
	}
	return "useless"
}

func (u *UI) handlePlayerScreen(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyEscape:
		u.g.State = StatePlaying
	default:
		switch e.Rune() {
		case 'z', 'q', '@':
			u.g.State = StatePlaying
		case 'x':
			u.g.State = StateSkills
			u.skillSel = 0
		}
	}
}

func (u *UI) drawSkills() {
	g := u.g
	p := g.Player
	u.clear(0)
	u.putStr(centerText("[ x ] character skill levels", screenW), 0, "[ x ] character skill levels", 2, 0)
	u.putStr(1, 2, fmt.Sprintf("[ %% ] level / experience...  level: %d  experience: %d  skill points: %d", p.Level, p.XP, g.FreePoints), 10, 0)
	u.putStr(1, 3, strings.Repeat("~", 76), 2, 0)
	for i, skill := range ActiveSkills {
		y := 5 + i*2
		mark := "  "
		fg := uint8(11)
		bg := uint8(0)
		if i == u.skillSel {
			mark = "> "
			fg = 0
			bg = 7
		}
		line := fmt.Sprintf("%s%-18s %3d  %-10s", mark, SkillNames[skill], p.Skill(skill), skillLevelName(p.Skill(skill)))
		u.putStr(2, y, fmt.Sprintf("%-74s", line), fg, bg)
		if i == u.skillSel {
			u.putStr(3, y+1, SkillDescs[skill], 8, 0)
		}
	}
	u.fill(0, 42, screenW, 2, ' ', 8, 0)
	u.putStr(1, 42, "up/down select   enter/+ increase skill (1 point)   z quit", 8, 0)
	u.putStr(1, 43, fmt.Sprintf("you currently have %d free skill points", g.FreePoints), 10, 0)
}

func (u *UI) handleSkills(e *tcell.EventKey) {
	g := u.g
	switch e.Key() {
	case tcell.KeyUp:
		u.skillSel = (u.skillSel + len(ActiveSkills) - 1) % len(ActiveSkills)
	case tcell.KeyDown:
		u.skillSel = (u.skillSel + 1) % len(ActiveSkills)
	case tcell.KeyEscape:
		g.State = StatePlaying
	case tcell.KeyEnter:
		u.spendSkillPoint()
	default:
		switch e.Rune() {
		case 'z', 'q', 'x':
			g.State = StatePlaying
		case '+':
			u.spendSkillPoint()
		case 'k':
			u.skillSel = (u.skillSel + len(ActiveSkills) - 1) % len(ActiveSkills)
		case 'j':
			u.skillSel = (u.skillSel + 1) % len(ActiveSkills)
		}
	}
}

func (u *UI) spendSkillPoint() {
	g := u.g
	if g.FreePoints <= 0 {
		g.Msg("you do not have any free skill points.", 8)
		return
	}
	skill := ActiveSkills[u.skillSel]
	if g.Player.Skill(skill) >= 95 {
		g.Msg("this skill is already at maximum.", 8)
		return
	}
	g.FreePoints--
	inc := 1 + g.Roll(3)
	g.Player.Skills[skill] += inc
	if g.Player.Skills[skill] > 95 {
		g.Player.Skills[skill] = 95
	}
	g.Msg(fmt.Sprintf("%s increased by %d.", SkillNames[skill], inc), 10)
}

func (u *UI) drawOptions() {
	g := u.g
	opts := ActiveOptionDefs()
	u.clear(0)
	u.putStr(centerText("[ = ] options", screenW), 0, "[ = ] options", 2, 0)
	if u.optSel >= len(opts) {
		u.optSel = len(opts) - 1
	}
	if u.optSel < 0 {
		u.optSel = 0
	}
	visible := 39
	if u.optSel < u.optScroll {
		u.optScroll = u.optSel
	}
	if u.optSel >= u.optScroll+visible {
		u.optScroll = u.optSel - visible + 1
	}
	for i := 0; i < visible; i++ {
		idx := u.optScroll + i
		if idx >= len(opts) {
			break
		}
		d := &opts[idx]
		mark := "  "
		fg := uint8(11)
		bg := uint8(0)
		if idx == u.optSel {
			mark = "> "
			fg = 0
			bg = 7
		}
		line := fmt.Sprintf("%s%-38s %4d (%d)", mark, d.Name, g.Options.Get(d.Name), d.Max)
		u.putStr(1, 2+i, fmt.Sprintf("%-76s", line), fg, bg)
	}
	d := &opts[u.optSel]
	u.putStr(1, 42, "[ ? ] description...", 10, 0)
	u.putStr(1, 43, strings.Repeat("~", 76), 2, 0)
	for i, l := range wrapText(d.Desc, 76) {
		if i > 2 {
			break
		}
		u.putStr(1, 44+i, l, 11, 0)
	}
	u.fill(0, 48, screenW, 2, ' ', 8, 0)
	u.putStr(1, 48, "z quit screen   <6> increase   <4> decrease", 8, 0)
	u.putStr(1, 49, "enter toggle", 8, 0)
}

func (u *UI) handleOptions(e *tcell.EventKey) {
	g := u.g
	opts := ActiveOptionDefs()
	if u.optSel >= len(opts) {
		u.optSel = len(opts) - 1
	}
	if u.optSel < 0 {
		u.optSel = 0
	}
	switch e.Key() {
	case tcell.KeyUp:
		u.optSel--
		if u.optSel < 0 {
			u.optSel = len(opts) - 1
		}
	case tcell.KeyDown:
		u.optSel++
		if u.optSel >= len(opts) {
			u.optSel = 0
		}
	case tcell.KeyLeft:
		g.ChangeOption(opts[u.optSel].Name, -1)
	case tcell.KeyRight:
		g.ChangeOption(opts[u.optSel].Name, 1)
	case tcell.KeyEscape:
		g.State = StatePlaying
	case tcell.KeyEnter:
		g.ToggleOption(opts[u.optSel].Name)
	default:
		switch e.Rune() {
		case 'z', 'q', 'O':
			g.State = StatePlaying
		case '4':
			g.ChangeOption(opts[u.optSel].Name, -1)
		case '6':
			g.ChangeOption(opts[u.optSel].Name, 1)
		case 'k':
			u.optSel--
			if u.optSel < 0 {
				u.optSel = len(opts) - 1
			}
		case 'j':
			u.optSel++
			if u.optSel >= len(opts) {
				u.optSel = 0
			}
		}
	}
}

func keyLabel(k KeyBind) string {
	switch k.Key {
	case 0x103:
		return "up"
	case 0x102:
		return "down"
	case 0x105:
		return "right"
	case 0x104:
		return "left"
	case 0x1cb:
		return "enter"
	case 0x1b:
		return "esc"
	case 0x09:
		return "tab"
	case 0x20:
		return "space"
	case 0x1d0:
		return "pgup"
	case 0x1d1:
		return "pgdn"
	case 0x1ca:
		return "home"
	case 0x1cf:
		return "end"
	}
	if k.Key >= 32 && k.Key < 127 {
		return string(rune(k.Key))
	}
	return fmt.Sprintf("k%d", k.Key)
}

func (u *UI) drawHelp() {
	u.clear(0)
	u.putStr(centerText("[ ? ] bindings...", screenW), 0, "[ ? ] bindings...", 2, 0)
	u.putStr(1, 2, "[ ? ] keybindings", 10, 0)
	u.putStr(1, 3, strings.Repeat("~", 76), 2, 0)
	visible := 38
	if u.helpSel < u.helpScroll {
		u.helpScroll = u.helpSel
	}
	if u.helpSel >= u.helpScroll+visible {
		u.helpScroll = u.helpSel - visible + 1
	}
	for i := 0; i < visible; i++ {
		idx := u.helpScroll + i
		if idx >= len(ActiveKeyBindings) {
			break
		}
		k := ActiveKeyBindings[idx]
		mark := "  "
		fg := uint8(11)
		bg := uint8(0)
		if idx == u.helpSel {
			mark = "> "
			fg = 0
			bg = 7
		}
		line := fmt.Sprintf("%s<%s> %s", mark, keyLabel(k), k.Name)
		u.putStr(1, 4+i, fmt.Sprintf("%-76s", line), fg, bg)
	}
	k := ActiveKeyBindings[u.helpSel]
	u.putStr(1, 43, "[ ? ] description...", 10, 0)
	u.putStr(1, 44, strings.Repeat("~", 76), 2, 0)
	u.putStr(1, 45, k.Desc, 11, 0)
	u.fill(0, 48, screenW, 2, ' ', 8, 0)
	u.putStr(1, 48, "z quit screen   up/down select   pgup/pgdn scroll", 8, 0)
}

func (u *UI) handleHelp(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyUp:
		u.helpSel--
		if u.helpSel < 0 {
			u.helpSel = len(ActiveKeyBindings) - 1
		}
	case tcell.KeyDown:
		u.helpSel++
		if u.helpSel >= len(ActiveKeyBindings) {
			u.helpSel = 0
		}
	case tcell.KeyPgUp:
		u.helpSel -= 20
		if u.helpSel < 0 {
			u.helpSel = 0
		}
	case tcell.KeyPgDn:
		u.helpSel += 20
		if u.helpSel >= len(ActiveKeyBindings) {
			u.helpSel = len(ActiveKeyBindings) - 1
		}
	case tcell.KeyEscape:
		u.g.State = StatePlaying
	default:
		switch e.Rune() {
		case 'z', 'q', '?', '/':
			u.g.State = StatePlaying
		case 'j':
			u.helpSel++
			if u.helpSel >= len(ActiveKeyBindings) {
				u.helpSel = 0
			}
		case 'k':
			u.helpSel--
			if u.helpSel < 0 {
				u.helpSel = len(ActiveKeyBindings) - 1
			}
		}
	}
}

func (u *UI) drawMessageBuffer() {
	u.clear(0)
	u.putStr(centerText("[ m ] message buffer", screenW), 0, "[ m ] message buffer", 2, 0)
	u.putStr(1, 2, "[ m ] message buffer", 10, 0)
	u.putStr(1, 3, strings.Repeat("~", 76), 2, 0)
	var lines []string
	for _, m := range u.g.Messages {
		lines = append(lines, wrapText(m.Text, 76)...)
	}
	visible := 40
	if u.msgScroll > len(lines)-visible {
		u.msgScroll = len(lines) - visible
	}
	if u.msgScroll < 0 {
		u.msgScroll = 0
	}
	for i := 0; i < visible; i++ {
		idx := u.msgScroll + i
		if idx >= len(lines) {
			break
		}
		u.putStr(1, 4+i, lines[idx], 11, 0)
	}
	u.fill(0, 48, screenW, 2, ' ', 8, 0)
	u.putStr(1, 48, "z quit screen   up/down scroll   pgup/pgdn page", 8, 0)
}

func (u *UI) handleMessages(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyUp:
		u.msgScroll--
	case tcell.KeyDown:
		u.msgScroll++
	case tcell.KeyPgUp:
		u.msgScroll -= 20
	case tcell.KeyPgDn:
		u.msgScroll += 20
	case tcell.KeyEscape:
		u.g.State = StatePlaying
	default:
		switch e.Rune() {
		case 'z', 'q', 'm', '`':
			u.g.State = StatePlaying
		}
	}
}
