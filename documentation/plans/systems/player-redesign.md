# Player Redesign — ADOM-Style Slots and Minimal HUD

Status: **in progress**
Parent: [Feature plans](../README.md)  
Related: [Command and scope reduction](command-scope-reduction.md), [Woods level revamp](../woods/level-revamp.md)  
Scope: player equipment, HP/AP, AV/DV, HUD, message log, item actions and key bindings. Throwing is deferred to a later plan.

Implementation status: complete.  Slots, HP/AP, the minimal HUD, the
player-slot overlay and directional melee are implemented.  Throwing is
deferred to a separate future plan.

The belt-start and HUD details below are superseded by
[Player experience adjustments](player-experience-adjustments.md): the player
will start with two belt slots and the permanent HUD becomes one line showing
only HP, AP and the weapon name.

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
| Body | `B` (unequip/drop) | armour | Worn automatically; AP is applied immediately. No action key. |
| Belt 1 | `c` | any non-body item | Available from the start. |
| Belt 2 | `d` | any non-body item | Available from the start. |
| Belt 3 | `r` | any non-body item | Unlocked at level 6. |

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

- level 1: Belt 1 (`c`) and Belt 2 (`d`)
- level 6: Belt 3 (`r`)

Unlock messages appear in the message line. Locked slots are shown in the HUD
as unavailable until earned.

## 3. Slot actions and key bindings

Every slot has one context-sensitive action key and one drop key.

### Movement and bump attack

- Arrows and `h j k l` move the player cardinally; `.` waits.
- Moving into an enemy attacks it (bump attack).
- `f` followed by a direction bashes an adjacent small obstacle without moving.
- There is no “nearest target” attack. The player chooses the direction.
- If there is no bashable obstacle in that direction, `f` reports that there is
  nothing to bash.

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
| `c` | Belt 1 | Use, swap, or consume the item |
| `d` | Belt 2 | Use, swap, or consume the item |
| `r` | Belt 3 | Use, swap, or consume the item |

Body armour has no action key. Equipping it immediately applies its AP and DV
penalty; `B` unequips and drops it.

### Uppercase = drop

| Key | Slot |
|---|---|
| `A` | Drop Hands item |
| `B` | Unequip and drop Body armour |
| `C` | Drop Belt 1 item |
| `D` | Drop Belt 2 item |
| `R` | Drop Belt 3 item |

Dropping succeeds only if the current tile has no item on it.

### Hands actions

`a` is context-sensitive:

- **tool** — use the tool (axe, lockpick, shovel, lantern, etc.)
- **weapon** — report that the weapon is ready; attack by moving into an enemy
- **throwable** — unavailable for now; throwing is deferred
- **empty** — report that the player has nothing ready

### Belt actions

`c`, `d`, `r` depend on the item:

| Belt item | Action |
|---|---|
| weapon or tool | Swap it into Hands; the old Hands item moves into that belt slot |
| throwable | Swap it into Hands; throwing itself is deferred |
| armour vial | Consume it and restore AP |
| health pack / food | Consume it and restore HP |
| other | Show a short message describing why it cannot be used |

Belt slots are quick access, not a second inventory.

### Body equipment

Armour is purely passive:

- picking up armour with an empty Body slot equips it automatically
- equipping immediately adds its AP to the player’s current/max AP and applies
  its DV penalty
- there is no armour activation, inspection or use action
- `B` removes the armour from Body and drops it on the ground (if the tile is
  empty)
- the overlay only displays armour stats and the `unequip/drop` action

With a full Body slot, a second armour item cannot be picked up until the
current armour is removed.

### Equipment overlay (`e`)

`e` opens a top overlay showing only the player slots:

- key bindings on the left
- a small ASCII character representation
- item name and stats for the focused item
- an action column on the right: `swap`, `consume`, `inspect`, `drop`,
  `deferred`

Suggested shape:

```text
PLAYER SLOTS                                      action
[a] Hands   chainsaw       dmg 7-44 av 5 dv -12    attack: move into enemy
Body      jacket         AP 20/20 absorb 50%     unequip (B)
[c] Belt 1  molotov        fire throwable          swap (deferred)
[d] Belt 2  medkit         heal 25, charges 5      consume
[r] Belt 3  ---                                   -
                 O
                /|\
                / \
```

`a`, `c`, `d`, `r` work inside the overlay; `A`–`D`, `R` drop the corresponding item; `e` or
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

Smash is intentional:

- `f` followed by a direction bashes an adjacent small obstacle
- a bash takes a turn whether it succeeds or fails
- a weapon in Hands is not required, but an axe, hatchet or saw improves the chance
- some containers such as barrels and crates drop supplies when smashed

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
row 47     HP 39/39   AP 20/20   chainsaw
rows 48-49 clear
```

Rules:

- The latest message always occupies one line at the top.
- `~` expands the message log to the last 8 lines. The expanded log overlays
  the top of the map; `~` or `Esc` collapses it.
- There is no permanent multi-line message buffer.
- HP is green/yellow/red depending on health.
- AP is shown in a steel/blue tone; `AP 0/20` is a warning.
- The HUD shows only HP, AP and the ready weapon name. Slot hotkeys and body
  armour details live in the player-slot overlay.
- No weight, burden, ammo, target, tactic, or XP bar in the main HUD.
- There is no permanent help line; `?` opens the help screen.

### Map overlays

- Look mode (`l`) shows a one-line description and a highlighted tile.
- Message expansion and the player-slot overlay share the top region; only one
  overlay is active at a time.
- Throwing overlays are deferred with the throwing feature.

### Player-slot overlay

The HUD is the permanent inventory summary; `e` opens the detailed overlay
described in section 3. It never scrolls and always shows all five slots.
## 6. Menus without arrow scrolling

Map movement keeps arrows and `h j k l`. Menus should not need arrow scrolling:

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
3. **Item actions** — `a`, `c`, `d`, `r` context actions and `A`–`D`, `R` drops.
4. **HUD** — top message line, bottom HP/AP/slot display, message expansion.
5. **Pickup/drop** — auto-placement and one-item ground rule.
6. **Menus** — fixed help/skills/player screens and non-scrolling options.
7. **Save migration and tests**.
8. **Playtest** — switching, AP attrition and chainsaw feel; throwables remain deferred.

## 9. Acceptance criteria

- [ ] Player has exactly five slots: Hands, Body, Belt 1–3.
- [ ] Belt 1 and 2 start available; Belt 3 unlocks at level 6.
- [ ] No weight, burden, ammo, target or tactic UI remains.
- [ ] No scrolling inventory.
- [ ] Movement into an enemy attacks; `f` + direction bashes small obstacles.
- [ ] `a`, `c`, `d`, `r` perform item actions; `A`–`D`, `R` drop items.
- [ ] Dropping fails when the tile already has an item.
- [ ] HP and AP are visible and AP absorbs damage.
- [ ] Equipping armour immediately applies its AP and DV penalty.
- [ ] Armour has no action key; `B` unequips/drops it.
- [ ] Armour vials and health packs work from belt slots.
- [ ] Chainsaw is a top-tier melee option; molotov data is preserved for the
      future throwing plan.
- [ ] Latest message is one line; `~` shows up to eight messages.
- [ ] `e` opens the non-scrolling equipment overlay.
- [ ] Old saves migrate into the new slots without crashing.

## 10. Resolved decisions

- One Hands slot.
- Attack is bump-only; `f` + direction is intentional bashing.
- Throwing is deferred; `t` reserved.
- Belt 2 available from the start, Belt 3 at level 6.
- Body accepts armour only; belt accepts non-body items.
- Pickup auto-places into Hands → Body (if armour) → first free belt slot.
- Belt throwables must be swapped into Hands before throwing.
- Smash is intentional: use `f` + direction; containers may drop supplies.
- HP regenerates slowly as before; AP is restored by kits, tools and vials.
- Acid bypasses AP completely for now.
- Armour has no action key: equip applies AP immediately; `B` unequips/drops it.
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
