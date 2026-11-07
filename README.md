# hulk infection — Go port

A from-scratch Go reimplementation of **hulk infection v0.3a — ultimate chainsaw
massacre** (2007, Michał Przybyłowicz), reconstructed by reverse-engineering the
original `hi030.exe`.

The original executable is a 16-bit DOS stub with a 32-bit DJGPP/COFF payload,
UPX-packed and built against PDCurses 2.8.  The Go version recreates the game
systems and uses the data tables recovered from the original binary: item
prototypes, critter prototypes, skill list, tile table, map sections, key
bindings, options, player classes, hints and the original texts.

## Build & run

```sh
cd hulk-infection-go
go build -o hulk-infection .
./hulk-infection
```

Requirements:

* Go 1.24+
* a terminal of at least 80×25 cells with UTF-8 support
* the `tcell` dependency is fetched by `go mod tidy` / `go build`

The game writes `hi_save.json` in the current directory when you save.

## Controls

The player-facing command set is intentionally melee-focused.

| Key | Action |
| --- | --- |
| `h` `j` `k` `l` | move west / south / north / east |
| `.` | wait / rest |
| arrows | move (or move the cursor in look/throw mode) |
| `f` + `h/j/k/l` | bash an adjacent small obstacle |
| `a` | use the Hands item / report the weapon ready |
| `c` `d` `r` | use, swap or consume Belt 1–3 |
| `A` `B` `C` `D` `R` | drop Hands, Body, Belt 1–3 |
| `g` | get the item from the ground |
| `e` | equipment overlay |
| `v` | view / look around |
| `~` | expand/collapse latest messages |
| `@` | character information |
| `x` | skill levels |
| `O` | options |
| `?` | help / key bindings |
| `S` | save and quit |
| `>` | use portal / stairs / elevator |
| `]` or `F7` | jump directly to the next map |
| `[` or `F8` | jump directly to the previous map |
| `Space` | close a modal screen without taking an action |
| `Q` | quit |

Walking into a monster attacks it. Bashing is intentional: use `f` and a
direction. A bash takes a turn, and some containers such as barrels or crates
may drop supplies when smashed.

Screens generally close with `z`, `q`, `Space`, or `Esc`.  The player is defenceless
without a melee weapon; there is no unarmed attack.  Throwing is deferred for
now.

The player has five equipment slots: Hands, Body, Belt 1–3.  Body armour is
worn passively and removed/dropped with `B`.  Belt 2 and 3 unlock at levels 3
and 6.

## Screens

* title, intro story, class selection
* compact 80×23 woods map view with unlimited in-map FOV and remembered terrain
* message log and message-buffer screen
* inventory, equipment and ground-item screen with item data
* character sheet, skill-level screen, options screen, key-binding screen
* look mode with tile highlighting
* save/load, death and victory screens

## Game systems

* five themed 80×23 woods maps, plus 80×46 section maps for the caves, temple
  and time portal
* the woods run is fixed at five maps: Church first, Meteor Field last, with the
  three middle themes shuffled once per run and a saved snake direction sequence
* each woods map has a dense tree border with one entry and one exit gap, a road
  from entry through the central landmark to the exit, lanterns, and optional
  themed water with bridges
* theme terrain fills the unused section space (caves/temple/plasteel), so
  explored maps no longer contain void holes
* portals are placed on reachable tiles and the starting room is connected to
  the main corridor grid
* unlimited woods FOV (the 80×23 map edge is the only limit), with local warm
  light pools around lanterns, fires and meteor rocks; section maps keep
  recursive shadow-casting FOV (the algorithm credited in `contributors.txt`)
* four level themes: woods, caves, temple of the underminer, time portal
* woods maps use outdoor terrain with compact props, a central themed landmark,
  road-facing huts, ponds/rivers, bridges, signature loot, villagers, animals
  and zombie conversion within a small content budget
* 24-bit woods palette with per-theme accents, warm light, deterministic tile
  variation, true-colour UI panels and 16-colour DOS fallback
* turn/energy based critters with faction-aware cardinal A* pathfinding,
  zombie, observer and animal AI
* melee combat with recovered weapon AV/DV and damage values
* throwing is deferred; HE grenade and molotov data are retained for a future pass
* recovered weapon table (attack/defense value and damage range per subtype)
* five survivor classes with distinct attributes, skills and starting kits
* HP and AP: body armour provides a depletable armour-point pool
* armour vials restore AP from belt slots; acid bypasses AP
* blood, corpses and limbs
* five player slots: Hands, Body and three belt slots; two belt slots start unlocked
* one-line HUD with HP/AP and the ready weapon; one-line message bar
* `~` expands the latest eight messages; `e` opens a non-scrolling equipment overlay
* one item maximum per ground tile; dropping onto an occupied tile is refused
* skills, skill checks, experience, levels and free skill points
* auto-opening doors, portals, bashing, acid/fire/smoke, rain
* options screen wired to the gameplay (AI, doors, god/ghost mode, …)
* JSON save/load for the current development format; older saves may break
  between iterations

## Fidelity notes

The data tables and most texts are recovered directly from `hi030.exe`; the
combat/AI/XP formulas are a playable reconstruction based on the recovered
options and disassembly.  See [the reverse-engineering notes](documentation/reference/reverse-engineering.md)
for addresses, struct layouts and what was exact vs. approximated.

The original game is copyright © 2006, 2007 by Michał Przybyłowicz.  This port
is a preservation/reimplementation effort and keeps the original notices.

## Documentation

- [Documentation index](documentation/README.md)
- [Feature plans and roadmap](documentation/plans/README.md)
- [Reverse-engineering notes](documentation/reference/reverse-engineering.md)

## License

MIT — see [LICENSE](LICENSE).
