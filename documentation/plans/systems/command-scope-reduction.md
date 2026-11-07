# Command & Scope Reduction — Melee-Focused Redesign

Status: **done**
Parent: [Feature plans](../README.md)  
Scope: global player command set and combat systems; first validated in the woods

Implementation status: first pass implemented and merged. Item names have since
been refitted by the [item theme review](../../reference/item-review.md): the HE
grenade and molotov data now appear as holy water and holy fire flasks.

This document proposes cutting the command surface down to what fits a
zombie-slashing, melee-focused roguelike. The goal is not to remove depth — it
is to remove systems that do not support the theme and to make the remaining
commands carry more weight.

---

## 1. Design goals

- Melee is the core combat language.
- The player is **defenceless without a weapon**; there is no unarmed attack.
- There are **no ranged weapons, ammo or reloading**.
- Targeting exists only for **throwing high-explosive grenades or molotov
  cocktails**.
- Throwing range is limited to the current FOV range; you cannot target
  arbitrary distant tiles.
- Gore, corpses, dismemberment, fire and environmental danger stay prominent.
- The command set should be small enough to learn quickly and fit the tone.

## 2. Recommendation summary

### Keep

| System / command | Why |
|---|---|
| Melee attack | Core of the game |
| Movement, wait/rest | Basic roguelike movement and regeneration |
| Bash / smash | Fits chainsaw-and-blunt-weapon fantasy; breaks doors and furniture |
| Throwing | Only remaining ranged interaction; HE grenades and molotovs |
| Inventory, pickup, drop, equip, use | Essential item loop |
| Armour | Melee survivability |
| Medical items, food, tools | Exploration and survival |
| Look, messages, player screen, skills | Information and character progression |
| Options, help, save/quit | Standard UI |
| Portals / stairs | Level progression |
| Factions, villagers, animals, zombie conversion | Supports the living-world feel planned for woods |
| Blood, corpses, limbs | Core gore identity |

### Keep but simplify

| System | Recommendation |
|---|---|
| Equipment slots | Reduce to melee, armour and a single throwable slot; no ranged slot |
| Skills | Keep only perception, dodging, melee, throwing, metabolism and focusing |
| AI | Keep faction/target layer, but remove ranged monster attacks for now |
| Options | Keep, but prune ranged/ammo/target-line options |

### Cut

| System / command | Reason |
|---|---|
| Ranged weapons | Not part of the melee-focused theme |
| Ammo, clips, reload, check ammo | No ranged weapons |
| Unarmed combat | Player is defenceless without a weapon |
| Aim mode / nearest target / clear target | Melee needs no target selection |
| Target line / current-target status | Only throwing needs a cursor |
| Weapon activation | Melee weapons are always ready |
| Swap weapons | Only one melee slot remains |
| Tactics / stance | Removed; combat is straightforward melee |
| Clean yourself | Blood can fade automatically |
| Search | Loot is visible; no hidden-search minigame for now |
| Operate | Doors auto-open on move; portals use the existing use-portal command |
| Highlight objects | Move to an option/debug toggle |
| Dump options | Developer-only; not a player command |
| Console key / duplicate help key | Keep the message buffer and one help screen |
| Ranged monster attacks | Remove for now; keep melee AI |
| Smoke / flashbang / acid grenades | Only HE grenades and molotovs remain |
| Energy weapons and ammo loot | Removed with ranged weapons |
| Debug keys from normal help | Keep in development builds only |

## 3. Proposed command set

This is the target player-facing command set:

| Key | Command |
|---|---|
| `8 9 7 6 5 4 3 2 1` | move / wait |
| arrows | move (and cursor movement in throw mode) |
| `f` | melee attack |
| `b` | smash adjacent door/wall/furniture |
| `t` | throw grenade / molotov, limited to FOV range |
| `,` | fast pick up |
| `i` | inventory |
| `l` | look |
| `m` | message buffer |
| `@` | player screen |
| `x` | skills |
| `O` | options |
| `?` | help / key bindings |
| `S` | save and quit |
| `>` | use portal / stairs / elevator |
| `Q` | quit |
| `Esc` | cancel / back |
| `Tab` / `Enter` | inventory panel / use item |

Cut from the player-facing set: `A`, `a`, `c`, `C`, `W`, `` ` ``, `=`, `h`,
`/`, `n`, `o`, `R`, `s`, `w`, `T`, and the six debug keys. The recovered bindings
remain available as historical data, but the help screen should show only the
active command set.

## 4. Combat model after the cut

### Melee

- `f` attacks the nearest adjacent valid enemy.
- Moving into an enemy attacks automatically if `optPlayerAttackOnMove` is on.
- Without a weapon, the player cannot attack:
  - message: “you are defenceless without a weapon.”
  - no damage, no improvised fist attack.
- Monsters still use natural attacks (claws, bites) — only the player loses
  unarmed combat.

### Throwing

- `t` enters a throw cursor.
- The cursor can only move to tiles within the current FOV range, in line of
  sight, and only visible tiles.
- Throwables:
  - **HE grenade** — explosion and shrapnel damage
  - **molotov cocktail** — fire patch, burning damage
- Throwing skill affects accuracy/scatter.
- No target line and no persistent target state.
- Throwing at an empty tile is allowed for area denial.
- Grenades/molotovs should be scarce.
- Remove smoke, flashbang, acid and incendiary variants. The molotov replaces
  the incendiary role.

### Environmental hazards

Keep:

- fire
- acid pools
- deep water
- smoke from fires, if already present

Remove:

- acid grenades
- smoke grenades
- flashbang grenades
- incendiary grenade item (replace with molotov)

## 5. Systems impact

### Items

- Remove `KindRanged` and `KindAmmo` from the active item set.
- Keep melee, armour, grenade/throwable, medical, tool, food and body-part
  item kinds.
- Keep only melee weapon subtypes in use: improvised weapons, blades, axes,
  hammers, chainsaw, etc.
- Do not spawn ranged weapons, clips, canisters or energy weapons.
- Old saves containing removed items should drop them with a migration message
  or convert them to scrap.

### Skills

Keep six skills:

- perception
- dodging
- weapon melee
- throwing
- metabolism
- focusing

Remove:

- unarmed combat
- weapon ranged
- searching

### Critters and AI

- Remove ranged attack code paths and ammo-equipped monsters for now.
- Zombies, animals and villagers are melee-only.
- Observer/tech-spider ranged behaviour can return later if those themes need
  it, but not in the woods pass.
- Keep faction-based target selection and villager-to-zombie conversion.

### UI

- Equipment slots: melee, armour, throwable (and optionally tool).
- Remove ranged slot, ammo display, target display and target line.
- Status line shows weapon, stance (if kept), armour and weight.
- Inventory item data shows damage/throw information, not ammo.
- Help screen shows only the active command table.
- Options screen hides or removes ranged/ammo/target-line options.

### Options

Remove or disable:

- `optMechanicBaseTileDV` — only relevant for tile-targeted ranged attacks
- `optMechanicLineOfFireHitChance` — line-of-fire shooting
- `optPlayerUseTargetLine` — melee/throwing do not need a line
- any ranged-only debug options

Keep:

- AI and map options
- melee/throw mechanics
- player convenience options
- theme and colour options

### Saves

- Bump save version.
- On load, drop removed ranged/ammo items and report what was removed.
- Reset removed skills to zero or omit them from the skill screen.
- Ensure old saves do not crash.

## 6. Implementation phases

1. **Curate the active command table** and hide the cut bindings from help.
2. **Remove player ranged/unarmed paths** from combat and input handling.
3. **Remove target state** except for throw-cursor mode.
4. **Prune items, skills, spawn tables and options**.
5. **Implement FOV-limited throwing** for HE grenades and molotovs.
6. **Simplify UI**: equipment slots, status line, inventory item data.
7. **Save migration and tests**.
8. **Woods playtest** to check that melee plus scarce throwables is enough.

## 7. Acceptance criteria

- [ ] No ranged weapons, ammo, clips or reloading exist in normal play.
- [ ] No player unarmed attack; a clear message appears when defenceless.
- [ ] No aim/nearest/clear-target commands or persistent target UI.
- [ ] `t` is the only targeting mode; it is limited to FOV range and visible tiles.
- [ ] Only HE grenades and molotovs are throwable.
- [ ] The player-facing help screen lists only active commands.
- [ ] Woods levels remain playable with melee and scarce throwables.
- [ ] Old saves containing removed content still load safely.
- [ ] Tests cover defenceless attacks, throw range, item pruning and save migration.

## 8. Open questions

- Keep bash as a separate command, or let melee attacks damage adjacent doors/furniture?
- Should molotovs require a bottle + fuel, or remain simple loot items?
- Should grenades explode immediately or have a short fuse?
- Should the player be able to throw while standing in fire/smoke?
- Should any enemy keep a ranged special attack later, or remain entirely melee?
- Should tools such as axes improve bashing and door-breaking, or stay flavour items?
