package game

import (
	"fmt"
	"unicode"
)

const (
	mapWidth  = 80
	mapHeight = 46
)

type Map struct {
	W, H      int
	Tiles     []uint8
	Seen      []bool
	Visible   []bool
	Blood     []uint8
	Fire      []int
	Acid      []int
	Smoke     []int
	Theme     int
	Name      string
	LandmarkX int
	LandmarkY int
}

func NewMap(w, h int) *Map {
	m := &Map{W: w, H: h}
	n := w * h
	m.Tiles = make([]uint8, n)
	m.Seen = make([]bool, n)
	m.Visible = make([]bool, n)
	m.Blood = make([]uint8, n)
	m.Fire = make([]int, n)
	m.Acid = make([]int, n)
	m.Smoke = make([]int, n)
	return m
}

func (m *Map) idx(x, y int) int { return y*m.W + x }

func (m *Map) InBounds(x, y int) bool { return x >= 0 && y >= 0 && x < m.W && y < m.H }

func (m *Map) Tile(x, y int) uint8 {
	if !m.InBounds(x, y) {
		return 0
	}
	return m.Tiles[m.idx(x, y)]
}

func (m *Map) SetTile(x, y int, t uint8) {
	if m.InBounds(x, y) {
		m.Tiles[m.idx(x, y)] = t
	}
}

func (m *Map) Def(x, y int) *TileDef {
	t := int(m.Tile(x, y))
	if t < 0 || t >= len(TileDefs) {
		return &TileDefs[0]
	}
	return &TileDefs[t]
}

func (m *Map) BlocksMove(x, y int) bool {
	d := m.Def(x, y)
	return d.BlockMove
}

func (m *Map) BlocksSight(x, y int) bool {
	d := m.Def(x, y)
	return d.BlockSight
}

func (m *Map) IsPassable(x, y int) bool { return m.InBounds(x, y) && !m.BlocksMove(x, y) }

func (m *Map) Free(x, y int) bool {
	return m.IsPassable(x, y)
}

func (m *Map) FirstFreeNear(x, y int) (int, int) {
	if m.Free(x, y) {
		return x, y
	}
	for r := 1; r < 20; r++ {
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				if abs(dx) != r && abs(dy) != r {
					continue
				}
				nx, ny := x+dx, y+dy
				if m.Free(nx, ny) {
					return nx, ny
				}
			}
		}
	}
	return x, y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m *Map) RandomFree(rng interface{ Intn(int) int }) (int, int) {
	for i := 0; i < 10000; i++ {
		x := rng.Intn(m.W)
		y := rng.Intn(m.H)
		if m.Free(x, y) {
			return x, y
		}
	}
	return m.W / 2, m.H / 2
}

// LineOfSight reports whether a straight line from (x0,y0) to (x1,y1) is not blocked.
func (m *Map) LineOfSight(x0, y0, x1, y1 int) bool {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	x, y := x0, y0
	for {
		if x == x1 && y == y1 {
			return true
		}
		if !(x == x0 && y == y0) && m.BlocksSight(x, y) {
			return false
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// Line returns all points on a Bresenham line, excluding the start.
func (m *Map) Line(x0, y0, x1, y1 int) [][2]int {
	var out [][2]int
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	x, y := x0, y0
	for {
		if x == x1 && y == y1 {
			return out
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
		out = append(out, [2]int{x, y})
	}
}

// FOV computes visible tiles using recursive shadowcasting (8 octants).
// This is the classic Bjorn Bergstrom algorithm credited in the original game.
func (m *Map) FOV(cx, cy, radius int) {
	for i := range m.Visible {
		m.Visible[i] = false
	}
	if radius < 1 {
		radius = 1
	}
	m.setVisible(cx, cy)
	octants := [8][4]int{
		{1, 0, 0, -1}, {0, 1, -1, 0}, {0, -1, -1, 0}, {-1, 0, 0, -1},
		{-1, 0, 0, 1}, {0, -1, 1, 0}, {0, 1, 1, 0}, {1, 0, 0, 1},
	}
	for _, oct := range octants {
		m.castLight(cx, cy, 1, 1.0, 0.0, radius, oct)
	}
}

func (m *Map) setVisible(x, y int) {
	if !m.InBounds(x, y) {
		return
	}
	m.Visible[m.idx(x, y)] = true
	m.Seen[m.idx(x, y)] = true
}

func (m *Map) castLight(cx, cy, row int, start, end float64, radius int, oct [4]int) {
	if start < end {
		return
	}
	newStart := start
	for j := row; j <= radius; j++ {
		dx := -j - 1
		dy := -j
		blocked := false
		for dx <= 0 {
			mx := cx + dx*oct[0] + dy*oct[1]
			my := cy + dx*oct[2] + dy*oct[3]
			lSlope := (float64(dx) - 0.5) / (float64(dy) + 0.5)
			rSlope := (float64(dx) + 0.5) / (float64(dy) - 0.5)
			if start < rSlope {
				dx++
				continue
			}
			if end > lSlope {
				break
			}
			if dx*dx+dy*dy <= radius*radius {
				m.setVisible(mx, my)
			}
			if blocked {
				if m.BlocksSight(mx, my) {
					newStart = rSlope
				} else {
					blocked = false
					start = newStart
				}
			} else if m.BlocksSight(mx, my) && j < radius {
				blocked = true
				m.castLight(cx, cy, j+1, start, lSlope, radius, oct)
				newStart = rSlope
			}
			dx++
		}
		if blocked {
			break
		}
	}
}

// Reachable computes the set of tiles reachable from (sx,sy).  Closed doors
// count as traversable because the player can open them.
func (m *Map) Reachable(sx, sy int) []int {
	seen := make([]bool, m.W*m.H)
	start := m.idx(sx, sy)
	seen[start] = true
	q := []int{start}
	for head := 0; head < len(q); head++ {
		cur := q[head]
		cx, cy := cur%m.W, cur/m.W
		for _, d := range dirs8 {
			nx, ny := cx+d[0], cy+d[1]
			if !m.InBounds(nx, ny) {
				continue
			}
			ni := m.idx(nx, ny)
			if seen[ni] {
				continue
			}
			if m.BlocksMove(nx, ny) && m.Def(nx, ny).Type != 3 {
				continue
			}
			seen[ni] = true
			q = append(q, ni)
		}
	}
	return q
}

// FarthestReachable returns the reachable floor tile farthest from (sx,sy).
func (m *Map) FarthestReachable(sx, sy int) (int, int) {
	bestX, bestY, bestD := -1, -1, -1
	for _, idx := range m.Reachable(sx, sy) {
		x, y := idx%m.W, idx/m.W
		if m.BlocksMove(x, y) {
			continue
		}
		d := abs(x-sx) + abs(y-sy)
		if d > bestD {
			bestD = d
			bestX, bestY = x, y
		}
	}
	return bestX, bestY
}

func (m *Map) Bresenham(x0, y0, x1, y1 int) [][2]int {
	return m.Line(x0, y0, x1, y1)
}

// GenerateLevel dispatches woods levels to the outdoor generator and all other
// themes to the recovered room-section generator.
func (g *Game) GenerateLevel(level int) {
	if levelTheme(level) == 0 {
		g.generateWoodsLevel(level)
		return
	}
	g.generateSectionLevel(level)
}

// generateSectionLevel keeps the recovered room-section generator for caves,
// the temple and the time portal.
func (g *Game) generateSectionLevel(level int) {
	m := NewMap(mapWidth, mapHeight)
	theme := levelTheme(level)
	m.Theme = theme
	m.Name = LevelThemes[theme]
	if theme == 3 {
		m.Name = "time portal"
	}
	// Fill with the theme's base terrain first.  The section templates use
	// spaces for "not part of this cube section"; those must become ordinary
	// terrain, not black void, otherwise the level looks full of holes.
	base := baseTile(theme)
	for i := range m.Tiles {
		m.Tiles[i] = base
	}
	sectioned := make([]bool, m.W*m.H)
	sectionRows := (m.H + 9) / 10
	// Place 2 columns of 40x10 sections, clipping the final row to the map.
	for row := 0; row < sectionRows; row++ {
		rowsHere := min(10, m.H-row*10)
		for col := 0; col < 2; col++ {
			sec := g.Roll(len(MapSections))
			txt := MapSections[sec]
			if g.Roll(2) == 1 {
				txt = flipH(txt)
			}
			if g.Roll(3) == 0 {
				txt = flipV(txt)
			}
			for y := 0; y < rowsHere; y++ {
				for x := 0; x < 40; x++ {
					ch := txt[y*40+x]
					mx, my := col*40+x, row*10+y
					if ch != ' ' {
						sectioned[m.idx(mx, my)] = true
					}
					m.SetTile(mx, my, tileForChar(ch, theme, g))
				}
			}
		}
	}
	// Ensure connections between section centers: carve simple corridors.
	for row := 0; row < sectionRows; row++ {
		cy := min(row*10+5, m.H-2)
		carveH(m, 0, m.W-1, cy, floorTile(theme))
		if row < sectionRows-1 {
			nextCy := min((row+1)*10+5, m.H-2)
			cx := m.W / 2
			carveV(m, cx, cy, nextCy, floorTile(theme))
		}
	}
	// Border walls.
	for x := 0; x < m.W; x++ {
		if m.IsPassable(x, 0) {
			m.SetTile(x, 0, wallTile(theme))
		}
		if m.IsPassable(x, m.H-1) {
			m.SetTile(x, m.H-1, wallTile(theme))
		}
	}
	for y := 0; y < m.H; y++ {
		if m.IsPassable(0, y) {
			m.SetTile(0, y, wallTile(theme))
		}
		if m.IsPassable(m.W-1, y) {
			m.SetTile(m.W-1, y, wallTile(theme))
		}
	}
	// Start position: first free tile near the upper-left area.
	px, py := m.FirstFreeNear(4, 4)
	if g.Player == nil {
		g.Player = &Critter{Proto: -1, Player: true, Align: AlignPlayer, Class: 4}
	}
	g.Player.X, g.Player.Y = px, py
	// Make sure the starting room is connected to the main corridor grid.
	carveH(m, min(px, 20), max(px, 20), py, floorTile(theme))
	carveV(m, 20, min(py, 5), max(py, 5), floorTile(theme))
	// Put the portal on the farthest reachable tile so it can always be found.
	qx, qy := m.FarthestReachable(px, py)
	if qx < 0 {
		qx, qy = m.FirstFreeNear(m.W-6, m.H-6)
	}
	m.SetTile(qx, qy, portalTile(theme))
	// Some water/fire/acid pools for the theme.
	if theme == 1 || theme == 2 {
		for i := 0; i < 10+g.Roll(20); i++ {
			x, y := m.RandomFree(g.RNG)
			if abs(x-px)+abs(y-py) < 12 {
				continue
			}
			switch g.Roll(3) {
			case 0:
				m.Acid[m.idx(x, y)] = 20 + g.Roll(40)
			case 1:
				m.Fire[m.idx(x, y)] = 10 + g.Roll(30)
			}
		}
	}
	if theme == 0 {
		// Overgrow the outdoor filler so the woods look like woods instead of
		// an empty plain.  Only outside the cube sections, and only passable
		// or shallow-water decorations so connectivity is preserved.
		for i := 0; i < 1400; i++ {
			x := g.Roll(m.W)
			y := g.Roll(m.H)
			if !m.InBounds(x, y) || sectioned[m.idx(x, y)] {
				continue
			}
			if m.Tile(x, y) != base {
				continue
			}
			if abs(x-px)+abs(y-py) < 6 || (x == qx && y == qy) {
				continue
			}
			switch {
			case g.Roll(100) < 42:
				m.SetTile(x, y, 3) // tree
			case g.Roll(100) < 14:
				m.SetTile(x, y, 4) // big tree
			case g.Roll(100) < 18:
				m.SetTile(x, y, uint8(50+g.Roll(3))) // plant
			case g.Roll(100) < 12:
				m.SetTile(x, y, uint8(6+g.Roll(2))) // swamp / shallow water
			}
		}
	}
	g.finishLevel(m, level, theme)
}

func (g *Game) finishLevel(m *Map, level, theme int) {
	g.Map = m
	g.Critters = nil
	g.Items = nil
	g.Level = level
	g.Options["optGameLevelID"] = level
	g.Theme = theme
	g.Weather = WeatherClear
	g.SpawnPopulation()
	if theme == 0 && level >= 0 && level < len(g.WoodLoot) {
		g.spawnSignatureLoot(level)
	}
	// Starting hint.
	if g.Options.Bool("optGameHints") {
		h := ActiveHints[g.HintsShown%len(ActiveHints)]
		g.HintsShown++
		g.Msg("hint: "+h, 10)
	}
	g.Msg(fmt.Sprintf("=== start of generating level === you are on %s level %d ===", m.Name, level), 7)
}

func levelTheme(level int) int {
	switch {
	case level < 9:
		return 0
	case level < 18:
		return 1
	case level < 27:
		return 2
	default:
		return 3
	}
}

func floorTile(theme int) uint8 {
	switch theme {
	case 0:
		return 2 // grass
	case 1:
		return 9 // stone floor
	case 2:
		return 20 // golden floor
	default:
		return 25 // plasteel floor
	}
}

func wallTile(theme int) uint8 {
	switch theme {
	case 0:
		return 14 // wooden wall
	case 1:
		return 10 // stone wall
	case 2:
		return 10 // stone wall
	default:
		return 28 // plasteel wall
	}
}

func baseTile(theme int) uint8 {
	switch theme {
	case 0:
		return 2 // grass
	case 1, 2:
		return 10 // solid rock / stone wall
	default:
		return 28 // plasteel wall
	}
}

func portalTile(theme int) uint8 {
	switch theme {
	case 0:
		return 23 // stone stairs
	case 1:
		return 23
	case 2:
		return 24 // time portal
	default:
		return 59 // elevator
	}
}

func shallowWaterTile() uint8 { return 7 }

func tileForChar(ch byte, theme int, g *Game) uint8 {
	switch ch {
	case '#':
		return wallTile(theme)
	case '.':
		return floorTile(theme)
	case '+', '|':
		if theme == 3 {
			return 33 // plasteel door
		}
		if theme == 2 {
			return 21 // golden door
		}
		return 15 // wooden door
	case 'I':
		if theme == 3 {
			return 28
		}
		return 12 // stone pillar
	case 'o':
		if theme == 0 {
			return 2
		}
		return floorTile(theme)
	case '@':
		return floorTile(theme)
	case ' ':
		return baseTile(theme)
	default:
		return floorTile(theme)
	}
}

func carveH(m *Map, x0, x1, y int, tile uint8) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	for x := x0; x <= x1; x++ {
		m.SetTile(x, y, tile)
	}
}

func carveV(m *Map, x, y0, y1 int, tile uint8) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	for y := y0; y <= y1; y++ {
		m.SetTile(x, y, tile)
	}
}

func flipH(txt string) string {
	out := make([]byte, len(txt))
	for y := 0; y < 10; y++ {
		for x := 0; x < 40; x++ {
			out[y*40+x] = txt[y*40+(39-x)]
		}
	}
	return string(out)
}

func flipV(txt string) string {
	out := make([]byte, len(txt))
	for y := 0; y < 10; y++ {
		for x := 0; x < 40; x++ {
			out[y*40+x] = txt[(9-y)*40+x]
		}
	}
	return string(out)
}

func (g *Game) SpawnPopulation() {
	theme := g.Theme
	pool := spawnPool(theme, g.Level)
	count := 8 + g.Level/2
	if count > 40 {
		count = 40
	}
	for i := 0; i < count; i++ {
		x, y := g.Map.RandomFree(g.RNG)
		if abs(x-g.Player.X)+abs(y-g.Player.Y) < 14 {
			continue
		}
		proto := pool[g.Roll(len(pool))]
		g.SpawnCritter(proto, x, y)
	}
	if theme == 0 {
		for i := 0; i < 2+g.Roll(3); i++ {
			proto := villagerProtos[g.Roll(len(villagerProtos))]
			g.spawnRandomCritter(proto, 10)
		}
		for i := 0; i < 3+g.Roll(4); i++ {
			proto := neutralWoodsAnimals[g.Roll(len(neutralWoodsAnimals))]
			g.spawnRandomCritter(proto, 8)
		}
	}
	// Ground items.
	itemPool := itemSpawnPool(theme, g.Level)
	for i := 0; i < 8+g.Roll(10); i++ {
		x, y := g.Map.RandomFree(g.RNG)
		if abs(x-g.Player.X)+abs(y-g.Player.Y) < 4 {
			continue
		}
		proto := itemPool[g.Roll(len(itemPool))]
		it := &Item{Proto: proto, X: x, Y: y}
		if d := it.Def(); d != nil && d.Kind == KindMedical {
			it.Charges = d.Charges
		}
		g.PlaceGroundItemNear(it, x, y)
	}
	// Guaranteed supplies near start.
	for i, proto := range []int{46, 48, 48} {
		x, y := g.Map.FirstFreeNear(g.Player.X+2+i, g.Player.Y+1)
		g.PlaceGroundItemNear(&Item{Proto: proto, X: -1, Y: -1, Charges: ItemDefs[proto].Charges}, x, y)
	}
}

func spawnPool(theme, level int) []int {
	switch theme {
	case 0:
		p := []int{0, 0, 1, 2, 3, 4, 14, 15}
		p = append(p, villagerProto, rabbitProto, foxProto, boarProto, crowProto, frogProto)
		if level > 3 {
			p = append(p, 5, 6, 14)
		}
		if level == 3 || level == 6 {
			p = append(p, villagerProto, militiaProto, villagerProto)
		}
		return p
	case 1:
		p := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		if level > 12 {
			p = append(p, 6, 7, 9)
		}
		return p
	case 2:
		p := []int{6, 7, 8, 9, 10, 11, 12, 13, 6}
		if level > 20 {
			p = append(p, 7, 9)
		}
		return p
	default:
		return []int{7, 8, 9, 10, 11, 12, 13}
	}
}

func itemSpawnPool(theme, level int) []int {
	base := []int{46, 47, 48, 8, 9, 10, 11, 12, 5}
	if level > 2 {
		base = append(base, 6, 15, 16, 17, 18, 19, 7, ArmourVialProto)
	}
	if level > 4 {
		base = append(base, 41)
	}
	if level > 6 {
		base = append(base, 13, 14, 18, 19, 42)
	}
	if level > 12 {
		base = append(base, 13, 14, 17, 18, 19, 5, 6, 7)
	}
	if theme == 2 {
		base = append(base, 14, 15, 16, 17)
	}
	return filterAvailableItems(base)
}

func (g *Game) spawnRandomCritter(proto, minDist int) *Critter {
	for i := 0; i < 200; i++ {
		x, y := g.Map.RandomFree(g.RNG)
		if abs(x-g.Player.X)+abs(y-g.Player.Y) < minDist {
			continue
		}
		if g.CritterAt(x, y) != nil {
			continue
		}
		return g.SpawnCritter(proto, x, y)
	}
	return nil
}

func (g *Game) SpawnCritter(proto, x, y int) *Critter {
	if proto < 0 || proto >= len(CritterDefs) {
		return nil
	}
	d := &CritterDefs[proto]
	c := &Critter{
		Proto:    proto,
		Name:     d.Name,
		X:        x,
		Y:        y,
		STR:      d.STR,
		TGH:      d.TGH,
		AGI:      d.AGI,
		INT:      d.INT,
		SPD:      d.SPD,
		Align:    d.Align,
		Class:    d.Class,
		AI:       d.AI,
		Color:    d.Color,
		Glyph:    d.Glyph,
		Skills:   d.Skills,
		Level:    d.Level,
		XP:       d.XP,
		IsAnimal: d.Class == 0,
	}
	c.MaxHP = c.TGH * 3
	if c.MaxHP < 1 {
		c.MaxHP = 1
	}
	c.HP = c.MaxHP
	// Scale a little with dungeon level.
	scale := 1 + g.Level/14
	if scale > 2 && c.Align != AlignNeutral {
		c.MaxHP += c.MaxHP * (scale - 1) / 2
		c.HP = c.MaxHP
		for i := range c.Skills {
			c.Skills[i] += c.Skills[i] * (scale - 1) / 4
			if c.Skills[i] > 95 {
				c.Skills[i] = 95
			}
		}
		c.XP += c.XP * (scale - 1) / 3
	}
	// Give some monsters melee equipment.
	if c.Class == 1 && g.Roll(3) == 0 {
		if g.Roll(2) == 0 {
			c.Equip[EquipMelee] = &Item{Proto: 8 + g.Roll(4), X: -1, Y: -1}
		} else {
			c.Equip[EquipArmor] = &Item{Proto: 5 + g.Roll(3), X: -1, Y: -1}
		}
	}
	if c.Align == AlignVillager {
		if c.Name == "militia" || g.Roll(100) < 40 {
			proto := itemProtoByName("club")
			if g.Roll(2) == 0 {
				proto = itemProtoByName("spear")
			}
			if proto >= 0 {
				c.Equip[EquipMelee] = &Item{Proto: proto, X: -1, Y: -1}
			}
		}
		if g.Roll(100) < 30 {
			c.Equip[EquipArmor] = &Item{Proto: 5, X: -1, Y: -1}
		}
	}
	g.Critters = append(g.Critters, c)
	return c
}

func (g *Game) SetupNewGame(class int, name string) {
	if class < 0 || class >= len(PlayerClasses) {
		class = 0
	}
	g.PlayerClass = class
	g.PlayerName = name
	g.WoodLoot = newWoodLootAssignment(g.RNG)
	g.Player = &Critter{
		Proto:  -1,
		Name:   name,
		Player: true,
		Align:  AlignPlayer,
		Class:  4,
		AI:     AIPlayer,
		Color:  14,
		Glyph:  '@',
		Dead:   false,
	}
	c := &PlayerClasses[class]
	g.Player.STR = c.STR
	g.Player.TGH = c.TGH
	g.Player.AGI = c.AGI
	g.Player.INT = c.INT
	g.Player.SPD = 2
	g.Player.Skills = c.Skills
	g.Player.MaxHP = g.Player.TGH * 3
	g.Player.HP = g.Player.MaxHP
	g.Player.Level = 1
	g.Player.XP = 0
	g.Player.Gold = 0
	g.Player.Blood = 0
	g.Player.Poison = 0
	g.Player.Slots = PlayerSlots{BeltUnlocked: 2}
	kit := ClassStartingKit(class)
	g.Player.Slots.Hands = startingItem(kit.Hands)
	if g.Player.Slots.Hands != nil {
		g.Player.Slots.Hands.Active = true
	}
	g.Player.Slots.Body = startingItem(kit.Body)
	ensureArmorAP(g.Player.Slots.Body)
	for i, proto := range kit.Belt {
		if i >= 2 {
			break
		}
		g.Player.Slots.Belt[i] = startingItem(proto)
	}
	g.Level = 0
	g.Turn = 0
	g.Messages = nil
	g.HintsShown = 0
	g.GenerateLevel(0)
}

func (g *Game) FOVRadius() int {
	if g.Player == nil {
		return 7
	}
	radius := 7 + g.Player.Skill(SkillPerception)/20
	if g.Options.Bool("optPlayerGhostMode") {
		radius += 2
	}
	return radius
}

func (g *Game) UpdateFOV() {
	if g.Map == nil || g.Player == nil {
		return
	}
	g.Map.FOV(g.Player.X, g.Player.Y, g.FOVRadius())
}

func (g *Game) TileDescription(x, y int) string {
	if !g.Map.InBounds(x, y) {
		return "there is unexplored space here"
	}
	if !g.Map.Seen[g.Map.idx(x, y)] {
		return "there is unexplored space here"
	}
	d := g.Map.Def(x, y)
	desc := "there is " + d.Name + " here"
	if d.Type == 3 {
		desc += " (door)"
	}
	if items := g.ItemsAt(x, y); len(items) > 0 {
		desc += "; " + capitalizeFirst(items[0].Name()) + " is lying here."
	}
	if c := g.CritterAt(x, y); c != nil && g.Map.Visible[g.Map.idx(x, y)] {
		desc += "; the " + c.NameOrUnknown() + " is standing here"
	}
	return desc
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
