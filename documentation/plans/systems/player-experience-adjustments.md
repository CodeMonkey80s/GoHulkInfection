# Player Experience Adjustments — Classes, Intro, HUD and Viewport

Status: **done**

Implementation status: complete. Classes, starting kits, intro, title screen,
one-line HUD, look-description colours and the 80×46 viewport are implemented.
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

## 2. Look mode description readability

The cursor highlight is acceptable. The problem is the description line: it is
hard to read against the current blue background.

Proposed fix:

- Keep the existing cursor highlight and tile glyph colours.
- Draw the look description on `RolePanelBackground` (`#141812`) with
  `RolePrimaryText` (`#D6CBB0`) in true-colour and 256-colour modes.
- Keep the current white-on-navy fallback in 16-colour mode.
- Use `RoleDimText` for the controls line instead of the map-grey fallback.
- Keep the description on row 45 and the controls on row 46.
- Name the item on the selected tile, for example `Crowbar is lying here.`
- Entering a tile with a ground item prints the same announcement in the
  message line.
- If the description is too long, wrap or truncate it deliberately rather than
  letting it run off-screen.

Acceptance: the look description is clearly readable in all three colour modes.

## 3. Starting belt slots

The player starts with **two unlocked belt slots**:

- Belt 1 (`c`) — available from the start.
- Belt 2 (`d`) — available from the start.
- Belt 3 (`e`) — unlock stays at level 6.

Starting kit placement:

- Belt 1: a healing item (medkit or bandage depending on class).
- Belt 2: starts empty so early loot has a free slot.
- Belt 3: empty until unlocked.

Update `PlayerSlots{BeltUnlocked: 2}`, the starting-kit setup, current save/load
round-trip, HUD/overlay text and tests.

## 4. Class roster

Five playable zombie-fighting classes. Values are first-pass tuning numbers, not
recovered data. Skill arrays use the existing order:
perception, dodging, unarmed, melee, ranged, throwing, searching, metabolism,
focusing.

| Class | STR | TGH | AGI | INT | Identity |
|---|---:|---:|---:|---:|---|
| The Butcher | 14 | 13 | 8 | 6 | Heavy melee bruiser; highest HP and melee, poor dodging and searching |
| The Hunter | 10 | 9 | 13 | 8 | Fast scout; best perception and dodging, decent throwing |
| The Medic | 8 | 11 | 9 | 13 | Sustain specialist; best metabolism and healing, weak melee |
| The Mechanic | 11 | 12 | 8 | 11 | Door/wall breaker; good melee and searching, slow |
| The Prepper | 9 | 10 | 11 | 10 | Balanced generalist; good searching and metabolism, no weakness |

Skills:

| Class | PER | DOD | UNA | MEL | RNG | THR | SRCH | MET | FOC |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Butcher | 35 | 40 | 70 | 90 | 0 | 0 | 35 | 55 | 15 |
| Hunter | 70 | 75 | 35 | 65 | 0 | 40 | 60 | 45 | 30 |
| Medic | 55 | 50 | 30 | 55 | 0 | 0 | 50 | 80 | 45 |
| Mechanic | 45 | 40 | 55 | 75 | 0 | 0 | 55 | 55 | 25 |
| Prepper | 60 | 60 | 40 | 60 | 0 | 30 | 60 | 65 | 40 |

Starting kits (existing prototypes first; dedicated class weapons can replace
the improvised ones later):

| Class | Hands | Body | Belt 1 | Belt 2 |
|---|---|---|---|---|
| Butcher | hammer | jacket | bandage | — |
| Hunter | crowbar | jacket | bandage | — |
| Medic | pipe | jacket | medkit | — |
| Mechanic | wrench | jacket | bandage | — |
| Prepper | crowbar | jacket | medkit | — |

Passive traits are deferred to a later pass. Candidates:

- Butcher: better bashing with heavy weapons.
- Hunter: +1 FOV radius or faster spotting.
- Medic: stronger healing from medical items.
- Mechanic: faster bashing of doors and wooden obstacles.
- Prepper: slightly better loot/search results.

Do not implement traits in this pass.

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

Before this change, maps were 80×80 and scrolled inside an 80×46 viewport.
Scrolling has been removed:

- Map size: **80×46**, the same as the visible map area.
- New maps fill the viewport; no scrolling occurs in normal play.
- Player, exit and all content are placed inside 80×46.
- `GenerateLevel` uses a 2×5 section grid, clipping the final section row to
  the map height, then carves the usual connections.
- `map_quality_test.go` and tests that assumed 80×80 were updated.
- The woods level-revamp acceptance criteria were updated to 80×46.
- The camera code path remains for compatibility with older saves and a future
  large-map option, but new maps never scroll.

Decision: apply 80×46 to **all themes** immediately. The map must never scroll
in normal play.

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

## 8. Title screen — Containment Briefing

Chosen design: **Option C, Containment Briefing**. Keep the words
`HULK INFECTION` and the New game / Continue / Quit menu. Copyright becomes
`copyright (c) 2026 Michał Przybyłowicz`.

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
|  copyright (c) 2026 Michał Przybyłowicz                                      |
+------------------------------------------------------------------------------+
```

The title remains the strongest element, the setting is stated immediately and
the menu stays obvious. No large ASCII illustration is needed.

## 9. Implementation order

1. Class data and starting kits, including two starting belt slots.
2. Map size and viewport changes, plus generator/test updates.
3. One-line HUD.
4. Look-mode description readability.
5. Intro text.
6. Chosen title screen.
7. Current save/load round-trip tests and documentation updates.

## 10. Acceptance criteria

- [x] Look-mode descriptions are clearly readable in all colour modes.
- [x] Player starts with two unlocked belt slots and a class-appropriate kit.
- [x] Five classes exist with distinct attributes, skills and starting equipment.
- [x] The intro states the infection source goal before play begins.
- [x] Maps are 80×46 and never scroll during normal play.
- [x] The permanent HUD is one line: HP, AP and weapon name.
- [x] The title screen keeps `HULK INFECTION` and the three menu options.
- [x] Title copyright reads `copyright (c) 2026 Michał Przybyłowicz`.
- [x] Current saves round-trip through load/save. Older development saves may
      break between iterations.
- [x] Caves, temple and time-portal levels still generate and remain reachable.

## 11. Resolved decisions

- Title design: **Containment Briefing**.
- Copyright: `copyright (c) 2026 Michał Przybyłowicz`.
- Belt 3 unlock stays at level 6.
- Map size 80×46 applies to all themes; no scrolling in normal play.
- Class names: Butcher, Hunter, Medic, Mechanic, Prepper.
- Class passives are deferred to a later pass.
- Look mode: keep the cursor highlight; fix the description line readability.
