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
	case SlotHead:
		return "Head"
	case SlotRobe1:
		return "Robe 1"
	case SlotRobe2:
		return "Robe 2"
	case SlotRobe3:
		return "Robe 3"
	case SlotRobe4:
		return "Robe 4"
	case SlotRobe5:
		return "Robe 5"
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

func (g *Game) EquipmentSlots() []int {
	slots := []int{SlotHands, SlotHead}
	for i := 0; i < g.Player.Slots.RobeSlots && i < MaxRobeSlots; i++ {
		slots = append(slots, SlotRobe1+i)
	}
	return slots
}

func slotKey(slot int) string {
	switch slot {
	case SlotHands:
		return "a"
	case SlotHead:
		return "h"
	case SlotRobe1:
		return "c"
	case SlotRobe2:
		return "d"
	case SlotRobe3:
		return "r"
	case SlotRobe4:
		return "4"
	case SlotRobe5:
		return "5"
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
		ensureArmorAP(it)
		if stats, ok := ArmorStatsFor(d); ok {
			return fmt.Sprintf("%s  AP %d/%d  absorb %d%%  dv %d", d.Name, it.AP, stats.MaxAP, stats.Absorb, stats.DV)
		}
		return d.Name
	case KindMedical:
		return fmt.Sprintf("%s  heal %d  charges %d", d.Name, d.Heal, it.Charges)
	case KindFood:
		return fmt.Sprintf("%s  heal %d  charges %d", d.Name, d.Heal, it.Charges)
	case KindShard:
		return fmt.Sprintf("%s  AP max %d", d.Name, d.Charges)
	case KindTool, KindRelic, KindHead:
		return d.Desc
	default:
		return d.Name
	}
}

func (u *UI) drawInventory() {
	g := u.g
	u.fill(0, 0, screenW, 24, ' ', 7, 0)
	u.putStr(centerText("[ e ] equipment", screenW), 0, "[ e ] equipment", 2, 0)
	u.putStr(1, 1, "key  slot     item              stats                         action", 8, 0)

	slots := g.EquipmentSlots()
	y := 3
	for i, slot := range slots {
		it := g.Player.SlotItem(slot)
		key := slotKey(slot)
		if key == "-" {
			key = " "
		}
		label := slotName(slot)
		name := "-"
		stats := ""
		if it != nil {
			name = it.Name()
			stats = slotStats(it)
		}
		action := slotAction(slot, it)
		fg, bg := uint8(11), uint8(0)
		if i == u.slotSel {
			fg, bg = 0, 7
		}
		u.putStr(1, y, fmt.Sprintf("[%s] %-7s %-16s %-28s %s", key, label, name, stats, action), fg, bg)
		y++
	}

	u.putStr(16, 10, "a  use hands          h  head", 8, 0)
	u.putStr(16, 11, "c/d/r  use robe slot 1-3", 8, 0)
	u.putStr(16, 12, "A/H/C/D/R  drop slot", 8, 0)
	u.putStr(16, 13, "j/k select  enter use  e/z/space close", 8, 0)
	u.putStr(16, 14, g.RiteSummary(), 11, 0)

	ground := g.ItemAt(g.Player.X, g.Player.Y)
	u.putStr(1, 17, "[ g ] item here: "+itemSummary(ground), 10, 0)

	u.fill(0, 22, screenW, 2, ' ', 8, 0)
	u.putStr(1, 22, "Robe slots hold vials, shards and small items; the robe cannot be removed.", 8, 0)
	u.putStr(1, 23, "g get   e/z/space close", 8, 0)
}

func slotStats(it *Item) string {
	if it == nil {
		return ""
	}
	d := it.Def()
	if d == nil {
		return ""
	}
	switch d.Kind {
	case KindMelee, KindGrenade:
		if wd := d.Weapon(); wd != nil {
			return fmt.Sprintf("dmg %d-%d av %d dv %d", wd.DmgMin, wd.DmgMax, wd.AV, wd.DV)
		}
	case KindArmor:
		ensureArmorAP(it)
		if stats, ok := ArmorStatsFor(d); ok {
			return fmt.Sprintf("AP %d/%d abs %d%% dv %d", it.AP, stats.MaxAP, stats.Absorb, stats.DV)
		}
	case KindMedical:
		if it.Proto == ArmourVialProto {
			return fmt.Sprintf("AP +10  charges %d", it.Charges)
		}
		return fmt.Sprintf("heal %d  charges %d", d.Heal, it.Charges)
	case KindFood:
		return fmt.Sprintf("heal %d  charges %d", d.Heal, it.Charges)
	case KindShard:
		ensureShardAP(it)
		return fmt.Sprintf("AP %d/%d", it.AP, shardMaxAP(it))
	case KindTool:
		return d.Desc
	case KindRelic:
		return "relic"
	case KindHead:
		return "head item"
	}
	return d.Name
}

func slotAction(slot int, it *Item) string {
	if it == nil {
		return "-"
	}
	if slot == SlotHands {
		return "act a"
	}
	d := it.Def()
	if d == nil {
		return "-"
	}
	switch d.Kind {
	case KindMelee, KindGrenade, KindTool:
		return "use " + slotKey(slot)
	case KindMedical, KindFood:
		return "consume " + slotKey(slot)
	case KindShard:
		return "inspect"
	case KindRelic, KindHead:
		return "inspect"
	}
	return "-"
}

func (u *UI) handleInventory(e *tcell.EventKey) {
	g := u.g
	slots := g.EquipmentSlots()
	if len(slots) == 0 {
		g.State = StatePlaying
		return
	}
	if u.slotSel >= len(slots) {
		u.slotSel = len(slots) - 1
	}
	switch e.Key() {
	case tcell.KeyEscape:
		g.State = StatePlaying
		return
	case tcell.KeyEnter:
		g.UseSlot(slots[u.slotSel])
		return
	}
	switch e.Rune() {
	case 'j':
		u.slotSel = (u.slotSel + 1) % len(slots)
	case 'k':
		u.slotSel = (u.slotSel + len(slots) - 1) % len(slots)
	case 'a':
		g.UseSlot(SlotHands)
	case 'h':
		g.UseSlot(SlotHead)
	case 'c':
		g.UseSlot(SlotRobe1)
	case 'd':
		g.UseSlot(SlotRobe2)
	case 'r':
		g.UseSlot(SlotRobe3)
	case 'A':
		g.DropSlot(SlotHands)
	case 'H':
		g.DropSlot(SlotHead)
	case 'C':
		g.DropSlot(SlotRobe1)
	case 'D':
		g.DropSlot(SlotRobe2)
	case 'R':
		g.DropSlot(SlotRobe3)
	case 'g':
		g.PickUp()
	case 'z', 'q', 'e', ' ':
		g.State = StatePlaying
	}
}

func (g *Game) playerArmorLine() string {
	cur, max := g.PlayerAP()
	return fmt.Sprintf("armour AP: %d/%d", cur, max)
}

func (u *UI) drawPlayerScreen() {
	g := u.g
	p := g.Player
	u.clear(0)
	u.putStr(centerText("[ @ ] character information", screenW), 0, "[ @ ] character information", 2, 0)
	u.putStr(1, 2, strings.Repeat("~", 76), 2, 0)
	u.putStr(2, 4, "name: "+p.Name, 14, 0)
	u.putStr(2, 5, fmt.Sprintf("health: %d/%d", p.HP, p.MaxHP), 11, 0)
	u.putStr(2, 6, g.playerArmorLine(), 11, 0)
	u.putStr(2, 7, fmt.Sprintf("experience: %d (stored)", p.XP), 11, 0)
	u.putStr(2, 8, fmt.Sprintf("gold: %d", p.Gold), 11, 0)
	u.putStr(2, 10, "[ rites ]", 10, 0)
	row := 11
	for i, rite := range g.PlayerRites {
		if rite < 0 || rite >= len(RiteDefs) {
			continue
		}
		u.putStr(3, row, fmt.Sprintf("%-16s %d charge", riteName(rite), g.RiteCharges[i]), 11, 0)
		row++
	}
	if row == 11 {
		u.putStr(3, row, "none", 8, 0)
		row++
	}
	row++
	u.putStr(2, row, "[ equipment ]", 10, 0)
	row++
	for _, slot := range g.EquipmentSlots() {
		it := p.SlotItem(slot)
		name := "-"
		if it != nil {
			name = it.Name()
		}
		u.putStr(3, row, fmt.Sprintf("%-8s %s", slotName(slot), name), 11, 0)
		row++
	}
	u.fill(0, 22, screenW, 3, ' ', 8, 0)
	u.putStr(1, 23, "z/space quit screen", 8, 0)
	u.putStr(1, 24, "e equipment", 8, 0)
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
		case 'z', 'q', '@', ' ':
			u.g.State = StatePlaying
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
		y := 4 + i*2
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
	u.fill(0, 22, screenW, 3, ' ', 8, 0)
	u.putStr(1, 23, "up/down select   enter/+ increase skill (1 point)   z quit", 8, 0)
	u.putStr(1, 24, fmt.Sprintf("you currently have %d free skill points", g.FreePoints), 10, 0)
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
		case 'z', 'q', 'x', ' ':
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
	visible := 17
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
	u.putStr(1, 20, "[ ? ] description...", 10, 0)
	u.putStr(1, 21, strings.Repeat("~", 76), 2, 0)
	for i, l := range wrapText(d.Desc, 76) {
		if i > 1 {
			break
		}
		u.putStr(1, 22+i, l, 11, 0)
	}
	u.fill(0, 24, screenW, 1, ' ', 8, 0)
	u.putStr(1, 24, "z quit   <6> increase   <4> decrease   enter toggle", 8, 0)
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
		case 'z', 'q', 'O', ' ':
			g.State = StatePlaying
		case '4', 'h':
			g.ChangeOption(opts[u.optSel].Name, -1)
		case '6', 'l':
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
	visible := 17
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
	u.putStr(1, 21, "[ ? ] description...", 10, 0)
	u.putStr(1, 22, strings.Repeat("~", 76), 2, 0)
	u.putStr(1, 23, k.Desc, 11, 0)
	u.fill(0, 24, screenW, 1, ' ', 8, 0)
	u.putStr(1, 24, "z quit   up/down select   pgup/pgdn page", 8, 0)
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
		case 'z', 'q', '?', '/', ' ':
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
	visible := 19
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
	u.fill(0, 24, screenW, 1, ' ', 8, 0)
	u.putStr(1, 24, "z quit   up/down scroll   pgup/pgdn page", 8, 0)
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
		case 'j':
			u.msgScroll++
		case 'k':
			u.msgScroll--
		case 'z', 'q', '~', '`', ' ':
			u.g.State = StatePlaying
		}
	}
}
