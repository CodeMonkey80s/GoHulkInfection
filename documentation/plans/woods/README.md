# Woods Plans — Levels 0–8

Status: **in progress** (core woods revamp and palette done; minimalist snake redesign proposed)

This folder contains the current plans for the woods theme.

## Documents

- [Level revamp](level-revamp.md) — replaces the reused room-section layout with an outdoor woods generator, adds unique landmarks, huts, roads, water, loot and factions.
- [Colour palette](color-palette.md) — true-colour proposal for the woods theme, including per-level accents and rendering rules.
- [Minimalist snake redesign](minimalist-snake-redesign.md) — proposed five-map 80×23 woods run with themed central landmarks and a random snake route.

## Recommended implementation order

1. **Palette foundation** — implemented: true-colour support with 16-colour fallback and the core woods palette.
2. **Terrain generator** — implemented: grass, forest clusters, clearings, bushes, roads and lanterns.
3. **Structures** — implemented: road-facing huts with doors, windows and furniture.
4. **Water and roads** — implemented: meandering rivers/ponds, banks, reeds, lily pads and road bridges.
5. **Level identities** — implemented: nine `WoodLevelSpec` definitions and unique landmarks.
6. **Items and signature loot** — implemented: wood-specific food, tools, weapons, armour, potions, relics and per-level signature categories.
7. **Factions and AI** — implemented: villagers, neutral animals, local zombie aggro and villager conversion.
8. **Exit placement and validation** — implemented: edge/centre/landmark exits, connectivity checks and visual smoke tests.

The level revamp and palette are implemented. The minimalist snake redesign is
the proposed next structural pass; it reuses the completed terrain, loot,
faction, movement and rendering systems.

## Related current code

- `internal/game/map.go` — generation, exits, spawn population
- `internal/game/ai.go` — critter turns
- `internal/game/combat.go` — damage, death, corpses
- `internal/game/ui.go` — rendering and palette mapping
- `internal/game/types.go` — items, critters, factions and options
- `internal/game/data_gen.go` — recovered tile/item/critter/weapon tables
