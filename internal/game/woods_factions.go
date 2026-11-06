package game

var (
	villagerProto       int
	militiaProto        int
	rabbitProto         int
	foxProto            int
	boarProto           int
	crowProto           int
	frogProto           int
	neutralWoodsAnimals []int
	villagerProtos      []int
	woodsCritterRoles   = map[int]ColorRole{}
)

func init() {
	add := func(d CritterDef) int {
		id := len(CritterDefs)
		CritterDefs = append(CritterDefs, d)
		return id
	}

	villagerProto = add(CritterDef{
		Name: "villager", Class: 2, Align: AlignVillager, Color: 7, Glyph: 'v',
		STR: 6, TGH: 6, AGI: 6, INT: 6, SPD: 2,
		Skills: [NumSkills]int{40, 45, 40, 45, 0, 0, 40, 50, 20},
		AI:     AIVillager, Level: 1, XP: 20,
	})
	militiaProto = add(CritterDef{
		Name: "militia", Class: 2, Align: AlignVillager, Color: 10, Glyph: 'V',
		STR: 8, TGH: 8, AGI: 7, INT: 7, SPD: 2,
		Skills: [NumSkills]int{45, 50, 55, 65, 0, 0, 45, 55, 25},
		AI:     AIVillager, Level: 3, XP: 35,
	})
	rabbitProto = add(CritterDef{
		Name: "rabbit", Class: 0, Align: AlignNeutral, Color: 7, Glyph: 'r',
		STR: 2, TGH: 2, AGI: 14, INT: 2, SPD: 4,
		Skills: [NumSkills]int{50, 80, 20, 10, 0, 0, 20, 15, 10},
		AI:     AIAnimal, Level: 1, XP: 10,
	})
	foxProto = add(CritterDef{
		Name: "fox", Class: 0, Align: AlignNeutral, Color: 6, Glyph: 'f',
		STR: 4, TGH: 3, AGI: 13, INT: 3, SPD: 3,
		Skills: [NumSkills]int{45, 70, 35, 25, 0, 0, 30, 20, 10},
		AI:     AIAnimal, Level: 2, XP: 20,
	})
	boarProto = add(CritterDef{
		Name: "boar", Class: 0, Align: AlignNeutral, Color: 6, Glyph: 'b',
		STR: 8, TGH: 7, AGI: 7, INT: 2, SPD: 2,
		Skills: [NumSkills]int{35, 50, 55, 35, 0, 0, 20, 25, 10},
		AI:     AIAnimal, Level: 4, XP: 40,
	})
	crowProto = add(CritterDef{
		Name: "crow", Class: 0, Align: AlignNeutral, Color: 8, Glyph: 'c',
		STR: 2, TGH: 2, AGI: 15, INT: 2, SPD: 4,
		Skills: [NumSkills]int{55, 85, 15, 10, 0, 0, 25, 10, 10},
		AI:     AIAnimal, Level: 1, XP: 10,
	})
	frogProto = add(CritterDef{
		Name: "frog", Class: 0, Align: AlignNeutral, Color: 10, Glyph: 'g',
		STR: 1, TGH: 1, AGI: 8, INT: 1, SPD: 1,
		Skills: [NumSkills]int{30, 60, 10, 5, 0, 0, 10, 10, 5},
		AI:     AIAnimal, Level: 1, XP: 5,
	})

	neutralWoodsAnimals = []int{rabbitProto, foxProto, boarProto, crowProto, frogProto, 15} // 15 = recovered deer
	villagerProtos = []int{villagerProto, militiaProto}
	woodsCritterRoles = map[int]ColorRole{
		villagerProto: RoleVillager,
		militiaProto:  RoleVillager,
		rabbitProto:   RoleRabbit,
		foxProto:      RoleFox,
		boarProto:     RoleBoar,
		crowProto:     RoleCrow,
		frogProto:     RoleFrog,
	}
}
