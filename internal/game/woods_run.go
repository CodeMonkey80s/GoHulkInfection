package game

import "math/rand"

// Map edges used by the woods snake route.
const (
	EdgeNorth = iota
	EdgeSouth
	EdgeEast
	EdgeWest
)

// Woods map themes. These are run-local identities, not level numbers.
const (
	WoodsChurch = iota
	WoodsPeasantHouses
	WoodsMillPond
	WoodsCemetery
	WoodsMeteorField
	numWoodsMaps
)

// WoodsRun stores the once-per-run woods layout: which middle themes are used,
// where each map is entered and where its exit lies.
type WoodsRun struct {
	Themes  [numWoodsMaps]int `json:"themes"`
	Entries [numWoodsMaps]int `json:"entries"`
	Exits   [numWoodsMaps]int `json:"exits"`
}

var woodThemeNames = [numWoodsMaps]string{
	"The Church",
	"Peasant Houses",
	"The Mill and Pond",
	"The Forgotten Cemetery",
	"The Meteor Field",
}

func woodsThemeName(theme int) string {
	if theme < 0 || theme >= len(woodThemeNames) {
		return woodThemeNames[0]
	}
	return woodThemeNames[theme]
}

func oppositeEdge(edge int) int {
	switch edge {
	case EdgeNorth:
		return EdgeSouth
	case EdgeSouth:
		return EdgeNorth
	case EdgeEast:
		return EdgeWest
	default:
		return EdgeEast
	}
}

func edgeAllowed(edge int, excluded ...int) bool {
	for _, e := range excluded {
		if edge == e {
			return false
		}
	}
	return true
}

func pickEdge(rng *rand.Rand, excluded ...int) int {
	options := make([]int, 0, 4)
	for _, edge := range []int{EdgeNorth, EdgeSouth, EdgeEast, EdgeWest} {
		if edgeAllowed(edge, excluded...) {
			options = append(options, edge)
		}
	}
	if len(options) == 0 {
		return EdgeNorth
	}
	return options[rng.Intn(len(options))]
}

// newWoodsRun shuffles the three middle themes once and builds a snake route
// whose next entry is always opposite the previous exit. Map 3 is forced to
// exit east so the fixed Meteor Field can enter from the west without a
// reversal.
func newWoodsRun(rng *rand.Rand) WoodsRun {
	run := WoodsRun{
		Themes: [numWoodsMaps]int{WoodsChurch, 0, 0, 0, WoodsMeteorField},
	}
	middle := []int{WoodsPeasantHouses, WoodsMillPond, WoodsCemetery}
	rng.Shuffle(len(middle), func(i, j int) {
		middle[i], middle[j] = middle[j], middle[i]
	})
	copy(run.Themes[1:4], middle)

	run.Entries[0] = EdgeWest
	run.Exits[0] = EdgeEast
	for i := 1; i <= 2; i++ {
		run.Entries[i] = oppositeEdge(run.Exits[i-1])
		excluded := []int{run.Entries[i]}
		if i == 2 {
			// Leave the east edge free for map 3's fixed exit.
			excluded = append(excluded, EdgeWest)
		}
		run.Exits[i] = pickEdge(rng, excluded...)
	}
	run.Entries[3] = oppositeEdge(run.Exits[2])
	run.Exits[3] = EdgeEast
	run.Entries[4] = EdgeWest
	run.Exits[4] = EdgeEast
	return run
}

func validWoodsRun(run *WoodsRun) bool {
	if run == nil {
		return false
	}
	if run.Themes[0] != WoodsChurch || run.Themes[4] != WoodsMeteorField {
		return false
	}
	seen := map[int]bool{WoodsChurch: true, WoodsMeteorField: true}
	for _, theme := range run.Themes[1:4] {
		if theme < WoodsPeasantHouses || theme > WoodsCemetery || seen[theme] {
			return false
		}
		seen[theme] = true
	}
	for i := 0; i < numWoodsMaps; i++ {
		if run.Entries[i] < EdgeNorth || run.Entries[i] > EdgeWest {
			return false
		}
		if run.Exits[i] < EdgeNorth || run.Exits[i] > EdgeWest {
			return false
		}
		if run.Exits[i] == run.Entries[i] {
			return false
		}
		if i > 0 && run.Entries[i] != oppositeEdge(run.Exits[i-1]) {
			return false
		}
	}
	return run.Entries[0] == EdgeWest && run.Exits[0] == EdgeEast &&
		run.Entries[4] == EdgeWest && run.Exits[4] == EdgeEast
}

func woodsRunIndex(level int) int {
	if level < 0 {
		return 0
	}
	if level >= numWoodsMaps {
		return numWoodsMaps - 1
	}
	return level
}

func (g *Game) currentWoodsTheme() int {
	if g == nil {
		return WoodsChurch
	}
	if validWoodsRun(&g.WoodsRun) {
		return g.WoodsRun.Themes[woodsRunIndex(g.Level)]
	}
	return woodsRunIndex(g.Level)
}
