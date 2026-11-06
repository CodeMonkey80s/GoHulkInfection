package game

import "fmt"

type point struct{ x, y int }

var dirs8 = [][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}}

func (g *Game) UpdateMonsters() {
	if !g.Options.Bool("optGameUseAI") {
		return
	}
	for _, c := range g.Critters {
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
	p := g.Player
	if p == nil || !p.IsAlive() {
		return
	}
	dist := chebyshev(c.X, c.Y, p.X, p.Y)
	canSee := dist <= 14 && g.Map.LineOfSight(c.X, c.Y, p.X, p.Y)
	if canSee {
		c.LastSeenX, c.LastSeenY = p.X, p.Y
		c.Alert = 10
		if c.Align == AlignNeutral && c.AI == AIAnimal {
			// Animals only fight if wounded or attacked.
			canSee = c.HP < c.MaxHP
		}
	}
	switch c.AI {
	case AIObserver:
		g.observerTurn(c, canSee, dist)
	case AIAnimal:
		g.animalTurn(c, canSee, dist)
	default:
		g.zombieTurn(c, canSee, dist)
	}
}

func (g *Game) zombieTurn(c *Critter, canSee bool, dist int) {
	switch {
	case canSee && dist <= 1:
		g.MonsterAttack(c, g.Player)
	case canSee:
		g.stepToward(c, g.Player.X, g.Player.Y)
	case c.Alert > 0:
		c.Alert--
		g.stepToward(c, c.LastSeenX, c.LastSeenY)
	default:
		g.wander(c)
	}
}

func (g *Game) animalTurn(c *Critter, canSee bool, dist int) {
	if c.Align == AlignEnemy {
		g.zombieTurn(c, canSee, dist)
		return
	}
	// Neutral animals flee when the player is close.
	if canSee && dist <= 3 {
		g.stepAway(c, g.Player.X, g.Player.Y)
		return
	}
	g.wander(c)
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
		return
	}
	nx, ny, ok := g.findStep(c, tx, ty)
	if !ok {
		g.wander(c)
		return
	}
	if other := g.CritterAt(nx, ny); other != nil {
		if other.Player {
			g.MonsterAttack(c, other)
		}
		return
	}
	if g.Map.BlocksMove(nx, ny) {
		// Try to open a door.
		d := g.Map.Def(nx, ny)
		if d.Type == 3 {
			g.OpenDoor(nx, ny)
		}
		return
	}
	c.X, c.Y = nx, ny
	g.applyMonsterTile(c)
}

func (g *Game) stepAway(c *Critter, fx, fy int) {
	bestX, bestY := c.X, c.Y
	bestD := chebyshev(c.X, c.Y, fx, fy)
	for _, d := range dirs8 {
		nx, ny := c.X+d[0], c.Y+d[1]
		if !g.Map.IsPassable(nx, ny) {
			continue
		}
		if g.CritterAt(nx, ny) != nil {
			continue
		}
		nd := chebyshev(nx, ny, fx, fy)
		if nd > bestD {
			bestD = nd
			bestX, bestY = nx, ny
		}
	}
	c.X, c.Y = bestX, bestY
	g.applyMonsterTile(c)
}

func (g *Game) wander(c *Critter) {
	for i := 0; i < 4; i++ {
		d := dirs8[g.Roll(len(dirs8))]
		nx, ny := c.X+d[0], c.Y+d[1]
		if g.Map.IsPassable(nx, ny) && g.CritterAt(nx, ny) == nil {
			c.X, c.Y = nx, ny
			g.applyMonsterTile(c)
			return
		}
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

// findStep does a breadth-first search and returns the first step from c to (tx,ty).
func (g *Game) findStep(c *Critter, tx, ty int) (int, int, bool) {
	if !g.Map.InBounds(tx, ty) {
		return 0, 0, false
	}
	w, h := g.Map.W, g.Map.H
	visited := make([]bool, w*h)
	prev := make([]int, w*h)
	for i := range prev {
		prev[i] = -1
	}
	start := c.Y*w + c.X
	goal := ty*w + tx
	q := []int{start}
	visited[start] = true
	head := 0
	for head < len(q) {
		cur := q[head]
		head++
		if cur == goal {
			break
		}
		cx, cy := cur%w, cur/w
		for _, d := range dirs8 {
			nx, ny := cx+d[0], cy+d[1]
			if nx < 0 || ny < 0 || nx >= w || ny >= h {
				continue
			}
			ni := ny*w + nx
			if visited[ni] {
				continue
			}
			blocked := g.Map.BlocksMove(nx, ny)
			isDoor := g.Map.Def(nx, ny).Type == 3
			if blocked && !isDoor {
				continue
			}
			if other := g.CritterAt(nx, ny); other != nil && !other.Player {
				continue
			}
			visited[ni] = true
			prev[ni] = cur
			q = append(q, ni)
		}
	}
	if !visited[goal] {
		return 0, 0, false
	}
	cur := goal
	for prev[cur] != start {
		if prev[cur] < 0 {
			return 0, 0, false
		}
		cur = prev[cur]
	}
	return cur % w, cur / w, true
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
