# Woods Plans — Five-Map Run

Status: **done** (core woods revamp, palette and minimalist snake redesign implemented)

This folder contains the current plans for the woods theme.

## Documents

- [Level revamp](level-revamp.md) — outdoor woods generator, landmarks, huts, roads, water, loot and factions reused by the five-map run.
- [Colour palette](color-palette.md) — true-colour palette for the woods theme, including per-theme accents and rendering rules.
- [Minimalist snake redesign](minimalist-snake-redesign.md) — implemented five-map 80×23 woods run with themed central landmarks and a random snake route.

## Recommended implementation order

1. **Palette foundation** — implemented: true-colour support with 16-colour fallback and the core woods palette.
2. **Terrain generator** — implemented: grass, forest clusters, clearings, bushes, roads and lanterns.
3. **Structures** — implemented: road-facing huts with doors, windows and furniture.
4. **Water and roads** — implemented: ponds/rivers, banks, reeds, lily pads and road bridges.
5. **Items and signature loot** — implemented: wood-specific food, tools, weapons, armour, potions, relics and signature categories.
6. **Factions and AI** — implemented: villagers, neutral animals, local zombie aggro and villager conversion.
7. **Screen reflow** — implemented: all screens fit 80×25, with an 80×23 woods map view.
8. **Snake run** — implemented: five themed maps, fixed Church/Meteor Field ends, shuffled middle themes and saved direction sequence.
9. **Density and validation** — implemented: compact props, population and loot budgets plus connectivity, theme, save/load and render tests.

The woods redesign is complete. Caves, temple and time portal remain on the
section-based generator until their own plans are written.

## Related current code

- `internal/game/woods_run.go` — theme order, snake entry/exit sequence and validity
- `internal/game/woods_gen.go` — compact 80×23 skeleton, roads, lanterns, water and props
- `internal/game/woods_landmarks.go` — the five themed central landmarks
- `internal/game/woods_population.go` — compact theme-specific critters and loot
- `internal/game/map.go` — level dispatch, FOV, spawning and map helpers
- `internal/game/ai.go` — critter turns
- `internal/game/combat.go` — damage, death, corpses
- `internal/game/ui.go` — 80×25 layout, rendering and palette mapping
- `internal/game/types.go` — items, critters, factions, run data and options
- `internal/game/data_gen.go` — recovered tile/item/critter/weapon tables
