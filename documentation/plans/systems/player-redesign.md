# Player Redesign — ADOM-Style Slots and Minimal HUD

Status: **proposed**  
Parent: [Feature plans](../README.md)  
Related: [Command and scope reduction](command-scope-reduction.md), [Woods level revamp](../woods/level-revamp.md)  
Scope: player equipment, HP/AP, AV/DV, HUD, message log, item actions and key bindings

This plan redesigns the player around a small, tactical set of ADOM-style
equipment slots. It assumes the melee-only command scope from the current
feature branch: no ranged weapons, no ammo, no unarmed attack, no tactics, and
throwing limited to HE grenades and molotovs.

---

## 1. Goals

- Limit the player to a handful of meaningful equipment choices.
- Make switching between weapon and tool a deliberate tactical decision.
- Remove weight, burden and large-inventory management.
- Introduce HP and AP as the two core survival resources.
- Simplify the HUD and message log so the map gets most of the screen.
- Make every action a direct key binding; no inventory scrolling.
- Keep the chainsaw as the signature melee weapon and the molotov as the
  signature throwable.
- Fit the grim, simple ASCII presentation.

## 2. Equipment slots

Use five player slots:

| Slot | Hotkey | Accepts | Notes |
|---|---|---|---|
| Hands | `a` | weapon, tool, throwable | The active item. One hand slot only. |
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

### Lowercase = use / act

| Key | Slot | Action |
|---|---|---|
| `a` | Hands | Attack with a weapon, use a tool, or throw a grenade/molotov |
| `b` | Body | Inspect armour / activate armour ability (future) |
| `c` | Belt 1 | Use, equip or swap the item |
| `d` | Belt 2 | Use, equip or swap the item |
| `e` | Belt 3 | Use, equip or swap the item |

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

- **weapon** — attack the nearest adjacent enemy
- **tool** — use the tool (axe, lockpick, shovel, lantern, etc.)
- **throwable** — enter throw targeting for HE grenade or molotov

If Hands is empty, `a` reports that the player has nothing ready.

This replaces the old separate `f` attack and `t` throw keys. Movement into an
enemy still auto-attacks when `optPlayerAttackOnMove` is enabled.

### Belt actions

`c`, `d`, `e` depend on the item:

| Belt item | Action |
|---|---|
| weapon or tool | Swap it into Hands; the old Hands item moves into that belt slot |
| throwable | Swap it into Hands (then press `a` to throw) |
| medical or food | Use it immediately; empty slot when consumed |
| other | Show a short message describing why it cannot be used |

Belt slots are quick access, not a second inventory.

### Body actions

`b` currently:

- reports armour name, current/max AP and absorb value
- can later activate powered armour or special body items

Armour is worn automatically when placed in Body. It is not “used” like a
consumable.

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

The old bash key `b` conflicts with Body. Move smash to **`v`**:

- `v` smashes an adjacent door, wall or furniture
- a weapon in Hands is not required, but an axe or chainsaw may improve it
- if we want an even smaller command set, smash can be removed and melee
  attacks can damage adjacent obstacles instead

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

Target 80×50 layout:

```
row 0      latest message only
row 1..45  map view
row 46     Name   HP 39/39   AP 20/20   Lvl 3
row 47     Hands: chainsaw        Body: jacket (AP 20/20)
row 48     [c] molotov  [d] medkit  [e] ---
row 49     a use  A drop  , pick up  m log  ? help
```

Rules:

- The latest message always occupies one line at the top.
- `m` expands the message log to the last 8 lines. The expanded log overlays
  the top of the map; `m` or `Esc` collapses it.
- There is no permanent multi-line message buffer.
- HP is green/yellow/red depending on health.
- AP is shown in a steel/blue tone; `AP 0/20` is a warning.
- Hands is the primary displayed weapon.
- Body shows armour and current/max AP.
- Belt shows only the three quick slots and their hotkeys.
- No weight, burden, ammo, target, tactic, or XP bar in the main HUD.
- Level and XP can stay on the compact player screen.

### Map overlays

- Look mode (`l`) shows a one-line description and a highlighted tile.
- Throw mode (`a` with a throwable) shows target coordinates, range and
  `enter`/`esc`.
- Message expansion and item info share the top region; only one overlay is
  active at a time.

### Item details

No scrolling inventory screen is required. The HUD is the inventory.

Optional: `i` opens a compact equipment-details overlay with the five slots and
their descriptions. It uses the same `a`–`e` keys and never scrolls. If the
command set must stay minimal, omit `i` and show item details through the
message line when an action is attempted.

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
2. **Combat** — HP/AP split, simplified AV/DV, chainsaw/molotov tuning.
3. **Item actions** — `a`–`e` context actions and `A`–`E` drops.
4. **HUD** — top message line, bottom HP/AP/slot display, message expansion.
5. **Pickup/drop** — auto-placement and one-item ground rule.
6. **Menus** — fixed help/skills/player screens and non-scrolling options.
7. **Save migration and tests**.
8. **Playtest** — switching, AP attrition, chainsaw and molotov feel.

## 9. Acceptance criteria

- [ ] Player has exactly five slots: Hands, Body, Belt 1–3.
- [ ] Belt 2 and 3 unlock through experience.
- [ ] No weight, burden, ammo, target or tactic UI remains.
- [ ] No scrolling inventory.
- [ ] `a`–`e` perform item actions; `A`–`E` drop items.
- [ ] Dropping fails when the tile already has an item.
- [ ] HP and AP are visible and AP absorbs damage.
- [ ] Chainsaw is a top-tier melee option; molotov is a top-tier throwable.
- [ ] Latest message is one line; `m` shows up to eight messages.
- [ ] Old saves migrate into the new slots without crashing.

## 10. Open questions

- Confirm one Hands slot rather than two.
- Confirm belt unlock levels 3 and 6.
- Should `b` do anything beyond inspecting armour for now?
- Should `i` exist as a non-scrolling equipment-details overlay?
- Should belt slots accept body armour, or should armour only fit Body?
- Should AP be repairable only with kits, or also slowly regenerate at rest?
- Should `v` remain smash, or should smash be folded into weapon attacks?
- Should `a` replace `f` entirely for attacking?
- Should throwing a belt grenade swap it into Hands first, or throw directly?
- Should acid bypass AP completely, or only a percentage?
