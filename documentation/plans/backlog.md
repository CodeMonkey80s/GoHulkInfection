# Post-Woods Backlog — Improvements to Tackle Later

Status: **Proposed**
Parent: [Feature plans](README.md)
Related: [Woods minimalist snake redesign](woods/minimalist-snake-redesign.md), [Monster AI and cardinal pathfinding](systems/monster-ai-pathfinding.md)
Scope: follow-up improvements identified after the woods snake redesign. Defer until after the next planned iteration; none of these are part of the current woods work.

Save compatibility is intentionally out of scope. Each iteration only needs to create and load saves produced by that same iteration.

## Recommended sequence

Point 1 is complete: the woods snake redesign was merged in PR #26. The remaining order is:

2. **Playtest the full five-map run** — done for the functional pass via automated walk-to-portal, full-run transition and render smoke tests. Manual feel/tuning remains part of step 4.
3. **Woods hardening iteration** — done: generator invariant validation, exhaustive theme × entry/exit matrix tests, end-to-end run test and the CI gate are implemented.
4. **Woods polish iteration** — theme-specific signature loot, difficulty curve, theme atmosphere, map name/exit-direction cues and remaining modal-screen tests.
5. **Next major feature selection** — choose from the roadmap only after the woods run is stable. Keep HUD depth and poisonous-food confirmation as smaller deferred extras.

The individual work items below are grouped by area; the sequence above is the intended order.

## Gameplay and generation

- [x] **Generator invariant validation** — after generation, verify the sealed border, exactly one entry/exit, connected roads, reachable portal, one landmark, no void tiles and valid bridges.
- [x] **Exhaustive map matrix tests** — cover every theme × entry/exit combination across many seeds instead of relying on random sampling.
- [x] **End-to-end run test** — walk Church → shuffled middle themes → Meteor Field → caves through `UsePortal`, not only direct level jumps.
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
- [x] **CI gate** — add GitHub Actions running format checks, tests, race tests, vet and build on every PR.
- [x] **Docs cleanup** — fix remaining historical 80×46 and nine-level references where they could confuse current behaviour.
- [x] **Full five-map run walkthrough** — automated walk-to-portal, portal-transition and render smoke tests cover entry/exit flow, connectivity, FOV blockers, lanterns, bridges, huts and the Meteor Field → caves transition.
