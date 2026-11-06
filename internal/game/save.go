package game

import (
	"encoding/json"
	"fmt"
	"os"
)

type saveItem struct {
	Proto   int  `json:"p"`
	Charges int  `json:"c,omitempty"`
	Active  bool `json:"a,omitempty"`
	X       int  `json:"x,omitempty"`
	Y       int  `json:"y,omitempty"`
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
	Tactics   int            `json:"tac"`
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
	Target      int            `json:"target"`
}

func itemToSave(it *Item) saveItem {
	if it == nil {
		return saveItem{Proto: -1}
	}
	return saveItem{Proto: it.Proto, Charges: it.Charges, Active: it.Active, X: it.X, Y: it.Y}
}

func itemFromSave(s saveItem) *Item {
	if s.Proto < 0 {
		return nil
	}
	return &Item{Proto: s.Proto, Charges: s.Charges, Active: s.Active, X: s.X, Y: s.Y}
}

func critterToSave(c *Critter) saveCritter {
	s := saveCritter{
		Proto: c.Proto, Name: c.Name, X: c.X, Y: c.Y, HP: c.HP, MaxHP: c.MaxHP,
		STR: c.STR, TGH: c.TGH, AGI: c.AGI, INT: c.INT, SPD: c.SPD, Skills: c.Skills,
		Align: c.Align, Class: c.Class, AI: c.AI, Color: c.Color, Glyph: c.Glyph,
		Player: c.Player, Energy: c.Energy, Level: c.Level, XP: c.XP, Gold: c.Gold, Tactics: c.Tactics, Blood: c.Blood,
		Poison: c.Poison, LastSeenX: c.LastSeenX, LastSeenY: c.LastSeenY,
		Alert: c.Alert, Dead: c.Dead,
	}
	for _, it := range c.Inventory {
		s.Inv = append(s.Inv, itemToSave(it))
	}
	for _, it := range c.Equip {
		s.Equip = append(s.Equip, itemToSave(it))
	}
	return s
}

func critterFromSave(s saveCritter) *Critter {
	c := &Critter{
		Proto: s.Proto, Name: s.Name, X: s.X, Y: s.Y, HP: s.HP, MaxHP: s.MaxHP,
		STR: s.STR, TGH: s.TGH, AGI: s.AGI, INT: s.INT, SPD: s.SPD, Skills: s.Skills,
		Align: s.Align, Class: s.Class, AI: s.AI, Color: s.Color, Glyph: s.Glyph,
		Player: s.Player, Energy: s.Energy, Level: s.Level, XP: s.XP, Gold: s.Gold, Tactics: s.Tactics, Blood: s.Blood,
		Poison: s.Poison, LastSeenX: s.LastSeenX, LastSeenY: s.LastSeenY,
		Alert: s.Alert, Dead: s.Dead, Target: -1,
	}
	for _, si := range s.Inv {
		if it := itemFromSave(si); it != nil {
			c.Inventory = append(c.Inventory, it)
		}
	}
	for i := 0; i < NumEquipSlots; i++ {
		if i < len(s.Equip) {
			c.Equip[i] = itemFromSave(s.Equip[i])
		}
	}
	return c
}

func (g *Game) SaveGame(path string) error {
	s := saveGame{
		Version:     1,
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
		Target:      g.Target,
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
	g.Player = critterFromSave(s.Player)
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
	g.Target = s.Target
	g.Map = &Map{
		W: s.Map.W, H: s.Map.H, Tiles: s.Map.Tiles, Seen: s.Map.Seen, Visible: s.Map.Visible,
		Blood: s.Map.Blood, Fire: s.Map.Fire, Acid: s.Map.Acid, Smoke: s.Map.Smoke,
		Theme: s.Map.Theme, Name: s.Map.Name,
	}
	g.Critters = nil
	for _, sc := range s.Critters {
		g.Critters = append(g.Critters, critterFromSave(sc))
	}
	g.Items = nil
	for _, si := range s.Items {
		if it := itemFromSave(si); it != nil {
			g.Items = append(g.Items, it)
		}
	}
	g.Messages = nil
	for _, m := range s.Messages {
		g.Messages = append(g.Messages, Message{Text: m.Text, Color: m.Color, Turn: m.Turn})
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
