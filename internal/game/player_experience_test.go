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

func TestSetupNewGameCreatesPriestProfileAndKit(t *testing.T) {
	for seed := int64(100); seed < 108; seed++ {
		g := NewGame(seed)
		g.SetupNewGame(0, "tester")
		if g.Player.Name != "tester" {
			t.Fatalf("player name = %q", g.Player.Name)
		}
		if len(g.KnownRites()) != 1 {
			t.Fatalf("seed %d has %d rites, want 1", seed, len(g.KnownRites()))
		}
		if g.RiteCharges[0] != 1 {
			t.Fatalf("seed %d rite charge = %d", seed, g.RiteCharges[0])
		}
		if g.Player.Slots.Hands == nil || g.Player.Slots.Hands.Name() != "sickle" {
			t.Fatalf("seed %d starting weapon missing", seed)
		}
		if g.Player.Slots.RobeSlots != 3 {
			t.Fatalf("seed %d robe slots = %d", seed, g.Player.Slots.RobeSlots)
		}
		if g.Player.Slots.Robe[0] == nil || g.Player.Slots.Robe[0].Name() != "leather shard" {
			t.Fatalf("seed %d starting shard missing", seed)
		}
		if g.Player.Slots.Robe[1] == nil || g.Player.Slots.Robe[1].Name() != "healing vial" {
			t.Fatalf("seed %d starting vial missing", seed)
		}
		if g.Player.Slots.Robe[2] != nil {
			t.Fatalf("seed %d third robe slot should start empty", seed)
		}
		if g.Player.Slots.Body != nil || g.Player.Slots.BeltUnlocked != 0 {
			t.Fatalf("seed %d still uses legacy body/belt slots", seed)
		}
	}
}

func TestAllThemesUseExpectedMapSizes(t *testing.T) {
	for level := 0; level < numWoodsMaps; level++ {
		g := NewGame(int64(level + 500))
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(level)
		if g.Map.W != screenW || g.Map.H != woodsMapHeight {
			t.Fatalf("woods level %d map = %dx%d, want %dx%d", level, g.Map.W, g.Map.H, screenW, woodsMapHeight)
		}
	}
	for _, level := range []int{9, 17, 18, 26, 27} {
		g := NewGame(int64(level + 500))
		g.SetupNewGame(0, "tester")
		g.GenerateLevel(level)
		if g.Map.W != screenW || g.Map.H != mapHeight {
			t.Fatalf("section level %d map = %dx%d, want %dx%d", level, g.Map.W, g.Map.H, screenW, mapHeight)
		}
	}
}

func TestGameUses80x25Layout(t *testing.T) {
	if screenW != 80 || screenH != 25 || mapViewH != 23 || woodsMapHeight != 23 {
		t.Fatalf("layout constants = %dx%d view %d woods %d, want 80x25 view 23 woods 23", screenW, screenH, mapViewH, woodsMapHeight)
	}
	g := NewGame(710)
	g.SetupNewGame(0, "tester")
	g.Msg("layout test message", 10)
	g.State = StatePlaying
	g.UpdateFOV()
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.drawGame()
	s.Show()
	rows := strings.Split(screenText(s), "\n")
	if len(rows) != screenH+1 {
		t.Fatalf("rendered %d rows, want %d plus trailing newline", len(rows)-1, screenH)
	}
	if !strings.Contains(rows[0], "layout test message") {
		t.Fatalf("row 0 is not the latest message: %q", rows[0])
	}
	if !strings.Contains(rows[hudRow], "HP ") || !strings.Contains(rows[hudRow], "AP ") {
		t.Fatalf("row %d is not the HUD: %q", hudRow, rows[hudRow])
	}
	found := false
	for y := 1; y <= mapViewH; y++ {
		if strings.Contains(rows[y], "@") {
			if y != g.Player.Y+1 {
				t.Fatalf("player rendered on row %d, want %d", y, g.Player.Y+1)
			}
			found = true
		}
	}
	if !found {
		t.Fatal("player is not rendered inside the 80x23 map area")
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
	if !strings.Contains(hud, "sickle") {
		t.Fatalf("HUD missing weapon name: %q", hud)
	}
	if !strings.Contains(hud, "Rite:") {
		t.Fatalf("HUD missing rite: %q", hud)
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

func TestIntroShowsGeneratedPriestAndRite(t *testing.T) {
	g := NewGame(703)
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.drawIntro()
	s.Show()
	text := screenText(s)
	if !strings.Contains(text, g.PlayerName) {
		t.Fatalf("intro missing priest name %q", g.PlayerName)
	}
	if !strings.Contains(text, "Rite:") || !strings.Contains(text, RiteDefs[g.PlayerRites[0]].Name) {
		t.Fatalf("intro missing generated rite: %q", text)
	}
	if strings.Contains(text, "select class") {
		t.Fatal("intro still shows class selection")
	}
	if !strings.Contains(text, "press enter to begin") {
		t.Fatal("intro missing start prompt")
	}
}

func TestInventoryGetPicksUpItem(t *testing.T) {
	g := NewGame(704)
	g.SetupNewGame(0, "tester")
	g.State = StateInventory
	it := &Item{Proto: 48, X: g.Player.X, Y: g.Player.Y, Charges: 5}
	g.Items = append(g.Items, it)
	u := &UI{g: g}
	u.handleInventory(tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone))
	if g.ItemAt(g.Player.X, g.Player.Y) != nil {
		t.Fatal("item was not removed from the ground")
	}
	if g.Player.Slots.Robe[2] == nil || g.Player.Slots.Robe[2].Proto != 48 {
		t.Fatal("picked-up item was not placed in the free robe slot")
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
