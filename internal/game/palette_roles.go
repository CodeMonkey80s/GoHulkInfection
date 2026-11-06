package game

// TileColorRole returns the woods palette role for a recovered tile ID.
// Tiles that belong to other themes still get a role so they can be recoloured
// consistently if they are reused by the woods generator later.
func TileColorRole(tile int) (ColorRole, bool) {
	if role, ok := woodsTileRoles[tile]; ok {
		return role, true
	}
	switch tile {
	case 0, 1: // space void, void
		return RoleMapBackground, true
	case 2: // grass
		return RoleGrassBase, true
	case 3: // tree
		return RoleCanopyMid, true
	case 4: // big tree
		return RoleCanopyLight, true
	case 5: // swamp tree
		return RoleDeadTree, true
	case 6: // swamp
		return RoleMarsh, true
	case 7: // shallow water
		return RoleShallowWater, true
	case 8: // deep water
		return RoleDeepWater, true
	case 9: // stone floor
		return RoleRock, true
	case 10, 11: // stone wall, solid wall
		return RoleStoneWall, true
	case 12: // stone pillar
		return RoleStoneLight, true
	case 13: // wooden floor
		return RoleWoodLight, true
	case 14: // wooden wall
		return RoleWoodWall, true
	case 15, 16: // wooden door closed/open
		return RoleDoor, true
	case 17, 18: // wooden window, broken window
		return RoleWindowFrame, true
	case 19: // hole in the floor
		return RoleMud, true
	case 20, 21, 22: // golden floor/door
		return RoleGoldRelic, true
	case 23: // stone stairs
		return RoleStoneLight, true
	case 24: // time portal
		return RoleGoldRelic, true
	case 25: // plasteel floor
		return RoleSilver, true
	case 26: // plasteel floor
		return RoleSilver, true
	case 27: // plasteel corridor floor
		return RoleIronDark, true
	case 28, 29, 30, 31: // plasteel walls
		return RoleSteel, true
	case 32: // stone wall
		return RoleStoneWall, true
	case 33, 34, 36, 37, 39, 40: // plasteel doors/hatches
		return RoleSteel, true
	case 35, 38: // plasteel air locks
		return RoleGoldRelic, true
	case 41, 42, 43, 44, 45, 46: // pipes
		return RoleIronDark, true
	case 47, 49: // fungus
		return RolePoison, true
	case 48: // fungus
		return RoleAsh, true
	case 50: // plant
		return RoleBush, true
	case 51: // plant
		return RoleFlowerPink, true
	case 52: // plant
		return RoleStump, true
	case 53: // pool of water
		return RoleShallowWater, true
	case 54: // pool of organic liquid
		return RolePoison, true
	case 55: // pool of fuel
		return RoleMud, true
	case 56: // cryogenic wall
		return RoleFoam, true
	case 57, 58: // plasteel door
		return RoleSteel, true
	case 59: // elevator
		return RoleSteel, true
	case 60, 61, 62, 63: // portals
		return RoleGoldRelic, true
	default:
		return "", false
	}
}

// ItemColorRole returns the palette role for an item prototype. Unknown or
// future prototypes fall back to a role based on their recovered kind.
func ItemColorRole(proto int) (ColorRole, bool) {
	if proto < 0 || proto >= len(ItemDefs) {
		return "", false
	}
	switch proto {
	case 0: // corpse
		return RoleDarkBlood, true
	case 1: // head
		return RoleBone, true
	case 2, 3, 4: // arm, leg, torso
		return RoleGuts, true
	case 5: // jacket
		return RoleLeather, true
	case 6: // flak armor
		return RoleIronDark, true
	case 7: // power jacket
		return RoleSteel, true
	case 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		return RoleSteel, true
	case 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33:
		return RoleSteel, true
	case 34, 35, 36, 37, 40:
		return RoleGoldRelic, true
	case 38: // flamer canister
		return RoleFireEmber, true
	case 39: // acid canister
		return RoleAcid, true
	case 41, 42: // HE and incendiary grenades
		return RoleFireEmber, true
	case 43: // acid grenade
		return RoleAcid, true
	case 44: // smoke grenade
		return RoleSmoke, true
	case 45: // flashbang
		return RoleFlameCore, true
	case 46, 48: // first aid kit, medkit
		return RoleHealthGood, true
	case 47: // bandage
		return RoleCloth, true
	}
	switch ItemDefs[proto].Kind {
	case KindArmor:
		return RoleLeather, true
	case KindMelee, KindRanged:
		return RoleSteel, true
	case KindAmmo:
		return RoleGoldRelic, true
	case KindGrenade:
		return RoleFireEmber, true
	case KindMedical:
		return RoleHealthGood, true
	case KindBodyPart:
		return RoleGuts, true
	}
	return "", false
}

// CritterColorRole returns the palette role for a critter prototype.
func CritterColorRole(proto int) (ColorRole, bool) {
	if proto < 0 || proto >= len(CritterDefs) {
		return "", false
	}
	switch proto {
	case 0: // zombie
		return RoleZombieSkin, true
	case 1: // toxic zombie
		return RolePoison, true
	case 2: // zombie punk
		return RoleRottenFlesh, true
	case 3: // zombie bomber
		return RoleZombieSkin, true
	case 4: // mutant zombie
		return RoleZombieSkin, true
	case 5: // biotech nightmare
		return RoleAcid, true
	case 6: // abomination
		return RoleRottenFlesh, true
	case 7, 8, 9, 10: // observers and tech-spiders
		return RoleSilver, true
	case 11: // warp wolf
		return RoleWolf, true
	case 12: // meta fly
		return RoleFreshBlood, true
	case 13: // meta worm
		return RolePoison, true
	case 14: // wolf
		return RoleWolf, true
	case 15: // deer
		return RoleDeer, true
	case 16: // you
		return RolePlayer, true
	}
	return "", false
}
