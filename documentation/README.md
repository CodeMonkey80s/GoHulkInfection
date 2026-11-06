# GoHulkInfection Documentation

This is the central documentation index. The repository root still keeps the
standard `README.md` as the GitHub landing page; everything else that is not a
build file lives here.

## Start here

- [Project README](../README.md) — build, controls and current gameplay features
- [Feature plans](plans/README.md) — roadmap and links to all planned work
- [Reverse-engineering reference](reference/reverse-engineering.md) — how the original binary was analysed and what was recovered

## Current focus

**Woods levels 0–8.** Two plans are active:

- [Woods level revamp](plans/woods/level-revamp.md) — generator, landmarks, huts, roads, loot, factions and AI
- [Woods colour palette](plans/woods/color-palette.md) — true-colour proposal for the woods theme

## Plan status

| Area | Status | Plan |
|---|---|---|
| Woods levels 0–8 | Planned | [woods/level-revamp.md](plans/woods/level-revamp.md) |
| Woods true-colour palette | Proposed | [woods/color-palette.md](plans/woods/color-palette.md) |
| Command and scope reduction (melee-only, limited throwing) | Done | [systems/command-scope-reduction.md](plans/systems/command-scope-reduction.md) |
| Player redesign (five slots, HP/AP, minimal HUD) | Proposed | [systems/player-redesign.md](plans/systems/player-redesign.md) |
| Caves levels 9–17 | Not planned yet | [plans/caves.md](plans/caves.md) |
| Temple of the Underminer levels 18–26 | Not planned yet | [plans/temple.md](plans/temple.md) |
| Time portal / lab level 27 | Not planned yet | [plans/time-portal.md](plans/time-portal.md) |

## Documentation conventions

- **Reference** documents describe what exists and what was recovered.
- **Plan** documents describe intended future work; they are not necessarily implemented.
- Statuses are `Not planned yet`, `Proposed`, `Planned`, `In progress`, `Done`.
- Code references use repository-relative paths such as `internal/game/map.go`.
- New feature plans should get their own folder or file under `documentation/plans/` and be linked from the plan index.
