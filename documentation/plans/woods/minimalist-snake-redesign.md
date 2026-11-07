# Woods Minimalist Snake Redesign — Five Maps at 80×23

Status: **done**
Parent: [Woods plans](README.md)
Related: [Woods level revamp](level-revamp.md), [Woods colour palette](color-palette.md), [Cardinal movement](../systems/cardinal-movement.md)
Scope: woods run structure, map size, map themes and content density. Reuses the existing terrain, landmark, loot, faction and palette systems. Caves, temple and time-portal levels are postponed until the woods/village redesign is complete.

This plan replaces the current nine 80×46 woods maps with a compact five-map
snake run. Each map is a single themed place rather than a miniature world.
The five-map run replaces the current woods entirely; caves, temple and
time-portal work is postponed.

Implementation status: the five-map 80×23 run, themed generators, snake
directions, unlimited woods FOV, compact population/loot budgets, save data and
render/connectivity tests are implemented. The caves, temple and time-portal
section generator remains unchanged.

---

## 1. Goals

- Make the woods feel intentionally small and focused.
- Give every map one central element and one clear identity.
- Fit the whole game in an **80×25 terminal**.
- Keep the map at **80×23**, with no scrolling.
- Generate a random snake route through the woods once per run.
- Reduce content density so landmarks, roads, lamps and ponds have room to breathe.
- Keep the existing palette, cardinal movement, pathfinding, loot and factions.

## 2. Screen layout at 80×25

```text
row 0      latest message
rows 1–23  map view (80×23)
row 24     one-line HUD
```

Title, intro, class selection, equipment, options, help and message screens must
be reflowed to 25 rows. Long lists use pages instead of scrolling.

## 3. Run structure

The woods run is five maps:

| Map | Theme | Fixed / random | Entry / exit |
|---:|---|---|---|
| 0 | Church | fixed first | entry west, exit east |
| 1 | Peasant Houses | shuffled middle | opposite previous exit → random |
| 2 | Mill and Pond | shuffled middle | opposite previous exit → random |
| 3 | Forgotten Cemetery | shuffled middle | opposite previous exit → random |
| 4 | Meteor Field | fixed last | entry west, exit east to caves |

The three middle themes are shuffled once at world creation. The first map is
always the Church; the last map is always the Meteor Field and its exit leads to
the caves.

### Snake directions

Generate a direction sequence at run start:

- Map 0 is fixed: the Church is horizontal, with entry on the west edge and exit
  on the east edge.
- Each map has exactly one exit direction: north, south, east or west.
- The next map’s entry is the opposite edge of the previous exit.
- The next exit is chosen from the three edges that are not the entry edge.
- Immediate 180° reversals are therefore impossible.
- All other map edges are blocked by dense trees.
- The final map is fixed: the Meteor Field enters from the west and exits east
  to the caves.

Example:

```text
Church:          enter west  → exit east
Peasant Houses:  enter west  → exit north
Mill and Pond:   enter south → exit east
Cemetery:        enter west  → exit east
Meteor Field:    enter west  → exit east to caves
```

Store the theme order and direction sequence in the save so loading returns to
the same woods layout.

## 4. Compact map generator

Each 80×23 map uses the same skeleton:

1. **Tree border**
   - one- to two-tile dense tree border
   - only the entry and exit roads break through it

2. **Entry and exit**
   - entry road starts at the chosen edge
   - exit road leaves through the chosen edge
   - both are always connected by the road

3. **Central landmark**
   - one themed structure or feature in the middle third of the map
   - reserved before generic props

4. **Road**
   - cardinal path from entry to landmark to exit
   - lamps every 5–7 tiles
   - no random branches unless the theme needs one

5. **Optional water**
   - small pond or short river segment
   - bridge only where the road crosses water

6. **Props and population**
   - a small number of theme-appropriate props
   - reduced critter and loot counts

### Content budget

Per map:

- one central landmark
- 0–3 small outbuildings
- one road
- one optional pond/river
- 4–8 lamps
- 6–12 props
- 2–5 loot spots
- roughly 4–10 critters, including villagers/animals where appropriate

No map should contain both a large landmark and a generic hut village.

## 5. Map themes

### 5.1 Church

The first map and the anchor of the woods.

- central church with a small bell tower
- small graveyard beside it
- one pond or rain pool
- lamps along the road
- a shed or two at the edge
- mostly quiet; a few zombies around the graves

The road enters from the west and leaves through the generated exit direction.

### 5.2 Peasant Houses

- 3–4 cottages maximum
- one well or small square
- fences, gardens, haystacks and carts
- villagers, dogs and chickens or other small animals
- low-to-moderate zombie pressure

This is the social map, not a full town.

### 5.3 Mill and Pond

- water mill building
- pond or short river
- bridge and reeds
- fishing spots and river loot
- mill stone, sacks and a cart
- a few villagers and animals

The road should cross the water at the bridge.

### 5.4 Forgotten Cemetery

- fenced graveyard
- mausoleum or crypt entrance
- gravestones and dead trees
- fog and undead
- grave goods and armour loot

This is the main undead map before the meteor field.

### 5.5 Meteor Field

The final woods map.

- crater in the centre
- glowing meteor rocks
- dead and twisted trees
- strange glowing creatures
- burned ground and scattered debris
- exit leads to the caves

The meteor glow is the map’s dominant accent colour.

## 5.6 Visibility and light

- Player FOV has **no radius limit**; the 80×23 map edge is the visibility limit.
- The whole current map is visible once the player is on it.
- Lanterns, fires and the meteor glow are visible across the map.
- Warm light pools remain local around the light source; they are atmosphere, not
  a visibility restriction.
- No fog-of-war radius tuning is needed for the compact maps.

## 6. Reuse and replacement

Keep:

- woods terrain and props
- water, bridges, reeds and lily pads
- huts, furniture and landmarks
- loot categories and signature loot
- villagers, animals and zombie conversion
- cardinal movement and A* pathfinding
- true-colour palette and lighting

Replace:

- nine woods levels with five
- 80×46 maps with 80×23 maps
- random landmarks per level with one fixed theme per map
- independent exits with a run-wide snake direction sequence

## 7. Implementation phases

- [x] **Run data** — `WoodsRun` with theme order and direction sequence; save/load.
- [x] **Screen reflow** — 80×25 layouts for all screens.
- [x] **Map skeleton** — 80×23 border, entry/exit, road, lamps.
- [x] **Church map** — first complete themed map.
- [x] **Middle themes** — Peasant Houses, Mill and Pond, Cemetery.
- [x] **Meteor Field** — final map and caves transition.
- [x] **Density pass** — spawn counts, loot and props.
- [x] **Tests** — snake sequence, connectivity, dimensions, themes and save/load.

## 8. Acceptance criteria

- [x] The woods run contains exactly five maps.
- [x] The whole game fits an 80×25 terminal.
- [x] Woods maps are 80×23 and never scroll.
- [x] Each map has exactly one entry edge and one exit edge.
- [x] FOV covers the whole map; map edges are the only visibility limit.
- [x] All other edges are blocked by dense trees.
- [x] The snake direction sequence is generated once and saved.
- [x] No map immediately reverses direction.
- [x] Entry, central landmark and exit are connected by road.
- [x] The Church is always first, with entry west and exit east.
- [x] The Meteor Field is always last, enters west and exits east to the caves.
- [x] The three middle themes are shuffled.
- [x] Each map has its themed landmark and appropriate population.
- [x] No void tiles and all exits are reachable under cardinal movement.
- [x] Caves, temple and time-portal maps are unaffected.

## 9. Resolved decisions

- Fourth theme: **Forgotten Cemetery**.
- The Church is fixed first: horizontal layout, entry west, exit east.
- The Meteor Field is fixed last: entry west, exit east to the caves.
- The five-map run replaces the current nine-map woods entirely.
- Caves, temple and time-portal work is postponed.
- FOV is unlimited within the 80×23 map; map edges are the only limit.
- The three middle themes are shuffled once per run.
