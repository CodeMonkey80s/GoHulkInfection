package game

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func slotName(slot int) string {
	switch slot {
	case SlotHands:
		return "Hands"
	case SlotBody:
		return "Body"
	case SlotBelt1:
		return "Belt 1"
	case SlotBelt2:
		return "Belt 2"
	case SlotBelt3:
		return "Belt 3"
	}
	return "?"
}

func slotKey(slot int) string {
	switch slot {
	case SlotHands:
		return "a"
	case SlotBelt1:
		return "c"
	case SlotBelt2:
		return "d"
	case SlotBelt3:
		return "e"
	}
	return "-"
}

func itemSummary(it *Item) string {
	if it == nil {
		return "empty"
	}
	d := it.Def()
	if d == nil {
		return it.Name()
	}
	switch d.Kind {
	case KindMelee, KindGrenade:
		out := d.Name
		if wd := d.Weapon(); wd != nil {
			out += fmt.Sprintf("  dmg %d-%d  av %d  dv %d", wd.DmgMin, wd.DmgMax, wd.AV, wd.DV)
		}
		return out
	case KindArmor:
		return fmt.Sprintf("%s  protection melee %d%% acid %d%%", d.Name, d.ArmorMelee, d.ArmorAcid)
	case KindMedical:
		return fmt.Sprintf("%s  heal %d  charges %d", d.Name, d.Heal, it.Charges)
	default:
		return d.Name
	}
}

func (u *UI) drawInventory() {
	g := u.g
	u.clear(0)
	u.putStr(centerText("[ i ] equipment", screenW), 0, "[ i ] equipment", 2, 0)

	y := 3
	slots := []int{SlotHands, SlotBody, SlotBelt1, SlotBelt2, SlotBelt3}
	for _, slot := range slots {
		it := g.Player.SlotItem(slot)
		locked := false
		if slot == SlotBelt2 && g.Player.Slots.BeltUnlocked < 2 {
			locked = true
		}
		if slot == SlotBelt3 && g.Player.Slots.BeltUnlocked < 3 {
			locked = true
		}
		key := slotKey(slot)
		label := slotName(slot)
		if locked {
			u.putStr(1, y, fmt.Sprintf("[%s] %-7s: locked (level %d)", key, label, 3+3*(slot-SlotBelt2)), 8, 0)
			y += 2
			continue
		}
		if key == "-" {
			u.putStr(1, y, fmt.Sprintf("   %-7s: %s", label, itemSummary(it)), 11, 0)
		} else {
			u.putStr(1, y, fmt.Sprintf("[%s] %-7s: %s", key, label, itemSummary(it)), 11, 0)
		}
		y += 2
	}

	ground := g.ItemAt(g.Player.X, g.Player.Y)
	u.putStr(1, y+1, "[ , ] item here: "+itemSummary(ground), 10, 0)

	u.fill(0, 46, screenW, 3, ' ', 8, 0)
	u.putStr(1, 46, "a hands   c/d/e belt   A-E drop   z quit", 8, 0)
	u.putStr(1, 47, "Body armour is worn; B drops it.", 8, 0)
	u.putStr(1, 48, "Pick up with , on the map.", 8, 0)
}

func (u *UI) handleInventory(e *tcell.EventKey) {
	g := u.g
	switch e.Key() {
	case tcell.KeyEscape:
		g.State = StatePlaying
		return
	}
	switch e.Rune() {
	case 'a':
		g.UseSlot(SlotHands)
	case 'c':
		g.UseSlot(SlotBelt1)
	case 'd':
		g.UseSlot(SlotBelt2)
	case 'e':
		g.UseSlot(SlotBelt3)
	case 'A':
		g.DropSlot(SlotHands)
	case 'B':
		g.DropSlot(SlotBody)
	case 'C':
		g.DropSlot(SlotBelt1)
	case 'D':
		g.DropSlot(SlotBelt2)
	case 'E':
		g.DropSlot(SlotBelt3)
	case 'z', 'q', 'i':
		g.State = StatePlaying
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
