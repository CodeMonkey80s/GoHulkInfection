package game

// The recovered item table contains modern military and science-fiction
// prototypes. This pass keeps the recovered IDs stable but refits the items
// that still exist in normal play to the priest/village/horror theme. Ranged
// weapons, ammunition and old body armour remain raw data only and are filtered
// out by ItemAvailable.
type itemThemeOverride struct {
	Name  string
	Desc  string
	Glyph byte
	Color uint8
}

var itemThemeOverrides = map[string]itemThemeOverride{
	"jacket": {
		Name: "padded jerkin", Desc: "a layered cloth jerkin with old stains.",
		Glyph: '[', Color: 6,
	},
	"flak armor": {
		Name: "reinforced hide coat", Desc: "a heavy coat reinforced with boiled hide.",
		Glyph: '[', Color: 6,
	},
	"power jacket": {
		Name: "penitent's mail", Desc: "a heavy mail shirt worn as an act of penance.",
		Glyph: '[', Color: 7,
	},
	"plasteel stick": {
		Name: "iron-shod staff", Desc: "a village staff capped with iron.",
		Glyph: '|', Color: 6,
	},
	"chainsaw": {
		Name: "relic chainblade", Desc: "a strange bladed chain relic decorated with demon faces.",
		Glyph: '/', Color: 8,
	},
	"power sword": {
		Name: "silvered longsword", Desc: "an old longsword with a silvered edge.",
		Glyph: '/', Color: 15,
	},
	"meta dagger": {
		Name: "consecrated dagger", Desc: "a small dagger marked with old prayers.",
		Glyph: '/', Color: 15,
	},
	"vibra knife": {
		Name: "blessed knife", Desc: "a knife etched with a blessing against the dead.",
		Glyph: '/', Color: 15,
	},
	"multiaxe": {
		Name: "penitent's axe", Desc: "a heavy axe carried by the village penitents.",
		Glyph: '/', Color: 7,
	},
	"shock mace": {
		Name: "censer mace", Desc: "a mace whose head still smells of incense.",
		Glyph: '/', Color: 14,
	},
	"battle hammer": {
		Name: "war hammer", Desc: "a heavy war hammer from the village militia.",
		Glyph: '/', Color: 7,
	},
	"he grenade": {
		Name: "holy water flask", Desc: "a flask of blessed water that bursts into purifying steam.",
		Glyph: '!', Color: 11,
	},
	"incendiary grenade": {
		Name: "holy fire flask", Desc: "a flask of blessed oil that burns with a pale flame.",
		Glyph: '!', Color: 12,
	},
	"acid grenade": {
		Name: "alchemical acid flask", Desc: "a stoppered flask of biting alchemical acid.",
		Glyph: '!', Color: 10,
	},
	"smoke grenade": {
		Name: "incense flask", Desc: "a flask of heavy incense for driving off the restless dead.",
		Glyph: '!', Color: 8,
	},
	"flashbang grenade": {
		Name: "radiant flask", Desc: "a flask of captured daylight, painful to creatures of the dark.",
		Glyph: '!', Color: 14,
	},
	"first aid kit": {
		Name: "healer's kit", Desc: "a worn kit of needles, thread and clean cloth.",
		Glyph: '!', Color: 10,
	},
	"medkit": {
		Name: "surgeon's kit", Desc: "a heavy kit of proper medical tools.",
		Glyph: '!', Color: 10,
	},
}

func init() {
	for i := range ItemDefs {
		override, ok := itemThemeOverrides[ItemDefs[i].Name]
		if !ok {
			continue
		}
		ItemDefs[i].Name = override.Name
		ItemDefs[i].Desc = override.Desc
		ItemDefs[i].Glyph = override.Glyph
		ItemDefs[i].Color = override.Color
	}
}
