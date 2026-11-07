# Woods Level Revamp — Levels 0–8

Status: **done**
Parent: [Woods plans](README.md)  
Related: [Woods colour palette](color-palette.md)  
Scope: **woods levels 0–8 only**

Superseded in part by [Woods Minimalist Snake Redesign](minimalist-snake-redesign.md):
the proposed five-map 80×23 run replaces the map count, size, progression and
per-level theme structure described here. The terrain, loot, faction and
rendering systems remain reusable.

Implementation status: the first outdoor terrain pass is implemented for woods
levels 0–8: grass, forest clusters, clearings, bushes, brambles, a dirt road,
lantern posts and scattered props. The first hut pass is also implemented:
road-facing huts with doors, windows and furniture. Rivers/ponds, reeds, lily
pads and road bridges are implemented as the first water pass. Each wood level
now has its named identity and a unique landmark. Wood-specific food, tools,
weapons, armour, potions and relics are appended, with a shuffled signature
loot category assigned to each level and guaranteed near its landmark. Villagers,
neutral animals, zombie target selection and villager-to-zombie conversion are
implemented. Exit placement now uses edge, centre and landmark candidates;
connectivity, hut-door/road, bridge, landmark, loot, faction and rendering
checks are covered by generator tests. Caves, temple and time-portal levels
still use the section generator.

This document describes the next major map-content pass for the Go port. The
goal is to stop treating every level as a reused room-section grid and give
each theme — starting with the woods — its own terrain language, landmarks,
loot identity and inhabitants.

---

## 1. Current baseline

The current Go port has:

- four level themes: `woods`, `caves`, `temple of underminer`, `time portal`
- 80×46 maps sized to fit the viewport; woods levels 0–8 use the outdoor
  terrain generator, other themes use the recovered 10×40 cube-section templates
- 64 recovered tile prototypes plus appended woods terrain tiles
- 49 recovered item prototypes and a 20-entry weapon table
- 17 recovered critter prototypes
- simple AI that generally treats the player as the only target
- theme terrain fill and reachable-portal validation

Relevant code:

- `internal/game/map.go` — `GenerateLevel`, terrain, exits, spawn population
- `internal/game/woods_gen.go` — woods terrain, forest noise, roads and lanterns
- `internal/game/woods_tiles.go` — appended woods tile definitions
- `internal/game/ai.go` — critter turns and target selection
- `internal/game/combat.go` — damage, death, corpses, explosions
- `internal/game/data_gen.go` — recovered tile/item/critter/weapon tables
- `internal/game/types.go` — critter, item, faction and option types

The woods revamp should replace the section-grid generator **only for levels
0–8**. Caves, temple and time-portal levels keep their current generator until
their own revamp plans are written.

---

## 2. Revamp status

| Theme | Levels | Status | Notes |
|---|---:|---|---|
| Woods | 0–8 | **Planned in this document** | Full outdoor generator, landmarks, new props, loot and factions |
| Caves | 9–17 | Not planned yet | Keep current section generator |
| Temple of the Underminer | 18–26 | Not planned yet | Keep current section generator |
| Time portal / lab | 27 | Not planned yet | Keep current generator |

---

## 3. Woods revamp (levels 0–8)

### 3.1 Goals

- Make the woods feel like an outdoor place, not a sequence of room sections.
- Drastically reduce random interior-style walls and doors.
- Use huts, fences, cemeteries, rivers, bridges, roads and vegetation for shape.
- Give every wood level a distinct identity and one unique landmark.
- Add wood-specific weapons, armour, tools and food.
- Give each wood level a random signature loot category.
- Add non-hostile people and varied animal life.
- Make zombies local threats with small awareness ranges and villager conversion.
- Randomize exits between map edges, the middle and landmark-related positions.

### 3.2 Generation pipeline

Replace the shared `2 × 8` section assembly for woods with this layered pipeline:

1. **Base terrain**
   - grass
   - forest-density noise
   - natural clearings
   - occasional stumps, flowers and fallen logs

2. **Water**
   - optional pond or edge-to-edge river
   - shallow banks, reeds, lily pads
   - deep-water sections impassable or risky

3. **Road network**
   - main dirt road from entry toward the exit
   - branches to huts and the unique landmark
   - trees, brambles and walls cleared from the road
   - lantern posts along roads every 6–8 tiles

4. **Landmark zone**
   - reserve one clear area for the level’s unique place
   - place its structure/props before generic decoration

5. **Huts**
   - 2–6 small wooden buildings per level
   - one road-facing door, one or two windows
   - wooden floor, furniture and small loot
   - no generic standalone walls outside buildings

6. **Props and vegetation**
   - bushes, brambles, fences, wells, gravestones, statues, lanterns
   - pillars only in cemeteries, ruins, stone circles or gateways

7. **Exit and entry**
   - select the exit using the candidate rules in §3.7
   - connect it to the road network

8. **Population and loot**
   - villagers, animals, zombies and hostile creatures
   - level signature loot and ordinary woodland loot

9. **Validation**
   - no void tiles
   - exit reachable
   - hut doors reachable from a road
   - bridge present wherever a road crosses water
   - landmark and signature loot present

### 3.3 New tiles and props

Append new tile definitions at the end of `TileDefs` so all recovered tile IDs
remain stable.

Suggested additions:

- **Ground:** dirt road, mud, ford, river bank, wooden bridge, stone bridge
- **Vegetation:** bush, bramble, reeds, flowers, mushroom ring, dead tree,
  orchard tree, stump, fallen log
- **Structures:** wooden fence, gate, well, lantern post, gravestone, grave,
  crypt wall, mausoleum door
- **Furniture:** table, chair, bed, crate, barrel, fireplace, shelf, altar,
  kiln, tent, cart, haystack, market stall
- **Landmark-specific:** chapel floor/wall/altar, bell tower, mill stone,
  stone circle, druid altar

Bush behaviour:

- **bush** — passable, blocks sight, can hide items, searchable
- **bramble** — blocks movement and sight, bashable, drops berries or wood
- some bushes may hide animals or small loot

Lanterns:

- placed beside roads
- a warm light marker and safe-road cue
- optional later extension: emit a small static FOV light radius

### 3.4 Huts and structures

Huts replace most interior walls/doors.

- size roughly 5×5 to 8×6
- wooden wall tiles, wooden floor
- exactly one door, facing the nearest road
- one or two windows
- furniture as blocking tiles where appropriate
- fireplace, table, chair, bed, crates, barrels, shelves
- loot: food, tools, clothes, occasionally a weapon
- villagers may spawn inside or near huts
- doors can be opened, bashed or lockpicked if a lock is added

Pillars are not generic decoration. They are allowed only in:

- cemeteries and mausoleums
- chapel ruins
- druid stone circles
- gateways or special ruins

### 3.5 Roads, lanterns, water and bridges

- Roads are generated with a cost-based pathfinder (A* or weighted BFS).
- Main road connects entry and exit.
- Side roads connect huts and the unique landmark.
- A two-tile-wide road is acceptable for the main route; side paths can be one tile.
- Trees and brambles are removed from the road.
- Lantern posts are placed every 6–8 tiles along the main road.
- A river is generated as a meandering edge-to-edge band.
- Bridges are placed at road/river intersections.
- A pond can replace the river in levels that need a calmer layout.
- Exit placement must account for water: if the exit is across the river, a bridge is mandatory.

### 3.6 Unique identities for woods levels 0–8

Each level gets one reserved landmark and a distinct feel.

| Level | Name | Unique place | Identity |
|---:|---|---|---|
| 0 | The Woodcutter’s Clearing | starting shed, woodpile, stumps | tutorial/quiet; few zombies; neutral animals |
| 1 | The Old Orchard | rows of fruit trees, beehives, cider press | edibles, wasps, deer and rabbits |
| 2 | The Forgotten Cemetery | fenced graveyard, mausoleum, dead trees | undead, grave goods, fog |
| 3 | The River Mill | river, bridge, mill house, mill pond | water, reeds, fishing, bridge crossing |
| 4 | The Poacher’s Camp | tents, campfire, cages, drying racks | hunters, traps, dogs/wolves, ranged loot |
| 5 | The Druid’s Grove | stone circle, giant trees, mushroom ring, altar | pillars only here; neutral beasts; herbs |
| 6 | The Abandoned Village | 5–8 huts, well, market square, carts | survivors, zombie attack, militia loot |
| 7 | The Charcoal Burner’s Clearing | burned forest, kilns, smoke/fire, charred huts | fire hazards, ash, torches and lanterns |
| 8 | The Ruined Chapel | chapel, bell tower, graveyard, crypt entrance | strongest undead, relics, entrance to caves |

Level 8’s exit is the chapel crypt and leads to the caves.

Other details:

- level 0 should remain the safest introduction
- levels 2 and 8 are the main undead levels
- levels 3 and 6 should have the most villagers
- level 4 should introduce hostile dogs/wolves and traps
- level 7 should use fire/smoke hazards carefully so paths remain usable
- level 5 may have the only stone-circle pillars in the woods

### 3.7 Exit placement

Current exits are placed on the farthest reachable tile. Replace that with a
candidate-based selection:

- **50%** random map edge zone (north/east/south/west)
- **25%** map center or central clearing
- **25%** beyond or inside the unique landmark

Rules:

- never inside a wall, hut, water or unreachable pocket
- always connected by a road or clear path
- on river levels, the exit may be across the water and requires a bridge
- level 8 exit is fixed at the chapel crypt
- exit is visually distinct: stairs, cave mouth, gate or portal
- do not spawn the player on the exit or within its immediate approach

### 3.8 Items, tools, food and signature loot

Add new item kinds and definitions:

- **Food:** apple, pear, berries, honey, bread, cheese, dried meat, mushroom
  (safe/toxic), raw/cooked fish, stew
- **Tools:** hatchet, woodcutter’s axe, saw, shovel, rope, lantern, torch,
  lockpick, fishing rod, whetstone
- **Weapons:** hatchet, wood axe, pitchfork, sickle, scythe, club, spear,
  hunting bow, crossbow, sling, old rifle, throwing knives
- **Armour:** leather jacket, hide vest, padded coat, chain shirt, ranger cloak

Food effects:

- restore health and/or regeneration
- some mushrooms are poisonous
- cooked food is safer and heals more than raw food

Tool effects (first pass):

- axe: better bashing of trees/brambles and wooden doors
- saw: faster work on wooden obstacles
- shovel: dig graves and soft ground
- rope: future climbing/descending hook
- lantern/torch: light source and fire starter
- lockpick: open locked hut/chapel doors
- fishing rod: catch fish at ponds/rivers
- whetstone: repair or improve weapon condition

Signature loot:

At world creation, shuffle **nine signature item categories** and assign one
unique category to each wood level. No category repeats across levels 0–8.
The level’s unique landmark guarantees at least one item from its category, and
ordinary loot on that level is biased toward it. The specific item within the
category is also randomized.

Proposed category mapping:

| Level | Signature category |
|---:|---|
| 0 | tools |
| 1 | edibles |
| 2 | grave goods / armour |
| 3 | fishing / river gear |
| 4 | hunting / ranged weapons |
| 5 | herbs / potions |
| 6 | militia melee weapons |
| 7 | fire / light tools |
| 8 | chapel relics |

Keep existing item IDs unchanged. Append new item and weapon definitions so old
saves either continue to work or can be migrated with a save-version bump.

### 3.9 People, animals, factions and zombies

Introduce a faction/target layer:

- player
- hostile undead
- neutral animals
- villagers/people

Behaviour:

**Villagers**

- wander roads, clearings and huts
- flee from zombies
- fight back if armed or cornered
- may give a hint or drop supplies when rescued
- may be killed and converted

**Animals**

- rabbits, deer, foxes: flee when approached or attacked
- boars: neutral until attacked, then charge
- wolves: hostile
- crows: flee
- frogs and pond life: ambient/non-combat
- not every animal should be hostile; most should be free-roaming scenery or food

**Zombies**

- only react when a living target is very close
- idle zombies usually stay still and occasionally take a random step
- small aggro radius (initial proposal: 2–3 tiles, or sight within 3)
- chase for a limited number of turns, then lose interest
- attack villagers, animals and the player
- if a zombie kills a villager, the villager rises as a zombie after 1–3 turns
- the original corpse/body remains if possible

This requires replacing the current “always target the player” logic with a
general nearest-valid-target selection based on faction and aggro state.

### 3.10 Implementation phases

- [x] Phase 1a — terrain tile definitions appended for roads, vegetation and lanterns
- [x] Phase 1b — item kinds, item prototypes and weapon entries
- [x] Phase 2 — woods terrain generator: grass, forest clusters, bushes, roads, lanterns
- [x] Phase 3 — huts, doors, windows, furniture, road-facing placement
- [x] Phase 4 — river/pond generation and bridges
- [x] Phase 5 — nine wood-level specifications and unique landmarks
- [x] Phase 6 — signature loot and wood-specific item pools
- [x] Phase 7 — villager/animal factions, zombie aggro and conversion
- [x] Phase 8 — exit placement, connectivity validation, visual polish and tests

### 3.11 Acceptance criteria

For each wood level and several deterministic seeds:

- [x] no void tiles
- [x] exactly one main landmark present
- [x] exit exists and is reachable
- [x] all hut doors are reachable from a road
- [x] every road/water crossing has a bridge
- [x] signature loot exists and is placed in the landmark or biased loot
- [x] villagers and at least one neutral animal are present
- [x] zombies can convert a killed villager
- [x] no generic walls/doors outside buildings, fences, cemeteries or ruins
- [x] caves, temple and time-portal levels still generate as before

---

## 4. Other level revamps

These sections are placeholders. They are intentionally not designed yet to
keep this pass focused on the woods.

### 4.1 Caves (levels 9–17)

Status: **not planned yet**

Keep the current section-based caves generator until a dedicated cave plan is
written.

Possible future topics:

- natural cave shapes instead of room sections
- underground rivers and bridges
- mines, shafts and collapsed tunnels
- fungus forests and subterranean lakes
- cave-specific creatures and loot

### 4.2 Temple of the Underminer (levels 18–26)

Status: **not planned yet**

Possible future topics:

- temple halls, chapels and crypts
- pillars, altars and ritual rooms
- golden doors and locked vaults
- undead priests and temple guardians
- relic-focused loot

### 4.3 Time portal / lab (level 27)

Status: **not planned yet**

Possible future topics:

- plasteel corridors, airlocks and laboratories
- observation rooms, drones and tech-spiders
- portal chamber and final encounter
- high-tech weapons and armour

---

## 5. Cross-cutting rules

- Do not renumber recovered tiles, items, critters or weapons; append new IDs.
- Keep `TileDefs` indices stable for save compatibility.
- Make new content data-driven where possible: level specs, loot categories,
  landmark definitions and spawn tables.
- Every map must be fully reachable under the same movement rules as the player
  (doors count as openable, bridges count as passable).
- Avoid void tiles in all generators.
- Unique landmarks should be visible or strongly hinted at from the road.
- Keep performance acceptable: 80×46 maps, A*/BFS pathfinding and a small
  number of static light sources are fine.
- Preserve the original game’s tone, text style and ASCII presentation.

---

## 6. Testing strategy

Add deterministic generator tests:

- map-size, no-void and reachability checks
- one-landmark-per-level check
- hut-door/road connectivity check
- bridge/river consistency check
- signature-loot presence check
- villager/animal population check
- zombie-to-villager conversion test

Add AI tests:

- idle zombie stays still most of the time
- zombie ignores distant targets
- zombie attacks a nearby villager
- killed villager becomes a zombie after the conversion delay
- neutral animals flee instead of attacking

Add visual smoke tests:

- one screenshot per wood level with the map revealed
- verify roads, landmark, huts, water and exit are visible
- verify the level’s signature landmark is readable at a glance

---

## 7. Open questions

- Should food require a hunger system, or only heal/buff for now?
- Should lanterns emit real FOV light immediately, or remain visual markers in the first pass?
- Should villagers be recruitable, or only ambient/rescue NPCs?
- Should bridges be destructible or repairable?
- Should huts be enterable with interiors on the same map, or kept as single-room buildings?
- Should the nine signature categories be fixed per level (as proposed) or fully shuffled across a run?
- How dangerous should zombie conversion be — immediate rise, delayed rise, or only on certain levels?
- Should level 8 always use the chapel crypt exit, even when the random exit would otherwise be an edge?
