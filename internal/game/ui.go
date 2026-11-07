package game

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	screenW  = mapWidth
	screenH  = 25
	mapViewH = 23
	hudRow   = screenH - 1
)

var cp437 = [128]rune{
	'Ç', 'ü', 'é', 'â', 'ä', 'à', 'å', 'ç', 'ê', 'ë', 'è', 'ï', 'î', 'ì', 'Ä', 'Å',
	'É', 'æ', 'Æ', 'ô', 'ö', 'ò', 'û', 'ù', 'ÿ', 'Ö', 'Ü', '¢', '£', '¥', '₧', 'ƒ',
	'á', 'í', 'ó', 'ú', 'ñ', 'Ñ', 'ª', 'º', '¿', '⌐', '¬', '½', '¼', '¡', '«', '»',
	'░', '▒', '▓', '│', '┤', '╡', '╢', '╖', '╕', '╣', '║', '╗', '╝', '╜', '╛', '┐',
	'└', '┴', '┬', '├', '─', '┼', '╞', '╟', '╚', '╔', '╩', '╦', '╠', '═', '╬', '╧',
	'╨', '╤', '╥', '╙', '╘', '╒', '╓', '╫', '╪', '┘', '┌', '█', '▄', '▌', '▐', '▀',
	'α', 'ß', 'Γ', 'π', 'Σ', 'σ', 'µ', 'τ', 'Φ', 'Θ', 'Ω', 'δ', '∞', 'φ', 'ε', '∩',
	'≡', '±', '≥', '≤', '⌠', '⌡', '÷', '≈', '°', '∙', '·', '√', 'ⁿ', '²', '■', '\u00a0',
}

func cp437Rune(b byte) rune {
	if b < 128 {
		return rune(b)
	}
	return cp437[b-128]
}

type UI struct {
	screen    tcell.Screen
	g         *Game
	ox, oy    int
	colorMode ColorMode

	titleSel    int
	classSel    int
	optSel      int
	optScroll   int
	helpSel     int
	helpScroll  int
	msgScroll   int
	msgExpanded bool
	bashPending bool
	skillSel    int
	slotSel     int
}

func Run() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}
	defer s.Fini()
	mode := DetectColorMode(s.Colors())
	s.SetStyle(tcell.StyleDefault.Background(screenBackgroundColor(mode)).Foreground(tcell.ColorWhite))
	s.Clear()

	u := &UI{
		screen:    s,
		g:         NewGame(time.Now().UnixNano()),
		colorMode: mode,
	}
	if len(os.Args) > 1 {
		u.g.PlayerName = os.Args[1]
	}
	u.g.State = StateTitle
	if _, err := os.Stat(u.g.SavePath); err == nil {
		u.titleSel = 1
	}
	for !u.g.Quit {
		sw, sh := s.Size()
		u.ox = (sw - screenW) / 2
		u.oy = (sh - screenH) / 2
		if u.ox < 0 {
			u.ox = 0
		}
		if u.oy < 0 {
			u.oy = 0
		}
		u.draw()
		s.Show()
		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			u.handleKey(e)
		}
	}
	return nil
}

func (u *UI) put(x, y int, r rune, fg, bg uint8) {
	u.putTcell(x, y, r, u.uiDOSColor(fg, false), u.uiDOSColor(bg, true))
}

func (u *UI) putTcell(x, y int, r rune, fg, bg tcell.Color) {
	if x < 0 || y < 0 || x >= screenW || y >= screenH {
		return
	}
	st := tcell.StyleDefault.Foreground(fg).Background(bg)
	u.screen.SetContent(u.ox+x, u.oy+y, r, nil, st)
}

func (u *UI) putStr(x, y int, s string, fg, bg uint8) {
	for i, r := range []rune(s) {
		u.put(x+i, y, r, fg, bg)
	}
}

func (u *UI) putStrTcell(x, y int, s string, fg, bg tcell.Color) {
	for i, r := range []rune(s) {
		u.putTcell(x+i, y, r, fg, bg)
	}
}

func (u *UI) fill(x, y, w, h int, r rune, fg, bg uint8) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			u.put(xx, yy, r, fg, bg)
		}
	}
}

func (u *UI) fillTcell(x, y, w, h int, r rune, fg, bg tcell.Color) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			u.putTcell(xx, yy, r, fg, bg)
		}
	}
}

// uiColor resolves a UI palette role, keeping the supplied DOS fallback for
// 16-colour terminals.
func (u *UI) uiColor(role ColorRole, fallback uint8) tcell.Color {
	if u.colorMode == ColorMode16 {
		return pal(fallback)
	}
	pc, ok := WoodsPalette.Lookup(role)
	if !ok {
		return pal(fallback)
	}
	return pc.Resolve(u.colorMode)
}

func uiForegroundRole(i uint8) (ColorRole, bool) {
	switch i {
	case 0:
		return RoleSelectionText, true
	case 1, 3, 9, 11:
		return RoleHintCyan, true
	case 2, 10:
		return RoleAccentGreen, true
	case 4:
		return RoleBloodAccent, true
	case 5, 13:
		return RolePotionPurple, true
	case 6:
		return RoleLeather, true
	case 7:
		return RolePrimaryText, true
	case 8:
		return RoleDimText, true
	case 12:
		return RoleDanger, true
	case 14:
		return RoleAccentGold, true
	case 15:
		return RoleSelectionText, true
	}
	return "", false
}

func uiBackgroundRole(i uint8) (ColorRole, bool) {
	switch i {
	case 0:
		return RoleUIBackground, true
	case 1, 3, 5, 6, 8, 9, 10, 11, 12, 13, 14, 15:
		return RolePanelBackground, true
	case 2:
		return RoleSelectionBackground, true
	case 4:
		return RoleDarkBlood, true
	case 7:
		return RoleSelectionBackground, true
	}
	return "", false
}

// uiDOSColor maps the old DOS colour arguments used by UI code onto semantic
// palette roles. 16-colour mode keeps the recovered DOS palette exactly.
func (u *UI) uiDOSColor(i uint8, background bool) tcell.Color {
	if u.colorMode == ColorMode16 {
		return pal(i)
	}
	var role ColorRole
	var ok bool
	if background {
		role, ok = uiBackgroundRole(i)
	} else {
		role, ok = uiForegroundRole(i)
	}
	if !ok {
		return pal(i)
	}
	pc, ok := WoodsPalette.Lookup(role)
	if !ok {
		return pal(i)
	}
	return pc.Resolve(u.colorMode)
}

func (u *UI) clear(bg uint8) {
	u.fill(0, 0, screenW, screenH, ' ', 7, bg)
}

func (u *UI) woodsColorsActive() bool {
	return u.g != nil && u.g.Theme == 0
}

func (u *UI) paletteFallback() bool {
	return u.colorMode == ColorMode16 || !u.woodsColorsActive()
}

// worldColor resolves a foreground colour. Remembered tiles are dimmed and
// slightly desaturated instead of being replaced by a flat grey.
func (u *UI) worldColor(role ColorRole, ok bool, fallback uint8, x, y int, visible bool) tcell.Color {
	if !ok || !u.woodsColorsActive() {
		if !visible {
			return pal(8)
		}
		return pal(fallback)
	}
	pc, found := WoodsPalette.Lookup(role)
	if !found {
		if !visible {
			return pal(8)
		}
		return pal(fallback)
	}
	if !visible {
		if u.colorMode == ColorMode16 {
			return pal(8)
		}
		c := u.applyWoodsAccent(pc.RGB, role, x, y)
		return ResolveRGB(ShadeRGB(c, 0.40, 0.15), fallback, u.colorMode)
	}
	c := VaryRGB(pc.RGB, u.g.Level, x, y)
	c = u.applyWoodsAccent(c, role, x, y)
	c = u.applyWarmLight(c, x, y)
	return ResolveRGB(c, fallback, u.colorMode)
}

// worldBackgroundColor resolves a background colour. It keeps the supplied
// fallback index instead of the remembered-tile grey used for glyphs.
func (u *UI) worldBackgroundColor(role ColorRole, ok bool, fallback uint8, x, y int, visible bool) tcell.Color {
	if !ok {
		return pal(fallback)
	}
	pc, found := WoodsPalette.Lookup(role)
	if !found {
		return pal(fallback)
	}
	if !visible {
		if u.colorMode == ColorMode16 {
			return pal(fallback)
		}
		c := u.applyWoodsAccent(pc.RGB, role, x, y)
		return ResolveRGB(ShadeRGB(c, 0.22, 0.15), fallback, u.colorMode)
	}
	c := VaryRGB(pc.RGB, u.g.Level, x, y)
	c = u.applyWoodsAccent(c, role, x, y)
	// Keep terrain backgrounds very dark so the tile glyphs read clearly.
	c = ShadeRGB(c, 0.30, 0.12)
	c = u.applyWarmLight(c, x, y)
	return ResolveRGB(c, fallback, u.colorMode)
}

func (u *UI) bloodColor(tileRole ColorRole, tileOK bool, fallback uint8, x, y int) tcell.Color {
	if !u.woodsColorsActive() || u.colorMode == ColorMode16 {
		return pal(4)
	}
	base := WoodsPalette[RoleGrassBase].RGB
	if tileOK {
		if pc, ok := WoodsPalette.Lookup(tileRole); ok {
			base = pc.RGB
		}
	}
	base = VaryRGB(base, u.g.Level, x, y)
	base = u.applyWoodsAccent(base, tileRole, x, y)
	blood := WoodsPalette[RoleDarkBlood].RGB
	return ResolveRGB(BlendRGB(base, blood, 0.55), fallback, u.colorMode)
}

func centerText(s string, width int) int {
	n := len([]rune(s))
	if n >= width {
		return 0
	}
	return (width - n) / 2
}

func (u *UI) draw() {
	switch u.g.State {
	case StateTitle:
		u.drawTitle()
	case StateIntro:
		u.drawIntro()
	case StatePlaying:
		u.drawGame()
	case StateInventory:
		u.drawGame()
		u.drawInventory()
	case StatePlayerScreen:
		u.drawPlayerScreen()
	case StateSkills:
		u.drawSkills()
	case StateOptions:
		u.drawOptions()
	case StateHelp:
		u.drawHelp()
	case StateMessages:
		u.drawMessageBuffer()
	case StateLook:
		u.drawGame()
		u.drawLookOverlay()
	case StateThrow:
		u.drawGame()
		u.drawThrowOverlay()
	case StateQuitConfirm:
		u.drawGame()
		u.drawConfirm()
	case StateGameOver:
		u.drawGameOver(false)
	case StateVictory:
		u.drawGameOver(true)
	}
}

func (u *UI) drawTitle() {
	u.clear(0)
	left, top, width, height := 1, 1, 78, 22
	border := uint8(2)
	for x := left + 1; x < left+width-1; x++ {
		u.put(x, top, '-', border, 0)
		u.put(x, top+height-1, '-', border, 0)
	}
	for y := top + 1; y < top+height-1; y++ {
		u.put(left, y, '|', border, 0)
		u.put(left+width-1, y, '|', border, 0)
	}
	u.put(left, top, '+', border, 0)
	u.put(left+width-1, top, '+', border, 0)
	u.put(left, top+height-1, '+', border, 0)
	u.put(left+width-1, top+height-1, '+', border, 0)

	u.putStr(left+3, top+1, "CASE FILE: HULK INFECTION", 14, 0)
	status := "STATUS: CONTAINMENT FAILED"
	u.putStr(left+width-3-len(status), top+1, status, 12, 0)
	u.putStr(left+3, top+2, "LOCATION: NORTH WOODS", 10, 0)
	date := "DATE: 2026"
	u.putStr(left+width-3-len(date), top+2, date, 10, 0)
	u.putStr(left+1, top+3, strings.Repeat("-", width-2), 2, 0)

	title := "H  U  L  K     I  N  F  E  C  T  I  O  N"
	u.putStr(left+(width-len(title))/2, top+5, title, 10, 0)
	rule := strings.Repeat("=", len(title))
	u.putStr(left+(width-len(rule))/2, top+6, rule, 2, 0)
	warning := "The dead are walking. The road is closed."
	u.putStr(left+(width-len(warning))/2, top+8, warning, 7, 0)

	menu := []string{"new game", "continue", "quit"}
	for i, s := range menu {
		prefix := "  "
		fg := uint8(7)
		if i == u.titleSel {
			prefix = "> "
			fg = 14
		}
		u.putStr(left+5, top+11+i*2, prefix+s, fg, 0)
	}

	copyright := "copyright (c) 2026 Michał Przybyłowicz"
	u.putStr(left+(width-len([]rune(copyright)))/2, top+height-2, copyright, 8, 0)
	u.putStr(centerText("arrows / enter", screenW), top+height, "arrows / enter", 8, 0)
}

func (u *UI) drawIntro() {
	u.clear(0)
	g := u.g
	rite := RiteDefs[0]
	if g.PlayerRites[0] >= 0 && g.PlayerRites[0] < len(RiteDefs) {
		rite = RiteDefs[g.PlayerRites[0]]
	}
	lines := []struct {
		text  string
		row   int
		color uint8
	}{
		{"T H E   W O O D S   A R E   D E A D", 2, 12},
		{g.PlayerName, 5, 14},
		{"The bishop assigned you to a village that no longer answers letters.", 7, 10},
		{"The dead are walking. Find the source and put it to rest.", 9, 14},
		{"Rite: " + rite.Name + " - " + rite.Desc, 12, 11},
		{"Hands: sickle   Robe: 3 slots   Shard: leather   Vial: healing", 14, 10},
		{"The first body rose seven days ago. By morning the town was gone.", 17, 10},
		{"Somewhere in these woods is the source of the infection.", 18, 10},
		{"press enter to begin", 21, 8},
	}
	for _, line := range lines {
		u.putStr(centerText(line.text, screenW), line.row, line.text, line.color, 0)
	}
}

func (u *UI) drawCharCreate() {
	u.clear(0)
	u.putStr(centerText("[ @ ] profession...", screenW), 0, "[ @ ] profession...", 11, 0)
	u.putStr(2, 2, "select class", 8, 0)
	for i, c := range PlayerClasses {
		mark := "  "
		fg := uint8(7)
		if i == u.classSel {
			mark = "> "
			fg = 15
		}
		u.putStr(4, 4+i, mark+c.Name, fg, 0)
	}
	c := &PlayerClasses[u.classSel]
	u.putStr(2, 10, "class: "+c.Name, 14, 0)
	descLines := wrapText(c.Desc, 74)
	for i, l := range descLines {
		if i > 2 {
			break
		}
		u.putStr(2, 11+i, l, 10, 0)
	}
	u.putStr(2, 15, fmt.Sprintf("strength:%d  toughness:%d  agility:%d  intelligence:%d", c.STR, c.TGH, c.AGI, c.INT), 7, 0)
	for i, skill := range ActiveSkills {
		col, row := i/3, i%3
		x := 2 + col*26
		y := 17 + row
		u.putStr(x, y, fmt.Sprintf("%-16s %3d", SkillNames[skill], c.Skills[skill]), 10, 0)
	}
	u.putStr(2, 21, startingKitSummary(u.classSel), 8, 0)
	u.putStr(2, 22, "name: "+u.g.PlayerName, 14, 0)
	u.putStr(2, 23, "arrows: select class   r: reroll name   enter: begin", 8, 0)
	u.putStr(2, 24, "esc: back", 8, 0)
}

func (u *UI) drawGame() {
	g := u.g
	if g.Map == nil || g.Player == nil {
		u.clear(0)
		return
	}
	u.clear(0)
	left, top := 0, 0
	if g.Map.W > screenW {
		left = clampInt(g.Player.X-screenW/2, 0, g.Map.W-screenW)
	}
	if g.Map.H > mapViewH {
		top = clampInt(g.Player.Y-mapViewH/2, 0, g.Map.H-mapViewH)
	}
	for sy := 0; sy < mapViewH; sy++ {
		for sx := 0; sx < screenW; sx++ {
			mx, my := left+sx, top+sy
			mapBG := u.worldColor(RoleMapBackground, true, 0, mx, my, true)
			if !g.Map.InBounds(mx, my) {
				u.putTcell(sx, sy+1, ' ', mapBG, mapBG)
				continue
			}
			i := g.Map.idx(mx, my)
			if !g.Map.Seen[i] {
				u.putTcell(sx, sy+1, ' ', mapBG, mapBG)
				continue
			}
			tileID := int(g.Map.Tile(mx, my))
			d := g.Map.Def(mx, my)
			vis := g.Map.Visible[i] || g.Options.Bool("optMapDrawHiddenTiles")
			tileRole, tileRoleOK := TileColorRole(tileID)
			fg := d.Color
			fgRole, fgRoleOK := tileRole, tileRoleOK
			fgVisible := vis
			if !vis {
				fg = 8
			}
			r := tileDisplayRune(tileID, d)
			if g.Options.Bool("optMapDrawPortals") && (d.Name == "time portal" || d.Name == "stone stairs" || d.Name == "elevator") {
				fg = 14
				fgRole, fgRoleOK = RoleGoldRelic, true
				fgVisible = true
			}
			bgRole, bgRoleOK := TileBackgroundRole(tileID)
			bg := d.Color
			if pc, found := WoodsPalette.Lookup(bgRole); found {
				bg = pc.Fallback
			}
			blooded := g.Map.Blood[i] > 0 && d.Type == 1
			if blooded {
				bgRole, bgRoleOK = RoleDarkBlood, true
				if pc, found := WoodsPalette.Lookup(bgRole); found {
					bg = pc.Fallback
				}
			}
			if g.Map.Fire[i] > 0 {
				r = '^'
				fg = 12
				fgRole, fgRoleOK = RoleFlame, true
				fgVisible = true
			} else if g.Map.Acid[i] > 0 {
				r = '~'
				fg = 10
				fgRole, fgRoleOK = RoleAcid, true
				fgVisible = true
			} else if g.Map.Smoke[i] > 0 && vis {
				r = '%'
				fg = 8
				fgRole, fgRoleOK = RoleSmoke, true
			}
			if vis {
				if it := g.ItemAt(mx, my); it != nil {
					if d2 := it.Def(); d2 != nil {
						r = cp437Rune(d2.Glyph)
						fg = d2.Color
						fgRole, fgRoleOK = ItemColorRole(it.Proto)
						if blooded && u.paletteFallback() {
							fg = 15 // keep corpses/items visible on blood
							fgRoleOK = false
						}
					}
				}
				if c := g.CritterAt(mx, my); c != nil && c.IsAlive() {
					if g.Map.Visible[i] || g.Options.Bool("optDrawDebugAI") {
						r = cp437Rune(c.Glyph)
						fg = c.Color
						fgRole, fgRoleOK = CritterColorRole(c.Proto)
						if c.Blood > 0 {
							fg = 12
							fgRole, fgRoleOK = RoleFreshBlood, true
							if blooded && u.paletteFallback() {
								fg = 15
								fgRoleOK = false
							}
						}
					}
				}
			}
			if mx == g.Player.X && my == g.Player.Y {
				r = '@'
				fg = 15
				fgRole, fgRoleOK = RolePlayer, true
				fgVisible = true
			}
			fgColor := u.worldColor(fgRole, fgRoleOK, fg, mx, my, fgVisible)
			var bgColor tcell.Color
			if blooded && vis {
				bgColor = u.bloodColor(tileRole, tileRoleOK, 4, mx, my)
			} else {
				bgColor = u.worldBackgroundColor(bgRole, bgRoleOK, bg, mx, my, vis)
			}
			u.putTcell(sx, sy+1, r, fgColor, bgColor)
		}
	}
	// Highlight the cursor tile in look and throw modes.
	hx, hy, hasHighlight := 0, 0, false
	switch g.State {
	case StateLook, StateThrow:
		hx, hy, hasHighlight = g.LookX, g.LookY, true
	}
	if hasHighlight {
		sx, sy := hx-left, hy-top
		if sx >= 0 && sy >= 0 && sx < screenW && sy < mapViewH {
			r, _ := g.CellForHighlight(hx, hy)
			u.put(sx, sy+1, r, 0, 7)
		}
	}
	u.drawTopMessages()
	u.drawStatus(hudRow)
	if u.bashPending {
		u.fill(0, 22, screenW, 2, ' ', 8, 0)
		u.fill(0, 22, screenW, 1, ' ', 15, 1)
		u.putStr(1, 22, "bash: choose a direction", 15, 1)
		u.putStr(1, 23, "h/j/k/l or arrows bash, esc cancels", 8, 0)
	}
}

// tileDisplayRune applies the Brogue-style presentation rules: plain ground is
// blank, walls use '#', and water uses '='.
func tileDisplayRune(tile int, d *TileDef) rune {
	if d == nil {
		return '?'
	}
	if TileIsGround(tile) {
		return ' '
	}
	if strings.Contains(strings.ToLower(d.Name), "wall") {
		return '#'
	}
	if role, ok := TileColorRole(tile); ok {
		switch role {
		case RoleDeepWater, RoleShallowWater, RoleRiver, RoleMarsh:
			return '='
		}
	}
	return cp437Rune(d.Glyph)
}

// CellForHighlight returns the glyph/colour of a world cell for the look/throw
// cursor.  Unexplored cells are shown as '?'.
func (g *Game) CellForHighlight(x, y int) (rune, uint8) {
	if g.Map == nil || !g.Map.InBounds(x, y) {
		return '?', 8
	}
	i := g.Map.idx(x, y)
	if !g.Map.Seen[i] {
		return '?', 8
	}
	d := g.Map.Def(x, y)
	r, fg := tileDisplayRune(int(g.Map.Tile(x, y)), d), d.Color
	if !g.Map.Visible[i] && !g.Options.Bool("optMapDrawHiddenTiles") {
		fg = 8
	}
	if it := g.ItemAt(x, y); it != nil {
		if d2 := it.Def(); d2 != nil {
			r, fg = cp437Rune(d2.Glyph), d2.Color
		}
	}
	if c := g.CritterAt(x, y); c != nil && c.IsAlive() {
		if g.Map.Visible[i] || g.Options.Bool("optDrawDebugAI") {
			r, fg = cp437Rune(c.Glyph), c.Color
		}
	}
	if x == g.Player.X && y == g.Player.Y {
		r, fg = '@', 15
	}
	return r, fg
}

func (u *UI) drawTopMessages() {
	g := u.g
	if u.msgExpanded {
		u.fill(0, 0, screenW, 8, ' ', 7, 0)
		u.drawMessageLines(0, 8)
		return
	}
	u.fill(0, 0, screenW, 1, ' ', 7, 0)
	if len(g.Messages) == 0 {
		return
	}
	m := g.Messages[len(g.Messages)-1]
	line := m.Text
	if len(line) > screenW {
		line = line[:screenW]
	}
	u.putStr(0, 0, line, m.Color, 0)
}

func (u *UI) drawMessageLines(y, rows int) {
	g := u.g
	type line struct {
		text  string
		color uint8
	}
	var lines []line
	for i := len(g.Messages) - 1; i >= 0 && len(lines) < rows; i-- {
		m := g.Messages[i]
		wrapped := wrapText(m.Text, screenW)
		for j := len(wrapped) - 1; j >= 0 && len(lines) < rows; j-- {
			lines = append([]line{{wrapped[j], m.Color}}, lines...)
		}
	}
	for i := 0; i < rows; i++ {
		if i < len(lines) {
			u.putStr(0, y+i, lines[i].text, lines[i].color, 0)
		}
	}
}

func slotBrief(it *Item) string {
	if it == nil {
		return "-"
	}
	return it.Name()
}

func (u *UI) drawStatus(y int) {
	g := u.g
	p := g.Player
	ac, am := g.PlayerAP()
	hpColor := uint8(10)
	if p.MaxHP > 0 {
		ratio := float64(p.HP) / float64(p.MaxHP)
		switch {
		case ratio <= 0.25:
			hpColor = 12
		case ratio <= 0.60:
			hpColor = 14
		}
	}
	apColor := uint8(11)
	if ac == 0 && am > 0 {
		apColor = 12
	}
	weapon := "none"
	if it := p.SlotItem(SlotHands); it != nil {
		weapon = it.Name()
	}
	hp := fmt.Sprintf("HP %d/%d", p.HP, p.MaxHP)
	ap := fmt.Sprintf("AP %d/%d", ac, am)
	u.putStr(0, y, hp, hpColor, 0)
	u.putStr(len(hp)+3, y, ap, apColor, 0)
	weaponX := len(hp) + len(ap) + 6
	u.putStr(weaponX, y, weapon, 7, 0)
	riteX := weaponX + len(weapon) + 3
	u.putStr(riteX, y, g.RiteSummary(), 11, 0)
	if target := g.CurrentTarget; target != nil && target.IsAlive() {
		text := fmt.Sprintf("%c %s %s HP %d/%d", cp437Rune(target.Glyph), target.NameOrUnknown(), target.DamageDescription(), target.HP, target.MaxHP)
		x := screenW - len([]rune(text))
		if x < riteX+len([]rune(g.RiteSummary()))+2 {
			x = riteX + len([]rune(g.RiteSummary())) + 2
		}
		u.putStr(x, y, text, 12, 0)
	}
}

func wrapText(s string, width int) []string {
	if width < 1 {
		return []string{s}
	}
	var out []string
	for len(s) > width {
		cut := strings.LastIndexByte(s[:width+1], ' ')
		if cut <= 0 {
			cut = width
		}
		out = append(out, s[:cut])
		s = strings.TrimLeft(s[cut:], " ")
	}
	if s != "" {
		out = append(out, s)
	}
	if len(out) == 0 {
		out = []string{""}
	}
	return out
}

func (u *UI) handleKey(e *tcell.EventKey) {
	g := u.g
	if e.Key() == tcell.KeyCtrlL {
		u.screen.Sync()
		return
	}
	switch g.State {
	case StateTitle:
		u.handleTitle(e)
	case StateIntro:
		u.handleIntro(e)
	case StatePlaying:
		u.handlePlaying(e)
	case StateInventory:
		u.handleInventory(e)
	case StatePlayerScreen:
		u.handlePlayerScreen(e)
	case StateSkills:
		u.handleSkills(e)
	case StateOptions:
		u.handleOptions(e)
	case StateHelp:
		u.handleHelp(e)
	case StateMessages:
		u.handleMessages(e)
	case StateLook:
		u.handleLook(e)
	case StateThrow:
		u.handleThrow(e)
	case StateQuitConfirm:
		u.handleQuitConfirm(e)
	case StateGameOver:
		u.handleGameOver(e)
	case StateVictory:
		u.handleGameOver(e)
	}
}

func isEnter(e *tcell.EventKey) bool {
	return e.Key() == tcell.KeyEnter || e.Rune() == '\n' || e.Rune() == '\r'
}

func (u *UI) handleIntro(e *tcell.EventKey) {
	if e.Key() == tcell.KeyEscape || e.Rune() == 'q' {
		u.g.State = StateTitle
		return
	}
	if e.Rune() == 'r' {
		// Development reroll: regenerate the priest, kit and world.
		u.g = NewGame(time.Now().UnixNano())
		u.g.State = StateIntro
		return
	}
	if !isEnter(e) {
		return
	}
	u.g.SetupNewGame(0, u.g.PlayerName)
	u.g.State = StatePlaying
	u.g.Msg("=== session started ===", 8)
	u.g.UpdateFOV()
}

func (u *UI) handleTitle(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyUp:
		u.titleSel--
		if u.titleSel < 0 {
			u.titleSel = 2
		}
	case tcell.KeyDown:
		u.titleSel++
		if u.titleSel > 2 {
			u.titleSel = 0
		}
	default:
		switch e.Rune() {
		case 'k':
			u.titleSel--
			if u.titleSel < 0 {
				u.titleSel = 2
			}
			return
		case 'j':
			u.titleSel++
			if u.titleSel > 2 {
				u.titleSel = 0
			}
			return
		}
		if e.Rune() == 'q' || e.Key() == tcell.KeyEscape {
			u.g.Quit = true
			return
		}
		if !isEnter(e) {
			return
		}
		switch u.titleSel {
		case 0:
			u.g = NewGame(time.Now().UnixNano())
			u.g.State = StateIntro
		case 1:
			if err := u.g.LoadGame(u.g.SavePath); err != nil {
				u.g.Msg("could not load game: "+err.Error(), 12)
				u.g.State = StateTitle
			}
		case 2:
			u.g.Quit = true
		}
	}
}

func (u *UI) handleCharCreate(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyUp:
		u.classSel = (u.classSel + len(PlayerClasses) - 1) % len(PlayerClasses)
	case tcell.KeyDown:
		u.classSel = (u.classSel + 1) % len(PlayerClasses)
	case tcell.KeyEscape:
		u.g.State = StateIntro
	default:
		switch e.Rune() {
		case 'k':
			u.classSel = (u.classSel + len(PlayerClasses) - 1) % len(PlayerClasses)
			return
		case 'j':
			u.classSel = (u.classSel + 1) % len(PlayerClasses)
			return
		}
		if e.Rune() == 'r' {
			u.g.PlayerName = RandomNames[u.g.Roll(len(RandomNames))]
		}
		if isEnter(e) {
			if u.g.PlayerName == "" {
				u.g.PlayerName = RandomNames[u.g.Roll(len(RandomNames))]
			}
			u.g.SetupNewGame(u.classSel, u.g.PlayerName)
			u.g.State = StatePlaying
			u.g.Msg("=== session started ===", 8)
			u.g.UpdateFOV()
		}
	}
}

func directionFromEvent(e *tcell.EventKey) (int, int, bool) {
	switch e.Key() {
	case tcell.KeyUp:
		return 0, -1, true
	case tcell.KeyDown:
		return 0, 1, true
	case tcell.KeyLeft:
		return -1, 0, true
	case tcell.KeyRight:
		return 1, 0, true
	}
	switch e.Rune() {
	case 'h':
		return -1, 0, true
	case 'j':
		return 0, 1, true
	case 'k':
		return 0, -1, true
	case 'l':
		return 1, 0, true
	}
	return 0, 0, false
}

func (u *UI) handlePlaying(e *tcell.EventKey) {
	g := u.g
	switch e.Key() {
	case tcell.KeyF1:
		u.debugSpawnCritter()
		return
	case tcell.KeyF2:
		u.debugSpawnItem()
		return
	case tcell.KeyF3:
		g.Player.HP = g.Player.MaxHP
		g.Msg("debug: health restored.", 10)
		return
	case tcell.KeyF4:
		for i := range g.Map.Seen {
			g.Map.Seen[i] = true
		}
		g.Msg("debug: map revealed.", 10)
		return
	case tcell.KeyF5:
		g.Player.XP += 200
		g.Msg("debug: experience granted.", 10)
		return
	case tcell.KeyF7:
		g.ChangeLevelRelative(1)
		return
	case tcell.KeyF8:
		g.ChangeLevelRelative(-1)
		return
	case tcell.KeyF6:
		for y := 0; y < g.Map.H; y++ {
			for x := 0; x < g.Map.W; x++ {
				if d := g.Map.Def(x, y); d.Name == "stone stairs" || d.Name == "time portal" || d.Name == "elevator" {
					g.Player.X, g.Player.Y = x, y
					g.UpdateFOV()
					g.Msg("debug: teleported to portal.", 10)
					return
				}
			}
		}
		return
	}
	if u.msgExpanded {
		if e.Key() == tcell.KeyEscape || e.Rune() == '~' {
			u.msgExpanded = false
		}
		return
	}
	if u.bashPending {
		if e.Key() == tcell.KeyEscape {
			u.bashPending = false
			g.Msg("bash cancelled.", 8)
			return
		}
		if dx, dy, ok := directionFromEvent(e); ok {
			u.bashPending = false
			g.PlayerBashDirection(dx, dy)
			return
		}
		u.bashPending = false
		g.Msg("bash cancelled.", 8)
		return
	}
	if dx, dy, ok := directionFromEvent(e); ok {
		g.MovePlayer(dx, dy)
		return
	}
	r := e.Rune()
	if slot, ok := slotForKey(r); ok {
		g.UseSlot(slot)
		return
	}
	switch r {
	case '.':
		g.WaitTurn()
	case 'r':
		g.UseRite(0)
	case 't':
		g.UseRite(1)
	case 'p':
		g.PickUp()
	case 'x':
		u.bashPending = true
		g.Msg("bash in which direction? (esc cancels)", 8)
	case '?':
		g.State = StateHelp
	case 'i':
		g.State = StateInventory
		u.slotSel = 0
	case 'v':
		g.State = StateLook
		g.LookX, g.LookY = g.Player.X, g.Player.Y
	case '~':
		u.msgExpanded = true
	case 'o':
		g.State = StateOptions
		u.optSel = 0
		u.optScroll = 0
	case '@':
		g.State = StatePlayerScreen
	case 's':
		g.SaveAndQuit()
	case '>':
		g.UsePortal()
	case ']':
		g.ChangeLevelRelative(1)
	case '[':
		g.ChangeLevelRelative(-1)
	case 'q', 'z':
		g.State = StateQuitConfirm
	}
	if e.Key() == tcell.KeyEscape {
		g.State = StateQuitConfirm
	}
}

func (u *UI) drawConfirm() {
	lines := []string{
		"are you sure you want to quit the game [y/n] ?",
		"",
		"y = save and quit    n = continue",
	}
	width := 0
	for _, l := range lines {
		if n := len([]rune(l)); n > width {
			width = n
		}
	}
	width += 6
	height := len(lines) + 4
	x := (screenW - width) / 2
	y := (screenH - height) / 2

	bg := u.uiColor(RolePanelBackground, 0)
	fg := u.uiColor(RolePrimaryText, 15)
	border := u.uiColor(RoleDanger, 12)

	u.fillTcell(x, y, width, height, ' ', fg, bg)
	for xx := x + 1; xx < x+width-1; xx++ {
		u.putTcell(xx, y, '-', border, bg)
		u.putTcell(xx, y+height-1, '-', border, bg)
	}
	for yy := y + 1; yy < y+height-1; yy++ {
		u.putTcell(x, yy, '|', border, bg)
		u.putTcell(x+width-1, yy, '|', border, bg)
	}
	u.putTcell(x, y, '+', border, bg)
	u.putTcell(x+width-1, y, '+', border, bg)
	u.putTcell(x, y+height-1, '+', border, bg)
	u.putTcell(x+width-1, y+height-1, '+', border, bg)
	for i, l := range lines {
		u.putStrTcell(x+3, y+2+i, l, fg, bg)
	}
}

func (u *UI) handleQuitConfirm(e *tcell.EventKey) {
	switch e.Rune() {
	case ' ':
		u.g.State = StatePlaying
	case 'y', 'Y':
		u.g.SaveAndQuit()
	case 'n', 'N':
		u.g.State = StatePlaying
	default:
		if e.Key() == tcell.KeyEscape {
			u.g.State = StatePlaying
		}
	}
}

func (u *UI) drawGameOver(victory bool) {
	u.clear(0)
	if victory {
		u.putStr(centerText("you survive", screenW), 7, "you survive", 14, 0)
		u.putStr(centerText("~~~~~~~~~~~", screenW), 8, "~~~~~~~~~~~", 14, 0)
		lines := wrapText("as you enter the time portal, you feel highly accelerated and years seem to be seconds ... after a while you were thrown out on the other side of the portal. you stand up slowly and look around, examining the cold, shiny walls of your cell ...", 70)
		for i, l := range lines {
			u.putStr(5, 11+i, l, 10, 0)
		}
		u.putStr(centerText("press enter to return to title", screenW), 22, "press enter to return to title", 8, 0)
	} else {
		u.putStr(centerText("you died.", screenW), 10, "you died.", 12, 0)
		u.putStr(centerText("=== session ended ===", screenW), 12, "=== session ended ===", 7, 0)
		u.putStr(centerText("press enter to return to title", screenW), 16, "press enter to return to title", 8, 0)
	}
}

func (u *UI) handleGameOver(e *tcell.EventKey) {
	if isEnter(e) || e.Key() == tcell.KeyEscape || e.Rune() == 'q' {
		u.g = NewGame(time.Now().UnixNano())
		u.g.State = StateTitle
		u.titleSel = 0
	}
}

func (u *UI) drawLookOverlay() {
	g := u.g
	desc := g.TileDescription(g.LookX, g.LookY)
	bg := u.uiColor(RolePanelBackground, 1)
	fg := u.uiColor(RolePrimaryText, 15)
	dim := u.uiColor(RoleDimText, 8)
	u.fillTcell(0, 22, screenW, 2, ' ', fg, bg)
	if lines := wrapText(desc, screenW-2); len(lines) > 0 {
		desc = lines[0]
	}
	u.putStrTcell(1, 22, desc, fg, bg)
	u.putStrTcell(1, 23, fmt.Sprintf("look: [%d:%d] move keys move  enter/esc/space quit", g.LookX, g.LookY), dim, bg)
}

func (u *UI) handleLook(e *tcell.EventKey) {
	g := u.g
	if dx, dy, ok := directionFromEvent(e); ok {
		g.LookX += dx
		g.LookY += dy
		return
	}
	if isEnter(e) || e.Key() == tcell.KeyEscape || e.Rune() == 'z' || e.Rune() == 'v' || e.Rune() == ' ' {
		g.State = StatePlaying
	}
}

func (u *UI) drawThrowOverlay() {
	g := u.g
	bg := u.uiColor(RolePanelBackground, 1)
	fg := u.uiColor(RolePrimaryText, 15)
	dim := u.uiColor(RoleDimText, 8)
	u.fillTcell(0, 22, screenW, 2, ' ', fg, bg)
	u.putStrTcell(1, 22, fmt.Sprintf("throw: target [%d:%d]  range %d  enter throw  esc/space cancel", g.LookX, g.LookY, g.FOVRadius()), fg, bg)
	u.putStrTcell(1, 23, "move keys move cursor", dim, bg)
}

func (g *Game) moveThrowCursor(dx, dy int) {
	nx, ny := g.LookX+dx, g.LookY+dy
	if !g.Map.InBounds(nx, ny) || !g.Map.Visible[g.Map.idx(nx, ny)] {
		g.Msg("you can only target visible tiles.", 8)
		return
	}
	if chebyshev(g.Player.X, g.Player.Y, nx, ny) > g.FOVRadius() {
		g.Msg("that is beyond your throwing range.", 8)
		return
	}
	g.LookX, g.LookY = nx, ny
}

func (u *UI) handleThrow(e *tcell.EventKey) {
	g := u.g
	if dx, dy, ok := directionFromEvent(e); ok {
		g.moveThrowCursor(dx, dy)
		return
	}
	switch e.Key() {
	case tcell.KeyEscape:
		g.State = StatePlaying
	case tcell.KeyEnter:
		g.State = StatePlaying
		g.ThrowGrenade(g.LookX, g.LookY)
	default:
		if e.Rune() == 'z' || e.Rune() == ' ' {
			g.State = StatePlaying
		}
	}
}

func (u *UI) debugSpawnCritter() {
	g := u.g
	spots := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {-1, -1}, {1, -1}, {-1, 1}}
	for _, s := range spots {
		x, y := g.Player.X+s[0], g.Player.Y+s[1]
		if g.Map.IsPassable(x, y) && g.CritterAt(x, y) == nil {
			c := g.SpawnCritter(0, x, y)
			g.Msg("debug: spawned "+c.NameOrUnknown()+".", 10)
			return
		}
	}
	g.Msg("debug: no free spot.", 12)
}

func (u *UI) debugSpawnItem() {
	g := u.g
	x, y := g.Player.X, g.Player.Y
	for {
		proto := g.Roll(len(ItemDefs))
		if !ItemAvailable(proto) {
			continue
		}
		it := &Item{Proto: proto, X: x, Y: y, Charges: ItemDefs[proto].Clip}
		if ItemDefs[proto].Kind == KindMedical {
			it.Charges = ItemDefs[proto].Charges
		}
		g.Items = append(g.Items, it)
		g.Msg("debug: spawned item "+it.Name()+".", 10)
		return
	}
}
