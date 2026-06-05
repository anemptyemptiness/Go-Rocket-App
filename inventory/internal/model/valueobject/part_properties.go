package valueobject

type PartProperties interface {
	Hull() HullProperties
	Engine() EngineProperties
	Shield() ShieldProperties
	Weapon() WeaponProperties
}

type partProperties struct {
	hull   HullProperties
	engine EngineProperties
	shield ShieldProperties
	weapon WeaponProperties
}

func (p *partProperties) Hull() HullProperties {
	return p.hull
}

func (p *partProperties) Engine() EngineProperties {
	return p.engine
}

func (p *partProperties) Shield() ShieldProperties {
	return p.shield
}

func (p *partProperties) Weapon() WeaponProperties {
	return p.weapon
}
