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
| `.` | wait one turn |
| arrows | move (or move the cursor in look/throw mode) |
| `x` + `h/j/k/l` | bash an adjacent small obstacle |
| `a` | use the Hands item / report the weapon ready |
| `b` | use the Head item |
| `c` `d` `e` `f` `g` | use Robe slots 1–5 |
| `r` / `t` | use the first / second known Rite |
| `p` | pick up the item from the ground |
| `i` | inventory / equipment screen |
| `v` | view / look around |
| `~` | expand/collapse latest messages |
| `@` | character information |
| `o` | options |
| `?` | help / key bindings |
| `s` | save and quit |
| `>` | use portal / stairs / elevator |
| `]` or `F7` | jump directly to the next map |
| `[` or `F8` | jump directly to the previous map |
| `Space` | close a modal screen without taking an action |
| `q` or `z` | quit |

Walking into a monster attacks it. Player and monster movement and attacks are
cardinal only (`hjkl`); there are no diagonal attacks. Bashing is intentional:
use `x` and a direction. A bash takes a turn, and some containers such as
barrels or crates may drop supplies when smashed.

On the inventory screen lowercase `a`–`g` use slots and uppercase `A`–`G` drop
them. Uppercase letters do nothing on the map. Screens generally close with
`i`, `z`, `q`, `Space`, or `Esc`. The player is
defenceless without a melee weapon; there is no unarmed attack. Flask throwing
is deferred for now.

The priest has Hands, Head and a permanently worn Robe with three slots. Robe
slots hold vials, shards or small items. Shards provide all AP and break when
drained. The first slice starts with one random Rite; the planned shrine that
grants a second Rite is not implemented yet.

## Screens

* title, generated priest intro and random Rite
* compact 80×23 woods map view with unlimited-radius FOV, line-of-sight blockers and remembered terrain
* message log and message-buffer screen
* inventory, equipment and ground-item screen with item data
* character sheet, equipment/robe screen, options screen, key-binding screen
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
* unlimited-radius woods FOV with line-of-sight occlusion by dense trees,
  walls and closed doors, plus local warm light pools around lanterns, fires and
  meteor rocks; section maps keep recursive shadow-casting FOV (the algorithm
  credited in `contributors.txt`)
* four level themes: woods, caves, temple of the underminer, time portal
* woods maps use outdoor terrain with compact props, a central themed landmark,
  road-facing huts, ponds/rivers, bridges, signature loot, villagers, animals
  and zombie conversion within a small content budget
* 24-bit woods palette with per-theme accents, warm light, deterministic tile
  variation, true-colour UI panels and 16-colour DOS fallback
* turn/energy based critters with faction-aware cardinal A* pathfinding,
  zombie, observer and animal AI
* melee combat with recovered weapon AV/DV and damage values
* monsters attack with explicit damage dice (for example `2d10+2`), always deal at least 1 damage, have HP only and carry no body armour
* throwing is deferred; thematic holy, fire, incense and alchemical flask data is retained for a future pass
* recovered weapon table (attack/defense value and damage range per subtype)
* one priest character with a randomly generated name and starting Rite
* HP and AP: equipped robe shards provide a depletable armour-point pool
* shards drain in slot order, break at 0 AP and can be replaced; acid bypasses AP
* blood, corpses and limbs
* Hands, Head and a permanently worn Robe with three starting slots
* one-line HUD with HP/AP, the ready weapon, Rite charges and the current target glyph, name, damage dice and HP; one-line message bar
* `~` expands the latest eight messages; `i` opens the non-scrolling inventory/equipment overlay
* one item maximum per ground tile; dropping onto an occupied tile is refused
* experience is stored for now; player levels and skill points are removed
* auto-opening doors, portals, bashing, acid/fire/smoke, rain
* options screen wired to the gameplay (AI, doors, god/ghost mode, …)
* JSON save/load for the current development format; older saves may break
  between iterations, and the save is removed when the priest dies

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
