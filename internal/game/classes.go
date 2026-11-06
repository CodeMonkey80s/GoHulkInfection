package game

// PlayerClasses are the playable survivor classes. ClassDefs in data_gen.go
// keeps the recovered original class table for reference; this table is the
// current player-facing roster.
var PlayerClasses = []ClassDef{
	{
		Name: "Butcher",
		Desc: "a slaughterhouse worker with a heavy arm and no patience left. built to stand in a doorway and make the dead come to him.",
		STR:  14, TGH: 13, AGI: 8, INT: 6,
		Skills: [NumSkills]int{35, 40, 70, 90, 0, 0, 35, 55, 15},
	},
	{
		Name: "Hunter",
		Desc: "a poacher who knows the woods better than the people in them. fast, quiet and always looking for an angle.",
		STR:  10, TGH: 9, AGI: 13, INT: 8,
		Skills: [NumSkills]int{70, 75, 35, 65, 0, 40, 60, 45, 30},
	},
	{
		Name: "Medic",
		Desc: "a field medic who has stitched more bites than bullet wounds. keeps the group alive when the line breaks.",
		STR:  8, TGH: 11, AGI: 9, INT: 13,
		Skills: [NumSkills]int{55, 50, 30, 55, 0, 0, 50, 80, 45},
	},
	{
		Name: "Mechanic",
		Desc: "a garage hand who learned that a locked door is only a suggestion. strong enough to swing a wrench and smart enough to keep it working.",
		STR:  11, TGH: 12, AGI: 8, INT: 11,
		Skills: [NumSkills]int{45, 40, 55, 75, 0, 0, 55, 55, 25},
	},
	{
		Name: "Prepper",
		Desc: "a survivalist who was right all along. no single specialty, just enough of everything to keep moving.",
		STR:  9, TGH: 10, AGI: 11, INT: 10,
		Skills: [NumSkills]int{60, 60, 40, 60, 0, 30, 60, 65, 40},
	},
}

// StartingKit is the equipment a class begins with. Proto values are resolved
// when the game starts so appended item IDs such as the armour vial are safe.
type StartingKit struct {
	Hands int
	Body  int
	Belt  []int
}

func ClassStartingKit(class int) StartingKit {
	switch class {
	case 0: // Butcher
		return StartingKit{Hands: 11, Body: 5, Belt: []int{47, ArmourVialProto}}
	case 1: // Hunter
		return StartingKit{Hands: 10, Body: 5, Belt: []int{47, ArmourVialProto}}
	case 2: // Medic
		return StartingKit{Hands: 8, Body: 5, Belt: []int{48, 47}}
	case 3: // Mechanic
		return StartingKit{Hands: 12, Body: 5, Belt: []int{47, ArmourVialProto}}
	case 4: // Prepper
		return StartingKit{Hands: 10, Body: 5, Belt: []int{48, 47}}
	default:
		return StartingKit{Hands: 10, Body: 5, Belt: []int{47, ArmourVialProto}}
	}
}

func startingItem(proto int) *Item {
	if proto < 0 || proto >= len(ItemDefs) {
		return nil
	}
	it := &Item{Proto: proto, X: -1, Y: -1}
	if d := it.Def(); d != nil && d.Kind == KindMedical {
		it.Charges = d.Charges
	}
	return it
}

func startingKitSummary(class int) string {
	kit := ClassStartingKit(class)
	name := func(proto int) string {
		if proto < 0 || proto >= len(ItemDefs) {
			return "-"
		}
		return ItemDefs[proto].Name
	}
	summary := "hands: " + name(kit.Hands) + "   body: " + name(kit.Body)
	if len(kit.Belt) > 0 {
		summary += "   belt: " + name(kit.Belt[0])
	}
	if len(kit.Belt) > 1 {
		summary += ", " + name(kit.Belt[1])
	}
	return summary
}
