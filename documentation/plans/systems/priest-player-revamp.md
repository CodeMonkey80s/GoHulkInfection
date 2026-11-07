# Priest Player Revamp — Random Rite, Robe Slots and Holy Gear

Status: **in progress**
Parent: [Feature plans](../README.md)
Related: [Player redesign](player-redesign.md), [Player experience adjustments](player-experience-adjustments.md), [Woods minimalist snake redesign](../woods/minimalist-snake-redesign.md), [Player and monster states](player-monster-states.md)
Scope: player identity, one random starting rite plus one shrine rite, intro, robe/head/slot equipment, priest-themed items and zombie/villager variants.

This plan replaces the current class-selection screen and stat-based survivor
roster with a single protagonist: a priest recently reassigned to the village.
The design target is a minimalist experience: few numbers, clear choices and
replayability through a randomly assigned rite, equipment and map/run variation
rather than perks or hidden statistics.

A separate end-game plan will cover the demon finale and the stone-circle map.
Player and monster states such as bleeding will also get a separate plan. Neither
is part of this revamp.

Implementation status: the first slice is implemented — no class selection,
priest/rite generation, development reroll, Hands/Head/Robe slots, vials and
shards with sequential breakage, no baseline auto-healing, XP stored without
levels and HUD rite charges. The shrine/second rite, flasks, head-item effects,
robe-slot upgrades, villager variants, states and demon end-game remain for
later slices.

---

## 1. Goals

- Make the player a priest, not a generic survivor class.
- Remove class selection entirely.
- Remove player-facing STR/TGH/AGI/INT and complicated derived statistics.
- Remove perks entirely.
- Give the priest one randomly assigned starting Rite.
- Allow exactly one additional random Rite from a single shrine per run.
- Keep Rites simple, active and limited to one charge per map.
- Show the available Rite and its remaining charge in the HUD.
- Remove player-facing levels; keep XP stored for now and decide progression
  later.
- Generate the intro text from the priest’s name, Rite and starting kit.
- Keep the priest’s past mysterious but hinted at through the intro.
- Replace body armour and belt slots with a permanent robe, a head slot and
  installable robe slots.
- Make vials, throwable flasks and armour shards compete for the same simple
  robe-slot resource.
- Make every item feel appropriate for a priest fighting corruption.
- Add living villager types and matching zombie versions.
- Keep the game readable at 80×25 and compatible with the completed woods run.

## 2. Locked design decisions

These decisions are settled and should drive the rewrite of the implementation
plan:

- The priest is always the player character; there is no class-selection screen.
- Traditional STR/TGH/AGI/INT statistics are removed from player-facing play.
- **Perks are removed completely.**
- Every priest can perform the basic actions: move, fight, search, use items,
  open doors, pick up/drop, throw flasks and use portals.
- The priest receives exactly **one random Rite at character creation**.
- A single shrine per run grants exactly **one additional random Rite**.
- Each Rite has **one charge per map** for now.
- The HUD must show that the priest has a usable Rite and its remaining charge.
- Player-facing levels are removed. XP remains stored for now; progression will
  be decided later.
- The robe is always worn and cannot be removed.
- The player has **Hands**, **Head** and **Robe** equipment areas.
- The robe starts with **three slots**.
- Additional robe slots can be installed later during the adventure, not
  necessarily through experience.
- Every robe slot holds exactly one item. Vials, flasks and shards all occupy
  one slot each; there is no size or weight calculation.
- Vials are consumables that restore a small flat amount of health.
- Flasks are throwable items and always occupy one robe slot.
- Shards are armour items that occupy one robe slot and provide a maximum AP
  value.
- Different shards have different maximum AP values.
- Damage drains equipped shards before HP.
- When a shard reaches 0 AP it breaks and disappears, freeing its slot.
- Excess damage carries to the next shard and then to HP.
- Drained shards are not repaired; they are replaced with new shards.
- The current Meteor Field remains the final woods map for this revamp. The
  stone circle, demons and final victory are deferred to a separate end-game
  plan.
- There is no baseline automatic healing. The earlier Convalescence perk is no
  longer part of the plan; any small regeneration must come from a later Rite,
  relic or separate progression decision.
- Player and monster states such as bleeding are deferred to a separate plan
  and are not part of this revamp.

## 3. Player model

Keep only the following player-facing resources:

- **HP** — health.
- **AP** — total armour points supplied by equipped shards.
- **Rite charges** — one use per map for each known Rite.
- **Robe slots** — the inventory limit for vials, flasks, shards and small
  relics.

No player-facing strength, toughness, agility or intelligence. No perks and no
skill sheet. Combat should remain readable through weapon values, HP, AP and
explicit item/Rite effects.

XP may still be awarded and stored internally, but it does not grant levels or
player-facing progression yet. The old level-up, free-point and belt-unlock flow
should be removed.

## 4. Rites

### 4.1 Rules

- Rites are active abilities, not passive perks.
- The priest starts with one random Rite.
- One shrine per run grants a second random Rite, once only.
- Each known Rite has one charge per map.
- Charges refresh when entering a new map.
- Using a Rite consumes its charge for the current map.
- Rites should have one clear effect and no percentage scaling.
- Rites should remain useful without a dedicated statistic such as Faith.

### 4.2 Starting Rite list

The first implementation should keep the list small and readable:

| Rite | Effect |
|---|---|
| **Blessing** | Restore 2 AP to the most damaged equipped shard. |
| **Healing Wounds** | Restore 5 HP to the priest. |
| **Shout** | Enemies within 6 tiles flee for 3 turns. |
| **Banish** | Target undead flees for 5 turns; weak undead are removed from the map. |
| **Immobilize** | Target enemy cannot move for 3 turns. |
| **Smite** | Deal 6 holy damage to one visible enemy, with a small bonus against undead. |
| **Sanctuary** | Enemies lose interest and do not attack for 3 turns. |

Exact values remain tuning targets, but the effects should stay simple.

### 4.3 Shrine Rite

- Exactly one shrine per run can grant a second Rite.
- The granted Rite is random.
- It should not duplicate an already known Rite; reroll until unique.
- The shrine is consumed after use.
- The exact placement rules and appearance are still open.

### 4.4 HUD

The one-line HUD currently shows HP, AP and the ready weapon. It must also make
the Rite usable at a glance.

Example compact form:

```text
HP 24/30   AP 6/8   sickle   Rite: Smite 1
```

The final layout and key binding are still open, but the Rite name and remaining
charge must be visible without opening a modal screen.

## 5. Random generation

At the start of a new game, generate:

1. **Name** — from a priest name table, e.g. Father Aleksy, Barnaba, Cyryl,
   Damian, Eliasz, Feliks, Grzegorz, Hieronim, Ignacy, Jakub, Kazimierz, Leon,
   Mateusz, Mikołaj, Piotr, Rafał, Szymon, Tomasz, Urban, Wiktor.
2. **Rite** — one random Rite from the current list.
3. **Starting kit** — a basic weapon, basic robe, one basic shard, one healing
   vial and an empty or Rite-appropriate third slot.
4. **Intro text** — assembled from the name, Rite and starting kit.

The player never sees a class-selection screen.

### Hidden reroll

For now the hidden `r` reroll regenerates **everything** generated for the new
game: name, Rite, starting kit, intro and run/world generation. This is a
development convenience while the design is still changing. It may be restricted
later.

The reroll remains undocumented in the UI and help screen for now.

## 6. Robe, head and slots

### 6.1 Equipment areas

| Area | Slots | Rules |
|---|---:|---|
| **Hands** | 1 | Weapon or tool. |
| **Head** | 1 | One item with one clear effect. |
| **Robe** | 3 initially | Always worn; cannot be removed. Holds small items. |

The robe can gain more slots during the adventure. The exact source is open:
rewards, rare finds, shrine upgrades or quests are all acceptable; it does not
have to come from experience.

### 6.2 Robe slot contents

Each robe slot holds exactly one item:

- healing vial
- throwable flask
- armour shard
- small relic or other thematic utility item

There is no separate belt, body armour or item-size system.

### 6.3 Armour shards

- Each shard has its own **current AP / maximum AP**.
- Different shard types have different maximum AP values.
- Total armour is the sum of the equipped shards.
- Damage drains the equipped shards one at a time.
- When a shard reaches 0 AP, it breaks and disappears, immediately freeing its
  slot.
- Excess damage carries to the next shard, then to HP.
- There is no repair system.
- Partially used shards remain equipped until they break or are deliberately
  dropped.
- A new shard can be attached to any free robe slot.

Suggested shard identities, with exact AP values still to tune:

- leather shard
- chain shard
- iron shard
- silver shard
- blessed shard

### 6.4 Vials and flasks

- **Vials** restore a small flat amount of health. Exact values are open.
- **Flasks** are throwable and occupy one slot each.
- Flask effect types should each have one clear identity, for example:
  - flammable flask — fire
  - holy smoke flask — smoke/anti-undead
  - other thematic flasks to be defined
- No slot-size distinction between vials and flasks: both are one slot.

### 6.5 Head items

The head slot should carry one clear effect rather than a bundle of statistics.
Examples:

- hat
- crown
- helmet
- monocle of true sight
- iron cage facemask

Exact head items and their effects remain open.

### 6.6 Starting kit

The priest starts with:

- a basic weapon in Hands
- a basic robe with three slots
- one basic armour shard
- one healing vial
- one empty slot, unless the generated Rite or starting kit fills it
- one random Rite

Example provisional kit:

| Area | Example |
|---|---|
| Hands | sickle / staff / military knife |
| Head | empty or simple hat |
| Robe slot 1 | basic leather shard |
| Robe slot 2 | healing vial |
| Robe slot 3 | empty |
| Rite | random: Smite / Blessing / Shout / … |

## 7. Priest-themed items

### 7.1 Robes and shards

The robe is a fixed identity item rather than a swappable armour tier. Armour
comes from shards. Do not carry over the old percentage-based armour table.

Shard values should stay small and readable, with no absorption percentages.

### 7.2 Weapons

- **Sickle** — basic harvest blade.
- **Crescent** — curved blade with a religious silhouette.
- **Holy cross** — blunt symbol, strong against undead.
- **Crozier** — staff with a hooked head.
- **Censer on chain** — flail-like, burning incense.
- **Silver dagger** — effective against undead.
- **War hammer** — heavy militia weapon.
- **Pitchfork / wood axe / club** — improvised village weapons.

### 7.3 Consumables, flasks and relics

- **Healing vial** — small flat heal.
- **Communion wine** — stronger heal, rarer.
- **Flammable flask** — throwable fire.
- **Holy smoke flask** — throwable smoke/anti-undead effect.
- **Holy water vial** — throwable or applied; exact behaviour open.
- **Lampion** — light source or throwable fire, to be decided.
- **Incense** — smoke, possibly repels weak undead.
- **Salt pouch** — temporary barrier or anti-undead ward.
- **Candle** — light source.
- **Oil flask** — fire starter or flask refill.
- **Healing herb / poultice** — healing.
- **Antidote** — poison cure.
- **Rosary** — focus or Rite aid.
- **Prayer book** — lore/searching aid.
- **Silver cross / bone relic / old tome** — chapel relics.

## 8. Zombies, villagers and corruption

Living villagers should have matching zombie versions. When a zombie kills a
living type, the corresponding zombie rises after the conversion delay.

The long-term priest goal is to purge the village from corruption by force,
banishment or curing select villagers. The exact victory/progress conditions
are still open.

| Living villager | Zombie version | Behaviour |
|---|---|---|
| Peasant | Zombie peasant | slow, weak, common |
| Militia | Zombie militia | armoured, armed, tougher |
| Runner / courier | Zombie runner | fast, low HP, persistent |
| Thug / village brute | Zombie thug | high HP, heavy hit, slow |

Additional enemy ideas:

- **Ghoul** — cemetery ambusher, faster than a zombie.
- **Wraith** — chapel/Meteor Field spirit, resistant to physical damage.
- **Plague crow** — fast, harasses from a distance.
- **Possessed villager** — living body controlled by the infection.
- **Church guardian** — armoured undead near the chapel.

Spawn by theme:

- Church: peasants, militia, a few zombies.
- Peasant Houses: peasants, militia, runners, thugs.
- Mill and Pond: peasants, fishers, zombies.
- Cemetery: zombies, ghouls, wraiths.
- Meteor Field: strange glowing creatures, wraiths, possessed villagers.

The number of living villager types in the first implementation pass is still
open.

## 9. End-game and demons — separate plan

The stone circle, the final demon and the end-game victory conditions are **not
part of this revamp**.

A separate end-game plan will cover:

- replacing the final woods map with a stone circle and pillars
- randomly selected demons, initially Azazel, Asmodeus and Astaroth
- unique demon abilities and counterplay
- the final purge/cure/banish victory conditions
- how the end-game connects to the caves/temple/time-portal roadmap

Until that plan is written, the current Meteor Field remains in place so this
revamp does not destabilise the completed woods run.

## 10. Player and monster states — separate plan

Status effects such as **bleeding** are **not part of this revamp**.

A separate [player and monster states](player-monster-states.md) plan will cover
conditions for the priest, villagers, animals, undead and demons. It will be
filled in and decided later.

The state plan is expected to cover at least:

- bleeding
- burning
- poison
- infection/corruption
- fear
- immobilize/stun
- holy or cursed ground effects
- how states stack, expire, are cured and interact with HP/AP

## 11. Other thematic proposals

These are ideas, not locked decisions:

- **Holy ground** — chapel floor and consecrated areas weaken undead.
- **Blessed water** — thrown holy water leaves a short-lived blessed puddle.
- **Flask fire** — flammable flasks create a small fire patch.
- **Searching checks** — hidden compartments in barrels, crates, graves and
  altars.
- **Village survivors** — peasants may ask for help, give hints or drop supplies.
- **Confession/mercy** — sparing or curing a possessed villager may yield a
  relic or hint.
- **Cursed relics** — powerful chapel items with a drawback.

## 12. Implementation phases

1. **Data model** — add `PriestProfile`, one or two known Rites, Hands/Head/Robe
   slots, shard AP state and robe-slot installation.
2. **Generation** — random name, random Rite, starting kit and intro text.
3. **Intro/UI** — remove class selection and player levels; add generated intro,
   HUD Rite display and development reroll.
4. **Equipment** — replace belt/body slots with the permanent robe and universal
   robe slots; implement vials, flasks and shards.
5. **Damage model** — sequential shard depletion, breakage at 0 AP and excess
   damage to HP.
6. **Rites** — one starting Rite, one shrine Rite, one charge per map and Rite
   activation UI.
7. **Items** — weapons, head items, shards, vials, flasks and relics.
8. **Villager/zombie variants** — living types, matching zombies and conversion.
9. **Tests and tuning** — Rite generation, shrine uniqueness, starting kits,
   slot limits, shard breakage, intro text and balance.

## 13. Acceptance criteria

- [x] There is no class-selection screen.
- [x] The player is always a priest with a randomly generated name.
- [x] Perks and numeric skill values are removed.
- [x] The priest starts with exactly one random Rite.
- [ ] A single shrine per run grants exactly one additional random Rite.
- [x] Each Rite has one charge per map, refreshed on entry.
- [x] The HUD shows the available Rite and its remaining charge.
- [x] The intro is generated from the assigned Rite and starting kit.
- [x] Hidden `r` rerolls everything during development.
- [x] STR/TGH/AGI/INT and player-facing levels are removed.
- [x] XP is still stored but does not grant levels yet.
- [x] There is no baseline automatic healing.
- [ ] Any later regeneration must be a Rite, relic or separate progression
      decision, not a perk.
- [x] Equipment is limited to Hands, Head and the permanently worn Robe.
- [ ] The robe starts with three slots and can gain more during the adventure.
- [x] Each robe slot holds exactly one vial, flask, shard or small item.
- [x] Vials restore a small flat amount of HP.
- [ ] Flasks are throwable and occupy one slot.
- [x] Shards provide different maximum AP values.
- [x] Damage drains shards in order; a shard at 0 AP breaks and frees its slot.
- [x] Excess shard damage carries to the next shard and then HP.
- [x] Drained shards are not repaired and can be replaced.
- [x] Every starting item fits the priest theme.
- [ ] Living peasants, militia, runners and thugs have matching zombie versions.
- [ ] Zombie kills convert the matching villager type.
- [x] The game remains readable at 80×25.
- [x] The Meteor Field remains the final woods map until the separate demon
      end-game plan is implemented.
- [x] Bleeding and other player/monster states are deferred to their separate plan.

## 14. Open questions

- What is the exact HUD layout and key binding for activating a Rite?
- How should the HUD show both Rites after the shrine grants a second one?
- What are the final tuning values for the seven starting Rites?
- Where does the single shrine appear, and how is it presented?
- What are the shard types and their maximum AP values?
- How much HP should each vial restore?
- What flask types exist, and what is their exact effect?
- Which head items exist, and what is each one’s single effect?
- How are new robe slots earned during the adventure?
- How many living villager types appear in the first implementation pass?
- What counts as purging corruption: killing, banishing, curing or a mix?
- What will XP eventually be used for, now that levels are removed?
- Where should small regeneration come from now that the Convalescence perk is
  gone?
- What are the final demon/end-game rules? Deferred to a separate plan.
