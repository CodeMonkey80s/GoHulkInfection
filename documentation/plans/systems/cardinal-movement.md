# Cardinal Movement and Unified hjkl Controls

Status: **done**
Parent: [Feature plans](../README.md)
Related: [Player redesign](player-redesign.md), [Player experience adjustments](player-experience-adjustments.md), [Monster AI and pathfinding](monster-ai-pathfinding.md)
Scope: player and monster movement, look/throw cursors, menus, key bindings and map connectivity.

This plan removes diagonal movement completely and makes `h j k l` the one
movement/navigation language across gameplay, overlays and menus.

Implementation status: complete. Cardinal movement, unified `h/j/k/l`
navigation, `.` wait, four-way connectivity and generator paths are in place.

---

## 1. Goals

- Remove diagonal player and monster movement.
- Remove numpad movement and the `u i o / j k l / m , .` grid.
- Use `h j k l` consistently for movement, cursor movement and menu navigation.
- Make map connectivity and generator pathfinding cardinal-only.
- Keep the control set small and predictable.
- Avoid diagonal-only gaps that look passable but are not reachable.

## 2. Gameplay controls

| Key | Action |
|---|---|
| `h` | move west |
| `j` | move south |
| `k` | move north |
| `l` | move east |
| `.` | wait / rest |
| `f` + `h/j/k/l` | bash an adjacent small obstacle |
| `g` | get item from the ground |
| `e` | equipment overlay |
| `v` | view / look around |
| `~` | expand/collapse messages |
| `Space` | close a modal screen without taking an action |
| `a` | use Hands |
| `c` / `d` / `r` | use Belt 1–3 |
| `A` / `B` / `C` / `D` / `R` | drop Hands, Body, Belt 1–3 |
| `@` / `x` / `O` / `?` | character / skills / options / help |
| `S` / `Q` / `>` | save / quit / use portal |

Arrow keys may remain as undocumented cardinal aliases for compatibility, but
the help screen and README should present `h j k l` as the canonical controls.

## 3. Menus and overlays

Use the same four keys everywhere:

- `j` — move selection/scroll down.
- `k` — move selection/scroll up.
- `h` / `l` — move left/right where a horizontal choice exists.
- `Enter` — confirm.
- `Esc` — back/close.

Specific screens:

- **Title:** `j/k` select New game / Continue / Quit; `Enter` confirm.
- **Class selection:** `j/k` choose class; `r` reroll name; `Enter` begin.
- **Equipment overlay:** `j/k` move the selected slot; `h/l` optionally cycle
  slots; `a/c/d/r` still perform slot actions.
- **Options:** `j/k` select option; `h/l` decrease/increase value; `Enter`
  toggle.
- **Skills:** `j/k` select skill; `Enter`/`+` spend a point.
- **Help/messages:** `j/k` scroll; `PgUp/PgDn` page.
- **Look/throw:** `h/j/k/l` move the cursor one cardinal tile.
- **Quit confirmation:** keep `y/n`; optionally `j/k` + `Enter` later.

## 4. Movement and connectivity rules

- Player movement is cardinal only. Moving into an enemy is a bump attack.
- `f` + `h/j/k/l` bashes an adjacent small obstacle without moving.
- Walking into an enemy is the only player attack.
- Monsters move cardinal only.
- Pathfinding, wandering, fleeing and map reachability all use four directions.
- The woods road pathfinder and hut/landmark connectors use cardinal movement.
- Connectivity tests use four-way movement, not eight-way.
- Diagonal gaps between walls no longer count as reachable shortcuts.

## 5. Implementation phases

1. **Input layer** — replace diagonal/numpad handling with `h/j/k/l`; move wait
   to `.`; update bash, look and throw direction handling.
2. **Menu navigation** — add `h/j/k/l` navigation to all screens and overlays.
3. **Movement primitives** — introduce `dirs4` for movement/reachability;
   remove movement uses of `dirs8`.
4. **Generators and validation** — switch woods roads, hut connectors and map
   reachability tests to four-way movement.
5. **Documentation and tests** — update README, help, plans and key-binding
   tests.

## 6. Acceptance criteria

- [x] No player or monster can move diagonally.
- [x] Numpad and `u i o / j k l / m , .` movement are removed.
- [x] `h/j/k/l` work consistently in gameplay, look/throw and menus.
- [x] `.` waits/rests.
- [x] `f` + cardinal direction bashes; moving into enemies attacks.
- [x] All generated maps remain fully reachable under four-way movement.
- [x] Caves, temple and time-portal maps still generate and connect correctly.
- [x] Help/README list the new controls.
- [x] Tests cover movement keys, menu navigation, bash direction and map
      reachability.

## 7. Resolved decisions

- Arrow keys remain hidden cardinal aliases; help/README show `h/j/k/l`.
- `h/l` are no-ops in vertical-only menus.
- `.` waits/rests in gameplay.
- The equipment overlay uses `j/k` to select a slot and `Enter` to use it,
  while `a/c/d/r` remain direct slot actions.
- `Space` closes modal screens without taking an action.
