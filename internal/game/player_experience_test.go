package game

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func newTestUI(g *Game, mode ColorMode) (*UI, tcell.SimulationScreen) {
	s := tcell.NewSimulationScreen("UTF-8")
	if err := s.Init(); err != nil {
		panic(err)
	}
	s.SetSize(screenW, screenH)
	return &UI{screen: s, g: g, colorMode: mode}, s
}

func screenText(s tcell.SimulationScreen) string {
	cells, w, h := s.GetContents()
	var b strings.Builder
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cell := cells[y*w+x]
			if len(cell.Runes) == 0 {
				b.WriteByte(' ')
				continue
			}
			b.WriteRune(cell.Runes[0])
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func TestPlayerClassRoster(t *testing.T) {
	wantNames := []string{"Butcher", "Hunter", "Medic", "Mechanic", "Prepper"}
	if len(PlayerClasses) != len(wantNames) {
		t.Fatalf("class count = %d, want %d", len(PlayerClasses), len(wantNames))
	}
	seen := map[string]bool{}
	for i, c := range PlayerClasses {
		if c.Name != wantNames[i] {
			t.Fatalf("class %d = %q, want %q", i, c.Name, wantNames[i])
		}
		if seen[c.Name] {
			t.Fatalf("duplicate class name %q", c.Name)
		}
		seen[c.Name] = true
		if c.STR <= 0 || c.TGH <= 0 || c.AGI <= 0 || c.INT <= 0 {
			t.Fatalf("%s has a non-positive attribute", c.Name)
		}
		kit := ClassStartingKit(i)
		if kit.Hands < 0 || kit.Body < 0 || len(kit.Belt) != 1 {
			t.Fatalf("%s has an invalid starting kit: %+v", c.Name, kit)
		}
		for _, proto := range append([]int{kit.Hands, kit.Body}, kit.Belt...) {
			if proto < 0 || proto >= len(ItemDefs) {
				t.Fatalf("%s kit references invalid item %d", c.Name, proto)
			}
			if !ItemAvailable(proto) {
				t.Fatalf("%s kit references unavailable item %d", c.Name, proto)
			}
		}
	}
}

func TestSetupNewGameUsesClassValuesAndKit(t *testing.T) {
	for class, def := range PlayerClasses {
		g := NewGame(int64(100 + class))
		g.SetupNewGame(class, "tester")
		if g.Player.STR != def.STR || g.Player.TGH != def.TGH ||
			g.Player.AGI != def.AGI || g.Player.INT != def.INT {
			t.Fatalf("%s attributes do not match class data", def.Name)
		}
		if g.Player.Skills != def.Skills {
			t.Fatalf("%s skills do not match class data", def.Name)
		}
		if g.Player.Slots.BeltUnlocked != 2 {
			t.Fatalf("%s starts with %d belt slots, want 2", def.Name, g.Player.Slots.BeltUnlocked)
		}
		kit := ClassStartingKit(class)
		if g.Player.Slots.Hands == nil || g.Player.Slots.Hands.Proto != kit.Hands {
			t.Fatalf("%s starting hands do not match kit", def.Name)
		}
		if g.Player.Slots.Body == nil || g.Player.Slots.Body.Proto != kit.Body {
			t.Fatalf("%s starting body does not match kit", def.Name)
		}
		if g.Player.Slots.Belt[0] == nil || g.Player.Slots.Belt[0].Proto != kit.Belt[0] {
			t.Fatalf("%s starting belt 1 does not match kit", def.Name)
		}
		if g.Player.Slots.Belt[1] != nil {
			t.Fatalf("%s should start with belt 2 empty", def.Name)
		}
		if g.Player.Slots.Belt[2] != nil {
			t.Fatalf("%s starts with an item in locked belt 3", def.Name)
		}
	}
}

func TestAllThemesUseViewportSizedMaps(t *testing.T) {
	for _, level := range []int{0, 8, 9, 17, 18, 26, 27} {
		g := NewGame(int64(level + 500))
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(level)
		if g.Map.W != screenW || (g.Map.H != mapHeight && g.Map.H != mapViewH) {
			t.Fatalf("level %d map = %dx%d, want width %d and height %d or %d", level, g.Map.W, g.Map.H, screenW, mapViewH, mapHeight)
		}
	}
}

func TestHUDIsOneLine(t *testing.T) {
	g := NewGame(700)
	g.SetupNewGame(0, "tester")
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.drawStatus(hudRow)
	s.Show()
	text := screenText(s)
	rows := strings.Split(text, "\n")
	if len(rows) < screenH+1 {
		t.Fatal("simulation screen has too few rows")
	}
	hud := rows[hudRow]
	if !strings.Contains(hud, "HP ") || !strings.Contains(hud, "AP ") {
		t.Fatalf("HUD missing HP/AP: %q", hud)
	}
	if !strings.Contains(hud, "hammer") {
		t.Fatalf("HUD missing weapon name: %q", hud)
	}
	if strings.Contains(hud, "[a]") || strings.Contains(hud, "Lvl") {
		t.Fatalf("HUD still contains old slot/level text: %q", hud)
	}
}

func TestLookOverlayUsesReadablePanelColours(t *testing.T) {
	g := NewGame(701)
	g.SetupNewGame(0, "tester")
	g.UpdateFOV()
	g.State = StateLook
	g.LookX, g.LookY = g.Player.X, g.Player.Y
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.drawLookOverlay()
	s.Show()
	text := screenText(s)
	rows := strings.Split(text, "\n")
	if !strings.Contains(rows[22], "there is") {
		t.Fatalf("look description missing: %q", rows[22])
	}
	if !strings.Contains(rows[23], "look:") {
		t.Fatalf("look controls missing: %q", rows[23])
	}
	cells, w, _ := s.GetContents()
	_, bg, _ := cells[22*w+0].Style.Decompose()
	want := u.uiColor(RolePanelBackground, 1)
	if bg != want || !bg.IsRGB() {
		t.Fatalf("look panel background = %v, want %v", bg, want)
	}
}

func TestCharCreateShowsFiveClassesAndKit(t *testing.T) {
	g := NewGame(703)
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.classSel = len(PlayerClasses) - 1
	u.drawCharCreate()
	s.Show()
	text := screenText(s)
	for _, class := range PlayerClasses {
		if !strings.Contains(text, class.Name) {
			t.Fatalf("class selection missing %q", class.Name)
		}
	}
	if !strings.Contains(text, "hands: crowbar") {
		t.Fatalf("class selection missing starting kit: %q", text)
	}
}

func TestInventoryGetPicksUpItem(t *testing.T) {
	g := NewGame(704)
	g.SetupNewGame(0, "tester")
	g.State = StateInventory
	it := &Item{Proto: 48, X: g.Player.X, Y: g.Player.Y, Charges: 5}
	g.Items = append(g.Items, it)
	u := &UI{g: g}
	u.handleInventory(tcell.NewEventKey(tcell.KeyRune, 'g', tcell.ModNone))
	if g.ItemAt(g.Player.X, g.Player.Y) != nil {
		t.Fatal("item was not removed from the ground")
	}
	if g.Player.Slots.Belt[1] == nil || g.Player.Slots.Belt[1].Proto != 48 {
		t.Fatal("picked-up item was not placed in the free belt slot")
	}
	if g.State != StateInventory {
		t.Fatalf("inventory closed after pickup: state %v", g.State)
	}
}

func TestMovingOntoItemAnnouncesIt(t *testing.T) {
	g := NewGame(706)
	g.SetupNewGame(0, "tester")
	g.State = StatePlaying
	var dx, dy int
	found := false
	for _, d := range dirs4 {
		nx, ny := g.Player.X+d[0], g.Player.Y+d[1]
		if g.Map.IsPassable(nx, ny) && g.CritterAt(nx, ny) == nil && g.ItemAt(nx, ny) == nil {
			dx, dy, found = d[0], d[1], true
			break
		}
	}
	if !found {
		t.Fatal("no free adjacent tile for movement test")
	}
	g.Items = append(g.Items, &Item{Proto: 10, X: g.Player.X + dx, Y: g.Player.Y + dy})
	g.MovePlayer(dx, dy)
	msgs := g.LastMessages(1)
	if len(msgs) != 1 || msgs[0].Text != "Crowbar is lying here." {
		t.Fatalf("last message = %q, want item announcement", msgs)
	}
}

func TestLookDescriptionNamesGroundItem(t *testing.T) {
	g := NewGame(705)
	g.SetupNewGame(0, "tester")
	g.UpdateFOV()
	g.Items = append(g.Items, &Item{Proto: 10, X: g.Player.X, Y: g.Player.Y})
	desc := g.TileDescription(g.Player.X, g.Player.Y)
	if !strings.Contains(desc, "Crowbar is lying here.") {
		t.Fatalf("tile description does not name the item: %q", desc)
	}

	g.State = StateLook
	g.LookX, g.LookY = g.Player.X, g.Player.Y
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.drawLookOverlay()
	s.Show()
	rows := strings.Split(screenText(s), "\n")
	if !strings.Contains(rows[22], "Crowbar is lying here.") {
		t.Fatalf("look overlay does not show the item: %q", rows[22])
	}
}

func TestSaveMigrationGrantsTwoBeltSlots(t *testing.T) {
	sc := saveCritter{
		Player: true,
		Slots: &saveSlots{
			BeltUnlocked: 1,
			Hands:        saveItem{Proto: 11},
			Body:         saveItem{Proto: 5},
		},
	}
	c, removed := critterFromSave(sc, 3)
	if removed != 0 {
		t.Fatalf("unexpected removed items: %d", removed)
	}
	if c.Slots.BeltUnlocked != 2 {
		t.Fatalf("migrated belt unlock = %d, want 2", c.Slots.BeltUnlocked)
	}
}

func TestAllClassesCanPickUpStartingLoot(t *testing.T) {
	for class, def := range PlayerClasses {
		g := NewGame(int64(900 + class))
		g.SetupNewGame(class, "tester")
		g.State = StatePlaying
		g.Items = append(g.Items, &Item{Proto: 48, X: g.Player.X, Y: g.Player.Y, Charges: 5})
		g.PickUp()
		if g.ItemAt(g.Player.X, g.Player.Y) != nil {
			t.Fatalf("%s could not pick up loot with the starting kit", def.Name)
		}
	}
}

func TestAllThemesRenderAtViewportSize(t *testing.T) {
	for _, level := range []int{0, 8, 9, 17, 18, 26, 27} {
		g := NewGame(int64(level + 800))
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(level)
		g.State = StatePlaying
		g.UpdateFOV()
		u, s := newTestUI(g, ColorModeTrueColor)
		u.drawGame()
		s.Show()
		text := screenText(s)
		if !strings.Contains(text, "@") {
			t.Fatalf("level %d did not render the player", level)
		}
		if len(strings.Split(text, "\n")) != screenH+1 {
			t.Fatalf("level %d rendered an unexpected screen height", level)
		}
		s.Fini()
	}
}

func TestQuitConfirmUsesReadablePanelColours(t *testing.T) {
	g := NewGame(707)
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.drawConfirm()
	s.Show()
	text := screenText(s)
	if !strings.Contains(text, "are you sure you want to quit") {
		t.Fatal("quit confirmation text missing")
	}
	wantBG := u.uiColor(RolePanelBackground, 0)
	wantFG := u.uiColor(RolePrimaryText, 15)
	if !wantBG.IsRGB() || !wantFG.IsRGB() {
		t.Fatalf("quit panel colours are not true colour: bg=%v fg=%v", wantBG, wantFG)
	}
	cells, w, h := s.GetContents()
	foundText := false
	for y := 0; y < h && !foundText; y++ {
		for x := 0; x < w; x++ {
			cell := cells[y*w+x]
			if len(cell.Runes) == 0 || cell.Runes[0] == ' ' {
				continue
			}
			fg, bg, _ := cell.Style.Decompose()
			if fg == wantFG && bg == wantBG {
				foundText = true
				break
			}
		}
	}
	if !foundText {
		t.Fatal("quit confirmation text does not use the readable panel colours")
	}
}

func TestTitleAndIntroText(t *testing.T) {
	g := NewGame(702)
	g.State = StateTitle
	u, s := newTestUI(g, ColorModeTrueColor)
	u.drawTitle()
	s.Show()
	title := screenText(s)
	if !strings.Contains(title, "CASE FILE: HULK INFECTION") {
		t.Fatal("title screen missing case-file title")
	}
	if !strings.Contains(title, "copyright (c) 2026 Michał Przybyłowicz") {
		t.Fatalf("title screen missing new copyright: %q", title)
	}
	for _, item := range []string{"new game", "continue", "quit"} {
		if !strings.Contains(title, item) {
			t.Fatalf("title menu missing %q", item)
		}
	}

	g.State = StateIntro
	u.drawIntro()
	s.Show()
	intro := screenText(s)
	if !strings.Contains(intro, "T H E   W O O D S   A R E   D E A D") {
		t.Fatal("intro missing new heading")
	}
	if !strings.Contains(intro, "source of the infection") {
		t.Fatal("intro missing the goal")
	}
	s.Fini()
}
