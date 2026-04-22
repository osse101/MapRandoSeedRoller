package models

type PresetFields struct{
	// Direct Flags
	Version int `path:"version"`
	EscapeMultiplier float64 `path:"skill_assumption_settings.escape_timer_multiplier"`
	StartingItems []StartingItem `path:"item_progression_settings.starting_items"` // starting_items: [{item,count}]
	ObjectiveOptions []ObjectiveOption `path:"objective_settings.objective_options"` // objective_options: [{objective, setting}]
	MinObjectives int `path:"objective_settings.min_objectives"` // TODO: Check number of set objectives is within bounds
	MaxObjectives int `path:"objective_settings.max_objectives"`
	MapLayout string `path:"map_layout"`
	SaveAnimals TriState `path:"save_animals"`
	WallJump string `path:"other_settings.wall_jump"`
	FreeShinesparks bool `path:"other_settings.energy_free_shinesparks"`
	SplitSpeed string `path:"other_settings.speed_booster"`
	IsRace bool `path:"other_settings.race_mode"`

	// Indirect/Resultant Fields
	Name string `path:"name"`
	SkillPreset string `path:"skill_assumption_settings.preset"`
	StartingPreset string `path:"item_progression_settings.starting_items_preset"`
	QualityOfLifePreset string `path:"quality_of_life_settings.preset"`
	ObjectivePreset string `path:"objective_settings.preset"`
	DoorPreset string `path:"doors_settings.preset"`
	AreaAssignmentPreset string `path:"other_settings.area_assignment.preset"`
}

type StartingItem struct{
	Item string `json:"item"`
	Count int `json:"count"`
}

type ObjectiveOption struct{
	Objective string `json:"objective"`
	Setting TriState `json:"setting"`
}

type AliasEntry struct{
	ShortName string // "dray"
	LongName string // "draygon"
}

var ObjectiveAliases = map[string]string{
	"kraid": "Kraid",
	"phan": "Phantoon",
	"dray": "Draygon",
	"ridley": "Ridley",
	"spore": "SporeSpawn",
	"croc": "Crocomire",
	"bot": "Botwoon",
	"gt": "GoldenTorizo",
	"bt": "BombTorizo",
	"bowling": "BowlingStatue",
	"acid": "AcidChozoStatue",
	"pit": "PitRoom",
	"babyk": "BabyKraidRoom",
	"plasma": "PlasmaRoom",
	"metal": "MetalPiratesRoom",
	"m1": "MetroidRoom1",
	"m2": "MetroidsRoom2",
	"m3": "MetroidsRoom3",
	"m4": "MetrpodsRoom4",
}

var ItemAliases = map[string]string{
	"missile": "Missile",
	"etank": "ETank",
	"rtank": "ReserveTank",
	"super": "Super",
	"pb": "PowerBomb",
	"charge": "Charge",
	"ice": "Ice",
	"wave": "Wave",
	"spazer": "Spazer",
	"plasma": "Plasma",
	"xray": "XRayScope",
	"morph": "Morph",
	"bomb": "Bombs",
	"grapple": "Grapple",
	"hjb": "HiJump",
	"speed": "SpeedBooster",
	"spring": "SpringBall",
	"space": "SpaceJump",
	"screw": "ScrewAttack",
	"varia": "Varia",
	"gravity": "Gravity",
}

var FlagAliases = map[string]string{
	"r":"race_mode",
	"d":"version",
	"x":"escape_timer_multiplier",
	"s":"starting_items",
	"o":"objective_options",
	"l":"map_layout",
}