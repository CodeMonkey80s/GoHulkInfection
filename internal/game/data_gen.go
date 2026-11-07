// Code generated from hi030.exe data tables. DO NOT EDIT.

package game

// ItemDef describes one recovered item prototype.
type ItemDef struct {
	Name, Desc                         string
	Kind                               ItemKind
	Weight                             int
	Damage                             int
	WeaponSub                          int
	Ranged                             bool
	AmmoType                           int
	Clip                               int
	MaxAmmo                            int
	Burst                              int
	Reload                             int
	Range                              int
	Accuracy                           int
	ArmorMelee, ArmorRanged, ArmorAcid int
	Charges                            int
	Heal                               int
	Radius                             int
	Glyph                              byte
	Color                              uint8
}

var ItemDefs = []ItemDef{
	{Name: "corpse", Desc: "", Kind: KindBodyPart, Weight: 370, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 37, Color: 4},
	{Name: "head", Desc: "", Kind: KindBodyPart, Weight: 390, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 104, Color: 4},
	{Name: "arm", Desc: "", Kind: KindBodyPart, Weight: 960, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 97, Color: 4},
	{Name: "leg", Desc: "", Kind: KindBodyPart, Weight: 960, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 108, Color: 4},
	{Name: "torso", Desc: "", Kind: KindBodyPart, Weight: 440, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 116, Color: 4},
	{Name: "jacket", Desc: "military jacket, provides the wearer with basic protection.", Kind: KindArmor, Weight: 5600, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 80, ArmorRanged: 80, ArmorAcid: 80, Charges: 0, Heal: 0, Radius: 0, Glyph: 91, Color: 6},
	{Name: "flak armor", Desc: "", Kind: KindArmor, Weight: 7800, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 70, ArmorRanged: 70, ArmorAcid: 70, Charges: 0, Heal: 0, Radius: 0, Glyph: 91, Color: 6},
	{Name: "power jacket", Desc: "", Kind: KindArmor, Weight: 9700, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 60, ArmorRanged: 60, ArmorAcid: 60, Charges: 0, Heal: 0, Radius: 0, Glyph: 91, Color: 6},
	{Name: "pipe", Desc: "", Kind: KindMelee, Weight: 2500, Damage: 47, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 124, Color: 7},
	{Name: "plasteel stick", Desc: "", Kind: KindMelee, Weight: 2700, Damage: 47, WeaponSub: 1, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 124, Color: 7},
	{Name: "crowbar", Desc: "", Kind: KindMelee, Weight: 1700, Damage: 47, WeaponSub: 2, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 124, Color: 7},
	{Name: "hammer", Desc: "", Kind: KindMelee, Weight: 2100, Damage: 47, WeaponSub: 3, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 124, Color: 7},
	{Name: "wrench", Desc: "", Kind: KindMelee, Weight: 1500, Damage: 47, WeaponSub: 4, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 124, Color: 7},
	{Name: "chainsaw", Desc: "strange sculptured chainsaw, it is decorated with figures of human heads, skulls and demon faces.", Kind: KindMelee, Weight: 8500, Damage: 47, WeaponSub: 5, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 2, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 124, Color: 7},
	{Name: "power sword", Desc: "", Kind: KindMelee, Weight: 1320, Damage: 124, WeaponSub: 6, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 47, Color: 15},
	{Name: "meta dagger", Desc: "", Kind: KindMelee, Weight: 550, Damage: 124, WeaponSub: 7, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 47, Color: 15},
	{Name: "vibra knife", Desc: "military knife, vibratng blade does more damage than standard knife.", Kind: KindMelee, Weight: 1050, Damage: 124, WeaponSub: 8, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 47, Color: 15},
	{Name: "multiaxe", Desc: "", Kind: KindMelee, Weight: 2910, Damage: 124, WeaponSub: 9, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 47, Color: 15},
	{Name: "shock mace", Desc: "", Kind: KindMelee, Weight: 2120, Damage: 124, WeaponSub: 10, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 47, Color: 15},
	{Name: "battle hammer", Desc: "", Kind: KindMelee, Weight: 3910, Damage: 124, WeaponSub: 11, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 1, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 47, Color: 15},
	{Name: "lasgun", Desc: "", Kind: KindRanged, Weight: 7000, Damage: 125, WeaponSub: 12, Ranged: true, AmmoType: 49, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 8, Accuracy: 4, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "plasma gun", Desc: "", Kind: KindRanged, Weight: 5300, Damage: 125, WeaponSub: 13, Ranged: true, AmmoType: 34, Clip: 320, MaxAmmo: 320, Burst: 7, Reload: 3, Range: 5, Accuracy: 5, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "personal pistol", Desc: "9mm military pistol.", Kind: KindRanged, Weight: 2500, Damage: 125, WeaponSub: 14, Ranged: true, AmmoType: 35, Clip: 18, MaxAmmo: 18, Burst: 1, Reload: 2, Range: 4, Accuracy: 3, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "flamer", Desc: "", Kind: KindRanged, Weight: 9500, Damage: 43, WeaponSub: 15, Ranged: true, AmmoType: 38, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 4, Range: 10, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "acid spitter", Desc: "", Kind: KindRanged, Weight: 8500, Damage: 43, WeaponSub: 16, Ranged: true, AmmoType: 39, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 4, Range: 11, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "lightning gun", Desc: "", Kind: KindRanged, Weight: 5600, Damage: 43, WeaponSub: 17, Ranged: true, AmmoType: 49, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 2, Range: 8, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "g-thrower", Desc: "", Kind: KindRanged, Weight: 6200, Damage: 125, WeaponSub: 18, Ranged: true, AmmoType: 40, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 3, Range: 12, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "black hole generator", Desc: "", Kind: KindRanged, Weight: 6200, Damage: 125, WeaponSub: 19, Ranged: true, AmmoType: 49, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 4, Range: 8, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "freeze gun", Desc: "", Kind: KindRanged, Weight: 8100, Damage: 125, WeaponSub: 19, Ranged: true, AmmoType: 49, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 2, Range: 8, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "ad-smg", Desc: "advanced submachine gun, fires with 9mm rounds. automatic burst fire mode only.", Kind: KindRanged, Weight: 3400, Damage: 125, WeaponSub: 14, Ranged: true, AmmoType: 36, Clip: 40, MaxAmmo: 40, Burst: 4, Reload: 2, Range: 4, Accuracy: 3, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "desintegrator", Desc: "", Kind: KindRanged, Weight: 13400, Damage: 125, WeaponSub: 19, Ranged: true, AmmoType: 49, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 2, Range: 8, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "subparticle gun", Desc: "", Kind: KindRanged, Weight: 9800, Damage: 43, WeaponSub: 19, Ranged: true, AmmoType: 49, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 2, Range: 8, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "warp cannon", Desc: "", Kind: KindRanged, Weight: 14800, Damage: 43, WeaponSub: 19, Ranged: true, AmmoType: 49, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 2, Range: 8, Accuracy: 9, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "mini gun", Desc: "", Kind: KindRanged, Weight: 14800, Damage: 43, WeaponSub: 14, Ranged: true, AmmoType: 37, Clip: 0, MaxAmmo: 0, Burst: 5, Reload: 3, Range: 4, Accuracy: 3, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 125, Color: 11},
	{Name: "plasma gun clip", Desc: "", Kind: KindAmmo, Weight: 780, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 34, Clip: 0, MaxAmmo: 320, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 61, Color: 14},
	{Name: "personal pistol clip", Desc: "9mm ammunition for military pistol, holds 18 rounds.", Kind: KindAmmo, Weight: 450, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 35, Clip: 0, MaxAmmo: 18, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 61, Color: 14},
	{Name: "ad-smg clip", Desc: "9mm ammunition for advanced submachine gun, holds 40 rounds.", Kind: KindAmmo, Weight: 950, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 36, Clip: 0, MaxAmmo: 40, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 61, Color: 14},
	{Name: "mini gun clip", Desc: "", Kind: KindAmmo, Weight: 7100, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 37, Clip: 0, MaxAmmo: 2000, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 61, Color: 14},
	{Name: "flamer canister", Desc: "", Kind: KindAmmo, Weight: 9100, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 38, Clip: 0, MaxAmmo: 11, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 61, Color: 14},
	{Name: "acid canister", Desc: "", Kind: KindAmmo, Weight: 9100, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 39, Clip: 0, MaxAmmo: 11, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 61, Color: 14},
	{Name: "g-thrower clip", Desc: "", Kind: KindAmmo, Weight: 3300, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 40, Clip: 0, MaxAmmo: 3, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 61, Color: 14},
	{Name: "he grenade", Desc: "defensive military high-explosive grenade, during explosion it cracks down to 800 deadly pieces", Kind: KindGrenade, Weight: 850, Damage: 59, WeaponSub: 18, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 12, Accuracy: 6, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 12, Glyph: 42, Color: 12},
	{Name: "incendiary grenade", Desc: "", Kind: KindGrenade, Weight: 950, Damage: 59, WeaponSub: 15, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 13, Accuracy: 7, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 13, Glyph: 42, Color: 12},
	{Name: "acid grenade", Desc: "", Kind: KindGrenade, Weight: 950, Damage: 59, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 42, Color: 12},
	{Name: "smoke grenade", Desc: "", Kind: KindGrenade, Weight: 850, Damage: 59, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 42, Color: 12},
	{Name: "flashbang grenade", Desc: "", Kind: KindGrenade, Weight: 650, Damage: 59, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 0, MaxAmmo: 0, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 0, Heal: 0, Radius: 0, Glyph: 42, Color: 12},
	{Name: "first aid kit", Desc: "small aid kit", Kind: KindMedical, Weight: 900, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 5, MaxAmmo: 13, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 5, Heal: 13, Radius: 0, Glyph: 43, Color: 10},
	{Name: "bandage", Desc: "heap of bandages", Kind: KindMedical, Weight: 300, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 1, MaxAmmo: 9, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 1, Heal: 9, Radius: 0, Glyph: 43, Color: 10},
	{Name: "medkit", Desc: "medkit", Kind: KindMedical, Weight: 1200, Damage: 0, WeaponSub: 0, Ranged: false, AmmoType: 0, Clip: 5, MaxAmmo: 25, Burst: 1, Reload: 1, Range: 0, Accuracy: 0, ArmorMelee: 0, ArmorRanged: 0, ArmorAcid: 0, Charges: 5, Heal: 25, Radius: 0, Glyph: 43, Color: 10},
}

// WeaponDef is one entry of the recovered weapon table at 0x898a0.
type WeaponDef struct {
	AV, DV, DmgMin, DmgMax int
}

var WeaponDefs = []WeaponDef{
	{AV: 2, DV: 2, DmgMin: 1, DmgMax: 5},
	{AV: 6, DV: 0, DmgMin: 1, DmgMax: 6},
	{AV: 2, DV: -1, DmgMin: 2, DmgMax: 4},
	{AV: 2, DV: -3, DmgMin: 3, DmgMax: 5},
	{AV: 4, DV: -2, DmgMin: 2, DmgMax: 4},
	{AV: 5, DV: -12, DmgMin: 7, DmgMax: 44},
	{AV: 5, DV: 0, DmgMin: 1, DmgMax: 8},
	{AV: 17, DV: 6, DmgMin: 3, DmgMax: 7},
	{AV: 15, DV: 5, DmgMin: 2, DmgMax: 6},
	{AV: 7, DV: 2, DmgMin: 3, DmgMax: 21},
	{AV: 9, DV: -5, DmgMin: 5, DmgMax: 18},
	{AV: 11, DV: -7, DmgMin: 4, DmgMax: 16},
	{AV: 11, DV: -2, DmgMin: 7, DmgMax: 13},
	{AV: 5, DV: -1, DmgMin: 6, DmgMax: 17},
	{AV: 7, DV: -2, DmgMin: 3, DmgMax: 11},
	{AV: 0, DV: 0, DmgMin: 1, DmgMax: 7},
	{AV: 0, DV: 0, DmgMin: 1, DmgMax: 7},
	{AV: 0, DV: 0, DmgMin: 4, DmgMax: 59},
	{AV: 0, DV: 0, DmgMin: 7, DmgMax: 55},
	{AV: 0, DV: 0, DmgMin: 0, DmgMax: 0},
}

// CritterDef describes one recovered critter prototype.
type CritterDef struct {
	Name                    string
	Class, Align            int
	Color                   uint8
	Glyph                   byte
	STR, TGH, AGI, INT, SPD int
	Skills                  [NumSkills]int
	AI                      int
	Level, XP               int
	DmgCount, DmgSides      int
	DmgBonus                int
}

var CritterDefs = []CritterDef{
	{Name: "zombie", Class: 1, Align: 1, Color: 7, Glyph: 122, STR: 13, TGH: 6, AGI: 4, INT: 2, SPD: 1, Skills: [9]int{20, 45, 86, 30, 0, 0, 10, 20, 10}, AI: 0, Level: 5, XP: 85},
	{Name: "toxic zombie", Class: 1, Align: 1, Color: 10, Glyph: 122, STR: 15, TGH: 8, AGI: 4, INT: 3, SPD: 1, Skills: [9]int{20, 35, 90, 30, 0, 0, 10, 20, 10}, AI: 0, Level: 5, XP: 89},
	{Name: "zombie punk", Class: 1, Align: 1, Color: 13, Glyph: 122, STR: 11, TGH: 5, AGI: 6, INT: 2, SPD: 3, Skills: [9]int{25, 42, 77, 55, 40, 0, 15, 20, 15}, AI: 0, Level: 5, XP: 91},
	{Name: "zombie bomber", Class: 1, Align: 1, Color: 14, Glyph: 122, STR: 10, TGH: 5, AGI: 4, INT: 2, SPD: 2, Skills: [9]int{20, 37, 72, 30, 30, 0, 15, 20, 10}, AI: 0, Level: 5, XP: 84},
	{Name: "mutant zombie", Class: 1, Align: 1, Color: 7, Glyph: 90, STR: 19, TGH: 18, AGI: 6, INT: 2, SPD: 1, Skills: [9]int{25, 57, 95, 60, 0, 0, 15, 25, 15}, AI: 0, Level: 5, XP: 121},
	{Name: "biotech nightmare", Class: 1, Align: 1, Color: 14, Glyph: 66, STR: 16, TGH: 15, AGI: 7, INT: 4, SPD: 2, Skills: [9]int{25, 50, 85, 60, 0, 0, 15, 25, 15}, AI: 0, Level: 25, XP: 272},
	{Name: "abomination", Class: 1, Align: 1, Color: 15, Glyph: 65, STR: 24, TGH: 36, AGI: 5, INT: 3, SPD: 1, Skills: [9]int{30, 47, 95, 70, 0, 0, 20, 40, 25}, AI: 0, Level: 40, XP: 444},
	{Name: "observer mainbrain", Class: 1, Align: 1, Color: 15, Glyph: 79, STR: 10, TGH: 12, AGI: 6, INT: 15, SPD: 2, Skills: [9]int{45, 40, 50, 40, 80, 0, 40, 30, 35}, AI: 2, Level: 30, XP: 314},
	{Name: "observer drone", Class: 1, Align: 1, Color: 15, Glyph: 111, STR: 6, TGH: 6, AGI: 10, INT: 8, SPD: 3, Skills: [9]int{40, 45, 45, 30, 65, 0, 30, 25, 25}, AI: 2, Level: 15, XP: 173},
	{Name: "tech-spider mainbrain", Class: 1, Align: 1, Color: 15, Glyph: 83, STR: 12, TGH: 14, AGI: 8, INT: 12, SPD: 2, Skills: [9]int{45, 50, 70, 60, 75, 0, 35, 30, 30}, AI: 2, Level: 35, XP: 363},
	{Name: "tech-spider drone", Class: 1, Align: 1, Color: 15, Glyph: 115, STR: 7, TGH: 7, AGI: 11, INT: 6, SPD: 3, Skills: [9]int{40, 55, 55, 35, 60, 0, 30, 25, 25}, AI: 2, Level: 18, XP: 201},
	{Name: "warp wolf", Class: 0, Align: 1, Color: 6, Glyph: 100, STR: 10, TGH: 9, AGI: 14, INT: 4, SPD: 4, Skills: [9]int{35, 60, 60, 40, 0, 0, 25, 25, 20}, AI: 3, Level: 22, XP: 232},
	{Name: "meta fly", Class: 0, Align: 1, Color: 12, Glyph: 102, STR: 3, TGH: 3, AGI: 18, INT: 1, SPD: 5, Skills: [9]int{30, 70, 45, 20, 0, 0, 10, 15, 10}, AI: 3, Level: 10, XP: 112},
	{Name: "meta worm", Class: 0, Align: 1, Color: 14, Glyph: 119, STR: 8, TGH: 10, AGI: 5, INT: 1, SPD: 2, Skills: [9]int{20, 35, 55, 20, 0, 0, 10, 20, 10}, AI: 3, Level: 15, XP: 165},
	{Name: "wolf", Class: 0, Align: 2, Color: 7, Glyph: 119, STR: 4, TGH: 3, AGI: 11, INT: 3, SPD: 3, Skills: [9]int{30, 60, 50, 25, 0, 0, 15, 15, 10}, AI: 3, Level: 5, XP: 73},
	{Name: "deer", Class: 0, Align: 2, Color: 15, Glyph: 100, STR: 3, TGH: 2, AGI: 13, INT: 3, SPD: 3, Skills: [9]int{35, 70, 20, 10, 0, 0, 20, 10, 10}, AI: 3, Level: 3, XP: 50},
	{Name: "you", Class: 4, Align: 0, Color: 14, Glyph: 64, STR: 5, TGH: 5, AGI: 5, INT: 5, SPD: 2, Skills: [9]int{25, 25, 25, 25, 25, 0, 25, 15, 0}, AI: 0, Level: 1, XP: 40},
}

// TileDef describes one recovered tile prototype.
type TileDef struct {
	Name                  string
	Type                  int
	Glyph                 byte
	Color                 uint8
	BlockMove, BlockSight bool
	Partner               int
}

var TileDefs = []TileDef{
	{Name: "space void", Type: 0, Glyph: 32, Color: 0, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "void", Type: 0, Glyph: 32, Color: 0, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "grass", Type: 1, Glyph: 46, Color: 2, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "tree", Type: 1, Glyph: 116, Color: 2, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "big tree", Type: 1, Glyph: 84, Color: 10, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "swamp tree", Type: 1, Glyph: 89, Color: 7, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "swamp", Type: 1, Glyph: 61, Color: 8, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "shallow water", Type: 1, Glyph: 61, Color: 9, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "deep water", Type: 1, Glyph: 61, Color: 1, BlockMove: true, BlockSight: false, Partner: -1},
	{Name: "stone floor", Type: 1, Glyph: 46, Color: 7, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "stone wall", Type: 2, Glyph: 176, Color: 7, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "solid wall", Type: 2, Glyph: 176, Color: 2, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "stone pillar", Type: 2, Glyph: 73, Color: 8, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "wooden floor", Type: 1, Glyph: 46, Color: 6, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "wooden wall", Type: 1, Glyph: 35, Color: 6, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "wooden door", Type: 3, Glyph: 254, Color: 6, BlockMove: true, BlockSight: true, Partner: 16},
	{Name: "wooden door", Type: 3, Glyph: 47, Color: 6, BlockMove: false, BlockSight: false, Partner: 15},
	{Name: "wooden window", Type: 6, Glyph: 43, Color: 6, BlockMove: true, BlockSight: false, Partner: 18},
	{Name: "wooden broken window", Type: 7, Glyph: 43, Color: 6, BlockMove: true, BlockSight: false, Partner: -1},
	{Name: "hole in the floor", Type: 1, Glyph: 111, Color: 6, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "golden floor", Type: 1, Glyph: 46, Color: 14, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "golden door", Type: 3, Glyph: 254, Color: 14, BlockMove: true, BlockSight: true, Partner: 22},
	{Name: "golden door", Type: 3, Glyph: 47, Color: 14, BlockMove: false, BlockSight: false, Partner: 21},
	{Name: "stone stairs", Type: 1, Glyph: 62, Color: 7, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "time portal", Type: 1, Glyph: 234, Color: 14, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "plasteel floor", Type: 1, Glyph: 46, Color: 3, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "plasteel floor", Type: 1, Glyph: 46, Color: 5, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "plasteel corridor floor", Type: 1, Glyph: 46, Color: 6, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "plasteel wall", Type: 2, Glyph: 178, Color: 11, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "plasteel wall", Type: 2, Glyph: 177, Color: 13, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "plasteel corridor wall", Type: 2, Glyph: 177, Color: 11, BlockMove: true, BlockSight: true, Partner: -1},
	{Name: "plasteel wall", Type: 4, Glyph: 177, Color: 11, BlockMove: true, BlockSight: true, Partner: 40},
	{Name: "stone wall", Type: 4, Glyph: 176, Color: 7, BlockMove: true, BlockSight: true, Partner: 15},
	{Name: "plasteel door", Type: 3, Glyph: 254, Color: 15, BlockMove: true, BlockSight: true, Partner: 36},
	{Name: "plasteel automatic door", Type: 3, Glyph: 254, Color: 13, BlockMove: true, BlockSight: true, Partner: 37},
	{Name: "plasteel air lock", Type: 254, Glyph: 3, Color: 14, BlockMove: true, BlockSight: true, Partner: 38},
	{Name: "plasteel door", Type: 3, Glyph: 47, Color: 11, BlockMove: false, BlockSight: false, Partner: 33},
	{Name: "plasteel automatic door", Type: 3, Glyph: 47, Color: 15, BlockMove: false, BlockSight: false, Partner: 34},
	{Name: "plasteel air lock", Type: 3, Glyph: 47, Color: 14, BlockMove: false, BlockSight: false, Partner: 35},
	{Name: "plasteel hatch", Type: 3, Glyph: 96, Color: 11, BlockMove: false, BlockSight: false, Partner: 40},
	{Name: "plasteel hatch", Type: 3, Glyph: 35, Color: 3, BlockMove: true, BlockSight: true, Partner: 39},
	{Name: "pipe", Type: 1, Glyph: 205, Color: 7, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "pipe", Type: 1, Glyph: 186, Color: 7, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "pipe", Type: 1, Glyph: 205, Color: 6, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "pipe", Type: 1, Glyph: 186, Color: 6, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "pipe", Type: 1, Glyph: 205, Color: 2, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "pipe", Type: 1, Glyph: 186, Color: 2, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "fungus", Type: 1, Glyph: 34, Color: 2, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "fungus", Type: 1, Glyph: 34, Color: 8, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "fungus", Type: 1, Glyph: 34, Color: 6, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "plant", Type: 1, Glyph: 118, Color: 2, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "plant", Type: 1, Glyph: 118, Color: 4, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "plant", Type: 1, Glyph: 118, Color: 6, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "pool of water", Type: 1, Glyph: 61, Color: 9, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "pool of organic liquid", Type: 1, Glyph: 61, Color: 10, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "pool of fuel", Type: 1, Glyph: 61, Color: 6, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "cryogenic wall", Type: 1, Glyph: 35, Color: 3, BlockMove: true, BlockSight: false, Partner: -1},
	{Name: "plasteel door", Type: 3, Glyph: 47, Color: 3, BlockMove: false, BlockSight: false, Partner: 58},
	{Name: "plasteel door", Type: 3, Glyph: 254, Color: 3, BlockMove: true, BlockSight: false, Partner: 57},
	{Name: "elevator", Type: 5, Glyph: 95, Color: 15, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "portal", Type: 5, Glyph: 143, Color: 14, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "portal", Type: 5, Glyph: 225, Color: 14, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "portal", Type: 5, Glyph: 228, Color: 14, BlockMove: false, BlockSight: false, Partner: -1},
	{Name: "portal", Type: 5, Glyph: 226, Color: 14, BlockMove: false, BlockSight: false, Partner: -1},
}

var KeyBindings = []KeyBind{
	{Key: 56, Mod: 1, Action: 1, Name: "move to north", Desc: "move player to north direction"},
	{Key: 57, Mod: 1, Action: 2, Name: "move to north-east", Desc: "move player to north-east direction"},
	{Key: 55, Mod: 1, Action: 3, Name: "move to north-west", Desc: "move player to north-west direction"},
	{Key: 54, Mod: 1, Action: 4, Name: "move to east", Desc: "move player to north direction"},
	{Key: 53, Mod: 1, Action: 5, Name: "wait / rest / search", Desc: "wait command, causes time to pass, useful for faster health recovery"},
	{Key: 52, Mod: 1, Action: 6, Name: "move to west", Desc: "moving player to west direction"},
	{Key: 51, Mod: 1, Action: 7, Name: "move to south-east", Desc: "moving player to south-east direction"},
	{Key: 50, Mod: 1, Action: 8, Name: "move to south", Desc: "moving player to south direction"},
	{Key: 49, Mod: 1, Action: 9, Name: "move to south-west", Desc: "moving player to south-west direction"},
	{Key: 65, Mod: 1, Action: 10, Name: "activate current weapon", Desc: "enables / activates currently using weapon, some weapons need to be activated before use"},
	{Key: 97, Mod: 1, Action: 11, Name: "aim / select new target", Desc: "selects new target for attack"},
	{Key: 98, Mod: 1, Action: 12, Name: "bash something", Desc: "bash attack, allows player to destroy the surroundings"},
	{Key: 67, Mod: 1, Action: 13, Name: "check ammo", Desc: "displays current ammunition status of using weapon"},
	{Key: 99, Mod: 1, Action: 14, Name: "clear current target", Desc: "cleares current player target"},
	{Key: 87, Mod: 1, Action: 15, Name: "clean yourself", Desc: "cleans player's clothes, face and hands"},
	{Key: 96, Mod: 1, Action: 16, Name: "console", Desc: "displays console / message buffer"},
	{Key: 61, Mod: 1, Action: 17, Name: "dump options", Desc: "dumps option variables to message buffer"},
	{Key: 44, Mod: 1, Action: 18, Name: "fast pick up", Desc: "fast pick up, when only one item is lying on the ground"},
	{Key: 102, Mod: 1, Action: 19, Name: "fight or fire / use weapon", Desc: "fight with hands, melee weapon, using weapon, fire with weapon"},
	{Key: 63, Mod: 1, Action: 20, Name: "help", Desc: "displays help screen"},
	{Key: 104, Mod: 1, Action: 21, Name: "highlight objects", Desc: "highlight spotted critters"},
	{Key: 105, Mod: 1, Action: 22, Name: "display inventory screen", Desc: "displays player inventory screen"},
	{Key: 47, Mod: 1, Action: 23, Name: "show key binding", Desc: "displays current set of keybindings"},
	{Key: 108, Mod: 1, Action: 24, Name: "look around", Desc: "look around / examine surroundings"},
	{Key: 109, Mod: 1, Action: 25, Name: "display message buffer", Desc: "displays message buffer"},
	{Key: 110, Mod: 1, Action: 26, Name: "select next nearest target", Desc: "selects nearest target to the player"},
	{Key: 111, Mod: 1, Action: 27, Name: "operate", Desc: "operates the triggers, doors, terminals and other"},
	{Key: 79, Mod: 1, Action: 28, Name: "display options screen", Desc: "displays options screen"},
	{Key: 64, Mod: 1, Action: 29, Name: "display player screen", Desc: "displays player screen"},
	{Key: 82, Mod: 1, Action: 30, Name: "fast reload", Desc: "quick realoding, without entering the inventory screen"},
	{Key: 83, Mod: 1, Action: 31, Name: "save game and quit", Desc: "saves current game state and quits"},
	{Key: 115, Mod: 1, Action: 32, Name: "search", Desc: "searches the surroundings for something hidden"},
	{Key: 119, Mod: 1, Action: 33, Name: "swap weapons", Desc: "quick swap of previously selected weapons"},
	{Key: 120, Mod: 1, Action: 34, Name: "display player skill levels", Desc: "displays player skill levels"},
	{Key: 84, Mod: 1, Action: 35, Name: "change tactic", Desc: "changes player tactics, causes to change the av and dv values"},
	{Key: 116, Mod: 1, Action: 36, Name: "throw", Desc: "throw grenade"},
	{Key: 62, Mod: 1, Action: 37, Name: "use portal or elevator", Desc: "changes current level to another one"},
	{Key: 259, Mod: 0, Action: 38, Name: "moves cursor north", Desc: "moves cursor north"},
	{Key: 258, Mod: 0, Action: 39, Name: "moves cursor south", Desc: "moves cursor south"},
	{Key: 261, Mod: 0, Action: 40, Name: "moves cursor east", Desc: "moves cursor east"},
	{Key: 260, Mod: 0, Action: 41, Name: "moves cursor west", Desc: "moves cursor west"},
	{Key: 265, Mod: 0, Action: 42, Name: "debug key 1", Desc: "debug key 1"},
	{Key: 266, Mod: 0, Action: 43, Name: "debug key 2", Desc: "debug key 2"},
	{Key: 267, Mod: 0, Action: 44, Name: "debug key 3", Desc: "debug key 3"},
	{Key: 268, Mod: 0, Action: 45, Name: "debug key 4", Desc: "debug key 4"},
	{Key: 269, Mod: 0, Action: 46, Name: "debug key 5", Desc: "debug key 5"},
	{Key: 270, Mod: 0, Action: 47, Name: "debug key 6", Desc: "debug key 6"},
	{Key: 121, Mod: 0, Action: 48, Name: "yes", Desc: "answer yes"},
	{Key: 110, Mod: 0, Action: 49, Name: "no", Desc: "answer no"},
	{Key: 32, Mod: 0, Action: 50, Name: "space", Desc: "space"},
	{Key: 9, Mod: 0, Action: 51, Name: "tab", Desc: "tab"},
	{Key: 13, Mod: 0, Action: 52, Name: "enter", Desc: "enter"},
	{Key: 459, Mod: 0, Action: 53, Name: "enter", Desc: "enter"},
	{Key: 27, Mod: 0, Action: 54, Name: "esc", Desc: "esc"},
	{Key: 113, Mod: 0, Action: 55, Name: "quit screen", Desc: "quit screen"},
	{Key: 122, Mod: 0, Action: 56, Name: "quit screen", Desc: "quit screen"},
	{Key: 90, Mod: 0, Action: 57, Name: "quit screen", Desc: "quit screen"},
	{Key: 81, Mod: 0, Action: 58, Name: "instant quit game", Desc: "instant quit game"},
	{Key: 464, Mod: 0, Action: 59, Name: "scroll page up", Desc: "scroll page up"},
	{Key: 465, Mod: 0, Action: 60, Name: "scroll page down", Desc: "scroll page down"},
	{Key: 458, Mod: 0, Action: 61, Name: "scroll up", Desc: "scroll up"},
	{Key: 463, Mod: 0, Action: 0, Name: "scroll down", Desc: "scroll down"},
}

var OptionDefs = []OptionDef{
	{Name: "optDrawDebugAI", Desc: "displays current state of ai of critters instead of ascii symbols", Value: 0, Min: 1, Max: 1},
	{Name: "optDrawDebugInfo", Desc: "displays debug information on the screen", Value: 0, Min: 0, Max: 1},
	{Name: "optDrawGameMechanicInfo", Desc: "displays game mechanic information on the screen", Value: 0, Min: 0, Max: 1},
	{Name: "optGameChangeLevel", Desc: "when set new level needs to be generated", Value: 0, Min: 0, Max: 1},
	{Name: "optGameFirstLevelCreated", Desc: "stores whether the first level is created", Value: 0, Min: 0, Max: 1},
	{Name: "optGameHintNumber", Desc: "stores number of displayed hint", Value: 0, Min: 0, Max: 0},
	{Name: "optGameHints", Desc: "enables in-game hints", Value: 0, Min: 1, Max: 1},
	{Name: "optGameLevelID", Desc: "stores the id of current level", Value: 0, Min: 0, Max: 27},
	{Name: "optGameSetOptionsForRelease", Desc: "sets proper options for release version", Value: 0, Min: 0, Max: 1},
	{Name: "optGameTheme", Desc: "sets color theme for the game", Value: 0, Min: 0, Max: 2},
	{Name: "optGameUseAI", Desc: "enables ai of critters", Value: 0, Min: 0, Max: 1},
	{Name: "optGameSendDamageToItems", Desc: "when set damage during combat is also send to armor", Value: 0, Min: 1, Max: 1},
	{Name: "optGameUseDamageDescriptions", Desc: "when set damage during combat is displayed", Value: 0, Min: 1, Max: 1},
	{Name: "optGameSpriteAnimationSpeed", Desc: "sets the speed of sprites animations", Value: 0, Min: 1, Max: 100},
	{Name: "optInvColorOutToHeavyItems", Desc: "when set items that too heavy for lift are colored out", Value: 0, Min: 1, Max: 1},
	{Name: "optInvQuitAfterItemUse", Desc: "quits inventory after item is used by player", Value: 0, Min: 1, Max: 1},
	{Name: "optMapCreateBlood", Desc: "when set the blood splashes on the map segments", Value: 0, Min: 1, Max: 1},
	{Name: "optMapCreateCorpses", Desc: "when set corpse is created after killing a critter", Value: 0, Min: 1, Max: 1},
	{Name: "optMapCreateLimbs", Desc: "when set limbs are created after killing a critter", Value: 0, Min: 1, Max: 1},
	{Name: "optMapDefaultExplored", Desc: "sets default value of map explored attribute", Value: 0, Min: 1, Max: 1},
	{Name: "optMapDefaultVisible", Desc: "sets default value of map visible attribute", Value: 0, Min: 1, Max: 1},
	{Name: "optMapDrawHiddenTiles", Desc: "shows all hidden tiles on map", Value: 0, Min: 1, Max: 1},
	{Name: "optMapDrawFullMapName", Desc: "when set a map name is show as full text", Value: 0, Min: 0, Max: 1},
	{Name: "optMapDrawPortals", Desc: "when set the portals are displayed on map", Value: 0, Min: 1, Max: 1},
	{Name: "optMapHighlightSpottedCritters", Desc: "when set spotted critters are highlighted", Value: 1, Min: 1, Max: 1},
	{Name: "optMapResetFOVAfterAction", Desc: "when set resets fov after every player action", Value: 0, Min: 0, Max: 1},
	{Name: "optMapResetCursorPosition", Desc: "when set resets cursor position before using aim and look command to the center of the screen", Value: 0, Min: 0, Max: 1},
	{Name: "optMapShadeColor", Desc: "sets default color of shade", Value: 0, Min: 1, Max: 15},
	{Name: "optMapUseMemoryMap", Desc: "when set map is using a memory map for critters an items found in level", Value: 0, Min: 1, Max: 1},
	{Name: "optMechanicAcidEatItemChance", Desc: "set chance for acid to eat item lying on ground", Value: 0, Min: 30, Max: 100},
	{Name: "optMechanicBaseHitChance", Desc: "sets basic hit chance for unarmed combat, melee combat and ranged attacks", Value: 0, Min: 0, Max: 100},
	{Name: "optMechanicBaseTileDV", Desc: "sets basic hit chance for targetting a tile", Value: 0, Min: 40, Max: 100},
	{Name: "optMechanicAVMod1", Desc: "sets numerator of av calculations", Value: 0, Min: 100, Max: 500},
	{Name: "optMechanicAVMod2", Desc: "sets nominator of av calculations", Value: 0, Min: 200, Max: 500},
	{Name: "optMechanicDVMod1", Desc: "sets numerator of dv calculations", Value: 0, Min: 100, Max: 500},
	{Name: "optMechanicDVMod2", Desc: "sets nominator of av calculations", Value: 0, Min: 200, Max: 500},
	{Name: "optMechanicAVModEnemyIsNear", Desc: "sets reduce factor for av of player when enemy is on adjacent tile when using ranged weapon or throwing objects", Value: 0, Min: 75, Max: 100},
	{Name: "optMechanicLineOfFireHitChance", Desc: "sets percent chance for hitting critter that is in the line of fire", Value: 0, Min: 12, Max: 100},
	{Name: "optPlayerAlwaysSpotCritters", Desc: "when set player automatically spot critter in sight range", Value: 0, Min: 1, Max: 1},
	{Name: "optPlayerAttackOnMove", Desc: "when set player attack critter automatically on move", Value: 0, Min: 1, Max: 1},
	{Name: "optPlayerAutoOpenDoor", Desc: "when set player automaticaly opens door on move", Value: 0, Min: 1, Max: 1},
	{Name: "optPlayerDisplayItemInfoInStatus", Desc: "when set item information is displayed on player status", Value: 0, Min: 1, Max: 1},
	{Name: "optPlayerDisplaySkillValues", Desc: "when set the skill values are shown on the player screen", Value: 1, Min: 1, Max: 1},
	{Name: "optPlayerFreePointsPerLevel", Desc: "sets the number of free points per gained level of experience", Value: 0, Min: 1, Max: 10},
	{Name: "optPlayerClearTargetResetTactics", Desc: "when set clearing the target also resets player tactics", Value: 0, Min: 1, Max: 1},
	{Name: "optPlayerGhostMode", Desc: "when set player is a ghost and can walk through impassable obstacles", Value: 0, Min: 0, Max: 1},
	{Name: "optPlayerGodMode", Desc: "when set player is a god and clound not be harmed", Value: 0, Min: 0, Max: 1},
	{Name: "optPlayerUseTargetLine", Desc: "whet set player use a targetting line for aim", Value: 0, Min: 1, Max: 1},
	{Name: "optPlayerShowAVDV", Desc: "when set av and dv values are shown in status", Value: 0, Min: 1, Max: 1},
	{Name: "optPlayerWins", Desc: "stores whether player wins the game or not", Value: 0, Min: 1, Max: 1},
}

var SkillNames = []string{"perception", "dodging", "unarmed combat", "weapon melee", "weapon ranged", "throwing", "searching", "metabolism", "focusing"}
var SkillDescs = []string{"determines how good you are in examinig the surroundings", "determines your defensive skill value (dv)", "determines how good you are in unarmed combat", "determines how good you are in handling melee weapons", "determines how good you are in handling ranged weapons", "determines how good you are in throwing objects", "determines your ability to search for something hidden", "determines how quickly you recover lost health", "determines how good you are in focusing during fight"}
var SkillLevelNames = []string{"useless", "terrible", "horrible", "bad", "very poor", "poor", "fair", "good", "very good", "great", "excellent", "percect"}
var Hints = []string{"use nanobot medkit to heal extensive wounds", "try to fight only one enemy at a time, use tight corridors as protection", "if fighting gets hard maybie its time to change tactics ?", "when fighting with ranged weapon try to focus your fire on one target", "use the surroundings for your advantage, f.ex. close door before enemy horde approaches", "remember to check your ammo before going into a fire fight", "entering inventory screen always take precious time", "toxic acid ignores the armor"}
var MapSections = []string{
	"###o##o###                              #........+                              #...##.###                              o...#....#                              #...#....#                              ######.###                              #...#....#                              o........o                              #...#....#                              ##o####o##                              ",
	"\u0000\u0000\u0000\u0000\n\u0000\u0000\u0000\n\u0000\u0000\u0000.\u0000\u0000\u0000\r\u0000\u0000\u0000#\u0000\u0000\u0000\u000e\u0000\u0000\u0000+\u0000\u0000\u0000\u000f\u0000\u0000\u0000o\u0000\u0000\u0000\u0011\u0000\u0000\u0000\u0000\u0000\u0000\u0000@\u0000\u0000\u0000\u0001\u0000\u0000\u0000  #########################               #.....#.....#.....#.....#               #.....#.....#.....#.....#               #.....#.....#.....#.....#               #.....#.....#.....#.....#             #####.#####.#####.#####.#####           #...#+#...#+#...#+#...#+#...#           #...........................#           ........................",
	".....           ..I.....I.....I.....I.....I..           \u0000\u0000\u0000\u0000\u001e\u0000\u0000\u0000\n\u0000\u0000\u0000.\u0000\u0000\u0000\t\u0000\u0000\u0000#\u0000\u0000\u0000\n\u0000\u0000\u0000+\u0000\u0000\u0000\u000f\u0000\u0000\u0000I\u0000\u0000\u0000\u000c\u0000\u0000\u0000\u0000\u0000\u0000\u0000@\u0000\u0000\u0000\u0002\u0000\u0000\u0000  ####                                   ##...                                  ##....                                  #.....                                  #.....                                  #.....                                  ##....                                   ##...  ",
	"                                  ####                                                                          \u0000\u0000\u0000\u0000\u0006\u0000\u0000\u0000\n\u0000\u0000\u0000.\u0000\u0000\u0000\t\u0000\u0000\u0000#\u0000\u0000\u0000\n\u0000\u0000\u0000+\u0000\u0000\u0000\u000f\u0000\u0000\u0000\u0000\u0000\u0000\u0000@\u0000\u0000\u0000\u0000\u0000\u0000\u0000@\u0000\u0000\u0000\u0003\u0000\u0000\u0000################                        ...#...........#                        ...#.I.......I.#                        ...#...........#                        ...|...........#                        ...#...........#                ",
}

var RandomNames = []string{"destroyer", "blade", "bone", "cypher", "emperor", "executor", "hector", "impostor", "jesse", "razor", "ripper", "novy", "shiv", "virus", "nuke"}
var LevelThemes = []string{"woods", "caves", "temple of underminer", "time portal"}
var ClassDefs = []ClassDef{
	{Name: "murderer", Desc: "violent punk, he is very strong, has excellent knowledge about using improvised weapons, he has also some skills in using modern weapons", STR: 13, TGH: 13, AGI: 9, INT: 7, Skills: [9]int{45, 50, 70, 90, 65, 0, 45, 50, 10}},
	{Name: "marine", Desc: "former marine, tough training made him deadly, he is very good at using modern weapons and some of melee weapons", STR: 10, TGH: 11, AGI: 10, INT: 8, Skills: [9]int{45, 50, 70, 85, 90, 0, 50, 85, 10}},
	{Name: "warpcrafter", Desc: "a member of mysterious cult named legion, he is trained in using warpcraft - a powerful force from outer universe, he also has skills in using melee weapons", STR: 7, TGH: 10, AGI: 10, INT: 13, Skills: [9]int{45, 40, 55, 50, 50, 0, 45, 45, 20}},
}
