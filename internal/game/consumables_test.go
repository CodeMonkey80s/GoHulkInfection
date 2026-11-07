package game

import (
	"strings"
	"testing"
)

func TestHealingConsumablesAreNotUsedAtFullHealth(t *testing.T) {
	for _, name := range []string{"apple", "bandage"} {
		t.Run(name, func(t *testing.T) {
			g := NewGame(16000)
			g.SetupNewGame(0, "tester")
			g.State = StatePlaying
			g.Player.HP = g.Player.MaxHP
			proto := itemProtoByName(name)
			if proto < 0 {
				t.Fatalf("item %q not found", name)
			}
			it := &Item{Proto: proto, X: -1, Y: -1, Charges: 1}
			g.Player.Slots.Hands = it
			turn := g.Turn
			g.UseSlot(SlotHands)
			if g.Player.Slots.Hands != it {
				t.Fatalf("%s was consumed or removed at full health", name)
			}
			if it.Charges != 1 {
				t.Fatalf("%s charges changed at full health: %d", name, it.Charges)
			}
			if g.Turn != turn {
				t.Fatalf("%s consumed a turn at full health: %d -> %d", name, turn, g.Turn)
			}
			msgs := g.LastMessages(1)
			if len(msgs) == 0 || !strings.Contains(msgs[0].Text, "full health") {
				t.Fatalf("%s did not report full health: %v", name, msgs)
			}
		})
	}
}
