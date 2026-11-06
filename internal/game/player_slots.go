package game

import (
	"fmt"
	"strings"
)

// Player slot helpers. Monsters keep using the legacy Inventory/Equip fields.

func (c *Critter) SlotItem(slot int) *Item {
	if c == nil {
		return nil
	}
	switch slot {
	case SlotHands:
		return c.Slots.Hands
	case SlotBody:
		return c.Slots.Body
	case SlotBelt1:
		if c.Slots.BeltUnlocked >= 1 {
			return c.Slots.Belt[0]
		}
	case SlotBelt2:
		if c.Slots.BeltUnlocked >= 2 {
			return c.Slots.Belt[1]
		}
	case SlotBelt3:
		if c.Slots.BeltUnlocked >= 3 {
			return c.Slots.Belt[2]
		}
	}
	return nil
}

func (c *Critter) SetSlotItem(slot int, it *Item) {
	if c == nil {
		return
	}
	switch slot {
	case SlotHands:
		c.Slots.Hands = it
	case SlotBody:
		c.Slots.Body = it
	case SlotBelt1:
		if c.Slots.BeltUnlocked >= 1 {
			c.Slots.Belt[0] = it
		}
	case SlotBelt2:
		if c.Slots.BeltUnlocked >= 2 {
			c.Slots.Belt[1] = it
		}
	case SlotBelt3:
		if c.Slots.BeltUnlocked >= 3 {
			c.Slots.Belt[2] = it
		}
	}
}

func (c *Critter) BeltSlotFree(i int) bool {
	if i < 0 || i > 2 || c.Slots.BeltUnlocked < i+1 {
		return false
	}
	return c.Slots.Belt[i] == nil
}

func (c *Critter) FirstFreeBelt() int {
	for i := 0; i < c.Slots.BeltUnlocked && i < 3; i++ {
		if c.Slots.Belt[i] == nil {
			return i
		}
	}
	return -1
}

func (c *Critter) HandItem() *Item {
	if c == nil {
		return nil
	}
	if c.Player {
		return c.Slots.Hands
	}
	return c.Equip[EquipMelee]
}

func (c *Critter) BodyItem() *Item {
	if c == nil {
		return nil
	}
	if c.Player {
		return c.Slots.Body
	}
	return c.Equip[EquipArmor]
}

func (c *Critter) ThrowableItem() *Item {
	if c == nil {
		return nil
	}
	if c.Player {
		return c.Slots.Hands
	}
	return c.Equip[EquipThrowable]
}

func ItemIsHandItem(d *ItemDef) bool {
	if d == nil {
		return false
	}
	switch d.Kind {
	case KindMelee, KindGrenade, KindTool:
		return true
	}
	return false
}

// AutoPlaceItem places a newly picked-up item in the first valid player slot.
func (c *Critter) AutoPlaceItem(it *Item) bool {
	d := it.Def()
	if d == nil {
		return false
	}
	if d.Kind == KindArmor {
		if c.Slots.Body == nil {
			ensureArmorAP(it)
			c.Slots.Body = it
			return true
		}
		return false
	}
	if ItemIsHandItem(d) && c.Slots.Hands == nil {
		c.Slots.Hands = it
		return true
	}
	if i := c.FirstFreeBelt(); i >= 0 {
		c.Slots.Belt[i] = it
		return true
	}
	return false
}

// PlaceGroundItem enforces the one-item-per-tile rule.
func (g *Game) PlaceGroundItem(it *Item, x, y int) bool {
	if it == nil || !g.Map.InBounds(x, y) {
		return false
	}
	ensureArmorAP(it)
	if g.ItemAt(x, y) != nil {
		return false
	}
	it.X, it.Y = x, y
	g.Items = append(g.Items, it)
	return true
}

func (g *Game) PlaceGroundItemNear(it *Item, x, y int) bool {
	for radius := 0; radius <= 6; radius++ {
		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				if abs(dx) != radius && abs(dy) != radius {
					continue
				}
				nx, ny := x+dx, y+dy
				if !g.Map.InBounds(nx, ny) || g.Map.BlocksMove(nx, ny) {
					continue
				}
				if g.PlaceGroundItem(it, nx, ny) {
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) playerConsumeMedical(it *Item) bool {
	if it == nil {
		return false
	}
	d := it.Def()
	if d == nil || d.Kind != KindMedical {
		return false
	}
	if it.Proto == ArmourVialProto {
		arm := g.Player.Slots.Body
		stats, ok := ArmorStatsFor(arm.Def())
		if !ok {
			g.Msg("you need to wear armour to use an armour vial.", 8)
			return false
		}
		ensureArmorAP(arm)
		if arm.AP >= stats.MaxAP {
			g.Msg("your armour is already fully repaired.", 8)
			return false
		}
		if it.Charges <= 0 {
			it.Charges = 1
		}
		arm.AP += 10
		if arm.AP > stats.MaxAP {
			arm.AP = stats.MaxAP
		}
		it.Charges--
		g.Msg(fmt.Sprintf("you apply the armour vial (%d/%d AP).", arm.AP, stats.MaxAP), 10)
		return true
	}
	if it.Charges <= 0 {
		it.Charges = d.Charges
	}
	heal := d.Heal
	if heal <= 0 {
		heal = 10
	}
	g.Player.HP += heal
	if g.Player.HP > g.Player.MaxHP {
		g.Player.HP = g.Player.MaxHP
	}
	it.Charges--
	g.Msg(fmt.Sprintf("you use %s and heal %d health.", d.Name, heal), 10)
	return true
}

func (g *Game) playerConsumeFood(it *Item) bool {
	if it == nil {
		return false
	}
	d := it.Def()
	if d == nil || d.Kind != KindFood {
		return false
	}
	if it.Charges <= 0 {
		it.Charges = d.Charges
		if it.Charges <= 0 {
			it.Charges = 1
		}
	}
	if d.Heal < 0 {
		dmg := -d.Heal
		g.Msg("you eat the "+d.Name+" and feel sick.", 12)
		g.DamageCritter(g.Player, dmg, nil, 2)
	} else {
		heal := d.Heal
		if heal <= 0 {
			heal = 1
		}
		g.Player.HP += heal
		if g.Player.HP > g.Player.MaxHP {
			g.Player.HP = g.Player.MaxHP
		}
		g.Msg(fmt.Sprintf("you eat %s and recover %d health.", d.Name, heal), 10)
	}
	it.Charges--
	return true
}

func (g *Game) playerUseTool(it *Item) bool {
	if it == nil {
		return false
	}
	d := it.Def()
	if d == nil || d.Kind != KindTool {
		return false
	}
	name := strings.ToLower(d.Name)
	switch {
	case strings.Contains(name, "torch"), strings.Contains(name, "lantern"):
		g.Msg("you light the "+d.Name+".", 10)
	case strings.Contains(name, "whetstone"):
		g.Msg("you run the whetstone along your weapon.", 8)
	case strings.Contains(name, "rope"):
		g.Msg("you coil the rope. climbing is not implemented yet.", 8)
	case strings.Contains(name, "shovel"):
		g.Msg("you test the shovel. digging is not implemented yet.", 8)
	case strings.Contains(name, "lockpick"):
		g.Msg("you test the lockpick. locks are not implemented yet.", 8)
	case strings.Contains(name, "fishing rod"):
		g.Msg("you ready the fishing rod. fishing is not implemented yet.", 8)
	default:
		g.Msg("you use the "+d.Name+".", 8)
	}
	return true
}

// UseSlot performs the context-sensitive action for a player slot.
func (g *Game) UseSlot(slot int) {
	it := g.Player.SlotItem(slot)
	if it == nil {
		g.Msg("that slot is empty.", 8)
		return
	}
	d := it.Def()
	if d == nil {
		return
	}
	switch slot {
	case SlotHands:
		switch d.Kind {
		case KindMelee:
			g.Msg("your "+d.Name+" is ready; use f to attack.", 8)
		case KindGrenade:
			g.Msg("throwing is not available yet.", 8)
		case KindMedical:
			if g.playerConsumeMedical(it) {
				if it.Charges <= 0 {
					g.Player.Slots.Hands = nil
				}
				g.EndPlayerTurn()
			}
		case KindFood:
			if g.playerConsumeFood(it) {
				if it.Charges <= 0 {
					g.Player.Slots.Hands = nil
				}
				g.EndPlayerTurn()
			}
		case KindTool:
			if g.playerUseTool(it) {
				g.EndPlayerTurn()
			}
		case KindRelic:
			g.Msg("you turn the "+d.Name+" over. it is a relic of the old world.", 8)
		default:
			g.Msg("you cannot use "+d.Name+" from your hands.", 8)
		}
	case SlotBody:
		g.Msg("your "+d.Name+" is worn.", 8)
	case SlotBelt1, SlotBelt2, SlotBelt3:
		switch d.Kind {
		case KindMelee, KindGrenade, KindTool:
			old := g.Player.Slots.Hands
			g.Player.Slots.Hands = it
			g.Player.SetSlotItem(slot, old)
			g.Msg("you swap your "+d.Name+" into your hands.", 8)
			g.EndPlayerTurn()
		case KindMedical:
			if g.playerConsumeMedical(it) {
				if it.Charges <= 0 {
					g.Player.SetSlotItem(slot, nil)
				}
				g.EndPlayerTurn()
			}
		case KindFood:
			if g.playerConsumeFood(it) {
				if it.Charges <= 0 {
					g.Player.SetSlotItem(slot, nil)
				}
				g.EndPlayerTurn()
			}
		case KindArmor:
			if g.Player.Slots.Body == nil {
				ensureArmorAP(it)
				g.Player.Slots.Body = it
				g.Player.SetSlotItem(slot, nil)
				g.Msg("you equip "+d.Name+".", 10)
				g.EndPlayerTurn()
			} else {
				g.Msg("your body slot is already occupied.", 8)
			}
		default:
			g.Msg("you cannot use "+d.Name+" from the belt.", 8)
		}
	}
}

// DropSlot drops one player slot item onto the current tile.
func (g *Game) DropSlot(slot int) {
	it := g.Player.SlotItem(slot)
	if it == nil {
		g.Msg("that slot is empty.", 8)
		return
	}
	if g.ItemAt(g.Player.X, g.Player.Y) != nil {
		g.Msg("there is already an item here.", 8)
		return
	}
	g.Player.SetSlotItem(slot, nil)
	if !g.PlaceGroundItem(it, g.Player.X, g.Player.Y) {
		g.Player.SetSlotItem(slot, it)
		g.Msg("there is already an item here.", 8)
		return
	}
	g.Msg("you drop "+it.Name()+".", 8)
	g.EndPlayerTurn()
}
