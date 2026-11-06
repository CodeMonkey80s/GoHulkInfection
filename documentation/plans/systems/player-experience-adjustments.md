# Player Experience Adjustments — Classes, Intro, HUD and Viewport

Status: **planned**
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
| The Medic | 8 | 11 | 9 | 13 | Sustain specialist; best metabolism and healing, weak melee |
| The Mechanic | 11 | 12 | 8 | 11 | Door/wall breaker; good melee and searching, slow |
| The Prepper | 9 | 10 | 11 | 10 | Balanced generalist; good searching and metabolism, no weakness |

Proposed skills:

| Class | PER | DOD | UNA | MEL | RNG | THR | SRCH | MET | FOC |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Butcher | 35 | 40 | 70 | 90 | 0 | 0 | 35 | 55 | 15 |
| Hunter | 70 | 75 | 35 | 65 | 0 | 40 | 60 | 45 | 30 |
| Medic | 55 | 50 | 30 | 55 | 0 | 0 | 50 | 80 | 45 |
| Mechanic | 45 | 40 | 55 | 75 | 0 | 0 | 55 | 55 | 25 |
| Prepper | 60 | 60 | 40 | 60 | 0 | 30 | 60 | 65 | 40 |

Proposed starting kits (existing prototypes first; dedicated class weapons can
replace the improvised ones later):

| Class | Hands | Body | Belt 1 | Belt 2 |
|---|---|---|---|---|
| Butcher | hammer | jacket | bandage | armour vial |
| Hunter | crowbar | jacket | bandage | armour vial |
| Medic | pipe | jacket | medkit | bandage |
| Mechanic | wrench | jacket | bandage | armour vial |
| Prepper | crowbar | jacket | medkit | bandage |

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
7. Save migration, tests and documentation updates.

## 10. Acceptance criteria

- [ ] Look-mode descriptions are clearly readable in all colour modes.
- [ ] Player starts with two unlocked belt slots and a class-appropriate kit.
- [ ] Five classes exist with distinct attributes, skills and starting equipment.
- [ ] The intro states the infection source goal before play begins.
- [ ] Maps are 80×46 and never scroll during normal play.
- [ ] The permanent HUD is one line: HP, AP and weapon name.
- [ ] The title screen keeps `HULK INFECTION` and the three menu options.
- [ ] Title copyright reads `copyright (c) 2026 Michał Przybyłowicz`.
- [ ] Existing saves migrate without crashing.
- [ ] Caves, temple and time-portal levels still generate and remain reachable.

## 11. Resolved decisions

- Title design: **Containment Briefing**.
- Copyright: `copyright (c) 2026 Michał Przybyłowicz`.
- Belt 3 unlock stays at level 6.
- Map size 80×46 applies to all themes; no scrolling in normal play.
- Class names: Butcher, Hunter, Medic, Mechanic, Prepper.
- Class passives are deferred to a later pass.
- Look mode: keep the cursor highlight; fix the description line readability.
