# Feature Plans

Status legend: `Not planned yet` · `Proposed` · `Planned` · `In progress` · `Done`

## Active focus: woods levels 0–8

| Feature | Status | Document |
|---|---|---|
| Woods level revamp — outdoor generator, landmarks, huts, roads, loot, factions | Done | [woods/level-revamp.md](woods/level-revamp.md) |
| Woods true-colour palette and per-level accents | Done | [woods/color-palette.md](woods/color-palette.md) |
| Woods minimalist snake redesign — five 80×23 themed maps | Proposed | [woods/minimalist-snake-redesign.md](woods/minimalist-snake-redesign.md) |

The woods folder has a short overview and implementation order here:
[woods/README.md](woods/README.md).

## Other level themes

| Theme | Levels | Status | Document |
|---|---:|---|---|
| Caves | 9–17 | Not planned yet | [caves.md](caves.md) |
| Temple of the Underminer | 18–26 | Not planned yet | [temple.md](temple.md) |
| Time portal / lab | 27 | Not planned yet | [time-portal.md](time-portal.md) |

The current section-based generator remains in use for these themes until their
own revamp plans are written.

## Cross-cutting systems

| Feature | Status | Document |
|---|---|---|
| Command and scope reduction — melee-only combat, no ranged/unarmed, FOV-limited throwing | Done | [systems/command-scope-reduction.md](systems/command-scope-reduction.md) |
| Player redesign — five ADOM-style slots, HP/AP, minimal HUD | Done | [systems/player-redesign.md](systems/player-redesign.md) |
| Player experience adjustments — classes, intro, one-line HUD, 80×46 viewport | Done | [systems/player-experience-adjustments.md](systems/player-experience-adjustments.md) |
| Cardinal movement — remove diagonals, unified hjkl controls | Done | [systems/cardinal-movement.md](systems/cardinal-movement.md) |
| Monster AI and cardinal pathfinding | Done | [systems/monster-ai-pathfinding.md](systems/monster-ai-pathfinding.md) |

The woods plan intentionally includes the first pass of:

- villager and animal factions
- zombie target selection and villager conversion
- wood-specific items, tools and food
- signature loot categories
- true-colour rendering with fallback

The command reduction applies globally and should be validated in the woods
first. When the remaining systems grow beyond the woods, they should get their
own documents under `documentation/plans/systems/` and be linked here.
