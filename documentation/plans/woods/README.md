# Woods Plans — Levels 0–8

Status: **in progress**

This folder contains the current plans for the woods theme.

## Documents

- [Level revamp](level-revamp.md) — replaces the reused room-section layout with an outdoor woods generator, adds unique landmarks, huts, roads, water, loot and factions.
- [Colour palette](color-palette.md) — true-colour proposal for the woods theme, including per-level accents and rendering rules.

## Recommended implementation order

1. **Palette foundation** — implemented: true-colour support with 16-colour fallback and the core woods palette.
2. **Terrain generator** — implemented: grass, forest clusters, clearings, bushes, roads and lanterns.
3. **Structures** — implemented: road-facing huts with doors, windows and furniture.
4. **Water and roads** — river/pond generation and bridges.
5. **Level identities** — nine `WoodLevelSpec` definitions and unique landmarks.
6. **Items and signature loot** — new wood-specific weapons, armour, tools, food and per-level signature categories.
7. **Factions and AI** — villagers, neutral animals, zombie aggro and villager conversion.
8. **Exit placement and validation** — random edge/middle/landmark exits, connectivity tests and polish.

The level revamp is the larger structural plan; the palette can be implemented
first because new tiles will need colour roles anyway.

## Related current code

- `internal/game/map.go` — generation, exits, spawn population
- `internal/game/ai.go` — critter turns
- `internal/game/combat.go` — damage, death, corpses
- `internal/game/ui.go` — rendering and palette mapping
- `internal/game/types.go` — items, critters, factions and options
- `internal/game/data_gen.go` — recovered tile/item/critter/weapon tables
