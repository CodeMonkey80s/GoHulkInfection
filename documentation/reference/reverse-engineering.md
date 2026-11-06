# Reverse-engineering report — hi030.exe → Go

This document records how the original 2007 DOS game was unpacked, analysed and
what was recovered.  The Go implementation lives in this directory.

## 1. Original binary

`hi030.exe` (225,152 bytes) reports as:

```
MS-DOS executable, COFF for MS-DOS, UPX compressed,
DJGPP go32 v2.04 DOS extender (stub), autoload "CWSDPMI.EXE"
```

`upx -d` expands it to 638,464 bytes.  The MZ stub is 0x800 bytes; the embedded
COFF starts at file offset **0x800**:

```
machine 0x014c (i386), 13 sections
.text  vaddr 0x1a38  raw 0x1238  size 0x879c8
.data  vaddr 0x89400 raw 0x88c00 size 0x12a00
.bss   vaddr 0x9be00
```

The symbol-table pointer in the COFF header is stale (it points past the file),
so the binary is effectively stripped.  Section virtual addresses equal file
offsets in the unpacked image, which makes disassembly easy:

* **game code and read-only strings:** `0x1a38 – ~0x318e0` (~196 KB)
* PDCurses 2.8: after that
* libstdc++ and DJGPP libc (built with gcc 3.3.2, Nov 30 2003): after that
* **game data:** `0x89400 – 0x9be00`

A leaked source path survives in a `__FILE__` string:

```
d:/djgpp/stuff/hi030/map.cpp
```

The credits in `contributors.txt` name Björn Bergström for the recursive
shadowcasting FOV and Takuji Nishimura / Makoto Matsumoto for the RNG, which
matches the code structure found in the binary.

## 2. Method

1. Unpack with UPX, carve the COFF at 0x800.
2. Recover all printable strings with their virtual addresses.
3. Disassemble the game range with Capstone / objdump.
4. Locate static tables in `.data` by following string pointers.
5. Locate switch tables in `.text` (GCC jump tables are stored inline after a
   `jmp eax`), then parse the case blocks for constructor/stat constants.
6. Confirm UI and behaviour by running the original under DOSBox with
   `CWSDPMI.EXE`, capturing screens and comparing messages.
7. Generate the Go tables from the parsed binary data; implement the systems in
   Go with the recovered constants.

## 3. Recovered tables

All virtual addresses below are in the unpacked image.

| Data | Address | Layout |
| --- | --- | --- |
| item names / descriptions | `0x899e0` | 49 × `{char* name; char* desc}` |
| item spawn switch | `0x1e728` | 49-entry jump table (targets `0x1e7ec…`) |
| critter names | `0x896c0` | 17 × `{char* name; char* unused}` |
| critter spawn switch | `0x18df8` | 17-entry jump table (targets `0x18e3c…`) |
| tile prototypes | `0x8ad60` | 64 × 32-byte structs |
| map cube sections | `0x8a554` | 4 × 400 chars (10 rows × 40 columns) |
| skill names / descriptions | `0x89580` | 9 × `{char* name; char* desc}` |
| skill level names | `0x89660` | 12 × `char*` |
| hints | `0x897e0` | 8 × `char*` |
| random player names | `0x89cf0` | 15 × `char*` |
| key bindings | `0x89b84` | 62 × 20-byte structs |
| options | `0x8a0b0` | 50 × 20-byte structs |

### 3.1 Items

`SpawnItem` allocates `0x84` bytes and dispatches on subtype.  Recovered
`ItemBasic` fields:

| Offset | Meaning |
| --- | --- |
| `+0x00` | item category (`0` body part, `1` armor, `2` melee, `3` ranged, …) |
| `+0x04` | subtype / item id |
| `+0x0c` | weight class / handling value |
| `+0x10` | item category/handling constant (NOT damage) |
| `+0x20` | **weight in grams** (confirmed in the inventory screen) |
| `+0x24` | equip class (`0` melee/ammo, `1` ranged, `2` armor, `3` grenade) |
| `+0x2c/+0x30` | durability current/max |
| `+0x34` | ranged flag |
| `+0x38` | ammo type (= ammo item index) |
| `+0x3c` | current ammo / charges |
| `+0x40` | clip capacity / heal amount |
| `+0x44` | burst size |
| `+0x4c` | reload time |
| `+0x50` | range / grenade radius |
| `+0x54` | weapon subtype (indexes the weapon table at `0x898a0`) |
| `+0x58` | accuracy |
| `+0x5c…+0x7c` | armor resistance values (nine values printed by the item-data screen) |
| `+0x80/+0x81` | flags |

A second table at **`0x898a0`** holds 20 weapon subtypes, 16 bytes each:

```
{ int attackValue; int defenseValue; int damageMin; int damageMax; }
```

The item-data screen reads exactly these fields (`av:`, `dv:`, ` dmg:`) and
indexes them by the item's `+0x54` subtype.  This is the authoritative weapon
damage/AV/DV table.  Examples: improvised pipe `av2 dv2 dmg1-5`, chainsaw
`av5 dv-12 dmg7-44`, vibra knife `av15 dv5 dmg2-6`, pistol
`av7 dv-2 dmg3-11`, HE grenade (g-thrower subtype) `av0 dv0 dmg7-55`.
The port embeds this table as `WeaponDefs`.

49 item prototypes were reconstructed, from `corpse` and body parts through
improvised melee weapons, energy guns, ammunition, grenades and medkits.
Weights like `jacket = 5600 g`, `vibra knife = 1050 g`, `bandage = 300 g`,
`medkit = 1200 g`, `he grenade = 850 g` match the original inventory screenshot
exactly (total 9600 g).

### 3.2 Critters

`SpawnCritter` allocates `0x13c` bytes.  Recovered `Critter` fields:

| Offset | Meaning |
| --- | --- |
| `+0x00` | alive/active flags |
| `+0x04` | critter prototype id |
| `+0x0c` | class (`0` animal, `1` monster, `4` player) |
| `+0x14` | alignment (`0` player, `1` enemy, `2` neutral) |
| `+0x18` | color |
| `+0x1c` | glyph |
| `+0x24/+0x28` | x / y |
| `+0x2c/+0x30` | HP current/max |
| `+0x38` | energy / action counter (1000 at spawn) |
| `+0x3c/+0x40/+0x44/+0x48/+0x4c` | strength / toughness / agility / intelligence / speed |
| `+0x5c…+0x7c` | nine skills |
| `+0x84…` | per-type experience counters |
| `+0x10c/+0x110` | level / spawn rank |
| `+0x114` | AI class (`0` zombie, `2` observer, `3` animal, `4` player) |
| `+0x134…` | carried-item list |

The critter spawn cases set stats and the two randomized skill values
(`dodging`, `unarmed`) directly.  For example the plain zombie case sets
strength 13, toughness 6, agility 4, intelligence 2, speed 1, dodging
`rand(25)+33` and unarmed `rand(15)+79`; max HP is then computed from toughness.

The three player professions are applied in `Map::IncreaseCritterAttributesBasedOnClass`
(`0x1fdf4`):

| Class | STR | TGH | AGI | INT |
| --- | --- | --- | --- | --- |
| murderer | 13 | 13 | 9 | 7 |
| marine | 10 | 11 | 10 | 8 |
| warpcrafter | 7 | 10 | 10 | 13 |

The same function adds the profession skill bonuses to the base player skills.

### 3.3 Tiles

64 tile prototypes, 32 bytes each:

| Offset | Meaning |
| --- | --- |
| `+0x00` | name pointer |
| `+0x04` | type (`0` void, `1` floor, `2` wall, `3` door, `4` diggable, `5` elevator, `6` window, `7` broken window) |
| `+0x08` | CP437 glyph |
| `+0x0c` | CP437 color |
| `+0x10` | flags: bit 0 = opaque, bit 8 = blocks movement |
| `+0x14` | door/operate state |
| `+0x18` | partner tile index (open/closed door pair) |

Examples: grass `.` green passable/transparent; stone wall `░` grey opaque and
blocking; closed wooden door CP437 `■` blocking/opaque with partner = open door;
open door `/` passable; deep water `=` blue blocking but transparent; window
`+` blocking but transparent.  The void tile is forced blocking/opaque.

### 3.4 Map generation

Four 400-character cube sections were recovered:

```
0x8a554  cross / room section
0x8a71c  grid of pillar rooms
0x8a8e4  diamond / ziggurat section
0x8aaac  square hall with pillars
```

Each is 10 rows × 40 columns.  An 80×80 level is assembled from 2 columns ×
8 rows of these sections (with flips in the Go port), and the characters map to
theme tiles:

```
'#' wall   '.' floor   '+'/'|' door   'I' pillar   'o' opening   ' ' void
```

The Go port fills unused section space with the theme's base terrain (grass
and trees outdoors, solid rock in caves/temple, plasteel walls in the lab),
adds connecting corridors, connects the starting room to the corridor grid and
places the portal on the farthest reachable tile.

### 3.5 Key bindings and options

Key bindings are 20-byte records:

```
{ char* name; char* description; int key; int modifier; int actionId; }
```

The recovered table matches the original help screen (`<8> move to north`,
`<f> fight or fire`, …).

Options are also 20-byte records:

```
{ char* name; char* description; int value; int min; int max; }
```

The Go port embeds all 50 names, descriptions, minimum and maximum values.
The displayed values in the original options screen (`50 (100)` for sprite
speed, `30 (100)` for acid chance, `100 (500)` for AV numerator, etc.) match
the recovered `min/max` fields.

### 3.6 Other recovered constants

* `maxHP = toughness * 3` — function `0x4508`.
* `max weight = strength * 5500 g` — the murderer (STR 13) shows 71500 g.
* Hit chance is assembled in `0xabda` as
  `base + AV*AVMod1/AVMod2 - DV*DVMod1/DVMod2` using the option globals.
* `AV` (`0x8932`) adds the relevant weapon skill, agility and weapon accuracy;
  `DV` (`0x8748`) uses dodging, agility/5, weapon defense and burden.
* Armor protection factors are stored in the item prototypes (jacket 80,
  flak 70, power jacket 60); damage is multiplied by factor/100.
* Skill values are clamped to 95 and mapped to the recovered level names
  (`useless`, `terrible`, … `percect`).
* Experience thresholds start at 200 for level 2, matching the original
  `exp:0/200` status line.

## 4. What is exact and what is reconstructed

**Recovered exactly (embedded in `data_gen.go`):**

* all item names/descriptions and their prototype stats
* the 20-entry weapon AV/DV/damage table
* all critter names, attributes, skills, colour/glyph and AI class
* all 64 tile names, glyphs, colours, passability/opacity and door partners
* the four 10×40 map sections
* the nine skills, twelve skill-level names and eight hints
* the three player professions and their stat/skill bonuses
* the 62 key bindings and 50 options
* the original story/title/dedication texts and many message strings

**Reconstructed for playability (documented approximations):**

* per-level spawn tables and population counts
* critical-hit rolls and the exact XP reward formula (the weapon damage,
  AV/DV and armour values themselves are recovered)
* monster AI state machine details beyond the recovered AI class ids
* map section placement/flips and level progression pacing
* weather, fire/acid spread and some screen layouts

## 5. Reproducing the extraction

The analysis scripts used to produce this port were kept outside the Go module
for size reasons.  They depend on Python 3, Capstone and the unpacked
`hi030.exe`; the key steps are:

```sh
upx -d hi030.exe -o hi030_unpacked.exe
dd if=hi030_unpacked.exe of=hi030.coff bs=1 skip=2048
# patch the stale symbol-table pointer, disassemble 0x1a38..0x318e0
# parse the tables listed in section 3
```

The generated Go file `internal/game/data_gen.go` contains the extracted
tables and is the authoritative result of that extraction.

## 6. Original copyright

The original game and its texts are copyright © 2006, 2007 by
Michał Przybyłowicz.  All original notices are preserved in the Go port.
