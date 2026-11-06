package game

import "math"

// WarmLight returns the strongest warm-light contribution at a tile from
// lantern posts, campfires, fireplaces and active fire effects.
func (m *Map) WarmLight(x, y int) float64 {
	if m == nil {
		return 0
	}
	best := 0.0
	for yy := y - 5; yy <= y+5; yy++ {
		for xx := x - 5; xx <= x+5; xx++ {
			if !m.InBounds(xx, yy) {
				continue
			}
			radius := 0
			switch m.Tile(xx, yy) {
			case tileLanternPost:
				radius = 4
			case tileCampfire, tileFireplace:
				radius = 5
			}
			if m.Fire[m.idx(xx, yy)] > 0 {
				radius = 5
			}
			if radius == 0 {
				continue
			}
			d := math.Hypot(float64(x-xx), float64(y-yy))
			if d > float64(radius) {
				continue
			}
			v := 1 - d/float64(radius)
			if v > best {
				best = v
			}
		}
	}
	return best
}
