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
	tileDirtRoad    uint8
	tileBush        uint8
	tileBramble     uint8
	tileFlowers     uint8
	tileStump       uint8
	tileFallenLog   uint8
	tileDeadTree    uint8
	tileLanternPost uint8
	tileMud         uint8
	tileReeds       uint8
	tileTable       uint8
	tileChair       uint8
	tileBed         uint8
	tileCrate       uint8
	tileBarrel      uint8
	tileFireplace   uint8
	tileShelf       uint8
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

	woodsTileRoles = map[int]ColorRole{
		int(tileDirtRoad):    RoleDirtRoad,
		int(tileBush):        RoleBush,
		int(tileBramble):     RoleBramble,
		int(tileFlowers):     RoleFlowerGold,
		int(tileStump):       RoleStump,
		int(tileFallenLog):   RoleStump,
		int(tileDeadTree):    RoleDeadTree,
		int(tileLanternPost): RoleLanternGlow,
		int(tileMud):         RoleMud,
		int(tileReeds):       RoleReeds,
		int(tileTable):       RoleFurniture,
		int(tileChair):       RoleFurniture,
		int(tileBed):         RoleCloth,
		int(tileCrate):       RoleFurniture,
		int(tileBarrel):      RoleFurniture,
		int(tileFireplace):   RoleFirelight,
		int(tileShelf):       RoleFurniture,
	}
}
