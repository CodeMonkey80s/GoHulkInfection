# Player Experience Adjustments — Classes, Intro, HUD and Viewport

Status: **proposed**
Parent: [Feature plans](../README.md)
Related: [Player redesign](player-redesign.md), [Woods level revamp](../woods/level-revamp.md), [Woods colour palette](../woods/color-palette.md)
Scope: onboarding, class roster, intro text, HUD, look mode, starting belt slots and map viewport. Combat scope and the woods generator are unchanged.

This plan collects the small player-facing adjustments agreed after the
true-colour foundation. It supersedes the HUD and belt-start sections of the
player redesign where they differ.

---

## 1. Goals

- Make the first screen and first minutes feel current and intentional.
- Give the player a clear, dark goal before the first turn.
- Make the five classes feel like different ways to survive zombies.
- Remove UI noise: no slot hotkeys in the permanent HUD.
- Make look mode readable without fighting the map colours.
- Keep the whole map on screen and remove scrolling.
- Start the player with enough belt capacity to make early loot choices.

## 2. Look mode visibility

Current look mode draws the selected glyph with a fixed black foreground on a
light-grey background. Items and terrain details can disappear against that
highlight, especially after the true-colour background change.

Proposed fix:

- Use `RoleSelectionBackground` (`#4A5A3A`) and `RoleSelectionText`
  (`#F2E6C9`) for the look/throw cursor in true-colour and 256-colour modes.
- Keep the current black-on-silver highlight in 16-colour fallback.
- Preserve the selected tile's item or critter foreground colour instead of
  forcing black; draw it on the selection background.
- Extend the look line to name what is there, for example:
  `look: grass; item: bandage [12:4] arrows move  enter/esc quit`.
- Keep the map background behind the cursor at the normal map background role;
  only the selected cell is highlighted.

Acceptance: terrain, items and critters remain readable on the selected cell in
all three colour modes.

## 3. Starting belt slots

The player starts with **two unlocked belt slots**:

- Belt 1 (`c`) — available from the start.
- Belt 2 (`d`) — available from the start.
- Belt 3 (`e`) — proposed unlock at level 4 (was level 6).

Starting kit placement:

- Belt 1: a healing item (medkit or bandage depending on class).
- Belt 2: a second consumable (bandage or armour vial).
- Belt 3: empty until unlocked.

Update `PlayerSlots{BeltUnlocked: 2}`, the starting-kit setup, save migration,
HUD/overlay text and tests.

## 4. Class roster

Five proposed zombie-fighting classes. Values are first-pass tuning numbers, not
recovered data. Skill arrays use the existing order:
perception, dodging, unarmed, melee, ranged, throwing, searching, metabolism,
focusing.

| Class | STR | TGH | AGI | INT | Identity |
|---|---:|---:|---:|---:|---|
| The Butcher | 14 | 13 | 8 | 6 | Heavy melee bruiser; highest HP and melee, poor dodging and searching |
| The Hunter | 10 | 9 | 13 | 8 | Fast scout; best perception and dodging, decent throwing |
| The Field Medic | 8 | 11 | 9 | 13 | Sustain specialist; best metabolism and healing, weak melee |
| The Mechanic | 11 | 12 | 8 | 11 | Door/wall breaker; good melee and searching, slow |
| The Prepper | 9 | 10 | 11 | 10 | Balanced generalist; good searching and metabolism, no weakness |

Proposed skills:

| Class | PER | DOD | UNA | MEL | RNG | THR | SRCH | MET | FOC |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Butcher | 35 | 40 | 70 | 90 | 0 | 0 | 35 | 55 | 15 |
| Hunter | 70 | 75 | 35 | 65 | 0 | 40 | 60 | 45 | 30 |
| Field Medic | 55 | 50 | 30 | 55 | 0 | 0 | 50 | 80 | 45 |
| Mechanic | 45 | 40 | 55 | 75 | 0 | 0 | 55 | 55 | 25 |
| Prepper | 60 | 60 | 40 | 60 | 0 | 30 | 60 | 65 | 40 |

Proposed starting kits (existing prototypes first; dedicated class weapons can
replace the improvised ones later):

| Class | Hands | Body | Belt 1 | Belt 2 |
|---|---|---|---|---|
| Butcher | hammer | jacket | bandage | armour vial |
| Hunter | crowbar | jacket | bandage | armour vial |
| Field Medic | pipe | jacket | medkit | bandage |
| Mechanic | wrench | jacket | bandage | armour vial |
| Prepper | crowbar | jacket | medkit | bandage |

Optional class traits, if we want more than stats:

- Butcher: better bashing with heavy weapons.
- Hunter: +1 FOV radius or faster spotting.
- Field Medic: stronger healing from medical items.
- Mechanic: faster bashing of doors and wooden obstacles.
- Prepper: slightly better loot/search results.

Keep traits out of the first pass unless they are cheap to implement.

Implementation note: add starting-kit fields to `ClassDef` or a separate class
kit table instead of hard-coding `switch class` cases in `SetupNewGame`.

## 5. Intro text

Replace the current 2007 intro with a short dark briefing that states the goal
up front. Draft:

```text
T H E   W O O D S   A R E   D E A D

The first body rose seven days ago. By morning the town was gone.
By nightfall the trees were full of things that used to be people.

You are the last one to leave the road.

Somewhere in these woods is the source of the infection.
Find it. Destroy it. Do not let it reach the cities.

You have a weapon, a handful of supplies, and no way back.

press enter to continue
```

The exact line breaks can be adjusted when the screen is implemented.

## 6. Map size and viewport

Current maps are 80×80 and scroll inside an 80×46 viewport. Remove scrolling:

- Map size: **80×46**, the same as the visible map area.
- Camera offset is fixed at `(0,0)`; the map always fills the viewport.
- Player, exit and all content must be placed inside 80×46.
- Adapt `GenerateLevel` from the 2×8 section grid to a 2×4 section grid
  (80×40) plus six rows of theme terrain, then carve the usual connections.
- Update `map_quality_test.go` and any tests that assume 80×80.
- Update the woods level-revamp acceptance criteria from 80×80 to 80×46.
- Keep the camera code path available for a future large-map option, but do not
  use it in normal play.

Open question: apply 80×46 to every theme now, or only to woods first. The
simplest implementation is every theme, because the map viewport is global.

## 7. HUD

Replace the two-line HUD with one line:

```text
HP 39/39   AP 20/20   chainsaw
```

Rules:

- Show current/max HP, current/max AP and the Hands weapon name.
- Do not show player name, level, slot hotkeys, belt items or body armour in
  the permanent HUD.
- If Hands is empty or holds a non-weapon tool, show `none` or the tool name as
  appropriate.
- Keep HP colour grading (good/wounded/critical) and an AP warning colour at 0.
- The player-slot overlay (`i`) remains the place for full slot and armour
  details.
- Row 47 holds the HUD; rows 48–49 stay clear.

## 8. Title screen options

Keep the words `HULK INFECTION` and the New game / Continue / Quit menu.
Copyright becomes `copyright (c) 2026 <author>`.

### Option A — The Dead Signal

Bold, minimal, high-contrast title with a short warning line. Closest to the
current screen but modernised.

```text
################################################################################
#                                                                              #
#      H  U  L  K     I  N  F  E  C  T  I  O  N                               #
#      ==========================================                              #
#                                                                              #
#                    no one is coming to save you                              #
#                                                                              #
#                         [ new game ]                                         #
#                         [ continue ]                                         #
#                         [   quit   ]                                         #
#                                                                              #
#                                                                              #
#                    copyright (c) 2026 <author>                               #
################################################################################
```

### Option B — The Last Clearing

An atmospheric ASCII scene: trees, a lone `@`, zombies and the title in the
sky. Most distinct from the other menus, but more art to maintain.

```text
                    .   .   .   HULK INFECTION   .   .   .
                 ~~~~~      the woods are listening      ~~~~~
                          t    t      T      t
                       t       @      z    t
                     ~~~~~   ..  ...  ..   ~~~~~

                          [ new game ]
                          [ continue ]
                          [   quit   ]

                    copyright (c) 2026 <author>
```

### Option C — Containment Briefing

A dark case-file/containment-report frame. Keeps the title prominent, gives the
setting immediately, and is easy to implement with box-drawing characters.

```text
+------------------------------------------------------------------------------+
|  CASE FILE: HULK INFECTION                    STATUS: CONTAINMENT FAILED      |
|  LOCATION: NORTH WOODS                        DATE: 2026                      |
|------------------------------------------------------------------------------|
|                                                                              |
|  The dead are walking. The road is closed.                                   |
|                                                                              |
|  > new game                                                                  |
|  > continue                                                                  |
|  > quit                                                                      |
|                                                                              |
|  copyright (c) 2026 <author>                                                 |
+------------------------------------------------------------------------------+
```

Recommendation: **Option C**. It fits the grim tone, keeps the menu obvious and
does not depend on a large ASCII illustration.

## 9. Implementation order

1. Class data and starting kits, including two starting belt slots.
2. Map size and viewport changes, plus generator/test updates.
3. One-line HUD.
4. Look-mode highlight and description.
5. Intro text.
6. Chosen title screen.
7. Save migration, tests and documentation updates.

## 10. Acceptance criteria

- [ ] Look mode keeps terrain, items and critters readable on the selected cell.
- [ ] Player starts with two unlocked belt slots and a class-appropriate kit.
- [ ] Five classes exist with distinct attributes, skills and starting equipment.
- [ ] The intro states the infection source goal before play begins.
- [ ] Maps are 80×46 and never scroll during normal play.
- [ ] The permanent HUD is one line: HP, AP and weapon name.
- [ ] The title screen keeps `HULK INFECTION` and the three menu options.
- [ ] Title copyright reads `copyright (c) 2026 <author>`.
- [ ] Existing saves migrate without crashing.
- [ ] Caves, temple and time-portal levels still generate and remain reachable.

## 11. Open decisions

- Which title design: A, B or C?
- What name should appear in the copyright line?
- Belt 3 unlock: level 4 as proposed, or keep level 6?
- Should the 80×46 map size apply to every theme immediately?
- Keep the five class names as proposed, or rename any of them?
- Implement class traits now, or keep the first pass to stats and kits only?
- Does “look mode barely visible” mean the cursor highlight, the description
  line, or both?
