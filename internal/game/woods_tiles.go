package game

// Woods terrain tiles are appended at the end of TileDefs so all recovered tile
// IDs remain stable for saves and older data.
const (
	tileWoodenFloor  = 13
	tileWoodenWall   = 14
	tileWoodenDoor   = 15
	tileWoodenWindow = 17
)

var (
	tileDirtRoad     uint8
	tileBush         uint8
	tileBramble      uint8
	tileFlowers      uint8
	tileStump        uint8
	tileFallenLog    uint8
	tileDeadTree     uint8
	tileLanternPost  uint8
	tileMud          uint8
	tileReeds        uint8
	tileTable        uint8
	tileChair        uint8
	tileBed          uint8
	tileCrate        uint8
	tileBarrel       uint8
	tileFireplace    uint8
	tileShelf        uint8
	tileRiverBank    uint8
	tileWoodenBridge uint8
	tileStoneBridge  uint8
	tileLilyPad      uint8
	tileOrchardTree  uint8
	tileBeehive      uint8
	tileCiderPress   uint8
	tileGravestone   uint8
	tileIronFence    uint8
	tileCryptDoor    uint8
	tileMillStone    uint8
	tileTent         uint8
	tileCampfire     uint8
	tileCage         uint8
	tileDryingRack   uint8
	tileMushroomRing uint8
	tileDruidAltar   uint8
	tileWell         uint8
	tileCart         uint8
	tileMarketStall  uint8
	tileHaystack     uint8
	tileKiln         uint8
	tileChapelAltar  uint8
	tileBellTower    uint8
	tileWoodpile     uint8
)

var woodsTileRoles = map[int]ColorRole{}

func init() {
	add := func(d TileDef) uint8 {
		id := uint8(len(TileDefs))
		TileDefs = append(TileDefs, d)
		return id
	}

	tileDirtRoad = add(TileDef{
		Name: "dirt road", Type: 1, Glyph: '.', Color: 6,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileBush = add(TileDef{
		Name: "bush", Type: 1, Glyph: '"', Color: 2,
		BlockMove: false, BlockSight: true, Partner: -1,
	})
	tileBramble = add(TileDef{
		Name: "bramble", Type: 1, Glyph: '#', Color: 6,
		BlockMove: true, BlockSight: true, Partner: -1,
	})
	tileFlowers = add(TileDef{
		Name: "flowers", Type: 1, Glyph: '*', Color: 14,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileStump = add(TileDef{
		Name: "stump", Type: 1, Glyph: 'o', Color: 6,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileFallenLog = add(TileDef{
		Name: "fallen log", Type: 1, Glyph: '=', Color: 6,
		BlockMove: true, BlockSight: true, Partner: -1,
	})
	tileDeadTree = add(TileDef{
		Name: "dead tree", Type: 1, Glyph: 't', Color: 8,
		BlockMove: true, BlockSight: true, Partner: -1,
	})
	tileLanternPost = add(TileDef{
		Name: "lantern post", Type: 2, Glyph: 'i', Color: 14,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileMud = add(TileDef{
		Name: "mud", Type: 1, Glyph: '~', Color: 6,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileReeds = add(TileDef{
		Name: "reeds", Type: 1, Glyph: '"', Color: 6,
		BlockMove: false, BlockSight: true, Partner: -1,
	})
	tileTable = add(TileDef{
		Name: "table", Type: 2, Glyph: 'T', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileChair = add(TileDef{
		Name: "chair", Type: 1, Glyph: 'h', Color: 6,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileBed = add(TileDef{
		Name: "bed", Type: 2, Glyph: 'b', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileCrate = add(TileDef{
		Name: "crate", Type: 2, Glyph: 'x', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileBarrel = add(TileDef{
		Name: "barrel", Type: 2, Glyph: '0', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileFireplace = add(TileDef{
		Name: "fireplace", Type: 2, Glyph: '&', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileShelf = add(TileDef{
		Name: "shelf", Type: 2, Glyph: '[', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileRiverBank = add(TileDef{
		Name: "river bank", Type: 1, Glyph: ',', Color: 6,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileWoodenBridge = add(TileDef{
		Name: "wooden bridge", Type: 1, Glyph: '=', Color: 6,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileStoneBridge = add(TileDef{
		Name: "stone bridge", Type: 1, Glyph: '=', Color: 7,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileLilyPad = add(TileDef{
		Name: "lily pad", Type: 1, Glyph: 'o', Color: 10,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileOrchardTree = add(TileDef{
		Name: "orchard tree", Type: 1, Glyph: 't', Color: 10,
		BlockMove: true, BlockSight: true, Partner: -1,
	})
	tileBeehive = add(TileDef{
		Name: "beehive", Type: 2, Glyph: 'B', Color: 14,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileCiderPress = add(TileDef{
		Name: "cider press", Type: 2, Glyph: 'P', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileGravestone = add(TileDef{
		Name: "gravestone", Type: 2, Glyph: '+', Color: 7,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileIronFence = add(TileDef{
		Name: "iron fence", Type: 2, Glyph: '|', Color: 8,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileCryptDoor = add(TileDef{
		Name: "crypt door", Type: 2, Glyph: '+', Color: 8,
		BlockMove: true, BlockSight: true, Partner: -1,
	})
	tileMillStone = add(TileDef{
		Name: "mill stone", Type: 2, Glyph: 'O', Color: 7,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileTent = add(TileDef{
		Name: "tent", Type: 2, Glyph: '^', Color: 6,
		BlockMove: true, BlockSight: true, Partner: -1,
	})
	tileCampfire = add(TileDef{
		Name: "campfire", Type: 2, Glyph: '&', Color: 12,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileCage = add(TileDef{
		Name: "cage", Type: 2, Glyph: 'X', Color: 8,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileDryingRack = add(TileDef{
		Name: "drying rack", Type: 2, Glyph: '=', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileMushroomRing = add(TileDef{
		Name: "mushroom ring", Type: 1, Glyph: '"', Color: 13,
		BlockMove: false, BlockSight: false, Partner: -1,
	})
	tileDruidAltar = add(TileDef{
		Name: "druid altar", Type: 2, Glyph: 'A', Color: 7,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileWell = add(TileDef{
		Name: "well", Type: 2, Glyph: 'W', Color: 7,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileCart = add(TileDef{
		Name: "cart", Type: 2, Glyph: 'D', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileMarketStall = add(TileDef{
		Name: "market stall", Type: 2, Glyph: 'M', Color: 4,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileHaystack = add(TileDef{
		Name: "haystack", Type: 2, Glyph: 'H', Color: 14,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileKiln = add(TileDef{
		Name: "kiln", Type: 2, Glyph: 'K', Color: 8,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileChapelAltar = add(TileDef{
		Name: "chapel altar", Type: 2, Glyph: 'A', Color: 14,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileBellTower = add(TileDef{
		Name: "bell tower", Type: 2, Glyph: 'B', Color: 7,
		BlockMove: true, BlockSight: false, Partner: -1,
	})
	tileWoodpile = add(TileDef{
		Name: "woodpile", Type: 2, Glyph: '=', Color: 6,
		BlockMove: true, BlockSight: false, Partner: -1,
	})

	woodsTileRoles = map[int]ColorRole{
		int(tileDirtRoad):     RoleDirtRoad,
		int(tileBush):         RoleBush,
		int(tileBramble):      RoleBramble,
		int(tileFlowers):      RoleFlowerGold,
		int(tileStump):        RoleStump,
		int(tileFallenLog):    RoleStump,
		int(tileDeadTree):     RoleDeadTree,
		int(tileLanternPost):  RoleLanternGlow,
		int(tileMud):          RoleMud,
		int(tileReeds):        RoleReeds,
		int(tileTable):        RoleFurniture,
		int(tileChair):        RoleFurniture,
		int(tileBed):          RoleCloth,
		int(tileCrate):        RoleFurniture,
		int(tileBarrel):       RoleFurniture,
		int(tileFireplace):    RoleFirelight,
		int(tileShelf):        RoleFurniture,
		int(tileRiverBank):    RoleMud,
		int(tileWoodenBridge): RoleBridgeWood,
		int(tileStoneBridge):  RoleStoneLight,
		int(tileLilyPad):      RoleFrog,
		int(tileOrchardTree):  RoleOrchardLeaf,
		int(tileBeehive):      RoleHay,
		int(tileCiderPress):   RoleFurniture,
		int(tileGravestone):   RoleGravestone,
		int(tileIronFence):    RoleIronFence,
		int(tileCryptDoor):    RoleCrypt,
		int(tileMillStone):    RoleStoneLight,
		int(tileTent):         RoleCloth,
		int(tileCampfire):     RoleFirelight,
		int(tileCage):         RoleIronDark,
		int(tileDryingRack):   RoleWoodWall,
		int(tileMushroomRing): RolePotionPurple,
		int(tileDruidAltar):   RoleAltarStone,
		int(tileWell):         RoleStoneWall,
		int(tileCart):         RoleWoodLight,
		int(tileMarketStall):  RoleMarketCloth,
		int(tileHaystack):     RoleHay,
		int(tileKiln):         RoleStoneWall,
		int(tileChapelAltar):  RoleAltarStone,
		int(tileBellTower):    RoleStoneLight,
		int(tileWoodpile):     RoleStump,
	}
}
