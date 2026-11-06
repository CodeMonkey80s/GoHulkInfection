package game

import (
	"encoding/json"
	"fmt"
	"os"
)

const saveVersion = 3

type saveItem struct {
	Proto   int  `json:"p"`
	Charges int  `json:"c,omitempty"`
	AP      int  `json:"ap,omitempty"`
	APInit  bool `json:"api,omitempty"`
	Active  bool `json:"a,omitempty"`
	X       int  `json:"x,omitempty"`
	Y       int  `json:"y,omitempty"`
}

type saveSlots struct {
	Hands        saveItem   `json:"hands"`
	Body         saveItem   `json:"body"`
	Belt         []saveItem `json:"belt,omitempty"`
	BeltUnlocked int        `json:"belt_unlocked"`
}

type saveCritter struct {
	Proto     int            `json:"p"`
	Name      string         `json:"n"`
	X         int            `json:"x"`
	Y         int            `json:"y"`
	HP        int            `json:"hp"`
	MaxHP     int            `json:"mhp"`
	STR       int            `json:"str"`
	TGH       int            `json:"tgh"`
	AGI       int            `json:"agi"`
	INT       int            `json:"int"`
	SPD       int            `json:"spd"`
	Skills    [NumSkills]int `json:"sk"`
	Align     int            `json:"al"`
	Class     int            `json:"cl"`
	AI        int            `json:"ai"`
	Color     uint8          `json:"co"`
	Glyph     byte           `json:"gl"`
	Player    bool           `json:"pl"`
	Energy    int            `json:"en"`
	Level     int            `json:"lvl"`
	XP        int            `json:"xp"`
	Gold      int            `json:"gold"`
	Inv       []saveItem     `json:"inv"`
	Equip     []saveItem     `json:"eq"`
	Slots     *saveSlots     `json:"slots,omitempty"`
	Blood     int            `json:"bl"`
	Poison    int            `json:"po"`
	LastSeenX int            `json:"lsx"`
	LastSeenY int            `json:"lsy"`
	Alert     int            `json:"alert"`
	Dead      bool           `json:"dead"`
}

type saveMessage struct {
	Text  string `json:"t"`
	Color uint8  `json:"c"`
	Turn  int    `json:"turn"`
}

type saveMap struct {
	W       int     `json:"w"`
	H       int     `json:"h"`
	Tiles   []uint8 `json:"t"`
	Seen    []bool  `json:"seen"`
	Visible []bool  `json:"vis"`
	Blood   []uint8 `json:"blood"`
	Fire    []int   `json:"fire"`
	Acid    []int   `json:"acid"`
	Smoke   []int   `json:"smoke"`
	Theme   int     `json:"theme"`
	Name    string  `json:"name"`
}

type saveGame struct {
	Version     int            `json:"version"`
	Map         saveMap        `json:"map"`
	Player      saveCritter    `json:"player"`
	PlayerName  string         `json:"player_name"`
	PlayerClass int            `json:"player_class"`
	Critters    []saveCritter  `json:"critters"`
	Items       []saveItem     `json:"items"`
	Messages    []saveMessage  `json:"messages"`
	Turn        int            `json:"turn"`
	Level       int            `json:"level"`
	Theme       int            `json:"theme"`
	Weather     int            `json:"weather"`
	WeatherLeft int            `json:"weather_left"`
	Options     map[string]int `json:"options"`
	HintsShown  int            `json:"hints"`
	FreePoints  int            `json:"free_points"`
	Score       int            `json:"score"`
}

func itemToSave(it *Item) saveItem {
	if it == nil {
		return saveItem{Proto: -1}
	}
	return saveItem{Proto: it.Proto, Charges: it.Charges, AP: it.AP, APInit: it.APInit, Active: it.Active, X: it.X, Y: it.Y}
}

func itemFromSave(s saveItem) *Item {
	if s.Proto < 0 {
		return nil
	}
	return &Item{Proto: s.Proto, Charges: s.Charges, AP: s.AP, APInit: s.APInit, Active: s.Active, X: s.X, Y: s.Y}
}

func critterToSave(c *Critter) saveCritter {
	s := saveCritter{
		Proto: c.Proto, Name: c.Name, X: c.X, Y: c.Y, HP: c.HP, MaxHP: c.MaxHP,
		STR: c.STR, TGH: c.TGH, AGI: c.AGI, INT: c.INT, SPD: c.SPD, Skills: c.Skills,
		Align: c.Align, Class: c.Class, AI: c.AI, Color: c.Color, Glyph: c.Glyph,
		Player: c.Player, Energy: c.Energy, Level: c.Level, XP: c.XP, Gold: c.Gold, Blood: c.Blood,
		Poison: c.Poison, LastSeenX: c.LastSeenX, LastSeenY: c.LastSeenY,
		Alert: c.Alert, Dead: c.Dead,
	}
	if c.Player {
		ss := &saveSlots{
			Hands:        itemToSave(c.Slots.Hands),
			Body:         itemToSave(c.Slots.Body),
			BeltUnlocked: c.Slots.BeltUnlocked,
		}
		for _, it := range c.Slots.Belt {
			ss.Belt = append(ss.Belt, itemToSave(it))
		}
		s.Slots = ss
	} else {
		for _, it := range c.Inventory {
			s.Inv = append(s.Inv, itemToSave(it))
		}
		for _, it := range c.Equip {
			s.Equip = append(s.Equip, itemToSave(it))
		}
	}
	return s
}

func critterFromSave(s saveCritter, version int) (*Critter, int) {
	removed := 0
	c := &Critter{
		Proto: s.Proto, Name: s.Name, X: s.X, Y: s.Y, HP: s.HP, MaxHP: s.MaxHP,
		STR: s.STR, TGH: s.TGH, AGI: s.AGI, INT: s.INT, SPD: s.SPD, Skills: s.Skills,
		Align: s.Align, Class: s.Class, AI: s.AI, Color: s.Color, Glyph: s.Glyph,
		Player: s.Player, Energy: s.Energy, Level: s.Level, XP: s.XP, Gold: s.Gold, Blood: s.Blood,
		Poison: s.Poison, LastSeenX: s.LastSeenX, LastSeenY: s.LastSeenY,
		Alert: s.Alert, Dead: s.Dead,
	}
	if s.Player {
		c.Slots = PlayerSlots{BeltUnlocked: 1}
		if version >= 3 && s.Slots != nil {
			if it := itemFromSave(s.Slots.Hands); it != nil {
				if ItemAvailable(it.Proto) {
					c.Slots.Hands = it
				} else {
					removed++
				}
			}
			if it := itemFromSave(s.Slots.Body); it != nil {
				if ItemAvailable(it.Proto) {
					c.Slots.Body = it
				} else {
					removed++
				}
			}
			for i := 0; i < 3 && i < len(s.Slots.Belt); i++ {
				it := itemFromSave(s.Slots.Belt[i])
				if it == nil {
					continue
				}
				if !ItemAvailable(it.Proto) {
					removed++
					continue
				}
				c.Slots.Belt[i] = it
			}
			c.Slots.BeltUnlocked = s.Slots.BeltUnlocked
		} else {
			// Migrate the old inventory/equipment layout into player slots.
			oldHand := -1
			oldBody := -1
			oldThrow := -1
			oldTool := -1
			if version < 2 {
				oldHand, oldBody, oldThrow, oldTool = 0, 2, 3, 4
			} else {
				oldHand, oldBody, oldThrow, oldTool = 0, 1, 2, 3
			}
			if oldHand >= 0 && oldHand < len(s.Equip) {
				if it := itemFromSave(s.Equip[oldHand]); it != nil && ItemAvailable(it.Proto) {
					c.Slots.Hands = it
				}
			}
			if oldBody >= 0 && oldBody < len(s.Equip) {
				if it := itemFromSave(s.Equip[oldBody]); it != nil && ItemAvailable(it.Proto) {
					c.Slots.Body = it
				}
			}
			c.Slots.BeltUnlocked = 3
			var candidates []saveItem
			for _, idx := range []int{oldThrow, oldTool} {
				if idx >= 0 && idx < len(s.Equip) {
					candidates = append(candidates, s.Equip[idx])
				}
			}
			candidates = append(candidates, s.Inv...)
			bi := 0
			for _, si := range candidates {
				if si.Proto < 0 {
					continue
				}
				if !ItemAvailable(si.Proto) {
					removed++
					continue
				}
				if bi >= 3 {
					removed++
					continue
				}
				if it := itemFromSave(si); it != nil {
					c.Slots.Belt[bi] = it
					bi++
				}
			}
		}
		ensureArmorAP(c.Slots.Body)
		if c.Slots.BeltUnlocked < 1 {
			c.Slots.BeltUnlocked = 1
		}
		if c.Slots.BeltUnlocked > 3 {
			c.Slots.BeltUnlocked = 3
		}
		return c, removed
	}
	for _, si := range s.Inv {
		if si.Proto < 0 {
			continue
		}
		if !ItemAvailable(si.Proto) {
			removed++
			continue
		}
		if it := itemFromSave(si); it != nil {
			c.Inventory = append(c.Inventory, it)
		}
	}
	if version < 2 {
		// Old layout: melee, ranged, armor, grenade, tool.
		if len(s.Equip) > 1 && s.Equip[1].Proto >= 0 && !ItemAvailable(s.Equip[1].Proto) {
			removed++
		}
		want := [NumEquipSlots]int{0, 2, 3, 4}
		for i, oldIndex := range want {
			if oldIndex < len(s.Equip) {
				si := s.Equip[oldIndex]
				if si.Proto < 0 {
					continue
				}
				if !ItemAvailable(si.Proto) {
					removed++
					continue
				}
				c.Equip[i] = itemFromSave(si)
			}
		}
	} else {
		for i := 0; i < NumEquipSlots; i++ {
			if i < len(s.Equip) {
				si := s.Equip[i]
				if si.Proto < 0 {
					continue
				}
				if !ItemAvailable(si.Proto) {
					removed++
					continue
				}
				c.Equip[i] = itemFromSave(si)
			}
		}
	}
	return c, removed
}

func (g *Game) SaveGame(path string) error {
	s := saveGame{
		Version:     saveVersion,
		Player:      critterToSave(g.Player),
		PlayerName:  g.PlayerName,
		PlayerClass: g.PlayerClass,
		Turn:        g.Turn,
		Level:       g.Level,
		Theme:       g.Theme,
		Weather:     int(g.Weather),
		WeatherLeft: g.WeatherLeft,
		Options:     map[string]int(g.Options),
		HintsShown:  g.HintsShown,
		FreePoints:  g.FreePoints,
		Score:       g.Score,
	}
	s.Map = saveMap{
		W: g.Map.W, H: g.Map.H, Tiles: g.Map.Tiles, Seen: g.Map.Seen, Visible: g.Map.Visible,
		Blood: g.Map.Blood, Fire: g.Map.Fire, Acid: g.Map.Acid, Smoke: g.Map.Smoke,
		Theme: g.Map.Theme, Name: g.Map.Name,
	}
	for _, c := range g.Critters {
		s.Critters = append(s.Critters, critterToSave(c))
	}
	for _, it := range g.Items {
		s.Items = append(s.Items, itemToSave(it))
	}
	for _, m := range g.Messages {
		s.Messages = append(s.Messages, saveMessage{Text: m.Text, Color: m.Color, Turn: m.Turn})
	}
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (g *Game) LoadGame(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var s saveGame
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	removed := 0
	g.Player, removed = critterFromSave(s.Player, s.Version)
	g.PlayerName = s.PlayerName
	g.PlayerClass = s.PlayerClass
	g.Turn = s.Turn
	g.Level = s.Level
	g.Theme = s.Theme
	g.Weather = Weather(s.Weather)
	g.WeatherLeft = s.WeatherLeft
	g.Options = Options(s.Options)
	if g.Options == nil {
		g.Options = NewOptions()
	}
	g.HintsShown = s.HintsShown
	g.FreePoints = s.FreePoints
	g.Score = s.Score
	g.Map = &Map{
		W: s.Map.W, H: s.Map.H, Tiles: s.Map.Tiles, Seen: s.Map.Seen, Visible: s.Map.Visible,
		Blood: s.Map.Blood, Fire: s.Map.Fire, Acid: s.Map.Acid, Smoke: s.Map.Smoke,
		Theme: s.Map.Theme, Name: s.Map.Name,
	}
	g.Critters = nil
	for _, sc := range s.Critters {
		c, n := critterFromSave(sc, s.Version)
		removed += n
		g.Critters = append(g.Critters, c)
	}
	g.Items = nil
	occupied := map[int]bool{}
	for _, si := range s.Items {
		if !ItemAvailable(si.Proto) {
			removed++
			continue
		}
		it := itemFromSave(si)
		if it == nil {
			continue
		}
		i := g.Map.idx(it.X, it.Y)
		if occupied[i] {
			removed++
			continue
		}
		occupied[i] = true
		g.Items = append(g.Items, it)
	}
	g.Messages = nil
	for _, m := range s.Messages {
		g.Messages = append(g.Messages, Message{Text: m.Text, Color: m.Color, Turn: m.Turn})
	}
	if removed > 0 {
		g.Msg(fmt.Sprintf("%d obsolete item(s) removed by the melee-only ruleset.", removed), 8)
	}
	g.State = StatePlaying
	g.UpdateFOV()
	return nil
}

func (g *Game) SaveAndQuit() {
	g.Msg("saving game...", 10)
	if err := g.SaveGame(g.SavePath); err != nil {
		g.Msg(fmt.Sprintf("could not save game: %v", err), 12)
		return
	}
	g.Msg("game saved.", 10)
	g.Quit = true
}
