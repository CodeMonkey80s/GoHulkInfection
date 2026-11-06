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
* a terminal of at least 80×50 cells with UTF-8 support
* the `tcell` dependency is fetched by `go mod tidy` / `go build`

The game writes `hi_save.json` in the current directory when you save.

## Controls

The player-facing command set is intentionally melee-focused.

| Key | Action |
| --- | --- |
| `8 9 7 6 5 4 3 2 1` | move in eight directions / wait |
| arrows | move (or move the cursor in look mode) |
| `f` | choose a direction and attack the adjacent enemy |
| `a` | use the Hands item / report the weapon ready |
| `c` `d` `e` | use, swap or consume Belt 1–3 |
| `A` `B` `C` `D` `E` | drop Hands, Body, Belt 1–3 |
| `v` | smash an adjacent door, wall or piece of furniture |
| `,` | fast pick up |
| `i` | player-slot overlay |
| `l` | look around |
| `m` | expand/collapse latest messages |
| `@` | character information |
| `x` | skill levels |
| `O` | options |
| `?` | help / key bindings |
| `S` | save and quit |
| `>` | use portal / stairs / elevator |
| `]` or `F7` | jump directly to the next map |
| `[` or `F8` | jump directly to the previous map |
| `Q` | quit |

Screens generally close with `z`, `q`, or `Esc`.  The player is defenceless
without a melee weapon; there is no unarmed attack.  Throwing is deferred for
now.

The player has five equipment slots: Hands, Body, Belt 1–3.  Body armour is
worn passively and removed/dropped with `B`.  Belt 2 and 3 unlock at levels 3
and 6.

## Screens

* title, intro story, class selection
* scrolling map window with FOV and remembered terrain
* message log and message-buffer screen
* inventory, equipment and ground-item screen with item data
* character sheet, skill-level screen, options screen, key-binding screen
* look mode with tile highlighting
* save/load, death and victory screens

## Game systems

* 80×46 maps assembled from the recovered 10×40 cube-section templates and
  sized to fit the viewport without scrolling
* theme terrain fills the unused section space (woods/caves/temple/plasteel), so
  explored maps no longer contain void holes
* portals are placed on reachable tiles and the starting room is connected to
  the main corridor grid
* recursive shadow-casting FOV (the algorithm credited in `contributors.txt`)
* four level themes: woods, caves, temple of the underminer, time portal
* woods levels 0–8 use an outdoor terrain pass with forest clusters, clearings,
  a dirt road, lantern posts, road-facing huts, rivers/ponds, bridges, one
  unique landmark per level and wood-specific signature loot
* 24-bit woods palette with deterministic tile variation, true-colour UI
  panels and 16-colour DOS fallback
* turn/energy based critters with zombie, observer and animal AI
* melee combat with recovered weapon AV/DV and damage values
* throwing is deferred; HE grenade and molotov data are retained for a future pass
* recovered weapon table (attack/defense value and damage range per subtype)
* five survivor classes with distinct attributes, skills and starting kits
* HP and AP: body armour provides a depletable armour-point pool
* armour vials restore AP from belt slots; acid bypasses AP
* blood, corpses and limbs
* five player slots: Hands, Body and three belt slots; two belt slots start unlocked
* one-line HUD with HP/AP and the ready weapon; one-line message bar
* `m` expands the latest eight messages; `i` opens a non-scrolling slot overlay
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
