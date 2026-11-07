# Item Theme Review

Status: **done**
Parent: [Documentation index](../README.md)
Related: [Reverse-engineering reference](reverse-engineering.md), [Priest player revamp](../plans/systems/priest-player-revamp.md)

The recovered item table contains modern military and science-fiction
prototypes. Recovered IDs are kept stable for data and save references, but
only thematic items may appear in normal play.

## Theme rules

- Village, woodland, medieval, priest, holy and horror items are allowed.
- Guns, ammunition, clips, energy weapons and old body armour are removed from
  normal play.
- Modern grenades are refitted as priest flasks.
- Raw removed prototypes remain in `ItemDefs` for historical reference but fail
  `ItemAvailable`.

## Renamed and refitted items

| Old recovered name | Current name | Role |
|---|---|---|
| jacket | padded jerkin | legacy armour data |
| flak armor | reinforced hide coat | legacy armour data |
| power jacket | penitent's mail | legacy armour data |
| plasteel stick | iron-shod staff | melee |
| chainsaw | relic chainblade | melee relic |
| power sword | silvered longsword | melee |
| meta dagger | consecrated dagger | melee |
| vibra knife | blessed knife | melee |
| multiaxe | penitent's axe | melee |
| shock mace | censer mace | melee |
| battle hammer | war hammer | melee |
| he grenade | holy water flask | throwable flask |
| incendiary grenade | holy fire flask | throwable flask |
| acid grenade | alchemical acid flask | throwable flask |
| smoke grenade | incense flask | throwable flask |
| flashbang grenade | radiant flask | throwable flask |
| first aid kit | healer's kit | healing |
| medkit | surgeon's kit | healing |

## Removed from normal play

- all recovered ranged weapons and energy weapons
- all ammunition and clips
- body armour and old armour pieces
- mobile-phone-era military equipment that has no thematic replacement

The priest uses the permanent robe and attachable shards instead of body armour.
The first-playable item families are woodland food/tools/weapons, priest
vials/shards/relics, holy flasks and ordinary medical supplies.

## Validation

`internal/game/item_theme_test.go` checks that:

- old out-of-theme names are gone from the active item set
- renamed thematic replacements exist
- available items do not contain banned terms such as grenade, gun, ammo, clip,
  plasma, laser, flak, plasteel or flashbang
- ranged weapons, ammunition and old body armour are unavailable
