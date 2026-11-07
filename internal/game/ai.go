package game

import (
	"container/heap"
	"fmt"
)

type point struct{ x, y int }

type pathNode struct {
	x, y int
	f    int
}

type pathHeap []pathNode

func (h pathHeap) Len() int            { return len(h) }
func (h pathHeap) Less(i, j int) bool  { return h[i].f < h[j].f }
func (h pathHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *pathHeap) Push(x interface{}) { *h = append(*h, x.(pathNode)) }
func (h *pathHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

var dirs4 = [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func (g *Game) UpdateMonsters() {
	if !g.Options.Bool("optGameUseAI") {
		return
	}
	for _, c := range append([]*Critter(nil), g.Critters...) {
		if !c.IsAlive() {
			continue
		}
		c.Energy += 50 + c.SPD*15
		if c.Energy < 100 {
			continue
		}
		c.Energy -= 100
		g.MonsterTurn(c)
		if g.State == StateGameOver || !g.Player.IsAlive() {
			return
		}
	}
}

func (g *Game) MonsterTurn(c *Critter) {
	if c.X == c.LastX && c.Y == c.LastY {
		c.Stuck++
	} else {
		c.Stuck = 0
	}
	c.LastX, c.LastY = c.X, c.Y
	if c.Stuck >= 3 {
		c.Path = nil
		c.Stuck = 0
	}
	switch c.AI {
	case AIObserver:
		p := g.Player
		if p == nil || !p.IsAlive() {
			return
		}
		dist, canSee := g.canSeeCritter(c, p, 14)
		if canSee {
			c.LastSeenX, c.LastSeenY = p.X, p.Y
			c.Alert = 10
		}
		g.observerTurn(c, canSee, dist)
	case AIVillager:
		g.villagerTurn(c)
	case AIAnimal:
		g.animalTurn(c)
	default:
		g.zombieTurn(c)
	}
}

func (g *Game) zombieTurn(c *Critter) {
	radius := 3
	if c.Align == AlignEnemy && c.IsAnimal {
		radius = 8
	}
	target, dist := g.nearestHostile(c, radius)
	if target != nil {
		c.LastSeenX, c.LastSeenY = target.X, target.Y
		c.Alert = 8
		if dist <= 1 {
			g.MonsterAttack(c, target)
			return
		}
		g.stepToward(c, target.X, target.Y)
		return
	}
	if c.Alert > 0 {
		c.Alert--
		g.stepToward(c, c.LastSeenX, c.LastSeenY)
		return
	}
	if c.Class == 1 && g.Roll(100) < 70 {
		return
	}
	g.wander(c)
}

func (g *Game) animalTurn(c *Critter) {
	if c.Align == AlignEnemy {
		g.zombieTurn(c)
		return
	}
	threat, _ := g.nearestAnimalThreat(c, 4)
	if threat != nil {
		g.stepAway(c, threat.X, threat.Y)
		return
	}
	if g.Roll(100) < 60 {
		return
	}
	g.wander(c)
}

func (g *Game) villagerTurn(c *Critter) {
	enemy, dist := g.nearestHostile(c, 6)
	if enemy == nil {
		if g.Roll(100) < 50 {
			return
		}
		g.wander(c)
		return
	}
	armed := c.HandItem() != nil
	if dist <= 1 {
		if armed {
			g.MonsterAttack(c, enemy)
			return
		}
		if g.stepAway(c, enemy.X, enemy.Y) {
			return
		}
		g.MonsterAttack(c, enemy)
		return
	}
	if armed && dist <= 3 {
		g.stepToward(c, enemy.X, enemy.Y)
		return
	}
	if dist <= 5 {
		g.stepAway(c, enemy.X, enemy.Y)
		return
	}
	if g.Roll(100) < 50 {
		return
	}
	g.wander(c)
}

func (g *Game) areHostile(a, b *Critter) bool {
	if a == nil || b == nil || a == b {
		return false
	}
	switch a.Align {
	case AlignEnemy:
		return b.Align == AlignPlayer || b.Align == AlignVillager || b.Align == AlignNeutral
	case AlignVillager:
		return b.Align == AlignEnemy
	case AlignPlayer:
		return b.Align == AlignEnemy
	}
	return false
}

// canSeeCritter applies the same sight blockers as player FOV: dense trees,
// brambles, bushes, walls and closed doors stop AI perception.
func (g *Game) canSeeCritter(c, other *Critter, maxDist int) (int, bool) {
	if g.Map == nil || c == nil || other == nil || !other.IsAlive() {
		return 0, false
	}
	dist := chebyshev(c.X, c.Y, other.X, other.Y)
	if dist > maxDist {
		return 0, false
	}
	if !g.Map.LineOfSight(c.X, c.Y, other.X, other.Y) {
		return 0, false
	}
	return dist, true
}

func (g *Game) nearestHostile(c *Critter, maxDist int) (*Critter, int) {
	best, bestDist := (*Critter)(nil), maxDist+1
	consider := func(other *Critter) {
		if other == nil || other == c || !other.IsAlive() || !g.areHostile(c, other) {
			return
		}
		d, canSee := g.canSeeCritter(c, other, maxDist)
		if !canSee || d >= bestDist {
			return
		}
		best, bestDist = other, d
	}
	if g.Player != nil && g.Player.IsAlive() {
		consider(g.Player)
	}
	for _, other := range g.Critters {
		consider(other)
	}
	return best, bestDist
}

func (g *Game) nearestAnimalThreat(c *Critter, maxDist int) (*Critter, int) {
	best, bestDist := (*Critter)(nil), maxDist+1
	consider := func(other *Critter) {
		if other == nil || other == c || !other.IsAlive() {
			return
		}
		if other.Align != AlignPlayer && other.Align != AlignEnemy {
			return
		}
		d, canSee := g.canSeeCritter(c, other, maxDist)
		if !canSee || d >= bestDist {
			return
		}
		best, bestDist = other, d
	}
	if g.Player != nil && g.Player.IsAlive() {
		consider(g.Player)
	}
	for _, other := range g.Critters {
		consider(other)
	}
	return best, bestDist
}

func (g *Game) observerTurn(c *Critter, canSee bool, dist int) {
	if !canSee {
		if c.Alert > 0 {
			c.Alert--
			g.stepToward(c, c.LastSeenX, c.LastSeenY)
		} else {
			g.wander(c)
		}
		return
	}
	if dist <= 1 {
		g.MonsterAttack(c, g.Player)
		return
	}
	g.stepToward(c, g.Player.X, g.Player.Y)
}

func (g *Game) stepToward(c *Critter, tx, ty int) {
	if c.X == tx && c.Y == ty {
		c.Path = nil
		return
	}
	if c.Stuck >= 2 {
		c.Path = nil
		if g.trySidestep(c, tx, ty) {
			c.Stuck = 0
			return
		}
	}
	path := g.pathFor(c, tx, ty)
	if len(path) == 0 {
		g.wander(c)
		return
	}
	next := path[0]
	if other := g.CritterAt(next.x, next.y); other != nil && other != c {
		if g.areHostile(c, other) {
			g.MonsterAttack(c, other)
			return
		}
		c.Path = nil
		if !g.trySidestep(c, tx, ty) {
			g.wander(c)
		}
		return
	}
	if g.Map.BlocksMove(next.x, next.y) {
		d := g.Map.Def(next.x, next.y)
		if d.Type == 3 {
			g.OpenDoor(next.x, next.y)
			return
		}
		c.Path = nil
		g.wander(c)
		return
	}
	c.X, c.Y = next.x, next.y
	c.Path = path[1:]
	g.applyMonsterTile(c)
}

func (g *Game) trySidestep(c *Critter, tx, ty int) bool {
	bestX, bestY := c.X, c.Y
	bestD := manhattan(c.X, c.Y, tx, ty)
	fallbackX, fallbackY := c.X, c.Y
	for _, d := range dirs4 {
		nx, ny := c.X+d[0], c.Y+d[1]
		if !g.Map.IsPassable(nx, ny) || g.CritterAt(nx, ny) != nil {
			continue
		}
		idx := g.Map.idx(nx, ny)
		if g.Map.Fire[idx] > 0 || g.Map.Acid[idx] > 0 {
			continue
		}
		if fallbackX == c.X && fallbackY == c.Y {
			fallbackX, fallbackY = nx, ny
		}
		nd := manhattan(nx, ny, tx, ty)
		if nd < bestD {
			bestD = nd
			bestX, bestY = nx, ny
		}
	}
	if bestX == c.X && bestY == c.Y {
		bestX, bestY = fallbackX, fallbackY
	}
	if bestX == c.X && bestY == c.Y {
		return false
	}
	c.X, c.Y = bestX, bestY
	g.applyMonsterTile(c)
	return true
}

func (g *Game) stepAway(c *Critter, fx, fy int) bool {
	bestX, bestY := c.X, c.Y
	bestD := manhattan(c.X, c.Y, fx, fy)
	for _, d := range dirs4 {
		nx, ny := c.X+d[0], c.Y+d[1]
		if !g.Map.IsPassable(nx, ny) || g.CritterAt(nx, ny) != nil {
			continue
		}
		idx := g.Map.idx(nx, ny)
		if g.Map.Fire[idx] > 0 || g.Map.Acid[idx] > 0 {
			continue
		}
		nd := manhattan(nx, ny, fx, fy)
		if nd > bestD {
			bestD = nd
			bestX, bestY = nx, ny
		}
	}
	moved := bestX != c.X || bestY != c.Y
	c.X, c.Y = bestX, bestY
	g.applyMonsterTile(c)
	return moved
}

func (g *Game) wander(c *Critter) {
	for i := 0; i < 4; i++ {
		d := dirs4[g.Roll(len(dirs4))]
		nx, ny := c.X+d[0], c.Y+d[1]
		if !g.Map.IsPassable(nx, ny) || g.CritterAt(nx, ny) != nil {
			continue
		}
		idx := g.Map.idx(nx, ny)
		if g.Map.Fire[idx] > 0 || g.Map.Acid[idx] > 0 {
			continue
		}
		c.X, c.Y = nx, ny
		g.applyMonsterTile(c)
		return
	}
}

func (g *Game) pathFor(c *Critter, tx, ty int) []point {
	if len(c.Path) > 0 && c.PathTargetX == tx && c.PathTargetY == ty && g.Turn <= c.PathExpires {
		return c.Path
	}
	path := g.findPath(c, tx, ty)
	c.Path = path
	c.PathTargetX, c.PathTargetY = tx, ty
	c.PathExpires = g.Turn + 3
	return path
}

func manhattan(x0, y0, x1, y1 int) int {
	return abs(x1-x0) + abs(y1-y0)
}

// findStep returns the first cardinal step of an A* path to (tx,ty).
func (g *Game) findStep(c *Critter, tx, ty int) (int, int, bool) {
	path := g.findPath(c, tx, ty)
	if len(path) == 0 {
		return 0, 0, false
	}
	return path[0].x, path[0].y, true
}

// findPath runs four-way A* and returns the path excluding the start tile.
func (g *Game) findPath(c *Critter, tx, ty int) []point {
	if !g.Map.InBounds(tx, ty) {
		return nil
	}
	const inf = 1 << 30
	const maxExpansions = 4096
	w, h := g.Map.W, g.Map.H
	n := w * h
	start := c.Y*w + c.X
	goal := ty*w + tx
	if start == goal {
		return nil
	}
	gScore := make([]int, n)
	closed := make([]bool, n)
	prev := make([]int, n)
	for i := range gScore {
		gScore[i] = inf
		prev[i] = -1
	}
	gScore[start] = 0
	open := &pathHeap{{x: c.X, y: c.Y, f: manhattan(c.X, c.Y, tx, ty)}}
	heap.Init(open)
	expansions := 0
	for open.Len() > 0 && expansions < maxExpansions {
		cur := heap.Pop(open).(pathNode)
		ci := cur.y*w + cur.x
		if closed[ci] {
			continue
		}
		closed[ci] = true
		if ci == goal {
			break
		}
		expansions++
		for _, d := range dirs4 {
			nx, ny := cur.x+d[0], cur.y+d[1]
			if nx < 0 || ny < 0 || nx >= w || ny >= h {
				continue
			}
			ni := ny*w + nx
			if closed[ni] {
				continue
			}
			step, ok := g.pathStepCost(c, nx, ny, tx, ty)
			if !ok {
				continue
			}
			tentative := gScore[ci] + step
			if tentative < gScore[ni] {
				prev[ni] = ci
				gScore[ni] = tentative
				heap.Push(open, pathNode{x: nx, y: ny, f: tentative + manhattan(nx, ny, tx, ty)})
			}
		}
	}
	if prev[goal] == -1 {
		return nil
	}
	var rev []point
	for at := goal; at != start; at = prev[at] {
		if at < 0 {
			return nil
		}
		rev = append(rev, point{x: at % w, y: at / w})
	}
	path := make([]point, len(rev))
	for i := range rev {
		path[i] = rev[len(rev)-1-i]
	}
	return path
}

func (g *Game) pathStepCost(c *Critter, nx, ny, tx, ty int) (int, bool) {
	if !g.Map.InBounds(nx, ny) {
		return 0, false
	}
	if nx == tx && ny == ty {
		return 1, true
	}
	t := g.Map.Tile(nx, ny)
	d := g.Map.Def(nx, ny)
	if d.BlockMove {
		if d.Type == 3 {
			return 3, true
		}
		return 0, false
	}
	if t == tileDeepWater {
		return 0, false
	}
	idx := g.Map.idx(nx, ny)
	if g.Map.Fire[idx] > 0 || g.Map.Acid[idx] > 0 || g.Map.Smoke[idx] > 0 {
		return 20, true
	}
	if t == tileShallowWater || t == tileLilyPad {
		return 8, true
	}
	if t == tileRiverBank {
		return 2, true
	}
	if other := g.CritterAt(nx, ny); other != nil && other != c {
		return 40, true
	}
	switch t {
	case tileBush, tileReeds:
		return 2, true
	default:
		return 1, true
	}
}

func (g *Game) applyMonsterTile(c *Critter) {
	if !g.Map.InBounds(c.X, c.Y) {
		return
	}
	i := g.Map.idx(c.X, c.Y)
	if g.Map.Fire[i] > 0 {
		g.DamageCritter(c, 10, nil, 2)
	}
	if g.Map.Acid[i] > 0 {
		g.DamageCritter(c, 12, nil, 2)
	}
}

func (g *Game) OpenDoor(x, y int) {
	d := g.Map.Def(x, y)
	if d.Type != 3 {
		return
	}
	if d.Partner >= 0 && d.Partner < len(TileDefs) {
		g.Map.SetTile(x, y, uint8(d.Partner))
		g.Msg("the "+d.Name+" opens.", 8)
	} else {
		g.Msg("you cannot open the "+d.Name+".", 8)
	}
}

func (g *Game) CloseDoor(x, y int) {
	d := g.Map.Def(x, y)
	if d.Type != 3 {
		return
	}
	if d.Partner >= 0 && d.Partner < len(TileDefs) {
		g.Map.SetTile(x, y, uint8(d.Partner))
		g.Msg("the "+d.Name+" closes.", 8)
	}
}

func (g *Game) MonsterMessage(c *Critter, text string) {
	g.Msg(fmt.Sprintf("%s %s", c.NameOrUnknown(), text), 8)
}
