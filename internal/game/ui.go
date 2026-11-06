package game

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	screenW  = 80
	screenH  = 50
	mapViewH = 42
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

// DOS/CP437 colour indices in the original data are not ANSI indices.
// Map them explicitly to the tcell 16-colour palette.
var dosPalette = [16]tcell.Color{
	tcell.ColorBlack,
	tcell.ColorNavy,    // 1 blue
	tcell.ColorGreen,   // 2 green
	tcell.ColorTeal,    // 3 cyan
	tcell.ColorMaroon,  // 4 red
	tcell.ColorPurple,  // 5 magenta
	tcell.ColorOlive,   // 6 brown
	tcell.ColorSilver,  // 7 light gray
	tcell.ColorGray,    // 8 dark gray
	tcell.ColorBlue,    // 9 light blue
	tcell.ColorLime,    // 10 light green
	tcell.ColorAqua,    // 11 light cyan
	tcell.ColorRed,     // 12 light red
	tcell.ColorFuchsia, // 13 light magenta
	tcell.ColorYellow,  // 14 yellow
	tcell.ColorWhite,   // 15 white
}

func pal(i uint8) tcell.Color {
	if int(i) < 16 {
		return dosPalette[i]
	}
	return tcell.ColorWhite
}

type UI struct {
	screen tcell.Screen
	g      *Game
	ox, oy int

	titleSel   int
	classSel   int
	invTab     int // 0 inventory, 1 ground
	invSel     int
	invScroll  int
	optSel     int
	optScroll  int
	helpSel    int
	helpScroll int
	msgScroll  int
	skillSel   int
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
	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	s.Clear()

	u := &UI{screen: s, g: NewGame(time.Now().UnixNano())}
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
	if x < 0 || y < 0 || x >= screenW || y >= screenH {
		return
	}
	st := tcell.StyleDefault.Foreground(pal(fg)).Background(pal(bg))
	u.screen.SetContent(u.ox+x, u.oy+y, r, nil, st)
}

func (u *UI) putStr(x, y int, s string, fg, bg uint8) {
	for i, r := range s {
		u.put(x+i, y, r, fg, bg)
	}
}

func (u *UI) fill(x, y, w, h int, r rune, fg, bg uint8) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			u.put(xx, yy, r, fg, bg)
		}
	}
}

func (u *UI) clear(bg uint8) {
	u.fill(0, 0, screenW, screenH, ' ', 7, bg)
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
	case StateCharCreate:
		u.drawCharCreate()
	case StatePlaying:
		u.drawGame()
	case StateInventory:
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
	u.putStr(centerText("==================================================", screenW), 6, "==================================================", 2, 0)
	u.putStr(centerText("H  U  L  K     I  N  F  E  C  T  I  O  N", screenW), 8, "H  U  L  K     I  N  F  E  C  T  I  O  N", 10, 0)
	u.putStr(centerText("==================================================", screenW), 9, "==================================================", 2, 0)
	u.putStr(centerText("u l t i m a t e   c h a i n s a w   m a s s a c r e", screenW), 11, "u l t i m a t e   c h a i n s a w   m a s s a c r e", 14, 0)
	lines := []struct {
		s  string
		fg uint8
	}{
		{"hulk infection (v0.3a)", 10},
		{"copyright (c) 2006, 2007 by michal przybylowicz", 8},
		{"all rights reserved", 8},
		{"", 7},
		{"( this edition is dedicated for Ryzy )", 14},
		{"www.przybyl666.prv.pl", 11},
	}
	for i, l := range lines {
		u.putStr(centerText(l.s, screenW), 15+i, l.s, l.fg, 0)
	}
	menu := []string{"New game", "Continue", "Quit"}
	for i, s := range menu {
		fg := uint8(7)
		if i == u.titleSel {
			fg = 15
			s = "> " + s + " <"
		}
		u.putStr(centerText(s, screenW), 28+i, s, fg, 0)
	}
	u.putStr(centerText("arrows / enter", screenW), 35, "arrows / enter", 8, 0)
}

func (u *UI) drawIntro() {
	u.clear(0)
	intro := []string{
		"year 2007",
		"~~~~~~~~~",
		"",
		"hungry for adventure you travel deep into nearest woods to",
		"find the truth about local gossip ... you heard that there is an",
		"ancient evil lurking deep underground somewhere in the forest ...",
		"",
		"after a while of meaningless walking you find a small wooden shed,",
		"tired and little exausted you head yourself toward it ...",
		"",
		"press enter to continue",
	}
	for i, s := range intro {
		fg := uint8(10)
		if i == 0 {
			fg = 14
		}
		if i == len(intro)-1 {
			fg = 8
		}
		u.putStr(centerText(s, screenW), 12+i, s, fg, 0)
	}
}

func (u *UI) drawCharCreate() {
	u.clear(0)
	u.putStr(centerText("[ @ ] proffession...", screenW), 1, "[ @ ] proffession...", 11, 0)
	u.putStr(2, 3, "select class", 8, 0)
	for i, c := range ClassDefs {
		mark := "  "
		fg := uint8(7)
		if i == u.classSel {
			mark = "> "
			fg = 15
		}
		u.putStr(4, 5+i*2, mark+c.Name, fg, 0)
	}
	c := &ClassDefs[u.classSel]
	u.putStr(2, 12, "class: "+c.Name, 14, 0)
	descLines := wrapText(c.Desc, 74)
	for i, l := range descLines {
		if i > 2 {
			break
		}
		u.putStr(2, 13+i, l, 10, 0)
	}
	u.putStr(2, 17, fmt.Sprintf("strength:%d  toughness:%d  agility:%d  intelligence:%d", c.STR, c.TGH, c.AGI, c.INT), 7, 0)
	for i, skill := range ActiveSkills {
		col, row := i/3, i%3
		x := 2 + col*26
		y := 19 + row
		u.putStr(x, y, fmt.Sprintf("%-16s %3d", SkillNames[skill], c.Skills[skill]), 10, 0)
	}
	u.putStr(2, 23, "name: "+u.g.PlayerName, 14, 0)
	u.putStr(2, 24, "arrows: select class   r: reroll name   enter: begin", 8, 0)
	u.putStr(2, 25, "esc: back", 8, 0)
}

func (u *UI) drawGame() {
	g := u.g
	if g.Map == nil || g.Player == nil {
		u.clear(0)
		return
	}
	u.clear(0)
	left := g.Player.X - screenW/2
	top := g.Player.Y - mapViewH/2
	if left < 0 {
		left = 0
	}
	if top < 0 {
		top = 0
	}
	if left > g.Map.W-screenW {
		left = g.Map.W - screenW
	}
	if top > g.Map.H-mapViewH {
		top = g.Map.H - mapViewH
	}
	for sy := 0; sy < mapViewH; sy++ {
		for sx := 0; sx < screenW; sx++ {
			mx, my := left+sx, top+sy
			if !g.Map.InBounds(mx, my) {
				u.put(sx, sy, ' ', 0, 0)
				continue
			}
			i := g.Map.idx(mx, my)
			if !g.Map.Seen[i] {
				u.put(sx, sy, ' ', 0, 0)
				continue
			}
			d := g.Map.Def(mx, my)
			vis := g.Map.Visible[i] || g.Options.Bool("optMapDrawHiddenTiles")
			fg := d.Color
			if !vis {
				fg = 8
			}
			r := cp437Rune(d.Glyph)
			if g.Options.Bool("optMapDrawPortals") && (d.Name == "time portal" || d.Name == "stone stairs" || d.Name == "elevator") {
				fg = 14
			}
			bg := uint8(0)
			if g.Map.Blood[i] > 0 && d.Type == 1 {
				bg = 4 // DOS red background
				if !vis {
					bg = 0
				}
			}
			if g.Map.Fire[i] > 0 {
				r = '^'
				fg = 12
			} else if g.Map.Acid[i] > 0 {
				r = '~'
				fg = 10
			} else if g.Map.Smoke[i] > 0 && vis {
				r = '%'
				fg = 8
			}
			if vis {
				if it := g.ItemAt(mx, my); it != nil {
					if d2 := it.Def(); d2 != nil {
						r = cp437Rune(d2.Glyph)
						fg = d2.Color
						if bg == 4 {
							fg = 15 // keep corpses/items visible on blood
						}
					}
				}
				if c := g.CritterAt(mx, my); c != nil && c.IsAlive() {
					if g.Map.Visible[i] || g.Options.Bool("optDrawDebugAI") {
						r = cp437Rune(c.Glyph)
						fg = c.Color
						if c.Blood > 0 {
							fg = 12
							if bg == 4 {
								fg = 15
							}
						}
					}
				}
			}
			if mx == g.Player.X && my == g.Player.Y {
				r = '@'
				fg = 15
			}
			u.put(sx, sy, r, fg, bg)
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
			u.put(sx, sy, r, 0, 7)
		}
	}
	u.drawMessages(42, 4)
	u.drawStatus(46)
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
	r, fg := cp437Rune(d.Glyph), d.Color
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

func (u *UI) drawMessages(y, rows int) {
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

func (u *UI) drawStatus(y int) {
	g := u.g
	p := g.Player
	need := 200 * p.Level
	line1 := fmt.Sprintf("%-12s health: (%s)   exp: %d/%d   [lvl %d]", p.Name, p.HealthState(), p.XP, need, p.Level)
	u.putStr(0, y, line1, 10, 0)
	wear := "none"
	if a := p.Equip[EquipArmor]; a != nil {
		wear = a.Name()
	}
	wep := "none"
	if w := p.Equip[EquipMelee]; w != nil {
		wep = w.Name()
	}
	tac := []string{"defensive", "normal", "aggressive"}[p.Tactics]
	line2 := fmt.Sprintf("weapon: %-16s wearing: %-14s tactic: %s", wep, wear, tac)
	u.putStr(0, y+1, line2, 7, 0)
	line3 := fmt.Sprintf("weight: %d/%d  burden: %s", p.Weight(), p.MaxWeight(), p.Burden())
	if g.Options.Bool("optPlayerShowAVDV") {
		line3 += fmt.Sprintf("  av:%d dv:%d", g.AttackValue(p, 1), g.DefenseValue(p))
	}
	u.putStr(0, y+2, line3, 8, 0)
	right := fmt.Sprintf("[%04d:%04d] [%s]", p.X, p.Y, g.Map.Name)
	u.putStr(screenW-len(right), y+2, right, 14, 0)
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
		if isEnter(e) || e.Key() == tcell.KeyEscape {
			g.State = StateCharCreate
			if g.PlayerName == "" {
				g.PlayerName = RandomNames[g.Roll(len(RandomNames))]
			}
		}
	case StateCharCreate:
		u.handleCharCreate(e)
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
		if e.Rune() == 'q' || e.Key() == tcell.KeyEscape {
			u.g.Quit = true
			return
		}
		if !isEnter(e) {
			return
		}
		switch u.titleSel {
		case 0:
			name := u.g.PlayerName
			u.g = NewGame(time.Now().UnixNano())
			u.g.PlayerName = name
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
		u.classSel = (u.classSel + len(ClassDefs) - 1) % len(ClassDefs)
	case tcell.KeyDown:
		u.classSel = (u.classSel + 1) % len(ClassDefs)
	case tcell.KeyEscape:
		u.g.State = StateIntro
	default:
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
		g.Player.XP += 200 * g.Player.Level
		g.CheckLevelUp()
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
	// movement
	moves := map[tcell.Key][2]int{
		tcell.KeyUp: {0, -1}, tcell.KeyDown: {0, 1}, tcell.KeyLeft: {-1, 0}, tcell.KeyRight: {1, 0},
		tcell.KeyHome: {-1, -1}, tcell.KeyEnd: {1, 1}, tcell.KeyPgUp: {1, -1}, tcell.KeyPgDn: {-1, 1},
	}
	if d, ok := moves[e.Key()]; ok {
		g.MovePlayer(d[0], d[1])
		return
	}
	r := e.Rune()
	switch r {
	case '8':
		g.MovePlayer(0, -1)
	case '9':
		g.MovePlayer(1, -1)
	case '7':
		g.MovePlayer(-1, -1)
	case '6':
		g.MovePlayer(1, 0)
	case '5':
		g.WaitTurn()
	case '4':
		g.MovePlayer(-1, 0)
	case '3':
		g.MovePlayer(1, 1)
	case '2':
		g.MovePlayer(0, 1)
	case '1':
		g.MovePlayer(-1, 1)
	case 'b':
		g.Bash()
	case ',':
		g.PickUp()
	case 'f':
		g.PlayerMeleeAttack()
	case '?':
		g.State = StateHelp
	case 'i':
		g.State = StateInventory
		u.invSel = 0
		u.invTab = 0
		u.invScroll = 0
	case 'l':
		g.State = StateLook
		g.LookX, g.LookY = g.Player.X, g.Player.Y
	case 'm':
		g.State = StateMessages
		u.msgScroll = 0
	case 'O':
		g.State = StateOptions
		u.optSel = 0
		u.optScroll = 0
	case '@':
		g.State = StatePlayerScreen
	case 'S':
		g.SaveAndQuit()
	case 'x':
		g.State = StateSkills
		u.skillSel = 0
	case 'T':
		g.ChangeTactics()
	case 't':
		if g.hasThrowable() {
			g.State = StateThrow
			g.LookX, g.LookY = g.Player.X, g.Player.Y
		} else {
			g.Msg("you do not have any throwable items ready.", 12)
		}
	case '>':
		g.UsePortal()
	case ']':
		g.ChangeLevelRelative(1)
	case '[':
		g.ChangeLevelRelative(-1)
	case 'Q':
		g.State = StateQuitConfirm
	case 'z':
		g.State = StateQuitConfirm
	case 'Z':
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
	width += 4
	height := len(lines) + 2
	x := (screenW - width) / 2
	y := (screenH - height) / 2
	const bg = uint8(4)
	u.fill(x, y, width, height, ' ', 15, bg)
	for i, l := range lines {
		u.putStr(x+2, y+1+i, l, 15, bg)
	}
}

func (u *UI) handleQuitConfirm(e *tcell.EventKey) {
	switch e.Rune() {
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
		u.putStr(centerText("you survive", screenW), 18, "you survive", 14, 0)
		u.putStr(centerText("~~~~~~~~~~~", screenW), 19, "~~~~~~~~~~~", 14, 0)
		lines := wrapText("as you enter the time portal, you feel highly accelerated and years seem to be seconds ... after a while you were thrown out on the other side of the portal. you stand up slowly and look around, examining the cold, shiny walls of your cell ...", 70)
		for i, l := range lines {
			u.putStr(5, 22+i, l, 10, 0)
		}
		u.putStr(centerText("press enter to return to title", screenW), 32, "press enter to return to title", 8, 0)
	} else {
		u.putStr(centerText("you died.", screenW), 20, "you died.", 12, 0)
		u.putStr(centerText("=== session ended ===", screenW), 22, "=== session ended ===", 7, 0)
		u.putStr(centerText("press enter to return to title", screenW), 26, "press enter to return to title", 8, 0)
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
	u.fill(0, 42, screenW, 2, ' ', 8, 0)
	u.fill(0, 42, screenW, 1, ' ', 15, 1)
	u.putStr(1, 42, desc, 15, 1)
	u.putStr(1, 43, fmt.Sprintf("look: [%d:%d] arrows move  enter/esc quit", g.LookX, g.LookY), 8, 0)
}

func (u *UI) handleLook(e *tcell.EventKey) {
	g := u.g
	switch e.Key() {
	case tcell.KeyUp:
		g.LookY--
	case tcell.KeyDown:
		g.LookY++
	case tcell.KeyLeft:
		g.LookX--
	case tcell.KeyRight:
		g.LookX++
	default:
		switch e.Rune() {
		case '8':
			g.LookY--
		case '2':
			g.LookY++
		case '4':
			g.LookX--
		case '6':
			g.LookX++
		default:
			if isEnter(e) || e.Key() == tcell.KeyEscape || e.Rune() == 'z' || e.Rune() == 'l' {
				g.State = StatePlaying
			}
		}
	}
}

func (u *UI) drawThrowOverlay() {
	g := u.g
	u.fill(0, 42, screenW, 2, ' ', 8, 0)
	u.fill(0, 42, screenW, 1, ' ', 15, 1)
	u.putStr(1, 42, fmt.Sprintf("throw: target [%d:%d]  range %d  enter throw  esc cancel", g.LookX, g.LookY, g.FOVRadius()), 15, 1)
	u.putStr(1, 43, "arrows move cursor", 8, 0)
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
	switch e.Key() {
	case tcell.KeyUp:
		g.moveThrowCursor(0, -1)
	case tcell.KeyDown:
		g.moveThrowCursor(0, 1)
	case tcell.KeyLeft:
		g.moveThrowCursor(-1, 0)
	case tcell.KeyRight:
		g.moveThrowCursor(1, 0)
	case tcell.KeyEscape:
		g.State = StatePlaying
	case tcell.KeyEnter:
		g.State = StatePlaying
		g.ThrowGrenade(g.LookX, g.LookY)
	default:
		switch e.Rune() {
		case '8':
			g.moveThrowCursor(0, -1)
		case '2':
			g.moveThrowCursor(0, 1)
		case '4':
			g.moveThrowCursor(-1, 0)
		case '6':
			g.moveThrowCursor(1, 0)
		case 'z':
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
