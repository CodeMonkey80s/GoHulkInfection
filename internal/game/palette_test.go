package game

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func testEnv(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestDetectColorMode(t *testing.T) {
	cases := []struct {
		name   string
		colors int
		env    map[string]string
		want   ColorMode
	}{
		{"no colour support", 0, nil, ColorMode16},
		{"plain sixteen", 16, nil, ColorMode16},
		{"eighty eight", 88, nil, ColorMode16},
		{"256", 256, nil, ColorMode256},
		{"direct colour", 1 << 24, nil, ColorModeTrueColor},
		{"colorterm truecolor", 256, map[string]string{"COLORTERM": "truecolor"}, ColorModeTrueColor},
		{"colorterm 24bit", 16, map[string]string{"COLORTERM": "24bit"}, ColorModeTrueColor},
		{"no color override", 1 << 24, map[string]string{"NO_COLOR": "1"}, ColorMode16},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := detectColorMode(tc.colors, testEnv(tc.env)); got != tc.want {
				t.Fatalf("detectColorMode(%d) = %s, want %s", tc.colors, got, tc.want)
			}
		})
	}
}

func TestColorModeString(t *testing.T) {
	if ColorMode16.String() != "16" || ColorMode256.String() != "256" || ColorModeTrueColor.String() != "truecolor" {
		t.Fatal("unexpected colour mode names")
	}
}

func TestWoodsPaletteCoreColours(t *testing.T) {
	cases := []struct {
		role     ColorRole
		hex      string
		fallback uint8
	}{
		{RoleMapBackground, "#050605", 0},
		{RoleGrassBase, "#4E5F3A", 2},
		{RoleGrassLight, "#5C7044", 10},
		{RoleGrassDark, "#3C4A2E", 2},
		{RoleDirtRoad, "#6E5A42", 6},
		{RoleCanopyMid, "#3E5734", 2},
		{RoleCanopyLight, "#577345", 10},
		{RoleDeepWater, "#1C3440", 1},
		{RoleRiver, "#2F5A66", 3},
		{RoleShallowWater, "#4E7C7A", 3},
		{RoleWoodWall, "#6E4F35", 6},
		{RoleDoor, "#8A5A32", 6},
		{RoleStoneWall, "#6E6E66", 7},
		{RoleFreshBlood, "#B32020", 4},
		{RoleZombieSkin, "#7F8F62", 10},
		{RolePlayer, "#F2E6C9", 15},
		{RoleSteel, "#B8C0C0", 7},
		{RoleGoldRelic, "#C9A227", 14},
		{RoleUIBackground, "#0B0D0A", 0},
	}
	for _, tc := range cases {
		pc, ok := WoodsPalette.Lookup(tc.role)
		if !ok {
			t.Fatalf("palette is missing %s", tc.role)
		}
		if got := pc.RGB.Hex(); got != tc.hex {
			t.Fatalf("%s = %s, want %s", tc.role, got, tc.hex)
		}
		if pc.Fallback != tc.fallback {
			t.Fatalf("%s fallback = %d, want %d", tc.role, pc.Fallback, tc.fallback)
		}
	}
	if len(WoodsPalette) < 100 {
		t.Fatalf("core woods palette has only %d roles", len(WoodsPalette))
	}
	for role, pc := range WoodsPalette {
		if role == "" {
			t.Fatal("palette contains an empty role")
		}
		if pc.Fallback > 15 {
			t.Fatalf("%s fallback %d is outside the DOS palette", role, pc.Fallback)
		}
	}
}

func TestPaletteResolution(t *testing.T) {
	pc := PaletteColor{RGB: RGB{R: 0x4E, G: 0x5F, B: 0x3A}, Fallback: 2}

	if got := pc.Resolve(ColorMode16); got != DOSColor(2) || got != tcell.ColorGreen {
		t.Fatalf("16-colour resolve = %v, want DOS green", got)
	}
	if got := pc.Resolve(ColorModeTrueColor); !got.IsRGB() || got.Hex() != 0x4E5F3A {
		t.Fatalf("true-colour resolve = %v, want #4E5F3A", got)
	}
	want256 := tcell.PaletteColor(Nearest256(pc.RGB))
	if got := pc.Resolve(ColorMode256); got != want256 {
		t.Fatalf("256-colour resolve = %v, want %v", got, want256)
	}
	if got := DOSColor(99); got != tcell.ColorWhite {
		t.Fatalf("out-of-range DOS colour = %v, want white", got)
	}
}

func TestNearest256Palette(t *testing.T) {
	cases := []struct {
		c    RGB
		want int
	}{
		{RGB{0, 0, 0}, 0},
		{RGB{255, 0, 0}, 9},
		{RGB{95, 135, 175}, 67},
		{RGB{8, 8, 8}, 232},
	}
	for _, tc := range cases {
		if got := Nearest256(tc.c); got != tc.want {
			t.Fatalf("Nearest256(%v) = %d, want %d", tc.c, got, tc.want)
		}
	}
	for _, c := range []RGB{{1, 2, 3}, {200, 10, 90}, {250, 250, 250}} {
		if got := Nearest256(c); got < 0 || got > 255 {
			t.Fatalf("Nearest256(%v) = %d, outside 0..255", c, got)
		}
	}
}

func TestVaryRGBIsDeterministic(t *testing.T) {
	base := RGB{R: 100, G: 120, B: 80}
	first := VaryRGB(base, 3, 17, 42)
	if second := VaryRGB(base, 3, 17, 42); second != first {
		t.Fatalf("variation is not deterministic: %v vs %v", first, second)
	}
	differences := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			c := VaryRGB(base, 3, x, y)
			if c != base {
				differences++
			}
			if abs(int(c.R)-int(base.R)) > 21 ||
				abs(int(c.G)-int(base.G)) > 21 ||
				abs(int(c.B)-int(base.B)) > 21 {
				t.Fatalf("variation at %d,%d changed a channel by more than 8%%: %v -> %v", x, y, base, c)
			}
		}
	}
	if differences == 0 {
		t.Fatal("per-tile variation never changed a colour")
	}
}

func TestShadeAndBlendRGB(t *testing.T) {
	dim := ShadeRGB(RGB{255, 255, 255}, 0.40, 0.15)
	if dim.R != 102 || dim.G != 102 || dim.B != 102 {
		t.Fatalf("shaded white = %v, want approximately 102,102,102", dim)
	}
	if black := ShadeRGB(RGB{0, 0, 0}, 0.40, 0.15); black != (RGB{}) {
		t.Fatalf("shaded black = %v, want black", black)
	}
	red := ShadeRGB(RGB{255, 0, 0}, 1, 1)
	if red.R != 76 || red.G != 76 || red.B != 76 {
		t.Fatalf("desaturated red = %v, want luminance grey", red)
	}

	a, b := RGB{10, 20, 30}, RGB{110, 120, 130}
	if got := BlendRGB(a, b, 0); got != a {
		t.Fatalf("blend at 0 = %v, want %v", got, a)
	}
	if got := BlendRGB(a, b, 1); got != b {
		t.Fatalf("blend at 1 = %v, want %v", got, b)
	}
	if got := BlendRGB(a, b, 0.5); got != (RGB{60, 70, 80}) {
		t.Fatalf("blend at 0.5 = %v, want 60,70,80", got)
	}
}

func TestAllRecoveredPrototypesHavePaletteRoles(t *testing.T) {
	for id := 0; id < len(TileDefs); id++ {
		role, ok := TileColorRole(id)
		if !ok {
			t.Fatalf("tile %d (%s) has no palette role", id, TileDefs[id].Name)
		}
		if _, ok := WoodsPalette.Lookup(role); !ok {
			t.Fatalf("tile %d uses unknown role %s", id, role)
		}
	}
	for proto := 0; proto < len(ItemDefs); proto++ {
		role, ok := ItemColorRole(proto)
		if !ok {
			t.Fatalf("item %d (%s) has no palette role", proto, ItemDefs[proto].Name)
		}
		if _, ok := WoodsPalette.Lookup(role); !ok {
			t.Fatalf("item %d uses unknown role %s", proto, role)
		}
	}
	for proto := 0; proto < len(CritterDefs); proto++ {
		role, ok := CritterColorRole(proto)
		if !ok {
			t.Fatalf("critter %d (%s) has no palette role", proto, CritterDefs[proto].Name)
		}
		if _, ok := WoodsPalette.Lookup(role); !ok {
			t.Fatalf("critter %d uses unknown role %s", proto, role)
		}
	}
}

func TestWoodsPaletteRoleMappings(t *testing.T) {
	if role, ok := TileColorRole(2); !ok || role != RoleGrassBase {
		t.Fatalf("grass role = %s, want %s", role, RoleGrassBase)
	}
	if role, ok := TileColorRole(14); !ok || role != RoleWoodWall {
		t.Fatalf("wooden wall role = %s, want %s", role, RoleWoodWall)
	}
	if role, ok := TileColorRole(8); !ok || role != RoleDeepWater {
		t.Fatalf("deep water role = %s, want %s", role, RoleDeepWater)
	}
	if role, ok := ItemColorRole(5); !ok || role != RoleLeather {
		t.Fatalf("jacket role = %s, want %s", role, RoleLeather)
	}
	if role, ok := ItemColorRole(46); !ok || role != RoleHealthGood {
		t.Fatalf("first aid role = %s, want %s", role, RoleHealthGood)
	}
	if role, ok := CritterColorRole(0); !ok || role != RoleZombieSkin {
		t.Fatalf("zombie role = %s, want %s", role, RoleZombieSkin)
	}
	if role, ok := CritterColorRole(16); !ok || role != RolePlayer {
		t.Fatalf("player role = %s, want %s", role, RolePlayer)
	}
	if _, ok := TileColorRole(-1); ok {
		t.Fatal("negative tile id should not have a role")
	}
	if _, ok := ItemColorRole(len(ItemDefs)); ok {
		t.Fatal("out-of-range item should not have a role")
	}
	if _, ok := CritterColorRole(len(CritterDefs)); ok {
		t.Fatal("out-of-range critter should not have a role")
	}
}

func TestUIWorldColorFallbackAndShading(t *testing.T) {
	u := &UI{
		g:         &Game{Theme: 0, Level: 2},
		colorMode: ColorMode16,
	}
	if got := u.worldColor(RoleGrassBase, true, 2, 10, 10, true); got != DOSColor(2) {
		t.Fatalf("16-colour world colour = %v, want DOS green", got)
	}
	if got := u.worldColor(RoleGrassBase, true, 2, 10, 10, false); got != DOSColor(8) {
		t.Fatalf("16-colour remembered colour = %v, want DOS grey", got)
	}

	u.colorMode = ColorModeTrueColor
	visible := u.worldColor(RoleGrassBase, true, 2, 10, 10, true)
	if !visible.IsRGB() {
		t.Fatalf("true-colour world colour is not RGB: %v", visible)
	}
	remembered := u.worldColor(RoleGrassBase, true, 2, 10, 10, false)
	if !remembered.IsRGB() {
		t.Fatalf("true-colour remembered colour is not RGB: %v", remembered)
	}
	vr, vg, vb := visible.RGB()
	rr, rg, rb := remembered.RGB()
	if rr >= vr || rg >= vg || rb >= vb {
		t.Fatalf("remembered colour is not dimmer: %v vs %v", remembered, visible)
	}

	u.g.Theme = 1
	if got := u.worldColor(RoleGrassBase, true, 2, 10, 10, true); got != DOSColor(2) {
		t.Fatalf("non-woods theme should keep DOS fallback, got %v", got)
	}
}

func TestScreenBackgroundColour(t *testing.T) {
	if got := screenBackgroundColor(ColorMode16); got != DOSColor(0) {
		t.Fatalf("16-colour screen background = %v, want DOS black", got)
	}
	if got := screenBackgroundColor(ColorModeTrueColor); !got.IsRGB() || got.Hex() != 0x0B0D0A {
		t.Fatalf("true-colour screen background = %v, want #0B0D0A", got)
	}
	want256 := tcell.PaletteColor(Nearest256(RGB{0x0B, 0x0D, 0x0A}))
	if got := screenBackgroundColor(ColorMode256); got != want256 {
		t.Fatalf("256-colour screen background = %v, want %v", got, want256)
	}
}

func TestDrawGameUsesTrueColourBackground(t *testing.T) {
	g := NewGame(8)
	g.SetupNewGame(0, "tester")
	g.UpdateFOV()

	s := tcell.NewSimulationScreen("UTF-8")
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	s.SetSize(screenW, screenH)
	u := &UI{screen: s, g: g, colorMode: ColorModeTrueColor}
	u.drawGame()
	s.Show()

	cells, w, h := s.GetContents()
	foundPlayer := false
	for y := 0; y < h && !foundPlayer; y++ {
		for x := 0; x < w; x++ {
			cell := cells[y*w+x]
			if len(cell.Runes) == 0 || cell.Runes[0] != '@' {
				continue
			}
			foundPlayer = true
			_, bg, _ := cell.Style.Decompose()
			if !bg.IsRGB() {
				t.Fatalf("true-colour map background is not RGB: %v", bg)
			}
			r, g, b := bg.RGB()
			if r > 16 || g > 16 || b > 16 {
				t.Fatalf("map background is not near-black: %v", bg)
			}
			break
		}
	}
	s.Fini()
	if !foundPlayer {
		t.Fatal("player glyph not found")
	}
	_, topBG, _ := cells[0].Style.Decompose()
	if !topBG.IsRGB() || topBG.Hex() != 0x0B0D0A {
		t.Fatalf("message/status background = %v, want #0B0D0A", topBG)
	}
}

func TestDrawGameUsesPalette(t *testing.T) {
	g := NewGame(7)
	g.SetupNewGame(0, "tester")
	g.UpdateFOV()

	for _, mode := range []ColorMode{ColorModeTrueColor, ColorMode16} {
		s := tcell.NewSimulationScreen("UTF-8")
		if err := s.Init(); err != nil {
			t.Fatal(err)
		}
		s.SetSize(screenW, screenH)
		u := &UI{screen: s, g: g, colorMode: mode}
		u.drawGame()
		s.Show()
		cells, w, h := s.GetContents()
		found := false
		for y := 0; y < h && !found; y++ {
			for x := 0; x < w; x++ {
				cell := cells[y*w+x]
				if len(cell.Runes) == 0 || cell.Runes[0] != '@' {
					continue
				}
				found = true
				fg, _, _ := cell.Style.Decompose()
				switch mode {
				case ColorModeTrueColor:
					if !fg.IsRGB() {
						t.Fatalf("true-colour player glyph is not RGB: %v (theme %d, player %d,%d, cell %d,%d, style %v)", fg, g.Theme, g.Player.X, g.Player.Y, x, y, cell.Style)
					}
				case ColorMode16:
					if fg != DOSColor(15) {
						t.Fatalf("16-colour player glyph = %v, want DOS white", fg)
					}
				}
				break
			}
		}
		s.Fini()
		if !found {
			t.Fatalf("player glyph not drawn in %s mode", mode)
		}
	}
}
