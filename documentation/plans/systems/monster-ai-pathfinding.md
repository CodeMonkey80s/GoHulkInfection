# Monster AI and Cardinal Pathfinding

Status: **done**
Parent: [Feature plans](../README.md)
Related: [Cardinal movement](cardinal-movement.md), [Woods level revamp](../woods/level-revamp.md)
Scope: monster target selection, pathfinding, local avoidance and villager/animal behaviour. Movement itself is covered by the cardinal movement plan.

This plan replaces the current first-step BFS with a cardinal pathfinding layer
that behaves well in the woods, around huts, water and crowds.

Implementation status: complete. Monsters use four-way A* with path caching,
hazard/occupancy costs, door opening, local sidesteps and stuck detection.

---

## 1. Current issues

The current AI:

- uses eight-way BFS and first-step-only movement
- ignores fire, acid, water and crowded tiles when choosing a route
- can get stuck behind corners, furniture or other critters
- paths through doors without a clear cost model
- uses simple “step away” logic for fleeing villagers and animals
- recomputes a full BFS for every active monster every turn
- has no stuck detection or local avoidance

The faction/target layer from the woods revamp is good and should be kept.

## 2. Goals

- Make monster movement cardinal and predictable.
- Route around walls, furniture, water and dangerous tiles.
- Keep villagers and animals believable without making them tactical geniuses.
- Avoid monsters piling up or blocking each other forever.
- Keep pathfinding cheap on compact 80×23 woods maps and 80×46 section maps.
- Make behaviour testable and deterministic where possible.

## 3. Pathfinding design

Use **A\*** with four-way movement and a Manhattan heuristic.

Suggested tile costs:

| Tile / state | Cost |
|---|---:|
| open floor / grass / road | 1 |
| bush / reeds | 2 |
| closed door | 3 |
| bramble / fallen log / small furniture | 6 |
| shallow water | 8 |
| fire / acid / smoke | 20 |
| deep water / solid wall / impassable | blocked |
| occupied by another critter | 40 |
| target tile | normal cost |

Rules:

- Paths may pass through a closed door; the critter opens it when stepping in.
- Other critters are high-cost rather than completely blocking, so a route can
  exist around a crowd.
- Fire and acid are strongly avoided but not absolutely forbidden.
- Deep water is blocked unless a future swimming trait is added.
- Use a maximum expansion budget so one bad target cannot stall a turn.
- Cache a path per critter/target for a few turns; invalidate it when the target
  moves, the next step becomes blocked, or the critter is stuck.
- If no path exists, use local cardinal sidesteps, wait, or wander.

## 4. Behaviour by faction

### Zombies

- Keep the small 3-tile aggro radius.
- Chase the last seen position for a limited number of turns.
- Idle most of the time when unaware.
- Use A* when chasing; use simple cardinal wandering otherwise.
- Attack any adjacent hostile target: player, villager or animal.
- Keep delayed villager-to-zombie conversion.

### Villagers

- Flee zombies using a short cardinal path away from the threat.
- Prefer routes that increase distance and avoid dead ends.
- Armed villagers and militia may advance and attack.
- Unarmed villagers only fight when cornered.
- Do not path through zombies or fire.

### Animals

- Rabbits, deer, foxes, crows and frogs flee using cardinal steps.
- Boars remain neutral until attacked, then charge using the hostile path layer.
- Wolves hunt with a larger sight radius.
- Ambient animals may idle instead of moving every turn.

### Observers and other themes

Caves, temple and time-portal AI should keep their current behaviour until their
own plans are written. The pathfinder can be shared, but their target rules
should not change in this pass.

## 5. Local avoidance and stuck detection

Track each critter’s last position and last progress turn.

- If the next path step is occupied, try a cardinal sidestep that still makes
  progress.
- If no progress is made for 2–3 turns, invalidate the cached path.
- If still stuck, wait one turn or wander.
- Prevent two monsters from swapping places indefinitely.
- Do not let fleeing critters oscillate between the same two tiles.

## 6. Performance

- 80×23 woods / 80×46 section maps, four-way A*, small expansion cap.
- Cache paths for a few turns.
- Only run full pathfinding for critters that are alert or have a target.
- Idle/wandering critters use a cheap cardinal step.
- Add a benchmark or timing test for many active monsters.

## 7. Testing

Add deterministic tests for:

- cardinal-only monster movement
- path around a wall or hut corner
- path through a door and opening it
- avoiding fire/acid when a safe route exists
- not pathing through deep water
- routing around allied critters
- stuck detection and recovery
- zombie chasing a villager/player
- villager fleeing a zombie
- neutral animal fleeing
- boar charging after being attacked
- villagers still converting after a zombie kill
- performance with crowded compact woods and section maps

## 8. Implementation phases

1. **Path primitives** — add `dirs4`, A*, cost function and path cache.
2. **Zombie/wolf movement** — replace BFS-first-step with A* chase.
3. **Villager/animals** — cardinal flee/charge paths.
4. **Local avoidance** — occupancy checks, stuck detection, sidesteps.
5. **Generators/tests** — ensure cardinal connectivity and add AI tests.
6. **Tuning** — aggro ranges, costs and idle rates after playtesting.

## 9. Acceptance criteria

- [x] Monsters move only cardinally.
- [x] Monsters route around walls, furniture and deep water.
- [x] Monsters strongly avoid fire and acid.
- [x] Monsters open doors instead of stopping at them.
- [x] Monsters do not permanently block or oscillate around each other.
- [x] Villagers flee zombies and fight when armed/cornered.
- [x] Animals flee; boars charge when attacked.
- [x] Zombies still convert killed villagers.
- [x] Pathfinding remains fast on crowded woods and section maps.
- [x] Caves, temple and time-portal AI remains unchanged.

## 10. Resolved decisions

- Monsters avoid fire/acid/smoke with a high cost but may cross if no safe route exists.
- Allies are high-cost obstacles; monsters path around them rather than swapping.
- Cached paths last three turns or until the target moves / the next step is blocked.
- Villagers use greedy cardinal fleeing; armed villagers use A* when advancing.
- Zombies use sight only for now; noise/hearing remains future work.
