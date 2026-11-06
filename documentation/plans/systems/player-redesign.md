# Player Redesign — ADOM-Style Slots and Minimal HUD

Status: **proposed**  
Parent: [Feature plans](../README.md)  
Related: [Command and scope reduction](command-scope-reduction.md), [Woods level revamp](../woods/level-revamp.md)  
Scope: player equipment, HP/AP, AV/DV, HUD, message log, item actions and key bindings. Throwing is deferred to a later plan.

This plan redesigns the player around a small, tactical set of ADOM-style
equipment slots. It assumes the melee-only command scope from the current
feature branch: no ranged weapons, no ammo, no unarmed attack, no tactics.
Throwing is deliberately deferred to a later plan; the molotov data is kept for
that future pass.

---

## 1. Goals

- Limit the player to a handful of meaningful equipment choices.
- Make switching between weapon and tool a deliberate tactical decision.
- Remove weight, burden and large-inventory management.
- Introduce HP and AP as the two core survival resources.
- Simplify the HUD and message log so the map gets most of the screen.
- Make every action a direct key binding; no inventory scrolling.
- Keep the chainsaw as the signature melee weapon; preserve molotov data for a
  later throwing pass.
- Fit the grim, simple ASCII presentation.

## 2. Equipment slots

Use five player slots:

| Slot | Hotkey | Accepts | Notes |
|---|---|---|---|
| Hands | `a` | weapon, tool, throwable | The active item. One hand slot only. Throwable actions are deferred. |
| Body | `b` | armour | Worn armour; provides AP and DV penalty. |
| Belt 1 | `c` | any non-body item | Available from the start. |
| Belt 2 | `d` | any non-body item | Unlocked at level 3. |
| Belt 3 | `e` | any non-body item | Unlocked at level 6. |

There is **no backpack and no general inventory**. Five slots total.

### Why one Hands slot?

Two hand slots would let the player carry a weapon and a tool at the same time,
which weakens the intended tactical choice. One Hands slot means:

- holding a chainsaw means no lockpick, shovel, lantern or axe in hand
- switching to a tool costs a turn and exposes the player
- switching back to a weapon before combat is a real decision

This is the core of the redesign.

### Belt unlock progression

Proposal:

- level 1: Belt 1 (`c`)
- level 3: Belt 2 (`d`)
- level 6: Belt 3 (`e`)

Unlock messages appear in the message line. Locked slots are shown in the HUD
as unavailable until earned.

## 3. Slot actions and key bindings

Every slot has one context-sensitive action key and one drop key.

### Movement and directional attack

- Arrows and the numpad move the player in eight directions.
- Moving into an enemy attacks it (bump attack).
- `f` followed by a direction performs a deliberate directional melee attack
  without moving.
- There is no “nearest target” attack. The player chooses the direction.
- If there is no adjacent enemy in that direction, `f` reports that there is
  nothing to attack.

### Throwing

Throwing is **deferred** to a later plan.

- `t` is reserved for throwing but is inactive for now.
- Throwables may still exist in item data and slots, but their action is shown
  as “throwing deferred” in the player-slot overlay.
- When throwing returns, a belt throwable is first swapped into Hands, then
  thrown with `t`. This preserves the deliberate switching decision.

### Lowercase = use / act

| Key | Slot | Action |
|---|---|---|
| `a` | Hands | Use a tool; a weapon reports that it is ready; throwables are deferred |
| `b` | Body | Open the player-slot overlay focused on Body |
| `c` | Belt 1 | Use, swap, or consume the item |
| `d` | Belt 2 | Use, swap, or consume the item |
| `e` | Belt 3 | Use, swap, or consume the item |

### Uppercase = drop

| Key | Slot |
|---|---|
| `A` | Drop Hands item |
| `B` | Drop Body armour |
| `C` | Drop Belt 1 item |
| `D` | Drop Belt 2 item |
| `E` | Drop Belt 3 item |

Dropping succeeds only if the current tile has no item on it.

### Hands actions

`a` is context-sensitive:

- **tool** — use the tool (axe, lockpick, shovel, lantern, etc.)
- **weapon** — report that the weapon is ready; attack with `f` + direction
- **throwable** — unavailable for now; throwing is deferred
- **empty** — report that the player has nothing ready

### Belt actions

`c`, `d`, `e` depend on the item:

| Belt item | Action |
|---|---|
| weapon or tool | Swap it into Hands; the old Hands item moves into that belt slot |
| throwable | Swap it into Hands; throwing itself is deferred |
| armour vial | Consume it and restore AP |
| health pack / food | Consume it and restore HP |
| other | Show a short message describing why it cannot be used |

Belt slots are quick access, not a second inventory.

### Body actions

`b` opens the player-slot overlay focused on Body. Armour is worn
automatically when placed in Body; it is not “used” like a consumable.

The overlay shows armour AP, absorption and DV penalty. There is no separate
armour activation command in this scope.

### Player-slot overlay (`i`)

`i` opens a top overlay showing only the player slots:

- key bindings on the left
- a small ASCII character representation
- item name and stats for the focused item
- an action column on the right: `swap`, `consume`, `inspect`, `drop`,
  `deferred`

Suggested shape:

```text
PLAYER SLOTS                                      action
[a] Hands   chainsaw       dmg 7-44 av 5 dv -12    attack: f + dir
[b] Body    jacket         AP 20/20 absorb 50%     inspect
[c] Belt 1  molotov        fire throwable          swap (deferred)
[d] Belt 2  medkit         heal 25, charges 5      consume
[e] Belt 3  ---                                   -
                 O
                /|\
                / \
```

`a`–`e` work inside the overlay; `A`–`E` drop the corresponding item; `i` or
`Esc` closes it.

### Pickup and ground items

- `,` picks up the item on the current tile.
- Placement order:
  1. Hands if empty
  2. Body if the item is armour and Body is empty
  3. first free belt slot
  4. otherwise “you have no free slot”
- One item maximum per map tile.
- If the current tile already contains an item, dropping is refused.
- There is no item stacking.

This keeps loot decisions immediate and avoids inventory Tetris.

### Smash

Smash stays on **`v`** for now:

- `v` smashes an adjacent door, wall or furniture
- a weapon in Hands is not required, but an axe or chainsaw may improve it
- may be revisited in a later plan

## 4. HP, AP, AV and DV

### HP — health points

- HP is the life pool; death at 0.
- Maximum HP comes from toughness as before.
- Healing items restore HP.

### AP — armour points

- AP is a depletable armour pool provided by the Body item.
- Each armour prototype has `MaxAP`, current `AP`, `Absorb%` and `DV penalty`.
- Incoming damage is split:

```
absorbed = min(AP, ceil(damage * absorbPercent / 100))
AP      -= absorbed
HP      -= damage - absorbed
```

- Acid bypasses AP entirely.
- Fire may partially bypass AP.
- At AP 0 the armour still occupies Body but provides no damage absorption.
- AP does not regenerate automatically. Armour repair kits or repair tools
  restore it.
- AP is shown in the HUD and on the Body slot.

Provisional armour values:

| Armour | Max AP | Absorb | DV penalty |
|---|---:|---:|---:|
| jacket | 20 | 50% | −1 |
| flak armour | 35 | 60% | −2 |
| power jacket | 50 | 70% | −3 |

These numbers are tuning placeholders; the recovered armour factors and the
new AP model should be reconciled during implementation.

### AV and DV

Keep AV/DV, but simplify them for one-hand melee:

- **AV** = melee skill + weapon AV + strength/3
- **DV** = dodging + agility/3 − armour DV penalty
- **Hit chance** = base 50 + AV − DV, clamped 5–95

No tactics command and no stance modifiers. The only defensive equipment
decision is which Body armour to wear.

### Chainsaw and molotov

- **Chainsaw** remains the intended top-tier melee weapon: high damage, poor DV,
  and the best answer to zombie groups if the player commits to it.
- **Molotov** remains the intended top-tier throwable: fire area denial, scarce,
  and strongest when the player has set up a chokepoint.
- Both should feel like rewards, not ordinary loot.

## 5. HUD and message layout

Target 80×50 layout — as little as possible:

```text
row 0      latest message only
row 1..46  map view
row 47     Name   HP 39/39   AP 20/20   Lvl 3
row 48     [a] chainsaw  [b] jacket 20/20  [c] molotov  [d] medkit  [e] ---
```

Rules:

- The latest message always occupies one line at the top.
- `m` expands the message log to the last 8 lines. The expanded log overlays
  the top of the map; `m` or `Esc` collapses it.
- There is no permanent multi-line message buffer.
- HP is green/yellow/red depending on health.
- AP is shown in a steel/blue tone; `AP 0/20` is a warning.
- The second HUD line shows the five slots with their hotkeys and item names.
- Body armour shows its current/max AP in the overlay, not in the permanent HUD.
- No weight, burden, ammo, target, tactic, or XP bar in the main HUD.
- There is no permanent help line; `?` opens the help screen.

### Map overlays

- Look mode (`l`) shows a one-line description and a highlighted tile.
- Message expansion and the player-slot overlay share the top region; only one
  overlay is active at a time.
- Throwing overlays are deferred with the throwing feature.

### Player-slot overlay

The HUD is the permanent inventory summary; `i` opens the detailed overlay
described in section 3. It never scrolls and always shows all five slots.
## 6. Menus without arrow scrolling

Map movement keeps arrows and numpad. Menus should not need arrow scrolling:

- Help: one fixed page of active commands, grouped by slot/movement/UI.
- Skills: six active skills, all visible on one page.
- Player screen: one page, no scrolling.
- Options: split into one or two fixed pages, or use direct letter hotkeys.
  If a list exceeds the screen, use page keys rather than scrolling arrows.

The important requirement is that **inventory management never scrolls**.

## 7. Data model changes

Replace the current inventory/equipment fields with a player-slot structure:

```go
type PlayerSlots struct {
    Hands        *Item
    Body         *Item
    Belt         [3]*Item
    BeltUnlocked int // 1, 2 or 3
}
```

- Remove `Critter.Inventory`, `Critter.Equip`, `Weight`, `MaxWeight` and
  `Burden` from player-facing logic.
- Item weight fields may remain in recovered data for compatibility but are no
  longer displayed or enforced.
- Body armour needs `AP`, `MaxAP`, `AbsorbPercent` and `DV`.
- Ground items become at most one per tile; enforce the invariant in pickup,
  drop and generation.
- Save format bumps again and migrates old saves:
  - old equipped melee → Hands
  - old equipped armour → Body
  - first up to three remaining carried items → Belt
  - extra items are dropped on free tiles or discarded with a message
  - weight fields are ignored

## 8. Implementation phases

1. **Data model** — player slots, belt unlock, AP fields, one-item-per-tile.
2. **Combat** — HP/AP split, simplified AV/DV and chainsaw tuning; preserve molotov data for later.
3. **Item actions** — `a`–`e` context actions and `A`–`E` drops.
4. **HUD** — top message line, bottom HP/AP/slot display, message expansion.
5. **Pickup/drop** — auto-placement and one-item ground rule.
6. **Menus** — fixed help/skills/player screens and non-scrolling options.
7. **Save migration and tests**.
8. **Playtest** — switching, AP attrition and chainsaw feel; throwables remain deferred.

## 9. Acceptance criteria

- [ ] Player has exactly five slots: Hands, Body, Belt 1–3.
- [ ] Belt 2 and 3 unlock through experience.
- [ ] No weight, burden, ammo, target or tactic UI remains.
- [ ] No scrolling inventory.
- [ ] Movement into an enemy attacks; `f` + direction attacks deliberately.
- [ ] `a`–`e` perform item actions; `A`–`E` drop items.
- [ ] Dropping fails when the tile already has an item.
- [ ] HP and AP are visible and AP absorbs damage.
- [ ] Armour vials and health packs work from belt slots.
- [ ] Chainsaw is a top-tier melee option; molotov data is preserved for the
      future throwing plan.
- [ ] Latest message is one line; `m` shows up to eight messages.
- [ ] `i` opens the non-scrolling player-slot overlay.
- [ ] Old saves migrate into the new slots without crashing.

## 10. Resolved decisions

- One Hands slot.
- Attack is directional: bump attack or `f` + direction.
- Throwing is deferred; `t` reserved.
- Belt 2 at level 3, Belt 3 at level 6.
- Body accepts armour only; belt accepts non-body items.
- Pickup auto-places into Hands → Body (if armour) → first free belt slot.
- Belt throwables must be swapped into Hands before throwing.
- Smash stays on `v` for now.
- HP regenerates slowly as before; AP is restored by kits, tools and vials.
- Acid bypasses AP completely for now.
- Body action opens the player-slot overlay; no armour activation command yet.
- Enemy AP is out of scope; monsters will be redesigned separately.
- Chainsaw is noisy and has no fuel requirement.
- Molotovs are loot-only.
- Drop is refused on a tile that already has an item.
- HUD is two lines, as minimal as possible.

## 11. Remaining tuning questions

These do not block implementation:

- How much AP should a vial restore?
- What is the chainsaw noise radius and how does it attract zombies?
- Exact overlay height and arrangement of the ASCII character silhouette.
- Whether the overlay pauses the game or is purely informational.
