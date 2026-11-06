package game

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// RGB is a 24-bit colour. The game keeps the recovered DOS colour index as a
// fallback so 16-colour terminals still get the original look.
type RGB struct {
	R, G, B uint8
}

func rgbFromHex(v uint32) RGB {
	return RGB{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v)}
}

func (c RGB) Hex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

func (c RGB) TCell() tcell.Color {
	return tcell.NewRGBColor(int32(c.R), int32(c.G), int32(c.B))
}

// ColorMode describes how much colour the terminal can display.
type ColorMode int

const (
	ColorMode16 ColorMode = iota
	ColorMode256
	ColorModeTrueColor
)

func (m ColorMode) String() string {
	switch m {
	case ColorModeTrueColor:
		return "truecolor"
	case ColorMode256:
		return "256"
	default:
		return "16"
	}
}

// DetectColorMode chooses a rendering mode from the tcell colour count and the
// conventional terminal environment variables.
func DetectColorMode(colors int) ColorMode {
	return detectColorMode(colors, os.Getenv)
}

func detectColorMode(colors int, getenv func(string) string) ColorMode {
	if getenv == nil {
		getenv = os.Getenv
	}
	if getenv("NO_COLOR") != "" || colors <= 0 {
		return ColorMode16
	}
	switch strings.ToLower(getenv("COLORTERM")) {
	case "truecolor", "24bit", "24-bit":
		return ColorModeTrueColor
	}
	switch {
	case colors >= 1<<24:
		return ColorModeTrueColor
	case colors >= 256:
		return ColorMode256
	default:
		return ColorMode16
	}
}

// ColorRole names a semantic slot in the woods palette. New tiles, items and
// critters should reference a role instead of a hard-coded colour.
type ColorRole string

const (
	RoleMapBackground ColorRole = "map.background"
	RoleGrassBase     ColorRole = "terrain.grass.base"
	RoleGrassLight    ColorRole = "terrain.grass.light"
	RoleGrassDark     ColorRole = "terrain.grass.dark"
	RoleLeafLitter    ColorRole = "terrain.leaf_litter"
	RoleDirtRoad      ColorRole = "terrain.dirt_road"
	RoleRoadLight     ColorRole = "terrain.road_light"
	RoleMud           ColorRole = "terrain.mud"
	RoleRock          ColorRole = "terrain.rock"
	RoleRockLight     ColorRole = "terrain.rock_light"
	RolePineCanopy    ColorRole = "vegetation.pine_canopy"
	RoleCanopyMid     ColorRole = "vegetation.canopy_mid"
	RoleCanopyLight   ColorRole = "vegetation.canopy_light"
	RoleBush          ColorRole = "vegetation.bush"
	RoleBramble       ColorRole = "vegetation.bramble"
	RoleThorn         ColorRole = "vegetation.thorn"
	RoleDeadTree      ColorRole = "vegetation.dead_tree"
	RoleOrchardLeaf   ColorRole = "vegetation.orchard_leaf"
	RoleStump         ColorRole = "vegetation.stump"
	RoleFlowerGold    ColorRole = "vegetation.flower_gold"
	RoleFlowerPink    ColorRole = "vegetation.flower_pink"
	RoleFlowerPale    ColorRole = "vegetation.flower_pale"

	RoleDeepWater    ColorRole = "water.deep"
	RoleRiver        ColorRole = "water.river"
	RoleShallowWater ColorRole = "water.shallow"
	RoleFoam         ColorRole = "water.foam"
	RoleMarsh        ColorRole = "water.marsh"
	RoleReeds        ColorRole = "water.reeds"
	RoleBridgeWood   ColorRole = "water.bridge_wood"
	RoleBridgeRope   ColorRole = "water.bridge_rope"
	RoleLanternMetal ColorRole = "light.lantern_metal"
	RoleLanternGlow  ColorRole = "light.lantern_glow"
	RoleFirelight    ColorRole = "light.firelight"
	RoleMoonlight    ColorRole = "light.moonlight"

	RoleWoodWallDark       ColorRole = "structure.wood_wall_dark"
	RoleWoodWall           ColorRole = "structure.wood_wall"
	RoleWoodLight          ColorRole = "structure.wood_light"
	RoleRoof               ColorRole = "structure.roof"
	RoleDoor               ColorRole = "structure.door"
	RoleWindowFrame        ColorRole = "structure.window_frame"
	RoleWarmWindowGlass    ColorRole = "structure.warm_window_glass"
	RoleFence              ColorRole = "structure.fence"
	RoleFurniture          ColorRole = "structure.furniture"
	RoleStoneWall          ColorRole = "structure.stone_wall"
	RoleStoneLight         ColorRole = "structure.stone_light"
	RoleGravestone         ColorRole = "structure.gravestone"
	RoleGraveDirt          ColorRole = "structure.grave_dirt"
	RoleCrypt              ColorRole = "structure.crypt"
	RoleIronFence          ColorRole = "structure.iron_fence"
	RoleBone               ColorRole = "structure.bone"
	RoleGraveFog           ColorRole = "structure.grave_fog"
	RolePlasterWall        ColorRole = "structure.plaster_wall"
	RoleThatch             ColorRole = "structure.thatch"
	RoleFadedPaint         ColorRole = "structure.faded_paint"
	RoleMarketCloth        ColorRole = "structure.market_cloth"
	RoleHay                ColorRole = "structure.hay"
	RoleAltarStone         ColorRole = "structure.altar_stone"
	RoleCandle             ColorRole = "structure.candle"
	RoleStainedGlassPurple ColorRole = "structure.stained_glass_purple"
	RoleStainedGlassBlue   ColorRole = "structure.stained_glass_blue"

	RoleFreshBlood   ColorRole = "gore.fresh_blood"
	RoleDarkBlood    ColorRole = "gore.dark_blood"
	RoleDriedBlood   ColorRole = "gore.dried_blood"
	RoleBloodGrass   ColorRole = "gore.blood_grass"
	RoleRottenFlesh  ColorRole = "gore.rotten_flesh"
	RoleZombieSkin   ColorRole = "gore.zombie_skin"
	RoleZombieShadow ColorRole = "gore.zombie_shadow"
	RoleBoneTeeth    ColorRole = "gore.bone_teeth"
	RoleGuts         ColorRole = "gore.guts"
	RolePoison       ColorRole = "gore.poison"
	RoleAcid         ColorRole = "gore.acid"
	RoleAcidDark     ColorRole = "gore.acid_dark"
	RoleFireEmber    ColorRole = "gore.fire_ember"
	RoleFlame        ColorRole = "gore.flame"
	RoleFlameCore    ColorRole = "gore.flame_core"
	RoleSmoke        ColorRole = "gore.smoke"
	RoleAsh          ColorRole = "gore.ash"
	RoleCharcoal     ColorRole = "gore.charcoal"

	RolePlayer                ColorRole = "creature.player"
	RolePlayerHighlight       ColorRole = "creature.player_highlight"
	RoleVillager              ColorRole = "creature.villager"
	RoleVillagerClothingBlue  ColorRole = "creature.villager_clothing_blue"
	RoleVillagerClothingBrown ColorRole = "creature.villager_clothing_brown"
	RoleVillagerClothingGrey  ColorRole = "creature.villager_clothing_grey"
	RoleDeer                  ColorRole = "creature.deer"
	RoleRabbit                ColorRole = "creature.rabbit"
	RoleFox                   ColorRole = "creature.fox"
	RoleBoar                  ColorRole = "creature.boar"
	RoleWolf                  ColorRole = "creature.wolf"
	RoleCrow                  ColorRole = "creature.crow"
	RoleFrog                  ColorRole = "creature.frog"
	RoleHostileEyes           ColorRole = "creature.hostile_eyes"

	RoleSteel        ColorRole = "item.steel"
	RoleIronDark     ColorRole = "item.iron_dark"
	RoleSilver       ColorRole = "item.silver"
	RoleGoldRelic    ColorRole = "item.gold_relic"
	RoleLeather      ColorRole = "item.leather"
	RoleCloth        ColorRole = "item.cloth"
	RoleApple        ColorRole = "item.apple"
	RoleBread        ColorRole = "item.bread"
	RoleCheese       ColorRole = "item.cheese"
	RoleMeat         ColorRole = "item.meat"
	RoleFish         ColorRole = "item.fish"
	RolePotionRed    ColorRole = "item.potion_red"
	RolePotionGreen  ColorRole = "item.potion_green"
	RolePotionBlue   ColorRole = "item.potion_blue"
	RolePotionPurple ColorRole = "item.potion_purple"
	RoleParchment    ColorRole = "item.parchment"
	RoleKey          ColorRole = "item.key"

	RoleUIBackground        ColorRole = "ui.background"
	RolePanelBackground     ColorRole = "ui.panel_background"
	RolePanelBorder         ColorRole = "ui.panel_border"
	RolePrimaryText         ColorRole = "ui.primary_text"
	RoleDimText             ColorRole = "ui.dim_text"
	RoleAccentGreen         ColorRole = "ui.accent_green"
	RoleAccentGold          ColorRole = "ui.accent_gold"
	RoleDanger              ColorRole = "ui.danger"
	RoleBloodAccent         ColorRole = "ui.blood_accent"
	RoleSelectionBackground ColorRole = "ui.selection_background"
	RoleSelectionText       ColorRole = "ui.selection_text"
	RoleHealthGood          ColorRole = "ui.health_good"
	RoleHealthWounded       ColorRole = "ui.health_wounded"
	RoleHealthCritical      ColorRole = "ui.health_critical"
	RoleHintCyan            ColorRole = "ui.hint_cyan"
)

// PaletteColor is one semantic colour plus the recovered 16-colour index used
// when the terminal cannot display true colour.
type PaletteColor struct {
	RGB      RGB
	Fallback uint8
}

func paletteColor(hex uint32, fallback uint8) PaletteColor {
	return PaletteColor{RGB: rgbFromHex(hex), Fallback: fallback}
}

// Palette maps semantic roles to colours.
type Palette map[ColorRole]PaletteColor

// WoodsPalette is the core woods palette from
// documentation/plans/woods/color-palette.md. Per-level accents are deliberately
// not included yet.
var WoodsPalette = Palette{
	RoleMapBackground: paletteColor(0x050605, 0),
	RoleGrassBase:     paletteColor(0x4E5F3A, 2),
	RoleGrassLight:    paletteColor(0x5C7044, 10),
	RoleGrassDark:     paletteColor(0x3C4A2E, 2),
	RoleLeafLitter:    paletteColor(0x5A4A32, 6),
	RoleDirtRoad:      paletteColor(0x6E5A42, 6),
	RoleRoadLight:     paletteColor(0x7E6A50, 6),
	RoleMud:           paletteColor(0x4A3C2E, 6),
	RoleRock:          paletteColor(0x5C5E56, 8),
	RoleRockLight:     paletteColor(0x7A7C72, 7),
	RolePineCanopy:    paletteColor(0x2D4029, 2),
	RoleCanopyMid:     paletteColor(0x3E5734, 2),
	RoleCanopyLight:   paletteColor(0x577345, 10),
	RoleBush:          paletteColor(0x556B3E, 2),
	RoleBramble:       paletteColor(0x4A3A2A, 6),
	RoleThorn:         paletteColor(0x8C7A5A, 6),
	RoleDeadTree:      paletteColor(0x4A4038, 8),
	RoleOrchardLeaf:   paletteColor(0x5E7A3E, 10),
	RoleStump:         paletteColor(0x5A4230, 6),
	RoleFlowerGold:    paletteColor(0xC9B45A, 14),
	RoleFlowerPink:    paletteColor(0xB87A8A, 13),
	RoleFlowerPale:    paletteColor(0xE0D8B0, 15),

	RoleDeepWater:    paletteColor(0x1C3440, 1),
	RoleRiver:        paletteColor(0x2F5A66, 3),
	RoleShallowWater: paletteColor(0x4E7C7A, 3),
	RoleFoam:         paletteColor(0xA6C4BC, 11),
	RoleMarsh:        paletteColor(0x3E4A34, 2),
	RoleReeds:        paletteColor(0x7A7A42, 6),
	RoleBridgeWood:   paletteColor(0x6E4F35, 6),
	RoleBridgeRope:   paletteColor(0x8C7A5A, 6),
	RoleLanternMetal: paletteColor(0x4A4A44, 8),
	RoleLanternGlow:  paletteColor(0xFFD27A, 14),
	RoleFirelight:    paletteColor(0xE08A3A, 12),
	RoleMoonlight:    paletteColor(0x223028, 8),

	RoleWoodWallDark:       paletteColor(0x4A3324, 6),
	RoleWoodWall:           paletteColor(0x6E4F35, 6),
	RoleWoodLight:          paletteColor(0x8A6A48, 6),
	RoleRoof:               paletteColor(0x5A3E2A, 6),
	RoleDoor:               paletteColor(0x8A5A32, 6),
	RoleWindowFrame:        paletteColor(0x6E4F35, 6),
	RoleWarmWindowGlass:    paletteColor(0xD9A441, 14),
	RoleFence:              paletteColor(0x6A5238, 6),
	RoleFurniture:          paletteColor(0x7A5A3A, 6),
	RoleStoneWall:          paletteColor(0x6E6E66, 7),
	RoleStoneLight:         paletteColor(0x8A8A80, 7),
	RoleGravestone:         paletteColor(0x8C8C84, 7),
	RoleGraveDirt:          paletteColor(0x4A3E32, 6),
	RoleCrypt:              paletteColor(0x4E4E4A, 8),
	RoleIronFence:          paletteColor(0x3E4648, 8),
	RoleBone:               paletteColor(0xD6CBB0, 7),
	RoleGraveFog:           paletteColor(0x8A9A92, 8),
	RolePlasterWall:        paletteColor(0xB8A88A, 7),
	RoleThatch:             paletteColor(0xA88A52, 6),
	RoleFadedPaint:         paletteColor(0x7A6A5A, 6),
	RoleMarketCloth:        paletteColor(0xA84A3A, 4),
	RoleHay:                paletteColor(0xC9A85A, 6),
	RoleAltarStone:         paletteColor(0x7A7A70, 7),
	RoleCandle:             paletteColor(0xF2C46A, 14),
	RoleStainedGlassPurple: paletteColor(0x7A4A8A, 5),
	RoleStainedGlassBlue:   paletteColor(0x3A6A8A, 1),

	RoleFreshBlood:   paletteColor(0xB32020, 4),
	RoleDarkBlood:    paletteColor(0x7A1A1A, 4),
	RoleDriedBlood:   paletteColor(0x5A2420, 4),
	RoleBloodGrass:   paletteColor(0x6B2A22, 4),
	RoleRottenFlesh:  paletteColor(0x7A8A5A, 10),
	RoleZombieSkin:   paletteColor(0x7F8F62, 10),
	RoleZombieShadow: paletteColor(0x4A553A, 2),
	RoleBoneTeeth:    paletteColor(0xD6CBB0, 7),
	RoleGuts:         paletteColor(0x8A4A4A, 4),
	RolePoison:       paletteColor(0x7A8A3A, 10),
	RoleAcid:         paletteColor(0x8FBF3F, 10),
	RoleAcidDark:     paletteColor(0x5A7A2A, 2),
	RoleFireEmber:    paletteColor(0xE06A2A, 12),
	RoleFlame:        paletteColor(0xF2A33A, 14),
	RoleFlameCore:    paletteColor(0xFFD27A, 14),
	RoleSmoke:        paletteColor(0x5A5A55, 8),
	RoleAsh:          paletteColor(0x7A7A72, 8),
	RoleCharcoal:     paletteColor(0x222222, 8),

	RolePlayer:                paletteColor(0xF2E6C9, 15),
	RolePlayerHighlight:       paletteColor(0xFFD27A, 14),
	RoleVillager:              paletteColor(0xC9B48A, 7),
	RoleVillagerClothingBlue:  paletteColor(0x6A7A8A, 1),
	RoleVillagerClothingBrown: paletteColor(0x8A5A3A, 6),
	RoleVillagerClothingGrey:  paletteColor(0x7A6A5A, 8),
	RoleDeer:                  paletteColor(0x8A6A4A, 6),
	RoleRabbit:                paletteColor(0xB8A88A, 7),
	RoleFox:                   paletteColor(0xB55A2A, 6),
	RoleBoar:                  paletteColor(0x5A4A3A, 6),
	RoleWolf:                  paletteColor(0x6A6A66, 8),
	RoleCrow:                  paletteColor(0x1E1E22, 8),
	RoleFrog:                  paletteColor(0x5A8A4A, 10),
	RoleHostileEyes:           paletteColor(0xC9A227, 14),

	RoleSteel:        paletteColor(0xB8C0C0, 7),
	RoleIronDark:     paletteColor(0x6A7070, 8),
	RoleSilver:       paletteColor(0xC0C8C8, 7),
	RoleGoldRelic:    paletteColor(0xC9A227, 14),
	RoleLeather:      paletteColor(0x7A4E2E, 6),
	RoleCloth:        paletteColor(0x8A7A5A, 6),
	RoleApple:        paletteColor(0xA83A32, 4),
	RoleBread:        paletteColor(0xC9A26A, 6),
	RoleCheese:       paletteColor(0xD9B85A, 14),
	RoleMeat:         paletteColor(0x8A3A2A, 4),
	RoleFish:         paletteColor(0x8AA8B0, 3),
	RolePotionRed:    paletteColor(0xA83A32, 4),
	RolePotionGreen:  paletteColor(0x4A8A3A, 10),
	RolePotionBlue:   paletteColor(0x3A6A8A, 1),
	RolePotionPurple: paletteColor(0x7A4A8A, 5),
	RoleParchment:    paletteColor(0xD8C9A3, 7),
	RoleKey:          paletteColor(0xC9A227, 14),

	RoleUIBackground:        paletteColor(0x0B0D0A, 0),
	RolePanelBackground:     paletteColor(0x141812, 0),
	RolePanelBorder:         paletteColor(0x3E4A38, 8),
	RolePrimaryText:         paletteColor(0xD6CBB0, 7),
	RoleDimText:             paletteColor(0x8A8F7A, 8),
	RoleAccentGreen:         paletteColor(0x7A9A5A, 10),
	RoleAccentGold:          paletteColor(0xC9A227, 14),
	RoleDanger:              paletteColor(0xB23A2E, 12),
	RoleBloodAccent:         paletteColor(0xB32020, 4),
	RoleSelectionBackground: paletteColor(0x4A5A3A, 2),
	RoleSelectionText:       paletteColor(0xF2E6C9, 15),
	RoleHealthGood:          paletteColor(0x6A9A4A, 10),
	RoleHealthWounded:       paletteColor(0xC9A227, 14),
	RoleHealthCritical:      paletteColor(0xB23A2E, 12),
	RoleHintCyan:            paletteColor(0x6A9A9A, 11),
}

func (p Palette) Lookup(role ColorRole) (PaletteColor, bool) {
	c, ok := p[role]
	return c, ok
}

// Resolve converts a palette colour to a tcell colour for the active mode.
func (pc PaletteColor) Resolve(mode ColorMode) tcell.Color {
	return ResolveRGB(pc.RGB, pc.Fallback, mode)
}

// ResolveRGB converts an RGB value to a tcell colour, using the supplied
// fallback index in 16-colour mode.
func ResolveRGB(c RGB, fallback uint8, mode ColorMode) tcell.Color {
	switch mode {
	case ColorModeTrueColor:
		return c.TCell()
	case ColorMode256:
		return tcell.PaletteColor(Nearest256(c))
	default:
		return DOSColor(fallback)
	}
}

// DOSColor maps a recovered DOS/CP437 colour index to the tcell 16-colour
// palette. Indices outside the original table fall back to white.
func DOSColor(i uint8) tcell.Color {
	if int(i) < len(dosPalette) {
		return dosPalette[i]
	}
	return tcell.ColorWhite
}

func pal(i uint8) tcell.Color { return DOSColor(i) }

// screenBackgroundColor returns the true background used by the UI and by
// cells that have no explicit background. Named tcell colours such as
// ColorBlack are terminal-palette entries, so a terminal theme can render them
// as grey or blue; true-colour and 256-colour modes use the palette RGB value
// instead.
func screenBackgroundColor(mode ColorMode) tcell.Color {
	if mode == ColorMode16 {
		return pal(0)
	}
	pc, ok := WoodsPalette.Lookup(RoleUIBackground)
	if !ok {
		return pal(0)
	}
	return ResolveRGB(pc.RGB, pc.Fallback, mode)
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

func clampByte(v float64) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 255 {
		return 255
	}
	return uint8(math.Round(v))
}

// tileColorHash produces stable per-tile noise for colour variation.
func tileColorHash(level, x, y int) uint32 {
	h := uint32(level)*0x9E3779B1 ^ uint32(x)*0x85EBCA77 ^ uint32(y)*0xC2B2AE3D
	h ^= h >> 16
	h *= 0x7FEB352D
	h ^= h >> 15
	h *= 0x846CA68B
	h ^= h >> 16
	return h
}

// VaryRGB applies a small deterministic brightness change of roughly ±8% to a
// tile colour. The same level and coordinate always produce the same result.
func VaryRGB(c RGB, level, x, y int) RGB {
	delta := int(tileColorHash(level, x, y)%17) - 8
	scale := 1 + float64(delta)/100
	return RGB{
		R: clampByte(float64(c.R) * scale),
		G: clampByte(float64(c.G) * scale),
		B: clampByte(float64(c.B) * scale),
	}
}

// ShadeRGB dims and slightly desaturates a colour, used for remembered tiles.
func ShadeRGB(c RGB, brightness, desaturate float64) RGB {
	if brightness < 0 {
		brightness = 0
	}
	if desaturate < 0 {
		desaturate = 0
	}
	if desaturate > 1 {
		desaturate = 1
	}
	lum := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
	r := float64(c.R) + (lum-float64(c.R))*desaturate
	g := float64(c.G) + (lum-float64(c.G))*desaturate
	b := float64(c.B) + (lum-float64(c.B))*desaturate
	return RGB{
		R: clampByte(r * brightness),
		G: clampByte(g * brightness),
		B: clampByte(b * brightness),
	}
}

// BlendRGB mixes two colours; t=0 returns a and t=1 returns b.
func BlendRGB(a, b RGB, t float64) RGB {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return RGB{
		R: clampByte(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: clampByte(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: clampByte(float64(a.B) + (float64(b.B)-float64(a.B))*t),
	}
}

// xterm256Palette is the standard xterm 256-colour palette. It is used for
// explicit 256-colour fallback instead of leaving conversion to the terminal.
var xterm256Palette = buildXterm256Palette()

func buildXterm256Palette() [256]RGB {
	var p [256]RGB
	base := [16]uint32{
		0x000000, 0x800000, 0x008000, 0x808000,
		0x000080, 0x800080, 0x008080, 0xC0C0C0,
		0x808080, 0xFF0000, 0x00FF00, 0xFFFF00,
		0x0000FF, 0xFF00FF, 0x00FFFF, 0xFFFFFF,
	}
	for i, v := range base {
		p[i] = rgbFromHex(v)
	}
	levels := [6]uint8{0, 95, 135, 175, 215, 255}
	idx := 16
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				p[idx] = RGB{R: levels[r], G: levels[g], B: levels[b]}
				idx++
			}
		}
	}
	for i := 0; i < 24; i++ {
		v := uint8(8 + i*10)
		p[232+i] = RGB{R: v, G: v, B: v}
	}
	return p
}

// Nearest256 returns the closest entry in the xterm 256-colour palette.
func Nearest256(c RGB) int {
	best, bestDist := 0, int(^uint(0)>>1)
	for i, p := range xterm256Palette {
		dr := int(c.R) - int(p.R)
		dg := int(c.G) - int(p.G)
		db := int(c.B) - int(p.B)
		d := dr*dr + dg*dg + db*db
		if d < bestDist {
			best, bestDist = i, d
			if d == 0 {
				return i
			}
		}
	}
	return best
}
