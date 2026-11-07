package game

// This file defines the reduced melee-focused scope: which skills, options,
// items and throwables remain active after the command-scope cut.

// ActiveSkills is the player-facing skill list.  The underlying [NumSkills]
// arrays are kept for save and data compatibility, but removed skills are no
// longer shown or trained.
var ActiveSkills = []int{
	SkillPerception,
	SkillDodging,
	SkillMelee,
	SkillThrowing,
	SkillMetabolism,
	SkillFocusing,
}

func SkillActive(skill int) bool {
	for _, s := range ActiveSkills {
		if s == skill {
			return true
		}
	}
	return false
}

var inactiveOptions = map[string]bool{
	"optMechanicBaseTileDV":            true, // tile-targeted ranged attacks
	"optMechanicLineOfFireHitChance":   true, // line-of-fire shooting
	"optPlayerUseTargetLine":           true, // ranged target line
	"optPlayerClearTargetResetTactics": true, // target command removed
	"optPlayerAttackOnMove":            true, // bump attacks are always enabled now
}

// ActiveOptionDefs returns the options that are still relevant to the
// melee-focused game, in display order.
func ActiveOptionDefs() []OptionDef {
	out := make([]OptionDef, 0, len(OptionDefs))
	for _, d := range OptionDefs {
		if inactiveOptions[d.Name] {
			continue
		}
		out = append(out, d)
	}
	return out
}

// ItemAvailable reports whether an item prototype may appear in normal play.
// Removed ranged/ammo items, guns, old body armour and other out-of-theme
// prototypes remain in the raw data for reference but are not spawned or used.
func ItemAvailable(proto int) bool {
	if proto < 0 || proto >= len(ItemDefs) {
		return false
	}
	d := &ItemDefs[proto]
	switch d.Kind {
	case KindRanged, KindAmmo, KindArmor:
		return false
	case KindGrenade:
		switch d.Name {
		case "holy water flask", "holy fire flask", "alchemical acid flask", "incense flask", "radiant flask":
			return true
		}
		return false
	}
	return true
}

func ItemIsThrowable(proto int) bool {
	if proto < 0 || proto >= len(ItemDefs) {
		return false
	}
	d := &ItemDefs[proto]
	return d.Kind == KindGrenade && ItemAvailable(proto)
}

func filterAvailableItems(ids []int) []int {
	out := ids[:0]
	for _, id := range ids {
		if ItemAvailable(id) {
			out = append(out, id)
		}
	}
	return out
}

// ArmourVialProto is appended at init so it does not disturb recovered IDs.
var ArmourVialProto int

func init() {
	ArmourVialProto = len(ItemDefs)
	ItemDefs = append(ItemDefs, ItemDef{
		Name:    "armour vial",
		Desc:    "a small vial of hardening compound that restores armour points.",
		Kind:    KindMedical,
		Charges: 1,
		Glyph:   '!',
		Color:   9,
	})
}

// ActiveKeyBindings is the reduced, player-facing help table.  The recovered
// KeyBindings data remains available as historical reference.
var ActiveKeyBindings = []KeyBind{
	{Key: 'h', Name: "move to west", Desc: "move player west"},
	{Key: 'j', Name: "move to south", Desc: "move player south"},
	{Key: 'k', Name: "move to north", Desc: "move player north"},
	{Key: 'l', Name: "move to east", Desc: "move player east"},
	{Key: '.', Name: "wait", Desc: "wait one turn"},
	{Key: 'x', Name: "bash", Desc: "choose a direction and bash an adjacent small obstacle"},
	{Key: 'a', Name: "hands action", Desc: "use the Hands item (tool) or report the weapon ready"},
	{Key: 'b', Name: "head action", Desc: "wear or inspect the Head item"},
	{Key: 'c', Name: "robe slot 1", Desc: "use or consume robe slot 1"},
	{Key: 'd', Name: "robe slot 2", Desc: "use or consume robe slot 2"},
	{Key: 'e', Name: "robe slot 3", Desc: "use or consume robe slot 3"},
	{Key: 'f', Name: "robe slot 4", Desc: "use or consume robe slot 4"},
	{Key: 'g', Name: "robe slot 5", Desc: "use or consume robe slot 5"},
	{Key: 'r', Name: "use rite", Desc: "use the first known rite"},
	{Key: 't', Name: "use second rite", Desc: "use the second known rite"},
	{Key: 'p', Name: "pick up item", Desc: "pick up items from the current tile"},
	{Key: 'i', Name: "inventory", Desc: "display equipment; uppercase A-G drop slots"},
	{Key: 'v', Name: "view / look", Desc: "look around and examine the surroundings"},
	{Key: '~', Name: "message log", Desc: "expand or collapse the latest eight messages"},
	{Key: '@', Name: "player screen", Desc: "display character information"},
	{Key: 'o', Name: "options", Desc: "display the options screen"},
	{Key: '?', Name: "help", Desc: "display this key binding screen"},
	{Key: ' ', Name: "close overlay", Desc: "close a modal screen without taking an action"},
	{Key: 's', Name: "save and quit", Desc: "save the game and quit"},
	{Key: '>', Name: "use portal", Desc: "use a portal, stairs or elevator"},
	{Key: ']', Name: "next map", Desc: "jump directly to the next map"},
	{Key: '[', Name: "previous map", Desc: "jump directly to the previous map"},
	{Key: 'q', Name: "quit", Desc: "quit the game"},
}

// ActiveHints is the melee-focused hint list.  The recovered hints that refer
// to ranged weapons, ammo or searching are no longer relevant.
var ActiveHints = []string{
	Hints[0], // medkit
	Hints[1], // fight one enemy at a time
	Hints[4], // use the surroundings
	Hints[6], // inventory takes time
	Hints[7], // acid ignores armor
	"you are defenceless without a weapon; keep a melee weapon equipped.",
	"holy flasks can only be thrown at tiles within your sight range.",
	"bash small obstacles with x and a direction; barrels and crates may hide supplies.",
}
