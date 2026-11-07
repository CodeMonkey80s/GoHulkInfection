# Priest Player Revamp — Weighted Perks, Robe Slots and Holy Gear

Status: **proposed**
Parent: [Feature plans](../README.md)
Related: [Player redesign](player-redesign.md), [Player experience adjustments](player-experience-adjustments.md), [Woods minimalist snake redesign](../woods/minimalist-snake-redesign.md)
Scope: player identity, weighted random perks, intro, robe/head/slot equipment, priest-themed items and zombie/villager variants.

This plan replaces the current class-selection screen and stat-based survivor
roster with a single protagonist: a priest recently reassigned to the village.
The design target is a minimalist experience: few numbers, clear choices and
replayability through randomly assigned perks and gear rather than hidden
statistics.

A separate end-game plan will cover the demon finale and the stone-circle map.
That work is deliberately not part of this revamp.

---

## 1. Goals

- Make the player a priest, not a generic survivor class.
- Remove class selection entirely.
- Remove player-facing STR/TGH/AGI/INT and complicated derived statistics.
- Generate the priest’s perks and starting kit with weighted randomness so
  different runs encourage different builds.
- Make more powerful perks rarer than common perks.
- Generate the intro text from the assigned perks and starting kit.
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
- Perks are assigned randomly at character creation using weighted probabilities.
- More powerful perks have a lower chance to appear.
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

## 3. Player model

Keep only the following player-facing resources:

- **HP** — health.
- **AP** — total armour points supplied by equipped shards.
- **Perks** — passive traits, randomly assigned and later earned.
- **Rites/spells** — active abilities, if kept separate from perks. This split
  is still open.
- **Robe slots** — the inventory limit for vials, flasks, shards and small
  relics.

No player-facing strength, toughness, agility or intelligence. Combat should
remain readable through weapon values, HP, AP and explicit perk effects rather
than hidden derived formulas.

## 4. Perks, rites and progression

### 4.1 Perks

Perks define a run’s build. They are passive unless explicitly stated otherwise.

- Assigned at character creation with weighted probabilities.
- Powerful perks are rarer; rarity should mean “more unusual or build-defining”,
  not strictly stronger in every situation.
- Perks must not duplicate within one run.
- The intro uses the assigned perks for background lines.
- Perks can be gained later through experience and a perk progression tree.
- The tree should stay shallow and readable; exact tiers and prerequisites are
  still open.
- Prefer binary effects or simple “+1 rank” bonuses over percentage stacking.

Example perk themes:

| Perk | Rough role | Intro hint |
|---|---|---|
| **Former Soldier** | better melee; starts with a military knife | “You wore a uniform before you wore the collar.” |
| **Gravedigger** | tougher; starts with a shovel; better grave loot | “You know the weight of wet earth.” |
| **Disgraced Exorcist** | better holy effects; starts with holy water | “Something followed you out of the last village.” |
| **Poacher’s Son** | better throwing; starts with a flask | “You learned to throw before you learned to pray.” |
| **Plague Orphan** | poison resistance; extra food | “The fever took your family, not you.” |
| **Pilgrim** | better searching; starts with a walking staff | “You have walked farther than most maps admit.” |

The exact perk pool, rarity weights and number of starting perks remain open.
The current candidate is three starting perks, with more gained during play.

### 4.2 Rites and spells

Active abilities such as **Shout**, **Banish**, **Blessing**, **Healing Wounds**
and **Immobilize** may be separate from passive perks.

The intended distinction is:

- **Perks** — passive, permanent, no resource cost.
- **Rites/spells** — active, limited by charges or uses, learned from relics,
  shrines or progression.
- Perks may modify rites, for example by adding one use of Banish.

Whether rites are completely separate from the perk tree, and whether the priest
starts with one rite, is still an open question.

## 5. Random generation

At the start of a new game, generate:

1. **Name** — from a priest name table, e.g. Father Aleksy, Barnaba, Cyryl,
   Damian, Eliasz, Feliks, Grzegorz, Hieronim, Ignacy, Jakub, Kazimierz, Leon,
   Mateusz, Mikołaj, Piotr, Rafał, Szymon, Tomasz, Urban, Wiktor.
2. **Perks** — weighted random assignment, no duplicates. The starting count is
   still open; three is the current preference, five is the upper option.
3. **Starting kit** — a basic weapon, basic robe, likely one basic shard and
   whichever vials/flasks fit the assigned perks.
4. **Intro text** — assembled from the name, assigned perks and starting kit.

The player never sees a class-selection screen.

### Hidden reroll

The previous proposal kept a hidden `r` reroll on the intro screen. It is still
under consideration because unlimited rerolling can undermine rare-perk
weights. If kept, it should regenerate name, perks, kit and intro without being
listed in the UI or help screen.

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
- a small starting kit chosen to fit the assigned perks
- likely one basic armour shard so the priest is not armourless

Example provisional kit:

| Area | Example |
|---|---|
| Hands | sickle / staff / military knife |
| Head | empty or simple hat |
| Robe slot 1 | basic leather shard |
| Robe slot 2 | healing vial |
| Robe slot 3 | empty or starting flask |

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
- **Rosary** — focus or rite aid.
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

## 10. Other thematic proposals

These are ideas, not locked decisions:

- **Holy ground** — chapel floor and consecrated areas weaken undead.
- **Blessed water** — thrown holy water leaves a short-lived blessed puddle.
- **Flask fire** — flammable flasks create a small fire patch.
- **Faith checks** — holy perks improve relics and blessed effects.
- **Searching checks** — hidden compartments in barrels, crates, graves and
  altars.
- **Village survivors** — peasants may ask for help, give hints or drop supplies.
- **Confession/mercy** — sparing or curing a possessed villager may yield a
  relic or hint.
- **Cursed relics** — powerful chapel items with a drawback.

## 11. Implementation phases

These phases should be updated again once the perk/rite and progression
questions are resolved.

1. **Data model** — add `PriestProfile`, weighted perks, Hands/Head/Robe slots,
   shard AP state and robe-slot installation.
2. **Generation** — random name, weighted perks, starting kit and intro text.
3. **Intro/UI** — remove class selection; add generated intro and decide on the
   hidden reroll.
4. **Equipment** — replace belt/body slots with the permanent robe and universal
   robe slots; implement vials, flasks and shards.
5. **Damage model** — sequential shard depletion, breakage at 0 AP and excess
   damage to HP.
6. **Perks and rites** — implement the shallow perk tree, perk gain cadence and
   active rites if kept separate.
7. **Items** — weapons, head items, shards, vials, flasks and relics.
8. **Villager/zombie variants** — living types, matching zombies and conversion.
9. **Tests and tuning** — weighted generation, perk invariants, starting kits,
   slot limits, shard breakage, intro text and balance.

## 12. Acceptance criteria

- [ ] There is no class-selection screen.
- [ ] The player is always a priest with a randomly generated name.
- [ ] Perks are assigned with weighted rarity and never duplicate.
- [ ] Powerful perks are rarer than common perks.
- [ ] The intro is generated from the assigned perks and starting kit.
- [ ] STR/TGH/AGI/INT are removed from player-facing play.
- [ ] Equipment is limited to Hands, Head and the permanently worn Robe.
- [ ] The robe starts with three slots and can gain more during the adventure.
- [ ] Each robe slot holds exactly one vial, flask, shard or small item.
- [ ] Vials restore a small flat amount of HP.
- [ ] Flasks are throwable and occupy one slot.
- [ ] Shards provide different maximum AP values.
- [ ] Damage drains shards in order; a shard at 0 AP breaks and frees its slot.
- [ ] Excess shard damage carries to the next shard and then HP.
- [ ] Drained shards are not repaired and can be replaced.
- [ ] Every starting item fits the priest theme.
- [ ] Living peasants, militia, runners and thugs have matching zombie versions.
- [ ] Zombie kills convert the matching villager type.
- [ ] The game remains readable at 80×25.
- [ ] The Meteor Field remains the final woods map until the separate demon
      end-game plan is implemented.

## 13. Open questions

- Should the priest start with three or five perks?
- What is the total perk pool, and what rarity weights should it use?
- How deep is the perk tree, what are its tiers and what cadence grants perks?
- Are rites/spells separate from perks, or part of the same progression?
- Does the hidden reroll remain, and can it fish for rare perks?
- What exactly is in the starting kit, including whether a basic shard is always
  equipped?
- How are new robe slots earned during the adventure?
- What are the shard types and their maximum AP values?
- How much HP should each vial restore?
- What flask types exist, and what gates throwing them?
- Which head items exist, and what is each one’s single effect?
- How many living villager types appear in the first implementation pass?
- What counts as purging corruption: killing, banishing, curing or a mix?
- What are the final demon/end-game rules? Deferred to a separate plan.
