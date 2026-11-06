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
		if kit.Hands < 0 || kit.Body < 0 || len(kit.Belt) != 2 {
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
		for i := 0; i < 2; i++ {
			if g.Player.Slots.Belt[i] == nil || g.Player.Slots.Belt[i].Proto != kit.Belt[i] {
				t.Fatalf("%s starting belt %d does not match kit", def.Name, i+1)
			}
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
		if g.Map.W != screenW || g.Map.H != mapViewH {
			t.Fatalf("level %d map = %dx%d, want %dx%d", level, g.Map.W, g.Map.H, screenW, mapViewH)
		}
	}
}

func TestHUDIsOneLine(t *testing.T) {
	g := NewGame(700)
	g.SetupNewGame(0, "tester")
	u, s := newTestUI(g, ColorModeTrueColor)
	defer s.Fini()
	u.drawStatus(47)
	s.Show()
	text := screenText(s)
	rows := strings.Split(text, "\n")
	if len(rows) < 48 {
		t.Fatal("simulation screen has too few rows")
	}
	hud := rows[47]
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
	if !strings.Contains(rows[45], "there is") {
		t.Fatalf("look description missing: %q", rows[45])
	}
	if !strings.Contains(rows[46], "look:") {
		t.Fatalf("look controls missing: %q", rows[46])
	}
	cells, w, _ := s.GetContents()
	_, bg, _ := cells[45*w+0].Style.Decompose()
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
