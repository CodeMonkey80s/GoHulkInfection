# Post-Woods Backlog — Improvements to Tackle Later

Status: **Proposed**
Parent: [Feature plans](README.md)
Related: [Woods minimalist snake redesign](woods/minimalist-snake-redesign.md), [Monster AI and cardinal pathfinding](systems/monster-ai-pathfinding.md)
Scope: follow-up improvements identified after the woods snake redesign. Defer until after the next planned iteration; none of these are part of the current woods work.

Save compatibility is intentionally out of scope. Each iteration only needs to create and load saves produced by that same iteration.

## Gameplay and generation

- [ ] **Generator invariant validation** — after generation, verify the sealed border, exactly one entry/exit, connected roads, reachable portal, one landmark, no void tiles and valid bridges.
- [ ] **Exhaustive map matrix tests** — cover every theme × entry/exit combination across many seeds instead of relying on random sampling.
- [ ] **End-to-end run test** — walk Church → shuffled middle themes → Meteor Field → caves through `UsePortal`, not only direct level jumps.
- [ ] **Theme-specific signature loot** — map loot categories to themes instead of using a random category assignment.
- [ ] **Difficulty curve** — make the Church clearly safest and the Meteor Field clearly hardest; tune zombie pressure, villager presence and loot by theme.
- [ ] **Theme atmosphere** — add cemetery fog/smoke, stronger meteor glow, and more mill/pond and church graveyard detail.
- [x] **AI line-of-sight** — make zombie/villager perception use the same sight blockers as the player so they cannot aggro through dense trees, walls or closed doors.
- [ ] **Dangerous consumables** — distinguish poisonous food more clearly, or add confirmation before eating it.

## UI and UX

- [ ] **UI location cues** — show the current map name and exit direction briefly on entry; consider run progress in messages.
- [ ] **HUD depth** — consider level/XP or map progress without breaking the one-line HUD.
- [ ] **Modal-screen tests** — render inventory, options, help, skills, equipment and message screens at 80×25 and assert nothing writes past row 24.

## Engineering and process

- [x] **Code cleanup** — remove leftover nine-level woods helpers, unused `ExitMode`, and obsolete water-generation paths if no longer used.
- [ ] **CI gate** — add GitHub Actions running format checks, tests, race tests, vet and build on every PR.
- [x] **Docs cleanup** — fix remaining historical 80×46 and nine-level references where they could confuse current behaviour.
- [ ] **Playtest pass** — verify lanterns, doors, fences, bridges and FOV blockers behave sensibly in actual play.
