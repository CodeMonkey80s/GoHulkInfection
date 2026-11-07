# Woods Colour Palette — True-Colour Proposal

Status: **done**
Parent: [Woods plans](README.md)  
Related: [Woods level revamp](level-revamp.md), [Minimalist snake redesign](minimalist-snake-redesign.md)  
Scope: woods rendering and UI palette only

This document proposes a 24-bit colour palette for the woods theme. The core
palette, true-colour/16-colour fallback, per-theme accents, warm lighting,
blood blending and true-colour UI are implemented.

The palette was originally tuned for the nine-level woods. The completed
[Minimalist snake redesign](minimalist-snake-redesign.md) now applies one accent
pair to each of the five themed 80×23 maps; the rendering rules and palette
roles are unchanged.

---

## 1. Current colour review

The port currently inherits the original 16-colour DOS palette. It works, but
has limits:

- woods are dominated by one or two saturated greens
- grass, trees, bushes and moss all read as “green dots”
- wooden structures are almost the same brown as roads and dirt
- water is a flat blue/cyan
- blood is a solid red background block
- items are coloured mostly by category, not material
- no lighting gradient, weather tint, or per-level variation
- danger and loot do not stand out strongly enough

The original palette is useful as a fallback and for preserving the old look,
but the woods can be much more readable and atmospheric with true colour.

## 2. Design direction

**Muted base, vivid punctuation.**

Do not make everything vivid: blood, fire, acid, loot and enemies stop standing
out when the whole screen is saturated.

Do not make everything muddy either: too much desaturation makes zombies,
corpses, items and blood hard to read.

Use:

- **environment:** muted, earthy, slightly desaturated, with subtle per-tile variation
- **danger:** high-saturation blood, fire, acid and poison
- **light:** warm lantern/candle pools against cool shadows
- **loot and interactables:** brighter than the background, but not neon
- **UI:** parchment/ink tones with blood-red and antique-gold accents

This gives a grim woodland atmosphere while keeping combat readable.

## 3. Core woods palette

### 3.1 Terrain and vegetation

| Role | Hex | Notes |
|---|---|---|
| Map background / unexplored | `#050605` | near-black, not pure black |
| Grass base | `#4E5F3A` | muted moss green |
| Grass light | `#5C7044` | sunlit patches |
| Grass dark | `#3C4A2E` | shade under trees |
| Leaf litter | `#5A4A32` | forest floor |
| Dirt road | `#6E5A42` | warm earth |
| Road light | `#7E6A50` | worn centre |
| Mud | `#4A3C2E` | wet ground |
| Rock | `#5C5E56` | cool grey-green |
| Rock light | `#7A7C72` | exposed stone |
| Pine canopy | `#2D4029` | dark treetops |
| Canopy mid | `#3E5734` | main tree colour |
| Canopy light | `#577345` | highlighted leaves |
| Bush | `#556B3E` | passable bush |
| Bramble | `#4A3A2A` | blocking brush |
| Thorn | `#8C7A5A` | bramble detail |
| Dead tree | `#4A4038` | cemetery / burned woods |
| Orchard leaf | `#5E7A3E` | fruit-tree rows |
| Stump | `#5A4230` | chopped wood |
| Flowers | `#C9B45A`, `#B87A8A`, `#E0D8B0` | small accents only |

### 3.2 Water, roads and light

| Role | Hex |
|---|---|
| Deep water | `#1C3440` |
| River | `#2F5A66` |
| Shallow water | `#4E7C7A` |
| Foam / bridge edge | `#A6C4BC` |
| Marsh | `#3E4A34` |
| Reeds | `#7A7A42` |
| Bridge wood | `#6E4F35` |
| Bridge rope/rail | `#8C7A5A` |
| Lantern metal | `#4A4A44` |
| Lantern glow | `#FFD27A` |
| Firelight | `#E08A3A` |
| Moonlight / shade | `#223028` |

### 3.3 Structures, huts, cemetery and village

| Role | Hex |
|---|---|
| Wood wall dark | `#4A3324` |
| Wood wall | `#6E4F35` |
| Wood light / trim | `#8A6A48` |
| Roof | `#5A3E2A` |
| Door | `#8A5A32` |
| Window frame | `#6E4F35` |
| Warm window glass | `#D9A441` |
| Fence | `#6A5238` |
| Table / furniture | `#7A5A3A` |
| Stone wall | `#6E6E66` |
| Stone light | `#8A8A80` |
| Gravestone | `#8C8C84` |
| Grave dirt | `#4A3E32` |
| Crypt / mausoleum | `#4E4E4A` |
| Iron fence | `#3E4648` |
| Bone | `#D6CBB0` |
| Grave fog | `#8A9A92` |
| Plaster wall | `#B8A88A` |
| Thatch | `#A88A52` |
| Faded paint | `#7A6A5A` |
| Market cloth | `#A84A3A` |
| Hay | `#C9A85A` |
| Altar stone | `#7A7A70` |
| Candle | `#F2C46A` |
| Stained glass | `#7A4A8A`, `#3A6A8A` |

### 3.4 Blood, gore and undead

| Role | Hex |
|---|---|
| Fresh blood | `#B32020` |
| Dark blood | `#7A1A1A` |
| Dried blood | `#5A2420` |
| Blood blended into grass | `#6B2A22` |
| Rotten flesh | `#7A8A5A` |
| Zombie skin | `#7F8F62` |
| Zombie shadow | `#4A553A` |
| Bone / teeth | `#D6CBB0` |
| Guts / open wound | `#8A4A4A` |
| Poison | `#7A8A3A` |
| Acid | `#8FBF3F` |
| Acid dark | `#5A7A2A` |
| Fire ember | `#E06A2A` |
| Flame | `#F2A33A` |
| Flame core | `#FFD27A` |
| Smoke | `#5A5A55` |
| Ash | `#7A7A72` |
| Charcoal | `#222222` |

### 3.5 Creatures

| Role | Hex |
|---|---|
| Player | `#F2E6C9` |
| Player highlight | `#FFD27A` |
| Villager | `#C9B48A` |
| Villager clothing | `#6A7A8A`, `#8A5A3A`, `#7A6A5A` |
| Deer | `#8A6A4A` |
| Rabbit | `#B8A88A` |
| Fox | `#B55A2A` |
| Boar | `#5A4A3A` |
| Wolf | `#6A6A66` |
| Crow | `#1E1E22` |
| Frog | `#5A8A4A` |
| Hostile animal eyes | `#C9A227` |

### 3.6 Items

| Role | Hex |
|---|---|
| Steel | `#B8C0C0` |
| Iron dark | `#6A7070` |
| Silver | `#C0C8C8` |
| Gold / relic | `#C9A227` |
| Leather | `#7A4E2E` |
| Cloth | `#8A7A5A` |
| Apple | `#A83A32` |
| Bread | `#C9A26A` |
| Cheese | `#D9B85A` |
| Meat | `#8A3A2A` |
| Fish | `#8AA8B0` |
| Potion red | `#A83A32` |
| Potion green | `#4A8A3A` |
| Potion blue | `#3A6A8A` |
| Potion purple | `#7A4A8A` |
| Parchment | `#D8C9A3` |
| Key | `#C9A227` |

### 3.7 UI

| Role | Hex |
|---|---|
| UI background | `#0B0D0A` |
| Panel background | `#141812` |
| Panel border | `#3E4A38` |
| Primary text | `#D6CBB0` |
| Dim text | `#8A8F7A` |
| Accent green | `#7A9A5A` |
| Accent gold | `#C9A227` |
| Danger / damage | `#B23A2E` |
| Blood accent | `#B32020` |
| Selection background | `#4A5A3A` |
| Selection text | `#F2E6C9` |
| Health good | `#6A9A4A` |
| Health wounded | `#C9A227` |
| Health critical | `#B23A2E` |
| Hint / cyan | `#6A9A9A` |

## 4. Per-theme accent palettes (five-map run)

Keep the core woods palette, then tint a small portion of the map. Accents
should affect 5–15% of the map, not the whole screen.

| Theme | Accent palette |
|---|---|
| Church | relic gold `#C9A227`, stained-glass blue `#3A6A8A` |
| Peasant Houses | thatch `#C9A85A`, timber `#8A5A3A` |
| Mill and Pond | river teal `#2F5A66`, mill gold `#C9A85A` |
| Forgotten Cemetery | bone `#D6CBB0`, grave fog `#8A9A92` |
| Meteor Field | meteor glow `#8AF2B0`, ember `#E06A2A` |

## 5. Rendering rules

- Use true colour when available; keep a 16-colour fallback.
- Add slight per-tile variation:
  - hash `(level, x, y)` into ±5–8% brightness
  - use 2–4 grass/canopy variants
- Shade remembered/unseen tiles instead of turning them grey:
  - visible: full colour
  - remembered: ~40% brightness, slightly desaturated
  - shadow/unseen: `#050605` or very dark desaturated
- Blend blood with the underlying tile instead of using a solid background block.
- Lanterns and fires create warm local light pools.
- Reserve high saturation for blood, fire, acid, loot and UI warnings.
- Keep glyphs distinct as a fallback for colour-blind players and monochrome terminals.

## 6. Implementation plan

- [x] Add an RGB palette type and a named woods palette.
- [x] Extend tile/item/critter colour data with RGB values and a fallback palette index.
- [x] Detect true-colour support (`COLORTERM` / terminal capability) and fall back to 256 or 16 colours.
- [x] Add deterministic per-tile colour variation.
- [x] Replace flat visibility shading with brightness/desaturation grading.
- [x] Blend blood and effects with tile colours instead of drawing solid background blocks.
- [x] Add warm light radii for lanterns, candles and fires.
- [x] Define the per-theme accent sets.
- [x] Apply the new UI palette.
- [x] Add automated render smoke checks for all five woods maps in true colour and 16-colour modes.

## 7. Acceptance criteria

- [x] All woods levels remain readable in 16-colour fallback mode.
- [x] Grass, trees, bushes, roads and wooden structures are visually distinct.
- [x] Blood, corpses, items, zombies and villagers remain easy to see on every terrain type.
- [x] All five woods maps are recognisable from their accent colour and landmark.
- [x] The UI remains legible and does not compete with the map.
- [x] Performance remains acceptable on the compact 80×23 woods maps.

## 8. Open questions

- Should lanterns provide real light/FOV radius immediately, or only warm colour?
- Should weather tint the whole palette (rain desaturation, fog wash, night blue)?
- How strong should per-tile colour variation be before it becomes noisy?
- Should the 16-colour fallback preserve the original DOS colours exactly?
- Should acid/fire colours pulse, or remain static for readability?
- Should true colour be selectable in the options screen?
