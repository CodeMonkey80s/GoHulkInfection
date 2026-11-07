# Priest Player Revamp — Random Backgrounds, Skills and Holy Gear

Status: **proposed**
Parent: [Feature plans](../README.md)
Related: [Player redesign](player-redesign.md), [Player experience adjustments](player-experience-adjustments.md), [Woods minimalist snake redesign](../woods/minimalist-snake-redesign.md)
Scope: player identity, random generation, stats/skills, intro, starting gear, priest-themed items, zombie/villager variants.

This plan replaces the current class-selection screen and stat-based survivor
roster with a single protagonist: a priest recently reassigned to the village.
The priest’s name, past, skills and starting kit are generated once, with a
hidden `[r]` reroll.

---

## 1. Goals

- Make the player a priest, not a generic survivor class.
- Remove class selection entirely.
- Generate name, background, skills and starting kit randomly.
- Generate the intro text from the assigned background and skills.
- Keep the priest’s past mysterious but hinted at through the intro.
- Replace STR/TGH/AGI/INT with a compact skill/perk system.
- Give the priest three belt slots from the start.
- Make every item feel appropriate for a priest fighting evil.
- Add living villager types and matching zombie versions.
- Keep the game readable at 80×25 and compatible with the completed woods run.

## 2. Do we still need player statistics?

**Recommendation: remove traditional STR/TGH/AGI/INT.**

Reasons:

- They duplicate the job of skills and perks.
- They hide important math from the player.
- Random stat rolls make the intro harder to explain and balance.
- The priest theme is stronger if abilities are explicit: “you can throw”,
  “you can search”, “you can bless”.
- The game is already minimalist; four invisible stats add noise.

Keep only:

- **HP** — fixed base plus Endurance bonuses.
- **AP** — armour points, as now.
- **Skills 0–3** — explicit abilities and gates.
- **Perks** — short binary background traits used for flavour, intro and small
  bonuses.

Combat values become:

- attack from Zeal + weapon AV + weapon skill
- defence from Reflex + armour DV
- HP from base + Endurance
- holy item power from Faith
- throwing range/accuracy from Throwing
- hidden loot from Searching

This is easier to explain on the intro screen and easier to balance.

## 3. Skills and perks

### 3.1 Skills

Six skills, each rated **0–3**:

| Skill | 0 | 1 | 2 | 3 |
|---|---|---|---|---|
| **Zeal** | weak melee | trained | strong | crusader |
| **Reflex** | clumsy | basic dodge | agile | untouchable |
| **Faith** | no holy bonus | use holy water | stronger blessings | relic mastery |
| **Searching** | cannot search containers | basic search | hidden compartments | secret caches and rare finds |
| **Throwing** | cannot throw | short range | medium range | long range and accuracy |
| **Endurance** | fragile | normal | tough | hardened |

Explicit gates:

- **Throwing 0** means the priest cannot throw holy water, lampions or other
  throwables at all.
- **Searching 0** means the priest cannot search barrels, crates or hidden
  compartments for extra loot.
- **Faith** improves holy water, lampion fire against undead, relics and
  blessings.
- **Endurance** raises HP and resistance to poison, fire and disease.
- **Zeal** and **Reflex** replace the old attack/defence stat math.

### 3.2 Perks

Two random binary perks per priest. They add flavour, a small bonus and intro
lines.

Suggested perks:

| Perk | Effect | Intro hint |
|---|---|---|
| **Former Soldier** | +1 Zeal; starts with a military knife | “You wore a uniform before you wore the collar.” |
| **Gravedigger** | +1 Endurance; starts with a shovel; better grave loot | “You know the weight of wet earth.” |
| **Disgraced Exorcist** | +1 Faith; starts with holy water | “Something followed you out of the last village.” |
| **Poacher’s Son** | +1 Throwing; starts with a lampion | “You learned to throw before you learned to pray.” |
| **Plague Orphan** | poison resistance; extra food | “The fever took your family, not you.” |
| **Pilgrim** | +1 Searching; starts with a walking staff | “You have walked farther than most maps admit.” |

Perks must not duplicate. The intro uses both perks and the two highest skills.

## 4. Random generation

At the start of a new game, generate:

1. **Name** — from a priest name table, e.g. Father Aleksy, Barnaba, Cyryl,
   Damian, Eliasz, Feliks, Grzegorz, Hieronim, Ignacy, Jakub, Kazimierz, Leon,
   Mateusz, Mikołaj, Piotr, Rafał, Szymon, Tomasz, Urban, Wiktor.
2. **Skills** — distribute 8 points across the six skills, max 3 each.
   - at least one of Zeal/Reflex must be at least 1
   - at least one of Faith/Searching/Throwing must be at least 1
3. **Perks** — two unique perks from the list.
4. **Starting kit** — based on skills and perks.
5. **Intro text** — assembled from name, assignment, perks and highest skills.

The player never sees a class-selection screen. The intro screen shows the
generated name, background, skills and starting kit.

### Hidden reroll

`r` on the intro screen regenerates name, skills, perks, starting kit and intro.
It is **not displayed anywhere** in the UI or help screen.

## 5. Intro screen

The intro replaces both the old intro and class selection. It should show:

- generated priest name
- assignment to the village
- two or three short background lines derived from perks
- current skills with levels
- the current threat/goal
- `press enter to begin`

Example structure:

```text
FATHER ALEKSY

The bishop assigned you to a village that no longer answers letters.
You wore a uniform before you wore the collar.
You learned to throw before you learned to pray.

Zeal 2   Reflex 1   Faith 1
Searching 0   Throwing 2   Endurance 1

The dead are walking. Find the source and put it to rest.

press enter to begin
```

The hidden `r` reroll is not listed. It should be implemented but undocumented.

## 6. Starting kit and slots

- The priest starts with **three belt slots unlocked**.
- No belt unlock progression.
- Hands: a basic weapon, usually a sickle, cross, staff or improvised tool.
- Body: a basic robe or cassock.
- Belt 1–3: holy water, lampion, bread or bandage depending on skills/perks.

Example kits:

| Background | Hands | Body | Belt 1 | Belt 2 | Belt 3 |
|---|---|---|---|---|---|
| Former Soldier | military knife | cassock | bandage | holy water | — |
| Gravedigger | shovel | cassock | bandage | holy water | — |
| Disgraced Exorcist | cross | cassock | holy water | bandage | — |
| Poacher’s Son | sickle | cassock | lampion | bandage | — |
| Plague Orphan | staff | cassock | bread | bandage | — |
| Pilgrim | staff | cassock | bandage | holy water | — |

## 7. Priest-themed items

### 7.1 Robes and armour

| Item | AP | Absorb | DV | Notes |
|---|---:|---:|---:|---|
| Novice robe | 8 | 30% | 0 | light, no penalty |
| Priest’s cassock | 15 | 40% | 0 | basic starting robe |
| Penitent’s robe | 22 | 50% | -1 | padded, uncomfortable |
| Exorcist’s vestment | 30 | 55% | -1 | better vs undead |
| Chain under-robe | 40 | 65% | -2 | heavy, hidden armour |
| Grave shroud | 12 | 35% | 0 | light, quiet |
| Leather jerkin | 25 | 50% | -1 | practical, not holy |

### 7.2 Weapons

- **Sickle** — basic harvest blade.
- **Crescent** — curved blade with a religious silhouette.
- **Holy cross** — blunt symbol, strong against undead.
- **Crozier** — staff with a hooked head.
- **Censer on chain** — flail-like, burning incense.
- **Silver dagger** — effective against undead.
- **War hammer** — heavy militia weapon.
- **Pitchfork / wood axe / club** — improvised village weapons.

### 7.3 Consumables and throwables

- **Holy water vial** — throwable, burn/steam damage, bonus against undead.
- **Lampion** — throwable oil lamp, fire/burn damage.
- **Blessed bread** — heals.
- **Communion wine** — heals and steadies.
- **Incense** — smoke, maybe repels weak undead.
- **Salt pouch** — temporary barrier or anti-undead ward.
- **Candle** — light source.
- **Oil flask** — fire starter or lampion refill.
- **Healing herb / poultice** — healing.
- **Antidote** — poison cure.
- **Rosary** — Faith focus.
- **Prayer book** — Faith/searching aid.
- **Silver cross / bone relic / old tome** — chapel relics.

## 8. Zombies and villagers

Living villagers should have matching zombie versions. When a zombie kills a
living type, the corresponding zombie rises after the conversion delay.

| Living villager | Zombie version | Behaviour |
|---|---|---|
| Peasant | Zombie peasant | slow, weak, common |
| Militia | Zombie militia | armoured, armed, tougher |
| Runner / courier | Zombie runner | fast, low HP, persistent |
| Thug / village brute | Zombie thug | high HP, heavy hit, slow |

Additional enemies that fit the theme:

- **Ghoul** — cemetery ambusher, faster than a zombie.
- **Wraith** — chapel/meteor spirit, resistant to physical damage.
- **Plague crow** — fast, harasses from a distance.
- **Possessed villager** — living body controlled by the infection.
- **Church guardian** — armoured undead near the chapel.

Spawn by theme:

- Church: peasants, militia, a few zombies.
- Peasant Houses: peasants, militia, runners, thugs.
- Mill and Pond: peasants, fishers, zombies.
- Cemetery: zombies, ghouls, wraiths.
- Meteor Field: strange glowing creatures, wraiths, possessed villagers.

## 9. Other thematic proposals

- **Holy ground** — chapel floor and consecrated areas weaken undead.
- **Blessed water** — thrown holy water leaves a short-lived blessed puddle.
- **Lampion fire** — lampions create a small fire patch and light.
- **Faith checks** — high Faith improves relics and holy item effects.
- **Searching checks** — hidden compartments in barrels, crates, graves and
  altars.
- **Village survivors** — peasants may ask for help, give hints or drop supplies.
- **Confession/mercy** — optional: sparing a possessed villager may yield a
  relic or hint.
- **Cursed relics** — powerful chapel items with a drawback.

## 10. Implementation phases

1. **Data model** — replace class/stat fields with `PriestProfile`, skills,
   perks and three starting belt slots.
2. **Generation** — random name, skills, perks, starting kit and intro text.
3. **Intro/UI** — remove class selection; add generated intro and hidden reroll.
4. **Combat formulas** — replace stats with Zeal/Reflex/Endurance and skill
   gates for throwing/searching/faith.
5. **Items** — robes, weapons, holy water, lampions, relics and consumables.
6. **Villager/zombie variants** — living types, matching zombies and conversion.
7. **Tests and tuning** — generation, skill gates, starting kits, intro text,
   conversion and balance.

## 11. Acceptance criteria

- [ ] There is no class-selection screen.
- [ ] The player is always a priest with a randomly generated name.
- [ ] The intro is generated from the assigned perks and skills.
- [ ] Hidden `r` rerolls name, skills, perks, kit and intro without being shown.
- [ ] STR/TGH/AGI/INT are removed from player-facing play.
- [ ] Skills are rated 0–3 and gate throwing, searching and holy abilities.
- [ ] The priest starts with three belt slots.
- [ ] Every starting item fits the priest theme.
- [ ] Living peasants, militia, runners and thugs have matching zombie versions.
- [ ] Zombie kills convert the matching villager type.
- [ ] Holy water and lampions are throwable burn weapons.
- [ ] The game remains readable at 80×25.

## 12. Open questions

- Should skills be fully random, or use a small point-budget with minimums?
- Should perks be binary or also have levels?
- Should hidden `r` also reroll the intro, or only name and skills?
- Should holy water damage all enemies or only undead?
- Should Searching be required for all container loot, or only hidden extras?
- Should the priest be able to use non-priest weapons without a penalty?
- How many living villager types should appear in the first pass?
