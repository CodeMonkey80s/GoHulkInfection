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
// Removed ranged/ammo items and unused grenade variants remain in the data for
// save compatibility but are not spawned, equipped or used.
func ItemAvailable(proto int) bool {
	if proto < 0 || proto >= len(ItemDefs) {
		return false
	}
	d := &ItemDefs[proto]
	switch d.Kind {
	case KindRanged, KindAmmo:
		return false
	case KindGrenade:
		return d.Name == "he grenade" || d.Name == "molotov cocktail"
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

func init() {
	// Repurpose the recovered incendiary grenade as a molotov cocktail.  Its
	// weapon subtype still gives the fire damage range, while ThrowGrenade
	// applies the burning ground effect.
	for i := range ItemDefs {
		if ItemDefs[i].Name == "incendiary grenade" {
			ItemDefs[i].Name = "molotov cocktail"
			ItemDefs[i].Desc = "a bottle filled with fuel and a burning rag, it spreads fire on impact."
			ItemDefs[i].Color = 12
		}
	}
}

// ActiveKeyBindings is the reduced, player-facing help table.  The recovered
// KeyBindings data remains available as historical reference.
var ActiveKeyBindings = []KeyBind{
	{Key: '8', Name: "move to north", Desc: "move player north"},
	{Key: '9', Name: "move to north-east", Desc: "move player north-east"},
	{Key: '7', Name: "move to north-west", Desc: "move player north-west"},
	{Key: '6', Name: "move to east", Desc: "move player east"},
	{Key: '5', Name: "wait / rest", Desc: "wait, rest and recover health"},
	{Key: '4', Name: "move to west", Desc: "move player west"},
	{Key: '3', Name: "move to south-east", Desc: "move player south-east"},
	{Key: '2', Name: "move to south", Desc: "move player south"},
	{Key: '1', Name: "move to south-west", Desc: "move player south-west"},
	{Key: 'f', Name: "attack", Desc: "attack an adjacent enemy with the equipped melee weapon"},
	{Key: 'b', Name: "smash", Desc: "smash an adjacent door, wall or piece of furniture"},
	{Key: 't', Name: "throw", Desc: "throw an HE grenade or molotov within sight range"},
	{Key: ',', Name: "fast pick up", Desc: "pick up items from the current tile"},
	{Key: 'i', Name: "inventory", Desc: "display inventory and equipment"},
	{Key: 'l', Name: "look", Desc: "look around and examine the surroundings"},
	{Key: 'm', Name: "message buffer", Desc: "display the message history"},
	{Key: '@', Name: "player screen", Desc: "display character information"},
	{Key: 'x', Name: "skills", Desc: "display and improve skills"},
	{Key: 'O', Name: "options", Desc: "display the options screen"},
	{Key: '?', Name: "help", Desc: "display this key binding screen"},
	{Key: 'S', Name: "save and quit", Desc: "save the game and quit"},
	{Key: '>', Name: "use portal", Desc: "use a portal, stairs or elevator"},
	{Key: ']', Name: "next map", Desc: "jump directly to the next map"},
	{Key: '[', Name: "previous map", Desc: "jump directly to the previous map"},
	{Key: 'Q', Name: "quit", Desc: "quit the game"},
}

// ActiveHints is the melee-focused hint list.  The recovered hints that refer
// to ranged weapons, ammo or searching are no longer relevant.
var ActiveHints = []string{
	Hints[0], // medkit
	Hints[1], // fight one enemy at a time
	Hints[2], // change tactics
	Hints[4], // use the surroundings
	Hints[6], // inventory takes time
	Hints[7], // acid ignores armor
	"you are defenceless without a weapon; keep a melee weapon equipped.",
	"grenades and molotovs can only be thrown at tiles within your sight range.",
}
