# GoHulkInfection Documentation

This is the central documentation index. The repository root still keeps the
standard `README.md` as the GitHub landing page; everything else that is not a
build file lives here.

## Start here

- [Project README](../README.md) — build, controls and current gameplay features
- [Feature plans](plans/README.md) — roadmap and links to all planned work
- [Reverse-engineering reference](reference/reverse-engineering.md) — how the original binary was analysed and what was recovered

## Current focus

**Woods minimalist snake redesign complete.** The woods run is now five themed
80×23 maps with a random snake route, unlimited-radius FOV with line-of-sight
blockers and compact content budgets. Caves, temple and time-portal work remains
postponed.

- [Woods level revamp](plans/woods/level-revamp.md) — generator, landmarks, huts, roads, loot, factions and AI (reused by the five-map run)
- [Woods colour palette](plans/woods/color-palette.md) — true-colour palette, accents, lighting and UI
- [Minimalist snake redesign](plans/woods/minimalist-snake-redesign.md) — five themed 80×23 maps with a random snake route
- [Post-woods backlog](plans/backlog.md) — deferred improvements for after the next iteration

## Plan status

| Area | Status | Plan |
|---|---|---|
| Woods level revamp (terrain, landmarks, loot, factions) | Done | [woods/level-revamp.md](plans/woods/level-revamp.md) |
| Woods true-colour palette | Done | [woods/color-palette.md](plans/woods/color-palette.md) |
| Woods minimalist snake redesign (five 80×23 maps) | Done | [woods/minimalist-snake-redesign.md](plans/woods/minimalist-snake-redesign.md) |
| Command and scope reduction (melee-only, limited throwing) | Done | [systems/command-scope-reduction.md](plans/systems/command-scope-reduction.md) |
| Player redesign (five slots, HP/AP, minimal HUD) | Done | [systems/player-redesign.md](plans/systems/player-redesign.md) |
| Player experience adjustments (classes, intro, HUD, viewport) | Done | [systems/player-experience-adjustments.md](plans/systems/player-experience-adjustments.md) |
| Cardinal movement (remove diagonals, unified hjkl controls) | Done | [systems/cardinal-movement.md](plans/systems/cardinal-movement.md) |
| Monster AI and cardinal pathfinding | Done | [systems/monster-ai-pathfinding.md](plans/systems/monster-ai-pathfinding.md) |
| Priest player revamp (random rites, robe slots, holy gear; demon end-game deferred) | In progress | [systems/priest-player-revamp.md](plans/systems/priest-player-revamp.md) |
| Player and monster states (bleeding, burning, poison, fear, etc.) | Not planned yet | [systems/player-monster-states.md](plans/systems/player-monster-states.md) |
| Caves levels 9–17 | Not planned yet | [plans/caves.md](plans/caves.md) |
| Temple of the Underminer levels 18–26 | Not planned yet | [plans/temple.md](plans/temple.md) |
| Time portal / lab level 27 | Not planned yet | [plans/time-portal.md](plans/time-portal.md) |

## Documentation conventions

- **Reference** documents describe what exists and what was recovered.
- **Plan** documents describe intended future work; they are not necessarily implemented.
- Statuses are `Not planned yet`, `Proposed`, `Planned`, `In progress`, `Done`.
- Code references use repository-relative paths such as `internal/game/map.go`.
- New feature plans should get their own folder or file under `documentation/plans/` and be linked from the plan index.
