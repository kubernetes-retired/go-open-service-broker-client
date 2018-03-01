package generator

// generator holds the parameters for generated responses.
type Generator struct {
	Services []Service
}

type Service struct {
	Plans []Plan
	Tags  int
}

type Plan struct {
}

// All dragon names from A Song of Ice and Fire series by George R.R. Martin
var ClassNames = []string{
	"Archonei",
	"Arrax",
	"Balerion",
	"Caraxes",
	"Dreamfyre",
	"Drogon",
	"Essovius",
	"Ghiscar",
	"Meleys",
	"Meraxes",
	"Morghul",
	"Rhaegal",
	"Seasmoke",
	"Sheepstealer",
	"Shrykos",
	"Silverwing",
	"Stormcloud",
	"Sunfyre",
	"Syrax",
	"Tyraxes",
	"Valryon",
	"Vermax",
	"Vermithor",
	"Vermithrax",
	"Vhagar",
	"Viserion",
}

// All ship names from A Song of Ice and Fire series by George R.R. Martin
var PlanNames = []string{
	"BlackWind",
	"BraveJoffrey",
	"Dagger",
	"DagonsFeast",
	"Esgred",
	"Fingerdancer",
	"Foamdrinker",
	"ForlornHope",
	"Fury",
	"GoldenRose",
	"GoldenStorm",
	"GreatKraken",
	"GreyGhost",
	"Grief",
	"Hardhand",
	"IronLady",
	"IronVengeance",
	"IronVictory",
	"IronWind",
	"IronWing",
	"KingRobertsHammer",
	"Kite",
	"KrakensKiss",
	"LadyJoanna",
	"LadyLyanna",
	"LadyOlenna",
	"Lamentation",
	"Leviathan",
	"Lioness",
	"Lionstar",
	"LordDagon",
	"LordQuellon",
	"LordRenly",
	"LordTywin",
	"LordVickon",
	"MaidensBane",
	"Nightflyer",
	"PrincessMarcella",
	"QueenMargaery",
	"ReapersWind",
	"RedJester",
	"RedTide",
	"SaltyWench",
	"SeaBitch",
	"SeaSong",
	"Seaswift",
	"SevenSkulls",
	"Shark",
	"Silence",
	"Silverfin",
	"Sparrowhawk",
	"SweetCersei",
	"Swiftin",
	"ThrallsBane",
	"Thunderer",
	"Warhammer",
	"WarriorWench",
	"WhiteWidow",
	"Woe",
}

// All castle names from A Song of Ice and Fire series by George R.R. Martin
var TagNames = []string{
	"AcornHall",
	"Antlers",
	"Ashemark",
	"Ashford",
	"Bandallon",
	"TheBanefort",
	"Bitterbridge",
	"Blackcrown",
	"Blackhaven",
	"Blackmont",
	"BloodyGate",
	"BrightwaterKeep",
	"Bronzegate",
	"Castamere",
	"CasterlyRock",
	"CastleBlack",
	"CastleCerwyn",
	"CastleStokeworth",
	"CiderHall",
	"TheCitadel",
	"CleganesKeep",
	"Coldwater",
	"TheCrag",
	"Crakehall",
	"CrowsNest",
	"DeepDen",
	"DeepLake",
	"DeepwoodMotte",
	"Dragonstone",
	"TheDreadfort",
	"Eastwatch-by-the-Sea",
	"EvenfallHall",
	"TheEyrie",
	"Faircastle",
	"Feastfires",
	"Felwood",
	"FlintsFinger",
	"GhostHill",
	"Godsgrace",
	"GoldenTooth",
	"Goldengrove",
	"GrassyVale",
	"Greyguard",
	"GreywaterWatch",
	"GriffinsRoost",
	"Hammerhorn",
	"Harrenhal",
	"HaystackHall",
	"HeartsHome",
	"Hellholt",
	"Highgarden",
	"Highpoint",
	"Honeyholt",
	"HornHill",
	"Hornvale",
	"Hornwood",
	"Ironoaks",
	"Ironrath",
	"Karhold",
	"Kingsgrave",
	"LastHearth",
	"Lemonwood",
	"LongBarrow",
	"LongTable",
	"LongbowHall",
	"Mistwood",
	"MoatCailin",
	"TheNightfort",
	"Nightsong",
	"OldOak",
	"Oldcastle",
	"PalaceofJustice",
	"Pinkmaiden",
	"Pyke",
	"Queensgate",
	"RainHouse",
	"Ramsgate",
	"RaventreeHall",
	"RedKeep",
	"RedLake",
	"TheRedfort",
	"RillwaterCrossing",
	"Riverrun",
	"Rosby",
	"Runestone",
	"Saltshore",
	"Sandstone",
	"Sarsfield",
	"Seagard",
	"SealordsPalace",
	"TheShadowTower",
	"SharpPoint",
	"Silverhill",
	"Skyreach",
	"Starfall",
	"StoneHedge",
	"Stonedance",
	"Stonehelm",
	"StormsEnd",
	"Summerhall",
	"SunflowerHall",
	"Sunspear",
	"TarbeckHall",
	"TenTowers",
	"ThreeTowers",
	"TheTor",
	"TorrhensSquare",
	"Tumbleton",
	"TheTwins",
	"UnnamedBaelishcastle",
	"Uplands",
	"Vaith",
	"VulturesRoost",
	"TheWhispers",
	"Whitewalls",
	"WidowsWatch",
	"Winterfell",
	"Wyl",
	"Yronwood",
}
