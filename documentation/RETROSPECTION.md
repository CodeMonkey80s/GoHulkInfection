# Retrospection — How the Project Drifted from Its Original Goal

Status: historical record
Related: [README](../README.md), [Documentation index](README.md),
[Reverse-engineering reference](reference/reverse-engineering.md)

This document is a retrospection over the git history and the repository
documentation. Its purpose is to record, honestly, how the project drifted away
from its initial goal: **to reverse engineer `hi030.exe` into Go source code.**

---

## 1. The original goal

The project started as a single commit, `911691b` "Initial Go port of hulk
infection" (2026-11-06). The initial `README.md` described it as:

> A from-scratch Go reimplementation of **hulk infection v0.3a — ultimate
> chainsaw massacre** (2007, Michał Przybyłowicz), reconstructed by
> reverse-engineering the original `hi030.exe`.

The first commit contained a complete, working Go port: ~6,300 lines across 18
files, embedding the data tables recovered from the binary — item prototypes,
critter prototypes, the skill list, tile prototypes, the four 10×40 map
sections, key bindings, options, player classes, hints and original texts.

The intent was **preservation through reimplementation**. The original game had:

- 8-directional movement via the numpad (62 recovered key bindings)
- ranged weapons with ammunition, clips, reloading and burst fire
- grenades of five types (HE, incendiary, acid, smoke, flashbang)
- three player professions (murderer, marine, warpcrafter) with recovered
  STR/TGH/AGI/INT bonuses
- nine skills, skill levels, experience, character levels and free skill points
- weight, burden, body armour with percentage absorption
- tactics (defensive / normal / aggressive), search, operate, clean, target
  selection and a target line
- 80×80 maps assembled from the four recovered cube sections

All of this was reconstructed or recovered in the initial commit, with the
approximations documented in `REVERSE_ENGINEERING.md`.

## 2. The drift, phase by phase

96 commits were made, all within 36 hours (2026-11-06 to 2026-11-07), all by
one author. The reverse engineering itself was completed *before* the
repository started; every subsequent commit was game design, not recovery.

**Phase 1 — command scope reduction (PR #1).** The
[command-scope-reduction plan](plans/systems/command-scope-reduction.md)
removed ranged weapons, ammunition, reloading, unarmed combat, aim mode, target
selection, the target line, weapon activation, weapon swapping, tactics,
cleaning, search and the operate command. This was the first step away from
fidelity: the recovered systems were now treated as material to cut, not
material to preserve.

**Phase 2 — player redesign (PRs #2–#6).** The
[player-redesign plan](plans/systems/player-redesign.md) replaced the recovered
weight/burden/armour model with five ADOM-style equipment slots and an
invented HP/AP (armour points) model. A minimal HUD replaced the recovered
status displays.

**Phase 3 — woods content (PRs #7–#18).** The largest block of work. The
original woods levels were assembled from the four recovered 10×40 sections;
the port replaced that with a fully original terrain generator, themed
landmarks, road-facing huts, rivers and ponds, signature loot, villager/animal
factions and zombie conversion, plus a 24-bit true-colour palette. None of
this was recovered from the binary — it is new creative content.

**Phase 4 — controls and AI redesign (PRs #19–#23).** The
[cardinal-movement plan](plans/systems/cardinal-movement.md) removed diagonal
movement entirely and replaced the recovered numpad bindings with `h j k l`.
The recovered 8-directional model of the original was deliberately abandoned
for a new four-directional design.

**Phase 5 — compact snake redesign (PRs #24–#26).** The
[minimalist-snake-redesign plan](plans/woods/minimalist-snake-redesign.md)
replaced the run structure with five invented 80×23 themed maps (Church,
Peasant Houses, Mill and Pond, Forgotten Cemetery, Meteor Field) connected by
a random "snake" route, and shrank the terminal requirement from 80×50 to
80×25. The caves, temple and time-portal maps — the parts of the original
game still using the recovered section data — were explicitly postponed.

**Phase 6 — priest revamp (PRs #28–#33, in progress).** The
[priest-player-revamp plan](plans/systems/priest-player-revamp.md) removed the
three recovered professions and replaced them with a single invented priest;
removed the recovered stats, skills, perks and player levels; and added
invented mechanics (random Rites, permanently worn Robe, armour shards,
one charge per map). The
[item theme review](reference/item-review.md) renamed recovered items — the
chainsaw became a "relic chainblade", the HE grenade a "holy water flask" —
and removed all ranged weapons, ammunition and body armour from normal play.
A later commit made the rendering "Brogue-like", an explicit reference to a
different game.

## 3. What drifted

| Aspect | Initial port (recovered) | Current state |
|---|---|---|
| Movement | 8-directional numpad | cardinal `hjkl` only |
| Combat | melee + ranged + ammo + reload + 5 grenade types | melee-only; throwing deferred |
| Player characters | 3 recovered professions | 1 invented priest |
| Statistics | recovered STR/TGH/AGI/INT | removed from play |
| Progression | skills, XP, levels, free points | XP stored but unused; levels removed |
| Equipment | inventory, weight, burden, body armour | Hands/Head/Robe slots, shard AP |
| Levels | 80×80 from recovered 10×40 sections | invented 80×23 generator (woods); recovered sections only in postponed cave/temple/lab maps |
| Items | 49 recovered prototypes | ranged/ammo/armour removed; ~20 renamed to a priest theme |
| Abilities | recovered skills and options | invented Rites (7 starting, shrine-granted) |
| Presentation | DOS/PDCurses colours | 24-bit true-colour, "Brogue-like" rendering |
| Genre | sci-fi/military horror | medieval/religious village horror |

Notable: of the 96 commit messages, **none mention fidelity, recovery,
preservation or reverse engineering**. There are 12 planning commits and 49
feature commits, all about new design. The project moved from *recreating the
original* to *using the recovered data tables as raw material for an original
game*.

## 4. Why it drifted — the pattern

The drift was **not** accidental scope creep. It was a sequence of deliberate,
documented redesign decisions. Each phase has a plan document, and each plan
explicitly cuts or replaces recovered systems:

1. **A design-driven process.** The workflow (plan → decision → implement →
   merge) treated the original game as a constraint to be removed, not a
   specification to be met. Every plan asks "what fits our theme" rather than
   "what does the original do".
2. **A thematic genre shift.** "Melee-focused zombie slashing" and later the
   priest/village theme conflicted with the original's military sci-fi content,
   so the recovered content that did not fit was cut or renamed.
3. **The ratchet of scope reduction.** Each plan removed systems "for now"
   (throwing, unarmed, ranged, states, caves, temple, end-game). Nothing
   removed has come back; the backlog keeps growing.
4. **Postponement as drift.** Caves, temple and the time portal — the remaining
   original content — have been "not planned yet" since the woods work began,
   while new original systems keep being added.
5. **The README lagged behind.** The root README still introduces the project
   as a faithful port "reconstructed by reverse-engineering the original
   `hi030.exe`", while most gameplay systems described beneath that sentence
   are new designs. Only the fidelity notes at the bottom admit the
   reconstruction is partial.

## 5. What remains faithful

The drift is real but not total. The following recovery work survives and is
still the foundation of the game:

- `internal/game/data_gen.go` still embeds the recovered item prototypes, the
  20-entry weapon AV/DV table, critter prototypes, all 64 tile prototypes, the
  four 10×40 map sections, skills, key bindings, options, hints and original
  texts.
- The weapon table is authoritative: attack/defence values and damage ranges
  still come from the binary.
- The 50 recovered options and their min/max values remain embedded.
- The caves, temple and time-portal maps still use the recovered section data
  (though they are postponed).
- The original copyright notices and texts are preserved.

## 6. Conclusion

The project began as a reverse-engineering port of `hi030.exe` and, within two
days, became the development of an original minimalist melee roguelike that
borrows the original's data tables. The recovered data remains in the code,
but the game's systems — movement, combat, equipment, progression, character,
levels, items and presentation — are now largely original designs, and the
themes that remain faithful are exactly the ones that have not been redesigned
yet.

This is a legitimate way for a project to evolve, but the repository should
be honest about it: the current game is **not** a faithful reimplementation of
hulk infection v0.3a. It is a new game built on recovered data.
