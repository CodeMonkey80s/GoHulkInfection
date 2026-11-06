package game

import "fmt"

func (g *Game) OptionValue(name string) int {
	return g.Options.Get(name)
}

func (g *Game) ChangeOption(name string, delta int) {
	for i := range OptionDefs {
		d := &OptionDefs[i]
		if d.Name != name {
			continue
		}
		v := g.Options.Get(name) + delta
		if v < d.Min {
			v = d.Min
		}
		if v > d.Max {
			v = d.Max
		}
		g.Options[name] = v
		return
	}
}

func (g *Game) ToggleOption(name string) {
	for i := range OptionDefs {
		d := &OptionDefs[i]
		if d.Name != name {
			continue
		}
		v := g.Options.Get(name)
		if v == 0 {
			v = 1
		} else {
			v = 0
		}
		if v < d.Min {
			v = d.Min
		}
		if v > d.Max {
			v = d.Max
		}
		g.Options[name] = v
		return
	}
}

func (g *Game) DumpOptions() {
	g.Msg("*** dumping options...", 8)
	for _, d := range OptionDefs {
		g.Msg(fmt.Sprintf("%s = %d", d.Name, g.Options.Get(d.Name)), 8)
	}
	g.Msg("*** end of options", 8)
	g.Msg("options dumped...", 8)
}
