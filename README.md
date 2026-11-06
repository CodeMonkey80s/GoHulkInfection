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

* Go 1.21+ (tested with Go 1.27)
* a terminal of at least 80×50 cells with UTF-8 support
* the `tcell` dependency is fetched by `go mod tidy` / `go build`

The game writes `hi_save.json` in the current directory when you save.

## Controls

These are the key bindings recovered from the original key-binding table.

| Key | Action |
| --- | --- |
| `8 9 7 6 5 4 3 2 1` | move in eight directions / wait |
| arrows | move (or move the cursor in look/aim/throw modes) |
| `A` | activate / deactivate current weapon |
| `a` | aim, select a target |
| `b` | bash a wall or door |
| `c` | check ammo |
| `C` | clear current target |
| `W` | clean yourself |
| `` ` `` | dump options to the message log |
| `=` | dump options to the message log |
| `,` | fast pick up |
| `f` | fight / fire / use weapon |
| `?` or `/` | help / key bindings |
| `h` | highlight spotted critters |
| `i` | inventory |
| `l` | look around |
| `m` | message buffer |
| `n` | select nearest target |
| `o` | operate door / portal |
| `O` | options |
| `@` | character information |
| `R` | fast reload |
| `S` | save and quit |
| `s` | search |
| `w` | swap weapons |
| `x` | skill levels |
| `T` | change tactic (defensive / normal / aggressive) |
| `t` | throw grenade |
| `>` | use portal / stairs / elevator |
| `]` or `F7` | jump directly to the next map |
| `[` or `F8` | jump directly to the previous map |
| `Q` | quit |

Screens generally close with `z`, `q`, or `Esc`.

## Screens

* title, intro story, class selection
* scrolling map window with FOV and remembered terrain
* message log and message-buffer screen
* inventory, equipment and ground-item screen with item data
* character sheet, skill-level screen, options screen, key-binding screen
* look mode, aim mode, grenade-throw cursor
* save/load, death and victory screens

## Game systems

* 80×80 maps assembled from the four recovered 10×40 cube-section templates
* theme terrain fills the unused section space (woods/caves/temple/plasteel), so
  explored maps no longer contain void holes
* portals are placed on reachable tiles and the starting room is connected to
  the main corridor grid
* recursive shadow-casting FOV (the algorithm credited in `contributors.txt`)
* four level themes: woods, caves, temple of the underminer, time portal
* turn/energy based critters with zombie, observer and animal AI
* melee, ranged burst fire, reloading, thrown grenades and explosions
* recovered weapon table (attack/defense value and damage range per subtype)
* armor damage reduction, armor wear, blood, corpses and limbs
* item inventory, equipment slots, weight and burden
* skills, skill checks, experience, levels and free skill points
* doors, portals, searching, bashing, acid/fire/smoke, rain
* options screen wired to the gameplay (AI, doors, target line, god/ghost mode, …)
* JSON save/load

## Fidelity notes

The data tables and most texts are recovered directly from `hi030.exe`; the
combat/AI/XP formulas are a playable reconstruction based on the recovered
options and disassembly.  See `REVERSE_ENGINEERING.md` for addresses, struct
layouts and what was exact vs. approximated.

The original game is copyright © 2006, 2007 by Michał Przybyłowicz.  This port
is a preservation/reimplementation effort and keeps the original notices.
