package record

type PartProperties struct {
	Hull   *HullProperties   `json:"hull,omitempty"`
	Engine *EngineProperties `json:"engine,omitempty"`
	Shield *ShieldProperties `json:"shield,omitempty"`
	Weapon *WeaponProperties `json:"weapon,omitempty"`
}

type HullProperties struct {
	Strength int32 `json:"strength"`
}

type EngineProperties struct {
	Class            string `json:"class"`
	RequiredStrength int32  `json:"required_strength"`
}

type ShieldProperties struct {
	ShieldType string `json:"shield_type"`
}

type WeaponProperties struct {
	WeaponType string `json:"weapon_type"`
}
