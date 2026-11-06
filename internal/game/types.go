package game

import (
	"math/rand"
)

type ItemKind int

const (
	KindBodyPart ItemKind = iota
	KindArmor
	KindMelee
	KindRanged
	KindAmmo
	KindGrenade
	KindMedical
)

const (
	SkillPerception = iota
	SkillDodging
	SkillUnarmed
	SkillMelee
	SkillRanged
	SkillThrowing
	SkillSearching
	SkillMetabolism
	SkillFocusing
	NumSkills
)

const (
	AlignPlayer = iota
	AlignEnemy
	AlignNeutral
)

const (
	AISimple = iota // zombie-style: chase, attack, lose interest
	AIObserver
	AIAnimal
	AIPlayer
)

const (
	EquipMelee = iota
	EquipArmor
	EquipThrowable
	EquipTool
	NumEquipSlots
)

const (
	SlotHands = iota
	SlotBody
	SlotBelt1
	SlotBelt2
	SlotBelt3
	NumPlayerSlots
)

// PlayerSlots is the complete player inventory. Monsters keep using the
// legacy Inventory/Equip fields.
type PlayerSlots struct {
	Hands        *Item
	Body         *Item
	Belt         [3]*Item
	BeltUnlocked int
}

type KeyBind struct {
	Key    int
	Mod    int
	Action int
	Name   string
	Desc   string
}

type OptionDef struct {
	Name            string
	Desc            string
	Value, Min, Max int
}

type ClassDef struct {
	Name, Desc         string
	STR, TGH, AGI, INT int
	Skills             [NumSkills]int
}

// Item is an instance of an item prototype.
type Item struct {
	Proto   int
	Charges int
	Active  bool
	X, Y    int // -1,-1 when carried
}

func (it *Item) Def() *ItemDef {
	if it == nil || it.Proto < 0 || it.Proto >= len(ItemDefs) {
		return nil
	}
	return &ItemDefs[it.Proto]
}

func (d *ItemDef) Weapon() *WeaponDef {
	if d == nil || d.WeaponSub < 0 || d.WeaponSub >= len(WeaponDefs) {
		return nil
	}
	return &WeaponDefs[d.WeaponSub]
}

func (it *Item) Name() string {
	if d := it.Def(); d != nil {
		return d.Name
	}
	return "unknown"
}

// Critter is an instance of a critter prototype.
type Critter struct {
	Proto     int
	Name      string
	X, Y      int
	HP        int
	MaxHP     int
	STR       int
	TGH       int
	AGI       int
	INT       int
	SPD       int
	Skills    [NumSkills]int
	Align     int
	Class     int
	AI        int
	Color     uint8
	Glyph     byte
	Player    bool
	Energy    int
	DmgMin    int
	DmgMax    int
	Level     int
	XP        int
	Gold      int
	Inventory []*Item
	Equip     [NumEquipSlots]*Item
	Slots     PlayerSlots
	Blood     int
	Poison    int
	LastSeenX int
	LastSeenY int
	Alert     int
	Dead      bool
	IsAnimal  bool
}

func (c *Critter) Def() *CritterDef {
	if c == nil || c.Proto < 0 || c.Proto >= len(CritterDefs) {
		return nil
	}
	return &CritterDefs[c.Proto]
}

func (c *Critter) NameOrUnknown() string {
	if c == nil {
		return "unknown"
	}
	if c.Name != "" {
		return c.Name
	}
	if d := c.Def(); d != nil {
		return d.Name
	}
	return "unknown"
}

func (c *Critter) Skill(s int) int {
	if s < 0 || s >= NumSkills {
		return 0
	}
	return c.Skills[s]
}

func (c *Critter) Weapon() *Item {
	return c.HandItem()
}

func (c *Critter) Throwable() *Item {
	return c.ThrowableItem()
}

func (c *Critter) Armor() *Item {
	return c.BodyItem()
}

func (c *Critter) IsAlive() bool {
	return c != nil && !c.Dead && c.HP > 0
}

func (c *Critter) HealthState() string {
	r := 1.0
	if c.MaxHP > 0 {
		r = float64(c.HP) / float64(c.MaxHP)
	}
	switch {
	case r <= 0:
		return "dying"
	case r < 0.1:
		return "almost dead"
	case r < 0.25:
		return "heavly wounded"
	case r < 0.4:
		return "badly wounded"
	case r < 0.55:
		return "moderately wounded"
	case r < 0.7:
		return "lightly wounded"
	case r < 0.9:
		return "slightly wounded"
	default:
		return "healthy"
	}
}

type Message struct {
	Text  string
	Color uint8
	Turn  int
}

type Options map[string]int

func NewOptions() Options {
	o := Options{}
	defaults := map[string]int{
		"optDrawDebugAI": 0, "optDrawDebugInfo": 0, "optDrawGameMechanicInfo": 0,
		"optGameChangeLevel": 0, "optGameFirstLevelCreated": 1, "optGameHintNumber": 0,
		"optGameHints": 1, "optGameLevelID": 0, "optGameSetOptionsForRelease": 1,
		"optGameTheme": 0, "optGameUseAI": 1, "optGameSendDamageToItems": 0,
		"optGameUseDamageDescriptions": 1, "optGameSpriteAnimationSpeed": 50,
		"optInvColorOutToHeavyItems": 1, "optInvQuitAfterItemUse": 1,
		"optMapCreateBlood": 1, "optMapCreateCorpses": 1, "optMapCreateLimbs": 1,
		"optMapDefaultExplored": 0, "optMapDefaultVisible": 1,
		"optMapDrawHiddenTiles": 0, "optMapDrawFullMapName": 1, "optMapDrawPortals": 1,
		"optMapHighlightSpottedCritters": 1, "optMapResetFOVAfterAction": 0,
		"optMapResetCursorPosition": 0, "optMapShadeColor": 1, "optMapUseMemoryMap": 1,
		"optMechanicAcidEatItemChance": 30, "optMechanicBaseHitChance": 50,
		"optMechanicAVMod1": 100, "optMechanicAVMod2": 200,
		"optMechanicDVMod1": 100, "optMechanicDVMod2": 200,
		"optMechanicAVModEnemyIsNear": 75,
		"optPlayerAlwaysSpotCritters": 0, "optPlayerAttackOnMove": 1,
		"optPlayerAutoOpenDoor": 1, "optPlayerDisplayItemInfoInStatus": 1,
		"optPlayerDisplaySkillValues": 1, "optPlayerFreePointsPerLevel": 1,
		"optPlayerGhostMode": 0, "optPlayerGodMode": 0, "optPlayerShowAVDV": 1,
		"optPlayerWins": 0,
	}
	for i := range OptionDefs {
		if OptionDefs[i].Name == "optGameHintNumber" {
			OptionDefs[i].Max = 8
		}
	}
	for _, d := range OptionDefs {
		if v, ok := defaults[d.Name]; ok {
			o[d.Name] = v
		} else {
			o[d.Name] = d.Min
		}
	}
	return o
}

func (o Options) Get(name string) int { return o[name] }

func (o Options) Bool(name string) bool { return o[name] != 0 }

type Weather int

const (
	WeatherClear Weather = iota
	WeatherRain
	WeatherAcidRain
)

type GameState int

const (
	StateTitle GameState = iota
	StateIntro
	StateCharCreate
	StatePlaying
	StateInventory
	StatePlayerScreen
	StateSkills
	StateOptions
	StateHelp
	StateMessages
	StateLook
	StateThrow
	StateQuitConfirm
	StateGameOver
	StateVictory
)

type Game struct {
	Map          *Map
	Player       *Critter
	PlayerName   string
	PlayerClass  int
	Critters     []*Critter
	Items        []*Item
	Messages     []Message
	Turn         int
	Level        int
	Theme        int
	Weather      Weather
	WeatherLeft  int
	Options      Options
	State        GameState
	CursorX      int
	CursorY      int
	LookX, LookY int
	SavePath     string
	RNG          *rand.Rand
	HintsShown   int
	FreePoints   int
	Running      bool
	Quit         bool
	Score        int
	StartTime    int64
	// temporary screen data
	ScrollMsg    int
	PendingThrow int
}

func NewGame(seed int64) *Game {
	g := &Game{
		RNG:      rand.New(rand.NewSource(seed)),
		Options:  NewOptions(),
		SavePath: "hi_save.json",
		Running:  true,
	}
	g.Player = &Critter{
		Proto:  -1,
		Name:   "you",
		Player: true,
		Align:  AlignPlayer,
		Class:  4,
		X:      -1,
		Y:      -1,
		Gold:   0,
	}
	return g
}

func (g *Game) Msg(text string, color uint8) {
	if text == "" {
		return
	}
	g.Messages = append(g.Messages, Message{Text: text, Color: color, Turn: g.Turn})
	if len(g.Messages) > 500 {
		g.Messages = g.Messages[len(g.Messages)-500:]
	}
}

func (g *Game) LastMessages(n int) []Message {
	if n <= 0 || len(g.Messages) == 0 {
		return nil
	}
	if n > len(g.Messages) {
		n = len(g.Messages)
	}
	return g.Messages[len(g.Messages)-n:]
}

func (g *Game) Roll(n int) int {
	if n <= 0 {
		return 0
	}
	return g.RNG.Intn(n)
}

func (g *Game) RollRange(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + g.RNG.Intn(hi-lo+1)
}

func (g *Game) Chance(percent int) bool {
	if percent <= 0 {
		return false
	}
	if percent >= 100 {
		return true
	}
	return g.RNG.Intn(100) < percent
}

func (g *Game) CritterAt(x, y int) *Critter {
	for _, c := range g.Critters {
		if c.IsAlive() && c.X == x && c.Y == y {
			return c
		}
	}
	return nil
}

func (g *Game) ItemAt(x, y int) *Item {
	for _, it := range g.Items {
		if it.X == x && it.Y == y {
			return it
		}
	}
	return nil
}

func (g *Game) ItemsAt(x, y int) []*Item {
	var out []*Item
	for _, it := range g.Items {
		if it.X == x && it.Y == y {
			out = append(out, it)
		}
	}
	return out
}

func (g *Game) RemoveItem(it *Item) {
	for i, x := range g.Items {
		if x == it {
			g.Items = append(g.Items[:i], g.Items[i+1:]...)
			return
		}
	}
}

func (g *Game) RemoveCritter(c *Critter) {
	for i, x := range g.Critters {
		if x == c {
			g.Critters = append(g.Critters[:i], g.Critters[i+1:]...)
			return
		}
	}
}
